// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules

// ReferenceFormID identifies a reference form of a verified fake.
type ReferenceFormID string

// State is the state of a reference form.
type State string

// The states of a reference form.
const (
	// StateMap is a map from the key of each call to the written value.
	StateMap State = "map"

	// StateList is a list of the appended values in the order of their
	// appends.
	StateList State = "list"
)

// States are the states of a reference form.
var States = []State{StateMap, StateList}

// Effect is what a callable of a reference form does to the state of the
// form.
type Effect string

// Effects are the effects of a callable of a reference form.
var Effects = []Effect{EffectStore, EffectLoad, EffectAppend}

// The effects of a callable of a reference form.
const (
	// EffectStore stores the written value of the callable under its key. A
	// callable whose first result has the type of the written value returns
	// that value.
	EffectStore Effect = "store"

	// EffectLoad returns the value under the key of the callable, or the zero
	// value beside the absent error where the map has none.
	EffectLoad Effect = "load"

	// EffectAppend appends the input of the callable to the list, and returns
	// the length of the list as the position of the new element.
	EffectAppend Effect = "append"
)

// ReferenceForm is the reference form of a verified fake: the state of a
// contract of the shape catalog, and the effect of each role and of each
// writer on that state. The generator fake writes a fake of a declaration
// whose callables the forms cover.
type ReferenceForm struct {
	// ID is the key of the entry in references.json.
	ID ReferenceFormID `json:"-"`

	// Contract is the contract of the shape catalog that the form covers.
	Contract string `json:"contract"`

	State State `json:"state"`

	// Roles contains the effect of each role of the contract.
	Roles map[string]Effect `json:"roles"`

	// Writers is the effect of a writer without a role of the contract, and
	// empty for a form that takes no such writer. A writer is a callable that
	// the catalog classifies with the shape writer or composite-writer.
	Writers Effect `json:"writers,omitempty"`

	Since string `json:"since"`
}

func (f *ReferenceForm) name(id ReferenceFormID) { f.ID = id }

var references = listed[ReferenceFormID, ReferenceForm]("references.json", "references")

// References returns every reference form, in the order of the identifiers.
// The slice is shared with every other call, so a caller does not modify it.
func References() []ReferenceForm { return references().entries }

// ReferenceOf returns the reference form with the identifier id, and false
// when the rules have no such form.
func ReferenceOf(id ReferenceFormID) (ReferenceForm, bool) { return references().of(id) }
