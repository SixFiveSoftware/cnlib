package cnlib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"

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
	txout := wire.NewTxOut(int64(data.Amount), decAddr.ScriptAddress())
	tx.AddTxOut(txout)

	// calculate change
	var transactionChangeMetadata *TransactionChangeMetadata
	if data.shouldAddChangeToTransaction() {
		changeAddr := tb.wallet.ChangeAddressForIndex(data.ChangePath.Index)
		decChange, decChangeErr := btcutil.DecodeAddress(changeAddr.Address, data.basecoin.defaultNetParams())
		if decChangeErr != nil {
			return nil, decChangeErr
		}
		changeOut := wire.NewTxOut(int64(data.ChangeAmount), decChange.ScriptAddress())
		tx.AddTxOut(changeOut)
		metadata := TransactionChangeMetadata{Address: changeAddr.Address, Path: data.ChangePath, VoutIndex: 1}
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
	err := tb.signInputsForTx(tx, data)
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

		sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
		if err != nil {
			return err
		}

		sigScript, err := txscript.SignatureScript(tx, i, sourcePkScript, txscript.SigHashAll, signer.derivedPrivateKey, false)
		if err != nil {
			return err
		}

		txin.SignatureScript = sigScript
	}

	return nil
}
