package cnlib

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"

	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"

	"github.com/lightningnetwork/lnd/zpay32"
)

/// Type Declarations

// HDWallet represents the user's current wallet.
type HDWallet struct {
	BaseCoin         *BaseCoin
	WalletWords      string // space-separated string of user's recovery words
	masterPrivateKey *hdkeychain.ExtendedKey
	accountPublicKey *hdkeychain.ExtendedKey
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

// NewHDWalletFromWords returns a pointer to an HDWallet, containing the BaseCoin, words, and unexported master private key.
func NewHDWalletFromWords(wordString string, basecoin *BaseCoin) *HDWallet {
	masterKey, err := masterPrivateKey(wordString, basecoin)
	if err != nil {
		return nil
	}
	kf := keyFactory{masterPrivateKey: masterKey}
	pubkey, _, err := kf.accountExtendedPublicKey(basecoin)
	if err != nil {
		return nil
	}
	wallet := HDWallet{BaseCoin: basecoin, WalletWords: wordString, masterPrivateKey: masterKey, accountPublicKey: pubkey}
	return &wallet
}

// NewHDWalletFromAccountExtendedPublicKey returns a pointer to an HDWallet, containing the BaseCoin, empty word list, nil master private key,
// and unexported pointer to extended key for account-level extended master private key. Returns error if unable to parse x/y/zpub.
func NewHDWalletFromAccountExtendedPublicKey(acctPubKeyStr string) (*HDWallet, error) {
	key, err := hdkeychain.NewKeyFromString(acctPubKeyStr)
	if err != nil {
		return nil, err
	}
	basecoin, err := NewBaseCoinFromAccountPubKey(acctPubKeyStr)
	if err != nil {
		return nil, err
	}
	wallet := HDWallet{BaseCoin: basecoin, WalletWords: "", masterPrivateKey: nil, accountPublicKey: key}
	return &wallet, nil
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
	kf := keyFactory{masterPrivateKey: wallet.masterPrivateKey}

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

// ReceiveAddressForIndex returns a receive MetaAddress derived from the current wallet, BaseCoin, and index.
func (wallet *HDWallet) ReceiveAddressForIndex(index int) (*MetaAddress, error) {
	return wallet.metaAddress(0, index)
}

// ChangeAddressForIndex returns a change MetaAddress derived from the current wallet, BaseCoin, and index.
func (wallet *HDWallet) ChangeAddressForIndex(index int) (*MetaAddress, error) {
	return wallet.metaAddress(1, index)
}

// UpdateCoin updates the pointer stored to a new instance of BaseCoin. Fetched MetaAddresses will reflect updated coin.
func (wallet *HDWallet) UpdateCoin(c *BaseCoin) {
	wallet.BaseCoin = c
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
	kf := keyFactory{masterPrivateKey: wallet.masterPrivateKey}
	return kf.signData(message)
}

// SignatureSigningData signs a given message and returns the signature in hex-encoded string format.
func (wallet *HDWallet) SignatureSigningData(message []byte) (string, error) {
	kf := keyFactory{masterPrivateKey: wallet.masterPrivateKey}
	return kf.signatureSigningData(message)
}

// EncryptWithEphemeralKey encrypts a given body (byte slice) using ECDH symmetric key encryption by creating an ephemeral keypair from entropy and given uncompressed public key.
func (wallet *HDWallet) EncryptWithEphemeralKey(entropy []byte, body []byte, recipientUncompressedPubkey string) ([]byte, error) {
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

	w := NewHDWalletFromWords(m, wallet.BaseCoin)
	privateKey, err := w.masterPrivateKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return encrypt(body, privateKey, publicKey)
}

// DecryptWithKeyFromDerivationPath decrypts a given payload with the key derived from given derivation path.
func (wallet *HDWallet) DecryptWithKeyFromDerivationPath(path *DerivationPath, body []byte) ([]byte, error) {
	kf := keyFactory{masterPrivateKey: wallet.masterPrivateKey}

	pk, err := kf.indexPrivateKey(path)
	if err != nil {
		return nil, err
	}

	ecpk, err := pk.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return decrypt(body, ecpk)
}

// EncryptMessage encrypts a payload using signing key (m/42) and recipient's public key.
func (wallet *HDWallet) EncryptMessage(body []byte, recipientUncompressedPubkey string) ([]byte, error) {
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

	return encrypt(body, signingKey, publicKey)
}

// DecryptMessage decrypts a payload using signing key (m/42) and included sender public key (expected to be last 65 bytes of payload).
func (wallet *HDWallet) DecryptMessage(body []byte) ([]byte, error) {
	signingKey, err := wallet.signingPrivateKey()
	if err != nil {
		return nil, err
	}

	return decrypt(body, signingKey)
}

// ImportPrivateKey accepts an encoded private key from a paper wallet/QR code, decodes it, and returns a ref to an ImportedPrivateKey struct, or error if failed.
func (wallet *HDWallet) ImportPrivateKey(encodedKey string) (*ImportedPrivateKey, error) {
	wif, err := btcutil.DecodeWIF(encodedKey)
	if err != nil {
		return nil, err
	}

	serializedPubkey := wif.SerializePubKey()
	hash160 := btcutil.Hash160(serializedPubkey)

	// legacy
	legacy := base58.CheckEncode(hash160, 0)

	// legacy segwit
	ls, err := bip49AddressFromPubkeyHash(hash160, wallet.BaseCoin)
	if err != nil {
		return nil, err
	}

	// native segwit
	ns, err := bip84AddressFromPubkeyHash(hash160, wallet.BaseCoin)
	if err != nil {
		return nil, err
	}

	addrs := []string{legacy, ls, ns}
	joined := strings.Join(addrs, " ")
	info := NewPreviousOutputInfo("", "", 0, 0)
	retval := ImportedPrivateKey{wif: wif, PossibleAddresses: joined, PrivateKeyAsWIF: wif.String(), PreviousOutputInfo: info}
	return &retval, nil
}

// AccountExtendedMasterPublicKey returns the stringified base58 encoded master extended public key.
func (wallet *HDWallet) AccountExtendedMasterPublicKey() (string, error) {
	kf := keyFactory{masterPrivateKey: wallet.masterPrivateKey}
	_, pubkeyString, err := kf.accountExtendedPublicKey(wallet.BaseCoin)
	if err != nil {
		return "", err
	}
	return pubkeyString, nil
}

// BuildTransactionMetadata will generate the tx metadata needed for client to consume.
func (wallet *HDWallet) BuildTransactionMetadata(data *TransactionData) (*TransactionMetadata, error) {
	builder := transactionBuilder{wallet: wallet}
	return builder.buildTxFromData(data)
}

// DecodeLightningInvoice returns a reference to an invoice.Invoice object if valid, or error if invalid.
func (wallet *HDWallet) DecodeLightningInvoice(invoice string) (*LightningInvoice, error) {
	inv, err := zpay32.Decode(invoice, wallet.BaseCoin.defaultNetParams())
	if err != nil {
		return nil, err
	}
	memo := ""
	if inv.Description != nil {
		memo = *inv.Description
	}

	sats := 0
	if inv.MilliSat != nil {
		sats = int(inv.MilliSat.ToSatoshis())
	}

	isExpired := false
	timestampPlusExpiry := inv.Timestamp.Add(inv.Expiry()).Unix()
	expiresAt := time.Unix(timestampPlusExpiry, 0)
	if time.Now().UTC().After(expiresAt) {
		isExpired = true
	}

	return &LightningInvoice{
		NumSatoshis: sats,
		Description: memo,
		IsExpired:   isExpired,
		ExpiresAt:   timestampPlusExpiry,
	}, nil
}

// CompressedPubKeyForPath returns a compressed public key byte slice for a given derivation path in a wallet.
func (wallet *HDWallet) CompressedPubKeyForPath(path *DerivationPath) ([]byte, error) {
	key, err := wallet.publicKey(path)
	if err != nil {
		return nil, err
	}
	return key.SerializeCompressed(), nil
}

// UncompressedPubKeyForPath returns a compressed public key byte slice for a given derivation path in a wallet.
func (wallet *HDWallet) UncompressedPubKeyForPath(path *DerivationPath) ([]byte, error) {
	key, err := wallet.publicKey(path)
	if err != nil {
		return nil, err
	}
	return key.SerializeUncompressed(), nil
}

/// Unexported functions

func (wallet *HDWallet) publicKey(path *DerivationPath) (*btcec.PublicKey, error) {
	if path == nil {
		return nil, errors.New("derivation path cannot be nil")
	}

	keyFactory := keyFactory{masterPrivateKey: wallet.masterPrivateKey}
	privKey, err := keyFactory.indexPrivateKey(path)
	if err != nil {
		return nil, err
	}
	pubKey, err := privKey.ECPubKey()
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func (wallet *HDWallet) metaAddress(change int, index int) (*MetaAddress, error) {
	if index < 0 {
		return nil, errors.New("index cannot be negative")
	}

	path := NewDerivationPath(wallet.BaseCoin, change, index)

	ua, err := newUsableAddressWithDerivationPath(wallet, path)
	if err != nil {
		return nil, err
	}

	return ua.MetaAddress()
}

func hardened(i int) uint32 {
	return hdkeychain.HardenedKeyStart + uint32(i)
}

func masterPrivateKey(wordString string, basecoin *BaseCoin) (*hdkeychain.ExtendedKey, error) {
	seed := bip39.NewSeed(wordString, "")
	defaultNet := basecoin.defaultNetParams()
	masterKey, err := hdkeychain.NewMaster(seed, defaultNet)
	if err != nil {
		return nil, err
	}
	return masterKey, nil
}

func (wallet *HDWallet) signingPrivateKey() (*btcec.PrivateKey, error) {
	kf := keyFactory{masterPrivateKey: wallet.masterPrivateKey}

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
