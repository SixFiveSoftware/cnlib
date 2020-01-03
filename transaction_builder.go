package cnlib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
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
		changeMetaAddr, err := tb.wallet.ChangeAddressForIndex(data.ChangePath.Index)
		if err != nil {
			return nil, err
		}

		changeAddr := changeMetaAddr.Address
		decChange, err := btcutil.DecodeAddress(changeAddr, data.basecoin.defaultNetParams())
		if err != nil {
			return nil, err
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
	for i := 0; i < data.UtxoCount(); i++ {
		utxo, utxoErr := data.RequiredUTXOAtIndex(i)
		if utxoErr != nil {
			return nil, utxoErr
		}

		// prev tx outpoint
		if utxo.Index < 0 || utxo.Index > int(math.MaxUint32) {
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

	// set locktime
	if data.Locktime < 0 || data.Locktime > int(math.MaxUint32) {
		return nil, errors.New("Locktime out of bounds")
	}
	tx.LockTime = uint32(data.Locktime)

	// sign inputs
	err = tb.signInputsForTx(tx, data)
	if err != nil {
		return nil, err
	}

	// encode and return
	txid := tx.TxHash().String()
	var encodedBytes bytes.Buffer
	err = tx.Serialize(&encodedBytes)
	if err != nil {
		return nil, err
	}

	tm := TransactionMetadata{Txid: txid, EncodedTx: hex.EncodeToString(encodedBytes.Bytes())}
	tm.TransactionChangeMetadata = transactionChangeMetadata
	return &tm, nil
}

func (tb transactionBuilder) signInputsForTx(tx *wire.MsgTx, data *TransactionData) error {
	prevPkScripts := make([][]byte, data.UtxoCount())
	inputValues := make([]btcutil.Amount, data.UtxoCount())
	secretsSource := cnSecretsSource{wallet: tb.wallet, usableAddresses: make(map[string]*UsableAddress)}

	for i := range tx.TxIn {
		utxo, _ := data.RequiredUTXOAtIndex(i)

		var address string
		if utxo.Path != nil {
			signer, err := NewUsableAddressWithDerivationPath(tb.wallet, utxo.Path)
			if err != nil {
				return err
			}
			meta, err := signer.MetaAddress()
			if err != nil {
				return err
			}
			address = meta.Address
			secretsSource.usableAddresses[address] = signer
		} else if utxo.ImportedPrivateKey != nil && utxo.ImportedPrivateKey.SelectedAddress != "" {
			signer := NewUsableAddressWithImportedPrivateKey(tb.wallet, utxo.ImportedPrivateKey)
			address = utxo.ImportedPrivateKey.SelectedAddress
			secretsSource.usableAddresses[address] = signer
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
	}

	scriptsErr := txauthor.AddAllInputScripts(tx, prevPkScripts, inputValues, secretsSource)
	if scriptsErr != nil {
		return scriptsErr
	}

	// verify
	err := validateMsgTx(tx, prevPkScripts, inputValues)
	if err != nil {
		return err
	}

	// success
	return nil
}

func validateMsgTx(tx *wire.MsgTx, prevScripts [][]byte, inputValues []btcutil.Amount) error {
	hashCache := txscript.NewTxSigHashes(tx)
	flags := txscript.StandardVerifyFlags
	for i, prevScript := range prevScripts {
		vm, err := txscript.NewEngine(prevScript, tx, i, flags, nil, hashCache, int64(inputValues[i]))
		if err != nil {
			return fmt.Errorf("cannot create script engine: %s", err)
		}
		err = vm.Execute()
		if err != nil {
			return fmt.Errorf("cannot validate transaction: %s", err)
		}
	}
	return nil
}
