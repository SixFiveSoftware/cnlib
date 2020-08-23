package cnlib

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
)

var (
	BaseCoinBip49MainNet = &BaseCoin{Purpose: 49, Coin: 0, Account: 0}
	BaseCoinBip49TestNet = &BaseCoin{Purpose: 49, Coin: 1, Account: 0}
	BaseCoinBip84MainNet = &BaseCoin{Purpose: 84, Coin: 0, Account: 0}
	BaseCoinBip84TestNet = &BaseCoin{Purpose: 84, Coin: 1, Account: 0}
)

const (
	mainnet = 0
	testnet = 1

	xpub = "xpub"
	ypub = "ypub"
	zpub = "zpub"
	tpub = "tpub"
	upub = "upub"
	vpub = "vpub"
)

var (
	// ErrInvalidPurposeValue describes an error in which the caller
	// passed an invalid purpose value.
	ErrInvalidPurposeValue = errors.New("invalid basecoin purpose value")

	// ErrInvalidCoinValue describes an error in which the caller
	// passed an invalid coin value.
	ErrInvalidCoinValue = errors.New("invalid basecoin coin value")
)

// BaseCoin is used to provide information about the current user's wallet.
type BaseCoin struct {
	Purpose int
	Coin    int
	Account int
}

// NewBaseCoin instantiates a new object and sets values
func NewBaseCoin(purpose int, coin int, account int) *BaseCoin {
	return &BaseCoin{Purpose: purpose, Coin: coin, Account: account}
}

// NewBaseCoinFromAccountPubKey returns a new BaseCoin pointer based on prefix, or error if unrecognized.
func NewBaseCoinFromAccountPubKey(key string) (*BaseCoin, error) {
	// prefix := key[:4]
	var bc *BaseCoin

	dec, _, err := base58.CheckDecode(key)
	if err != nil {
		return nil, err
	}

	decoded := []byte{0x04}
	decoded = append(decoded, dec...)

	acctBytes := decoded[9:13]                                   // bytes 9 through 12 are the 4 bytes of the account "child number", but first byte 0x04 is dropped in CheckDecode
	acctIndex := binary.BigEndian.Uint32(acctBytes) & 0x0FFFFFFF // if hardened, remove hardened offset
	acct := int(acctIndex)

	prefix := decoded[:4]
	if bytes.Equal(prefix, pubkeyIDs[xpub]) {
		bc = &BaseCoin{Purpose: bip44purpose, Coin: mainnet, Account: acct}
	} else if bytes.Equal(prefix, pubkeyIDs[ypub]) {
		bc = &BaseCoin{Purpose: bip49purpose, Coin: mainnet, Account: acct}
	} else if bytes.Equal(prefix, pubkeyIDs[zpub]) {
		bc = &BaseCoin{Purpose: bip84purpose, Coin: mainnet, Account: acct}
	} else if bytes.Equal(prefix, pubkeyIDs[tpub]) {
		bc = &BaseCoin{Purpose: bip44purpose, Coin: testnet, Account: acct}
	} else if bytes.Equal(prefix, pubkeyIDs[upub]) {
		bc = &BaseCoin{Purpose: bip49purpose, Coin: testnet, Account: acct}
	} else if bytes.Equal(prefix, pubkeyIDs[vpub]) {
		bc = &BaseCoin{Purpose: bip84purpose, Coin: testnet, Account: acct}
	}

	if bc != nil {
		return bc, nil
	}

	return nil, errors.New("unrecognized account key prefix")
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

func (bc *BaseCoin) defaultExtendedPubkeyType() (string, error) {
	if bc.Purpose == bip44purpose {
		if bc.Coin == mainnet {
			return xpub, nil
		}
		if bc.Coin == testnet {
			return tpub, nil
		}
		return "", ErrInvalidCoinValue
	}
	if bc.Purpose == bip49purpose {
		if bc.Coin == mainnet {
			return ypub, nil
		}
		if bc.Coin == testnet {
			return upub, nil
		}
		return "", ErrInvalidCoinValue
	}
	if bc.Purpose == bip84purpose {
		if bc.Coin == mainnet {
			return zpub, nil
		}
		if bc.Coin == testnet {
			return vpub, nil
		}
		return "", ErrInvalidCoinValue
	}
	return "", ErrInvalidPurposeValue
}

func (bc *BaseCoin) defaultNetParams() *chaincfg.Params {
	if bc.isTestNet() {
		return &chaincfg.TestNet3Params
	}
	return &chaincfg.MainNetParams
}
