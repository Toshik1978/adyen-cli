package commands

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
