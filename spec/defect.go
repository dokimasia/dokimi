// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

// DefectFormID identifies a form of a planted defect.
type DefectFormID string

// DefectForm is a form of the planted defects. The template tree of each
// language writes a form as a type that wraps a fresh subject and overrides
// the callable that the form changes.
type DefectForm struct {
	ID DefectFormID `json:"id"`

	// Changes is the callable that the form overrides: callable for the
	// callable of the check, or partner:<name> for a partner of its claim.
	Changes string `json:"changes"`

	// Does is what the override does, as one sentence.
	Does string `json:"does"`

	Since string `json:"since"`

	// MinSteps is the length of the shortest sequence of calls that shows
	// the defect. Falsification of a property check generates sequences of
	// at least this length.
	MinSteps int `json:"min-steps"`
}

func (f *DefectForm) key() DefectFormID { return f.ID }

var defectForms = section[DefectFormID, DefectForm]("defects.json", "forms")

// DefectForms returns every defect form in the order of defects.json. The
// slice is shared with every other call, so a caller does not modify it.
func DefectForms() []DefectForm { return defectForms().entries }

// DefectFormOf returns the defect form with the identifier id, and false
// when the spec has no such form.
func DefectFormOf(id DefectFormID) (DefectForm, bool) { return defectForms().of(id) }
