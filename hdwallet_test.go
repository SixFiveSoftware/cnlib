package cnlib

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"
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
	w := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
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
	w := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
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
	w := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)

	pk := wallet.SigningPublicKey()

	pkString := hex.EncodeToString(pk)
	expected := "024458596b5c97e716e82015a72c37b5d3fe0c5dc70a4b83d72e7d2eb65920633e"

	if pkString != expected {
		t.Errorf("Expected private key hex to be %v, got %v", expected, pkString)
	}
}
