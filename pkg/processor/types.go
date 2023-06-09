package processor

// Config declare processor's configuration.
type Config struct {
	AdyenCalKey      string `env:"ADYEN_CAL_KEY,required"`
	AdyenCalTestKey  string `env:"ADYEN_CAL_TEST_KEY,required"`
	AdyenMgmtKey     string `env:"ADYEN_MGMT_KEY,required"`
	AdyenMgmtTestKey string `env:"ADYEN_MGMT_TEST_KEY,required"`
	AdyenCalURL      string `env:"ADYEN_CAL_URL,required"`
	AdyenCalTestURL  string `env:"ADYEN_CAL_TEST_URL,required"`
	AdyenMgmtURL     string `env:"ADYEN_MGMT_URL,required"`
	AdyenMgmtTestURL string `env:"ADYEN_MGMT_TEST_URL,required"`
}

// LinkRecord declare one split configuration record.
type LinkRecord struct {
	MerchantID        string `csv:"MERCHANT ID"`
	AccountHolderCode string `csv:"ACCOUNT HOLDER CODE"`
	ToastGUID         string `csv:"TOAST GUID"`
	StoreID           string `csv:"STORE ID"`
	SplitID           string `csv:"SPLIT ID"`
}
