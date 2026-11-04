// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"slices"

	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"
)

// consumed are the machineries that a consumer supplies through a door of the
// application, in the order of the spec. The param machineries of the
// directive suite names each by its param, such as faults.
var consumed = func() []spec.Machinery {
	var out []spec.Machinery
	for _, m := range spec.Machineries() {
		if m.By == spec.ByConsumer {
			out = append(out, m)
		}
	}
	return out
}()

// consumedParams are the params of the machineries of consumed, in their
// order.
var consumedParams = func() []string {
	out := make([]string, 0, len(consumed))
	for _, m := range consumed {
		out = append(out, m.Param)
	}
	return out
}()

// machineriesOfInterface stamps the machineries of an interface under the
// directive suite, as [generator.machineriesOf] states.
func (g *generator) machineriesOfInterface(m *sdk.InterfaceMatch, _ part, st *sdk.Stamper) error {
	g.machineriesOf(m, m.Interface.Name, st)
	return nil
}

// machineriesOfStruct stamps the machineries of a struct under the directive
// suite, as [generator.machineriesOf] states.
func (g *generator) machineriesOfStruct(m *sdk.StructMatch, _ part, st *sdk.Stamper) error {
	g.machineriesOf(m, m.Struct.Name, st)
	return nil
}

// machineriesOfFunction stamps the machineries of a function under the
// directive suite, as [generator.machineriesOf] states.
func (g *generator) machineriesOfFunction(m *sdk.FunctionMatch, _ part, st *sdk.Stamper) error {
	g.machineriesOf(m, m.Function.Name, st)
	return nil
}

// machineriesOf stamps the identity of each machinery that the param
// machineries of the directive suite of m names on the declaration named
// decl, under [KeyMachineries]. It reports [UnknownMachinery] at the directive
// for each value that names no machinery of a consumer, and stamps the other
// values. It stamps nothing for a directive without such a value, such as a
// directive without the param.
func (g *generator) machineriesOf(m sdk.SubjectMatch, decl string, st *sdk.Stamper) {
	d := m.Directive()
	param, _ := d.Param(MachineriesParam)
	var out []string
	for i := range param.List {
		value := param.List[i].Str
		at := slices.IndexFunc(consumed, func(c spec.Machinery) bool { return c.Param == value })
		if at < 0 {
			m.ErrorfAt(UnknownMachinery, d.Pos, "no consumer supplies the machinery %s, so the suite of %s derives no "+
				"claim of it. The values of the param machineries are %s", value, decl, listed(consumedParams, "and"))
			continue
		}
		out = append(out, string(consumed[at].ID))
	}
	if len(out) > 0 {
		sdk.Stamp(st, g.machineries, out)
	}
}

// machineriesOfDeclaration returns the machineries that the declaration id
// states through the param machineries of its directive suite.
func (g *generator) machineriesOfDeclaration(m sdk.Matcher, id *symbol.Identity) []spec.MachineryID {
	stated, _ := sdk.FactOf(m, *id, g.machineries)
	out := make([]spec.MachineryID, 0, len(stated))
	for _, s := range stated {
		out = append(out, spec.MachineryID(s))
	}
	return out
}
