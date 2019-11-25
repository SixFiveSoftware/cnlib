package cnlib

/// Type Definition

// UTXO is a type used to manage an unspent transaction output.
type UTXO struct {
	Txid        string
	Index       int // must be in UInt32 range
	Amount      int // must be in UInt32 range
	Path        *DerivationPath
	IsConfirmed bool
}

/// Constructor

// NewUTXO instantiates a new UTXO object and returns a ref to it.
func NewUTXO(txid string, index int, amount int, path *DerivationPath, isConfirmed bool) *UTXO {
	u := UTXO{
		Txid:        txid,
		Index:       index,
		Amount:      amount,
		Path:        path,
		IsConfirmed: isConfirmed,
	}
	return &u
}
