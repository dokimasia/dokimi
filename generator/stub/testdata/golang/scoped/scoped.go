// Package scoped declares at package scope the names of the packages that the
// double of Store imports.
package scoped

// The variables that take the names of the packages of the double.
var (
	assert = "the clock and the test of the double"
	fault  = "the retry schedule of a method"
	rand   = "the random source of the faults"
	stub   = "the runtime of the doubles"
	time   = "the latency of the methods"
)

// Store lists the items of the shop.
//
//dokimi:stub
type Store interface {
	// Get returns the item of an ID, and selects the calls of a partition
	// fault by the ID.
	//
	//dokimi:fault partition=id
	Get(id string) (string, error)
	// List returns every item, and recovers on the third attempt.
	//
	//dokimi:fault retry=3
	List() ([]string, error)
}
