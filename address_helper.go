package cnlib

import (
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

// AddressIsBase58CheckEncoded decodes the address to determine if address is valid.
func AddressIsBase58CheckEncoded(addr string) bool {
	result, version, err := base58.CheckDecode(addr)

	if err != nil {
		return false
	}

	if len(result) > 0 && version >= 0 {
		return true
	}

	return false
}
