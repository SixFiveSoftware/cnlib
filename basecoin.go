package cnlib

import (
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
)

var (
	BaseCoinBip49MainNet = &BaseCoin{Purpose: 49, Coin: 0, Account: 0}
	BaseCoinBip49TestNet = &BaseCoin{Purpose: 49, Coin: 1, Account: 0}
	BaseCoinBip84MainNet = &BaseCoin{Purpose: 84, Coin: 0, Account: 0}
	BaseCoinBip84TestNet = &BaseCoin{Purpose: 84, Coin: 1, Account: 0}
)

const (
	hardenedOffset = 0x80000000

	mainnet = 0
	testnet = 1
)

var (
	// ErrInvalidPurposeValue describes an error in which the caller
	// passed an invalid purpose value.
	ErrInvalidPurposeValue = errors.New("invalid purpose value")

	// ErrInvalidCoinValue describes an error in which the caller
	// passed an invalid coin value.
	ErrInvalidCoinValue = errors.New("invalid coin value")
)

// BaseCoin is used to provide information about the current user's wallet.
type BaseCoin struct {
	Purpose int
	Coin    int
	Account int
	params  *chaincfg.Params
}

// NewBaseCoin instantiates a new object and sets values
func NewBaseCoin(purpose int, coin int, account int) *BaseCoin {
	return &BaseCoin{Purpose: purpose, Coin: coin, Account: account, params: nil}
}

func (bc *BaseCoin) isHardened() bool {
	return (bc.Purpose > hardenedOffset) && (bc.Coin > hardenedOffset) && (bc.Account > hardenedOffset)
}

// UpdatePurpose updates the purpose value on the BaseCoin receiver.
func (bc *BaseCoin) UpdatePurpose(purpose int) {
	bc.Purpose = purpose
}

// UpdateCoin updates the coin value on the BaseCoin receiver.
func (bc *BaseCoin) UpdateCoin(coin int) {
	bc.Coin = coin
}

// UpdateAccount updates the coin account on the BaseCoin receiver.
func (bc *BaseCoin) UpdateAccount(account int) {
	bc.Account = account
}

// GetBech32HRP returns a Bech32 HRP string derived from Purpose and Coin
func (bc *BaseCoin) GetBech32HRP() (string, error) {
	if bc == nil {
		return "", errors.New("no basecoin provided")
	}

	if bc.Purpose != 84 {
		return "", errors.New("basecoin purpose is not a segwit purpose")
	}
	if bc.Coin == 0 {
		return "bc", nil
	}
	return "bcrt", nil
}

func (bc *BaseCoin) isTestNet() bool {
	return bc.Coin != 0
}

func (bc *BaseCoin) defaultNetParams() *chaincfg.Params {
	if bc.params != nil {
		return bc.params
	}
	params := &chaincfg.MainNetParams
	if bc.isTestNet() {
		params = &chaincfg.RegressionNetParams
	}

	privKeyID, pubKeyID, err := bc.extendedKeyPrefixPairs()
	if err != nil {
		return params
	}

	var privCopy [4]byte
	copy(privCopy[:], privKeyID[:4])
	params.HDPrivateKeyID = privCopy

	var pubCopy [4]byte
	copy(pubCopy[:], pubKeyID[:4])
	params.HDPublicKeyID = pubCopy

	params.Net++
	err = chaincfg.Register(params)
	if err != nil {
		return params
	}
	params.Net--

	bc.params = params
	return bc.params
}

func (bc *BaseCoin) extendedKeyPrefixPairs() ([]byte, []byte, error) {
	return nil, nil, nil
}
