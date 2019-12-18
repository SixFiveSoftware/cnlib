package cnlib

import (
	"errors"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
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
func AddressIsBase58CheckEncoded(addr string) (bool, error) {
	result, version, err := base58.CheckDecode(addr)

	if err != nil {
		return false, err
	}

	if len(result) > 0 && version >= 0 {
		return true, nil
	}

	return false, errors.New("address is not base58 encoded")
}

// AddressIsValidSegwitAddress decodes the address, returns true if is a witness type.
func AddressIsValidSegwitAddress(addr string) (bool, error) {
	params := &chaincfg.MainNetParams
	if strings.HasPrefix(strings.ToLower(addr), "bcrt") {
		params = &chaincfg.RegressionNetParams
	}

	if !strings.HasPrefix(strings.ToLower(addr), "bc") {
		return false, errors.New("address is not a bech32 encoded segwit address")
	}

	address, err := btcutil.DecodeAddress(addr, params)

	if err != nil {
		return false, err
	}

	_, okWpkh := address.(*btcutil.AddressWitnessPubKeyHash)
	_, okWsh := address.(*btcutil.AddressWitnessScriptHash)

	if okWpkh || okWsh {
		return true, nil
	}

	return false, errors.New("address is not a bech32 encoded segwit address")
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

func (ah *AddressHelper) bytesPerInput(utxo *UTXO) int {
	purpose := bip49purpose
	if utxo == nil {
		if ah.Basecoin.Purpose == bip84purpose {
			purpose = bip84purpose
		}
	}

	if utxo.Path.Purpose == bip84purpose {
		purpose = bip84purpose
	}

	if purpose == bip49purpose {
		return p2shSegwitInputSize
	}
	return p2wpkhSegwitInputSize
}

func (ah *AddressHelper) bytesPerChangeOuptut() int {
	if ah.Basecoin.Purpose == bip84purpose {
		return p2wpkhOutputSize
	}
	return p2shOutputSize
}

// totalBytes computes number of bytes a tx will be, given number of inputs, destination address, and if includes change or not.
func (ah *AddressHelper) totalBytes(utxos []*UTXO, address string, includeChange bool) (int, error) {
	total := baseSize

	for _, utxo := range utxos {
		total += ah.bytesPerInput(utxo)
	}

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

func (ah *AddressHelper) bytesPerOutputAddress(addr string) (int, error) {
	dec, decErr := btcutil.DecodeAddress(addr, ah.Basecoin.defaultNetParams())
	if decErr != nil {
		return 0, decErr
	}

	switch dec.(type) {
	case *btcutil.AddressPubKey:
		return p2DefaultOutputSize, nil
	case *btcutil.AddressPubKeyHash:
		return p2pkhOutputSize, nil
	case *btcutil.AddressScriptHash:
		return p2shOutputSize, nil
	case *btcutil.AddressWitnessPubKeyHash:
		return p2wpkhOutputSize, nil
	case *btcutil.AddressWitnessScriptHash:
		return p2wpkhOutputSize, nil
	}

	return 0, errors.New("address not supported")
}
