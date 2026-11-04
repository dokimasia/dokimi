// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules

import "go.dokimi.dev/spec"

// DefectForm is a defect form of the spec with the callable that it
// changes and the length of the shortest sequence that shows it.
type DefectForm struct {
	spec.DefectForm `json:"-"`

	// Changes is the callable that the form overrides: callable for the
	// callable of the check, partner:<name> for a partner of its claim,
	// declaration for the first callable of the declaration with a value
	// result, or medium for the medium that the form opens the subject over.
	Changes string `json:"changes"`

	// MinSteps is the length of the shortest sequence of calls that shows
	// the defect.
	MinSteps int `json:"min-steps"`
}

func (f *DefectForm) join(s *spec.DefectForm) { f.DefectForm = *s }

var defectForms = joined[spec.DefectFormID, spec.DefectForm, DefectForm](
	"defects.json", "forms", spec.DefectForms, func(f *spec.DefectForm) spec.DefectFormID { return f.ID },
)

// DefectForms returns each defect form of the spec with its rule, in the
// order of defects.json of the spec. The slice is shared with every other
// call, so a caller does not modify it.
func DefectForms() []DefectForm { return defectForms().entries }

// DefectFormOf returns the defect form with the identifier id and its rule,
// and false when the spec has no such form.
func DefectFormOf(id spec.DefectFormID) (DefectForm, bool) { return defectForms().of(id) }
