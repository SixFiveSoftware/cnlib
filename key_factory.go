package cnlib

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/hdkeychain"
)

/// Type Definition

// KeyFactory is a struct holding a ref to an HDWallet, with receiver methods to obtain keys relative to the wallet.
type keyFactory struct {
	Wallet *HDWallet
}

/// Receiver methods

func (kf keyFactory) indexPrivateKey(path *DerivationPath) *hdkeychain.ExtendedKey {
	purposeKey, _ := kf.Wallet.masterPrivateKey.Child(hardened(path.Purpose))
	coinKey, _ := purposeKey.Child(hardened(path.Coin))
	accountKey, _ := coinKey.Child(hardened(path.Account))
	changeKey, _ := accountKey.Child(uint32(path.Change))
	indexKey, _ := changeKey.Child(uint32(path.Index))
	return indexKey
}

func (kf keyFactory) indexPublicKey(path *DerivationPath) *btcec.PublicKey {
	ecpub, _ := kf.indexPrivateKey(path).ECPubKey()
	return ecpub
}

func (kf keyFactory) signingMasterKey() *hdkeychain.ExtendedKey {
	masterKey := kf.Wallet.masterPrivateKey
	if masterKey == nil {
		return nil
	}
	childKey, childErr := masterKey.Child(42)
	if childErr != nil {
		return nil
	}
	return childKey
}

// SignData signs a given message and returns the signature in bytes.
func (kf keyFactory) SignData(message []byte) []byte {
	messageHash := chainhash.DoubleHashB(message)
	key, keyErr := kf.signingMasterKey().ECPrivKey()
	if keyErr != nil {
		return nil
	}

	signature, signErr := key.Sign(messageHash)

	if signErr != nil {
		return nil
	}

	verified := signature.Verify(messageHash, key.PubKey())

	if !verified {
		return nil
	}

	return signature.Serialize()
}
