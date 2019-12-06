package cnlib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
)

type transactionBuilder struct {
	wallet *HDWallet
}

type cnSecretsSource struct {
	wallet          *HDWallet
	usableAddresses map[string]*UsableAddress
}

func (s cnSecretsSource) GetKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	script := s.usableAddresses[addr.EncodeAddress()]
	return script.derivedPrivateKey, true, nil
}

func (s cnSecretsSource) GetScript(addr btcutil.Address) ([]byte, error) {
	script := s.usableAddresses[addr.EncodeAddress()]
	pk := script.derivedPrivateKey
	hash := btcutil.Hash160(pk.PubKey().SerializeCompressed())
	scriptSig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(hash).Script()
	if err != nil {
		return nil, err
	}
	addrHash, err := btcutil.NewAddressScriptHash(scriptSig, s.wallet.Basecoin.defaultNetParams())
	if err != nil {
		return nil, err
	}
	return addrHash.ScriptAddress(), nil
}

func (s cnSecretsSource) ChainParams() *chaincfg.Params {
	return s.wallet.Basecoin.defaultNetParams()
}

func (tb transactionBuilder) buildTxFromData(data *TransactionData) (*TransactionMetadata, error) {
	// create transaction with version
	tx := wire.NewMsgTx(wire.TxVersion)

	// populate tx with payment data
	decAddr, decAddrErr := btcutil.DecodeAddress(data.PaymentAddress, data.basecoin.defaultNetParams())
	if decAddrErr != nil {
		return nil, decAddrErr
	}
	destPkScript, err := txscript.PayToAddrScript(decAddr)
	if err != nil {
		return nil, err
	}
	txout := wire.NewTxOut(int64(data.Amount), destPkScript)
	tx.AddTxOut(txout)

	// calculate change
	var transactionChangeMetadata *TransactionChangeMetadata
	if data.shouldAddChangeToTransaction() {
		changeMetaAddr := tb.wallet.ChangeAddressForIndex(data.ChangePath.Index)
		changeAddr := changeMetaAddr.Address
		decChange, decChangeErr := btcutil.DecodeAddress(changeAddr, data.basecoin.defaultNetParams())
		if decChangeErr != nil {
			return nil, decChangeErr
		}
		changePkScript, err := txscript.PayToAddrScript(decChange)
		if err != nil {
			return nil, err
		}
		changeOut := wire.NewTxOut(int64(data.ChangeAmount), changePkScript)
		tx.AddTxOut(changeOut)
		metadata := TransactionChangeMetadata{Address: changeAddr, Path: data.ChangePath, VoutIndex: 1}
		transactionChangeMetadata = &metadata
	}

	// populate utxos as inputs
	for i := 0; i < data.utxoCount(); i++ {
		utxo, utxoErr := data.requiredUTXOAtIndex(i)
		if utxoErr != nil {
			return nil, utxoErr
		}

		// prev tx outpoint
		if utxo.Index < 0 || utxo.Index > math.MaxUint32 {
			return nil, errors.New("previous utxo index out of bounds")
		}
		newHash, newHashErr := chainhash.NewHashFromStr(utxo.Txid)
		if newHashErr != nil {
			return nil, newHashErr
		}
		outpoint := wire.NewOutPoint(newHash, uint32(utxo.Index))

		// build input from previous output
		txIn := wire.NewTxIn(outpoint, nil, nil)

		// set sequence
		txIn.Sequence = data.getSuggestedSequence()

		// add input to tx inputs
		tx.AddTxIn(txIn)
	}

	// sign inputs
	err = tb.signInputsForTx(tx, data)
	if err != nil {
		return nil, err
	}

	// set locktime
	if data.Locktime < 0 || data.Locktime > math.MaxUint32 {
		return nil, errors.New("Locktime out of bounds")
	}
	tx.LockTime = uint32(data.Locktime)

	// encode and return
	txid := tx.TxHash().String()
	var encodedBytes bytes.Buffer
	tx.Serialize(&encodedBytes)
	tm := TransactionMetadata{Txid: txid, EncodedTx: hex.EncodeToString(encodedBytes.Bytes())}
	tm.TransactionChangeMetadata = transactionChangeMetadata
	return &tm, nil
}

func (tb transactionBuilder) signInputsForTx(tx *wire.MsgTx, data *TransactionData) error {
	prevPkScripts := make([][]byte, data.utxoCount())
	inputValues := make([]btcutil.Amount, data.utxoCount())
	secretsSource := cnSecretsSource{wallet: tb.wallet, usableAddresses: make(map[string]*UsableAddress)}

	for i := range tx.TxIn {
		utxo, _ := data.requiredUTXOAtIndex(i)

		var signer *UsableAddress
		if utxo.Path != nil {
			signer = NewUsableAddressWithDerivationPath(tb.wallet, utxo.Path)
		} else if utxo.ImportedPrivateKey != nil {
			signer = NewUsableAddressWithImportedPrivateKey(tb.wallet, utxo.ImportedPrivateKey)
		} else {
			return errors.New("no private key available to sign input")
		}

		var address string
		if utxo.Path != nil {
			address = signer.MetaAddress().Address
		} else if utxo.ImportedPrivateKey != nil && utxo.ImportedPrivateKey.SelectedAddress != "" {
			address = utxo.ImportedPrivateKey.SelectedAddress
		} else {
			return errors.New("no source address available to sign input")
		}

		sourceAddress, err := btcutil.DecodeAddress(address, tb.wallet.Basecoin.defaultNetParams())
		if err != nil {
			return err
		}

		pkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return err
		}

		prevPkScripts[i] = pkScript
		inputValues[i] = btcutil.Amount(utxo.Amount)
		secretsSource.usableAddresses[address] = signer
	}

	scriptsErr := txauthor.AddAllInputScripts(tx, prevPkScripts, inputValues, secretsSource)
	if scriptsErr != nil {
		return scriptsErr
	}

	// verify
	for i := range tx.TxIn {
		flags := txscript.StandardVerifyFlags
		pkScript := prevPkScripts[i]
		inputValue := inputValues[i]
		vm, verErr := txscript.NewEngine(pkScript, tx, i, flags, nil, nil, int64(inputValue))
		if verErr != nil {
			return verErr
		}
		if vmErr := vm.Execute(); vmErr != nil {
			return vmErr
		}
	}

	return nil
}
