// Package scoped declares at package scope the names of the packages that the
// suite of Hyphenate imports.
package scoped

import "strings"

// The variables that take the names of the packages of the suite.
var (
	conform = "the runtime of the checks"
	context = "the context of a call"
	probe   = "the probes of the derived checks"
	prop    = "the generators of the inputs"
	spec    = "the sequential specification of the kit"
	testing = "the test of the suite"
	time    = "the ceilings of the budget"
)

// Hyphenate returns the words of a title in lower case, joined by hyphens.
//
//dokimi:suite
//dokimi:budget p99=2ms mean=1ms
func Hyphenate(title string) string {
	return strings.Join(strings.Fields(strings.ToLower(title)), "-")
}
