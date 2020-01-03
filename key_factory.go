package cnlib

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/hdkeychain"
)

/// Type Definition

// KeyFactory is a struct holding a ref to an HDWallet, with receiver methods to obtain keys relative to the wallet.
type keyFactory struct {
	Wallet *HDWallet
}

/// Receiver methods

func (kf keyFactory) indexPrivateKey(path *DerivationPath) (*hdkeychain.ExtendedKey, error) {
	purposeKey, err := kf.Wallet.masterPrivateKey.Child(hardened(path.Purpose))
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

func (kf keyFactory) signingMasterKey() (*hdkeychain.ExtendedKey, error) {
	masterKey := kf.Wallet.masterPrivateKey
	if masterKey == nil {
		return nil, errors.New("missing master private key")
	}
	childKey, err := masterKey.Child(42)
	if err != nil {
		return nil, err
	}
	return childKey, nil
}

func (kf keyFactory) signData(message []byte) ([]byte, error) {
	messageHash := chainhash.DoubleHashB(message)

	key, err := kf.signingMasterKey()
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

func (kf keyFactory) signatureSigningData(message []byte) (string, error) {
	sign, err := kf.signData(message)
	if err != nil {
		return "", err
	}

	if len(sign) == 0 {
		return "", errors.New("signature is empty")
	}

	return hex.EncodeToString(sign), nil
}
