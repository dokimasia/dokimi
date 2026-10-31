// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// Derived is one check that a claim of the conformance spec derives for a
// callable, which the template writes as a check that calls its probe.
//
// Field and Type are the field and the type of the callable, as its
// [Callable] contains them, and Method and Segment are the method and the
// segment of its identity. Claim, Probe, Reason and Law are the names of the
// constants of the package spec of the Go runtime, Machineries those of the
// machineries of the claim, the deepest first, and Entails the identities of
// the shallower checks that the law of the check implies. Args are the
// arguments of the probe after the handle. Defect is the name of the function
// that returns the planted defect, which every check with the same break of
// the same callable shares. Clock reports a probe of the clock, which runs
// inside a testing/synctest bubble.
type Derived struct {
	Field       string
	Type        string
	Method      string
	Segment     string
	Statement   string
	Claim       string
	Probe       string
	Reason      string
	Law         string
	Defect      string
	Machineries []string
	Entails     []Identity
	Args        []Arg
	Shallow     bool
	Clock       bool
}

// Arg is one argument of a probe: a closure that calls the subject, an
// expression, a projection of a field, or a duration.
type Arg struct {
	Closure    *Closure
	Expr       *Expr
	Projection *Projection
	Duration   time.Duration
}

// Projection is a function that returns the field Field of a value of the
// type Type, which has the type Result.
type Projection struct {
	Type   *emit.TypeRef
	Result *emit.TypeRef
	Field  string
}

// Expr is one expression of a derived check: the name of a variable, a
// value, the field Select of another expression, a pair of two expressions,
// or a word of Go such as nil or true. Where Hook is set, the expression is a
// function of the function type Hook that calls the function Name and
// ignores its arguments.
type Expr struct {
	Value  *emit.Value
	Of     *Expr
	Pair   *PairExpr
	Hook   *emit.TypeRef
	Name   string
	Select string
	Word   string
}

// PairExpr is a probe.Pair of two expressions of the types Types.
type PairExpr struct {
	Types  [2]*emit.TypeRef
	Values [2]Expr
}

// Closure is a function that a derived check passes to its probe. It takes
// the context of the call and Param, declares Lets, and calls the callable
// Method on Receiver, or the subject itself where Method is empty, with Args.
// It assigns the results of the call to Assign and returns Results. Without
// Assign it returns the results of the call, and with Discard it makes the
// call a statement of its own.
//
// Returns are the types of its results before the error, which it returns
// where Error is set, and before the flag, which it returns where Flag is
// set. Pair reports that its two types are the two values of a probe.Pair.
// Later are the assignments that it makes on each call after its first.
// Stamps is the field of each element of the result that it returns in place
// of the elements.
type Closure struct {
	Param    *ClosureParam
	Receiver *Expr
	Method   string
	Stamps   string
	Returns  []*emit.TypeRef
	Lets     []Let
	Later    []Let
	Args     []Expr
	Assign   []string
	Results  []Expr
	Discard  bool
	Error    bool
	Flag     bool
	Pair     bool
}

// ClosureParam is the parameter of a closure after the context. Hook and
// Body report the hook func() of a registrar and the body probe.Act of a
// scope, and Type is the type of any other parameter.
type ClosureParam struct {
	Type *emit.TypeRef
	Name string
	Hook bool
	Body bool
}

// Let is a variable that a closure declares before its call, or the
// assignment of its field Select where Select is set.
type Let struct {
	Name   string
	Select string
	Value  Expr
}

// form is the Go form of one input of a probe.
type form int

// The forms of the inputs of the probes, by RFC-0027 of the Go probes.
const (
	// formAct calls the callable for its error.
	formAct form = iota
	// formCall returns the value, or the pair of values, beside the error.
	formCall
	// formLookup returns the value, or the pair of values, beside the flag.
	formLookup
	// formWritten passes a sample to the written input and returns it.
	formWritten
	// formRewritten is formWritten with the first sample on its first call
	// and the second sample on every later call.
	formRewritten
	// formActs calls the callable for its error with the first sample on its
	// first call and the second sample on every later call, so two writes
	// add two elements.
	formActs
	// formTaking takes a value for the first input and returns the error.
	formTaking
	// formInverse takes a value for the first input and returns the value
	// beside the error.
	formInverse
	// formRoundTrip returns the sample of the first input and the result.
	formRoundTrip
	// formHook takes the hook for the first input.
	formHook
	// formScope takes the body for the first input.
	formScope
	// formHandle takes the handle of a call for the first input, or calls
	// the method of the handle that a param names.
	formHandle
	// formDelay passes a sample to the first input and returns it.
	formDelay
	// formStamps returns the stamp of each entry of the callable.
	formStamps
	// formZero is the zero value, or the pair of zero values, of the
	// results of another input.
	formZero
	// formValue is the value of a param.
	formValue
	// formInt is the int of a param.
	formInt
	// formDuration is the duration of a param.
	formDuration
	// formProjection returns the field that a param names of the result of
	// another input.
	formProjection
	// formMatched calls a writer that a predicate guards for its error, with
	// the first sample on its first call and the field of the param field of
	// the second sample on every later call.
	formMatched
	// formMatch calls the predicate of a guarded writer with the value of the
	// later calls of formMatched, and returns its verdict beside the error.
	formMatch
)

// slot is one input of a probe: its name in the entry of the spec, its form,
// and the input whose results a zero or a projection reads.
type slot struct {
	name string
	of   string
	form form
}

// The names of the inputs of the probes that more than one probe takes.
const (
	inCall     = "call"
	inZero     = "zero"
	inSentinel = "sentinel"
	inRead     = "read"
	inWrite    = "write"
	inObserve  = "observe"
	inPut      = "put"
	inDrain    = "drain"
)

// The inputs that many probes share.
var (
	act      = slot{name: inCall, form: formAct}
	call     = slot{name: inCall, form: formCall}
	zero     = slot{name: inZero, of: inCall, form: formZero}
	sentinel = slot{name: inSentinel, form: formValue}
	read     = slot{name: inRead, form: formCall}
	drain    = slot{name: inDrain, form: formCall}
	writes   = slot{name: inWrite, form: formActs}
	put      = slot{name: inPut, form: formAct}
	observe  = slot{name: inObserve, form: formCall}
)

// probes contains the inputs of each probe that a suite calls, in the order
// of the parameters of its Go function after the handle. The probes of the
// enumerations and the errors belong to the generators enum and sentinel.
var probes = map[spec.ProbeID][]slot{
	spec.Survives:              {act},
	spec.UnderCancelled:        {act},
	spec.ZeroBesideError:       {call, zero},
	spec.ZeroPairBesideError:   {call, zero},
	spec.ZeroBesideFalse:       {{name: inCall, form: formLookup}, zero},
	spec.ZeroPairBesideFalse:   {{name: inCall, form: formLookup}, zero},
	spec.NilBesideAMiss:        {call},
	spec.ZeroBesideAMiss:       {call, zero},
	spec.RepeatAgrees:          {call},
	spec.SurvivesRepeat:        {act},
	spec.ReportsSentinel:       {call, sentinel},
	spec.RepeatIsQuiet:         {act},
	spec.NeverDecreases:        {call},
	spec.RepeatAscends:         {call},
	spec.EscapesTheUnsafe:      {{name: inCall, of: "unsafe", form: formCall}, {name: "unsafe", form: formValue}},
	spec.AcceptsTheEdge:        {{name: inCall, form: formTaking}, {name: "edge", form: formValue}},
	spec.RemainsPoisoned:       {act, {name: "induce", form: formAct}},
	spec.SucceedsWithin:        {act, {name: "attempts", form: formInt}},
	spec.RefusesTheSecondWrite: {act, sentinel},
	spec.RefusesPastTheBurst:   {act, {name: "burst", form: formInt}, sentinel},
	spec.RefusesEarly:          {act, sentinel},
	spec.WriteThenRead:         {read, {name: inWrite, form: formWritten}},
	spec.PinnedAcrossWrite:     {read, {name: inWrite, form: formAct}},
	spec.SecondWriteReplaces:   {{name: inCall, form: formRewritten}, read},
	spec.KeepsPartitionsApart:  {{name: inWrite, form: formRewritten}, read},
	spec.RoundTrips:            {{name: inCall, form: formRoundTrip}, {name: "inverse", form: formInverse}},
	spec.WorksAfter:            {act, {name: "ready", form: formAct}},
	spec.ReportsSealed:         {act, {name: "close", form: formAct}, sentinel},
	spec.SealsAfterClose: {
		call, {name: "close", of: inCall, form: formHandle}, {name: "after", of: inCall, form: formHandle}, sentinel,
	},
	spec.MovesObservation:    {act, observe},
	spec.EachCallCompounds:   {act, observe},
	spec.SettlesThenObserves: {act, {name: "settle", form: formAct}, observe},
	spec.DetectsTamper:       {act, {name: "verify", form: formAct}, {name: "tamper", form: formAct}},
	spec.ConvergesAfter:      {act, {name: "begin", form: formAct}, {name: "attempts", form: formInt}},
	spec.RunsTheHook:         {act, {name: "register", form: formHook}},
	spec.RefusesTheInvalid: {
		{name: inCall, form: formTaking}, {name: "validate", form: formTaking}, {name: "invalid", form: formValue},
	},
	spec.RefusesWhatThePredicateRejects: {{name: inCall, form: formMatched}, {name: "match", form: formMatch}},
	spec.ResignsLeadership:              {act, {name: "resign", form: formAct}, {name: "leading", form: formCall}},
	spec.RollsBackOnFailure: {
		{name: inCall, form: formScope}, {name: inWrite, form: formAct}, {name: inRead, form: formAct}, sentinel,
	},
	spec.AgreesWithTheBacking: {call, {name: "backing", form: formCall}, put},
	spec.TakesBackWhatItLent:  {call, {name: inPut, form: formTaking}},
	spec.ContinuesFromTheCursor: {
		{name: inCall, form: formInverse},
		{name: "cursor", of: inCall, form: formProjection},
	},
	spec.CompensatesTheStep: {act, {name: "compensate", form: formAct}},
	spec.RefusedLeavesNothing: {
		{name: inCall, form: formTaking},
		{name: "reader", form: formCall},
		{name: "refused", form: formValue},
		{name: inZero, of: "reader", form: formZero},
	},
	spec.RemovesOnDelete:   {act, put, {name: inRead, form: formAct}, sentinel},
	spec.NoElementTwice:    {drain, writes},
	spec.DrainsInOneOrder:  {drain, writes},
	spec.GainsWhatWentIn:   {drain, writes},
	spec.KeepsWhatWentIn:   {drain, writes},
	spec.ReadsNeverRegress: {read, {name: inWrite, form: formAct}, {name: "version", of: inRead, form: formProjection}},
	spec.WritesAdvance: {
		read, {name: inWrite, form: formCall}, {name: "version", of: inRead, form: formProjection},
	},
	spec.ReplaysInStampOrder:      {{name: inCall, form: formActs}, {name: "replay", form: formStamps}},
	spec.ReplayGrows:              {act, {name: "replay", form: formCall}},
	spec.ReturnsWithin:            {act, {name: "budget", form: formDuration}},
	spec.CountsWithinWindow:       {act, {name: "count", form: formCall}, {name: "window", form: formDuration}},
	spec.RunsOnSchedule:           {{name: inCall, form: formDelay}, {name: "count", form: formCall}},
	spec.MovesWithTheClock:        {call, put},
	spec.SteadyWhileTheClockStops: {call, put},
	spec.ExpiresAfter: {
		read, {name: inWrite, form: formAct}, {name: "lifetime", form: formDuration}, sentinel,
	},
}

// clockProbes are the probes that advance the clock of the handle.
var clockProbes = []spec.ProbeID{
	spec.ReturnsWithin, spec.CountsWithinWindow, spec.RunsOnSchedule, spec.MovesWithTheClock,
	spec.SteadyWhileTheClockStops, spec.ExpiresAfter,
}

// goName returns the name of the constant of id in the package spec of the
// Go runtime, as the naming table of the spec spells it.
func goName[T ~string](id T) string { return naming.Pascal(string(id)) }

// identityOf returns the identity of the derived check c among the callables
// of s: its field, its method and its segment. A claim that binds once for
// each instance of a contract has the id of the instance after its segment.
func identityOf(s *Suite, c *derive.Check) (Callable, string, string) {
	segment := c.Identity[strings.Index(c.Identity, "/")+1:]
	method := naming.Pascal(strings.ReplaceAll(segment, "#", "-"))
	i := slices.IndexFunc(s.Facts, func(f derive.Facts) bool { return f.Callable == c.Facts.Callable })
	return s.Callables[i], method, segment
}

// derived records into s the derived check of each check of the derivation of
// s, with its planted defect.
func (k *caller) derived(s *Suite) {
	for i := range s.Derivation.Checks {
		c := &s.Derivation.Checks[i]
		callable, method, segment := identityOf(s, c)
		d := Derived{
			Field: callable.Field, Type: callable.Type, Method: method, Segment: segment, Statement: c.Statement,
			Claim: goName(c.Claim.ID), Probe: goName(c.Claim.Probe), Reason: goName(c.Claim.Defect.Reason),
			Law: goName(c.Claim.Law), Shallow: c.Shallow, Clock: slices.Contains(clockProbes, c.Claim.Probe),
		}
		for _, m := range c.Claim.Machineries {
			d.Machineries = append(d.Machineries, goName(m))
		}
		for _, e := range s.Derivation.Edges {
			if e.Deeper != c.Identity {
				continue
			}
			j := slices.IndexFunc(s.Derivation.Checks, func(o derive.Check) bool { return o.Identity == e.Shallower })
			shallow, shallowMethod, _ := identityOf(s, &s.Derivation.Checks[j])
			d.Entails = append(d.Entails, Identity{Field: shallow.Field, Method: shallowMethod, ID: e.Shallower})
		}
		b := binder{k: k, check: c}
		for _, sl := range probes[c.Claim.Probe] {
			d.Args = append(d.Args, b.arg(&sl))
		}
		defect, known := k.defect(s, c, &b)
		if defect.Body == bodyLiveHandle && defect.Handle == nil {
			s.Derivation.Refusals = append(s.Derivation.Refusals, derive.Refusal{
				Claim:    c.Claim.ID,
				Callable: c.Facts.Callable,
				Key:      c.Key,
				Why:      "the handle that " + c.Facts.Callable + " returns is not an interface, so no planted handle can replace it",
				Remedy:   "return the handle of " + c.Facts.Callable + " as an interface",
			})
			continue
		}
		if !k.referableArgs(d.Args) || !k.referableDefect(defect) {
			s.Derivation.Refusals = append(s.Derivation.Refusals, derive.Refusal{
				Claim: c.Claim.ID, Callable: c.Facts.Callable, Key: c.Key,
				Why:    "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy: "author a sample of each input of " + c.Facts.Callable + " with the directive sample",
			})
			continue
		}
		if !known {
			s.Defects = append(s.Defects, defect)
		}
		d.Defect = defect.Name
		s.Derived = append(s.Derived, d)
	}
}

// referableArgs reports whether the file of the suite can refer to each
// sample that a closure of args passes or declares. The other values of a
// check are referable wherever its declaration is: a zero value is a literal
// of a type of the callable, and a param names a declaration in the scope of
// the declaration. A read selects a field of the value that its write passes
// or declares, so the sample of the write covers the field.
func (k *caller) referableArgs(args []Arg) bool {
	var exprs []Expr
	for i := range args {
		if c := args[i].Closure; c != nil {
			exprs = append(exprs, c.Args...)
			for _, l := range slices.Concat(c.Lets, c.Later) {
				exprs = append(exprs, l.Value)
			}
		}
	}
	return !slices.ContainsFunc(exprs, func(e Expr) bool { return !k.referableExpr(&e) })
}

// referableExpr reports whether the file of the suite can refer to the value
// of the expression e, and reports true for an expression without one.
func (k *caller) referableExpr(e *Expr) bool { return e.Value == nil || k.referable(e.Value) }

// referableDefect reports whether the file of the suite can refer to the
// sample of each value that the planted defect d returns. The zero value of a
// type is a literal of the type, which the file can refer to with the type.
func (k *caller) referableDefect(d *Defect) bool {
	return !slices.ContainsFunc(d.Results, func(o Outcome) bool {
		return o.Sample != nil && !k.referableExpr(o.Sample)
	})
}

// binder binds the inputs of a derived check to the arguments of its probe.
type binder struct {
	k     *caller
	check *derive.Check
}

// binding returns the binding of the input name of the claim of the check.
func (b *binder) binding(name string) *spec.Binding {
	i := slices.IndexFunc(b.check.Claim.Inputs, func(in spec.Binding) bool { return in.Input == name })
	return &b.check.Claim.Inputs[i]
}

// facts returns the facts of the callable of the input name: the callable of
// the check, or the partner that the binding names.
func (b *binder) facts(name string) (*derive.Facts, *derive.Pair) {
	in := b.binding(name)
	if in.Kind == spec.FromPartner {
		pair := b.check.Pairs[in.Partner]
		return pair.Partner, pair
	}
	return b.check.Facts, b.ownPair()
}

// ownPair returns the pair of the check in which the callable of the check
// reads or writes, and nil where it does neither. The partners of a claim are
// visited in the order of the claim, so one pair is chosen for every run.
func (b *binder) ownPair() *derive.Pair {
	for _, p := range b.check.Claim.Partners {
		pair := b.check.Pairs[p.Name]
		if same(pair.Read, b.check.Facts) || same(pair.Write, b.check.Facts) {
			return pair
		}
	}
	return nil
}

// same reports whether a is the facts of the callable of b. a is nil for the
// read and the write of a pair where neither writes.
func same(a, b *derive.Facts) bool { return a != nil && a.ID == b.ID }

// param returns the value of the param key of the classifications of the
// callable of the check, and nil where none sets it.
func (b *binder) param(key string) any {
	for _, params := range b.check.Facts.Class.Params {
		if v, ok := params[directive.ParamKey(key)]; ok {
			return v
		}
	}
	return nil
}

// arg returns the argument of the probe for the input sl.
func (b *binder) arg(sl *slot) Arg {
	switch sl.form {
	case formZero:
		f, _ := b.facts(sl.of)
		return Arg{Expr: b.k.zeroOf(f)}
	case formValue:
		return Arg{Expr: b.k.valueOfParam(b.param(b.binding(sl.name).Param))}
	case formInt:
		return Arg{Expr: &Expr{Word: intOf(b.param(b.binding(sl.name).Param))}}
	case formDuration:
		d, _ := time.ParseDuration(nameOf(b.param(b.binding(sl.name).Param)))
		return Arg{Duration: d}
	case formProjection:
		f, _ := b.facts(sl.of)
		field := nameOf(b.param(b.binding(sl.name).Param))
		ref := b.k.resultRef(f, 0)
		return Arg{Projection: &Projection{Type: b.k.typeOf(ref), Result: b.k.fieldRef(ref, field), Field: field}}
	default:
	}
	in := &input{sl: sl}
	switch sl.form {
	case formHandle:
		call, _ := b.facts(sl.of)
		in.handle = b.k.resultRef(call, 0)
		if bound := b.binding(sl.name); bound.Kind == spec.FromParam {
			in.method = nameOf(b.param(bound.Param))
			return Arg{Closure: b.k.closure(in)}
		}
	case formStamps:
		// The stamp of an entry is the field that the param version of the
		// mixin causal of the callable names.
		in.stamps = nameOf(b.param(stampParam))
	case formCall:
		if sl.of != "" {
			in.with = b.k.valueOfParam(b.param(b.binding(sl.of).Param))
		}
	case formMatched, formMatch:
		in.changed = nameOf(b.param(fieldParam))
	default:
	}
	in.f, in.pair = b.facts(sl.name)
	return Arg{Closure: b.k.closure(in)}
}

// stampParam is the param of the mixin causal that names the field of the
// stamp of an entry, and fieldParam the param of the contract if-match that
// names the field that the predicate judges.
const (
	stampParam = "version"
	fieldParam = "field"
)

// nameOf returns the name of the declaration that a param refers to, such as
// the field of version=Rev, or the text of a param that is a string. It
// returns the empty string for a param that is not set.
func nameOf(v any) string {
	switch p := v.(type) {
	case symbol.Identity:
		return p.Name
	case string:
		return p
	}
	return ""
}

// fieldRef returns the type of the field name of the struct that ref refers
// to, at the witnesses of the declaration.
func (k *caller) fieldRef(ref *node.TypeRef, name string) *emit.TypeRef {
	if ref.Form == symbol.FormOptional {
		ref = ref.Elems[0]
	}
	sym, _ := k.m.Reader().Lookup(ref.Target)
	st, _ := sym.(*node.Struct)
	// The derivation binds a field param only to a field of the struct.
	at := slices.IndexFunc(st.Fields, func(f *node.Field) bool { return f.Name == name })
	return k.typeOf(st.Fields[at].Type)
}

// resultRef returns the type of the i-th result in the role value of the
// callable of f. Each value of a Go callable precedes its flag and its
// error, so the i-th value is its i-th result.
func (k *caller) resultRef(f *derive.Facts, i int) *node.TypeRef {
	return k.instance(f.Signature.Returns[i].Ref)
}

// zeroOf returns the zero value of the results of the callable of f: the
// zero of its one value, or the pair of the zeros of its two values.
func (k *caller) zeroOf(f *derive.Facts) *Expr {
	if len(f.Results) == 1 {
		return &Expr{Value: typed(k.typeOf(k.resultRef(f, 0)), &f.Results[0].Zero)}
	}
	p := &PairExpr{}
	for i := range 2 {
		p.Types[i] = k.typeOf(k.resultRef(f, i))
		p.Values[i] = Expr{Value: typed(p.Types[i], &f.Results[i].Zero)}
	}
	return &Expr{Pair: p}
}

// typed returns v as a value of the type ref: the conversion of a literal or
// a name, whose constant can be untyped, and any other value as it is, which
// already has its type.
func typed(ref *emit.TypeRef, v *emit.Value) *emit.Value {
	if v.Kind == emit.ValueLiteral || v.Kind == emit.ValueName {
		return new(emit.Conversion(ref, *v))
	}
	return v
}

// valueOfParam returns the expression of the value of a param, which names a
// declaration of the package, such as sentinel=ErrMissing.
func (*caller) valueOfParam(v any) *Expr {
	id, _ := v.(symbol.Identity)
	return &Expr{Value: new(emit.Name(id))}
}

// intOf returns the decimal text of an int param, which the classification
// gives as an int64 or as its text.
func intOf(v any) string {
	if n, ok := v.(int64); ok {
		return strconv.FormatInt(n, 10)
	}
	return nameOf(v)
}
