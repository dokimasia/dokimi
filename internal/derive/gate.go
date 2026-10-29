// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive

import (
	"slices"

	"go.dokimi.dev/eidos/plugin/shape"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/spec"
)

// conditions contains the condition of each gate of the spec over the facts
// of a callable. The derivation implements each condition once for every
// language, because each reads only the roles, the forms and the samples that
// the rules of every language assign.
var conditions = map[spec.GateID]func(f *Facts) bool{
	spec.Failure:            func(f *Facts) bool { return f.Signature.Errors != rules.ErrorsNone },
	spec.ReturnsBesideError: func(f *Facts) bool { return f.Signature.Errors == rules.ErrorsLastReturn },
	spec.Cancellation:       func(f *Facts) bool { return hasParam(f, rules.ParamContext) },
	spec.OneValue:           func(f *Facts) bool { return len(f.Results) == 1 },
	spec.TwoValues:          func(f *Facts) bool { return len(f.Results) == 2 },
	spec.NoValue:            func(f *Facts) bool { return len(f.Results) == 0 },
	spec.PresenceFlag:       func(f *Facts) bool { return hasReturn(f, rules.ReturnOkBool) },
	spec.OptionalValue:      func(f *Facts) bool { return firstResultIs(f, symbol.FormOptional) },
	spec.ListValue: func(f *Facts) bool {
		return firstResultIs(f, symbol.FormList) || hasReturn(f, rules.ReturnStream)
	},
	spec.OrderedValue: func(f *Facts) bool {
		return firstResultIs(f, symbol.FormScalar) || firstResultIs(f, symbol.FormText)
	},
	spec.ValueSamples: func(f *Facts) bool {
		return !slices.ContainsFunc(f.Results, func(v Value) bool { return !v.sampled() })
	},
	spec.OptionalInput: func(f *Facts) bool {
		return slices.ContainsFunc(f.Inputs, func(v Value) bool { return v.Shape.Form == symbol.FormOptional })
	},
	spec.ReturnsUnaided: func(f *Facts) bool { return !slices.Contains(f.Class.Mixins, shape.Waits) },
}

// OtherGates are the gates of the spec that another part of the derivation
// judges: the partner resolution judges the partner of a claim, and the
// derivation of an enumeration and of the errors of a package judges their
// facts.
var OtherGates = []spec.GateID{
	spec.PartnerFound, spec.OnePartner, spec.PartnerArguments, spec.Rewritable,
	spec.Members, spec.TwoMembers, spec.ListFunction, spec.ValidFunction, spec.OutsideValue, spec.ParseFunction,
	spec.TextEncoding,
	spec.ErrorValues, spec.TwoErrorValues, spec.ErrorTypes, spec.TwoErrorTypes,
	spec.SelfComparingType, spec.ErrorCause,
}

// Passes reports whether the facts f of a callable pass the gate id, and
// whether the derivation has a condition of the gate over the facts of a
// callable.
func Passes(id spec.GateID, f *Facts) (passes, defined bool) {
	condition, ok := conditions[id]
	if !ok {
		return false, false
	}
	return condition(f), true
}

// hasParam reports whether a parameter of the callable of f has the role.
func hasParam(f *Facts, role rules.ParamRole) bool {
	return slices.ContainsFunc(f.Signature.Params, func(p rules.ParamView) bool { return p.Role == role })
}

// hasReturn reports whether a return of the callable of f has the role.
func hasReturn(f *Facts, role rules.ReturnRole) bool {
	return slices.ContainsFunc(f.Signature.Returns, func(r rules.ReturnView) bool { return r.Role == role })
}

// firstResultIs reports whether the first value result of f has the form.
func firstResultIs(f *Facts, form symbol.TypeForm) bool {
	return len(f.Results) > 0 && f.Results[0].Shape.Form == form
}
