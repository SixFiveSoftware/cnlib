package cnlib

import "github.com/btcsuite/btcd/chaincfg"

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
func (bc *Basecoin) GetBech32HRP() string {
	if bc == nil {
		return ""
	}

	basecoin := *bc
	if basecoin.Purpose != 84 {
		return ""
	}
	if basecoin.Coin == 0 {
		return "bc"
	}
	return "bcrt"
}

func (bc *Basecoin) isTestNet() bool {
	if bc.Coin == 0 {
		return false
	}
	return true
}

func (bc *Basecoin) defaultNetParams() *chaincfg.Params {
	if bc.isTestNet() {
		return &chaincfg.RegressionNetParams
	}
	return &chaincfg.MainNetParams
}
