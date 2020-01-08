package cnlib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"

	"github.com/btcsuite/btcd/btcec"
)

const minPayloadSize = 131

// decrypt data using public/private keypair
func decrypt(data []byte, privateKey *btcec.PrivateKey) ([]byte, error) {

	if len(data) < minPayloadSize {
		return nil, errors.New("insufficient data")
	}

	version := data[:1]
	options := data[1:2]
	iv := data[2:18]
	cipherText := data[18:(len(data) - 32 - 65)]
	hmacVal := data[len(data)-32-65 : len(data)-65]
	publicKeyUncomp := data[len(data)-65:]

	if options[0] != byte(0) {
		return nil, errors.New("invalid payload option")
	}

	msg := make([]byte, 0)
	msg = append(msg, version...)
	msg = append(msg, options...)
	msg = append(msg, iv...)
	msg = append(msg, cipherText...)

	publicKey, err := btcec.ParsePubKey(publicKeyUncomp, btcec.S256())
	if err != nil {
		return nil, err
	}

	secret := generateSharedSecretRFC4753(privateKey, publicKey)
	keyData := sha512.Sum512(secret)

	encKey := keyData[:32]
	hmacKey := keyData[32:]

	testHmac := hmac.New(sha256.New, hmacKey)
	testHmac.Write(msg)
	testHmacVal := testHmac.Sum(nil)

	// its important to use hmac.Equal to not leak time
	// information. See https://github.com/RNCryptor/RNCryptor-Spec
	if verified := hmac.Equal(testHmacVal, hmacVal); !verified {
		return nil, errors.New("invalid hmac")
	}

	cipherBlock, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(cipherText))
	copy(decrypted, cipherText)
	decrypter := cipher.NewCBCDecrypter(cipherBlock, iv)
	decrypter.CryptBlocks(decrypted, decrypted)

	// un-padd decrypted data
	length := len(decrypted)
	unpadding := int(decrypted[length-1])

	return decrypted[:(length - unpadding)], nil
}

// encrypt Data using public/private keypair
func encrypt(data []byte, privateKey *btcec.PrivateKey, publicKey *btcec.PublicKey) ([]byte, error) {

	secret := generateSharedSecretRFC4753(privateKey, publicKey)
	keyData := sha512.Sum512(secret)
	encKey := keyData[:32]
	hmacKey := keyData[32:]

	iv, err := randBytes(16)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(data))
	copy(cipherText, data)

	version := byte(3)
	options := byte(0) // No Password, No HMAC Salt, No Enryption Salt

	msg := make([]byte, 0)
	msg = append(msg, version)
	msg = append(msg, options)
	msg = append(msg, iv...)

	cipherBlock, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	// padd text for encryption
	blockSize := cipherBlock.BlockSize()
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	cipherText = append(cipherText, padText...)

	encrypter := cipher.NewCBCEncrypter(cipherBlock, iv)
	encrypter.CryptBlocks(cipherText, cipherText)

	msg = append(msg, cipherText...)

	hmacSrc := hmac.New(sha256.New, hmacKey)
	hmacSrc.Write(msg)
	hmacVal := hmacSrc.Sum(nil)

	msg = append(msg, hmacVal...)

	msg = append(msg, privateKey.PubKey().SerializeUncompressed()...)

	return msg, nil
}

func randBytes(num int64) ([]byte, error) {
	bits := make([]byte, num)
	_, err := rand.Read(bits)
	if err != nil {
		return nil, err
	}
	return bits, nil
}

// generateSharedSecretRFC4753 generates the shared secret by multiplying the public and private key and
// combining the X and Y component to get the shared secret. This is the old way of generating a shared
// secret used by libbitcoin
func generateSharedSecretRFC4753(privkey *btcec.PrivateKey, pubkey *btcec.PublicKey) []byte {
	x, y := btcec.S256().ScalarMult(pubkey.X, pubkey.Y, privkey.D.Bytes())
	sharedSecretPublicKey := btcec.PublicKey{Curve: btcec.S256(), X: x, Y: y}
	return sharedSecretPublicKey.SerializeUncompressed()
}
