// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive

import (
	"reflect"
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
	spec.ReturnsBeforeCancel: func(f *Facts) bool {
		return !slices.Contains(f.Class.Mixins, shape.RunsUntilCancelled)
	},
	spec.TwoInputsOfOneType: func(f *Facts) bool {
		return len(f.Inputs) == 2 && reflect.DeepEqual(f.Inputs[0].Shape, f.Inputs[1].Shape)
	},
	spec.ResultOfTheInputType: func(f *Facts) bool {
		return len(f.Results) == 1 && len(f.Inputs) > 0 && reflect.DeepEqual(f.Results[0].Shape, f.Inputs[0].Shape)
	},
	spec.ContextOrInput: func(f *Facts) bool { return hasParam(f, rules.ParamContext) || len(f.Inputs) > 0 },
	spec.ExactFidelity: func(f *Facts) bool {
		for _, params := range f.Class.Params {
			if params[shape.CodecFidelity] == lossy {
				return false
			}
		}
		return true
	},
	spec.ReturnsAValue: func(f *Facts) bool {
		return len(f.Results) > 0 && !slices.ContainsFunc(f.Class.Mixins, func(m shape.Mixin) bool {
			return m == shape.Waits || m == shape.RunsUntilCancelled
		})
	},
}

// lossy is the fidelity of a codec whose round trip drops a part of a value.
const lossy = "lossy"

// OtherGates are the gates of the spec that another part of dokimi judges.
// The binding of a claim judges the samples of the inputs of the callable, the
// params that the claim reads and the values that its partners return. The
// partner resolution judges the partner of a claim. The derivation of an
// enumeration and of the errors of a package judges their facts. The suite
// generator of a language judges the code that it writes for a derived check,
// and the types of the partners and of the variables that the params of a
// classification name. [ReferenceOf] judges the classifications, the keys and
// the results of the callables of a declaration under the directive fake.
var OtherGates = []spec.GateID{
	spec.InputSamples, spec.ClassificationParams, spec.InterfaceHandle, spec.ImportableSamples,
	spec.PartnerFound, spec.OnePartner, spec.PartnerArguments, spec.Rewritable, spec.WriteReturnsTheValue,
	spec.Members, spec.TwoMembers, spec.ListFunction, spec.ValidFunction, spec.OutsideValue, spec.ParseFunction,
	spec.TextEncoding,
	spec.ErrorValues, spec.TwoErrorValues, spec.ErrorTypes, spec.TwoErrorTypes,
	spec.SelfComparingType, spec.ErrorCause,
	spec.CountableUse, spec.NumericTotal, spec.TakesAReplica, spec.GraphOfStates, spec.DeliversOnAChannel,
	spec.AtLeastOnceMode, spec.StreamMutations, spec.SubscriptionMethods, spec.ReferenceRole,
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

// returnsTheValue reports whether the callable of partner returns first a
// value of the type that the callable of f returns first. A claim with the
// gate write-returns-the-value has the gate one-value, so f returns a value.
func returnsTheValue(partner, f *Facts) bool {
	return len(partner.Results) > 0 && reflect.DeepEqual(partner.Results[0].Shape, f.Results[0].Shape)
}

// firstResultIs reports whether the first value result of f has the form.
func firstResultIs(f *Facts, form symbol.TypeForm) bool {
	return len(f.Results) > 0 && f.Results[0].Shape.Form == form
}
