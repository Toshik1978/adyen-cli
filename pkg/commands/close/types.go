package close

// Record declare one close store record.
type Record struct {
	AccountHolderCode string `csv:"ACCOUNT HOLDER CODE"`
	StoreID           string `csv:"STORE ID"`
}
