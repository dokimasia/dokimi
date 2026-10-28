// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

// LawID identifies a law. names.gen.go declares one constant for each.
type LawID string

// Law is a statement about a subject that is true or false. The service
// counts coverage per law.
type Law struct {
	ID LawID `json:"id"`

	// Statement is the law as one sentence.
	Statement string `json:"statement"`

	// Depth is the machinery at which a check witnesses the law in full. A
	// check whose claim does not list it witnesses the law at its shallow
	// end.
	Depth MachineryID `json:"depth"`

	// Since is the version of the spec that added the law.
	Since string `json:"since"`

	// DistinctFrom contains each law that can be false while this one is
	// true, with the case that separates them.
	DistinctFrom []Distinction `json:"distinct-from"`
}

// Distinction is a law that can be false while another law is true.
type Distinction struct {
	Law LawID `json:"law"`

	// Case is the situation in which the two laws disagree.
	Case string `json:"case"`
}

func (l *Law) key() LawID { return l.ID }

var laws = section[LawID, Law]("laws.json", "laws")

// Laws returns every law in the order of laws.json. The slice is shared
// with every other call, so a caller does not modify it.
func Laws() []Law { return laws().entries }

// LawOf returns the law with the identifier id, and false when the spec has
// no such law.
func LawOf(id LawID) (Law, bool) { return laws().of(id) }
