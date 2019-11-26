package cnlib

import (
	"errors"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

const (
	checksumSize      = 4
	p2wpkhProgramSize = 20
	p2wshProgramSize  = 32
)

const (
	bip49purpose = 49
	bip84purpose = 84
)

// constants for size in bytes of pieces of a transaction
const (
	p2pkhOutputSize       = 34
	p2shOutputSize        = 32
	p2wpkhOutputSize      = 31
	p2DefaultOutputSize   = 32
	p2shSegwitInputSize   = 91
	p2wpkhSegwitInputSize = 68
	baseSize              = 11
)

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

/// Unexposed methods

func (ah *AddressHelper) bytesPerInput() uint {
	if ah.Basecoin.Purpose == bip84purpose {
		return p2wpkhSegwitInputSize
	}
	return p2shSegwitInputSize
}

func (ah *AddressHelper) bytesPerChangeOuptut() uint {
	if ah.Basecoin.Purpose == bip84purpose {
		return p2wpkhOutputSize
	}
	return p2shOutputSize
}

// totalBytes computes number of bytes a tx will be, given number of inputs, destination address, and if includes change or not.
func (ah *AddressHelper) totalBytes(numInputs uint16, address string, includeChange bool) (uint, error) {
	total := uint(baseSize)

	total = total + (ah.bytesPerInput() * uint(numInputs))

	if includeChange {
		total = total + ah.bytesPerChangeOuptut()
	}

	outBytes, err := ah.bytesPerOutputAddress(address)
	if err != nil {
		return 0, err
	}
	total += outBytes

	return total, nil
}

func (ah *AddressHelper) bytesPerOutputAddress(addr string) (uint, error) {
	dec, decErr := btcutil.DecodeAddress(addr, ah.Basecoin.defaultNetParams())
	if decErr != nil {
		return 0, decErr
	}

	if _, ok := dec.(*btcutil.AddressPubKey); ok {
		return p2DefaultOutputSize, nil
	}

	if _, ok := dec.(*btcutil.AddressPubKeyHash); ok {
		return p2pkhOutputSize, nil
	}

	if _, ok := dec.(*btcutil.AddressScriptHash); ok {
		return p2shOutputSize, nil
	}

	if _, ok := dec.(*btcutil.AddressWitnessPubKeyHash); ok {
		return p2wpkhOutputSize, nil
	}

	if _, ok := dec.(*btcutil.AddressWitnessScriptHash); ok {
		return p2wpkhOutputSize, nil
	}

	return 0, errors.New("address not supported")
}
