package cnlib

import "testing"

func TestTransactionBuilderBuildsTxCorrect(t *testing.T) {
	basecoin := NewBaseCoin(49, 0, 0)
	inputPath := NewDerivationPath(49, 0, 0, 1, 53)
	utxo := NewUTXO("1a08dafe993fdc17fdc661988c88f97a9974013291e759b9b5766b8e97c78f87", 1, 2788424, inputPath, nil, true)
	amount := 13584
	feeAmount := 3000
	changeAmount := 2771840
	changePath := NewDerivationPath(49, 0, 0, 1, 56)
	toAddress := "3BgxxADLtnoKu9oytQiiVzYUqvo8weCVy9"

	data := NewTransactionDataFlatFee(toAddress, basecoin, amount, feeAmount, changePath, 539943)
	data.AddUTXO(utxo)
	data.Generate()

	expectedEncodedTx := "01000000000101878fc7978e6b76b5b959e791320174997af9888c9861c6fd17dc3f99feda081a0100000017160014509060a6bedf13087124c0aeafc6e3db4e1e9a08ffffffff02103500000000000017a9146daec6ddb6faaf01f83f515045822a94d0c2331e87804b2a000000000017a914e0bc3e6f5f4080b4f007c6307ba579595e459a06870247304402205b2d50ca2b20fa290323687c3e60bfd4702f9082544afeeb62d849437d04092002204d6dbdef48a992e20700452eff01966d08dcc767b4e7a205c78617d8b5faa1f7012103d447f34dd13359a8fc64ed3977fcecea3f6802f842f9a9f857de07453b715735273d0800"
	expectedTxid := "9ea15d4a60c33a1be64da5805c399663831f7aee13724bfa702db2c3cfafd5bb"
	expectedChangeAddress := "3NBJnvo9U5YbJnr1pALFqQEur1wXWJrjoM"

	wallet := NewHDWalletFromWords(w, basecoin)

	builder := transactionBuilder{wallet: wallet}
	meta, err := builder.buildTxFromData(data.TransactionData)
	if err != nil {
		t.Errorf("Expected to build tx metadata, got error: %v", err)
	}

	if data.TransactionData.PaymentAddress != toAddress {
		t.Errorf("Expected toAddress to be %v, got %v", toAddress, data.TransactionData.PaymentAddress)
	}
	if data.TransactionData.ChangeAmount != changeAmount {
		t.Errorf("Expected change amount to be %v, got %v", changeAmount, data.TransactionData.ChangeAmount)
	}
	if meta.EncodedTx != expectedEncodedTx {
		t.Errorf("Expected encoded tx to be %v,\ngot %v", expectedEncodedTx, meta.EncodedTx)
	}
	if meta.Txid != expectedTxid {
		t.Errorf("Expected txid to be %v,\ngot %v", expectedTxid, meta.Txid)
	}
	if meta.TransactionChangeMetadata.VoutIndex != 1 {
		t.Errorf("Expected change vout index to be %v, got %v", 1, meta.VoutIndex)
	}
	if meta.TransactionChangeMetadata.Path.Index != 56 {
		t.Errorf("Expected change path index to be %v, got %v", 56, meta.TransactionChangeMetadata.Path.Index)
	}
	if meta.TransactionChangeMetadata.Address != expectedChangeAddress {
		t.Errorf("Expected change address to be %v, got %v", expectedChangeAddress, meta.TransactionChangeMetadata.Address)
	}
}
