package cnlib

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
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
	addr := ua.generateAddress()
	path := *ua.DerivationPath
	indexKey := ua.Wallet.privateKey(path)
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

/// Unexposed methods

func (ua *UsableAddress) generateAddress() string {
	purpose := ua.DerivationPath.Purpose

	if purpose == 84 {
		return ua.buildSegwitAddress(ua.DerivationPath)
	} else if purpose == 49 {
		return ua.buildBIP49Address(ua.DerivationPath)
	}
	return ""
}

func (ua *UsableAddress) buildBIP49Address(path *DerivationPath) string {
	indexKey := ua.Wallet.privateKey(*path)
	ecPriv, _ := indexKey.ECPrivKey()
	ecPub := ecPriv.PubKey()
	pubkeyBytes := ecPub.SerializeCompressed()
	keyHash := btcutil.Hash160(pubkeyBytes)

	scriptSig, _ := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(keyHash).Script()
	addrHash, _ := btcutil.NewAddressScriptHash(scriptSig, &chaincfg.MainNetParams)

	return addrHash.EncodeAddress()
}

func (ua *UsableAddress) buildSegwitAddress(path *DerivationPath) string {
	indexKey := ua.Wallet.privateKey(*path)
	ecPriv, _ := indexKey.ECPrivKey()
	ecPub := ecPriv.PubKey()
	pubkeyBytes := ecPub.SerializeCompressed()
	keyHash := btcutil.Hash160(pubkeyBytes)
	addrHash, _ := btcutil.NewAddressWitnessPubKeyHash(keyHash, &chaincfg.MainNetParams)
	return addrHash.EncodeAddress()
}
