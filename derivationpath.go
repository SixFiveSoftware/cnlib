package cnlib

/// Type Declaration

// DerivationPath is used to provide information about an address to be generated.
type DerivationPath struct {
	Purpose int
	Coin    int
	Account int
	Change  int
	Index   int
}

/// Constructors

// NewDerivationPath instantiates a new object and sets values.
func NewDerivationPath(purpose int, coin int, account int, change int, index int) *DerivationPath {
	dp := DerivationPath{
		Purpose: purpose,
		Coin:    coin,
		Account: account,
		Change:  change,
		Index:   index,
	}
	return &dp
}
