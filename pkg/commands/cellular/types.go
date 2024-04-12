package cellular

// Record declare one cellular record.
type Record struct {
	Serial     string `csv:"SERIAL"`
	TerminalID string `csv:"TERMINAL ID"`
}
