package inputs

import (
	"context"
	"io"
	"sync/atomic"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Pipeline:
//
// Pipeline.Each/smoke: Each returns without a panic for derived inputs.
//
// Pipeline.Do/smoke: Do returns without a panic for derived inputs.
//
// Pipeline.Notify/smoke: Notify returns without a panic for derived
// inputs.
//
// Pipeline.Apply/smoke: Apply returns without a panic for derived inputs.
//
// Pipeline.Put/smoke: Put returns without a panic for derived inputs.
//
// Pipeline.Label/smoke: Label returns without a panic for derived inputs.

// The derived checks of Pipeline that witness their laws at the shallow
// end:
//
// Pipeline.Each/smoke
//
// Pipeline.Do/smoke
//
// Pipeline.Notify/smoke
//
// Pipeline.Apply/smoke
//
// Pipeline.Put/smoke
//
// Pipeline.Label/smoke

// The refused claims of Pipeline:
//
// Pipeline.Pump/cancel: eidos derives no sample of the parameter from of
// Pump. The remedy is to author a sample of the parameter from with the
// directive sample.
//
// Pipeline.Pump/zero-on-error: eidos derives no sample of the parameter
// from of Pump. The remedy is to author a sample of the parameter from
// with the directive sample.
//
// Pipeline.Fold/cancel: eidos derives no sample of the parameter items of
// Fold. The remedy is to author a sample of the parameter items with the
// directive sample.
//
// Pipeline.Count/deterministic: eidos derives no sample of the parameter
// ints of Count. The remedy is to author a sample of the parameter ints
// with the directive sample.

// The unclaimed classifications of Pipeline:
//
// The shape:composite-writer of Pipeline.Merge has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Pipeline.Each has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:writer of Pipeline.Handle has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Pipeline.Do has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Pipeline.Notify has the stance silent. A mutation
// is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Pipeline.Apply has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:writer of Pipeline.Fold has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:multi-arg-writer of Pipeline.Put has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Pipeline.Spin has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:writer of Pipeline.Copy has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Pipeline.Subscribe has the stance silent. A
// mutation is an input of the other checks, and its result is an error
// alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Pipeline writes.
type (
	PipelineHandle = conform.Handle[Pipeline]
	PipelineCheck  = conform.Check[Pipeline]
	PipelineOption = conform.Option[Pipeline]
)

// PipelineChecks contains the name of each callable of Pipeline, with the
// identities of its checks.
var PipelineChecks = struct {
	Pump      pipelinePump
	Merge     pipelineMerge
	Each      pipelineEach
	Handle    pipelineHandle
	Do        pipelineDo
	Notify    pipelineNotify
	Apply     pipelineApply
	Fold      pipelineFold
	Put       pipelinePut
	Count     pipelineCount
	Spin      pipelineSpin
	Copy      pipelineCopy
	Configure pipelineConfigure
	Subscribe pipelineSubscribe
	Label     pipelineLabel
}{
	Pump:      pipelinePump{"Pipeline.Pump"},
	Merge:     pipelineMerge{"Pipeline.Merge"},
	Each:      pipelineEach{"Pipeline.Each"},
	Handle:    pipelineHandle{"Pipeline.Handle"},
	Do:        pipelineDo{"Pipeline.Do"},
	Notify:    pipelineNotify{"Pipeline.Notify"},
	Apply:     pipelineApply{"Pipeline.Apply"},
	Fold:      pipelineFold{"Pipeline.Fold"},
	Put:       pipelinePut{"Pipeline.Put"},
	Count:     pipelineCount{"Pipeline.Count"},
	Spin:      pipelineSpin{"Pipeline.Spin"},
	Copy:      pipelineCopy{"Pipeline.Copy"},
	Configure: pipelineConfigure{"Pipeline.Configure"},
	Subscribe: pipelineSubscribe{"Pipeline.Subscribe"},
	Label:     pipelineLabel{"Pipeline.Label"},
}

// pipelinePump is the name of Pipeline.Pump, with the identities of its
// checks.
type pipelinePump struct {
	conform.Callable[Pipeline]
}

// pipelineMerge is the name of Pipeline.Merge, with the identities of its
// checks.
type pipelineMerge struct {
	conform.Callable[Pipeline]
}

// pipelineEach is the name of Pipeline.Each, with the identities of its
// checks.
type pipelineEach struct {
	conform.Callable[Pipeline]
}

// Smoke returns the identity of the check that Each returns without a
// panic for derived inputs.
func (c pipelineEach) Smoke() conform.ID[Pipeline] {
	return c.Check("smoke")
}

// pipelineHandle is the name of Pipeline.Handle, with the identities of
// its checks.
type pipelineHandle struct {
	conform.Callable[Pipeline]
}

// pipelineDo is the name of Pipeline.Do, with the identities of its
// checks.
type pipelineDo struct {
	conform.Callable[Pipeline]
}

// Smoke returns the identity of the check that Do returns without a panic
// for derived inputs.
func (c pipelineDo) Smoke() conform.ID[Pipeline] {
	return c.Check("smoke")
}

// pipelineNotify is the name of Pipeline.Notify, with the identities of
// its checks.
type pipelineNotify struct {
	conform.Callable[Pipeline]
}

// Smoke returns the identity of the check that Notify returns without a
// panic for derived inputs.
func (c pipelineNotify) Smoke() conform.ID[Pipeline] {
	return c.Check("smoke")
}

// pipelineApply is the name of Pipeline.Apply, with the identities of its
// checks.
type pipelineApply struct {
	conform.Callable[Pipeline]
}

// Smoke returns the identity of the check that Apply returns without a
// panic for derived inputs.
func (c pipelineApply) Smoke() conform.ID[Pipeline] {
	return c.Check("smoke")
}

// pipelineFold is the name of Pipeline.Fold, with the identities of its
// checks.
type pipelineFold struct {
	conform.Callable[Pipeline]
}

// pipelinePut is the name of Pipeline.Put, with the identities of its
// checks.
type pipelinePut struct {
	conform.Callable[Pipeline]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c pipelinePut) Smoke() conform.ID[Pipeline] {
	return c.Check("smoke")
}

// pipelineCount is the name of Pipeline.Count, with the identities of its
// checks.
type pipelineCount struct {
	conform.Callable[Pipeline]
}

// pipelineSpin is the name of Pipeline.Spin, with the identities of its
// checks.
type pipelineSpin struct {
	conform.Callable[Pipeline]
}

// pipelineCopy is the name of Pipeline.Copy, with the identities of its
// checks.
type pipelineCopy struct {
	conform.Callable[Pipeline]
}

// pipelineConfigure is the name of Pipeline.Configure, with the identities
// of its checks.
type pipelineConfigure struct {
	conform.Callable[Pipeline]
}

// pipelineSubscribe is the name of Pipeline.Subscribe, with the identities
// of its checks.
type pipelineSubscribe struct {
	conform.Callable[Pipeline]
}

// pipelineLabel is the name of Pipeline.Label, with the identities of its
// checks.
type pipelineLabel struct {
	conform.Callable[Pipeline]
}

// Smoke returns the identity of the check that Label returns without a
// panic for derived inputs.
func (c pipelineLabel) Smoke() conform.ID[Pipeline] {
	return c.Check("smoke")
}

// PipelineSuite returns the checks of Pipeline, and the calls and the
// budgets of its callables.
func PipelineSuite() conform.Suite[Pipeline] {
	return conform.Suite[Pipeline]{
		Declaration: "Pipeline",
		Checks: []conform.Check[Pipeline]{
			{
				ID:          PipelineChecks.Each.Smoke(),
				Statement:   "Each returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PipelineHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Each(func(Item) bool { return false })
						return nil
					})
				},
				Defect: pipelineEachPanics,
			},
			{
				ID:          PipelineChecks.Do.Smoke(),
				Statement:   "Do returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PipelineHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Do(func() (int, error) { return 0, nil })
						return nil
					})
				},
				Defect: pipelineDoPanics,
			},
			{
				ID:          PipelineChecks.Notify.Smoke(),
				Statement:   "Notify returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PipelineHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Notify(func(event string) {})
						return nil
					})
				},
				Defect: pipelineNotifyPanics,
			},
			{
				ID:          PipelineChecks.Apply.Smoke(),
				Statement:   "Apply returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PipelineHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Apply(func(counts ...int) {})
						return nil
					})
				},
				Defect: pipelineApplyPanics,
			},
			{
				ID:          PipelineChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PipelineHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(Item{Count: 1}, Line{Qty: 1}, DefaultNote())
					})
				},
				Defect: pipelinePutPanics,
			},
			{
				ID:          PipelineChecks.Label.Smoke(),
				Statement:   "Label returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PipelineHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().
							Label(Tag("new"), DefaultCode(), Level(42), Text("text-a"), Bag[int]{42})
						return nil
					})
				},
				Defect: pipelineLabelPanics,
			},
		},
		Calls: []conform.Call[Pipeline]{
			{
				Calls:    pipelinePumpCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Pump.Callable,
			},
			{
				Calls:    pipelineMergeCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Merge.Callable,
			},
			{
				Calls:    pipelineEachCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Each.Callable,
			},
			{
				Calls:    pipelineHandleCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Handle.Callable,
			},
			{
				Calls:    pipelineDoCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Do.Callable,
			},
			{
				Calls:    pipelineNotifyCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Notify.Callable,
			},
			{
				Calls:    pipelineApplyCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Apply.Callable,
			},
			{
				Calls:    pipelineFoldCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Fold.Callable,
			},
			{
				Calls:    pipelinePutCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Put.Callable,
			},
			{
				Calls:    pipelineCountCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Count.Callable,
			},
			{
				Calls:    pipelineSpinCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Spin.Callable,
			},
			{
				Calls:    pipelineCopyCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Copy.Callable,
			},
			{
				Calls:    pipelineConfigureCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Configure.Callable,
			},
			{
				Calls:    pipelineSubscribeCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Subscribe.Callable,
			},
			{
				Calls:    pipelineLabelCalls,
				Over:     newPipelineOver,
				Callable: PipelineChecks.Label.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.SampledFrom(Raw, Cooked)),
			prop.Using(pipelineItemValues),
			prop.Using(pipelineLineValues),
			prop.Using(pipelineNoteValues),
			prop.Using(pipelineTagValues),
			prop.Using(pipelineCodeValues),
			prop.Using(prop.Composite(func(c *prop.Case) io.Reader {
				return conform.NewBytes(c.Draw(prop.Of[[]byte](), "reader"))
			})),
		},
	}
}

// RunPipelineSuite applies the checks of Pipeline, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunPipelineSuite[T Pipeline](t *testing.T, newSubject func() T, opts ...PipelineOption) {
	t.Helper()
	RunPipelineSuiteFromHandle(t, func(*PipelineHandle) T { return newSubject() }, opts...)
}

// RunPipelineSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPipelineSuiteFromHandle[T Pipeline](
	t *testing.T,
	newSubject func(h *PipelineHandle) T,
	opts ...PipelineOption,
) {
	t.Helper()
	conform.Run(t, PipelineSuite(), pipelineSubjects(newSubject), opts...)
}

// FuzzPipelineSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPipelineSuite[T Pipeline](
	f *testing.F,
	callable conform.Callable[Pipeline],
	newSubject func(h *PipelineHandle) T,
	opts ...PipelineOption,
) {
	f.Helper()
	conform.Fuzz(f, PipelineSuite(), callable, pipelineSubjects(newSubject), opts...)
}

// RunPipelineBudgets applies the allocation ceilings of Pipeline and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunPipelineBudgets[T Pipeline](
	t *testing.T,
	newSubject func(h *PipelineHandle) T,
	opts ...PipelineOption,
) {
	t.Helper()
	conform.RunBudgets(t, PipelineSuite(), pipelineSubjects(newSubject), opts...)
}

// BenchPipelineSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchPipelineSuite[T Pipeline](
	b conform.B,
	callable conform.Callable[Pipeline],
	newSubject func(h *PipelineHandle) T,
	opts ...PipelineOption,
) {
	b.Helper()
	conform.Bench(b, PipelineSuite(), callable, pipelineSubjects(newSubject), opts...)
}

// pipelineSubjects returns the factory of the subjects of Pipeline that
// calls newSubject with the handle of each check.
func pipelineSubjects[T Pipeline](
	newSubject func(h *PipelineHandle) T,
) func(h *PipelineHandle) Pipeline {
	return func(h *PipelineHandle) Pipeline { return newSubject(h) }
}

// pipelinePumpCalls generates the calls of Pipeline.Pump with the inputs
// that each case draws.
var pipelinePumpCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	fromNextResults0 := c.Draw(prop.List(prop.Of[Item]()), "from-next-0")
	fromNextResults1 := c.Draw(
		prop.List(
			prop.Of[bool](),
			prop.MinSize(len(fromNextResults0)),
			prop.MaxSize(len(fromNextResults0)),
		),
		"from-next-1",
	)
	fromNextResults2 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, ErrDrained),
			prop.MinSize(len(fromNextResults0)),
			prop.MaxSize(len(fromNextResults0)),
		),
		"from-next-2",
	)
	var fromNextAt atomic.Int64
	fromGetResults0 := c.Draw(prop.List(prop.Of[Item]()), "from-get-0")
	fromGetResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, ErrMissing),
			prop.MinSize(len(fromGetResults0)),
			prop.MaxSize(len(fromGetResults0)),
		),
		"from-get-1",
	)
	var fromGetAt atomic.Int64
	fromCloseResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"from-close",
	)
	var fromCloseAt atomic.Int64
	from := &pipelineFeedInput{
		onNext: func(context.Context) (r0 Item, r1 bool, r2 error) {
			i := int(fromNextAt.Add(1) - 1)
			if i >= len(fromNextResults0) {
				r2 = ErrDrained
				return r0, r1, r2
			}
			return fromNextResults0[i],
				fromNextResults1[i],
				fromNextResults2[i]
		},
		onGet: func(context.Context, string) (r0 Item, r1 error) {
			i := int(fromGetAt.Add(1) - 1)
			if i >= len(fromGetResults0) {
				r1 = ErrMissing
				return r0, r1
			}
			return fromGetResults0[i],
				fromGetResults1[i]
		},
		onSeek: func(...string) {},
		onClose: func() (r0 error) {
			i := int(fromCloseAt.Add(1) - 1)
			if i >= len(fromCloseResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return fromCloseResults[i]
		},
	}
	return func(ctx context.Context, subject Pipeline) error {
		_, err := subject.Pump(ctx, from)
		return err
	}
})

// pipelineMergeCalls generates the calls of Pipeline.Merge with the inputs
// that each case draws.
var pipelineMergeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	leftNextResults0 := c.Draw(prop.List(prop.Of[Item]()), "left-next-0")
	leftNextResults1 := c.Draw(
		prop.List(
			prop.Of[bool](),
			prop.MinSize(len(leftNextResults0)),
			prop.MaxSize(len(leftNextResults0)),
		),
		"left-next-1",
	)
	leftNextResults2 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, ErrDrained),
			prop.MinSize(len(leftNextResults0)),
			prop.MaxSize(len(leftNextResults0)),
		),
		"left-next-2",
	)
	var leftNextAt atomic.Int64
	leftGetResults0 := c.Draw(prop.List(prop.Of[Item]()), "left-get-0")
	leftGetResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, ErrMissing),
			prop.MinSize(len(leftGetResults0)),
			prop.MaxSize(len(leftGetResults0)),
		),
		"left-get-1",
	)
	var leftGetAt atomic.Int64
	leftCloseResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"left-close",
	)
	var leftCloseAt atomic.Int64
	left := &pipelineFeedInput{
		onNext: func(context.Context) (r0 Item, r1 bool, r2 error) {
			i := int(leftNextAt.Add(1) - 1)
			if i >= len(leftNextResults0) {
				r2 = ErrDrained
				return r0, r1, r2
			}
			return leftNextResults0[i],
				leftNextResults1[i],
				leftNextResults2[i]
		},
		onGet: func(context.Context, string) (r0 Item, r1 error) {
			i := int(leftGetAt.Add(1) - 1)
			if i >= len(leftGetResults0) {
				r1 = ErrMissing
				return r0, r1
			}
			return leftGetResults0[i],
				leftGetResults1[i]
		},
		onSeek: func(...string) {},
		onClose: func() (r0 error) {
			i := int(leftCloseAt.Add(1) - 1)
			if i >= len(leftCloseResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return leftCloseResults[i]
		},
	}
	rightNextResults0 := c.Draw(prop.List(prop.Of[Item]()), "right-next-0")
	rightNextResults1 := c.Draw(
		prop.List(
			prop.Of[bool](),
			prop.MinSize(len(rightNextResults0)),
			prop.MaxSize(len(rightNextResults0)),
		),
		"right-next-1",
	)
	rightNextResults2 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, ErrDrained),
			prop.MinSize(len(rightNextResults0)),
			prop.MaxSize(len(rightNextResults0)),
		),
		"right-next-2",
	)
	var rightNextAt atomic.Int64
	rightGetResults0 := c.Draw(prop.List(prop.Of[Item]()), "right-get-0")
	rightGetResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, ErrMissing),
			prop.MinSize(len(rightGetResults0)),
			prop.MaxSize(len(rightGetResults0)),
		),
		"right-get-1",
	)
	var rightGetAt atomic.Int64
	rightCloseResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"right-close",
	)
	var rightCloseAt atomic.Int64
	right := &pipelineFeedInput{
		onNext: func(context.Context) (r0 Item, r1 bool, r2 error) {
			i := int(rightNextAt.Add(1) - 1)
			if i >= len(rightNextResults0) {
				r2 = ErrDrained
				return r0, r1, r2
			}
			return rightNextResults0[i],
				rightNextResults1[i],
				rightNextResults2[i]
		},
		onGet: func(context.Context, string) (r0 Item, r1 error) {
			i := int(rightGetAt.Add(1) - 1)
			if i >= len(rightGetResults0) {
				r1 = ErrMissing
				return r0, r1
			}
			return rightGetResults0[i],
				rightGetResults1[i]
		},
		onSeek: func(...string) {},
		onClose: func() (r0 error) {
			i := int(rightCloseAt.Add(1) - 1)
			if i >= len(rightCloseResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return rightCloseResults[i]
		},
	}
	return func(_ context.Context, subject Pipeline) error {
		return subject.Merge(left, right)
	}
})

// pipelineEachCalls generates the calls of Pipeline.Each with the inputs
// that each case draws.
var pipelineEachCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	fnResults := c.Draw(prop.List(prop.Of[bool]()), "fn")
	var fnAt atomic.Int64
	fn := func(Item) (r0 bool) {
		i := int(fnAt.Add(1) - 1)
		if i >= len(fnResults) {
			return r0
		}
		return fnResults[i]
	}
	return func(_ context.Context, subject Pipeline) error {
		subject.Each(fn)
		return nil
	}
})

// pipelineHandleCalls generates the calls of Pipeline.Handle with the
// inputs that each case draws.
var pipelineHandleCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	hResults := c.Draw(prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)), "h")
	var hAt atomic.Int64
	h := func(Item) (r0 error) {
		i := int(hAt.Add(1) - 1)
		if i >= len(hResults) {
			r0 = conform.ErrGenerated
			return r0
		}
		return hResults[i]
	}
	return func(_ context.Context, subject Pipeline) error {
		return subject.Handle(h)
	}
})

// pipelineDoCalls generates the calls of Pipeline.Do with the inputs that
// each case draws.
var pipelineDoCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	opResults0 := c.Draw(prop.List(prop.Of[int]()), "op-0")
	opResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(opResults0)),
			prop.MaxSize(len(opResults0)),
		),
		"op-1",
	)
	var opAt atomic.Int64
	op := func() (r0 int, r1 error) {
		i := int(opAt.Add(1) - 1)
		if i >= len(opResults0) {
			r1 = conform.ErrGenerated
			return r0, r1
		}
		return opResults0[i],
			opResults1[i]
	}
	return func(_ context.Context, subject Pipeline) error {
		subject.Do(op)
		return nil
	}
})

// pipelineNotifyCalls generates the calls of Pipeline.Notify with the
// inputs that each case draws.
var pipelineNotifyCalls = prop.Composite(func(*prop.Case) func(context.Context, Pipeline) error {
	fn := func(string) {}
	return func(_ context.Context, subject Pipeline) error {
		subject.Notify(fn)
		return nil
	}
})

// pipelineApplyCalls generates the calls of Pipeline.Apply with the inputs
// that each case draws.
var pipelineApplyCalls = prop.Composite(func(*prop.Case) func(context.Context, Pipeline) error {
	fn := func(...int) {}
	return func(_ context.Context, subject Pipeline) error {
		subject.Apply(fn)
		return nil
	}
})

// pipelineFoldCalls generates the calls of Pipeline.Fold with the inputs
// that each case draws.
var pipelineFoldCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	itemsValues := c.Draw(prop.List(prop.Of[Item]()), "items")
	return func(ctx context.Context, subject Pipeline) error {
		return subject.Fold(ctx, func() <-chan Item {
			ch := make(chan Item, len(itemsValues))
			for _, v := range itemsValues {
				ch <- v
			}
			close(ch)
			return ch
		}())
	}
})

// pipelinePutCalls generates the calls of Pipeline.Put with the inputs
// that each case draws.
var pipelinePutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	item := c.Draw(prop.Of[Item](), "item")
	line := c.Draw(prop.Of[Line](), "line")
	note := c.Draw(prop.Of[Note](), "note")
	return func(_ context.Context, subject Pipeline) error {
		return subject.Put(item, line, note)
	}
})

// pipelineCountCalls generates the calls of Pipeline.Count with the inputs
// that each case draws.
var pipelineCountCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	intsReadResults := c.Draw(prop.List(prop.Of[int]()), "ints-read")
	var intsReadAt atomic.Int64
	ints := &pipelineSourceInput{
		onRead: func() (r0 int) {
			i := int(intsReadAt.Add(1) - 1)
			if i >= len(intsReadResults) {
				return r0
			}
			return intsReadResults[i]
		},
	}
	textsReadResults := c.Draw(prop.List(prop.Of[string]()), "texts-read")
	var textsReadAt atomic.Int64
	texts := &pipelineSourceInput0{
		onRead: func() (r0 string) {
			i := int(textsReadAt.Add(1) - 1)
			if i >= len(textsReadResults) {
				return r0
			}
			return textsReadResults[i]
		},
	}
	return func(_ context.Context, subject Pipeline) error {
		_ = subject.Count(ints, texts)
		return nil
	}
})

// pipelineSpinCalls generates the calls of Pipeline.Spin with the inputs
// that each case draws.
var pipelineSpinCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	l := c.Draw(prop.Of[Loop1](), "l")
	return func(_ context.Context, subject Pipeline) error {
		subject.Spin(l)
		return nil
	}
})

// pipelineCopyCalls generates the calls of Pipeline.Copy with the inputs
// that each case draws.
var pipelineCopyCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	src := c.Draw(prop.Of[io.Reader](), "src")
	return func(_ context.Context, subject Pipeline) error {
		return subject.Copy(src)
	}
})

// pipelineConfigureCalls generates the calls of Pipeline.Configure with
// the inputs that each case draws.
var pipelineConfigureCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, Pipeline) error {
		name := c.Draw(prop.Of[string](), "name")
		return func(_ context.Context, subject Pipeline) error {
			subject.Configure(name)
			return nil
		}
	},
)

// pipelineSubscribeCalls generates the calls of Pipeline.Subscribe with
// the inputs that each case draws.
var pipelineSubscribeCalls = prop.Composite(func(*prop.Case) func(context.Context, Pipeline) error {
	return func(_ context.Context, subject Pipeline) error {
		subject.Subscribe()
		return nil
	}
})

// pipelineLabelCalls generates the calls of Pipeline.Label with the inputs
// that each case draws.
var pipelineLabelCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pipeline) error {
	tag := c.Draw(prop.Of[Tag](), "tag")
	code := c.Draw(prop.Of[Code](), "code")
	level := c.Draw(prop.Of[Level](), "level")
	text := c.Draw(prop.Of[Text](), "text")
	bag := c.Draw(prop.Of[Bag[int]](), "bag")
	return func(_ context.Context, subject Pipeline) error {
		subject.Label(tag, code, level, text, bag)
		return nil
	}
})

// pipelineItemValues generates the values of Item that keep or change each
// field of its starting value.
var pipelineItemValues = prop.Composite(func(c *prop.Case) Item {
	changed := c.Draw(prop.Of[Item](), "item")
	value := NewItemBuilder().Build()
	if c.Draw(prop.Boolean(), "item.id") {
		value.ID = changed.ID
	}
	if c.Draw(prop.Boolean(), "item.count") {
		value.Count = changed.Count
	}
	if c.Draw(prop.Boolean(), "item.kind") {
		value.Kind = changed.Kind
	}
	if c.Draw(prop.Boolean(), "item.base") {
		value.Base = changed.Base
	}
	return value
})

// pipelineLineValues generates the values of Line that keep or change each
// field of its starting value.
var pipelineLineValues = prop.Composite(func(c *prop.Case) Line {
	changed := c.Draw(prop.Of[Line](), "line")
	value := Line{Qty: 1}
	if c.Draw(prop.Boolean(), "line.sku") {
		value.SKU = changed.SKU
	}
	if c.Draw(prop.Boolean(), "line.qty") {
		value.Qty = changed.Qty
	}
	if c.Draw(prop.Boolean(), "line.rank") {
		value.Rank = changed.Rank
	}
	return value
})

// pipelineNoteValues generates the values of Note that keep or change each
// field of its starting value.
var pipelineNoteValues = prop.Composite(func(c *prop.Case) Note {
	changed := c.Draw(prop.Of[Note](), "note")
	value := DefaultNote()
	if c.Draw(prop.Boolean(), "note.text") {
		value.Text = changed.Text
	}
	return value
})

// pipelineTagValues generates the values of Tag that keep or change its
// starting value.
var pipelineTagValues = prop.Composite(func(c *prop.Case) Tag {
	changed := c.Draw(prop.Of[Tag](), "tag")
	if c.Draw(prop.Boolean(), "tag.changed") {
		return changed
	}
	return "new"
})

// pipelineCodeValues generates the values of Code that keep or change its
// starting value.
var pipelineCodeValues = prop.Composite(func(c *prop.Case) Code {
	changed := c.Draw(prop.Of[Code](), "code")
	if c.Draw(prop.Boolean(), "code.changed") {
		return changed
	}
	return DefaultCode()
})

// pipelineFeedInputEmbedded is Feed under the name of the field that
// pipelineFeedInput embeds.
type pipelineFeedInputEmbedded = Feed

// pipelineFeedInput implements Feed. Its methods return the results that
// the generators of the suite of Pipeline draw for them, in turn.
type pipelineFeedInput struct {
	pipelineFeedInputEmbedded
	onNext  func(a0 context.Context) (Item, bool, error)
	onGet   func(a0 context.Context, a1 string) (Item, error)
	onSeek  func(a0 ...string)
	onClose func() error
}

// Next calls the function of Next.
func (in pipelineFeedInput) Next(a0 context.Context) (Item, bool, error) {
	return in.onNext(a0)
}

// Get calls the function of Get.
func (in pipelineFeedInput) Get(a0 context.Context, a1 string) (Item, error) {
	return in.onGet(a0, a1)
}

// Seek calls the function of Seek.
func (in pipelineFeedInput) Seek(a0 ...string) {
	in.onSeek(a0...)
}

// Close calls the function of Close.
func (in pipelineFeedInput) Close() error {
	return in.onClose()
}

// pipelineSourceInputEmbedded is Source under the name of the field that
// pipelineSourceInput embeds.
type pipelineSourceInputEmbedded = Source[int]

// pipelineSourceInput implements Source. Its methods return the results
// that the generators of the suite of Pipeline draw for them, in turn.
type pipelineSourceInput struct {
	pipelineSourceInputEmbedded
	onRead func() int
}

// Read calls the function of Read.
func (in pipelineSourceInput) Read() int {
	return in.onRead()
}

// pipelineSourceInput0Embedded is Source under the name of the field that
// pipelineSourceInput0 embeds.
type pipelineSourceInput0Embedded = Source[string]

// pipelineSourceInput0 implements Source. Its methods return the results
// that the generators of the suite of Pipeline draw for them, in turn.
type pipelineSourceInput0 struct {
	pipelineSourceInput0Embedded
	onRead func() string
}

// Read calls the function of Read.
func (in pipelineSourceInput0) Read() string {
	return in.onRead()
}

// pipelineOverEmbedded is Pipeline under the name of the field that
// pipelineOver embeds.
type pipelineOverEmbedded = Pipeline

// pipelineOver is a subject of Pipeline whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type pipelineOver struct {
	pipelineOverEmbedded
	spend func()
}

// Pump calls spend and returns its zero results.
func (o pipelineOver) Pump(context.Context, Feed) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// Merge calls spend and returns its zero results.
func (o pipelineOver) Merge(Feed, Feed) (r0 error) {
	o.spend()
	return r0
}

// Each calls spend and returns its zero results.
func (o pipelineOver) Each(func(Item) bool) {
	o.spend()
}

// Handle calls spend and returns its zero results.
func (o pipelineOver) Handle(Handler) (r0 error) {
	o.spend()
	return r0
}

// Do calls spend and returns its zero results.
func (o pipelineOver) Do(func() (int, error)) {
	o.spend()
}

// Notify calls spend and returns its zero results.
func (o pipelineOver) Notify(func(event string)) {
	o.spend()
}

// Apply calls spend and returns its zero results.
func (o pipelineOver) Apply(func(counts ...int)) {
	o.spend()
}

// Fold calls spend and returns its zero results.
func (o pipelineOver) Fold(context.Context, <-chan Item) (r0 error) {
	o.spend()
	return r0
}

// Put calls spend and returns its zero results.
func (o pipelineOver) Put(Item, Line, Note) (r0 error) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o pipelineOver) Count(Source[int], Source[string]) (r0 int) {
	o.spend()
	return r0
}

// Spin calls spend and returns its zero results.
func (o pipelineOver) Spin(Loop1) {
	o.spend()
}

// Copy calls spend and returns its zero results.
func (o pipelineOver) Copy(io.Reader) (r0 error) {
	o.spend()
	return r0
}

// Configure calls spend and returns its zero results.
func (o pipelineOver) Configure(string, ...Option) {
	o.spend()
}

// Subscribe calls spend and returns its zero results.
func (o pipelineOver) Subscribe(...Feed) {
	o.spend()
}

// Label calls spend and returns its zero results.
func (o pipelineOver) Label(Tag, Code, Level, Text, Bag[int]) {
	o.spend()
}

// newPipelineOver returns the planted defect of a ceiling of Pipeline,
// whose callables call spend.
func newPipelineOver(_ *PipelineHandle, spend func()) Pipeline {
	return pipelineOver{spend: spend}
}

// pipelineEachPanicsDefect is the planted defect of a Pipeline whose Each
// panics.
type pipelineEachPanicsDefect struct {
	pipelineOverEmbedded
}

// Each is the override of the planted defect.
func (*pipelineEachPanicsDefect) Each(p0 func(Item) bool) {
	panic("the planted defect of a Pipeline whose Each panics")
}

// pipelineEachPanics returns the planted defect of a Pipeline whose Each
// panics over a fresh subject.
func pipelineEachPanics(h *PipelineHandle) Pipeline {
	return &pipelineEachPanicsDefect{pipelineOverEmbedded: h.Fresh()}
}

// pipelineDoPanicsDefect is the planted defect of a Pipeline whose Do
// panics.
type pipelineDoPanicsDefect struct {
	pipelineOverEmbedded
}

// Do is the override of the planted defect.
func (*pipelineDoPanicsDefect) Do(p0 func() (int, error)) {
	panic("the planted defect of a Pipeline whose Do panics")
}

// pipelineDoPanics returns the planted defect of a Pipeline whose Do
// panics over a fresh subject.
func pipelineDoPanics(h *PipelineHandle) Pipeline {
	return &pipelineDoPanicsDefect{pipelineOverEmbedded: h.Fresh()}
}

// pipelineNotifyPanicsDefect is the planted defect of a Pipeline whose
// Notify panics.
type pipelineNotifyPanicsDefect struct {
	pipelineOverEmbedded
}

// Notify is the override of the planted defect.
func (*pipelineNotifyPanicsDefect) Notify(p0 func(event string)) {
	panic("the planted defect of a Pipeline whose Notify panics")
}

// pipelineNotifyPanics returns the planted defect of a Pipeline whose
// Notify panics over a fresh subject.
func pipelineNotifyPanics(h *PipelineHandle) Pipeline {
	return &pipelineNotifyPanicsDefect{pipelineOverEmbedded: h.Fresh()}
}

// pipelineApplyPanicsDefect is the planted defect of a Pipeline whose
// Apply panics.
type pipelineApplyPanicsDefect struct {
	pipelineOverEmbedded
}

// Apply is the override of the planted defect.
func (*pipelineApplyPanicsDefect) Apply(p0 func(counts ...int)) {
	panic("the planted defect of a Pipeline whose Apply panics")
}

// pipelineApplyPanics returns the planted defect of a Pipeline whose Apply
// panics over a fresh subject.
func pipelineApplyPanics(h *PipelineHandle) Pipeline {
	return &pipelineApplyPanicsDefect{pipelineOverEmbedded: h.Fresh()}
}

// pipelinePutPanicsDefect is the planted defect of a Pipeline whose Put
// panics.
type pipelinePutPanicsDefect struct {
	pipelineOverEmbedded
}

// Put is the override of the planted defect.
func (*pipelinePutPanicsDefect) Put(p0 Item, p1 Line, p2 Note) (r0 error) {
	panic("the planted defect of a Pipeline whose Put panics")
}

// pipelinePutPanics returns the planted defect of a Pipeline whose Put
// panics over a fresh subject.
func pipelinePutPanics(h *PipelineHandle) Pipeline {
	return &pipelinePutPanicsDefect{pipelineOverEmbedded: h.Fresh()}
}

// pipelineLabelPanicsDefect is the planted defect of a Pipeline whose
// Label panics.
type pipelineLabelPanicsDefect struct {
	pipelineOverEmbedded
}

// Label is the override of the planted defect.
func (*pipelineLabelPanicsDefect) Label(p0 Tag, p1 Code, p2 Level, p3 Text, p4 Bag[int]) {
	panic("the planted defect of a Pipeline whose Label panics")
}

// pipelineLabelPanics returns the planted defect of a Pipeline whose Label
// panics over a fresh subject.
func pipelineLabelPanics(h *PipelineHandle) Pipeline {
	return &pipelineLabelPanicsDefect{pipelineOverEmbedded: h.Fresh()}
}
