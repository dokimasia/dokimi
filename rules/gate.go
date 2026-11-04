// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules

import "go.dokimi.dev/spec"

// Gate is a gate of the spec with its condition, and the reason and the
// remedy of a refusal. A licensed claim whose gate fails is a refusal.
type Gate struct {
	spec.Gate `json:"-"`

	// Passes is the condition of the gate, in words. The derivation
	// implements the condition of a gate on the facts of a declaration once
	// for every language. The suite generator of a language implements the
	// condition of a gate on the code that it writes, such as
	// importable-samples.
	Passes string `json:"passes"`

	// Why is the reason of a refusal, and Remedy is what the author changes
	// to remove it. Both write the callable as {callable}, its declaration as
	// {declaration}, the name of a partner as {partner} and the name of a
	// parameter or of a param of a classification as {param}.
	Why    string `json:"why"`
	Remedy string `json:"remedy"`
}

func (g *Gate) join(s *spec.Gate) { g.Gate = *s }

var gates = joined[spec.GateID, spec.Gate, Gate](
	"gates.json",
	"gates",
	spec.Gates,
	func(g *spec.Gate) spec.GateID { return g.ID },
)

// Gates returns each gate of the spec with its rule, in the order of
// gates.json of the spec. The slice is shared with every other call, so a
// caller does not modify it.
func Gates() []Gate { return gates().entries }

// GateOf returns the gate with the identifier id and its rule, and false
// when the spec has no such gate.
func GateOf(id spec.GateID) (Gate, bool) { return gates().of(id) }
