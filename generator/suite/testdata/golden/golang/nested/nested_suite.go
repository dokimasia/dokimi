package nested

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

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Courier writes.
type (
	CourierHandle = conform.Handle[Courier]
	CourierCheck  = conform.Check[Courier]
	CourierOption = conform.Option[Courier]
)

// CourierChecks contains the name of each callable of Courier, with the
// identities of its checks.
var CourierChecks = struct {
	Ship courierShip
	Send courierSend
}{
	Ship: courierShip{"Courier.Ship"},
	Send: courierSend{"Courier.Send"},
}

// courierShip is the name of Courier.Ship, with the identities of its
// checks.
type courierShip struct {
	conform.Callable[Courier]
}

// Smoke returns the identity of the check that Ship returns without a
// panic for derived inputs.
func (c courierShip) Smoke() conform.ID[Courier] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Ship returns an error for
// a cancelled context.
func (c courierShip) Cancel() conform.ID[Courier] {
	return c.Check("cancel")
}

// courierSend is the name of Courier.Send, with the identities of its
// checks.
type courierSend struct {
	conform.Callable[Courier]
}

// Smoke returns the identity of the check that Send returns without a
// panic for derived inputs.
func (c courierSend) Smoke() conform.ID[Courier] {
	return c.Check("smoke")
}

// CourierSuite returns the checks of Courier, and the calls and the
// budgets of its callables.
func CourierSuite() conform.Suite[Courier] {
	return conform.Suite[Courier]{
		Declaration: "Courier",
		Checks: []conform.Check[Courier]{
			{
				ID:          CourierChecks.Ship.Smoke(),
				Statement:   "Ship returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CourierHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Ship(ctx, &ShipRequest{ID: "ID-a"})
					})
				},
				Defect: courierShipPanics,
			},
			{
				ID:          CourierChecks.Ship.Cancel(),
				Statement:   "Ship returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CourierHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Ship(ctx, &ShipRequest{ID: "ID-a"})
					})
				},
				Defect: courierShipSucceedsAlways,
			},
			{
				ID:          CourierChecks.Send.Smoke(),
				Statement:   "Send returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CourierHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Send(Batch{Requests: []ShipRequest{{ID: "ID-a"}}})
					})
				},
				Defect: courierSendPanics,
			},
		},
		Calls: []conform.Call[Courier]{
			{
				Calls:    courierShipCalls,
				Over:     newCourierOver,
				Callable: CourierChecks.Ship.Callable,
			},
			{
				Calls:    courierSendCalls,
				Over:     newCourierOver,
				Callable: CourierChecks.Send.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) Labeler {
				labelerLabelResults0 := c.Draw(prop.List(prop.Of[Receipt]()), "labeler-label-0")
				labelerLabelResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, ErrTorn),
						prop.MinSize(len(labelerLabelResults0)),
						prop.MaxSize(len(labelerLabelResults0)),
					),
					"labeler-label-1",
				)
				var labelerLabelAt atomic.Int64
				return &courierLabelerInput{
					onLabel: func(string) (r0 Receipt, r1 error) {
						i := int(labelerLabelAt.Add(1) - 1)
						if i >= len(labelerLabelResults0) {
							r1 = ErrTorn
							return r0, r1
						}
						return labelerLabelResults0[i],
							labelerLabelResults1[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) func() error {
				funcResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"func",
				)
				var funcAt atomic.Int64
				return func() (r0 error) {
					i := int(funcAt.Add(1) - 1)
					if i >= len(funcResults) {
						r0 = conform.ErrGenerated
						return r0
					}
					return funcResults[i]
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) Ack {
				ackResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"ack",
				)
				var ackAt atomic.Int64
				return func(string) (r0 error) {
					i := int(ackAt.Add(1) - 1)
					if i >= len(ackResults) {
						r0 = conform.ErrGenerated
						return r0
					}
					return ackResults[i]
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) Scale[int] {
				scaleWeighResults := c.Draw(prop.List(prop.Of[int]()), "scale-weigh")
				var scaleWeighAt atomic.Int64
				return &courierScaleInput{
					onWeigh: func() (r0 int) {
						i := int(scaleWeighAt.Add(1) - 1)
						if i >= len(scaleWeighResults) {
							return r0
						}
						return scaleWeighResults[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(*prop.Case) func() {
				return func() {}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) io.Reader {
				return conform.NewBytes(c.Draw(prop.Of[[]byte](), "reader"))
			})),
		},
	}
}

// RunCourierSuite applies the checks of Courier, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCourierSuite[T Courier](t *testing.T, newSubject func() T, opts ...CourierOption) {
	t.Helper()
	RunCourierSuiteFromHandle(t, func(*CourierHandle) T { return newSubject() }, opts...)
}

// RunCourierSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCourierSuiteFromHandle[T Courier](
	t *testing.T,
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	t.Helper()
	conform.Run(t, CourierSuite(), courierSubjects(newSubject), opts...)
}

// FuzzCourierSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCourierSuite[T Courier](
	f *testing.F,
	callable conform.Callable[Courier],
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	f.Helper()
	conform.Fuzz(f, CourierSuite(), callable, courierSubjects(newSubject), opts...)
}

// RunCourierBudgets applies the allocation ceilings of Courier and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCourierBudgets[T Courier](
	t *testing.T,
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	t.Helper()
	conform.RunBudgets(t, CourierSuite(), courierSubjects(newSubject), opts...)
}

// BenchCourierSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchCourierSuite[T Courier](
	b conform.B,
	callable conform.Callable[Courier],
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	b.Helper()
	conform.Bench(b, CourierSuite(), callable, courierSubjects(newSubject), opts...)
}

// courierSubjects returns the factory of the subjects of Courier that
// calls newSubject with the handle of each check.
func courierSubjects[T Courier](
	newSubject func(h *CourierHandle) T,
) func(h *CourierHandle) Courier {
	return func(h *CourierHandle) Courier { return newSubject(h) }
}

// courierShipCalls generates the calls of Courier.Ship with the inputs
// that each case draws.
var courierShipCalls = prop.Composite(func(c *prop.Case) func(context.Context, Courier) error {
	req := c.Draw(prop.Of[ShipRequest](), "req")
	return func(ctx context.Context, subject Courier) error {
		return subject.Ship(ctx, &req)
	}
})

// courierSendCalls generates the calls of Courier.Send with the inputs
// that each case draws.
var courierSendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Courier) error {
	b := c.Draw(prop.Of[Batch](), "b")
	return func(_ context.Context, subject Courier) error {
		return subject.Send(b)
	}
})

// courierLabelerInputEmbedded is Labeler under the name of the field that
// courierLabelerInput embeds.
type courierLabelerInputEmbedded = Labeler

// courierLabelerInput implements Labeler. Its methods return the results
// that the generators of the suite of Courier draw for them, in turn.
type courierLabelerInput struct {
	courierLabelerInputEmbedded
	onLabel func(a0 string) (Receipt, error)
}

// Label calls the function of Label.
func (in courierLabelerInput) Label(a0 string) (Receipt, error) {
	return in.onLabel(a0)
}

// courierScaleInputEmbedded is Scale under the name of the field that
// courierScaleInput embeds.
type courierScaleInputEmbedded = Scale[int]

// courierScaleInput implements Scale. Its methods return the results that
// the generators of the suite of Courier draw for them, in turn.
type courierScaleInput struct {
	courierScaleInputEmbedded
	onWeigh func() int
}

// Weigh calls the function of Weigh.
func (in courierScaleInput) Weigh() int {
	return in.onWeigh()
}

// courierOverEmbedded is Courier under the name of the field that
// courierOver embeds.
type courierOverEmbedded = Courier

// courierOver is a subject of Courier whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type courierOver struct {
	courierOverEmbedded
	spend func()
}

// Ship calls spend and returns its zero results.
func (o courierOver) Ship(context.Context, *ShipRequest) (r0 error) {
	o.spend()
	return r0
}

// Send calls spend and returns its zero results.
func (o courierOver) Send(Batch) (r0 error) {
	o.spend()
	return r0
}

// newCourierOver returns the planted defect of a ceiling of Courier, whose
// callables call spend.
func newCourierOver(_ *CourierHandle, spend func()) Courier {
	return courierOver{spend: spend}
}

// courierShipPanicsDefect is the planted defect of a Courier whose Ship
// panics.
type courierShipPanicsDefect struct {
	courierOverEmbedded
}

// Ship is the override of the planted defect.
func (*courierShipPanicsDefect) Ship(p0 context.Context, p1 *ShipRequest) (r0 error) {
	panic("the planted defect of a Courier whose Ship panics")
}

// courierShipPanics returns the planted defect of a Courier whose Ship
// panics over a fresh subject.
func courierShipPanics(h *CourierHandle) Courier {
	return &courierShipPanicsDefect{courierOverEmbedded: h.Fresh()}
}

// courierShipSucceedsAlwaysDefect is the planted defect of a Courier whose
// Ship ignores the context.
type courierShipSucceedsAlwaysDefect struct {
	courierOverEmbedded
}

// Ship is the override of the planted defect.
func (*courierShipSucceedsAlwaysDefect) Ship(p0 context.Context, p1 *ShipRequest) (r0 error) {
	return r0
}

// courierShipSucceedsAlways returns the planted defect of a Courier whose
// Ship ignores the context over a fresh subject.
func courierShipSucceedsAlways(h *CourierHandle) Courier {
	return &courierShipSucceedsAlwaysDefect{courierOverEmbedded: h.Fresh()}
}

// courierSendPanicsDefect is the planted defect of a Courier whose Send
// panics.
type courierSendPanicsDefect struct {
	courierOverEmbedded
}

// Send is the override of the planted defect.
func (*courierSendPanicsDefect) Send(p0 Batch) (r0 error) {
	panic("the planted defect of a Courier whose Send panics")
}

// courierSendPanics returns the planted defect of a Courier whose Send
// panics over a fresh subject.
func courierSendPanics(h *CourierHandle) Courier {
	return &courierSendPanicsDefect{courierOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Tracker writes.
type (
	TrackerHandle = conform.Handle[Tracker]
	TrackerCheck  = conform.Check[Tracker]
	TrackerOption = conform.Option[Tracker]
)

// TrackerChecks contains the name of each callable of Tracker, with the
// identities of its checks.
var TrackerChecks = struct {
	Route  trackerRoute
	Follow trackerFollow
	Walk   trackerWalk
	Volley trackerVolley
	Track  trackerTrack
	Pick   trackerPick
	Play   trackerPlay
	Link   trackerLink
}{
	Route:  trackerRoute{"Tracker.Route"},
	Follow: trackerFollow{"Tracker.Follow"},
	Walk:   trackerWalk{"Tracker.Walk"},
	Volley: trackerVolley{"Tracker.Volley"},
	Track:  trackerTrack{"Tracker.Track"},
	Pick:   trackerPick{"Tracker.Pick"},
	Play:   trackerPlay{"Tracker.Play"},
	Link:   trackerLink{"Tracker.Link"},
}

// trackerRoute is the name of Tracker.Route, with the identities of its
// checks.
type trackerRoute struct {
	conform.Callable[Tracker]
}

// trackerFollow is the name of Tracker.Follow, with the identities of its
// checks.
type trackerFollow struct {
	conform.Callable[Tracker]
}

// trackerWalk is the name of Tracker.Walk, with the identities of its
// checks.
type trackerWalk struct {
	conform.Callable[Tracker]
}

// trackerVolley is the name of Tracker.Volley, with the identities of its
// checks.
type trackerVolley struct {
	conform.Callable[Tracker]
}

// trackerTrack is the name of Tracker.Track, with the identities of its
// checks.
type trackerTrack struct {
	conform.Callable[Tracker]
}

// trackerPick is the name of Tracker.Pick, with the identities of its
// checks.
type trackerPick struct {
	conform.Callable[Tracker]
}

// trackerPlay is the name of Tracker.Play, with the identities of its
// checks.
type trackerPlay struct {
	conform.Callable[Tracker]
}

// trackerLink is the name of Tracker.Link, with the identities of its
// checks.
type trackerLink struct {
	conform.Callable[Tracker]
}

// TrackerSuite returns the checks of Tracker, and the calls and the
// budgets of its callables.
func TrackerSuite() conform.Suite[Tracker] {
	return conform.Suite[Tracker]{
		Declaration: "Tracker",
		Calls: []conform.Call[Tracker]{
			{
				Calls:    trackerRouteCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Route.Callable,
			},
			{
				Calls:    trackerFollowCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Follow.Callable,
			},
			{
				Calls:    trackerWalkCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Walk.Callable,
			},
			{
				Calls:    trackerVolleyCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Volley.Callable,
			},
			{
				Calls:    trackerTrackCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Track.Callable,
			},
			{
				Calls:    trackerPickCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Pick.Callable,
			},
			{
				Calls:    trackerPlayCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Play.Callable,
			},
			{
				Calls:    trackerLinkCalls,
				Over:     newTrackerOver,
				Callable: TrackerChecks.Link.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Recursive(prop.Just[Node](&trackerNodeInput{
				onName: func() (r0 string) {
					return r0
				},
				onChildren: func() (r0 []Node, r1 error) {
					r1 = conform.ErrGenerated
					return r0, r1
				},
			}),
				func(prop.Generator[Node]) prop.Generator[Node] {
					return prop.Composite(func(c *prop.Case) Node {
						nodeNameResults := c.Draw(prop.List(prop.Of[string]()), "node-name")
						var nodeNameAt atomic.Int64
						nodeChildrenResults0 := c.Draw(
							prop.List(prop.Of[[]Node]()),
							"node-children-0",
						)
						nodeChildrenResults1 := c.Draw(
							prop.List(
								prop.SampledFrom[error](nil, conform.ErrGenerated),
								prop.MinSize(len(nodeChildrenResults0)),
								prop.MaxSize(len(nodeChildrenResults0)),
							),
							"node-children-1",
						)
						var nodeChildrenAt atomic.Int64
						return &trackerNodeInput{
							onName: func() (r0 string) {
								i := int(nodeNameAt.Add(1) - 1)
								if i >= len(nodeNameResults) {
									return r0
								}
								return nodeNameResults[i]
							},
							onChildren: func() (r0 []Node, r1 error) {
								i := int(nodeChildrenAt.Add(1) - 1)
								if i >= len(nodeChildrenResults0) {
									r1 = conform.ErrGenerated
									return r0, r1
								}
								return nodeChildrenResults0[i],
									nodeChildrenResults1[i]
							},
						}
					})
				}, prop.MaxLeaves(10))),
			prop.Using(prop.Recursive(prop.Just[Branch](&trackerBranchInput{
				onNext: func() (r0 Hop) {
					return r0
				},
			}),
				func(prop.Generator[Branch]) prop.Generator[Branch] {
					return prop.Composite(func(c *prop.Case) Branch {
						branchNextResults := c.Draw(prop.List(prop.Of[Hop]()), "branch-next")
						var branchNextAt atomic.Int64
						return &trackerBranchInput{
							onNext: func() (r0 Hop) {
								i := int(branchNextAt.Add(1) - 1)
								if i >= len(branchNextResults) {
									return r0
								}
								return branchNextResults[i]
							},
						}
					})
				}, prop.MaxLeaves(10))),
			prop.Using(prop.Recursive(prop.Just[Step](func(string) (r0 Step) {
				return r0
			}),
				func(prop.Generator[Step]) prop.Generator[Step] {
					return prop.Composite(func(c *prop.Case) Step {
						stepResults := c.Draw(prop.List(prop.Of[Step]()), "step")
						var stepAt atomic.Int64
						return func(string) (r0 Step) {
							i := int(stepAt.Add(1) - 1)
							if i >= len(stepResults) {
								return r0
							}
							return stepResults[i]
						}
					})
				}, prop.MaxLeaves(10))),
			prop.Using(prop.Recursive(prop.Just[Pong](&trackerPongInput{
				onPing: func() (r0 Ping) {
					return r0
				},
			}),
				func(prop.Generator[Pong]) prop.Generator[Pong] {
					return prop.Composite(func(c *prop.Case) Pong {
						pongPingResults := c.Draw(prop.List(prop.Of[Ping]()), "pong-ping")
						var pongPingAt atomic.Int64
						return &trackerPongInput{
							onPing: func() (r0 Ping) {
								i := int(pongPingAt.Add(1) - 1)
								if i >= len(pongPingResults) {
									return r0
								}
								return pongPingResults[i]
							},
						}
					})
				}, prop.MaxLeaves(10))),
			prop.Using(prop.Composite(func(*prop.Case) func() {
				return func() {}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) Pair {
				pairLeftResults := c.Draw(prop.List(prop.Of[Leaf]()), "pair-left")
				var pairLeftAt atomic.Int64
				pairRightResults := c.Draw(prop.List(prop.Of[Leaf]()), "pair-right")
				var pairRightAt atomic.Int64
				return &trackerPairInput{
					onLeft: func() (r0 Leaf) {
						i := int(pairLeftAt.Add(1) - 1)
						if i >= len(pairLeftResults) {
							return r0
						}
						return pairLeftResults[i]
					},
					onRight: func() (r0 Leaf) {
						i := int(pairRightAt.Add(1) - 1)
						if i >= len(pairRightResults) {
							return r0
						}
						return pairRightResults[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) Rally {
				rallyServeResults := c.Draw(prop.List(prop.Of[Ping]()), "rally-serve")
				var rallyServeAt atomic.Int64
				return &trackerRallyInput{
					onServe: func() (r0 Ping) {
						i := int(rallyServeAt.Add(1) - 1)
						if i >= len(rallyServeResults) {
							return r0
						}
						return rallyServeResults[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) Linker {
				linkerLinkResults := c.Draw(prop.List(prop.Of[Chain]()), "linker-link")
				var linkerLinkAt atomic.Int64
				return &trackerLinkerInput{
					onLink: func() (r0 Chain) {
						i := int(linkerLinkAt.Add(1) - 1)
						if i >= len(linkerLinkResults) {
							return r0
						}
						return linkerLinkResults[i]
					},
				}
			})),
			prop.Using(prop.Recursive(prop.Just[Ping](&trackerPingInput{
				onPong: func() (r0 Pong) {
					return r0
				},
			}),
				func(prop.Generator[Ping]) prop.Generator[Ping] {
					return prop.Composite(func(c *prop.Case) Ping {
						pingPongResults := c.Draw(prop.List(prop.Of[Pong]()), "ping-pong")
						var pingPongAt atomic.Int64
						return &trackerPingInput{
							onPong: func() (r0 Pong) {
								i := int(pingPongAt.Add(1) - 1)
								if i >= len(pingPongResults) {
									return r0
								}
								return pingPongResults[i]
							},
						}
					})
				}, prop.MaxLeaves(10))),
			prop.Using(prop.Composite(func(c *prop.Case) Leaf {
				leafResults := c.Draw(prop.List(prop.Of[string]()), "leaf")
				var leafAt atomic.Int64
				return func() (r0 string) {
					i := int(leafAt.Add(1) - 1)
					if i >= len(leafResults) {
						return r0
					}
					return leafResults[i]
				}
			})),
		},
	}
}

// RunTrackerSuite applies the checks of Tracker, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunTrackerSuite[T Tracker](t *testing.T, newSubject func() T, opts ...TrackerOption) {
	t.Helper()
	RunTrackerSuiteFromHandle(t, func(*TrackerHandle) T { return newSubject() }, opts...)
}

// RunTrackerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTrackerSuiteFromHandle[T Tracker](
	t *testing.T,
	newSubject func(h *TrackerHandle) T,
	opts ...TrackerOption,
) {
	t.Helper()
	conform.Run(t, TrackerSuite(), trackerSubjects(newSubject), opts...)
}

// FuzzTrackerSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTrackerSuite[T Tracker](
	f *testing.F,
	callable conform.Callable[Tracker],
	newSubject func(h *TrackerHandle) T,
	opts ...TrackerOption,
) {
	f.Helper()
	conform.Fuzz(f, TrackerSuite(), callable, trackerSubjects(newSubject), opts...)
}

// RunTrackerBudgets applies the allocation ceilings of Tracker and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunTrackerBudgets[T Tracker](
	t *testing.T,
	newSubject func(h *TrackerHandle) T,
	opts ...TrackerOption,
) {
	t.Helper()
	conform.RunBudgets(t, TrackerSuite(), trackerSubjects(newSubject), opts...)
}

// BenchTrackerSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchTrackerSuite[T Tracker](
	b conform.B,
	callable conform.Callable[Tracker],
	newSubject func(h *TrackerHandle) T,
	opts ...TrackerOption,
) {
	b.Helper()
	conform.Bench(b, TrackerSuite(), callable, trackerSubjects(newSubject), opts...)
}

// trackerSubjects returns the factory of the subjects of Tracker that
// calls newSubject with the handle of each check.
func trackerSubjects[T Tracker](
	newSubject func(h *TrackerHandle) T,
) func(h *TrackerHandle) Tracker {
	return func(h *TrackerHandle) Tracker { return newSubject(h) }
}

// trackerRouteCalls generates the calls of Tracker.Route with the inputs
// that each case draws.
var trackerRouteCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	rootNameResults := c.Draw(prop.List(prop.Of[string]()), "root-name")
	var rootNameAt atomic.Int64
	rootChildrenResults0 := c.Draw(prop.List(prop.Of[[]Node]()), "root-children-0")
	rootChildrenResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(rootChildrenResults0)),
			prop.MaxSize(len(rootChildrenResults0)),
		),
		"root-children-1",
	)
	var rootChildrenAt atomic.Int64
	root := &trackerNodeInput{
		onName: func() (r0 string) {
			i := int(rootNameAt.Add(1) - 1)
			if i >= len(rootNameResults) {
				return r0
			}
			return rootNameResults[i]
		},
		onChildren: func() (r0 []Node, r1 error) {
			i := int(rootChildrenAt.Add(1) - 1)
			if i >= len(rootChildrenResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return rootChildrenResults0[i],
				rootChildrenResults1[i]
		},
	}
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Route(root)
		return nil
	}
})

// trackerFollowCalls generates the calls of Tracker.Follow with the inputs
// that each case draws.
var trackerFollowCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	bNextResults := c.Draw(prop.List(prop.Of[Hop]()), "b-next")
	var bNextAt atomic.Int64
	b := &trackerBranchInput{
		onNext: func() (r0 Hop) {
			i := int(bNextAt.Add(1) - 1)
			if i >= len(bNextResults) {
				return r0
			}
			return bNextResults[i]
		},
	}
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Follow(b)
		return nil
	}
})

// trackerWalkCalls generates the calls of Tracker.Walk with the inputs
// that each case draws.
var trackerWalkCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	sResults := c.Draw(prop.List(prop.Of[Step]()), "s")
	var sAt atomic.Int64
	s := func(string) (r0 Step) {
		i := int(sAt.Add(1) - 1)
		if i >= len(sResults) {
			return r0
		}
		return sResults[i]
	}
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Walk(s)
		return nil
	}
})

// trackerVolleyCalls generates the calls of Tracker.Volley with the inputs
// that each case draws.
var trackerVolleyCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	pPongResults := c.Draw(prop.List(prop.Of[Pong]()), "p-pong")
	var pPongAt atomic.Int64
	p := &trackerPingInput{
		onPong: func() (r0 Pong) {
			i := int(pPongAt.Add(1) - 1)
			if i >= len(pPongResults) {
				return r0
			}
			return pPongResults[i]
		},
	}
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Volley(p)
		return nil
	}
})

// trackerTrackCalls generates the calls of Tracker.Track with the inputs
// that each case draws.
var trackerTrackCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	lLabelResults0 := c.Draw(prop.List(prop.Of[Receipt]()), "l-label-0")
	lLabelResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, ErrTorn),
			prop.MinSize(len(lLabelResults0)),
			prop.MaxSize(len(lLabelResults0)),
		),
		"l-label-1",
	)
	var lLabelAt atomic.Int64
	l := &trackerLabelerInput{
		onLabel: func(string) (r0 Receipt, r1 error) {
			i := int(lLabelAt.Add(1) - 1)
			if i >= len(lLabelResults0) {
				r1 = ErrTorn
				return r0, r1
			}
			return lLabelResults0[i],
				lLabelResults1[i]
		},
	}
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Track(l)
		return nil
	}
})

// trackerPickCalls generates the calls of Tracker.Pick with the inputs
// that each case draws.
var trackerPickCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	c0 := c.Draw(prop.Of[Choice](), "c")
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Pick(c0)
		return nil
	}
})

// trackerPlayCalls generates the calls of Tracker.Play with the inputs
// that each case draws.
var trackerPlayCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	m := c.Draw(prop.Of[Match](), "m")
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Play(m)
		return nil
	}
})

// trackerLinkCalls generates the calls of Tracker.Link with the inputs
// that each case draws.
var trackerLinkCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tracker) error {
	l := c.Draw(prop.Of[Links](), "l")
	return func(_ context.Context, subject Tracker) error {
		_ = subject.Link(l)
		return nil
	}
})

// trackerNodeInputEmbedded is Node under the name of the field that
// trackerNodeInput embeds.
type trackerNodeInputEmbedded = Node

// trackerNodeInput implements Node. Its methods return the results that
// the generators of the suite of Tracker draw for them, in turn.
type trackerNodeInput struct {
	trackerNodeInputEmbedded
	onName     func() string
	onChildren func() ([]Node, error)
}

// Name calls the function of Name.
func (in trackerNodeInput) Name() string {
	return in.onName()
}

// Children calls the function of Children.
func (in trackerNodeInput) Children() ([]Node, error) {
	return in.onChildren()
}

// trackerBranchInputEmbedded is Branch under the name of the field that
// trackerBranchInput embeds.
type trackerBranchInputEmbedded = Branch

// trackerBranchInput implements Branch. Its methods return the results
// that the generators of the suite of Tracker draw for them, in turn.
type trackerBranchInput struct {
	trackerBranchInputEmbedded
	onNext func() Hop
}

// Next calls the function of Next.
func (in trackerBranchInput) Next() Hop {
	return in.onNext()
}

// trackerPingInputEmbedded is Ping under the name of the field that
// trackerPingInput embeds.
type trackerPingInputEmbedded = Ping

// trackerPingInput implements Ping. Its methods return the results that
// the generators of the suite of Tracker draw for them, in turn.
type trackerPingInput struct {
	trackerPingInputEmbedded
	onPong func() Pong
}

// Pong calls the function of Pong.
func (in trackerPingInput) Pong() Pong {
	return in.onPong()
}

// trackerLabelerInputEmbedded is Labeler under the name of the field that
// trackerLabelerInput embeds.
type trackerLabelerInputEmbedded = Labeler

// trackerLabelerInput implements Labeler. Its methods return the results
// that the generators of the suite of Tracker draw for them, in turn.
type trackerLabelerInput struct {
	trackerLabelerInputEmbedded
	onLabel func(a0 string) (Receipt, error)
}

// Label calls the function of Label.
func (in trackerLabelerInput) Label(a0 string) (Receipt, error) {
	return in.onLabel(a0)
}

// trackerPongInputEmbedded is Pong under the name of the field that
// trackerPongInput embeds.
type trackerPongInputEmbedded = Pong

// trackerPongInput implements Pong. Its methods return the results that
// the generators of the suite of Tracker draw for them, in turn.
type trackerPongInput struct {
	trackerPongInputEmbedded
	onPing func() Ping
}

// Ping calls the function of Ping.
func (in trackerPongInput) Ping() Ping {
	return in.onPing()
}

// trackerPairInputEmbedded is Pair under the name of the field that
// trackerPairInput embeds.
type trackerPairInputEmbedded = Pair

// trackerPairInput implements Pair. Its methods return the results that
// the generators of the suite of Tracker draw for them, in turn.
type trackerPairInput struct {
	trackerPairInputEmbedded
	onLeft  func() Leaf
	onRight func() Leaf
}

// Left calls the function of Left.
func (in trackerPairInput) Left() Leaf {
	return in.onLeft()
}

// Right calls the function of Right.
func (in trackerPairInput) Right() Leaf {
	return in.onRight()
}

// trackerRallyInputEmbedded is Rally under the name of the field that
// trackerRallyInput embeds.
type trackerRallyInputEmbedded = Rally

// trackerRallyInput implements Rally. Its methods return the results that
// the generators of the suite of Tracker draw for them, in turn.
type trackerRallyInput struct {
	trackerRallyInputEmbedded
	onServe func() Ping
}

// Serve calls the function of Serve.
func (in trackerRallyInput) Serve() Ping {
	return in.onServe()
}

// trackerLinkerInputEmbedded is Linker under the name of the field that
// trackerLinkerInput embeds.
type trackerLinkerInputEmbedded = Linker

// trackerLinkerInput implements Linker. Its methods return the results
// that the generators of the suite of Tracker draw for them, in turn.
type trackerLinkerInput struct {
	trackerLinkerInputEmbedded
	onLink func() Chain
}

// Link calls the function of Link.
func (in trackerLinkerInput) Link() Chain {
	return in.onLink()
}

// trackerOverEmbedded is Tracker under the name of the field that
// trackerOver embeds.
type trackerOverEmbedded = Tracker

// trackerOver is a subject of Tracker whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type trackerOver struct {
	trackerOverEmbedded
	spend func()
}

// Route calls spend and returns its zero results.
func (o trackerOver) Route(Node) (r0 int) {
	o.spend()
	return r0
}

// Follow calls spend and returns its zero results.
func (o trackerOver) Follow(Branch) (r0 int) {
	o.spend()
	return r0
}

// Walk calls spend and returns its zero results.
func (o trackerOver) Walk(Step) (r0 int) {
	o.spend()
	return r0
}

// Volley calls spend and returns its zero results.
func (o trackerOver) Volley(Ping) (r0 int) {
	o.spend()
	return r0
}

// Track calls spend and returns its zero results.
func (o trackerOver) Track(Labeler) (r0 int) {
	o.spend()
	return r0
}

// Pick calls spend and returns its zero results.
func (o trackerOver) Pick(Choice) (r0 int) {
	o.spend()
	return r0
}

// Play calls spend and returns its zero results.
func (o trackerOver) Play(Match) (r0 int) {
	o.spend()
	return r0
}

// Link calls spend and returns its zero results.
func (o trackerOver) Link(Links) (r0 int) {
	o.spend()
	return r0
}

// newTrackerOver returns the planted defect of a ceiling of Tracker, whose
// callables call spend.
func newTrackerOver(_ *TrackerHandle, spend func()) Tracker {
	return trackerOver{spend: spend}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Hire writes.
type (
	HireHandle = conform.Handle[HireFunc]
	HireCheck  = conform.Check[HireFunc]
	HireOption = conform.Option[HireFunc]
)

// HireFunc is the type of the function Hire. Each subject of the suite of
// Hire has it.
type HireFunc func(hire func() Labeler) int

// HireChecks contains the name of each callable of Hire, with the
// identities of its checks.
var HireChecks = hireCallable{"Hire"}

// hireCallable is the name of Hire, with the identities of its checks.
type hireCallable struct {
	conform.Callable[HireFunc]
}

// HireSuite returns the checks of Hire, and the calls and the budgets of
// its callables.
func HireSuite() conform.Suite[HireFunc] {
	return conform.Suite[HireFunc]{
		Declaration: "Hire",
		Calls: []conform.Call[HireFunc]{
			{
				Calls:    hireCalls,
				Over:     newHireOver,
				Callable: HireChecks.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) Labeler {
				labelerLabelResults0 := c.Draw(prop.List(prop.Of[Receipt]()), "labeler-label-0")
				labelerLabelResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, ErrTorn),
						prop.MinSize(len(labelerLabelResults0)),
						prop.MaxSize(len(labelerLabelResults0)),
					),
					"labeler-label-1",
				)
				var labelerLabelAt atomic.Int64
				return &hireLabelerInput{
					onLabel: func(string) (r0 Receipt, r1 error) {
						i := int(labelerLabelAt.Add(1) - 1)
						if i >= len(labelerLabelResults0) {
							r1 = ErrTorn
							return r0, r1
						}
						return labelerLabelResults0[i],
							labelerLabelResults1[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(*prop.Case) func() {
				return func() {}
			})),
		},
	}
}

// RunHireSuite applies the checks of Hire, and the checks that opts add,
// to subject.
func RunHireSuite(t *testing.T, subject HireFunc, opts ...HireOption) {
	t.Helper()
	RunHireSuiteFromHandle(t, hireSubjects(subject), opts...)
}

// RunHireSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunHireSuiteFromHandle(
	t *testing.T,
	newSubject func(h *HireHandle) HireFunc,
	opts ...HireOption,
) {
	t.Helper()
	conform.Run(t, HireSuite(), newSubject, opts...)
}

// FuzzHireSuite registers on f the fuzz target of the property checks of
// Hire on subject.
func FuzzHireSuite(f *testing.F, subject HireFunc, opts ...HireOption) {
	f.Helper()
	conform.Fuzz(f, HireSuite(), HireChecks.Callable, hireSubjects(subject), opts...)
}

// RunHireBudgets applies the allocation ceiling of Hire and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunHireBudgets(t *testing.T, subject HireFunc, opts ...HireOption) {
	t.Helper()
	conform.RunBudgets(t, HireSuite(), hireSubjects(subject), opts...)
}

// BenchHireSuite measures subject under the ceilings of the budget of
// Hire.
func BenchHireSuite(b conform.B, subject HireFunc, opts ...HireOption) {
	b.Helper()
	conform.Bench(b, HireSuite(), HireChecks.Callable, hireSubjects(subject), opts...)
}

// hireSubjects returns the factory of the subjects of Hire that returns
// subject for each handle.
func hireSubjects(subject HireFunc) func(*HireHandle) HireFunc {
	return func(*HireHandle) HireFunc { return subject }
}

// hireCalls generates the calls of Hire with the inputs that each case
// draws.
var hireCalls = prop.Composite(func(c *prop.Case) func(context.Context, HireFunc) error {
	hireResults := c.Draw(prop.List(prop.Of[Labeler]()), "hire")
	var hireAt atomic.Int64
	hire := func() (r0 Labeler) {
		i := int(hireAt.Add(1) - 1)
		if i >= len(hireResults) {
			return r0
		}
		return hireResults[i]
	}
	return func(_ context.Context, subject HireFunc) error {
		_ = subject(hire)
		return nil
	}
})

// hireLabelerInputEmbedded is Labeler under the name of the field that
// hireLabelerInput embeds.
type hireLabelerInputEmbedded = Labeler

// hireLabelerInput implements Labeler. Its methods return the results that
// the generators of the suite of Hire draw for them, in turn.
type hireLabelerInput struct {
	hireLabelerInputEmbedded
	onLabel func(a0 string) (Receipt, error)
}

// Label calls the function of Label.
func (in hireLabelerInput) Label(a0 string) (Receipt, error) {
	return in.onLabel(a0)
}

// newHireOver returns the planted defect of a ceiling of Hire, which calls
// spend and returns its zero results.
func newHireOver(_ *HireHandle, spend func()) HireFunc {
	return func(func() Labeler) (r0 int) {
		spend()
		return r0
	}
}
