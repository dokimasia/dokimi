package suited

// Register opens the ledgers of the shop under the directive stub. Its method
// returns a Ledger of another file, whose double the directive suite asks for.
//
//dokimi:stub
type Register interface {
	// Open returns the ledger of the day.
	Open(day string) Ledger
}
