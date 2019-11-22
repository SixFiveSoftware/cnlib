package cnlib

import (
	"encoding/hex"
	"testing"
)

func TestSignData(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	message := []byte("Hello World")

	signature := wallet.SignData(message)
	signString := hex.EncodeToString(signature)
	expectedSignString := "3045022100c515fc2ed70810f6b1383cfe8e81b9b41b08682511e92d557f1b1719391b521d02200d9d734fd09ce60586ac48b0a7eb587a50958cd9fa548ffa39088fc6ada12eec"

	if signString != expectedSignString {
		t.Errorf("Expected signature %v, got %v", expectedSignString, signString)
	}
}

func TestSignatureSigningData(t *testing.T) {
	bc := NewBaseCoin(84, 0, 0)
	wallet := NewHDWalletFromWords(w, bc)
	message := []byte("Hello World")

	str := wallet.SignatureSigningData(message)

	expectedSignString := "3045022100c515fc2ed70810f6b1383cfe8e81b9b41b08682511e92d557f1b1719391b521d02200d9d734fd09ce60586ac48b0a7eb587a50958cd9fa548ffa39088fc6ada12eec"

	if str != expectedSignString {
		t.Errorf("Expected signature %v, got %v", expectedSignString, str)
	}
}
