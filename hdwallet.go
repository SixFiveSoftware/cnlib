package cnlib

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
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

// ImportedPrivateKey encapsulates the three possible receive addresses to check for funds.
type ImportedPrivateKey struct {
	wif          *btcutil.WIF
	Legacy       string
	LegacySegwit string
	NativeSegwit string
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
	masterKey, err := masterPrivateKey(wordString, basecoin)
	if err != nil {
		return nil
	}
	wallet := HDWallet{Basecoin: basecoin, WalletWords: wordString, masterPrivateKey: masterKey}
	return &wallet
}

/// Receiver functions

// SigningKey returns the private key at the m/42 path.
func (wallet *HDWallet) SigningKey() []byte {
	kf := keyFactory{Wallet: wallet}
	ec, _ := kf.signingMasterKey().ECPrivKey()
	return ec.Serialize()
}

// SigningPublicKey returns the public key at the m/42 path.
func (wallet *HDWallet) SigningPublicKey() []byte {
	kf := keyFactory{Wallet: wallet}
	ec, _ := kf.signingMasterKey().ECPubKey()
	return ec.SerializeCompressed()
}

// CoinNinjaVerificationKeyHexString returns the hex-encoded string of the signing pubkey byte slice.
func (wallet *HDWallet) CoinNinjaVerificationKeyHexString() string {
	return hex.EncodeToString(wallet.SigningPublicKey())
}

// ReceiveAddressForIndex returns a receive MetaAddress derived from the current wallet, Basecoin, and index.
func (wallet *HDWallet) ReceiveAddressForIndex(index int) *MetaAddress {
	return wallet.metaAddress(0, index)
}

// ChangeAddressForIndex returns a change MetaAddress derived from the current wallet, Basecoin, and index.
func (wallet *HDWallet) ChangeAddressForIndex(index int) *MetaAddress {
	return wallet.metaAddress(1, index)
}

// UpdateCoin updates the pointer stored to a new instance of Basecoin. Fetched MetaAddresses will reflect updated coin.
func (wallet *HDWallet) UpdateCoin(c *Basecoin) {
	wallet.Basecoin = c
}

// CheckForAddress scans the wallet for a given address up to a given index on both receive/change chains.
func (wallet *HDWallet) CheckForAddress(a string, upTo int) (*MetaAddress, error) {
	for i := 0; i < upTo; i++ {
		rma := wallet.ReceiveAddressForIndex(i)
		cma := wallet.ChangeAddressForIndex(i)
		if rma.Address == a {
			return rma, nil
		}
		if cma.Address == a {
			return cma, nil
		}
	}
	return nil, errors.New("address not found")
}

// SignData signs a given message and returns the signature in bytes.
func (wallet *HDWallet) SignData(message []byte) []byte {
	kf := keyFactory{Wallet: wallet}
	signature := kf.signData(message)
	return signature
}

// SignatureSigningData signs a given message and returns the signature in hex-encoded string format.
func (wallet *HDWallet) SignatureSigningData(message []byte) string {
	kf := keyFactory{Wallet: wallet}
	str := kf.signatureSigningData(message)
	return str
}

// EncryptWithEphemeralKey encrypts a given body (byte slice) using ECDH symmetric key encryption by creating an ephemeral keypair from entropy and given uncompressed public key.
func (wallet *HDWallet) EncryptWithEphemeralKey(body []byte, entropy []byte, recipientUncompressedPubkey string) ([]byte, error) {
	pubkeyBytes, _ := hex.DecodeString(recipientUncompressedPubkey)
	return Encrypt(body, pubkeyBytes, entropy) // Encrypt(body, pubkeyBytes, privkeyBytes.Serialize())
}

// DecryptWithKeyFromDerivationPath decrypts a given payload with the key derived from given derivation path.
func (wallet *HDWallet) DecryptWithKeyFromDerivationPath(body []byte, path *DerivationPath) ([]byte, error) {
	kf := keyFactory{Wallet: wallet}
	pk := kf.indexPrivateKey(path)
	ecpk, _ := pk.ECPrivKey()
	privkeyBytes := ecpk.Serialize()

	return Decrypt(body, privkeyBytes)
}

// EncryptWithDefaultKey encrypts a payload using signing key (m/42) and recipient's public key.
func (wallet *HDWallet) EncryptWithDefaultKey(body []byte, recipientUncompressedPubkey string) ([]byte, error) {
	pubkeyBytes, _ := hex.DecodeString(recipientUncompressedPubkey)
	privkeyBytes := wallet.SigningKey()
	return Encrypt(body, pubkeyBytes, privkeyBytes)
}

// DecryptWithDefaultKey decrypts a payload using signing key (m/42) and included sender public key (expected to be last 65 bytes of payload).
func (wallet *HDWallet) DecryptWithDefaultKey(body []byte) ([]byte, error) {
	privkeyBytes := wallet.SigningKey()
	return Decrypt(body, privkeyBytes)
}

/// Unexported functions

func (wallet *HDWallet) metaAddress(change int, index int) *MetaAddress {
	if index < 0 {
		return nil
	}
	c := wallet.Basecoin
	path := NewDerivationPath(c.Purpose, c.Coin, c.Account, change, index)
	ua := NewUsableAddress(wallet, path)
	ma := ua.MetaAddress()
	return ma
}

func hardened(i int) uint32 {
	return hdkeychain.HardenedKeyStart + uint32(i)
}

func masterPrivateKey(wordString string, basecoin *Basecoin) (*hdkeychain.ExtendedKey, error) {
	seed := bip39.NewSeed(wordString, "")
	defaultNet := basecoin.defaultNetParams()
	masterKey, err := hdkeychain.NewMaster(seed, defaultNet)
	if err != nil {
		return nil, err
	}
	return masterKey, nil
}

func (wallet *HDWallet) importPrivateKey(encodedKey string) (*ImportedPrivateKey, error) {
	wif, err := btcutil.DecodeWIF(encodedKey)
	if err != nil {
		return nil, err
	}

	serializedPubkey := wif.SerializePubKey()
	hash160 := btcutil.Hash160(serializedPubkey)
	ua := NewUsableAddress(wallet, nil)

	// legacy
	legacy := base58.CheckEncode(hash160, 0)

	// legacy segwit
	ls := ua.BIP49AddressFromPubkeyHash(hash160)

	// native segwit
	ns := ua.BIP84AddressFromPubkeyHash(hash160)

	retval := ImportedPrivateKey{wif: wif, Legacy: legacy, LegacySegwit: ls, NativeSegwit: ns}
	return &retval, nil
}
