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
	if ok, err := data.Generate(); !ok {
		t.Errorf("Expected to generate transaction, got error: %v", err)
	}

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

func TestTransactionBuilder_TwoInputs_BuildsTransaction(t *testing.T) {
	basecoin := NewBaseCoin(49, 0, 0)
	path1 := NewDerivationPath(49, 0, 0, 1, 56)
	path2 := NewDerivationPath(49, 0, 0, 1, 57)
	utxo1 := NewUTXO("24cc9150963a2369d7f413af8b18c3d0243b438ba742d6d083ec8ed492d312f9", 1, 2769977, path1, nil, true)
	utxo2 := NewUTXO("ed611c20fc9088aa5ec1c86de88dd017965358c150c58f71eda721cdb2ac0a48", 1, 314605, path2, nil, true)
	amount := 3000000
	feeAmount := 4000
	changeAmount := 80582
	changePath := NewDerivationPath(49, 0, 0, 1, 58)
	toAddress := "3CkiUcj5vU4TGZJeDcrmYGWH8GYJ5vKcQq"

	data := NewTransactionDataFlatFee(toAddress, basecoin, amount, feeAmount, changePath, 540220)
	data.AddUTXO(utxo1)
	data.AddUTXO(utxo2)
	if ok, err := data.Generate(); !ok {
		t.Errorf("Expected to generate transaction, got error: %v", err)
	}

	expectedEncodedTx := "01000000000102f912d392d48eec83d0d642a78b433b24d0c3188baf13f4d769233a965091cc24010000001716001436386ac950d557ae06bfffc51e7b8fa08474c05ffdffffff480aacb2cd21a7ed718fc550c158539617d08de86dc8c15eaa8890fc201c61ed010000001716001480e1e7dc2f6436a60abec5e9e7f6b62b0b9985c4fdffffff02c0c62d000000000017a914795c7bc23aebac7ddea222bb13c5357b32ed0cd487c63a01000000000017a914a4a2fab6264d22efbfc997f30738ccc6db0f8c058702483045022100be58ec4344d27e7a9a014703f3c2b7ac2c284d7fb933c0fea71a266b3d19c7980220149fb4345f5612f080ee2f4e3c45138d2e054141f567e2b3ce024daa909efbec0121027c3fde52baba263e526ee5acc051f7fd69000eb633b8cf7decd1334db8fb44ee02483045022100c8303887f614e851c9dafe26a952ec1593af9f88b7587f66fa18a6089c59be7b02200fcee28275114efde7ca62f6c24d14552c38f211b4c3790c679f48a9ab1972d4012103cbd9a8066a39e1d05ec26b72116e84b8b852b6784a6359ebb35f5794445245883c3e0800"
	expectedTxid := "f94e7111736dd2a5fd1c5bbcced153f90d17ee1b032f166dda785354f4063651"
	expectedChangeAddress := "3GhXz1NGhwQusEiBYKKhTqQYE6MKt2utDN"

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
	if meta.TransactionChangeMetadata.Path.Index != 58 {
		t.Errorf("Expected change path index to be %v, got %v", 58, meta.TransactionChangeMetadata.Path.Index)
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
	if ok, err := data.Generate(); !ok {
		t.Errorf("Expected to generate transaction, got error: %v", err)
	}

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

func TestTransactionBuilder_BuildP2KH_NoChange(t *testing.T) {
	basecoin := NewBaseCoin(49, 0, 0)
	path := NewDerivationPath(49, 0, 0, 1, 7)
	utxo := NewUTXO("f14914f76ad26e0c1aa5a68c82b021b854c93850fde12f8e3188c14be6dc384e", 1, 33253, path, nil, true)
	amount := 23147
	feeAmount := 10108
	changePath := NewDerivationPath(49, 0, 0, 1, 2)
	toAddress := "1HT6WtD5CAToc8wZdacCgY4XjJR4jV5Q5d"

	data := NewTransactionDataFlatFee(toAddress, basecoin, amount, feeAmount, changePath, 500000)
	data.AddUTXO(utxo)
	success, err := data.Generate()
	if !success {
		t.Errorf("Expected to generate transaction, got error: %v", err)
	}

	expectedEncodedTx := "010000000001014e38dce64bc188318e2fe1fd5038c954b821b0828ca6a51a0c6ed26af71449f10100000017160014b4381165b195b3286079d46eb2dc8058e6f02241ffffffff016b5a0000000000001976a914b4716e71b900b957e49f749c8432b910417788e888ac02483045022100f8a78ff2243c591ffb7af46ed670b173e5e5dd3f19853493f5c3bda85425f8ef02203d152fdc632388da527c4a58b796a8a40d1a9d15176d80dedfef96a38ecc9ae7012103a45ef894ab9e6f2e55683561181be9e69b20207af746d60b95fab33476dc932420a10700"
	expectedTxid := "77cf4bddf3d133fc37a08e18c47607702e0aec095606f364081d22a4680c3e97"

	wallet := NewHDWalletFromWords(w, basecoin)

	builder := transactionBuilder{wallet: wallet}
	meta, err := builder.buildTxFromData(data.TransactionData)

	if err != nil {
		t.Errorf("Expected to build transaction, got error: %v", err)
		return
	}

	if data.TransactionData.PaymentAddress != toAddress {
		t.Errorf("Expected toAddress to be %v, got %v", toAddress, data.TransactionData.PaymentAddress)
	}
	if meta.EncodedTx != expectedEncodedTx {
		t.Errorf("Expected encoded tx to be %v,\ngot %v", expectedEncodedTx, meta.EncodedTx)
	}
	if meta.Txid != expectedTxid {
		t.Errorf("Expected txid to be %v,\ngot %v", expectedTxid, meta.Txid)
	}
	if meta.TransactionChangeMetadata != nil {
		t.Errorf("Expected change metadata to be nil")
	}
}
