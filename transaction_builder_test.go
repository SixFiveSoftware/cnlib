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

	expectedEncodedTx := "01000000000101878fc7978e6b76b5b959e791320174997af9888c9861c6fd17dc3f99feda081a0100000017160014509060a6bedf13087124c0aeafc6e3db4e1e9a08fdffffff02103500000000000017a9146daec6ddb6faaf01f83f515045822a94d0c2331e87804b2a000000000017a914e0bc3e6f5f4080b4f007c6307ba579595e459a0687024730440220031851e7fc75043bfa4bb7234478408fd024a50088fee8e16953d347bcfc37ae022050604330a862f1e6d3d2941e0bc5911d3b2f55c8e396e4d0d8c43acbf7e66f16012103d447f34dd13359a8fc64ed3977fcecea3f6802f842f9a9f857de07453b715735273d0800"
	expectedTxid := "20d9d7eae4283573e042de272c0fc6af7df5a1100c4871127fa07c9022da1945"
	expectedChangeAddress := "3NBJnvo9U5YbJnr1pALFqQEur1wXWJrjoM"

	wallet := NewHDWalletFromWords(w, basecoin)

	builder := transactionBuilder{wallet: wallet}
	meta, err := builder.buildTxFromData(data.TransactionData)
	if err != nil {
		t.Errorf("Expected to build tx metadata, got error: %v", err)
		return
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

func TestTransactionBuilder_BuildsNativeSegwitTransaction(t *testing.T) {
	basecoin := NewBaseCoin(84, 0, 0)
	path := NewDerivationPath(84, 0, 0, 0, 1)
	utxo := NewUTXO("a89a9bed1f2daca01a0dca58f7fd0f2f0bf114d762b38e65845c5d1489339a69", 0, 96537, path, nil, true)
	amount := 9755
	feeAmount := 846
	changeAmount := 85936
	changePath := NewDerivationPath(84, 0, 0, 1, 1)
	toAddress := "bc1qjv79zewlvyyyd5y0qfk3svexzrqnammllj7mw6"

	data := NewTransactionDataFlatFee(toAddress, basecoin, amount, feeAmount, changePath, 590582)
	data.AddUTXO(utxo)
	data.Generate()

	expectedEncodedTx := "01000000000101699a3389145d5c84658eb362d714f10b2f0ffdf758ca0d1aa0ac2d1fed9b9aa80000000000fdffffff021b26000000000000160014933c5165df610846d08f026d18332610c13eef7fb04f0100000000001600144227d834f1aae95273f0c87495f4ff0cb366545202483045022100b232240638739a01414442f38f5e2747c891746597edaffbb0120b89120d12fd02201f5de6f8b938492c28459d07f5824fdddd0b869e522680429ca7b08515cd6eaf012103e775fd51f0dfb8cd865d9ff1cca2a158cf651fe997fdc9fee9c1d3b5e995ea77f6020900"
	expectedTxid := "fe7f9a6de3203eb300cc66159e762251d675b5555dbd215c3574e75a762ca402"
	expectedChangeAddress := "bc1qggnasd834t54yulsep6fta8lpjekv4zj6gv5rf"

	wallet := NewHDWalletFromWords(w, basecoin)

	builder := transactionBuilder{wallet: wallet}
	meta, err := builder.buildTxFromData(data.TransactionData)
	if err != nil {
		t.Errorf("Expected to build tx metadata, got error: %v", err)
		return
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
	if meta.TransactionChangeMetadata.Path.Index != 1 {
		t.Errorf("Expected change path index to be %v, got %v", 56, meta.TransactionChangeMetadata.Path.Index)
	}
	if meta.TransactionChangeMetadata.Address != expectedChangeAddress {
		t.Errorf("Expected change address to be %v, got %v", expectedChangeAddress, meta.TransactionChangeMetadata.Address)
	}
}
