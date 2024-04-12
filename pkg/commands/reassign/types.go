package reassign

// Record declare one reassign record.
type Record struct {
	Serial     string `csv:"SERIAL"`
	TerminalID string `csv:"TERMINAL ID"`
	MerchantID string `csv:"MERCHANT ID"`
	StoreID    string `csv:"STORE ID"`
}
