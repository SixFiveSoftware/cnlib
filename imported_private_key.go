package cnlib

import "github.com/btcsuite/btcutil"

// ImportedPrivateKey encapsulates the possible receive addresses to check for funds. When found, set that address to `SelectedAddress`.
type ImportedPrivateKey struct {
	wif               *btcutil.WIF
	PossibleAddresses string // space-separated list of addresses
	PrivateKeyAsWIF   string
	*PreviousOutputInfo
}

// PreviousOutputInfo contains selectedAddress, txid, index about the funding utxo.
type PreviousOutputInfo struct {
	SelectedAddress string
	Txid            string
	Index           int
}

// NewPreviousOutputInfo exposes an initializer to the client to provide previous output info to ImportedPrivateKey.
func NewPreviousOutputInfo(selectedAddress string, txid string, index int) *PreviousOutputInfo {
	return &PreviousOutputInfo{SelectedAddress: selectedAddress, Txid: txid, Index: index}
}
