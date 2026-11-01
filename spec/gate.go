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
