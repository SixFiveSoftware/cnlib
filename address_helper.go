package cnlib

import (
	"errors"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

const (
	p2pk   = 0
	p2pkh  = 1
	p2sh   = 2
	p2wpkh = 3
	p2wsh  = 4
)

const checksumSize = 4
// AddressHelper is a struct with helper functions to provide info about addresses.
type AddressHelper struct {
	Basecoin *Basecoin
}

// NewAddressHelper returns a ref to a new AddressHelper object, given a *Basecoin.
func NewAddressHelper(basecoin *Basecoin) *AddressHelper {
	ah := AddressHelper{Basecoin: basecoin}
	return &ah
}

// AddressIsBase58CheckEncoded decodes the address, returns true if address is base58check encoded.
func (*AddressHelper) AddressIsBase58CheckEncoded(addr string) bool {
	result, version, err := base58.CheckDecode(addr)

	if err != nil {
		return false
	}

	if len(result) > 0 && version >= 0 {
		return true
	}

	return false
}

// AddressIsValidSegwitAddress decodes the address, returns true if is a witness type.
func (ah *AddressHelper) AddressIsValidSegwitAddress(addr string) bool {

	address, addrErr := btcutil.DecodeAddress(addr, ah.Basecoin.defaultNetParams())

	if addrErr != nil {
		return false
	}

	_, okWpkh := address.(*btcutil.AddressWitnessPubKeyHash)
	_, okWsh := address.(*btcutil.AddressWitnessScriptHash)

	return okWpkh || okWsh
}

// HRPFromAddress decodes the given address, and if a SegWit address, returns the HRP.
func (ah *AddressHelper) HRPFromAddress(addr string) (string, error) {
	address, addrErr := btcutil.DecodeAddress(addr, ah.Basecoin.defaultNetParams())

	if addrErr != nil {
		return "", errors.New("failed to decode address")
	}

	wpkhAddr, okWpkh := address.(*btcutil.AddressWitnessPubKeyHash)
	if okWpkh {
		return wpkhAddr.Hrp(), nil
	}

	wshAddr, okWsh := address.(*btcutil.AddressWitnessScriptHash)
	if okWsh {
		return wshAddr.Hrp(), nil
	}

	return "", errors.New("invalid segwit address")
}
