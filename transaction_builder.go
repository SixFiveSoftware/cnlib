package cnlib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type transactionBuilder struct {
	wallet *HDWallet
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

	// populate utxos
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
	for i, txin := range tx.TxIn {
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

		if _, ok := sourceAddress.(*btcutil.AddressWitnessPubKeyHash); ok {
			err := spendSegwitInput(tx, txin, i, int64(utxo.Amount), sourceAddress, signer.derivedPrivateKey)
			if err != nil {
				return err
			}
		}
		if _, ok := sourceAddress.(*btcutil.AddressPubKeyHash); ok {
			err := spendKeyHashInput(tx, txin, i, int64(utxo.Amount), sourceAddress, signer.derivedPrivateKey)
			if err != nil {
				return err
			}
		}
		if _, ok := sourceAddress.(*btcutil.AddressScriptHash); ok {
			err := spendSegwitInput(tx, txin, i, int64(utxo.Amount), sourceAddress, signer.derivedPrivateKey)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func spendSegwitInput(tx *wire.MsgTx, txIn *wire.TxIn, idx int, amount int64, address btcutil.Address, privKey *btcec.PrivateKey) error {
	witnessProgram, err := txscript.PayToAddrScript(address)
	if err != nil {
		return err
	}

	if _, ok := address.(*btcutil.AddressScriptHash); ok {
		bldr := txscript.NewScriptBuilder()
		bldr.AddData(witnessProgram)
		sigScript, err := bldr.Script()
		if err != nil {
			return err
		}
		txIn.SignatureScript = sigScript
	}

	txSigHashes := txscript.NewTxSigHashes(tx)
	witnessScript, err := txscript.WitnessSignature(tx, txSigHashes, idx, amount, witnessProgram, txscript.SigHashAll, privKey, false)
	if err != nil {
		return err
	}

	txIn.Witness = witnessScript
	return nil
}

func spendKeyHashInput(tx *wire.MsgTx, txIn *wire.TxIn, idx int, amount int64, address btcutil.Address, privKey *btcec.PrivateKey) error {
	sourcePkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return err
	}

	sigScript, err := txscript.SignatureScript(tx, idx, sourcePkScript, txscript.SigHashAll, privKey, false)
	if err != nil {
		return err
	}

	txIn.SignatureScript = sigScript
	return nil
}
