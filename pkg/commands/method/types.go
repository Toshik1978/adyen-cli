package method

// Record declare one add payment methods record.
type Record struct {
	StoreID        string `csv:"STORE ID"`
	PaymentMethods string `csv:"PAYMENT METHODS"`
	Currency       string `csv:"CURRENCY"`
}
