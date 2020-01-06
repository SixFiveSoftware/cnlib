package cnlib

import "testing"
import "github.com/stretchr/testify/assert"

func TestTransactionBuilderBuildsTxCorrect(t *testing.T) {
	inputPath := NewDerivationPath(BaseCoinBip49MainNet, 1, 53)
	utxo := NewUTXO("1a08dafe993fdc17fdc661988c88f97a9974013291e759b9b5766b8e97c78f87", 1, 2788424, inputPath, nil, true)
	amount := 13584
	feeAmount := 3000
	changeAmount := 2771840
	changePath := NewDerivationPath(BaseCoinBip49MainNet, 1, 56)
	toAddress := "3BgxxADLtnoKu9oytQiiVzYUqvo8weCVy9"

	data := NewTransactionDataFlatFee(toAddress, BaseCoinBip49MainNet, amount, feeAmount, changePath, 539943)
	data.AddUTXO(utxo)
	err := data.Generate()

	assert.Nil(t, err)

	expectedEncodedTx := "01000000000101878fc7978e6b76b5b959e791320174997af9888c9861c6fd17dc3f99feda081a0100000017160014509060a6bedf13087124c0aeafc6e3db4e1e9a08fdffffff02103500000000000017a9146daec6ddb6faaf01f83f515045822a94d0c2331e87804b2a000000000017a914e0bc3e6f5f4080b4f007c6307ba579595e459a06870247304402205a9d97a269cefe296a746dc07e898d19889567e910339f31e12268703079a45a0220537145228842a020a16894006c7e50ae5109672ea13135a02b354f66838f9676012103d447f34dd13359a8fc64ed3977fcecea3f6802f842f9a9f857de07453b715735273d0800"
	expectedTxid := "20d9d7eae4283573e042de272c0fc6af7df5a1100c4871127fa07c9022da1945"
	expectedChangeAddress := "3NBJnvo9U5YbJnr1pALFqQEur1wXWJrjoM"

	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)

	meta, err := wallet.BuildTransactionMetadata(data.TransactionData)

	assert.Nil(t, err)

	assert.Equal(t, toAddress, data.TransactionData.PaymentAddress)
	assert.Equal(t, changeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedEncodedTx, meta.EncodedTx)
	assert.Equal(t, expectedTxid, meta.Txid)
	assert.Equal(t, 1, meta.TransactionChangeMetadata.VoutIndex)
	assert.Equal(t, 56, meta.TransactionChangeMetadata.Path.Index)
	assert.Equal(t, expectedChangeAddress, meta.TransactionChangeMetadata.Address)
}

func TestTransactionBuilder_TwoInputs_BuildsTransaction(t *testing.T) {
	path1 := NewDerivationPath(BaseCoinBip49MainNet, 1, 56)
	path2 := NewDerivationPath(BaseCoinBip49MainNet, 1, 57)
	utxo1 := NewUTXO("24cc9150963a2369d7f413af8b18c3d0243b438ba742d6d083ec8ed492d312f9", 1, 2769977, path1, nil, true)
	utxo2 := NewUTXO("ed611c20fc9088aa5ec1c86de88dd017965358c150c58f71eda721cdb2ac0a48", 1, 314605, path2, nil, true)
	amount := 3000000
	feeAmount := 4000
	changeAmount := 80582
	changePath := NewDerivationPath(BaseCoinBip49MainNet, 1, 58)
	toAddress := "3CkiUcj5vU4TGZJeDcrmYGWH8GYJ5vKcQq"

	data := NewTransactionDataFlatFee(toAddress, BaseCoinBip49MainNet, amount, feeAmount, changePath, 540220)
	data.AddUTXO(utxo1)
	data.AddUTXO(utxo2)
	err := data.Generate()

	assert.Nil(t, err)

	expectedEncodedTx := "01000000000102f912d392d48eec83d0d642a78b433b24d0c3188baf13f4d769233a965091cc24010000001716001436386ac950d557ae06bfffc51e7b8fa08474c05ffdffffff480aacb2cd21a7ed718fc550c158539617d08de86dc8c15eaa8890fc201c61ed010000001716001480e1e7dc2f6436a60abec5e9e7f6b62b0b9985c4fdffffff02c0c62d000000000017a914795c7bc23aebac7ddea222bb13c5357b32ed0cd487c63a01000000000017a914a4a2fab6264d22efbfc997f30738ccc6db0f8c05870247304402202a1dfa92a9dba16fa476c738197316009665f1b705e5626b2729b136bb64aaa102203041d91270d91124cb9341c6d1bfb2c7aa3372ef85f412fa00b8bf4fa7091f2b0121027c3fde52baba263e526ee5acc051f7fd69000eb633b8cf7decd1334db8fb44ee02483045022100a3843ddb39dd088e8d9657eaed5454a27737112c821eb6f674414e02f295d39402206de16b7c5b1ff054d102451a9242b10fccf81828003c377046bd11fa6c025179012103cbd9a8066a39e1d05ec26b72116e84b8b852b6784a6359ebb35f5794445245883c3e0800"
	expectedTxid := "f94e7111736dd2a5fd1c5bbcced153f90d17ee1b032f166dda785354f4063651"
	expectedChangeAddress := "3GhXz1NGhwQusEiBYKKhTqQYE6MKt2utDN"

	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)

	meta, err := wallet.BuildTransactionMetadata(data.TransactionData)

	assert.Nil(t, err)

	assert.Equal(t, toAddress, data.TransactionData.PaymentAddress)
	assert.Equal(t, changeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedEncodedTx, meta.EncodedTx)
	assert.Equal(t, expectedTxid, meta.Txid)
	assert.Equal(t, 1, meta.TransactionChangeMetadata.VoutIndex)
	assert.Equal(t, 58, meta.TransactionChangeMetadata.Path.Index)
	assert.Equal(t, expectedChangeAddress, meta.TransactionChangeMetadata.Address)
}

func TestTransactionBuilder_BuildsNativeSegwitTransaction(t *testing.T) {
	path := NewDerivationPath(BaseCoinBip84MainNet, 0, 1)
	utxo := NewUTXO("a89a9bed1f2daca01a0dca58f7fd0f2f0bf114d762b38e65845c5d1489339a69", 0, 96537, path, nil, true)
	amount := 9755
	feeAmount := 846
	changeAmount := 85936
	changePath := NewDerivationPath(BaseCoinBip49MainNet, 1, 1)
	toAddress := "bc1qjv79zewlvyyyd5y0qfk3svexzrqnammllj7mw6"

	data := NewTransactionDataFlatFee(toAddress, BaseCoinBip49MainNet, amount, feeAmount, changePath, 590582)
	data.AddUTXO(utxo)
	err := data.Generate()

	assert.Nil(t, err)

	expectedEncodedTx := "01000000000101699a3389145d5c84658eb362d714f10b2f0ffdf758ca0d1aa0ac2d1fed9b9aa80000000000fdffffff021b26000000000000160014933c5165df610846d08f026d18332610c13eef7fb04f0100000000001600144227d834f1aae95273f0c87495f4ff0cb366545202473044022024b8f49fddcc119fc30990d6c970d8a1e0fa56d951d31591bed76c0867dbd11d0220755bb57af82993facbf413e523a8fa6fbccf8055ec95d1764da5e98b54e16bf2012103e775fd51f0dfb8cd865d9ff1cca2a158cf651fe997fdc9fee9c1d3b5e995ea77f6020900"
	expectedTxid := "fe7f9a6de3203eb300cc66159e762251d675b5555dbd215c3574e75a762ca402"
	expectedChangeAddress := "bc1qggnasd834t54yulsep6fta8lpjekv4zj6gv5rf"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	meta, err := wallet.BuildTransactionMetadata(data.TransactionData)

	assert.Nil(t, err)

	assert.Equal(t, toAddress, data.TransactionData.PaymentAddress)
	assert.Equal(t, changeAmount, data.TransactionData.ChangeAmount)
	assert.Equal(t, expectedEncodedTx, meta.EncodedTx)
	assert.Equal(t, expectedTxid, meta.Txid)
	assert.Equal(t, 1, meta.TransactionChangeMetadata.VoutIndex)
	assert.Equal(t, 1, meta.TransactionChangeMetadata.Path.Index)
	assert.Equal(t, expectedChangeAddress, meta.TransactionChangeMetadata.Address)
}

func TestTransactionBuilder_BuildP2KH_NoChange(t *testing.T) {
	path := NewDerivationPath(BaseCoinBip49MainNet, 1, 7)
	utxo := NewUTXO("f14914f76ad26e0c1aa5a68c82b021b854c93850fde12f8e3188c14be6dc384e", 1, 33255, path, nil, true)
	amount := 23147
	feeAmount := 10108
	changePath := NewDerivationPath(BaseCoinBip49MainNet, 1, 2)
	toAddress := "1HT6WtD5CAToc8wZdacCgY4XjJR4jV5Q5d"

	data := NewTransactionDataFlatFee(toAddress, BaseCoinBip49MainNet, amount, feeAmount, changePath, 500000)
	data.AddUTXO(utxo)
	err := data.Generate()

	assert.Nil(t, err)

	expectedEncodedTx := "010000000001014e38dce64bc188318e2fe1fd5038c954b821b0828ca6a51a0c6ed26af71449f10100000017160014b4381165b195b3286079d46eb2dc8058e6f02241fdffffff016b5a0000000000001976a914b4716e71b900b957e49f749c8432b910417788e888ac0247304402204147d25961e7ea6f88df58878aa38167fe6f8ae04c3625485dc594ff716f18a002200c08aabefae62d59568155cfb7ca8df1a4d54c01e5abd767d59e7b982663db23012103a45ef894ab9e6f2e55683561181be9e69b20207af746d60b95fab33476dc932420a10700"
	expectedTxid := "86a9dc5bef7933df26d2b081376084e456a5bd3c2f2df28e758ff062b05a8c17"

	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)

	meta, err := wallet.BuildTransactionMetadata(data.TransactionData)

	assert.Nil(t, err)

	assert.Equal(t, toAddress, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedEncodedTx, meta.EncodedTx)
	assert.Equal(t, expectedTxid, meta.Txid)
	assert.Nil(t, meta.TransactionChangeMetadata)
}

func TestTransationBuilder_BuildSingleUTXO(t *testing.T) {
	path := NewDerivationPath(BaseCoinBip49MainNet, 0, 0)
	utxo := NewUTXO("3480e31ea00efeb570472983ff914694f62804e768a6c6b4d1b6cd70a1cd3efa", 1, 449893, path, nil, true)
	amount := 218384
	feeAmount := 668
	changeAmount := 230841
	changePath := NewDerivationPath(BaseCoinBip49MainNet, 1, 0)
	toAddress := "3ERQiyXSeUYmxxqKyg8XwqGo4W7utgDrTR"

	data := NewTransactionDataFlatFee(toAddress, BaseCoinBip49MainNet, amount, feeAmount, changePath, 500000)
	data.AddUTXO(utxo)
	err := data.Generate()

	assert.Nil(t, err)

	expectedEncodedTx := "01000000000101fa3ecda170cdb6d1b4c6a668e70428f6944691ff83294770b5fe0ea01ee380340100000017160014f990679acafe25c27615373b40bf22446d24ff44fdffffff02105503000000000017a9148ba60342bf59f73327fecab2bef17c1612888c3587b98503000000000017a9141cc1e09a63d1ae795a7130e099b28a0b1d8e4fae8702473044022026f508a317df64f935c43f135280f9f0e95617c22d0f80df77c333656d9303a802206a1c16bd7957e49ddac990f6151065cab326e55d011418e24333d2a979f963d60121039b3b694b8fc5b5e07fb069c783cac754f5d38c3e08bed1960e31fdb1dda35c2420a10700"
	expectedTxid := "221ced4e8784290dea336afa1b0a06fa868812e51abbdca3126ce8d99335a6e2"
	expectedChangeAddress := "34K56kSjgUCUSD8GTtuF7c9Zzwokbs6uZ7"

	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)
	meta, err := wallet.BuildTransactionMetadata(data.TransactionData)

	assert.Nil(t, err)
	assert.Equal(t, expectedEncodedTx, meta.EncodedTx)
	assert.Equal(t, expectedTxid, meta.Txid)
	assert.Equal(t, expectedChangeAddress, meta.TransactionChangeMetadata.Address)
	assert.Equal(t, 1, meta.TransactionChangeMetadata.VoutIndex)
	assert.Equal(t, 0, meta.TransactionChangeMetadata.Path.Index)
	assert.Equal(t, changeAmount, data.TransactionData.ChangeAmount)
}

func TestTransactionBuilder_TestNet(t *testing.T) {
	path := NewDerivationPath(BaseCoinBip49TestNet, 0, 0)
	utxo := NewUTXO("1cfd000efbe248c48b499b0a5d76ea7687ee76cad8481f71277ee283df32af26", 0, 1250000000, path, nil, true)
	amount := 9523810
	feeAmount := 830
	changeAmount := 1240475360
	changePath := NewDerivationPath(BaseCoinBip49MainNet, 1, 0)
	toAddress := "2N8o4Mu5PRAR27TC2eai62CRXarTbQmjyCx"

	data := NewTransactionDataFlatFee(toAddress, BaseCoinBip49TestNet, amount, feeAmount, changePath, 644)
	data.AddUTXO(utxo)
	err := data.Generate()

	assert.Nil(t, err)

	expectedEncodedTx := "0100000000010126af32df83e27e27711f48d8ca76ee8776ea765d0a9b498bc448e2fb0e00fd1c000000001716001438971f73930f6c141d977ac4fd4a727c854935b3fdffffff02625291000000000017a914aa8f293a04a7df8794b743e14ffb96c2a30a1b2787e026f0490000000017a914251dd11457a259c3ba47e5cca3717fe4214e02988702483045022100f24650e94fd022459920770af43f7b630654a85caca68fa73060a7c2422840fc022079267209fb416538e3d471d108f95c90e71e23d7628448f8a3e8c036e93849a1012103a1af804ac108a8a51782198c2d034b28bf90c8803f5a53f76276fa69a4eae77f84020000"
	expectedTxid := "5eb44c7faaa9c17c886588a1e20461d60fbfe1e504e7bac5af3469fdd9039837"
	expectedChangeAddress := "2MvdUi5o3f2tnEFh9yGvta6FzptTZtkPJC8"

	wallet := NewHDWalletFromWords(w, BaseCoinBip49TestNet)
	meta, err := wallet.BuildTransactionMetadata(data.TransactionData)

	assert.Nil(t, err)
	assert.Equal(t, toAddress, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedEncodedTx, meta.EncodedTx)
	assert.Equal(t, expectedTxid, meta.Txid)
	assert.Equal(t, expectedChangeAddress, meta.TransactionChangeMetadata.Address)
	assert.Equal(t, 1, meta.TransactionChangeMetadata.VoutIndex)
	assert.Equal(t, 0, meta.TransactionChangeMetadata.Path.Index)
	assert.Equal(t, changeAmount, data.TransactionData.ChangeAmount)
}

func TestTransactionBuilder_SendToNativeSegwit_BuildsProperly(t *testing.T) {
	path := NewDerivationPath(BaseCoinBip49MainNet, 0, 80)
	utxo := NewUTXO("94b5bcfbd52a405b291d906e636c8e133407e68a75b0a1ccc492e131ff5d8f90", 0, 10261, path, nil, true)
	amount := 5000
	feeAmount := 1000
	changeAmount := 4261
	changePath := NewDerivationPath(BaseCoinBip49MainNet, 1, 102)
	toAddress := "bc1ql2sdag2nm9csz4wmlj735jxw88ym3yukyzmrpj"

	data := NewTransactionDataFlatFee(toAddress, BaseCoinBip49MainNet, amount, feeAmount, changePath, 500000)
	data.AddUTXO(utxo)
	err := data.Generate()

	assert.Nil(t, err)

	expectedEncodedTx := "01000000000101908f5dff31e192c4cca1b0758ae60734138e6c636e901d295b402ad5fbbcb594000000001716001442288ee31111f7187e8cfe8c82917c4734da4c2efdffffff028813000000000000160014faa0dea153d9710155dbfcbd1a48ce39c9b89396a51000000000000017a914aa71651e8f7c618a4576873254ec80c4dfaa068b8702483045022100b3c3d02e7f455503447e70138bcf2f3e928af0d7b9640631e086a56d43740199022018906455f9f7314109e73489bb12c169b3a59302c8456b1b154e894466039f8d01210270d4003d27b5340df1895ef3a5aee2ae2fe3ed7383c01ba623723e702b6c83c120a10700"
	expectedTxid := "1f1ffca0eda219b09116743d2c9b9dcf8eefd10d240bdc4e66678d72a6e4614d"
	expectedChangeAddress := "3HEEdyeVwoGZf86jq8ovUhw9FiXkwCdY79"

	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)
	meta, err := wallet.BuildTransactionMetadata(data.TransactionData)

	assert.Nil(t, err)
	assert.Equal(t, toAddress, data.TransactionData.PaymentAddress)
	assert.Equal(t, expectedEncodedTx, meta.EncodedTx)
	assert.Equal(t, expectedTxid, meta.Txid)
	assert.Equal(t, expectedChangeAddress, meta.TransactionChangeMetadata.Address)
	assert.Equal(t, 1, meta.TransactionChangeMetadata.VoutIndex)
	assert.Equal(t, 102, meta.TransactionChangeMetadata.Path.Index)
	assert.Equal(t, changeAmount, data.TransactionData.ChangeAmount)
}
