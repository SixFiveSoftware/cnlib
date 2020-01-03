package cnlib

import "github.com/btcsuite/btcd/chaincfg"

import "errors"

/// Type Declaration

// Basecoin is used to provide information about the current user's wallet.
type Basecoin struct {
	Purpose int
	Coin    int
	Account int
}

/// Constructors

// NewBaseCoin instantiates a new object and sets values
func NewBaseCoin(purpose int, coin int, account int) *Basecoin {
	bc := Basecoin{Purpose: purpose, Coin: coin, Account: account}
	return &bc
}

/// Receiver methods

// UpdatePurpose updates the purpose value on the BaseCoin receiver.
func (bc *Basecoin) UpdatePurpose(purpose int) {
	bc.Purpose = purpose
}

// UpdateCoin updates the coin value on the BaseCoin receiver.
func (bc *Basecoin) UpdateCoin(coin int) {
	bc.Coin = coin
}

// UpdateAccount updates the account value on the BaseCoin receiver.
func (bc *Basecoin) UpdateAccount(account int) {
	bc.Account = account
}

// GetBech32HRP returns a Bech32 HRP string derived from Purpose and Coin
func (bc *Basecoin) GetBech32HRP() (string, error) {
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

func (bc *Basecoin) isTestNet() bool {
	return bc.Coin != 0
}

func (bc *Basecoin) defaultNetParams() *chaincfg.Params {
	if bc.isTestNet() {
		return &chaincfg.RegressionNetParams
	}
	return &chaincfg.MainNetParams
}
