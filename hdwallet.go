package cnlib

import (
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"

	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

/// Type Declarations

// HDWallet represents the user's current wallet.
type HDWallet struct {
	Basecoin         *Basecoin
	WalletWords      string // space-separated string of user's recovery words
	masterPrivateKey *hdkeychain.ExtendedKey
}

// GetFullBIP39WordListString returns all 2,048 BIP39 mnemonic words as a space-separated string.
func GetFullBIP39WordListString() string {
	return strings.Join(wordlists.English, " ")
}

// NewWordListFromEntropy returns a space-separated list of mnemonic words from entropy.
func NewWordListFromEntropy(entropy []byte) string {
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

// NewHDWalletFromWords returns a pointer to an HDWallet, containing the Basecoin, words, and unexported master private key.
func NewHDWalletFromWords(wordString string, basecoin *Basecoin) *HDWallet {
	masterKey, err := masterPrivateKey(wordString)
	if err != nil {
		return nil
	}
	wallet := HDWallet{Basecoin: basecoin, WalletWords: wordString, masterPrivateKey: masterKey}
	return &wallet
}

/// Receiver functions

// SigningKey returns the private key at the m/42 path.
func (wallet *HDWallet) SigningKey() []byte {
	ec, _ := wallet.signingMasterKey().ECPrivKey()
	return ec.Serialize()
}

// SigningPublicKey returns the public key at the m/42 path.
func (wallet *HDWallet) SigningPublicKey() []byte {
	ec, _ := wallet.signingMasterKey().ECPubKey()
	return ec.SerializeCompressed()
}

// ReceiveAddressAtIndex returns a receive address at given path based on wallet's Basecoin.
func (wallet *HDWallet) ReceiveAddressAtIndex(index int) string {
	path := DerivationPath{wallet.Basecoin.Purpose, wallet.Basecoin.Coin, wallet.Basecoin.Account, 0, index}
	indexKey := privateKey(wallet.masterPrivateKey, path)
	ecPriv, _ := indexKey.ECPrivKey()
	ecPub := ecPriv.PubKey()
	pubkeyBytes := ecPub.SerializeCompressed()
	keyHash := btcutil.Hash160(pubkeyBytes)
	defaultNet := &chaincfg.MainNetParams
	if wallet.Basecoin.Purpose == 84 {
		addrHash, _ := btcutil.NewAddressWitnessPubKeyHash(keyHash, defaultNet)
		return addrHash.EncodeAddress()
	}
	return ""
}

/// Unexported functions

func hardened(i int) uint32 {
	return hdkeychain.HardenedKeyStart + uint32(i)
}

func privateKey(masterKey *hdkeychain.ExtendedKey, derivationPath DerivationPath) *hdkeychain.ExtendedKey {
	purposeKey, _ := masterKey.Child(hardened(derivationPath.Purpose))
	coinKey, _ := purposeKey.Child(hardened(derivationPath.Coin))
	accountKey, _ := coinKey.Child(hardened(derivationPath.Account))
	changeKey, _ := accountKey.Child(uint32(derivationPath.Change))
	indexKey, _ := changeKey.Child(uint32(derivationPath.Index))
	return indexKey
}

func masterPrivateKey(wordString string) (*hdkeychain.ExtendedKey, error) {
	seed := bip39.NewSeed(wordString, "")
	defaultNet := chaincfg.MainNetParams
	masterKey, err := hdkeychain.NewMaster(seed, &defaultNet)
	if err != nil {
		return nil, err
	}
	return masterKey, nil
}

func (wallet *HDWallet) signingMasterKey() *hdkeychain.ExtendedKey {
	masterKey := wallet.masterPrivateKey
	if masterKey == nil {
		return nil
	}
	childKey, childErr := masterKey.Child(42)
	if childErr != nil {
		return nil
	}
	return childKey
}
