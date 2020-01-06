package cnlib

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// SigningPrivateKey returns the private key used for signing
func (wallet *HDWallet) SigningPrivateKey() (*btcec.PrivateKey, error) {
	smk, err := wallet.SigningMasterKey()
	if err != nil {
		return nil, err
	}

	ec, err := smk.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return ec, nil
}

// SigningKey returns the private key at the m/42 path.
func (wallet *HDWallet) SigningKey() ([]byte, error) {
	ec, err := wallet.SigningPrivateKey()
	if err != nil {
		return nil, err
	}
	return ec.Serialize(), nil
}

// SigningPublicKey returns the public key at the m/42 path.
func (wallet *HDWallet) SigningPublicKey() ([]byte, error) {
	smk, err := wallet.SigningMasterKey()
	if err != nil {
		return nil, err
	}

	ec, err := smk.ECPubKey()
	if err != nil {
		return nil, err
	}

	return ec.SerializeCompressed(), nil
}

func (wallet *HDWallet) IndexPrivateKey(path *DerivationPath) (*hdkeychain.ExtendedKey, error) {
	purposeKey, err := wallet.masterPrivateKey.Child(hardened(path.Purpose))
	if err != nil {
		return nil, err
	}
	coinKey, err := purposeKey.Child(hardened(path.Coin))
	if err != nil {
		return nil, err
	}
	accountKey, err := coinKey.Child(hardened(path.Account))
	if err != nil {
		return nil, err
	}
	changeKey, err := accountKey.Child(uint32(path.Change))
	if err != nil {
		return nil, err
	}
	indexKey, err := changeKey.Child(uint32(path.Index))
	if err != nil {
		return nil, err
	}
	return indexKey, nil
}

func (wallet *HDWallet) SigningMasterKey() (*hdkeychain.ExtendedKey, error) {
	masterKey := wallet.masterPrivateKey
	if masterKey == nil {
		return nil, errors.New("missing master private key")
	}
	childKey, err := masterKey.Child(42)
	if err != nil {
		return nil, err
	}
	return childKey, nil
}

func (wallet *HDWallet) SignData(message []byte) ([]byte, error) {
	messageHash := chainhash.DoubleHashB(message)

	key, err := wallet.SigningMasterKey()
	if err != nil {
		return nil, err
	}

	privKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}

	signature, err := privKey.Sign(messageHash)

	if err != nil {
		return nil, err
	}

	verified := signature.Verify(messageHash, privKey.PubKey())

	if !verified {
		return nil, errors.New("failed to sign data")
	}

	return signature.Serialize(), nil
}

func (wallet *HDWallet) SignatureSigningData(message []byte) (string, error) {
	sign, err := wallet.SignData(message)
	if err != nil {
		return "", err
	}

	if len(sign) == 0 {
		return "", errors.New("signature is empty")
	}

	return hex.EncodeToString(sign), nil
}
