// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/rules"
)

// covered returns each classification that a claim covers, in the form of
// an entry of unclaimed: the classification of each licence and of each
// licence With, and each role of a contract that a claim binds as a
// partner.
func covered() []string {
	var out []string
	add := func(c rules.Classification) {
		switch {
		case c.Shape != "":
			out = append(out, "shape:"+string(c.Shape))
		case c.Mixin != "":
			out = append(out, "mixin:"+string(c.Mixin))
		case c.Contract != "":
			out = append(out, "contract:"+string(c.Contract)+"/"+string(c.Role))
		}
	}
	for _, c := range rules.Claims() {
		var contracts []shape.Contract
		for _, l := range c.Licences {
			add(l.Classification)
			if l.With != nil {
				add(*l.With)
				contracts = append(contracts, l.With.Contract)
			}
			contracts = append(contracts, l.Contract)
		}
		for _, p := range c.Partners {
			if p.Role == "" {
				continue
			}
			for _, contract := range contracts {
				s, ok := shape.SpecOf(contract)
				if ok && slices.ContainsFunc(s.Roles, func(r shape.RoleSpec) bool { return r.Name == p.Role }) {
					add(rules.Classification{Contract: contract, Role: p.Role})
				}
			}
		}
	}
	return out
}

// catalog returns each classification of the shape catalog in the form of
// an entry of unclaimed, with each role of a contract.
func catalog() []string {
	var out []string
	for _, s := range shape.Specs() {
		switch s.Form {
		case shape.FormShape:
			out = append(out, "shape:"+s.Name)
		case shape.FormMixin:
			out = append(out, "mixin:"+s.Name)
		default:
			for _, r := range s.Roles {
				out = append(out, "contract:"+s.Name+"/"+string(r.Name))
			}
		}
	}
	return out
}

func TestUnclaimed(t *testing.T) {
	t.Parallel()

	t.Run("UnclaimedClassifications", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the classifications of unclaimed.json in sorted order", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(rules.UnclaimedClassifications(), func(u rules.Unclaimed) string { return u.Classification }),
				ruleIDs(t, "unclaimed", "unclaimed"),
				"the classifications are the keys of unclaimed.json",
			)
		})

		t.Run("returns a stance for each classification of the shape catalog that no claim covers", func(t *testing.T) {
			t.Parallel()
			claimed := covered()

			for _, c := range catalog() {
				_, unclaimed := rules.UnclaimedOf(c)
				expect.True(
					t,
					slices.Contains(claimed, c) != unclaimed,
					"the classification "+c+" has a claim or a stance, and not both",
				)
			}
		})

		t.Run("returns classifications of the shape catalog with a stance of the spec and a note", func(t *testing.T) {
			t.Parallel()
			all := catalog()

			for _, u := range rules.UnclaimedClassifications() {
				expect.True(
					t,
					slices.Contains(all, u.Classification),
					"the shape catalog has the classification "+u.Classification,
				)
				_, ok := spec.StanceOf(u.Stance)
				expect.True(t, ok, "the stance of "+u.Classification+" is a stance of the spec")
				expect.True(t, u.Note != "", "the classification "+u.Classification+" has a note")
			}
		})
	})

	t.Run("UnclaimedOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the stance of the classification", func(t *testing.T) {
			t.Parallel()

			got, ok := rules.UnclaimedOf("mixin:default-on-error")

			assert.True(t, ok, "the classification has a stance")
			assert.Equal(t, got.Stance, spec.Blocked, "the classification is blocked")
		})

		t.Run("returns false for a classification without a stance", func(t *testing.T) {
			t.Parallel()

			_, ok := rules.UnclaimedOf("mixin:read-after-write")

			assert.False(t, ok, "a licensed classification has no stance")
		})
	})
}
