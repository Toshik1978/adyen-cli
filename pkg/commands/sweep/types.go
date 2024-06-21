package sweep

// Record declare one fix sweep configuration record.
type Record struct {
	AccountHolderID string `csv:"ACCOUNT HOLDER ID"`
	BalanceID       string `csv:"BALANCE ID"`
}
