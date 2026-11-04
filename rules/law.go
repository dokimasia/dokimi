// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules

import "go.dokimi.dev/spec"

// Law is a law of the spec with the machinery of its depth and the laws
// that it differs from.
type Law struct {
	spec.Law `json:"-"`

	// Depth is the machinery at which a check witnesses the law in full. A
	// check whose claim does not list it witnesses the law at its shallow
	// end.
	Depth spec.MachineryID `json:"depth"`

	// DistinctFrom contains each law that can be false while this one is
	// true, with the case that separates them.
	DistinctFrom []Distinction `json:"distinct-from"`
}

func (l *Law) join(s *spec.Law) { l.Law = *s }

// Distinction is a law that can be false while another law is true.
type Distinction struct {
	Law spec.LawID `json:"law"`

	// Case is the situation in which the two laws disagree.
	Case string `json:"case"`
}

var laws = joined[spec.LawID, spec.Law, Law](
	"laws.json",
	"laws",
	spec.Laws,
	func(l *spec.Law) spec.LawID { return l.ID },
)

// Laws returns each law of the spec with its rule, in the order of laws.json
// of the spec. The slice is shared with every other call, so a caller does
// not modify it.
func Laws() []Law { return laws().entries }

// LawOf returns the law with the identifier id and its rule, and false when
// the spec has no such law.
func LawOf(id spec.LawID) (Law, bool) { return laws().of(id) }
