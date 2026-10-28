// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

// MachineryID identifies a machinery.
type MachineryID string

// Supplier is the part of a run that supplies a machinery.
type Supplier string

// The suppliers of a machinery.
const (
	// ByRuntime is the runtime, which supplies the machinery to every
	// check.
	ByRuntime Supplier = "runtime"

	// ByConsumer is the consumer, through a door of the application. The
	// declaration states the machinery with the param machineries of the
	// directive suite.
	ByConsumer Supplier = "consumer"

	// BySim is the generator sim.
	BySim Supplier = "sim"
)

// Machinery is what a run supplies before a check exercises its claim.
type Machinery struct {
	ID MachineryID `json:"id"`

	// Supplies is what the run supplies, as one sentence.
	Supplies string `json:"supplies"`

	By Supplier `json:"by"`

	// Param is the value of the param machineries of the directive suite
	// that states the machinery, and empty for a machinery that the
	// consumer does not supply.
	Param string `json:"param,omitempty"`

	Since string `json:"since"`
}

func (m *Machinery) key() MachineryID { return m.ID }

var machineries = section[MachineryID, Machinery]("machineries.json", "machineries")

// Machineries returns every machinery in the order of machineries.json. The
// slice is shared with every other call, so a caller does not modify it.
func Machineries() []Machinery { return machineries().entries }

// MachineryOf returns the machinery with the identifier id, and false when
// the spec has no such machinery.
func MachineryOf(id MachineryID) (Machinery, bool) { return machineries().of(id) }
