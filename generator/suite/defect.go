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
	sdkrules "go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/rules"
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
// slow callable waits, twice its budget, and Input and Second the positions
// of the first and the second input. Handle is the handle that a callable of the form opens-a-live-handle
// returns in place of its own, and Handled the type of the handle that it
// wraps. Delegates reports a body that calls the callable of the fresh
// subject, and Counts a body that counts its calls. Keeps reports a body that
// keeps the handle that builds the defect, to read the injector of the handle,
// to build a fresh subject for each call, or to yield to the scheduler of the
// handle. Acquire is the callable that the
// param open of the mixin leak-free names, which a body of the form
// opens-without-closing calls once more on a dependency of its own, and empty
// for any other body. Context is the parameter of the context of the
// callable, and empty for a callable without one. Guards reports a body that
// returns the error of a cancelled context first, so the defect reports a
// cancelled context as the callable does. A panic, a body that breaks
// cancellation and a body that reads the error of the fresh subject have no
// guard.
//
// After is the number of calls that a body of generated sequences makes as
// the callable does before its break. Limit is the most use that the param
// limit of the mixin bounded states, and Grow the count of the use that the
// body of grows-past-the-limit-after-three returns above it: int or len.
// Stream is the name of the result of the stream that the body of
// iterates-a-copy copies, and Elems the types of the values of each element of
// the stream.
//
// Overrides are the overrides of a planted defect that changes every callable
// of its declaration, each with the signature and the body of one callable.
// The template writes them in place of the one override of Method.
//
// Door is the name of the function that opens the planted defect over a
// medium, for a defect of a durable claim, and empty for any other defect.
// The defect then wraps a subject that the door Open of the application opens
// over the medium, through Handle.FreshOpen. Medium is the function of the
// package probe that wraps the medium of a defect that changes the medium,
// such as Unsynced, and empty for a defect that changes a callable. Sentinel
// is the declared error that a body of the form ignores-checksum swallows.
//
// Hosted reports the defect of a claim of a world, whose Door is the door
// Host. The world of its check opens each incarnation of the subject through
// the door, which opens the subject of the application through
// Handle.FreshHost. The defect changes the subject from the second
// incarnation of its host on, and is the subject of the application outside a
// world. Opener is the function of the package probe that opens such a
// defect, such as Amnesiac, and empty for a defect that changes its callables.
//
// Txn is what an override of a planted defect of transactions writes, and nil
// for any other override. Wrapper is the type that wraps the handle of a
// transaction whose methods write, read, commit and roll back, with the
// override of each, and nil for a defect without one.
type Defect struct {
	Method    *emit.Method
	Handle    *Defect
	Handled   *emit.TypeRef
	Sentinel  *Expr
	Txn       *TxnPlant
	Wrapper   *Defect
	Overrides []*Defect
	Name      string
	Type      string
	Doc       string
	Body      string
	Version   string
	Context   string
	Acquire   string
	Limit     string
	Grow      string
	Stream    string
	Door      string
	Medium    string
	Opener    string
	Args      []string
	Results   []Outcome
	Elems     []*emit.TypeRef
	Budget    time.Duration
	Input     int
	Second    int
	After     int
	Delegates bool
	Counts    bool
	Keeps     bool
	Guards    bool
	Hosted    bool
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
	bodyZeroBefore     = "zero-before-failure"
	bodySwallow        = "swallow"
	bodyFailAfter      = "fail-after-failure"
	bodySampleAfter    = "sample-after-failure"
	bodyFreshEachCall  = "fresh-each-call"
	bodyOpenOnFailure  = "open-on-failure"
	bodyZeroSecond     = "zero-second"
	bodyKeepFirst      = "keep-first"
	bodyDropLast       = "drop-last"
	bodyFirstAfter     = "first-after"
	bodyReverseAfter   = "reverse-after"
	bodyVersionAfter   = "version-after"
	bodyZeroAfter      = "zero-after"
	bodyRepeatAfter    = "repeat-after"
	bodyDropAfter      = "drop-after"
	bodyDivergeAfter   = "diverge-after"
	bodyGrowAfter      = "grow-after"
	bodyTwiceAfter     = "twice-after"
	bodyHalfAfter      = "half-applied-after"
	bodyLeak           = "leak"
	bodyCopy           = "copy"
	bodyRacy           = "racy"
	bodyMark           = "mark"
	bodyMarkedSample   = "marked-sample"
	bodySharedResult   = "shared-result"
	bodyCheckThenGrant = "check-then-grant"
	bodyZeroFirst      = "zero-first"
	bodyZeroOnSentinel = "zero-on-sentinel"
	bodyZeroBeside     = "zero-beside-error"
	bodyPark           = "park-after-recovery"
)

// media are the functions of the package probe that wrap the medium of each
// defect form that changes the medium.
var media = map[spec.DefectFormID]string{
	spec.AcknowledgesBeforeSync: "Unsynced",
	spec.HidesAFullMedium:       "Unfilled",
}

// openers are the functions of the package probe that open the defect of each
// defect form of a world that changes the disk of an incarnation.
var openers = map[spec.DefectFormID]string{
	spec.AmnesiacRecovery: "Amnesiac",
}

// The words that the name of the door Open and of the door Host of a planted
// defect add to the name of its function, and the param of the mixin durable
// that refers to the error of a corrupted record.
const (
	openWord     = "Open"
	hostWord     = "Host"
	corruptParam = "corrupt"
)

// afterBodies are the bodies of the defect forms of generated sequences,
// which call the callable as it is until the call that the min-steps of their
// form counts.
var afterBodies = []string{
	bodyFirstAfter, bodyReverseAfter, bodyVersionAfter, bodyZeroAfter, bodyRepeatAfter, bodyDropAfter,
	bodyDivergeAfter, bodyGrowAfter, bodyTwiceAfter, bodyHalfAfter,
}

// limitParam is the param of the mixin bounded that states the most use, and
// growText the growth of a text past the limit.
const (
	limitParam = "limit"
	growText   = "text"
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
	spec.KeepsFailedWrites:      bodyZeroBefore,
	spec.PutRefuses:             bodyFail,
	spec.CompensationRefuses:    bodyFail,
	spec.FilesRefusedWrite:      bodySample,
	spec.ServesOnePage:          bodyFirstPage,
	spec.SwallowsFailure:        bodySwallow,
	spec.GivesUpEarly:           bodyFailAfter,
	spec.HalfApplied:            bodySampleAfter,
	spec.CommitsOnFailure:       bodySampleAfter,
	spec.NeverOpens:             bodyFreshEachCall,
	spec.BalancesEachCall:       bodyFreshEachCall,
	spec.OpensWithoutClosing:    bodyOpenOnFailure,
	spec.IgnoresTheSecondInput:  bodyEcho,
	spec.KeepsItsFirstResult:    bodyKeepFirst,
	spec.DropsTheLastKey:        bodyDropLast,
	spec.ReportsPoison:          bodyFail,
	spec.AdmitsTheSize:          bodySample,
	spec.StopsDelegating:        bodySample,
	spec.ReportsTamper:          bodyFail,
	spec.ReleaseKeepsTheLease:   bodyZero,
	spec.FrozenStats:            bodyZero,

	spec.GoesBackAfterThree:            bodyFirstAfter,
	spec.ReadsThePastAfterThree:        bodyFirstAfter,
	spec.ReadsAStaleVersionAfterThree:  bodyFirstAfter,
	spec.ReadsAStaleValueAfterThree:    bodyFirstAfter,
	spec.ServesAStaleValueAfterThree:   bodyFirstAfter,
	spec.ReordersAfterThree:            bodyReverseAfter,
	spec.ReversesAfterThree:            bodyReverseAfter,
	spec.ReordersWritesAfterThree:      bodyZeroAfter,
	spec.StampsThePastAfterThree:       bodyVersionAfter,
	spec.DescendsAfterThree:            bodyZeroAfter,
	spec.DropsARepeatAfterThree:        bodyZeroAfter,
	spec.DropsAWriteAfterThree:         bodyZeroAfter,
	spec.AdmitsAnyTransitionAfterThree: bodyTwiceAfter,
	spec.RepeatsAfterThree:             bodyRepeatAfter,
	spec.DropsAfterThree:               bodyDropAfter,
	spec.DivergesAfterThree:            bodyDivergeAfter,
	spec.GrowsPastTheLimitAfterThree:   bodyGrowAfter,
	spec.IgnoresTheOtherReplica:        bodyZero,
	spec.LeaksOneUnit:                  bodyLeak,
	spec.HalfAppliedAfterThree:         bodyHalfAfter,
	spec.IteratesACopy:                 bodyCopy,
	spec.DropsAppendsAfterThree:        bodyZeroAfter,
	spec.BreaksTheModeAfterThree:       bodyZeroAfter,
	spec.SkipsARedeliveryAfterThree:    bodyZeroAfter,
	spec.MissesUpdatesAfterThree:       bodyZeroAfter,
	spec.RacyReadWrite:                 bodyRacy,
	spec.RacyRead:                      bodySharedResult,
	spec.RacyPin:                       bodySharedResult,
	spec.RacyCausalRead:                bodySharedResult,
	spec.ComputesPerCall:               bodyFreshEachCall,
	spec.CheckThenGrant:                bodyCheckThenGrant,
	spec.TrustsAPartialRecord:          bodyZeroFirst,
	spec.IgnoresChecksum:               bodyZeroOnSentinel,
	spec.DropsTheField:                 bodyZeroBeside,
	spec.BlocksAfterRecovery:           bodyPark,
}

// partnerOf returns the pair of the partner name of the check c: the partner
// with that name, or the partner with the role of that name, as a defect form
// names the partner is-leader of the claim leader-election.
func partnerOf(c *derive.Check, name string) *derive.Pair {
	if pair, ok := c.Pairs[name]; ok {
		return pair
	}
	i := slices.IndexFunc(c.Claim.Partners, func(p rules.Partner) bool { return string(p.Role) == name })
	return c.Pairs[c.Claim.Partners[i].Name]
}

// delegating are the bodies that call the callable of the fresh subject,
// counting the bodies that count the calls of the override, and keeping the
// bodies that keep the handle that builds the defect.
var (
	delegating = slices.Concat([]string{
		bodyPanicLater, bodyFailLater, bodyFailSuccess, bodyZeroLater, bodyBesideFalse, bodyBesideError,
		bodySampleSuccess, bodyUncancelled, bodyVersionZero, bodyDrainTwice, bodyReorder, bodyReverse, bodySlow,
		bodyFirstPage, bodyLiveHandle, bodyZeroBefore, bodySwallow, bodyFailAfter, bodySampleAfter, bodyOpenOnFailure,
		bodyZeroSecond, bodyKeepFirst, bodyDropLast, bodyLeak, bodyCopy, bodySharedResult, bodyCheckThenGrant,
		bodyZeroFirst, bodyZeroOnSentinel, bodyZeroBeside, bodyPark,
	}, afterBodies)
	counting = slices.Concat([]string{
		bodyPanicLater, bodyFailLater, bodyZeroLater, bodyDrift, bodyDescend, bodyVersionDescend, bodyReorder, bodyLeak,
	}, afterBodies)
	keeping = []string{
		bodyZeroBefore, bodySwallow, bodyFailAfter, bodySampleAfter, bodyFreshEachCall, bodyHalfAfter, bodySharedResult,
		bodyCheckThenGrant,
	}
)

// cancelling are the bodies that do not return the error of a cancelled
// context first: a panic, which needs none, and the bodies whose break is the
// answer to a cancelled context.
var cancelling = []string{bodyPanic, bodyBesideError, bodyUncancelled}

// budgetParam is the param of the mixin timeout that states the budget, and
// openParam the param of the mixin leak-free whose value is the callable that
// acquires a resource. modeParam is the param of the contract publisher that
// states the delivery mode. atLeastOnce is the mode that allows a second
// delivery, and atMostOnce the mode that allows a lost event.
const (
	budgetParam = "duration"
	openParam   = "open"
	modeParam   = "mode"
	atLeastOnce = "at-least-once"
	atMostOnce  = "at-most-once"
)

// defect returns the planted defect of the check c, and whether s already has
// it, because each break of each callable has one planted defect.
func (k *caller) defect(s *Suite, c *derive.Check, b *binder) (*Defect, bool) {
	form, _ := rules.DefectFormOf(c.Claim.Defect.Form)
	f := c.Facts
	if partner, ok := strings.CutPrefix(form.Changes, "partner:"); ok {
		f = partnerOf(c, partner).Partner
	}
	racy := form.ID == spec.RacyReadWrite
	words := []string{s.Declaration}
	if f.Declaration != "" && !slices.Contains(wholeForms, form.ID) && c.Transaction == nil {
		words = append(words, f.Callable)
	}
	name := k.e.JoinName(naming.Pascal(string(form.ID)), strings.Join(words, ""))
	for _, d := range s.Defects {
		if d.Name == name {
			return d, true
		}
	}
	d := &Defect{Name: name, Doc: c.Clause, Body: bodies[form.ID], Medium: media[form.ID], Opener: openers[form.ID]}
	switch slots := probes[c.Claim.Probe]; {
	case slices.ContainsFunc(slots, func(sl slot) bool { return sl.form == formWorld }):
		d.Door, d.Hosted = k.e.JoinName(hostWord, name), true
	case slices.ContainsFunc(slots, func(sl slot) bool { return sl.form == formDurable }):
		d.Door = k.e.JoinName(openWord, name)
	}
	switch {
	case form.ID == spec.IgnoresTheSecondInput && !reflect.DeepEqual(f.Results[0].Shape, f.Inputs[0].Shape):
		// A callable whose result has another type cannot return its first
		// input, so it calls with a zero second input.
		d.Body = bodyZeroSecond
	case form.ID == spec.DrainsNothing && streamAt(f) >= 0:
		// A range over a nil stream panics, so a callable that returns a
		// stream returns an empty one.
		d.Body = bodySample
	case form.ID == spec.BreaksTheModeAfterThree && nameOf(b.param(modeParam)) == atMostOnce:
		// A publisher at most once may drop an event, so it breaks its mode
		// by a second delivery of each event.
		d.Body = bodyTwiceAfter
	default:
	}
	switch d.Body {
	case bodyEcho, bodyFirstPage, bodyKeepFirst:
		d.Input = inputAt(f, 0)
	case bodyZeroSecond:
		d.Second = inputAt(f, 1)
	case bodyCopy:
		at := streamAt(f)
		d.Stream = "r" + strconv.Itoa(at)
		for _, e := range k.elemsOf(f.Signature.Returns[at].Ref) {
			d.Elems = append(d.Elems, k.typeOf(e))
		}
	case bodyGrowAfter:
		d.Limit = intOf(b.param(limitParam))
		if use, ok := k.countableOf(s, s.Sequences, f); ok {
			d.Grow = use.Count
			if f.Results[0].Shape.Form == symbol.FormText {
				d.Grow = growText
			}
		}
	case bodyZeroFirst:
		_, pair := b.facts(durable.name)
		d.Input = inputAt(f, pair.Written)
	case bodyZeroOnSentinel:
		d.Sentinel = k.valueOfParam(b.param(corruptParam))
	default:
	}
	if slices.Contains(afterBodies, d.Body) {
		d.After = form.MinSteps - 1
	}
	d.Delegates = slices.Contains(delegating, d.Body)
	d.Counts = slices.Contains(counting, d.Body)
	d.Keeps = slices.Contains(keeping, d.Body)
	if d.Body == bodyOpenOnFailure {
		d.Acquire = nameOf(b.param(openParam))
	}
	if f.Declaration != "" && d.Medium == "" && d.Opener == "" {
		d.Type = k.e.JoinName(defectWord, name)
	}
	d.Version = nameOf(b.param(stampParam))
	if v, ok := b.param(budgetParam).(string); ok {
		budget, _ := time.ParseDuration(v)
		d.Budget = 2 * budget
	}
	d.Method, d.Args = k.signatureOf(f)
	d.Results = k.outcomesOf(f)
	guard(d, f)
	if form.ID == spec.OpensALiveHandle {
		d.Handle = k.liveHandle(f, nameOf(b.param(nextParam)), k.e.JoinName(handleWord, name))
	}
	if racy {
		d.Keeps, d.Overrides = true, k.racy(s)
	}
	if form.ID == spec.BlocksAfterRecovery {
		d.Overrides = k.parked(s)
	}
	if c.Transaction != nil {
		k.txnDefect(s, d, form.ID, c.Transaction)
	}
	return d, false
}

// guard sets the parameter of the context of the callable of f on d, and
// whether d returns the error of a cancelled context first.
func guard(d *Defect, f *derive.Facts) {
	if ctx := slices.IndexFunc(f.Signature.Params, func(p sdkrules.ParamView) bool {
		return p.Role == sdkrules.ParamContext
	}); ctx >= 0 {
		d.Context = d.Args[ctx]
		d.Guards = hasError(f) && !slices.Contains(cancelling, d.Body)
	}
}

// racy returns the overrides of the planted defect racy-read-write of the
// methods of s. The first method with a value result returns its sample while
// another call is marked. Every other method marks its call, yields to the
// scheduler of the case, and then calls the subject.
func (k *caller) racy(s *Suite) []*Defect {
	read := slices.IndexFunc(s.Facts, func(f derive.Facts) bool { return len(f.Results) > 0 })
	out := make([]*Defect, 0, len(s.Facts))
	for i := range s.Facts {
		f := &s.Facts[i]
		o := &Defect{Body: bodyMark, Delegates: true}
		if i == read {
			o.Body = bodyMarkedSample
		}
		o.Method, o.Args = k.signatureOf(f)
		o.Results = k.outcomesOf(f)
		guard(o, f)
		out = append(out, o)
	}
	return out
}

// parked returns the overrides of the planted defect blocks-after-recovery of
// the callables of s. Each override calls the subject as it is in the first
// incarnation of its host. From the second incarnation on, it blocks until
// its context ends: the context of the call, or the lifetime of the
// incarnation for a callable without one.
func (k *caller) parked(s *Suite) []*Defect {
	out := make([]*Defect, 0, len(s.Facts))
	for i := range s.Facts {
		f := &s.Facts[i]
		o := &Defect{Body: bodyPark, Delegates: true}
		o.Method, o.Args = k.signatureOf(f)
		o.Results = k.outcomesOf(f)
		guard(o, f)
		out = append(out, o)
	}
	return out
}

// wholeForms are the defect forms whose planted defect changes the whole
// declaration, so the name of the defect names the declaration alone. Such a
// defect overrides every callable of its declaration, or opens each
// incarnation of the subject after a restart over an empty disk.
var wholeForms = []spec.DefectFormID{spec.RacyReadWrite, spec.BlocksAfterRecovery, spec.AmnesiacRecovery}

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

// spellsZero reports whether the value v spells the zero value of its type:
// false, the number 0, the empty string, nil, a conversion of such a value,
// or a composite literal whose entries are named fields that each spell a
// zero value. eidos writes a number as its shortest decimal text, so 0 is the
// only text of a zero number. A list or a map literal with an entry is not
// nil, so it does not spell a zero value. A name and a call do not spell a
// zero value either, because their text does not show their value.
func spellsZero(v *emit.Value) bool {
	switch v.Kind {
	case emit.ValueLiteral:
		switch v.Literal {
		case emit.LiteralBool:
			return v.Text == falseSpelling
		case emit.LiteralInt, emit.LiteralFloat:
			return v.Text == zeroSpelling
		default:
			return v.Text == ""
		}
	case emit.ValueConversion:
		return spellsZero(v.Inner)
	case emit.ValueComposite:
		return !slices.ContainsFunc(v.Fields, func(f emit.ValueField) bool {
			return f.Name == "" || !spellsZero(&f.Value)
		})
	default:
		return false
	}
}

// falseSpelling is the text of the literal false, and zeroSpelling the text
// of the number 0.
const (
	falseSpelling = "false"
	zeroSpelling  = "0"
)

// outcomesOf returns the results of the callable of f, each with its role
// and a sample that is not its zero value.
func (k *caller) outcomesOf(f *derive.Facts) []Outcome {
	out := make([]Outcome, 0, len(f.Signature.Returns))
	n := 0
	for i, r := range f.Signature.Returns {
		res := Outcome{Name: "r" + strconv.Itoa(i)}
		switch r.Role {
		case sdkrules.ReturnError:
			res.Role = roleError
		case sdkrules.ReturnOkBool:
			res.Role = roleFlag
		case sdkrules.ReturnStream:
			// The sample of a stream is an empty stream of its elements.
			res.Role = roleValue
			res.Sample = &Expr{}
			for _, e := range k.elemsOf(r.Ref) {
				res.Sample.Stream = append(res.Sample.Stream, k.typeOf(e))
			}
		default:
			res.Role = roleValue
			ref := k.typeOf(k.instance(r.Ref))
			v := &f.Results[n]
			// The sample of a defect differs from the zero value, so a
			// defect that returns it is seen. The two samples differ, so
			// one of them is not the zero value.
			sample := &v.Samples[1].Value
			if spellsZero(sample) {
				sample = &v.Samples[0].Value
			}
			res.Sample = &Expr{Value: typed(ref, sample)}
			n++
		}
		out = append(out, res)
	}
	return out
}
