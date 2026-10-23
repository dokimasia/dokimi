// Package suited keeps the ledgers of the shop, whose declarations are under
// the directive suite.
package suited

import "context"

// Ledger records the entries of the shop under the directive suite alone.
//
//dokimi:suite
type Ledger interface {
	// Record records the entry.
	Record(ctx context.Context, entry string) error
}

// Journal records the lines of the shop under the directives stub and suite.
//
//dokimi:stub
//dokimi:suite
type Journal interface {
	// Write records the line.
	Write(ctx context.Context, line string) error
}

// Cashbox counts the coins of the shop under the directive suite alone.
//
//dokimi:suite
type Cashbox struct {
	coins int
}

// Count returns the number of coins.
func (c *Cashbox) Count() int {
	return c.coins
}

// Tally returns the sum of two counts of coins under the directive suite,
// and a function has no double.
//
//dokimi:suite
func Tally(a, b int) int {
	return a + b
}

// Notebook keeps the notes of the shop without a directive.
type Notebook interface {
	// Note records the note.
	Note(text string)
}
