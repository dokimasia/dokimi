// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// Defect is the planted defect of the derived checks with one break of one
// callable. The template writes it as a function that wraps a fresh subject
// of the application: a type that embeds the subject type and overrides the
// callable for an interface or a struct, and a function for a function.
//
// Name is the name of the function, and Type the name of the type, empty for
// a function. Doc is the clause of the planted defect. Body is the body of the
// override, which [bodies] names for each defect form. Method is the
// signature of the callable, with the name p and its position for each
// parameter and r and its position for each result, and Args are the names of
// the parameters in a call. Results are the results of the callable.
//
// Version is the field of the version of a value, Budget the time that a
// slow callable waits, twice its budget, and Input the position of the first
// input. Handle is the handle that a callable of the form opens-a-live-handle
// returns in place of its own, and Handled the type of the handle that it
// wraps. Delegates reports a body that calls the callable of the fresh
// subject, and Counts a body that counts its calls. Context is the parameter
// of the context of the callable, and empty for a callable without one.
// Guards reports a body that returns the error of a cancelled context first,
// so the defect reports a cancelled context as the callable does. A panic, a
// body that breaks cancellation and a body that reads the error of the fresh
// subject have no guard.
type Defect struct {
	Method    *emit.Method
	Handle    *Defect
	Handled   *emit.TypeRef
	Name      string
	Type      string
	Doc       string
	Body      string
	Version   string
	Context   string
	Args      []string
	Results   []Outcome
	Budget    time.Duration
	Input     int
	Delegates bool
	Counts    bool
	Guards    bool
}

// Outcome is one result of the callable of a planted defect: its name, its
// role and a sample that is not its zero value. Role is value, flag or error.
// A body returns the zero value of a result through a bare return.
type Outcome struct {
	Sample *Expr
	Name   string
	Role   string
}

// The roles of the results of the callable of a planted defect.
const (
	roleValue = "value"
	roleFlag  = "flag"
	roleError = "error"
)

// The bodies of the overrides of the planted defects, which the template of a
// defect names.
const (
	bodyPanic          = "panic"
	bodyPanicLater     = "panic-later"
	bodyZero           = "zero"
	bodyZeroLater      = "zero-later"
	bodyFail           = "fail"
	bodyFailSuccess    = "fail-success"
	bodyFailLater      = "fail-later"
	bodyBesideError    = "beside-error"
	bodyUncancelled    = "uncancelled"
	bodyBesideFalse    = "beside-false"
	bodySample         = "sample"
	bodySampleSuccess  = "sample-success"
	bodyDrift          = "drift"
	bodyDescend        = "descend"
	bodyVersionDescend = "version-descend"
	bodyVersionZero    = "version-zero"
	bodyDrainTwice     = "drain-twice"
	bodyReorder        = "reorder"
	bodyEcho           = "echo"
	bodySlow           = "slow"
	bodyReverse        = "reverse"
	bodyLiveHandle     = "live-handle"
	bodyLead           = "lead"
	bodyFirstPage      = "first-page"
)

// bodies contains the body of the override of each defect form that a suite
// plants. Forms with one body share it.
var bodies = map[spec.DefectFormID]string{
	spec.Panics:                 bodyPanic,
	spec.PanicsOnRepeat:         bodyPanicLater,
	spec.SucceedsAlways:         bodyZero,
	spec.FailsAlways:            bodyFailSuccess,
	spec.ReportsAnotherError:    bodyFail,
	spec.RefusesRepeat:          bodyFailLater,
	spec.IgnoresThePredicate:    bodyZeroLater,
	spec.IgnoresCancel:          bodyUncancelled,
	spec.ValueBeside:            bodyBesideError,
	spec.ValuesBeside:           bodyBesideError,
	spec.ValueBesideAFalseFlag:  bodyBesideFalse,
	spec.ValuesBesideAFalseFlag: bodyBesideFalse,
	spec.InventsAPointer:        bodySample,
	spec.InventsAValue:          bodySample,
	spec.Drifts:                 bodyDrift,
	spec.Descends:               bodyDescend,
	spec.VersionDescends:        bodyVersionDescend,
	spec.StampsThePast:          bodyVersionZero,
	spec.AnswersASample:         bodySample,
	spec.ReadsAnotherValue:      bodySampleSuccess,
	spec.StoresNothing:          bodyZero,
	spec.DrainsTwice:            bodyDrainTwice,
	spec.ReordersDrain:          bodyReorder,
	spec.DrainsNothing:          bodyZero,
	spec.Echoes:                 bodyEcho,
	spec.AcceptsTheInvalid:      bodyZero,
	spec.ObservationUnmoved:     bodyZero,
	spec.VerifySilent:           bodyZero,
	spec.SettlesTwice:           bodyZero,
	spec.CountDrifts:            bodyDrift,
	spec.Slow:                   bodySlow,
	spec.ReadSurvivesDelete:     bodySample,
	spec.SwallowsHooks:          bodyZero,
	spec.ReplaysReversed:        bodyReverse,
	spec.OpensALiveHandle:       bodyLiveHandle,
	spec.LeadsOn:                bodyLead,
	spec.KeepsFailedWrites:      bodyZero,
	spec.PutRefuses:             bodyFail,
	spec.CompensationRefuses:    bodyFail,
	spec.FilesRefusedWrite:      bodySample,
	spec.ServesOnePage:          bodyFirstPage,
}

// partnerOf returns the pair of the partner name of the check c: the partner
// with that name, or the partner with the role of that name, as a defect form
// names the partner is-leader of the claim leader-election.
func partnerOf(c *derive.Check, name string) *derive.Pair {
	if pair, ok := c.Pairs[name]; ok {
		return pair
	}
	i := slices.IndexFunc(c.Claim.Partners, func(p spec.Partner) bool { return string(p.Role) == name })
	return c.Pairs[c.Claim.Partners[i].Name]
}

// delegating are the bodies that call the callable of the fresh subject, and
// counting the bodies that count the calls of the override.
var (
	delegating = []string{
		bodyPanicLater, bodyFailLater, bodyFailSuccess, bodyZeroLater, bodyBesideFalse, bodyBesideError,
		bodySampleSuccess, bodyUncancelled, bodyVersionZero, bodyDrainTwice, bodyReorder, bodyReverse, bodySlow,
		bodyFirstPage, bodyLiveHandle,
	}
	counting = []string{
		bodyPanicLater, bodyFailLater, bodyZeroLater, bodyDrift, bodyDescend, bodyVersionDescend, bodyReorder,
	}
)

// cancelling are the bodies that do not return the error of a cancelled
// context first: a panic, which needs none, and the bodies whose break is the
// answer to a cancelled context.
var cancelling = []string{bodyPanic, bodyBesideError, bodyUncancelled}

// budgetParam is the param of the mixin timeout that states the budget.
const budgetParam = "duration"

// defect returns the planted defect of the check c, and whether s already has
// it, because each break of each callable has one planted defect.
func (k *caller) defect(s *Suite, c *derive.Check, b *binder) (*Defect, bool) {
	form, _ := spec.DefectFormOf(c.Claim.Defect.Form)
	f := c.Facts
	if partner, ok := strings.CutPrefix(form.Changes, "partner:"); ok {
		f = partnerOf(c, partner).Partner
	}
	words := []string{s.Declaration}
	if f.Declaration != "" {
		words = append(words, f.Callable)
	}
	name := k.e.JoinName(naming.Pascal(string(form.ID)), strings.Join(words, ""))
	for _, d := range s.Defects {
		if d.Name == name {
			return d, true
		}
	}
	d := &Defect{Name: name, Doc: c.Clause, Body: bodies[form.ID]}
	if d.Body == bodyEcho || d.Body == bodyFirstPage {
		d.Input = inputAt(f, 0)
	}
	d.Delegates = slices.Contains(delegating, d.Body)
	d.Counts = slices.Contains(counting, d.Body)
	if f.Declaration != "" {
		d.Type = k.e.JoinName(defectWord, name)
	}
	d.Version = nameOf(b.param(stampParam))
	if v, ok := b.param(budgetParam).(string); ok {
		budget, _ := time.ParseDuration(v)
		d.Budget = 2 * budget
	}
	d.Method, d.Args = k.signatureOf(f)
	d.Results = k.outcomesOf(f)
	if ctx := slices.IndexFunc(f.Signature.Params, func(p rules.ParamView) bool {
		return p.Role == rules.ParamContext
	}); ctx >= 0 {
		d.Context = d.Args[ctx]
		d.Guards = hasError(f) && !slices.Contains(cancelling, d.Body)
	}
	if form.ID == spec.OpensALiveHandle {
		d.Handle = k.liveHandle(f, nameOf(b.param(nextParam)), k.e.JoinName(handleWord, name))
	}
	return d, false
}

// nextParam is the param of the contract cursor that names the method that
// drains the cursor.
const nextParam = "next-method"

// liveHandle returns the handle of a planted defect of the form
// opens-a-live-handle: a type named name that embeds the handle that the
// callable of f returns and overrides its method next, which returns its zero
// results and no error. It returns nil for a handle that is not an
// interface, because no other type can take the place of such a handle.
func (k *caller) liveHandle(f *derive.Facts, next, name string) *Defect {
	ref := k.resultRef(f, 0)
	decl, _ := k.m.Reader().Lookup(ref.Target)
	if _, isInterface := decl.(*node.Interface); !isInterface {
		return nil
	}
	b := k.b
	methods := methodsOf(&b, decl)
	m := methods[slices.IndexFunc(methods, func(m *node.Method) bool { return m.Name == next })]
	facts := k.factsOf(m)
	h := &Defect{Name: name, Type: name, Body: bodyZero, Handled: k.typeOf(ref)}
	h.Method, h.Args = k.signatureOf(&facts)
	h.Results = k.outcomesOf(&facts)
	return h
}

// defectWord is the word that the name of the type of a planted defect adds
// to the name of its function. The type of its live handle adds handleWord.
const defectWord = "defect"

// signatureOf returns the signature of the callable of f in a planted defect,
// and the names of its parameters in a call: the name p and its position of
// each parameter, and of each variadic parameter with the ellipsis.
func (k *caller) signatureOf(f *derive.Facts) (*emit.Method, []string) {
	m := &emit.Method{Origin: f.ID, Name: f.Callable, Visibility: symbol.VisibilityPublic}
	args := make([]string, 0, len(f.Signature.Params))
	for i, p := range f.Signature.Params {
		name := "p" + strconv.Itoa(i)
		param := &emit.Param{Name: name, Type: k.typeOf(p.Ref)}
		arg := name
		if p.Variadic {
			param.Variadic = symbol.VariadicPositional
			arg += "..."
		}
		m.Params.Append(param)
		args = append(args, arg)
	}
	for i, r := range f.Signature.Returns {
		m.Returns.Append(&emit.Return{Name: "r" + strconv.Itoa(i), Type: k.typeOf(r.Ref)})
	}
	return m, args
}

// outcomesOf returns the results of the callable of f, each with its role
// and a sample that is not its zero value.
func (k *caller) outcomesOf(f *derive.Facts) []Outcome {
	out := make([]Outcome, 0, len(f.Signature.Returns))
	n := 0
	for i, r := range f.Signature.Returns {
		res := Outcome{Name: "r" + strconv.Itoa(i)}
		switch r.Role {
		case rules.ReturnError:
			res.Role = roleError
		case rules.ReturnOkBool:
			res.Role = roleFlag
		case rules.ReturnStream:
			res.Role = roleValue
			res.Sample = &Expr{Word: nilWord}
		default:
			res.Role = roleValue
			ref := k.typeOf(k.instance(r.Ref))
			v := &f.Results[n]
			// The sample of a defect differs from the zero value, so a
			// defect that returns it is seen. The two samples differ, so
			// one of them is not the zero value.
			sample := &v.Samples[1].Value
			if reflect.DeepEqual(*sample, v.Zero) {
				sample = &v.Samples[0].Value
			}
			res.Sample = &Expr{Value: typed(ref, sample)}
			n++
		}
		out = append(out, res)
	}
	return out
}
