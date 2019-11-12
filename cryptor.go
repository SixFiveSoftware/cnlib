package cnlib

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"

	"git.coinninja.net/engineering/cryptor"
)

// Encrypt will encrypt a memo using a public key and optionally a private key (pass nil to generate an ephemeral key)
func Encrypt(data []byte, publicKeyBytes []byte, privateKeyBytes []byte) ([]byte, error) {

	// Get or generate public key
	publicKey, err := btcec.ParsePubKey(publicKeyBytes, btcec.S256())
	if err != nil {
		return nil, errors.New("invalid pulic key")
	}

	// Get or generate private key
	var privateKey *btcec.PrivateKey
	if len(privateKeyBytes) == 0 {
		privateKey, _ = btcec.NewPrivateKey(btcec.S256())
	} else {
		privateKey, _ = btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	}
	if err != nil {
		return nil, errors.New("invalid private key")
	}

	return cryptor.Encrypt(data, privateKey, publicKey)

}

// Decrypt will dencrypt a memo using a private key (bytes)
func Decrypt(data []byte, privateKeyBytes []byte) ([]byte, error) {

	privateKey, err := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	if err != nil {
		return nil, errors.New("invalid private key")
	}

	return cryptor.Decrypt(data, privateKey)

}
