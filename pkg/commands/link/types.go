package link

// Record declare one split configuration record.
type Record struct {
	MerchantID        string `csv:"MERCHANT ID"`
	AccountHolderCode string `csv:"ACCOUNT HOLDER CODE"`
	StoreID           string `csv:"STORE ID"`
	SplitID           string `csv:"SPLIT ID"`
}
