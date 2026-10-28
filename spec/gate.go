// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

// GateID identifies a gate.
type GateID string

// Gate is a condition on the facts of the subject of a claim. A licensed
// claim whose gate fails is a refusal with the reason and the remedy of the
// gate.
type Gate struct {
	ID GateID `json:"id"`

	// Passes is the condition of the gate, in words. The derivation
	// implements it once for every language.
	Passes string `json:"passes"`

	// Why is the reason of a refusal, and Remedy is what the author changes
	// to remove it. Both name the callable as {callable}.
	Why    string `json:"why"`
	Remedy string `json:"remedy"`

	Since string `json:"since"`
}

func (g *Gate) key() GateID { return g.ID }

var gates = section[GateID, Gate]("gates.json", "gates")

// Gates returns every gate in the order of gates.json. The slice is shared
// with every other call, so a caller does not modify it.
func Gates() []Gate { return gates().entries }

// GateOf returns the gate with the identifier id, and false when the spec
// has no such gate.
func GateOf(id GateID) (Gate, bool) { return gates().of(id) }
