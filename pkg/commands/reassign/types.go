package reassign

// Record declare one split configuration record.
type Record struct {
	TerminalID string `csv:"TERMINAL ID"`
	MerchantID string `csv:"MERCHANT ID"`
	StoreID    string `csv:"STORE ID"`
}
