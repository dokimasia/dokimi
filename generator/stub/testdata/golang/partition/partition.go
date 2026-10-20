// Package partition declares the partition faults whose parameters a double
// cannot compare with a key.
package partition

// Ledger appends lines.
//
//dokimi:stub
type Ledger interface {
	// Log appends the lines of a format, and its partition is variadic.
	//
	//dokimi:fault partition=args
	Log(format string, args ...string) error
}

// Box keeps a value.
//
//dokimi:stub
type Box[T comparable] interface {
	// Put keeps a value, and its partition has the type of a type parameter.
	//
	//dokimi:fault partition=v
	Put(v T) error
}
