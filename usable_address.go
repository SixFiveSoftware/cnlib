package cnlib

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

/// Type Definition

// UsableAddress is a wrapper struct that can provide a usable output address.
type UsableAddress struct {
	Wallet            *HDWallet
	DerivationPath    *DerivationPath
	derivedPrivateKey *btcec.PrivateKey // derived from master along a derivation path, or specific pk from sweep.
}

/// Constructors

// NewUsableAddressWithDerivationPath accepts a wallet and derivation path, and returns a pointer to a UsableAddress.
func NewUsableAddressWithDerivationPath(wallet *HDWallet, derivationPath *DerivationPath) *UsableAddress {
	kf := keyFactory{Wallet: wallet}
	indexKey := kf.indexPrivateKey(derivationPath)
	ecPriv, _ := indexKey.ECPrivKey()
	ua := UsableAddress{Wallet: wallet, DerivationPath: derivationPath, derivedPrivateKey: ecPriv}
	return &ua
}

// NewUsableAddressWithImportedPrivateKey accepts a wallet and imported private key, and returns a pointer to a UsableAddress.
func NewUsableAddressWithImportedPrivateKey(wallet *HDWallet, importedPrivateKey *ImportedPrivateKey) *UsableAddress {
	ecPriv := importedPrivateKey.wif.PrivKey
	ua := UsableAddress{Wallet: wallet, DerivationPath: nil, derivedPrivateKey: ecPriv}
	return &ua
}

/// Receiver methods

// MetaAddress returns a meta address with a given path based on wallet's Basecoin, and uncompressed pubkey if a receive address. UsableAddress's DerivationPath must not be nil.
func (ua *UsableAddress) MetaAddress() *MetaAddress {
	addr, addrErr := ua.generateAddress()

	if addrErr != nil {
		return nil
	}

	path := ua.DerivationPath
	if path == nil {
		return nil
	}
	ecPub := ua.derivedPrivateKey.PubKey()
	pubkeyBytes := ecPub.SerializeUncompressed()
	pubkey := ""
	if path.Change == 0 {
		pubkey = hex.EncodeToString(pubkeyBytes)
	}

	ma := MetaAddress{Address: addr, DerivationPath: path, UncompressedPublicKey: pubkey}
	return &ma
}

// BIP49AddressFromPubkeyHash returns a P2SH-P2WPKH address from a pubkey's Hash160.
func bip49AddressFromPubkeyHash(hash []byte, basecoin *Basecoin) string {
	scriptSig, _ := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(hash).Script()
	addrHash, _ := btcutil.NewAddressScriptHash(scriptSig, basecoin.defaultNetParams())
	return addrHash.EncodeAddress()
}

// BIP84AddressFromPubkeyHash returns a native P2WPKH address from a pubkey's Hash160.
func bip84AddressFromPubkeyHash(hash []byte, basecoin *Basecoin) string {
	addrHash, _ := btcutil.NewAddressWitnessPubKeyHash(hash, basecoin.defaultNetParams())
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
	ecPub := ua.derivedPrivateKey.PubKey()
	pubkeyBytes := ecPub.SerializeCompressed()
	keyHash := btcutil.Hash160(pubkeyBytes)
	return bip49AddressFromPubkeyHash(keyHash, ua.Wallet.Basecoin)
}

func (ua *UsableAddress) buildSegwitAddress(path *DerivationPath) string {
	ecPub := ua.derivedPrivateKey.PubKey()
	pubkeyBytes := ecPub.SerializeCompressed()
	keyHash := btcutil.Hash160(pubkeyBytes)
	return bip84AddressFromPubkeyHash(keyHash, ua.Wallet.Basecoin)
}

func (ua *UsableAddress) buildCompressedPublicKey() []byte {
	ecPub := ua.derivedPrivateKey.PubKey()
	return ecPub.SerializeCompressed()
}
