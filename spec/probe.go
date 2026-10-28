// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

// ProbeID identifies a probe.
type ProbeID string

// InputKind is the kind of an input of a probe. The template tree of each
// language passes each kind in one form.
type InputKind string

// The kinds of the inputs of a probe.
const (
	// KindCall is a function that calls the subject with bound arguments
	// and returns the results of the subject.
	KindCall InputKind = "call"

	// KindValue is a sample value of a type.
	KindValue InputKind = "value"

	// KindError is an error value or an error type that the probe matches
	// by identity.
	KindError InputKind = "error"

	// KindDuration is a length of time, such as the lifetime of an entry.
	KindDuration InputKind = "duration"

	// KindCount is a whole number, such as the attempts of a retry.
	KindCount InputKind = "count"

	// KindProjection is a function that returns one part of a value, such
	// as its version.
	KindProjection InputKind = "projection"

	// KindScope is a function that runs a body inside a scope of the
	// subject, such as a transaction.
	KindScope InputKind = "scope"

	// KindHook is a function that installs a hook through the subject.
	KindHook InputKind = "hook"

	// KindPredicate is a function that returns whether the subject is in a
	// state, such as leadership.
	KindPredicate InputKind = "predicate"

	// KindClock is the controlled clock of the run.
	KindClock InputKind = "clock"

	// KindValues are named values that the check covers as one set, such as
	// the members of an enumeration or the error values of a package.
	KindValues InputKind = "values"
)

// Input is one input of a probe.
type Input struct {
	Name string    `json:"name"`
	Kind InputKind `json:"kind"`
}

// Probe is the procedure of a check. Each runtime implements every probe in
// its own language, and the probe asserts through the assert library of the
// language.
type Probe struct {
	ID ProbeID `json:"id"`

	Summary string `json:"summary"`

	// Vacuous is the condition under which the probe does not exercise its
	// claim.
	Vacuous string `json:"vacuous"`

	Since string `json:"since"`

	Inputs []Input `json:"inputs"`

	// Reasons lists the reason codes that the probe reports. The handle of
	// a runtime fails a check whose probe reports another code.
	Reasons []ReasonID `json:"reasons"`
}

func (p *Probe) key() ProbeID { return p.ID }

var probes = section[ProbeID, Probe]("probes.json", "probes")

// Probes returns every probe in the order of probes.json. The slice is
// shared with every other call, so a caller does not modify it.
func Probes() []Probe { return probes().entries }

// ProbeOf returns the probe with the identifier id, and false when the spec
// has no such probe.
func ProbeOf(id ProbeID) (Probe, bool) { return probes().of(id) }
