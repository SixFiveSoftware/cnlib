package cnlib

/// Type Definition

// UTXO is a type used to manage an unspent transaction output. Use `Path` if deriving a private key from wallet's derivation path, or `ImportedPrivateKey` if sweeping a direct private key.
type UTXO struct {
	Txid               string
	Index              int // must be in UInt32 range
	Amount             int // must be in UInt32 range
	Path               *DerivationPath
	ImportedPrivateKey *ImportedPrivateKey
	IsConfirmed        bool
}

/// Constructor

// NewUTXO instantiates a new UTXO object and returns a ref to it.
func NewUTXO(txid string, index int, amount int, path *DerivationPath, importedPrivateKey *ImportedPrivateKey, isConfirmed bool) *UTXO {
	u := UTXO{
		Txid:               txid,
		Index:              index,
		Amount:             amount,
		Path:               path,
		ImportedPrivateKey: importedPrivateKey,
		IsConfirmed:        isConfirmed,
	}
	return &u
}
