package cnlib

/// Type Definitions

// TransactionChangeMetadata holds info about the change back to the user's wallet as an output of a transaction.
type TransactionChangeMetadata struct {
	Address   string
	Path      *DerivationPath
	VoutIndex int
}

// TransactionMetadata is the main object containing the txid and encoded tx for an outgoing transaction, with associated change metadata, if necessary.
type TransactionMetadata struct {
	Txid      string
	EncodedTx string
	*TransactionChangeMetadata
}
