package cnlib

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	w = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
)

func TestGetFullBIP39WordListString(t *testing.T) {
	wl := GetFullBIP39WordListString()
	words := strings.Split(wl, " ")

	assert.Equal(t, 2048, len(words))
	assert.Equal(t, "abandon", words[0])
	assert.Equal(t, "zoo", words[len(words)-1])
}

func TestNewWordListFromEntropy(t *testing.T) {
	size := 16
	expectedWordLen := 12

	// first set
	bs1 := make([]byte, size)
	n1, err := rand.Read(bs1)
	assert.Nil(t, err)
	assert.Equal(t, size, n1)

	wordString1, err := NewWordListFromEntropy(bs1)
	assert.Nil(t, err)

	words1 := strings.Split(wordString1, " ")

	assert.Equal(t, expectedWordLen, len(words1))

	// second set
	bs2 := make([]byte, size)
	n2, err := rand.Read(bs2)
	assert.Nil(t, err)
	assert.Equal(t, size, n2)

	wordString2, err := NewWordListFromEntropy(bs2)
	assert.Nil(t, err)

	words2 := strings.Split(wordString2, " ")

	assert.Equal(t, expectedWordLen, len(words2))
	assert.NotEqual(t, wordString1, wordString2)
}

func TestSigningKey(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	sk, err := wallet.SigningKey()
	assert.Nil(t, err)

	skString := hex.EncodeToString(sk)
	expected := "8eca986c3aeb26f5ce7717b6c246ebee58ff490ee74c43ce3c4021bb723bd750"

	assert.Equal(t, expected, skString)
}

func TestSigningPublicKey(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	pk, err := wallet.SigningPublicKey()
	assert.Nil(t, err)

	pkString := hex.EncodeToString(pk)
	expected := "024458596b5c97e716e82015a72c37b5d3fe0c5dc70a4b83d72e7d2eb65920633e"

	assert.Equal(t, expected, pkString)
}

func TestCoinNinjaVerificationKeyHexString(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	pkString, err := wallet.CoinNinjaVerificationKeyHexString()
	assert.Nil(t, err)

	expected := "024458596b5c97e716e82015a72c37b5d3fe0c5dc70a4b83d72e7d2eb65920633e"

	assert.Equal(t, expected, pkString)
}

func TestReceiveAddressForIndex_ValidIndex(t *testing.T) {
	i := 0
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	ma, err := wallet.ReceiveAddressForIndex(i)
	assert.Nil(t, err)

	expectedAddress := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"
	expectedPath := NewDerivationPath(BaseCoinBip84MainNet, 0, i)
	expectedKey := "0430d54fd0dd420a6e5f8d3624f5f3482cae350f79d5f0753bf5beef9c2d91af3c04717159ce0828a7f686c2c7510b7aa7d4c685ebc2051642ccbebc7099e2f679"

	assert.Equal(t, expectedAddress, ma.Address)

	// dereference both to compare values, not pointers
	assert.Equal(t, *expectedPath, *ma.DerivationPath)
	assert.Equal(t, expectedKey, ma.UncompressedPublicKey)
}

func TestReceiveAddressForIndex_InvalidIndex(t *testing.T) {
	i := -1
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	ma, err := wallet.ReceiveAddressForIndex(i)
	assert.NotNil(t, err)
	// assert equal error?
	assert.Nil(t, ma)
}

func TestChangeAddressForIndex_ValidIndex(t *testing.T) {
	i := 0
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	ma, err := wallet.ChangeAddressForIndex(i)
	assert.Nil(t, err)

	expectedAddress := "bc1q8c6fshw2dlwun7ekn9qwf37cu2rn755upcp6el"
	expectedPath := NewDerivationPath(BaseCoinBip84MainNet, 1, i)
	expectedKey := ""

	assert.Equal(t, expectedAddress, ma.Address)

	// dereference both to compare values, not pointers
	assert.Equal(t, *expectedPath, *ma.DerivationPath)
	assert.Equal(t, expectedKey, ma.UncompressedPublicKey)
}

func TestChangeAddressForIndex_InvalidIndex(t *testing.T) {
	i := -1
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)

	ma, err := wallet.ChangeAddressForIndex(i)
	assert.EqualError(t, errors.New("index cannot be negative"), err.Error())

	assert.Nil(t, ma)
}

func TestUpdateCoin(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip49MainNet)
	expAddr1 := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	expAddr2 := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"

	ma1, err := wallet.ReceiveAddressForIndex(0)
	assert.Nil(t, err)

	assert.Equal(t, expAddr1, ma1.Address)

	wallet.UpdateCoin(BaseCoinBip84MainNet)

	ma2, err := wallet.ReceiveAddressForIndex(0)
	assert.Nil(t, err)

	assert.Equal(t, expAddr2, ma2.Address)
}

func TestCheckForAddress_AddressExistsInRange(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	expectedAddrAt10 := "bc1qd30z5a5e50jtgx28rvt64483tq65r9pkj623wh"

	ma, err := wallet.CheckForAddress(expectedAddrAt10, 20)

	assert.Nil(t, err)
	assert.Equal(t, expectedAddrAt10, ma.Address)
}

func TestCheckForAddress_AddressDoesNotExist(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	expectedAddrAt30 := "bc1qvy9t2k673tsp6wdwpym3m29sz829nuac9jccc9"

	ma, err := wallet.CheckForAddress(expectedAddrAt30, 20)

	assert.EqualError(t, errors.New("address not found"), err.Error())
	assert.Nil(t, ma)
}

func TestEncyptWithEphemeralKey(t *testing.T) {
	aliceWords := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	bobWords := "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong"
	messageString := "hey dude"
	message := []byte(messageString)
	entropy, err := hex.DecodeString("01010101010101010101010101010101")
	assert.Nil(t, err)
	assert.Equal(t, 16, len(entropy))

	aliceWallet := NewHDWalletFromWords(aliceWords, BaseCoinBip84MainNet)
	bobWallet := NewHDWalletFromWords(bobWords, BaseCoinBip84MainNet)
	bobAddr, err := bobWallet.ReceiveAddressForIndex(0)
	assert.Nil(t, err)
	bobUCPK := bobAddr.UncompressedPublicKey

	assert.Equal(t, 130, len(bobUCPK))

	enc, encErr := aliceWallet.EncryptWithEphemeralKey(entropy, message, bobUCPK)
	assert.Nil(t, encErr)

	bobPath := NewDerivationPath(BaseCoinBip84MainNet, 0, 0)
	dec, err := bobWallet.DecryptWithKeyFromDerivationPath(bobPath, enc)
	assert.Nil(t, err)

	decryptedString := string(dec)
	assert.Equal(t, messageString, decryptedString)
}

func TestEncryptionWithDefaultKeysEndToEnd(t *testing.T) {
	aliceWords := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	bobWords := "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong"
	messageString := "hey dude"
	message := []byte(messageString)

	aliceWallet := NewHDWalletFromWords(aliceWords, BaseCoinBip84MainNet)
	bobWallet := NewHDWalletFromWords(bobWords, BaseCoinBip84MainNet)
	bobCPK, err := bobWallet.CoinNinjaVerificationKeyHexString()
	assert.Nil(t, err)

	enc, err := aliceWallet.EncryptMessage(message, bobCPK)
	assert.Nil(t, err)

	dec, err := bobWallet.DecryptMessage(enc)
	assert.Nil(t, err)

	decryptedString := string(dec)
	assert.Equal(t, messageString, decryptedString)
}

func TestImportPrivateKey(t *testing.T) {
	encodedKey := "L2uv4eejGywPPmsESp3N9Vum9HGX6gBg6RTWJ5oakN9HFTiSKB8i"
	expectedAddress := "1Ad4RSbPrFvo4T5eRMFCoieYf9AuhYdL3h"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	imported, err := wallet.ImportPrivateKey(encodedKey)

	assert.Nil(t, err)
	assert.Equal(t, encodedKey, imported.wif.String())

	addrs := strings.Split(imported.PossibleAddresses, " ")
	found := false
	for _, item := range addrs {
		if item == expectedAddress {
			found = true
		}
	}
	assert.Truef(t, found, "Expected base58check p2pkh address %v, from %v", expectedAddress, imported.PossibleAddresses)
}

func TestImportPrivateKeyP2SHSegwit(t *testing.T) {
	encodedKey := "L3mDwYGp77Zjvqse4YPwbJ7R2M1Zh4vp1RM69JXhbzutVjKwwx9s"
	expectedAddress := "3CFfFMGHUc6rj1JHuTjQYbEmDngnPQF9ev"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	imported, err := wallet.ImportPrivateKey(encodedKey)

	assert.Nil(t, err)
	assert.Equal(t, encodedKey, imported.wif.String())

	addrs := strings.Split(imported.PossibleAddresses, " ")
	found := false
	for _, item := range addrs {
		if item == expectedAddress {
			found = true
		}
	}
	assert.Truef(t, found, "Expected base58check p2sh-p2wkph address %v, from %v", expectedAddress, imported.PossibleAddresses)
}

func TestImportPrivateKeyNativeSegwit(t *testing.T) {
	encodedKey := "L2hgQ3HC3Ru88Jkn5TDwReqeZPhWW4AePebUVFnEQCGJnTPQLgAv"
	expectedAddress := "bc1q2ef8pkkefnamef2sv97dls5ktrq3jlg2ru8ceu"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	imported, err := wallet.ImportPrivateKey(encodedKey)

	assert.Nil(t, err)
	assert.Equal(t, encodedKey, imported.wif.String())

	addrs := strings.Split(imported.PossibleAddresses, " ")
	found := false
	for _, item := range addrs {
		if item == expectedAddress {
			found = true
		}
	}
	assert.Truef(t, found, "Expected segwit address %v, from %v", expectedAddress, imported.PossibleAddresses)
}

func TestDecodeLightningInvoice_WithMemo_WithSats(t *testing.T) {
	invoice := "lnbc2500u1pvjluezpp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpuaztrnwngzn3kdzw5hydlzf03qdgm2hdq27cqv3agm2awhz5se903vruatfhq77w3ls4evs3ch9zw97j25emudupq63nyw24cg27h2rspfj9srp"

	expectedAmount := 250000
	expectedDescription := "1 cup coffee"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	di, err := wallet.DecodeLightningInvoice(invoice)
	assert.Nil(t, err)
	assert.Equal(t, expectedAmount, di.NumSatoshis)
	assert.Equal(t, expectedDescription, di.Description)
	assert.True(t, di.IsExpired)
}

func TestDecodeLightningInvoice_Expired(t *testing.T) {
	invoiceString := "lnbc10u1p0znnwnpp590yg0vlxfd5lsdvn4m3fxzr939nk7ewzsz43fenc3j97exnrurrsdq8w3jhxaqcqzpgxqzfvydc7jvtacl4c7z4jygdhq76whzhve6q6qqdmecahw7nsx0rw0gas3g9ufcej55thffjawjxuga62nphexptuh92fwrumy3z7xng3zwsplmgqah"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	invoice, err := wallet.DecodeLightningInvoice(invoiceString)

	assert.Nil(t, err)
	assert.True(t, invoice.IsExpired)
	assert.Equal(t, int64(1579798271), invoice.ExpiresAt)
}

func TestDecodeLightningInvoice_NoMemo_NoSats(t *testing.T) {
	invoice := "lnbc1p0puj0spp5scnsdawzvccv3mct6x0egfwy28xdu04qt5pmat4z3g52als7kwwqdqqcqzpgxqy9gcqcx7z7vqj4n2yx8wjkyeljjja0xdp2n66f7s4h736mckcpudde0nkjzn03k2qre0lwt3vqt6pj57pkwqkz8elekkp3xhmqjzqkxdgvvcq2gv79f"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	di, err := wallet.DecodeLightningInvoice(invoice)
	assert.Nil(t, err)
	assert.Equal(t, 0, di.NumSatoshis)
	assert.Equal(t, "", di.Description)
}

func TestDecodeLightningInvoice_NoMemo_WithSats(t *testing.T) {
	invoice := "lnbc5u1p0puntgpp583kp34svkjgjzh3r6plgj37ukqzy3m5kjnv8uxrdwd3j5mdtrrmsdqqcqzpgxqy9gcqdqhl9439lxp5fzjgfhj0thnd87s6r5crjqm0qtsdla27hql05x7hpf8s055n5rv9c4dv7488xflg49gj65r9kvgxyhdsqm8d5h634sgpevp8er"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	di, err := wallet.DecodeLightningInvoice(invoice)
	assert.Nil(t, err)
	assert.Equal(t, 500, di.NumSatoshis)
	assert.Equal(t, "", di.Description)
}

func TestDecodeLightningInvoice_WithMemo_NoSats(t *testing.T) {
	invoice := "lnbc1p0punsepp5ae28vtazjqdzqln9whtv3hn55q59eys75rzlv4f7muty8958dwrqdq0fpjhjgreyaskcmqcqzpgxqy9gcq234shpy9k2kflhmdmah3xn6m7s0avk840hzxkfydaurrugxyl78pa80x5x8emncje7ftjsh09q2t7443wdxn07h9gnep3uzdppw5xpgp83xq7q"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	di, err := wallet.DecodeLightningInvoice(invoice)
	assert.Nil(t, err)
	assert.Equal(t, 0, di.NumSatoshis)
	assert.Equal(t, "Hey y'all", di.Description)
}

func TestDecodeLightningInvoice_Malformed(t *testing.T) {
	invoice := "lnbc1p0punsepp5ae28vtazjqdzhtv3hn55q59eys75rzlv4f7muty8958dwrqdq0fpjhjgreyaskcmqcqzpgxqy9gcq234shpy9k2kflhmdmah3xn6m7s0avk840hzxkfydaurrugxyl78pa80x5x8emncje7ftjsh09q2t7443wdxn07h9gnep3uzdppw5xpgp83xq7q"

	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	di, err := wallet.DecodeLightningInvoice(invoice)
	assert.Contains(t, err.Error(), "checksum failed")
	assert.Nil(t, di)
}

func TestExtendedAccountPublicKey_BIP84(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	actualKey, err := wallet.AccountExtendedMasterPublicKey()
	assert.Nil(t, err)
	expectedKey := "zpub6rFR7y4Q2AijBEqTUquhVz398htDFrtymD9xYYfG1m4wAcvPhXNfE3EfH1r1ADqtfSdVCToUG868RvUUkgDKf31mGDtKsAYz2oz2AGutZYs"
	assert.Equal(t, expectedKey, actualKey)
}

func TestExtendedAccountPublicKey_BIP44(t *testing.T) {
	bc := NewBaseCoin(44, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	actualKey, err := wallet.AccountExtendedMasterPublicKey()
	assert.Nil(t, err)
	expectedKey := "xpub6BosfCnifzxcFwrSzQiqu2DBVTshkCXacvNsWGYJVVhhawA7d4R5WSWGFNbi8Aw6ZRc1brxMyWMzG3DSSSSoekkudhUd9yLb6qx39T9nMdj"
	assert.Equal(t, expectedKey, actualKey)
}

func TestExtendedAccountPublicKey_BIP49(t *testing.T) {
	bc := NewBaseCoin(49, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	actualKey, err := wallet.AccountExtendedMasterPublicKey()
	assert.Nil(t, err)
	expectedKey := "ypub6Ww3ibxVfGzLrAH1PNcjyAWenMTbbAosGNB6VvmSEgytSER9azLDWCxoJwW7Ke7icmizBMXrzBx9979FfaHxHcrArf3zbeJJJUZPf663zsP"
	assert.Equal(t, expectedKey, actualKey)
}
