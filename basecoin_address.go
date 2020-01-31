package cnlib

import (
	"errors"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

const (
	bip44purpose = 44
	bip49purpose = 49
	bip84purpose = 84
)

// constants for size in bytes of pieces of a transaction
const (
	p2pkhOutputSize       = 34
	p2shOutputSize        = 32
	p2wpkhOutputSize      = 31
	p2DefaultOutputSize   = 32
	p2pkhInputSize        = 147
	p2shSegwitInputSize   = 91
	p2wpkhSegwitInputSize = 68
	baseSize              = 11
)

// AddressIsBase58CheckEncoded decodes the address, returns true if address is base58check encoded.
func AddressIsBase58CheckEncoded(addr string) error {
	result, _, err := base58.CheckDecode(addr)

	if err != nil {
		return err
	}

	if len(result) > 0 {
		return nil
	}

	return errors.New("address is not base58check encoded")
}

// AddressIsValidSegwitAddress decodes the address, returns true if is a witness type.
func AddressIsValidSegwitAddress(addr string) error {
	params := &chaincfg.MainNetParams
	if strings.HasPrefix(strings.ToLower(addr), "bcrt") {
		params = &chaincfg.RegressionNetParams
	}

	if !strings.HasPrefix(strings.ToLower(addr), "bc") {
		return errors.New("address is not a bech32 encoded segwit address")
	}

	address, err := btcutil.DecodeAddress(addr, params)

	if err != nil {
		return err
	}

	_, okWpkh := address.(*btcutil.AddressWitnessPubKeyHash)
	_, okWsh := address.(*btcutil.AddressWitnessScriptHash)

	if okWpkh || okWsh {
		return nil
	}

	return errors.New("address is not a bech32 encoded segwit address")
}

// HRPFromAddress decodes the given address, and if a SegWit address, returns the HRP.
func (bc *BaseCoin) HRPFromAddress(addr string) (string, error) {
	address, addrErr := btcutil.DecodeAddress(addr, bc.defaultNetParams())

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

func (bc *BaseCoin) bytesPerInput(utxo *UTXO) (int, error) {
	if utxo == nil {
		if bc.Purpose == bip84purpose {
			return p2wpkhSegwitInputSize, nil
		}
		return p2shSegwitInputSize, nil
	}

	if utxo.ImportedPrivateKey != nil {
		addr, err := btcutil.DecodeAddress(utxo.ImportedPrivateKey.SelectedAddress, bc.defaultNetParams())
		if err != nil {
			return 0, err
		}
		switch addr.(type) {
		case *btcutil.AddressPubKeyHash:
			return p2pkhInputSize, nil
		case *btcutil.AddressScriptHash:
			return p2shSegwitInputSize, nil
		case *btcutil.AddressWitnessPubKeyHash:
			return p2wpkhSegwitInputSize, nil
		case *btcutil.AddressWitnessScriptHash:
			return p2wpkhSegwitInputSize, nil
		}
	}

	if utxo.Path != nil {
		if utxo.Path.Purpose == bip84purpose {
			return p2wpkhSegwitInputSize, nil
		}
		return p2shSegwitInputSize, nil
	}

	return 0, errors.New("invalid destination address")
}

func (bc *BaseCoin) bytesPerChangeOuptut() int {
	if bc.Purpose == bip84purpose {
		return p2wpkhOutputSize
	}
	return p2shOutputSize
}

// totalBytes computes number of bytes a tx will be, given number of inputs, destination address, and if includes change or not.
func (bc *BaseCoin) totalBytes(utxos []*UTXO, address string, includeChange bool) (int, error) {
	total := baseSize

	for _, utxo := range utxos {
		bytes, err := bc.bytesPerInput(utxo)
		if err != nil {
			return 0, err
		}
		total += bytes
	}

	if includeChange {
		total = total + bc.bytesPerChangeOuptut()
	}

	addressForSizeEstimation := address
	if address == PlaceholderDestination {
		addressForSizeEstimation = "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"
	}

	outBytes, err := bc.bytesPerOutputAddress(addressForSizeEstimation)
	if err != nil {
		return 0, err
	}
	total += outBytes

	return total, nil
}

func (bc *BaseCoin) bytesPerOutputAddress(addr string) (int, error) {
	dec, decErr := btcutil.DecodeAddress(addr, bc.defaultNetParams())
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
