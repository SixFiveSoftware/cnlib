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

// extendedKeyPrefixPairs returns bitcoin-only prefixes for private/public extended keys.
// Supports BIP32 paths: m/44'/0', m/44'/1', m/49'/0', m/49'/1', m/84'/0', m/84'/1'.
func (bc *BaseCoin) extendedKeyPrefixPairs() ([]byte, []byte, error) {
	purpose := bc.Purpose
	coin := bc.Coin

	if bc.isHardened() {
		purpose -= hardenedOffset
		coin -= hardenedOffset
	}

	if purpose == bip44purpose {
		if coin == mainnet {
			// 0x0488ade4, 0x0488b21e
			return []byte{0x04, 0x88, 0xad, 0xe4}, []byte{0x04, 0x88, 0xb2, 0x1e}, nil
		}
		if coin == testnet {
			// 0x04358394, 0x043587cf
			return []byte{0x04, 0x35, 0x83, 0x94}, []byte{0x04, 0x35, 0x87, 0xcf}, nil
		}
		return nil, nil, ErrInvalidCoinValue
	}

	if purpose == bip49purpose {
		if coin == mainnet {
			// 0x049d7878, 0x049d7cb2
			return []byte{0x04, 0x9d, 0x78, 0x78}, []byte{0x04, 0x9d, 0x7c, 0xb2}, nil
		}
		if coin == testnet {
			// 0x044a4e28, 0x044a5262
			return []byte{0x04, 0x4a, 0x4e, 0x28}, []byte{0x04, 0x4a, 0x52, 0x62}, nil
		}
		return nil, nil, ErrInvalidCoinValue
	}

	if purpose == bip84purpose {
		if coin == mainnet {
			// 0x04b2430c, 0x04b24746
			return []byte{0x04, 0xb2, 0x43, 0x0c}, []byte{0x04, 0xb2, 0x47, 0x46}, nil
		}
		if coin == testnet {
			// 0x045f18bc, 0x045f1cf6
			return []byte{0x04, 0x5f, 0x18, 0xbc}, []byte{0x04, 0x5f, 0x1c, 0xf6}, nil
		}
		return nil, nil, ErrInvalidCoinValue
	}

	return nil, nil, ErrInvalidPurposeValue
}
