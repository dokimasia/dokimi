// Package sequence declares an interface whose methods return sequences.
package sequence

import "iter"

// Item is what a row contains.
type Item struct {
	// Name names the item.
	Name string
}

// Lines is a sequence of lines under a name of its own.
type Lines = iter.Seq[string]

// Table returns its contents as sequences.
//
//dokimi:stub
type Table interface {
	// Names returns a sequence of values.
	Names() iter.Seq[string]
	// Rows returns a sequence of values that can fail.
	Rows(limit int) iter.Seq2[Item, error]
	// Counts returns a sequence of pairs.
	Counts() iter.Seq2[string, int]
	// Tail returns a sequence through the alias Lines.
	Tail() Lines
}
