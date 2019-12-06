package cnlib

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
)

func TestCryptorEphemeralKey(t *testing.T) {

	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	assert.Nil(t, err)
	privateKeyBytes := privateKey.Serialize()

	publicKeyBytes := privateKey.PubKey().SerializeUncompressed()

	testPayload := []byte("This is my test payload! 🚀")

	encryptedData, err := Encrypt(testPayload, publicKeyBytes, nil)
	assert.Nil(t, err)

	unencryptedPayload, err := Decrypt(encryptedData, privateKeyBytes)
	assert.Nil(t, err)

	assert.Equal(t, testPayload, unencryptedPayload)

}

func TestCryptorPrivateKey(t *testing.T) {

	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	assert.Nil(t, err)
	privateKeyBytes := privateKey.Serialize()

	publicKeyBytes := privateKey.PubKey().SerializeUncompressed()

	ephemeralKey, err := btcec.NewPrivateKey(btcec.S256())
	assert.Nil(t, err)
	ephemeralKeyBytes := ephemeralKey.Serialize()

	testPayload := []byte("This is my test payload! 🚀")

	encryptedData, err := Encrypt(testPayload, publicKeyBytes, ephemeralKeyBytes)
	assert.Nil(t, err)

	unencryptedPayload, err := Decrypt(encryptedData, privateKeyBytes)
	assert.Nil(t, err)

	assert.Equal(t, testPayload, unencryptedPayload)

}
