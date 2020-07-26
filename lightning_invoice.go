package cnlib

// LightningInvoice is a wrapper type for returning a decoded LN invoice
type LightningInvoice struct {
	NumSatoshis int
	Description string
	IsExpired   bool
	ExpiresAt   int64 // seconds since unix epoch
}
