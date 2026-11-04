// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules

import "go.dokimi.dev/spec"

// Unclaimed is the stance of a classification of the shape catalog that
// licenses no claim and that no claim binds as a partner.
type Unclaimed struct {
	// Classification is shape:<name>, mixin:<name> or
	// contract:<name>/<role>, with the names of the shape catalog. It is the
	// key of the entry in unclaimed.json.
	Classification string `json:"-"`

	Stance spec.StanceID `json:"stance"`

	// Note is the reason of the stance for this classification.
	Note string `json:"note"`
}

func (u *Unclaimed) name(c string) { u.Classification = c }

var unclaimed = listed[string, Unclaimed]("unclaimed.json", "unclaimed")

// UnclaimedClassifications returns the stance of every classification
// without a claim, in the order of the classifications. The slice is shared
// with every other call, so a caller does not modify it.
func UnclaimedClassifications() []Unclaimed { return unclaimed().entries }

// UnclaimedOf returns the stance of the classification c, and false when
// the rules have no stance for it.
func UnclaimedOf(c string) (Unclaimed, bool) { return unclaimed().of(c) }
