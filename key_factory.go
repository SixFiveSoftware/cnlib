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
	masterPrivateKey *hdkeychain.ExtendedKey
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
	extendedPublicKey, err := accountKey.Neuter()
	if err != nil {
		return "", err
	}
	return extendedPublicKey.String(), nil
}

// Neuter returns a new extended public key from this extended private key.  The
// same extended key will be returned unaltered if it is already an extended
// public key.
//
// As the name implies, an extended public key does not have access to the
// private key, so it is not capable of signing transactions or deriving
// child extended private keys.  However, it is capable of deriving further
// child extended public keys.
// func (k *hdkeychain.ExtendedKey) Neuter(bc *BaseCoin) (*hdkeychain.ExtendedKey, error) {
// 	// Already an extended public key.
// 	if !k.IsPrivate() {
// 		return k, nil
// 	}

// 	// Get the associated public extended key version bytes.
// 	// version, err := chaincfg.HDPrivateKeyToPublicKeyID(k.version)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	version, err := bc.extendedPubKeyVersionPrefix()
// 	if err != nil {
// 		return nil, err
// 	}

// 	pubkey, err := k.ECPubKey()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Convert it to an extended public key.  The key for the new extended
// 	// key will simply be the pubkey of the current extended private key.
// 	//
// 	// This is the function N((k,c)) -> (K, c) from [BIP32].
// 	// return NewExtendedKey(version, k.pubKeyBytes(), k.chainCode, k.parentFP,
// 	// 	k.depth, k.childNum, false), nil
// 	return hdkeychain.NewExtendedKey(version, pubkey.SerializeCompressed(), )
// }

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
