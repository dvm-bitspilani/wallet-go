package helpers

// GetValidTransactionPairs check with seniors if this is the combo they're looking for
func GetValidTransactionPairs() [][]string {
	return [][]string{
		{"bitsian", "bitsian"},
		{"bitsian", "vendor"},
		{"participant", "participant"},
		{"participant", "vendor"},
		{"teller", "participant"},
	}
}
