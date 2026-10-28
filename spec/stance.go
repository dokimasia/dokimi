// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

// StanceID identifies a stance.
type StanceID string

// Stance is the reason why a classification has no check.
type Stance struct {
	ID StanceID `json:"id"`

	// Meaning is the reason, as one sentence.
	Meaning string `json:"meaning"`

	// Closes is the change that gives the classification a check.
	Closes string `json:"closes"`

	Since string `json:"since"`
}

func (s *Stance) key() StanceID { return s.ID }

// Unclaimed is the stance of a classification of the shape catalog that
// licenses no claim.
type Unclaimed struct {
	// Classification is shape:<name>, mixin:<name> or
	// contract:<name>/<role>, with the names of the shape catalog.
	Classification string `json:"id"`

	Stance StanceID `json:"stance"`

	// Note is the reason of the stance for this classification.
	Note string `json:"note"`
}

func (u *Unclaimed) key() string { return u.Classification }

var (
	stances   = section[StanceID, Stance]("stances.json", "stances")
	unclaimed = section[string, Unclaimed]("stances.json", "unclaimed")
)

// Stances returns every stance in the order of stances.json. The slice is
// shared with every other call, so a caller does not modify it.
func Stances() []Stance { return stances().entries }

// StanceOf returns the stance with the identifier id, and false when the
// spec has no such stance.
func StanceOf(id StanceID) (Stance, bool) { return stances().of(id) }

// UnclaimedClassifications returns the stance of every classification
// without a claim, in the order of stances.json. The slice is shared with
// every other call, so a caller does not modify it.
func UnclaimedClassifications() []Unclaimed { return unclaimed().entries }

// UnclaimedOf returns the stance of the classification c, and false when the
// spec has no stance for it.
func UnclaimedOf(c string) (Unclaimed, bool) { return unclaimed().of(c) }
