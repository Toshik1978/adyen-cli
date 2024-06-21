package sales

// Record declare one change sales close time record.
type Record struct {
	AccountHolderID string `csv:"ACCOUNT HOLDER ID"`
	BalanceID       string `csv:"BALANCE ID"`
	CloseTime       string `csv:"CLOSE TIME"`
	Delays          int    `csv:"DELAYS"`
}
