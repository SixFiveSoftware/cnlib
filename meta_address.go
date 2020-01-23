package cnlib

/// Type Definition

// MetaAddress is a model object which holds meta info about an address.
type MetaAddress struct {
	Address               string
	DerivationPath        *DerivationPath
	UncompressedPublicKey string
}

/// Constructors

// NewMetaAddress creates and returns a pointer to a MetaAddress object.
func NewMetaAddress(address string, path *DerivationPath, uncompressedPublicKey string) *MetaAddress {
	ma := MetaAddress{Address: address, DerivationPath: path, UncompressedPublicKey: uncompressedPublicKey}
	return &ma
}

// IsReceiveAddress computes if the address is a receive address, vs change address.
func (ma *MetaAddress) IsReceiveAddress() bool {
	change := ma.DerivationPath.Change
	return change == 0
}
