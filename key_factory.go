package cnlib

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
)

/// Type Definition

// KeyFactory is a struct holding a ref to an HDWallet, with receiver methods to obtain keys relative to the wallet.
type keyFactory struct {
	masterPrivateKey *hdkeychain.ExtendedKey
}

var pubkeyIDs = map[string][]byte{
	"xpub": []byte{0x04, 0x88, 0xb2, 0x1e}, // m/44'/0'
	"ypub": []byte{0x04, 0x9d, 0x7c, 0xb2}, // m/49'/0'
	"zpub": []byte{0x04, 0xb2, 0x47, 0x46}, // m/84'/0'
	"tpub": []byte{0x04, 0x35, 0x87, 0xcf}, // m/44'/1'
	"upub": []byte{0x04, 0x4a, 0x52, 0x62}, // m/49'/1'
	"vpub": []byte{0x04, 0x5f, 0x1c, 0xf6}, // m/84'/1'
}

/// Receiver methods

func (kf keyFactory) indexPrivateKey(path *DerivationPath) (*hdkeychain.ExtendedKey, error) {
	purposeKey, err := kf.masterPrivateKey.Child(hardened(path.Purpose))
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

func (kf keyFactory) accountExtendedPublicKey(bc *BaseCoin) (string, error) {
	// derive account child
	purposeKey, err := kf.masterPrivateKey.Child(hardened(bc.Purpose))
	if err != nil {
		return "", err
	}
	coinKey, err := purposeKey.Child(hardened(bc.Coin))
	if err != nil {
		return "", err
	}
	accountKey, err := coinKey.Child(hardened(bc.Account))
	if err != nil {
		return "", err
	}

	// get extended pubkey
	extendedPublicKey, err := accountKey.Neuter()
	if err != nil {
		return "", err
	}

	// base58check encode extended pubkey
	neutered := extendedPublicKey.String()

	// get appropriate prefix
	idType, err := bc.defaultExtendedPubkeyType()
	if err != nil {
		return "", err
	}
	newPrefix := pubkeyIDs[idType]

	// decode
	decoded, version, err := base58.CheckDecode(neutered)
	if err != nil {
		return "", err
	}

	if version != newPrefix[0] {
		return "", errors.New("version mismatch when decoding account pubkey")
	}

	// swap bytes. `version` has first byte, and needs to match first byte of prefix.
	// `temp` does not need `version` to be first byte, CheckEncode will do that.
	temp := make([]byte, len(decoded))
	copy(temp[:3], newPrefix[1:4])
	copy(temp[3:], decoded[3:])

	// re-encode
	encoded := base58.CheckEncode(temp, version)

	return encoded, nil
}

func (kf keyFactory) signingMasterKey() (*hdkeychain.ExtendedKey, error) {
	masterKey := kf.masterPrivateKey
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
