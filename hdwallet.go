package cnlib

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"

	"git.coinninja.net/engineering/cryptor"
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

// ImportedPrivateKey encapsulates the possible receive addresses to check for funds. When found, set that address to `SelectedAddress`.
type ImportedPrivateKey struct {
	wif               *btcutil.WIF
	PossibleAddresses string // space-separated list of addresses
	PrivateKeyAsWIF   string
	SelectedAddress   string
}

// GetFullBIP39WordListString returns all 2,048 BIP39 mnemonic words as a space-separated string.
func GetFullBIP39WordListString() string {
	return strings.Join(wordlists.English, " ")
}

// NewWordListFromEntropy returns a space-separated list of mnemonic words from entropy.
func NewWordListFromEntropy(entropy []byte) (string, error) {
	words, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	valid := bip39.IsMnemonicValid(words)
	if !valid {
		return "", errors.New("invalid mnemonic")
	}
	return words, nil
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
func (wallet *HDWallet) SigningKey() ([]byte, error) {
	ec, err := wallet.signingPrivateKey()
	if err != nil {
		return nil, err
	}
	return ec.Serialize(), nil
}

// SigningPublicKey returns the public key at the m/42 path.
func (wallet *HDWallet) SigningPublicKey() ([]byte, error) {
	kf := keyFactory{Wallet: wallet}

	smk, err := kf.signingMasterKey()
	if err != nil {
		return nil, err
	}

	ec, err := smk.ECPubKey()
	if err != nil {
		return nil, err
	}

	return ec.SerializeCompressed(), nil
}

// CoinNinjaVerificationKeyHexString returns the hex-encoded string of the signing pubkey byte slice.
func (wallet *HDWallet) CoinNinjaVerificationKeyHexString() (string, error) {
	key, err := wallet.SigningPublicKey()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// ReceiveAddressForIndex returns a receive MetaAddress derived from the current wallet, Basecoin, and index.
func (wallet *HDWallet) ReceiveAddressForIndex(index int) (*MetaAddress, error) {
	return wallet.metaAddress(0, index)
}

// ChangeAddressForIndex returns a change MetaAddress derived from the current wallet, Basecoin, and index.
func (wallet *HDWallet) ChangeAddressForIndex(index int) (*MetaAddress, error) {
	return wallet.metaAddress(1, index)
}

// UpdateCoin updates the pointer stored to a new instance of Basecoin. Fetched MetaAddresses will reflect updated coin.
func (wallet *HDWallet) UpdateCoin(c *Basecoin) {
	wallet.Basecoin = c
}

// CheckForAddress scans the wallet for a given address up to a given index on both receive/change chains.
func (wallet *HDWallet) CheckForAddress(a string, upTo int) (*MetaAddress, error) {
	for i := 0; i < upTo; i++ {
		rma, err := wallet.ReceiveAddressForIndex(i)
		if err != nil {
			return nil, err
		}
		cma, err := wallet.ChangeAddressForIndex(i)
		if err != nil {
			return nil, err
		}
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
func (wallet *HDWallet) SignData(message []byte) ([]byte, error) {
	kf := keyFactory{Wallet: wallet}
	return kf.signData(message)
}

// SignatureSigningData signs a given message and returns the signature in hex-encoded string format.
func (wallet *HDWallet) SignatureSigningData(message []byte) (string, error) {
	kf := keyFactory{Wallet: wallet}
	return kf.signatureSigningData(message)
}

// EncryptWithEphemeralKey encrypts a given body (byte slice) using ECDH symmetric key encryption by creating an ephemeral keypair from entropy and given uncompressed public key.
func (wallet *HDWallet) EncryptWithEphemeralKey(body []byte, entropy []byte, recipientUncompressedPubkey string) ([]byte, error) {
	pubkeyBytes, err := hex.DecodeString(recipientUncompressedPubkey)
	if err != nil {
		return nil, err
	}

	publicKey, err := btcec.ParsePubKey(pubkeyBytes, btcec.S256())
	if err != nil {
		return nil, err
	}

	m, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	w := NewHDWalletFromWords(m, wallet.Basecoin)
	privateKey, err := w.masterPrivateKey.ECPrivKey()

	return cryptor.Encrypt(body, privateKey, publicKey)
}

// DecryptWithKeyFromDerivationPath decrypts a given payload with the key derived from given derivation path.
func (wallet *HDWallet) DecryptWithKeyFromDerivationPath(body []byte, path *DerivationPath) ([]byte, error) {
	kf := keyFactory{Wallet: wallet}

	pk, err := kf.indexPrivateKey(path)
	if err != nil {
		return nil, err
	}

	ecpk, err := pk.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return cryptor.Decrypt(body, ecpk)
}

// EncryptWithDefaultKey encrypts a payload using signing key (m/42) and recipient's public key.
func (wallet *HDWallet) EncryptWithDefaultKey(body []byte, recipientUncompressedPubkey string) ([]byte, error) {
	pubkeyBytes, err := hex.DecodeString(recipientUncompressedPubkey)
	if err != nil {
		return nil, err
	}

	publicKey, err := btcec.ParsePubKey(pubkeyBytes, btcec.S256())
	if err != nil {
		return nil, err
	}

	signingKey, err := wallet.signingPrivateKey()
	if err != nil {
		return nil, err
	}

	return cryptor.Encrypt(body, signingKey, publicKey)
}

// DecryptWithDefaultKey decrypts a payload using signing key (m/42) and included sender public key (expected to be last 65 bytes of payload).
func (wallet *HDWallet) DecryptWithDefaultKey(body []byte) ([]byte, error) {
	signingKey, err := wallet.signingPrivateKey()
	if err != nil {
		return nil, err
	}

	return cryptor.Decrypt(body, signingKey)
}

// ImportPrivateKey accepts an encoded private key from a paper wallet/QR code, decodes it, and returns a ref to an ImportedPrivateKey struct, or error if failed.
func (wallet *HDWallet) ImportPrivateKey(encodedKey string) (*ImportedPrivateKey, error) {
	wif, err := btcutil.DecodeWIF(encodedKey)
	if err != nil {
		return nil, err
	}

	serializedPubkey := wif.SerializePubKey()
	hash160 := btcutil.Hash160(serializedPubkey)
	basecoin := NewBaseCoin(84, 0, 0)

	// legacy
	legacy := base58.CheckEncode(hash160, 0)

	// legacy segwit
	ls, err := bip49AddressFromPubkeyHash(hash160, basecoin)
	if err != nil {
		return nil, err
	}

	// native segwit
	ns, err := bip84AddressFromPubkeyHash(hash160, basecoin)
	if err != nil {
		return nil, err
	}

	addrs := []string{legacy, ls, ns}
	joined := strings.Join(addrs, " ")
	retval := ImportedPrivateKey{wif: wif, PossibleAddresses: joined, PrivateKeyAsWIF: wif.String(), SelectedAddress: ""}
	return &retval, nil
}

/// Unexported functions

func (wallet *HDWallet) metaAddress(change int, index int) (*MetaAddress, error) {
	if index < 0 {
		return nil, errors.New("index cannot be negative")
	}

	c := wallet.Basecoin
	path := NewDerivationPath(c.Purpose, c.Coin, c.Account, change, index)

	ua, err := NewUsableAddressWithDerivationPath(wallet, path)
	if err != nil {
		return nil, err
	}

	return ua.MetaAddress()
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

func (wallet *HDWallet) signingPrivateKey() (*btcec.PrivateKey, error) {
	kf := keyFactory{Wallet: wallet}

	smk, err := kf.signingMasterKey()
	if err != nil {
		return nil, err
	}

	ec, err := smk.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return ec, nil
}
