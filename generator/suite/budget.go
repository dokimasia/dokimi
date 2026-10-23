// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"fmt"
	"time"

	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// BudgetDirective states the budget of a callable: the most heap allocations
// and bytes of one call, the highest p99 and mean latency of one call, and the
// warm-up of its benchmarks. On a method of a declaration under
// [SuiteDirective], and on a function under it, the budget applies to every
// subject of the declaration. On a subject, its positional param
// [CallableParam] contains the callable, such as Store.Get, and the budget
// applies to that subject alone. AllocsParam, BytesParam, LatencyParam and
// MeanParam are the params of the four ceilings. WarmupParam is the param of
// the number of iterations that each benchmark of the callable runs before it
// measures.
const (
	BudgetDirective directive.Name     = "budget"
	AllocsParam     directive.ParamKey = "allocs"
	BytesParam      directive.ParamKey = "bytes"
	LatencyParam    directive.ParamKey = "p99"
	MeanParam       directive.ParamKey = "mean"
	WarmupParam     directive.ParamKey = "warmup"
)

// KeyAllocs, KeyBytes, KeyLatency, KeyMean and KeyWarmup are the keys of the
// facts that the annotate half stamps on a callable under [BudgetDirective],
// one for each ceiling and for the warm-up that its budget states. The facts
// contain the most heap allocations and the most heap bytes of one call, the
// highest p99 and mean latency of one call in nanoseconds, and the number of
// iterations of the warm-up.
const (
	KeyAllocs  meta.KeyName = "suite.allocs"
	KeyBytes   meta.KeyName = "suite.bytes"
	KeyLatency meta.KeyName = "suite.p99"
	KeyMean    meta.KeyName = "suite.mean"
	KeyWarmup  meta.KeyName = "suite.warmup"
)

// The segments of the identities of the checks of the four ceilings, as the
// runtime derives the checks from a budget.
const (
	allocationSegment  = "allocation-budget"
	byteSegment        = "byte-budget"
	latencySegment     = "latency-budget"
	meanLatencySegment = "mean-latency-budget"
)

// ceilingSegments are the segments of the checks of the ceilings, which no
// check of a consumer may have.
var ceilingSegments = []string{allocationSegment, byteSegment, latencySegment, meanLatencySegment}

// Budget is the budget of one call of a callable. A nil Allocs or Bytes, and a
// zero Latency or Mean, state no ceiling. Warmup is the number of iterations
// that each benchmark of the callable runs before it measures, and nil for the
// warm-up of the runtime, each call of the benchmark once. Warmup states no
// ceiling.
type Budget struct {
	Allocs  *int64
	Bytes   *int64
	Warmup  *int64
	Latency time.Duration
	Mean    time.Duration
}

// Ceiling is the budget that a subject states for one callable of its
// declaration. Field is the field of the callable in the names of the
// callables, as Get.
type Ceiling struct {
	Field  string
	Budget Budget
}

// budgetSchema is the directive of a budget. A subject states a budget for
// each of its callables, so the directive repeats.
var budgetSchema = directive.Schema{
	Name:       BudgetDirective,
	Repeatable: true,
	Kinds:      []symbol.Kind{symbol.KindMethod, symbol.KindFunction, symbol.KindStruct},
	Positional: []directive.ParamSpec{{
		Key: CallableParam, Type: directive.TypeString,
		Doc: "the callable of the budget of a subject, such as Store.Get. The budget of a callable names none",
	}},
	Params: []directive.ParamSpec{
		{Key: AllocsParam, Type: directive.TypeInt, Doc: "the most heap allocations of one call"},
		{Key: BytesParam, Type: directive.TypeInt, Doc: "the most heap bytes that one call allocates"},
		{
			Key: LatencyParam, Type: directive.TypeString,
			Doc: "the highest p99 latency of one call, as a duration such as 2us",
		},
		{
			Key: MeanParam, Type: directive.TypeString,
			Doc: "the highest mean latency of one call, as a duration such as 1us. A budget with a mean states a p99",
		},
		{
			Key: WarmupParam, Type: directive.TypeInt,
			Doc: "the number of iterations that each benchmark of the callable runs before it measures, each call " +
				"of the benchmark once without the param",
		},
	},
	Doc: "states the ceilings of the allocations, the bytes and the latency of one call of a callable, each a " +
		"check with a planted defect, and the warm-up of the benchmarks of the callable",
}

// budgetOfMethod stamps the facts of the budget of a method of a declaration
// under the directive suite. It reports [BadBudget] at the budget of a method
// of another declaration and of a method that is not exported, because no
// suite has such a callable.
func (g *generator) budgetOfMethod(m *sdk.MethodMatch, _ part, st *sdk.Stamper) error {
	method := m.Method
	name := method.ID.Owner + "." + method.Name
	switch {
	case !suited(m, &method.Host):
		m.ErrorfAt(BadBudget, m.Directive().Pos, "%s has no directive suite, so the budget of %s has no check",
			method.ID.Owner, name)
	case method.Visibility != symbol.VisibilityPublic:
		m.ErrorfAt(BadBudget, m.Directive().Pos, "%s is not exported, so its budget has no check", name)
	default:
		g.stamp(m, name, st)
	}
	return nil
}

// budgetOfFunction stamps the facts of the budget of a function under the
// directive suite. It reports [BadBudget] at the budget of another function.
func (g *generator) budgetOfFunction(m *sdk.FunctionMatch, _ part, st *sdk.Stamper) error {
	fn := m.Function
	if !suited(m, &fn.ID) {
		m.ErrorfAt(BadBudget, m.Directive().Pos, "%s has no directive suite, so the budget of %s has no check",
			fn.Name, fn.Name)
		return nil
	}
	g.stamp(m, fn.Name, st)
	return nil
}

// budgetOfStruct reports [BadBudget] at the budget of a struct without the
// directive subject. The conformance test of a subject reads the budgets of
// its struct, and no other struct has a test that reads them.
func (*generator) budgetOfStruct(m *sdk.StructMatch, _ part, _ *sdk.Emitter) error {
	if instances(m, SubjectDirective) == 0 {
		m.ErrorfAt(BadBudget, m.Directive().Pos, "%s has no directive subject, so its budget has no check",
			m.Struct.Name)
	}
	return nil
}

// stamp stamps a fact of each ceiling and of the warm-up of the budget of the
// callable name on the subject of m. A callable has one budget, and its budget
// names no callable. stamp reports [BadBudget] at the budget otherwise, and at
// each mistake that [budgetOf] reports, and stamps nothing then.
func (g *generator) stamp(m sdk.SubjectMatch, name string, st *sdk.Stamper) {
	d := m.Directive()
	if n := instances(m, BudgetDirective); n > 1 {
		m.ErrorfAt(BadBudget, d.Pos, "%s states %d budgets, and a callable has one budget", name, n)
		return
	}
	if len(d.Args) > 0 {
		m.ErrorfAt(BadBudget, d.Pos, "the budget of %s names the callable %s, and only the budget of a subject "+
			"names its callable", name, d.Args[0].Str)
		return
	}
	b, valid := budgetOf(m, d, name)
	if !valid {
		return
	}
	if b.Allocs != nil {
		sdk.Stamp(st, g.allocs, *b.Allocs)
	}
	if b.Bytes != nil {
		sdk.Stamp(st, g.bytes, *b.Bytes)
	}
	if b.Warmup != nil {
		sdk.Stamp(st, g.warmup, *b.Warmup)
	}
	if b.Latency != 0 {
		sdk.Stamp(st, g.latency, int64(b.Latency))
	}
	if b.Mean != 0 {
		sdk.Stamp(st, g.mean, int64(b.Mean))
	}
}

// instances returns the number of the directives named name on the subject of
// m.
func instances(m sdk.SubjectMatch, name directive.Name) int {
	n := 0
	for range m.Directives(name) {
		n++
	}
	return n
}

// budgetOf returns the budget that the directive d states for the callable
// name, and whether the budget is valid. A count of allocations, bytes or
// iterations of the warm-up is 0 or more, and a p99 and a mean are durations
// above 0. A budget with a mean states a p99, because a mean alone passes a
// tail that grows, and a budget states a ceiling or a warm-up. budgetOf
// reports [BadBudget] at d for each mistake. eidos refuses a count that is no
// integer before suite reads the directive.
func budgetOf(m sdk.SubjectMatch, d *directive.Directive, name string) (Budget, bool) {
	r := reading{m: m, d: d, name: name, valid: true}
	b := Budget{
		Allocs:  r.count(AllocsParam, "a ceiling of allocations is 0 or more"),
		Bytes:   r.count(BytesParam, "a ceiling of bytes is 0 or more"),
		Warmup:  r.count(WarmupParam, "a warm-up is 0 iterations or more"),
		Latency: r.duration(LatencyParam),
		Mean:    r.duration(MeanParam),
	}
	_, latency := d.Param(LatencyParam)
	_, mean := d.Param(MeanParam)
	_, allocs := d.Param(AllocsParam)
	_, bytes := d.Param(BytesParam)
	_, warmup := d.Param(WarmupParam)
	switch {
	case mean && !latency:
		r.refuse("states a mean and no p99, and a mean alone passes a tail that grows")
	case !allocs && !bytes && !latency && !warmup:
		r.refuse("states no ceiling and no warm-up. State allocs, bytes, p99, mean or warmup")
	}
	return b, r.valid
}

// reading is the read of the budget of one directive d of the callable name.
// valid reports that the read found no mistake, and m reports each mistake.
type reading struct {
	m     sdk.SubjectMatch
	d     *directive.Directive
	name  string
	valid bool
}

// refuse reports [BadBudget] at the directive with the reason that format and
// args state, and marks the read invalid.
func (r *reading) refuse(format string, args ...any) {
	r.m.ErrorfAt(BadBudget, r.d.Pos, "the budget of %s %s", r.name, fmt.Sprintf(format, args...))
	r.valid = false
}

// count returns the count of the param key, and nil where the directive does
// not state the param. It refuses a count below 0 with the rule of the param,
// such as "a ceiling of bytes is 0 or more".
func (r *reading) count(key directive.ParamKey, rule string) *int64 {
	v, stated := r.d.Param(key)
	switch {
	case !stated:
		return nil
	case v.Int < 0:
		r.refuse("states %s=%d, and %s", key, v.Int, rule)
		return nil
	}
	return new(v.Int)
}

// duration returns the duration of the param key, and 0 where the directive
// does not state the param. It refuses a value that is no duration, and a
// duration of 0 or less.
func (r *reading) duration(key directive.ParamKey) time.Duration {
	v, stated := r.d.Param(key)
	if !stated {
		return 0
	}
	d, err := time.ParseDuration(v.Str)
	switch {
	case err != nil:
		r.refuse("states %s=%q, and a %s is a duration such as 2us", key, v.Str, key)
	case d <= 0:
		r.refuse("states %s=%q, and a %s is above 0", key, v.Str, key)
	default:
		return d
	}
	return 0
}

// budgetFact returns the budget that the facts of the callable id contain.
func (g *generator) budgetFact(m sdk.Matcher, id *symbol.Identity) Budget {
	var b Budget
	if n, stated := sdk.FactOf(m, *id, g.allocs); stated {
		b.Allocs = new(n)
	}
	if n, stated := sdk.FactOf(m, *id, g.bytes); stated {
		b.Bytes = new(n)
	}
	if n, stated := sdk.FactOf(m, *id, g.warmup); stated {
		b.Warmup = new(n)
	}
	latency, _ := sdk.FactOf(m, *id, g.latency)
	mean, _ := sdk.FactOf(m, *id, g.mean)
	b.Latency, b.Mean = time.Duration(latency), time.Duration(mean)
	return b
}

// segments returns the segment of the check of each ceiling that the budget
// states, in the order of its fields.
func (b Budget) segments() []string {
	var out []string
	if b.Allocs != nil {
		out = append(out, allocationSegment)
	}
	if b.Bytes != nil {
		out = append(out, byteSegment)
	}
	if b.Latency != 0 {
		out = append(out, latencySegment)
	}
	if b.Mean != 0 {
		out = append(out, meanLatencySegment)
	}
	return out
}

// ceilingsOf returns the budget that the subject of m states for each callable
// of the declaration declared, in the order of the directives. methods are
// the callables of the declaration. A budget of a subject names one callable
// of the declaration, as in Store.Get, whose declaration states no budget,
// and a subject states one budget of a callable. ceilingsOf reports
// [BadBudget] at each budget otherwise, and at each mistake of the ceilings
// that [budgetOf] reports. It returns false where it reports a budget. A
// subject of a struct with several directives subject reads the budgets of
// the callables of its own declaration, as [skips] states.
func (g *generator) ceilingsOf(
	m *sdk.StructMatch, declared *symbol.Identity, methods []*node.Method,
) ([]Ceiling, bool) {
	subject := m.Struct.Name
	names := make([]string, 0, len(methods))
	for _, method := range methods {
		names = append(names, method.Name)
	}
	written := map[string]int{}
	for d := range m.Directives(BudgetDirective) {
		if len(d.Args) > 0 {
			written[d.Args[0].Str]++
		}
	}
	var out []Ceiling
	stated := true
	for d := range m.Directives(BudgetDirective) {
		named := ""
		if len(d.Args) > 0 {
			named = d.Args[0].Str
		}
		if skipped, refused := skips(m, named); skipped {
			stated = stated && !refused
			continue
		}
		if len(d.Args) == 0 {
			m.ErrorfAt(BadBudget, d.Pos, "the budget of %s names no callable of %s, so it has no check%s",
				subject, declared.Name, callablesOf(declared.Name, names))
			stated = false
			continue
		}
		callable := d.Args[0].Str
		at := -1
		for i, name := range names {
			if declared.Name+"."+name == callable {
				at = i
			}
		}
		var ofDeclaration Budget
		if at >= 0 {
			ofDeclaration = g.budgetFact(m, &methods[at].ID)
		}
		switch {
		case at < 0:
			m.ErrorfAt(BadBudget, d.Pos, "%s is no callable of %s, so the budget of %s has no check%s",
				callable, declared.Name, subject, callablesOf(declared.Name, names))
			stated = false
		case written[callable] > 1:
			m.ErrorfAt(BadBudget, d.Pos, "%s states %d budgets of %s, and a callable has one budget",
				subject, written[callable], callable)
			stated = false
		case ofDeclaration != (Budget{}):
			advice := ""
			if len(ofDeclaration.segments()) > 0 {
				advice = ". Excuse the checks of the declaration instead"
			}
			m.ErrorfAt(BadBudget, d.Pos, "%s states a budget of %s, whose declaration states a budget%s",
				subject, callable, advice)
			stated = false
		default:
			b, valid := budgetOf(m, d, callable)
			if !valid {
				stated = false
				continue
			}
			out = append(out, Ceiling{Field: names[at], Budget: b})
		}
	}
	return out, stated
}
