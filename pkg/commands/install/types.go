package install

// Record declare one app install record.
type Record struct {
	CompanyID      string `csv:"COMPANY ID"`
	StoreID        string `csv:"STORE ID"`
	TerminalFilter string `csv:"FILTER"`
	TerminalID     string `csv:"TERMINAL ID"`
	PackageName    string `csv:"PACKAGE NAME"`
	VersionName    string `csv:"VERSION NAME"`
	Date           string `csv:"DATE"`
}
