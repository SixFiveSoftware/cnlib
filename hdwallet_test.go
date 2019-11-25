package cnlib

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

const (
	w = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
)

func TestGetFullBIP39WordListString(t *testing.T) {
	wl := GetFullBIP39WordListString()
	words := strings.Split(wl, " ")

	if len(words) != 2048 {
		t.Errorf("Expected 2048 words, got %v", len(words))
	}

	if words[0] != "abandon" {
		t.Errorf("Expected first word to be 'abandon', got %v", words[0])
	}

	if words[len(words)-1] != "zoo" {
		t.Errorf("Expected last word to be 'zoo', got %v", words[len(words)-1])
	}
}

func TestNewWordListFromEntropy(t *testing.T) {
	size := 16
	expectedWordLen := 12

	// first set
	bs1 := make([]byte, size)
	n1, err := rand.Read(bs1)
	if err != nil {
		t.Errorf("Expected bytes1 to be created properly. error: %v", err)
		return
	}

	if n1 != size {
		t.Errorf("Expected %v bytes, got %v", size, n1)
	}

	wordString1 := NewWordListFromEntropy(bs1)
	words1 := strings.Split(wordString1, " ")

	if len(words1) != expectedWordLen {
		t.Errorf("Expected word list len to be %v, got %v", expectedWordLen, len(words1))
	}

	// second set
	bs2 := make([]byte, size)
	n2, err := rand.Read(bs2)
	if err != nil {
		t.Errorf("Expected bytes2 to be created properly. error: %v", err)
		return
	}

	if n2 != size {
		t.Errorf("Expected %v bytes, got %v", size, n2)
	}

	wordString2 := NewWordListFromEntropy(bs2)
	words2 := strings.Split(wordString2, " ")

	if len(words2) != expectedWordLen {
		t.Errorf("Expected word list len to be %v, got %v", expectedWordLen, len(words2))
	}

	if wordString1 == wordString2 {
		t.Errorf("Expected two wordStrings to be different, but were the same.")
	}
}

func TestNewHDWalletFromWords(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	if wallet.WalletWords != w {
		t.Errorf("Expected wallet words to equal provided words: %v,\n...but got: %v", w, wallet.WalletWords)
	}

	if wallet.Basecoin.Purpose != bc.Purpose {
		t.Errorf("Expected purpose %v to equal provided purpose %v.", wallet.Basecoin.Purpose, bc.Purpose)
	}

	if wallet.Basecoin.Coin != bc.Coin {
		t.Errorf("Expected coin %v to equal provided coin %v.", wallet.Basecoin.Coin, bc.Coin)
	}

	if wallet.Basecoin.Account != bc.Account {
		t.Errorf("Expected account %v to equal provided account %v.", wallet.Basecoin.Account, bc.Account)
	}
}

func TestSigningKey(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	sk := wallet.SigningKey()

	skString := hex.EncodeToString(sk)
	expected := "8eca986c3aeb26f5ce7717b6c246ebee58ff490ee74c43ce3c4021bb723bd750"

	if skString != expected {
		t.Errorf("Expected private key hex to be %v, got %v", expected, skString)
	}
}

func TestSigningPublicKey(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	pk := wallet.SigningPublicKey()

	pkString := hex.EncodeToString(pk)
	expected := "024458596b5c97e716e82015a72c37b5d3fe0c5dc70a4b83d72e7d2eb65920633e"

	if pkString != expected {
		t.Errorf("Expected private key hex to be %v, got %v", expected, pkString)
	}
}

func TestCoinNinjaVerificationKeyHexString(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	pkString := wallet.CoinNinjaVerificationKeyHexString()
	expected := "024458596b5c97e716e82015a72c37b5d3fe0c5dc70a4b83d72e7d2eb65920633e"

	if pkString != expected {
		t.Errorf("Expected private key hex to be %v, got %v", expected, pkString)
	}
}

func TestReceiveAddressForIndex_ValidIndex(t *testing.T) {
	i := 0
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	ma := wallet.ReceiveAddressForIndex(i)
	expectedAddress := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"
	expectedPath := NewDerivationPath(84, 0, 0, 0, i)
	expectedKey := "0430d54fd0dd420a6e5f8d3624f5f3482cae350f79d5f0753bf5beef9c2d91af3c04717159ce0828a7f686c2c7510b7aa7d4c685ebc2051642ccbebc7099e2f679"

	if ma.Address != expectedAddress {
		t.Errorf("Expected address %v, got %v", expectedAddress, ma.Address)
	}

	// dereference both to compare values, not pointers
	if *ma.DerivationPath != *expectedPath {
		t.Errorf("Expected path %v, got %v", expectedPath, ma.DerivationPath)
	}

	if ma.UncompressedPublicKey != expectedKey {
		t.Errorf("Expected key %v, got %v", expectedKey, ma.UncompressedPublicKey)
	}
}

func TestReceiveAddressForIndex_InvalidIndex(t *testing.T) {
	i := -1
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	ma := wallet.ReceiveAddressForIndex(i)

	if ma != nil {
		t.Errorf("Expected MetaAddress to be nil.")
	}
}

func TestChangeAddressForIndex_ValidIndex(t *testing.T) {
	i := 0
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	ma := wallet.ChangeAddressForIndex(i)
	expectedAddress := "bc1q8c6fshw2dlwun7ekn9qwf37cu2rn755upcp6el"
	expectedPath := NewDerivationPath(84, 0, 0, 1, i)
	expectedKey := ""

	if ma.Address != expectedAddress {
		t.Errorf("Expected address %v, got %v", expectedAddress, ma.Address)
	}

	// dereference both to compare values, not pointers
	if *ma.DerivationPath != *expectedPath {
		t.Errorf("Expected path %v, got %v", expectedPath, ma.DerivationPath)
	}

	if ma.UncompressedPublicKey != expectedKey {
		t.Errorf("Expected key %v, got %v", expectedKey, ma.UncompressedPublicKey)
	}
}

func TestChangeAddressForIndex_InvalidIndex(t *testing.T) {
	i := -1
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	ma := wallet.ChangeAddressForIndex(i)

	if ma != nil {
		t.Errorf("Expected MetaAddress to be nil.")
	}
}

func TestUpdateCoin(t *testing.T) {
	bc := NewBaseCoin(49, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	newCoin := NewBaseCoin(84, 0, 0)
	expAddr1 := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	expAddr2 := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"

	ma1 := wallet.ReceiveAddressForIndex(0)
	if ma1.Address != expAddr1 {
		t.Errorf("Expected address %v, got %v", expAddr1, ma1.Address)
	}

	wallet.UpdateCoin(newCoin)

	ma2 := wallet.ReceiveAddressForIndex(0)
	if ma2.Address != expAddr2 {
		t.Errorf("Expected address %v, got %v", expAddr2, ma2.Address)
	}
}

func TestCheckForAddress_AddressExistsInRange(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	expectedAddrAt10 := "bc1qd30z5a5e50jtgx28rvt64483tq65r9pkj623wh"

	ma, err := wallet.CheckForAddress(expectedAddrAt10, 20)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ma.Address != expectedAddrAt10 {
		t.Errorf("Expected to find %v, got %v", expectedAddrAt10, ma.Address)
	}
}

func TestCheckForAddress_AddressDoesNotExist(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	expectedAddrAt30 := "bc1qvy9t2k673tsp6wdwpym3m29sz829nuac9jccc9"

	ma, err := wallet.CheckForAddress(expectedAddrAt30, 20)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if ma != nil {
		t.Errorf("Expected MetaAddress to be nil, got %v", ma.Address)
	}
}

func TestEncyptWithEphemeralKey(t *testing.T) {
	// Skip for now.
	t.Skip()
	aliceWords := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	bobWords := "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong"
	messageString := "hey dude"
	message := []byte(messageString)
	c := NewBaseCoin(84, 0, 0)
	entropy := []byte("01010101010101010101010101010101")

	aliceWallet := NewHDWalletFromWords(aliceWords, c)
	fmt.Println("alice privkey: ", hex.EncodeToString(aliceWallet.SigningKey()))
	bobWallet := NewHDWalletFromWords(bobWords, c)
	bobUCPK := bobWallet.ReceiveAddressForIndex(0).UncompressedPublicKey

	enc, encErr := aliceWallet.EncryptWithEphemeralKey(message, entropy, bobUCPK)
	fmt.Println("enc: ", hex.EncodeToString(enc))
	if encErr != nil {
		t.Errorf("Should have received an encrypted payload, got %v.", encErr)
	}

	bobPath := NewDerivationPath(84, 0, 0, 0, 0)
	dec, decErr := bobWallet.DecryptWithKeyFromDerivationPath(enc, bobPath)
	fmt.Println("dec: ", hex.EncodeToString(dec))
	if decErr != nil {
		t.Errorf("Should have decrypted an encrypted payload, got %v", decErr)
	}
	decryptedString := hex.EncodeToString(dec)
	if decryptedString != messageString {
		t.Errorf("Expected %v, got %v", messageString, decryptedString)
	}
}

func TestEncryptionWithDefaultKeysEndToEnd(t *testing.T) {
	// Skip for now.
	t.Skip()
	aliceWords := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	bobWords := "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong"
	messageString := "hey dude"
	message := []byte(messageString)
	c := NewBaseCoin(84, 0, 0)

	aliceWallet := NewHDWalletFromWords(aliceWords, c)
	fmt.Println("alice privkey: ", hex.EncodeToString(aliceWallet.SigningKey()))
	bobWallet := NewHDWalletFromWords(bobWords, c)
	bobCPK := bobWallet.CoinNinjaVerificationKeyHexString()

	enc, encErr := aliceWallet.EncryptWithDefaultKey(message, bobCPK)
	fmt.Println("enc: ", hex.EncodeToString(enc))
	if encErr != nil {
		t.Errorf("Should have received an encrypted payload, got %v.", encErr)
	}

	dec, decErr := bobWallet.DecryptWithDefaultKey(enc)
	fmt.Println("dec: ", hex.EncodeToString(dec))
	if decErr != nil {
		t.Errorf("Should have decrypted an encrypted payload, got %v", decErr)
	}
	decryptedString := hex.EncodeToString(dec)
	if decryptedString != messageString {
		t.Errorf("Expected %v, got %v", messageString, decryptedString)
	}
}
