package cnlib

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignData(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	message := []byte("Hello World")

	signature, err := wallet.SignData(message)
	assert.Nil(t, err)

	signString := hex.EncodeToString(signature)
	expectedSignString := "3045022100c515fc2ed70810f6b1383cfe8e81b9b41b08682511e92d557f1b1719391b521d02200d9d734fd09ce60586ac48b0a7eb587a50958cd9fa548ffa39088fc6ada12eec"

	assert.Equal(t, expectedSignString, signString)
}

func TestSignatureSigningData(t *testing.T) {
	wallet := NewHDWalletFromWords(w, BaseCoinBip84MainNet)
	message := []byte("Hello World")

	str, err := wallet.SignatureSigningData(message)
	assert.Nil(t, err)

	expectedSignString := "3045022100c515fc2ed70810f6b1383cfe8e81b9b41b08682511e92d557f1b1719391b521d02200d9d734fd09ce60586ac48b0a7eb587a50958cd9fa548ffa39088fc6ada12eec"

	assert.Equal(t, expectedSignString, str)
}
