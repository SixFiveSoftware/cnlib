package cnlib

// DerivationPath is used to provide information about an address to be generated.
type DerivationPath struct {
	*BaseCoin // Embedded
	Change    int
	Index     int
}

// NewDerivationPath instantiates a new object and sets values.
func NewDerivationPath(bc *BaseCoin, change int, index int) *DerivationPath {
	return &DerivationPath{
		BaseCoin: bc,
		Change:   change,
		Index:    index,
	}
}
