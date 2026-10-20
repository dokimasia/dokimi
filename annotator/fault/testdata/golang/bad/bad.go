// Package bad declares the faults that the annotator reports.
package bad

// Ledger records entries.
type Ledger interface {
	// Put keeps an entry, and its retry states an attempt that no failed
	// call can precede.
	//
	//dokimi:fault retry=1 partition=key
	Put(key, entry string) error
	// Len returns the count of the entries, and cannot fail.
	//
	//dokimi:fault retry=3
	Len() int
}
