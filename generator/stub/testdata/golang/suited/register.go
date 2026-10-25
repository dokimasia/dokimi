package suited

// Register opens the ledgers of the shop under the directive stub. Its method
// returns a Ledger of another file, which gets its double as a handle.
//
//dokimi:stub
type Register interface {
	// Open returns the ledger of the day.
	Open(day string) Ledger
}
