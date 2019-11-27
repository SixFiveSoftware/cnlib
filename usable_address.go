package cnlib

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

/// Type Definition

// UsableAddress is a wrapper struct that can provide a usable output address.
type UsableAddress struct {
	Wallet         *HDWallet
	DerivationPath *DerivationPath
}

/// Constructors

// NewUsableAddress accepts a wallet and derivation path, and returns a pointer to a UsableAddress.
func NewUsableAddress(wallet *HDWallet, derivationPath *DerivationPath) *UsableAddress {
	ua := UsableAddress{Wallet: wallet, DerivationPath: derivationPath}
	return &ua
}

/// Receiver methods

// MetaAddress returns a meta address with a given path based on wallet's Basecoin, and uncompressed pubkey if a receive address.
func (ua *UsableAddress) MetaAddress() *MetaAddress {
	addr, addrErr := ua.generateAddress()

	if addrErr != nil {
		return nil
	}

	path := ua.DerivationPath
	kf := keyFactory{Wallet: ua.Wallet}
	indexKey := kf.indexPrivateKey(path)
	ecPriv, _ := indexKey.ECPrivKey()
	ecPub := ecPriv.PubKey()
	pubkeyBytes := ecPub.SerializeUncompressed()
	pubkey := ""
	if path.Change == 0 {
		pubkey = hex.EncodeToString(pubkeyBytes)
	}

	ma := MetaAddress{Address: addr, DerivationPath: ua.DerivationPath, UncompressedPublicKey: pubkey}
	return &ma
}

// BIP49AddressFromPubkeyHash returns a P2SH-P2WPKH address from a pubkey's Hash160.
func (ua *UsableAddress) BIP49AddressFromPubkeyHash(hash []byte) string {
	scriptSig, _ := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(hash).Script()
	addrHash, _ := btcutil.NewAddressScriptHash(scriptSig, ua.Wallet.Basecoin.defaultNetParams())
	return addrHash.EncodeAddress()
}

// BIP84AddressFromPubkeyHash returns a native P2WPKH address from a pubkey's Hash160.
func (ua *UsableAddress) BIP84AddressFromPubkeyHash(hash []byte) string {
	addrHash, _ := btcutil.NewAddressWitnessPubKeyHash(hash, ua.Wallet.Basecoin.defaultNetParams())
	return addrHash.EncodeAddress()
}

/// Unexposed methods

func (ua *UsableAddress) generateAddress() (string, error) {
	purpose := ua.DerivationPath.Purpose

	if purpose == 84 {
		return ua.buildSegwitAddress(ua.DerivationPath), nil
	} else if purpose == 49 {
		return ua.buildBIP49Address(ua.DerivationPath), nil
	}
	return "", errors.New("Unrecognized Address Purpose")
}

func (ua *UsableAddress) buildBIP49Address(path *DerivationPath) string {
	kf := keyFactory{Wallet: ua.Wallet}
	indexKey := kf.indexPrivateKey(path)
	ecPriv, _ := indexKey.ECPrivKey()
	ecPub := ecPriv.PubKey()
	pubkeyBytes := ecPub.SerializeCompressed()
	keyHash := btcutil.Hash160(pubkeyBytes)
	return ua.BIP49AddressFromPubkeyHash(keyHash)
}

func (ua *UsableAddress) buildSegwitAddress(path *DerivationPath) string {
	kf := keyFactory{Wallet: ua.Wallet}
	indexKey := kf.indexPrivateKey(path)
	ecPriv, _ := indexKey.ECPrivKey()
	ecPub := ecPriv.PubKey()
	pubkeyBytes := ecPub.SerializeCompressed()
	keyHash := btcutil.Hash160(pubkeyBytes)
	return ua.BIP84AddressFromPubkeyHash(keyHash)
}
