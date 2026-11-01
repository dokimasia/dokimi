package drawn

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Tallies:
//
// Tallies.Merge/smoke: Merge returns without a panic for derived inputs.
//
// Tallies.Merge/smoke#generated: Merge returns without a panic for derived
// inputs.
//
// Tallies.Merge/deterministic: Merge returns the same value for the same
// inputs.
//
// Tallies.Merge/deterministic#generated: Merge returns the same value for
// the same inputs.
//
// Tallies.Merge/commutative: Merge returns one value for two inputs in
// either order.
//
// Tallies.Merge/associative: Merge returns one value for three inputs in
// either grouping.
//
// Tallies.Distance/smoke: Distance returns without a panic for derived
// inputs.
//
// Tallies.Distance/smoke#generated: Distance returns without a panic for
// derived inputs.
//
// Tallies.Distance/deterministic: Distance returns the same value for the
// same inputs.
//
// Tallies.Distance/deterministic#generated: Distance returns the same
// value for the same inputs.
//
// Tallies.Distance/commutative: Distance returns one value for two inputs
// in either order.
//
// Tallies.Each/smoke: Each returns without a panic for derived inputs.
//
// Tallies.Each/smoke#generated: Each returns without a panic for derived
// inputs.
//
// Tallies.Each/cancel: Each returns an error for a cancelled context.
//
// Tallies.Each/cancel#generated: Each returns an error for a cancelled
// context.
//
// Tallies.Save/smoke: Save returns without a panic for derived inputs.
//
// Tallies.Save/smoke#generated: Save returns without a panic for derived
// inputs.
//
// Tallies.Save/cancel: Save returns an error for a cancelled context.
//
// Tallies.Save/cancel#generated: Save returns an error for a cancelled
// context.
//
// Tallies.Limit/smoke: Limit returns without a panic for derived inputs.
//
// Tallies.Limit/smoke#generated: Limit returns without a panic for derived
// inputs.
//
// Tallies.Limit/cancel: Limit returns an error for a cancelled context.
//
// Tallies.Limit/cancel#generated: Limit returns an error for a cancelled
// context.
//
// Tallies.Count/smoke: Count returns without a panic for derived inputs.
//
// Tallies.Count/cancel: Count returns an error for a cancelled context.
//
// Tallies.Count/zero-on-error: Count returns the zero value beside an
// error.
//
// Tallies.Visit/smoke: Visit returns without a panic for derived inputs.
//
// Tallies.Visit/smoke#generated: Visit returns without a panic for derived
// inputs.
//
// Tallies.Visit/cancel: Visit returns an error for a cancelled context.
//
// Tallies.Visit/cancel#generated: Visit returns an error for a cancelled
// context.
//
// Tallies.Schedule/smoke: Schedule returns without a panic for derived
// inputs.
//
// Tallies.Schedule/smoke#generated: Schedule returns without a panic for
// derived inputs.
//
// Tallies.Schedule/cancel: Schedule returns an error for a cancelled
// context.
//
// Tallies.Schedule/cancel#generated: Schedule returns an error for a
// cancelled context.
//
// Tallies.Link/smoke: Link returns without a panic for derived inputs.
//
// Tallies.Link/smoke#generated: Link returns without a panic for derived
// inputs.
//
// Tallies.Link/cancel: Link returns an error for a cancelled context.
//
// Tallies.Link/cancel#generated: Link returns an error for a cancelled
// context.
//
// Tallies.Subscribe/smoke: Subscribe returns without a panic for derived
// inputs.
//
// Tallies.Subscribe/cancel: Subscribe returns an error for a cancelled
// context.

// The derived checks of Tallies that witness their laws at the shallow
// end:
//
// Tallies.Merge/smoke
//
// Tallies.Merge/deterministic
//
// Tallies.Distance/smoke
//
// Tallies.Distance/deterministic
//
// Tallies.Each/smoke
//
// Tallies.Save/smoke
//
// Tallies.Limit/smoke
//
// Tallies.Count/smoke
//
// Tallies.Visit/smoke
//
// Tallies.Schedule/smoke
//
// Tallies.Link/smoke
//
// Tallies.Subscribe/smoke

// The unclaimed classifications of Tallies:
//
// The shape:writer of Tallies.Each has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Tallies.Save has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Tallies.Limit has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Tallies.Visit has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Tallies.Schedule has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:writer of Tallies.Link has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Tallies.Subscribe has the stance silent. A write is
// an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Tallies writes.
type (
	TalliesHandle = conform.Handle[Tallies]
	TalliesCheck  = conform.Check[Tallies]
	TalliesOption = conform.Option[Tallies]
)

// TalliesChecks contains the name of each callable of Tallies, with the
// identities of its checks.
var TalliesChecks = struct {
	Merge     talliesMerge
	Distance  talliesDistance
	Each      talliesEach
	Save      talliesSave
	Limit     talliesLimit
	Count     talliesCount
	Visit     talliesVisit
	Schedule  talliesSchedule
	Link      talliesLink
	Subscribe talliesSubscribe
}{
	Merge:     talliesMerge{"Tallies.Merge"},
	Distance:  talliesDistance{"Tallies.Distance"},
	Each:      talliesEach{"Tallies.Each"},
	Save:      talliesSave{"Tallies.Save"},
	Limit:     talliesLimit{"Tallies.Limit"},
	Count:     talliesCount{"Tallies.Count"},
	Visit:     talliesVisit{"Tallies.Visit"},
	Schedule:  talliesSchedule{"Tallies.Schedule"},
	Link:      talliesLink{"Tallies.Link"},
	Subscribe: talliesSubscribe{"Tallies.Subscribe"},
}

// talliesMerge is the name of Tallies.Merge, with the identities of its
// checks.
type talliesMerge struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Merge returns without a
// panic for derived inputs.
func (c talliesMerge) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Merge returns
// without a panic for derived inputs.
func (c talliesMerge) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Deterministic returns the identity of the check that Merge returns the
// same value for the same inputs.
func (c talliesMerge) Deterministic() conform.ID[Tallies] {
	return c.Check("deterministic")
}

// DeterministicGenerated returns the identity of the check that Merge
// returns the same value for the same inputs.
func (c talliesMerge) DeterministicGenerated() conform.ID[Tallies] {
	return c.Check("deterministic#generated")
}

// Commutative returns the identity of the check that Merge returns one
// value for two inputs in either order.
func (c talliesMerge) Commutative() conform.ID[Tallies] {
	return c.Check("commutative")
}

// Associative returns the identity of the check that Merge returns one
// value for three inputs in either grouping.
func (c talliesMerge) Associative() conform.ID[Tallies] {
	return c.Check("associative")
}

// talliesDistance is the name of Tallies.Distance, with the identities of
// its checks.
type talliesDistance struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Distance returns without a
// panic for derived inputs.
func (c talliesDistance) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Distance returns
// without a panic for derived inputs.
func (c talliesDistance) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Deterministic returns the identity of the check that Distance returns
// the same value for the same inputs.
func (c talliesDistance) Deterministic() conform.ID[Tallies] {
	return c.Check("deterministic")
}

// DeterministicGenerated returns the identity of the check that Distance
// returns the same value for the same inputs.
func (c talliesDistance) DeterministicGenerated() conform.ID[Tallies] {
	return c.Check("deterministic#generated")
}

// Commutative returns the identity of the check that Distance returns one
// value for two inputs in either order.
func (c talliesDistance) Commutative() conform.ID[Tallies] {
	return c.Check("commutative")
}

// talliesEach is the name of Tallies.Each, with the identities of its
// checks.
type talliesEach struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Each returns without a
// panic for derived inputs.
func (c talliesEach) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Each returns
// without a panic for derived inputs.
func (c talliesEach) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Each returns an error for
// a cancelled context.
func (c talliesEach) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Each returns an
// error for a cancelled context.
func (c talliesEach) CancelGenerated() conform.ID[Tallies] {
	return c.Check("cancel#generated")
}

// talliesSave is the name of Tallies.Save, with the identities of its
// checks.
type talliesSave struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Save returns without a
// panic for derived inputs.
func (c talliesSave) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Save returns
// without a panic for derived inputs.
func (c talliesSave) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Save returns an error for
// a cancelled context.
func (c talliesSave) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Save returns an
// error for a cancelled context.
func (c talliesSave) CancelGenerated() conform.ID[Tallies] {
	return c.Check("cancel#generated")
}

// talliesLimit is the name of Tallies.Limit, with the identities of its
// checks.
type talliesLimit struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Limit returns without a
// panic for derived inputs.
func (c talliesLimit) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Limit returns
// without a panic for derived inputs.
func (c talliesLimit) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Limit returns an error for
// a cancelled context.
func (c talliesLimit) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Limit returns an
// error for a cancelled context.
func (c talliesLimit) CancelGenerated() conform.ID[Tallies] {
	return c.Check("cancel#generated")
}

// talliesCount is the name of Tallies.Count, with the identities of its
// checks.
type talliesCount struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Count returns without a
// panic for derived inputs.
func (c talliesCount) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Count returns an error for
// a cancelled context.
func (c talliesCount) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Count returns the
// zero value beside an error.
func (c talliesCount) ZeroOnError() conform.ID[Tallies] {
	return c.Check("zero-on-error")
}

// talliesVisit is the name of Tallies.Visit, with the identities of its
// checks.
type talliesVisit struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Visit returns without a
// panic for derived inputs.
func (c talliesVisit) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Visit returns
// without a panic for derived inputs.
func (c talliesVisit) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Visit returns an error for
// a cancelled context.
func (c talliesVisit) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Visit returns an
// error for a cancelled context.
func (c talliesVisit) CancelGenerated() conform.ID[Tallies] {
	return c.Check("cancel#generated")
}

// talliesSchedule is the name of Tallies.Schedule, with the identities of
// its checks.
type talliesSchedule struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Schedule returns without a
// panic for derived inputs.
func (c talliesSchedule) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Schedule returns
// without a panic for derived inputs.
func (c talliesSchedule) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Schedule returns an error
// for a cancelled context.
func (c talliesSchedule) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Schedule returns
// an error for a cancelled context.
func (c talliesSchedule) CancelGenerated() conform.ID[Tallies] {
	return c.Check("cancel#generated")
}

// talliesLink is the name of Tallies.Link, with the identities of its
// checks.
type talliesLink struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Link returns without a
// panic for derived inputs.
func (c talliesLink) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Link returns
// without a panic for derived inputs.
func (c talliesLink) SmokeGenerated() conform.ID[Tallies] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Link returns an error for
// a cancelled context.
func (c talliesLink) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Link returns an
// error for a cancelled context.
func (c talliesLink) CancelGenerated() conform.ID[Tallies] {
	return c.Check("cancel#generated")
}

// talliesSubscribe is the name of Tallies.Subscribe, with the identities
// of its checks.
type talliesSubscribe struct {
	conform.Callable[Tallies]
}

// Smoke returns the identity of the check that Subscribe returns without a
// panic for derived inputs.
func (c talliesSubscribe) Smoke() conform.ID[Tallies] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Subscribe returns an error
// for a cancelled context.
func (c talliesSubscribe) Cancel() conform.ID[Tallies] {
	return c.Check("cancel")
}

// TalliesSuite returns the checks of Tallies, and the calls and the
// budgets of its callables.
func TalliesSuite() conform.Suite[Tallies] {
	return conform.Suite[Tallies]{
		Declaration: "Tallies",
		Package:     "example.com/shop/drawn",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Tallies]{
			{
				ID:          TalliesChecks.Merge.Smoke(),
				Statement:   "Merge returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Merge(Tally{Items: 42}, Tally{Items: 42})
						return nil
					})
				},
				Defect: talliesMergePanics,
			},
			{
				ID:          TalliesChecks.Merge.SmokeGenerated(),
				Statement:   "Merge returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnA := c.Draw(prop.Of[Tally](), "a")
					drawnB := c.Draw(prop.Of[Tally](), "b")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Merge(drawnA, drawnB)
						return nil
					})
				},
				Defect: talliesMergePanics,
			},
			{
				ID:          TalliesChecks.Merge.Deterministic(),
				Statement:   "Merge returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Tallies]{TalliesChecks.Merge.Smoke()},
				Apply: func(h *TalliesHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (Tally, error) {
						r0 := h.Subject().Merge(Tally{Items: 42}, Tally{Items: 42})
						return r0, nil
					})
				},
				Defect: talliesMergeDrifts,
			},
			{
				ID:          TalliesChecks.Merge.DeterministicGenerated(),
				Statement:   "Merge returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Entails:     []conform.ID[Tallies]{TalliesChecks.Merge.SmokeGenerated()},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnA := c.Draw(prop.Of[Tally](), "a")
					drawnB := c.Draw(prop.Of[Tally](), "b")
					probe.RepeatAgrees(h, func(ctx context.Context) (Tally, error) {
						r0 := h.Subject().Merge(drawnA, drawnB)
						return r0, nil
					})
				},
				Defect: talliesMergeDrifts,
			},
			{
				ID:          TalliesChecks.Merge.Commutative(),
				Statement:   "Merge returns one value for two inputs in either order",
				Claim:       spec.Commutative,
				Probe:       spec.EitherOrderAgrees,
				Reason:      spec.OrderDependent,
				Law:         spec.CombiningIgnoresTheOrder,
				Machineries: []spec.MachineryID{spec.GeneratedInputs},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnA := c.Draw(prop.Of[Tally](), "a")
					drawnB := c.Draw(prop.Of[Tally](), "b")
					probe.EitherOrderAgrees(
						h,
						func(ctx context.Context, x, y Tally) (Tally, error) {
							r0 := h.Subject().Merge(x, y)
							return r0, nil
						},
						drawnA,
						drawnB,
					)
				},
				Defect: talliesMergeIgnoresTheSecondInput,
			},
			{
				ID:          TalliesChecks.Merge.Associative(),
				Statement:   "Merge returns one value for three inputs in either grouping",
				Claim:       spec.Associative,
				Probe:       spec.EitherGroupingAgrees,
				Reason:      spec.GroupingDependent,
				Law:         spec.CombiningIgnoresTheGrouping,
				Machineries: []spec.MachineryID{spec.GeneratedInputs},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnA := c.Draw(prop.Of[Tally](), "a")
					drawnB := c.Draw(
						prop.List(
							prop.Of[Tally](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"b",
					)
					probe.EitherGroupingAgrees(
						h,
						func(ctx context.Context, x, y Tally) (Tally, error) {
							r0 := h.Subject().Merge(x, y)
							return r0, nil
						},
						drawnA,
						drawnB[0],
						drawnB[1],
					)
				},
				Defect: talliesMergeKeepsItsFirstResult,
			},
			{
				ID:          TalliesChecks.Distance.Smoke(),
				Statement:   "Distance returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Distance(Point{X: 42}, Point{X: 42})
						return nil
					})
				},
				Defect: talliesDistancePanics,
			},
			{
				ID:          TalliesChecks.Distance.SmokeGenerated(),
				Statement:   "Distance returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnA := c.Draw(prop.Of[Point](), "a")
					drawnB := c.Draw(prop.Of[Point](), "b")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Distance(drawnA, drawnB)
						return nil
					})
				},
				Defect: talliesDistancePanics,
			},
			{
				ID:          TalliesChecks.Distance.Deterministic(),
				Statement:   "Distance returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Tallies]{TalliesChecks.Distance.Smoke()},
				Apply: func(h *TalliesHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Distance(Point{X: 42}, Point{X: 42})
						return r0, nil
					})
				},
				Defect: talliesDistanceDrifts,
			},
			{
				ID:          TalliesChecks.Distance.DeterministicGenerated(),
				Statement:   "Distance returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Entails:     []conform.ID[Tallies]{TalliesChecks.Distance.SmokeGenerated()},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnA := c.Draw(prop.Of[Point](), "a")
					drawnB := c.Draw(prop.Of[Point](), "b")
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Distance(drawnA, drawnB)
						return r0, nil
					})
				},
				Defect: talliesDistanceDrifts,
			},
			{
				ID:          TalliesChecks.Distance.Commutative(),
				Statement:   "Distance returns one value for two inputs in either order",
				Claim:       spec.Commutative,
				Probe:       spec.EitherOrderAgrees,
				Reason:      spec.OrderDependent,
				Law:         spec.CombiningIgnoresTheOrder,
				Machineries: []spec.MachineryID{spec.GeneratedInputs},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnA := c.Draw(prop.Of[Point](), "a")
					drawnB := c.Draw(prop.Of[Point](), "b")
					probe.EitherOrderAgrees(h, func(ctx context.Context, x, y Point) (int, error) {
						r0 := h.Subject().Distance(x, y)
						return r0, nil
					}, drawnA, drawnB)
				},
				Defect: talliesDistanceIgnoresTheSecondInput,
			},
			{
				ID:          TalliesChecks.Each.Smoke(),
				Statement:   "Each returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Each(ctx, func(sku string) error { return nil })
					})
				},
				Defect: talliesEachPanics,
			},
			{
				ID:          TalliesChecks.Each.SmokeGenerated(),
				Statement:   "Each returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					fnResults := c.Draw(
						prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
						"fn",
					)
					drawnFn := func() func(sku string) error {
						var fnAt atomic.Int64
						return func(string) (r0 error) {
							i := int(fnAt.Add(1) - 1)
							if i >= len(fnResults) {
								r0 = conform.ErrGenerated
								return r0
							}
							return fnResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Each(ctx, drawnFn())
					})
				},
				Defect: talliesEachPanics,
			},
			{
				ID:          TalliesChecks.Each.Cancel(),
				Statement:   "Each returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Each(ctx, func(sku string) error { return nil })
					})
				},
				Defect: talliesEachIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Each.CancelGenerated(),
				Statement:   "Each returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					fnResults := c.Draw(
						prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
						"fn",
					)
					drawnFn := func() func(sku string) error {
						var fnAt atomic.Int64
						return func(string) (r0 error) {
							i := int(fnAt.Add(1) - 1)
							if i >= len(fnResults) {
								r0 = conform.ErrGenerated
								return r0
							}
							return fnResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Each(ctx, drawnFn())
					})
				},
				Defect: talliesEachIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Save.Smoke(),
				Statement:   "Save returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Save(ctx, &Item{SKU: "SKU-a"})
					})
				},
				Defect: talliesSavePanics,
			},
			{
				ID:          TalliesChecks.Save.SmokeGenerated(),
				Statement:   "Save returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[Item](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Save(ctx, &drawnItem)
					})
				},
				Defect: talliesSavePanics,
			},
			{
				ID:          TalliesChecks.Save.Cancel(),
				Statement:   "Save returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Save(ctx, &Item{SKU: "SKU-a"})
					})
				},
				Defect: talliesSaveIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Save.CancelGenerated(),
				Statement:   "Save returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[Item](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Save(ctx, &drawnItem)
					})
				},
				Defect: talliesSaveIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Limit.Smoke(),
				Statement:   "Limit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Limit(ctx, int(7), string("sku-a"))
					})
				},
				Defect: talliesLimitPanics,
			},
			{
				ID:          TalliesChecks.Limit.SmokeGenerated(),
				Statement:   "Limit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Limit(ctx, int(7), drawnSku)
					})
				},
				Defect: talliesLimitPanics,
			},
			{
				ID:          TalliesChecks.Limit.Cancel(),
				Statement:   "Limit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Limit(ctx, int(7), string("sku-a"))
					})
				},
				Defect: talliesLimitIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Limit.CancelGenerated(),
				Statement:   "Limit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Limit(ctx, int(7), drawnSku)
					})
				},
				Defect: talliesLimitIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Count.Smoke(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: talliesCountPanics,
			},
			{
				ID:          TalliesChecks.Count.Cancel(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: talliesCountIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Count.ZeroOnError(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx)
					}, int(0))
				},
				Defect: talliesCountValueBeside,
			},
			{
				ID:          TalliesChecks.Visit.Smoke(),
				Statement:   "Visit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Visit(ctx, string("sku-a"), func(sku string) error { return nil })
					})
				},
				Defect: talliesVisitPanics,
			},
			{
				ID:          TalliesChecks.Visit.SmokeGenerated(),
				Statement:   "Visit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Visit(ctx, drawnSku, func(sku string) error { return nil })
					})
				},
				Defect: talliesVisitPanics,
			},
			{
				ID:          TalliesChecks.Visit.Cancel(),
				Statement:   "Visit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Visit(ctx, string("sku-a"), func(sku string) error { return nil })
					})
				},
				Defect: talliesVisitIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Visit.CancelGenerated(),
				Statement:   "Visit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Visit(ctx, drawnSku, func(sku string) error { return nil })
					})
				},
				Defect: talliesVisitIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Schedule.Smoke(),
				Statement:   "Schedule returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Schedule(ctx, string("sku-a"), Plan{Name: "Name-a"})
					})
				},
				Defect: talliesSchedulePanics,
			},
			{
				ID:          TalliesChecks.Schedule.SmokeGenerated(),
				Statement:   "Schedule returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Schedule(ctx, drawnSku, Plan{Name: "Name-a"})
					})
				},
				Defect: talliesSchedulePanics,
			},
			{
				ID:          TalliesChecks.Schedule.Cancel(),
				Statement:   "Schedule returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Schedule(ctx, string("sku-a"), Plan{Name: "Name-a"})
					})
				},
				Defect: talliesScheduleIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Schedule.CancelGenerated(),
				Statement:   "Schedule returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Schedule(ctx, drawnSku, Plan{Name: "Name-a"})
					})
				},
				Defect: talliesScheduleIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Link.Smoke(),
				Statement:   "Link returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Link(ctx, Node{SKU: "SKU-a"})
					})
				},
				Defect: talliesLinkPanics,
			},
			{
				ID:          TalliesChecks.Link.SmokeGenerated(),
				Statement:   "Link returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[Node](), "n")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Link(ctx, drawnN)
					})
				},
				Defect: talliesLinkPanics,
			},
			{
				ID:          TalliesChecks.Link.Cancel(),
				Statement:   "Link returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Link(ctx, Node{SKU: "SKU-a"})
					})
				},
				Defect: talliesLinkIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Link.CancelGenerated(),
				Statement:   "Link returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TalliesHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[Node](), "n")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Link(ctx, drawnN)
					})
				},
				Defect: talliesLinkIgnoresCancel,
			},
			{
				ID:          TalliesChecks.Subscribe.Smoke(),
				Statement:   "Subscribe returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TalliesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Subscribe(ctx, func(sku string) {})
					})
				},
				Defect: talliesSubscribePanics,
			},
			{
				ID:          TalliesChecks.Subscribe.Cancel(),
				Statement:   "Subscribe returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TalliesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Subscribe(ctx, func(sku string) {})
					})
				},
				Defect: talliesSubscribeIgnoresCancel,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Tallies.Each",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Tallies.Save",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Tallies.Limit",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Tallies.Visit",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Tallies.Schedule",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Tallies.Link",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Tallies.Subscribe",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Tallies]{
			{
				Calls:    talliesMergeCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Merge.Callable,
			},
			{
				Calls:    talliesDistanceCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Distance.Callable,
			},
			{
				Calls:    talliesEachCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Each.Callable,
			},
			{
				Calls:    talliesSaveCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Save.Callable,
			},
			{
				Calls:    talliesLimitCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Limit.Callable,
			},
			{
				Calls:    talliesCountCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Count.Callable,
			},
			{
				Calls:    talliesVisitCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Visit.Callable,
			},
			{
				Calls:    talliesScheduleCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Schedule.Callable,
			},
			{
				Calls:    talliesLinkCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Link.Callable,
			},
			{
				Calls:    talliesSubscribeCalls,
				Over:     newTalliesOver,
				Callable: TalliesChecks.Subscribe.Callable,
			},
		},
		Using: []prop.Option{
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
		},
	}
}

// RunTalliesSuite applies the checks of Tallies, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunTalliesSuite[T Tallies](t *testing.T, newSubject func() T, opts ...TalliesOption) {
	t.Helper()
	RunTalliesSuiteFromHandle(t, func(*TalliesHandle) T { return newSubject() }, opts...)
}

// RunTalliesSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTalliesSuiteFromHandle[T Tallies](
	t *testing.T,
	newSubject func(h *TalliesHandle) T,
	opts ...TalliesOption,
) {
	t.Helper()
	conform.Run(t, TalliesSuite(), talliesSubjects(newSubject), opts...)
}

// FuzzTalliesSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTalliesSuite[T Tallies](
	f *testing.F,
	callable conform.Callable[Tallies],
	newSubject func(h *TalliesHandle) T,
	opts ...TalliesOption,
) {
	f.Helper()
	conform.Fuzz(f, TalliesSuite(), callable, talliesSubjects(newSubject), opts...)
}

// RunTalliesBudgets applies the allocation ceilings of Tallies and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunTalliesBudgets[T Tallies](
	t *testing.T,
	newSubject func(h *TalliesHandle) T,
	opts ...TalliesOption,
) {
	t.Helper()
	conform.RunBudgets(t, TalliesSuite(), talliesSubjects(newSubject), opts...)
}

// BenchTalliesSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchTalliesSuite[T Tallies](
	b conform.B,
	callable conform.Callable[Tallies],
	newSubject func(h *TalliesHandle) T,
	opts ...TalliesOption,
) {
	b.Helper()
	conform.Bench(b, TalliesSuite(), callable, talliesSubjects(newSubject), opts...)
}

// talliesSubjects returns the factory of the subjects of Tallies that
// calls newSubject with the handle of each check.
func talliesSubjects[T Tallies](
	newSubject func(h *TalliesHandle) T,
) func(h *TalliesHandle) Tallies {
	return func(h *TalliesHandle) Tallies { return newSubject(h) }
}

// talliesMergeCalls generates the calls of Tallies.Merge with the inputs
// that each case draws.
var talliesMergeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	a := c.Draw(prop.Of[Tally](), "a")
	b := c.Draw(prop.Of[Tally](), "b")
	return func(_ context.Context, subject Tallies) error {
		_ = subject.Merge(a, b)
		return nil
	}
})

// talliesDistanceCalls generates the calls of Tallies.Distance with the
// inputs that each case draws.
var talliesDistanceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	a := c.Draw(prop.Of[Point](), "a")
	b := c.Draw(prop.Of[Point](), "b")
	return func(_ context.Context, subject Tallies) error {
		_ = subject.Distance(a, b)
		return nil
	}
})

// talliesEachCalls generates the calls of Tallies.Each with the inputs
// that each case draws.
var talliesEachCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	fnResults := c.Draw(prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)), "fn")
	var fnAt atomic.Int64
	fn := func(string) (r0 error) {
		i := int(fnAt.Add(1) - 1)
		if i >= len(fnResults) {
			r0 = conform.ErrGenerated
			return r0
		}
		return fnResults[i]
	}
	return func(ctx context.Context, subject Tallies) error {
		return subject.Each(ctx, fn)
	}
})

// talliesSaveCalls generates the calls of Tallies.Save with the inputs
// that each case draws.
var talliesSaveCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	item := c.Draw(prop.Of[Item](), "item")
	return func(ctx context.Context, subject Tallies) error {
		return subject.Save(ctx, &item)
	}
})

// talliesLimitCalls generates the calls of Tallies.Limit with the inputs
// that each case draws.
var talliesLimitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	n := int(7)
	sku := c.Draw(prop.Of[string](), "sku")
	return func(ctx context.Context, subject Tallies) error {
		return subject.Limit(ctx, n, sku)
	}
})

// talliesCountCalls generates the calls of Tallies.Count with the inputs
// that each case draws.
var talliesCountCalls = prop.Composite(func(*prop.Case) func(context.Context, Tallies) error {
	return func(ctx context.Context, subject Tallies) error {
		_, err := subject.Count(ctx)
		return err
	}
})

// talliesVisitCalls generates the calls of Tallies.Visit with the inputs
// that each case draws.
var talliesVisitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	sku := c.Draw(prop.Of[string](), "sku")
	return func(ctx context.Context, subject Tallies) error {
		return subject.Visit(ctx, sku)
	}
})

// talliesScheduleCalls generates the calls of Tallies.Schedule with the
// inputs that each case draws.
var talliesScheduleCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	sku := c.Draw(prop.Of[string](), "sku")
	p := c.Draw(prop.Of[Plan](), "p")
	return func(ctx context.Context, subject Tallies) error {
		return subject.Schedule(ctx, sku, p)
	}
})

// talliesLinkCalls generates the calls of Tallies.Link with the inputs
// that each case draws.
var talliesLinkCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tallies) error {
	n := c.Draw(prop.Of[Node](), "n")
	return func(ctx context.Context, subject Tallies) error {
		return subject.Link(ctx, n)
	}
})

// talliesSubscribeCalls generates the calls of Tallies.Subscribe with the
// inputs that each case draws.
var talliesSubscribeCalls = prop.Composite(func(*prop.Case) func(context.Context, Tallies) error {
	fn := func(string) {}
	return func(ctx context.Context, subject Tallies) error {
		return subject.Subscribe(ctx, fn)
	}
})

// talliesOverEmbedded is Tallies under the name of the field that
// talliesOver embeds.
type talliesOverEmbedded = Tallies

// talliesOver is a subject of Tallies whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type talliesOver struct {
	talliesOverEmbedded
	spend func()
}

// Merge calls spend and returns its zero results.
func (o talliesOver) Merge(Tally, Tally) (r0 Tally) {
	o.spend()
	return r0
}

// Distance calls spend and returns its zero results.
func (o talliesOver) Distance(Point, Point) (r0 int) {
	o.spend()
	return r0
}

// Each calls spend and returns its zero results.
func (o talliesOver) Each(context.Context, func(sku string) error) (r0 error) {
	o.spend()
	return r0
}

// Save calls spend and returns its zero results.
func (o talliesOver) Save(context.Context, *Item) (r0 error) {
	o.spend()
	return r0
}

// Limit calls spend and returns its zero results.
func (o talliesOver) Limit(context.Context, int, string) (r0 error) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o talliesOver) Count(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// Visit calls spend and returns its zero results.
func (o talliesOver) Visit(context.Context, string, ...func(sku string) error) (r0 error) {
	o.spend()
	return r0
}

// Schedule calls spend and returns its zero results.
func (o talliesOver) Schedule(context.Context, string, Plan) (r0 error) {
	o.spend()
	return r0
}

// Link calls spend and returns its zero results.
func (o talliesOver) Link(context.Context, Node) (r0 error) {
	o.spend()
	return r0
}

// Subscribe calls spend and returns its zero results.
func (o talliesOver) Subscribe(context.Context, func(sku string)) (r0 error) {
	o.spend()
	return r0
}

// newTalliesOver returns the planted defect of a ceiling of Tallies, whose
// callables call spend.
func newTalliesOver(_ *TalliesHandle, spend func()) Tallies {
	return talliesOver{spend: spend}
}

// talliesMergePanicsDefect is the planted defect of a Tallies whose Merge
// panics.
type talliesMergePanicsDefect struct {
	talliesOverEmbedded
}

// Merge is the override of the planted defect.
func (*talliesMergePanicsDefect) Merge(p0, p1 Tally) (r0 Tally) {
	panic("the planted defect of a Tallies whose Merge panics")
}

// talliesMergePanics returns the planted defect of a Tallies whose Merge
// panics over a fresh subject.
func talliesMergePanics(h *TalliesHandle) Tallies {
	return &talliesMergePanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesMergeDriftsDefect is the planted defect of a Tallies whose Merge
// returns another value on the repeat.
type talliesMergeDriftsDefect struct {
	talliesOverEmbedded
	calls int
}

// Merge is the override of the planted defect.
func (d *talliesMergeDriftsDefect) Merge(p0, p1 Tally) (r0 Tally) {
	d.calls++
	if d.calls > 1 {
		r0 = Tally{Items: 7}
	}
	return r0
}

// talliesMergeDrifts returns the planted defect of a Tallies whose Merge
// returns another value on the repeat over a fresh subject.
func talliesMergeDrifts(h *TalliesHandle) Tallies {
	return &talliesMergeDriftsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesMergeIgnoresTheSecondInputDefect is the planted defect of a
// Tallies whose Merge ignores its second input.
type talliesMergeIgnoresTheSecondInputDefect struct {
	talliesOverEmbedded
}

// Merge is the override of the planted defect.
func (*talliesMergeIgnoresTheSecondInputDefect) Merge(p0, p1 Tally) (r0 Tally) {
	r0 = p0
	return r0
}

// talliesMergeIgnoresTheSecondInput returns the planted defect of a
// Tallies whose Merge ignores its second input over a fresh subject.
func talliesMergeIgnoresTheSecondInput(h *TalliesHandle) Tallies {
	return &talliesMergeIgnoresTheSecondInputDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesMergeKeepsItsFirstResultDefect is the planted defect of a Tallies
// whose Merge keeps the value of its first call.
type talliesMergeKeepsItsFirstResultDefect struct {
	talliesOverEmbedded
	first *Tally
}

// Merge is the override of the planted defect.
func (d *talliesMergeKeepsItsFirstResultDefect) Merge(p0, p1 Tally) (r0 Tally) {
	if d.first != nil && reflect.DeepEqual(p0, *d.first) {
		r0 = p0
		return r0
	}
	r0 = d.talliesOverEmbedded.Merge(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	return r0
}

// talliesMergeKeepsItsFirstResult returns the planted defect of a Tallies
// whose Merge keeps the value of its first call over a fresh subject.
func talliesMergeKeepsItsFirstResult(h *TalliesHandle) Tallies {
	return &talliesMergeKeepsItsFirstResultDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesDistancePanicsDefect is the planted defect of a Tallies whose
// Distance panics.
type talliesDistancePanicsDefect struct {
	talliesOverEmbedded
}

// Distance is the override of the planted defect.
func (*talliesDistancePanicsDefect) Distance(p0, p1 Point) (r0 int) {
	panic("the planted defect of a Tallies whose Distance panics")
}

// talliesDistancePanics returns the planted defect of a Tallies whose
// Distance panics over a fresh subject.
func talliesDistancePanics(h *TalliesHandle) Tallies {
	return &talliesDistancePanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesDistanceDriftsDefect is the planted defect of a Tallies whose
// Distance returns another value on the repeat.
type talliesDistanceDriftsDefect struct {
	talliesOverEmbedded
	calls int
}

// Distance is the override of the planted defect.
func (d *talliesDistanceDriftsDefect) Distance(p0, p1 Point) (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// talliesDistanceDrifts returns the planted defect of a Tallies whose
// Distance returns another value on the repeat over a fresh subject.
func talliesDistanceDrifts(h *TalliesHandle) Tallies {
	return &talliesDistanceDriftsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesDistanceIgnoresTheSecondInputDefect is the planted defect of a
// Tallies whose Distance ignores its second input.
type talliesDistanceIgnoresTheSecondInputDefect struct {
	talliesOverEmbedded
}

// Distance is the override of the planted defect.
func (d *talliesDistanceIgnoresTheSecondInputDefect) Distance(p0, p1 Point) (r0 int) {
	return d.talliesOverEmbedded.Distance(p0, *new(Point))
}

// talliesDistanceIgnoresTheSecondInput returns the planted defect of a
// Tallies whose Distance ignores its second input over a fresh subject.
func talliesDistanceIgnoresTheSecondInput(h *TalliesHandle) Tallies {
	return &talliesDistanceIgnoresTheSecondInputDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesEachPanicsDefect is the planted defect of a Tallies whose Each
// panics.
type talliesEachPanicsDefect struct {
	talliesOverEmbedded
}

// Each is the override of the planted defect.
func (*talliesEachPanicsDefect) Each(p0 context.Context, p1 func(sku string) error) (r0 error) {
	panic("the planted defect of a Tallies whose Each panics")
}

// talliesEachPanics returns the planted defect of a Tallies whose Each
// panics over a fresh subject.
func talliesEachPanics(h *TalliesHandle) Tallies {
	return &talliesEachPanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesEachIgnoresCancelDefect is the planted defect of a Tallies whose
// Each ignores the context.
type talliesEachIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Each is the override of the planted defect.
func (d *talliesEachIgnoresCancelDefect) Each(
	p0 context.Context,
	p1 func(sku string) error,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.talliesOverEmbedded.Each(p0, p1)
}

// talliesEachIgnoresCancel returns the planted defect of a Tallies whose
// Each ignores the context over a fresh subject.
func talliesEachIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesEachIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesSavePanicsDefect is the planted defect of a Tallies whose Save
// panics.
type talliesSavePanicsDefect struct {
	talliesOverEmbedded
}

// Save is the override of the planted defect.
func (*talliesSavePanicsDefect) Save(p0 context.Context, p1 *Item) (r0 error) {
	panic("the planted defect of a Tallies whose Save panics")
}

// talliesSavePanics returns the planted defect of a Tallies whose Save
// panics over a fresh subject.
func talliesSavePanics(h *TalliesHandle) Tallies {
	return &talliesSavePanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesSaveIgnoresCancelDefect is the planted defect of a Tallies whose
// Save ignores the context.
type talliesSaveIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Save is the override of the planted defect.
func (d *talliesSaveIgnoresCancelDefect) Save(p0 context.Context, p1 *Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.talliesOverEmbedded.Save(p0, p1)
}

// talliesSaveIgnoresCancel returns the planted defect of a Tallies whose
// Save ignores the context over a fresh subject.
func talliesSaveIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesSaveIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesLimitPanicsDefect is the planted defect of a Tallies whose Limit
// panics.
type talliesLimitPanicsDefect struct {
	talliesOverEmbedded
}

// Limit is the override of the planted defect.
func (*talliesLimitPanicsDefect) Limit(p0 context.Context, p1 int, p2 string) (r0 error) {
	panic("the planted defect of a Tallies whose Limit panics")
}

// talliesLimitPanics returns the planted defect of a Tallies whose Limit
// panics over a fresh subject.
func talliesLimitPanics(h *TalliesHandle) Tallies {
	return &talliesLimitPanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesLimitIgnoresCancelDefect is the planted defect of a Tallies whose
// Limit ignores the context.
type talliesLimitIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Limit is the override of the planted defect.
func (d *talliesLimitIgnoresCancelDefect) Limit(p0 context.Context, p1 int, p2 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.talliesOverEmbedded.Limit(p0, p1, p2)
}

// talliesLimitIgnoresCancel returns the planted defect of a Tallies whose
// Limit ignores the context over a fresh subject.
func talliesLimitIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesLimitIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesCountPanicsDefect is the planted defect of a Tallies whose Count
// panics.
type talliesCountPanicsDefect struct {
	talliesOverEmbedded
}

// Count is the override of the planted defect.
func (*talliesCountPanicsDefect) Count(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Tallies whose Count panics")
}

// talliesCountPanics returns the planted defect of a Tallies whose Count
// panics over a fresh subject.
func talliesCountPanics(h *TalliesHandle) Tallies {
	return &talliesCountPanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesCountIgnoresCancelDefect is the planted defect of a Tallies whose
// Count ignores the context.
type talliesCountIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Count is the override of the planted defect.
func (d *talliesCountIgnoresCancelDefect) Count(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.talliesOverEmbedded.Count(p0)
}

// talliesCountIgnoresCancel returns the planted defect of a Tallies whose
// Count ignores the context over a fresh subject.
func talliesCountIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesCountIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesCountValueBesideDefect is the planted defect of a Tallies whose
// Count returns a value beside an error.
type talliesCountValueBesideDefect struct {
	talliesOverEmbedded
}

// Count is the override of the planted defect.
func (d *talliesCountValueBesideDefect) Count(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.talliesOverEmbedded.Count(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// talliesCountValueBeside returns the planted defect of a Tallies whose
// Count returns a value beside an error over a fresh subject.
func talliesCountValueBeside(h *TalliesHandle) Tallies {
	return &talliesCountValueBesideDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesVisitPanicsDefect is the planted defect of a Tallies whose Visit
// panics.
type talliesVisitPanicsDefect struct {
	talliesOverEmbedded
}

// Visit is the override of the planted defect.
func (*talliesVisitPanicsDefect) Visit(
	p0 context.Context,
	p1 string,
	p2 ...func(sku string) error,
) (r0 error) {
	panic("the planted defect of a Tallies whose Visit panics")
}

// talliesVisitPanics returns the planted defect of a Tallies whose Visit
// panics over a fresh subject.
func talliesVisitPanics(h *TalliesHandle) Tallies {
	return &talliesVisitPanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesVisitIgnoresCancelDefect is the planted defect of a Tallies whose
// Visit ignores the context.
type talliesVisitIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Visit is the override of the planted defect.
func (d *talliesVisitIgnoresCancelDefect) Visit(
	p0 context.Context,
	p1 string,
	p2 ...func(sku string) error,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.talliesOverEmbedded.Visit(p0, p1, p2...)
}

// talliesVisitIgnoresCancel returns the planted defect of a Tallies whose
// Visit ignores the context over a fresh subject.
func talliesVisitIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesVisitIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesSchedulePanicsDefect is the planted defect of a Tallies whose
// Schedule panics.
type talliesSchedulePanicsDefect struct {
	talliesOverEmbedded
}

// Schedule is the override of the planted defect.
func (*talliesSchedulePanicsDefect) Schedule(p0 context.Context, p1 string, p2 Plan) (r0 error) {
	panic("the planted defect of a Tallies whose Schedule panics")
}

// talliesSchedulePanics returns the planted defect of a Tallies whose
// Schedule panics over a fresh subject.
func talliesSchedulePanics(h *TalliesHandle) Tallies {
	return &talliesSchedulePanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesScheduleIgnoresCancelDefect is the planted defect of a Tallies
// whose Schedule ignores the context.
type talliesScheduleIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Schedule is the override of the planted defect.
func (d *talliesScheduleIgnoresCancelDefect) Schedule(
	p0 context.Context,
	p1 string,
	p2 Plan,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.talliesOverEmbedded.Schedule(p0, p1, p2)
}

// talliesScheduleIgnoresCancel returns the planted defect of a Tallies
// whose Schedule ignores the context over a fresh subject.
func talliesScheduleIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesScheduleIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesLinkPanicsDefect is the planted defect of a Tallies whose Link
// panics.
type talliesLinkPanicsDefect struct {
	talliesOverEmbedded
}

// Link is the override of the planted defect.
func (*talliesLinkPanicsDefect) Link(p0 context.Context, p1 Node) (r0 error) {
	panic("the planted defect of a Tallies whose Link panics")
}

// talliesLinkPanics returns the planted defect of a Tallies whose Link
// panics over a fresh subject.
func talliesLinkPanics(h *TalliesHandle) Tallies {
	return &talliesLinkPanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesLinkIgnoresCancelDefect is the planted defect of a Tallies whose
// Link ignores the context.
type talliesLinkIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Link is the override of the planted defect.
func (d *talliesLinkIgnoresCancelDefect) Link(p0 context.Context, p1 Node) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.talliesOverEmbedded.Link(p0, p1)
}

// talliesLinkIgnoresCancel returns the planted defect of a Tallies whose
// Link ignores the context over a fresh subject.
func talliesLinkIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesLinkIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesSubscribePanicsDefect is the planted defect of a Tallies whose
// Subscribe panics.
type talliesSubscribePanicsDefect struct {
	talliesOverEmbedded
}

// Subscribe is the override of the planted defect.
func (*talliesSubscribePanicsDefect) Subscribe(p0 context.Context, p1 func(sku string)) (r0 error) {
	panic("the planted defect of a Tallies whose Subscribe panics")
}

// talliesSubscribePanics returns the planted defect of a Tallies whose
// Subscribe panics over a fresh subject.
func talliesSubscribePanics(h *TalliesHandle) Tallies {
	return &talliesSubscribePanicsDefect{talliesOverEmbedded: h.Fresh()}
}

// talliesSubscribeIgnoresCancelDefect is the planted defect of a Tallies
// whose Subscribe ignores the context.
type talliesSubscribeIgnoresCancelDefect struct {
	talliesOverEmbedded
}

// Subscribe is the override of the planted defect.
func (d *talliesSubscribeIgnoresCancelDefect) Subscribe(
	p0 context.Context,
	p1 func(sku string),
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.talliesOverEmbedded.Subscribe(p0, p1)
}

// talliesSubscribeIgnoresCancel returns the planted defect of a Tallies
// whose Subscribe ignores the context over a fresh subject.
func talliesSubscribeIgnoresCancel(h *TalliesHandle) Tallies {
	return &talliesSubscribeIgnoresCancelDefect{talliesOverEmbedded: h.Fresh()}
}

// The derived checks of Notes:
//
// Notes.Put/smoke: Put returns without a panic for derived inputs.
//
// Notes.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Notes.Put/cancel: Put returns an error for a cancelled context.
//
// Notes.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Notes.Put/updater: a second Put for one key replaces the first.
//
// Notes.Put/updater#generated: a second Put for one key replaces the
// first.
//
// Notes.Get/smoke: Get returns without a panic for derived inputs.
//
// Notes.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Notes.Get/cancel: Get returns an error for a cancelled context.
//
// Notes.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Notes.Get/zero-on-error: Get returns the zero value beside an error.
//
// Notes.Get/zero-on-error#generated: Get returns the zero value beside an
// error.

// The derived checks of Notes that witness their laws at the shallow end:
//
// Notes.Put/smoke
//
// Notes.Get/smoke

// The unclaimed classifications of Notes:
//
// The shape:composite-writer of Notes.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Notes writes.
type (
	NotesHandle = conform.Handle[Notes]
	NotesCheck  = conform.Check[Notes]
	NotesOption = conform.Option[Notes]
)

// NotesChecks contains the name of each callable of Notes, with the
// identities of its checks.
var NotesChecks = struct {
	Put notesPut
	Get notesGet
}{
	Put: notesPut{"Notes.Put"},
	Get: notesGet{"Notes.Get"},
}

// notesPut is the name of Notes.Put, with the identities of its checks.
type notesPut struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c notesPut) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c notesPut) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c notesPut) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c notesPut) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// Updater returns the identity of the check that a second Put for one key
// replaces the first.
func (c notesPut) Updater() conform.ID[Notes] {
	return c.Check("updater")
}

// UpdaterGenerated returns the identity of the check that a second Put for
// one key replaces the first.
func (c notesPut) UpdaterGenerated() conform.ID[Notes] {
	return c.Check("updater#generated")
}

// notesGet is the name of Notes.Get, with the identities of its checks.
type notesGet struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c notesGet) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c notesGet) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c notesGet) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c notesGet) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c notesGet) ZeroOnError() conform.ID[Notes] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c notesGet) ZeroOnErrorGenerated() conform.ID[Notes] {
	return c.Check("zero-on-error#generated")
}

// NotesSuite returns the checks of Notes, and the calls and the budgets of
// its callables.
func NotesSuite() conform.Suite[Notes] {
	return conform.Suite[Notes]{
		Declaration: "Notes",
		Package:     "example.com/shop/drawn",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Notes]{
			{
				ID:          NotesChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("shelf-a"), FirstNote())
					})
				},
				Defect: notesPutPanics,
			},
			{
				ID:          NotesChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnShelf := c.Draw(prop.Of[string](), "shelf")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnShelf, FirstNote())
					})
				},
				Defect: notesPutPanics,
			},
			{
				ID:          NotesChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("shelf-a"), FirstNote())
					})
				},
				Defect: notesPutIgnoresCancel,
			},
			{
				ID:          NotesChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnShelf := c.Draw(prop.Of[string](), "shelf")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnShelf, FirstNote())
					})
				},
				Defect: notesPutIgnoresCancel,
			},
			{
				ID:          NotesChecks.Put.Updater(),
				Statement:   "a second Put for one key replaces the first",
				Claim:       spec.Updater,
				Probe:       spec.SecondWriteReplaces,
				Reason:      spec.UnreplacedWrite,
				Law:         spec.SecondWriteReplacesTheFirst,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.SecondWriteReplaces(h, func() func(ctx context.Context) (Note, error) {
						calls := 0
						return func(ctx context.Context) (Note, error) {
							calls++
							w := FirstNote()
							w.Text = string("Text-a")
							if calls > 1 {
								w.Text = string("Text-b")
							}
							err := h.Subject().Put(ctx, string("shelf-a"), w)
							return w, err
						}
					}(), func(ctx context.Context) (Note, error) {
						return h.Subject().Get(ctx, FirstNote().SKU)
					})
				},
				Defect: notesPutStoresNothing,
			},
			{
				ID:          NotesChecks.Put.UpdaterGenerated(),
				Statement:   "a second Put for one key replaces the first",
				Claim:       spec.Updater,
				Probe:       spec.SecondWriteReplaces,
				Reason:      spec.UnreplacedWrite,
				Law:         spec.SecondWriteReplacesTheFirst,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnShelf := c.Draw(prop.Of[string](), "shelf")
					probe.SecondWriteReplaces(h, func() func(ctx context.Context) (Note, error) {
						calls := 0
						return func(ctx context.Context) (Note, error) {
							calls++
							w := FirstNote()
							w.Text = string("Text-a")
							if calls > 1 {
								w.Text = string("Text-b")
							}
							err := h.Subject().Put(ctx, drawnShelf, w)
							return w, err
						}
					}(), func(ctx context.Context) (Note, error) {
						return h.Subject().Get(ctx, FirstNote().SKU)
					})
				},
				Defect: notesPutStoresNothing,
			},
			{
				ID:          NotesChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("sku-a"))
						return err
					})
				},
				Defect: notesGetPanics,
			},
			{
				ID:          NotesChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnSku)
						return err
					})
				},
				Defect: notesGetPanics,
			},
			{
				ID:          NotesChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("sku-a"))
						return err
					})
				},
				Defect: notesGetIgnoresCancel,
			},
			{
				ID:          NotesChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnSku)
						return err
					})
				},
				Defect: notesGetIgnoresCancel,
			},
			{
				ID:          NotesChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Note, error) {
						return h.Subject().Get(ctx, string("sku-a"))
					}, Note{})
				},
				Defect: notesGetValueBeside,
			},
			{
				ID:          NotesChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnSku := c.Draw(prop.Of[string](), "sku")
					probe.ZeroBesideError(h, func(ctx context.Context) (Note, error) {
						return h.Subject().Get(ctx, drawnSku)
					}, Note{})
				},
				Defect: notesGetValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Notes.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Notes]{
			{
				Calls:    notesPutCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Put.Callable,
			},
			{
				Calls:    notesGetCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Get.Callable,
			},
		},
	}
}

// RunNotesSuite applies the checks of Notes, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunNotesSuite[T Notes](t *testing.T, newSubject func() T, opts ...NotesOption) {
	t.Helper()
	RunNotesSuiteFromHandle(t, func(*NotesHandle) T { return newSubject() }, opts...)
}

// RunNotesSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunNotesSuiteFromHandle[T Notes](
	t *testing.T,
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	t.Helper()
	conform.Run(t, NotesSuite(), notesSubjects(newSubject), opts...)
}

// FuzzNotesSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzNotesSuite[T Notes](
	f *testing.F,
	callable conform.Callable[Notes],
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	f.Helper()
	conform.Fuzz(f, NotesSuite(), callable, notesSubjects(newSubject), opts...)
}

// RunNotesBudgets applies the allocation ceilings of Notes and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunNotesBudgets[T Notes](
	t *testing.T,
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	t.Helper()
	conform.RunBudgets(t, NotesSuite(), notesSubjects(newSubject), opts...)
}

// BenchNotesSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchNotesSuite[T Notes](
	b conform.B,
	callable conform.Callable[Notes],
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	b.Helper()
	conform.Bench(b, NotesSuite(), callable, notesSubjects(newSubject), opts...)
}

// notesSubjects returns the factory of the subjects of Notes that calls
// newSubject with the handle of each check.
func notesSubjects[T Notes](newSubject func(h *NotesHandle) T) func(h *NotesHandle) Notes {
	return func(h *NotesHandle) Notes { return newSubject(h) }
}

// notesPutCalls generates the calls of Notes.Put with the inputs that each
// case draws.
var notesPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	shelf := c.Draw(prop.Of[string](), "shelf")
	n := Note(FirstNote())
	return func(ctx context.Context, subject Notes) error {
		return subject.Put(ctx, shelf, n)
	}
})

// notesGetCalls generates the calls of Notes.Get with the inputs that each
// case draws.
var notesGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	sku := c.Draw(prop.Of[string](), "sku")
	return func(ctx context.Context, subject Notes) error {
		_, err := subject.Get(ctx, sku)
		return err
	}
})

// notesOverEmbedded is Notes under the name of the field that notesOver
// embeds.
type notesOverEmbedded = Notes

// notesOver is a subject of Notes whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type notesOver struct {
	notesOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o notesOver) Put(context.Context, string, Note) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o notesOver) Get(context.Context, string) (r0 Note, r1 error) {
	o.spend()
	return r0, r1
}

// newNotesOver returns the planted defect of a ceiling of Notes, whose
// callables call spend.
func newNotesOver(_ *NotesHandle, spend func()) Notes {
	return notesOver{spend: spend}
}

// notesPutPanicsDefect is the planted defect of a Notes whose Put panics.
type notesPutPanicsDefect struct {
	notesOverEmbedded
}

// Put is the override of the planted defect.
func (*notesPutPanicsDefect) Put(p0 context.Context, p1 string, p2 Note) (r0 error) {
	panic("the planted defect of a Notes whose Put panics")
}

// notesPutPanics returns the planted defect of a Notes whose Put panics
// over a fresh subject.
func notesPutPanics(h *NotesHandle) Notes {
	return &notesPutPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesPutIgnoresCancelDefect is the planted defect of a Notes whose Put
// ignores the context.
type notesPutIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Put is the override of the planted defect.
func (d *notesPutIgnoresCancelDefect) Put(p0 context.Context, p1 string, p2 Note) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Put(p0, p1, p2)
}

// notesPutIgnoresCancel returns the planted defect of a Notes whose Put
// ignores the context over a fresh subject.
func notesPutIgnoresCancel(h *NotesHandle) Notes {
	return &notesPutIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesPutStoresNothingDefect is the planted defect of a Notes whose Put
// keeps the first write.
type notesPutStoresNothingDefect struct {
	notesOverEmbedded
}

// Put is the override of the planted defect.
func (*notesPutStoresNothingDefect) Put(p0 context.Context, p1 string, p2 Note) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// notesPutStoresNothing returns the planted defect of a Notes whose Put
// keeps the first write over a fresh subject.
func notesPutStoresNothing(h *NotesHandle) Notes {
	return &notesPutStoresNothingDefect{notesOverEmbedded: h.Fresh()}
}

// notesGetPanicsDefect is the planted defect of a Notes whose Get panics.
type notesGetPanicsDefect struct {
	notesOverEmbedded
}

// Get is the override of the planted defect.
func (*notesGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Note, r1 error) {
	panic("the planted defect of a Notes whose Get panics")
}

// notesGetPanics returns the planted defect of a Notes whose Get panics
// over a fresh subject.
func notesGetPanics(h *NotesHandle) Notes {
	return &notesGetPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesGetIgnoresCancelDefect is the planted defect of a Notes whose Get
// ignores the context.
type notesGetIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Get is the override of the planted defect.
func (d *notesGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Note, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.notesOverEmbedded.Get(p0, p1)
}

// notesGetIgnoresCancel returns the planted defect of a Notes whose Get
// ignores the context over a fresh subject.
func notesGetIgnoresCancel(h *NotesHandle) Notes {
	return &notesGetIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesGetValueBesideDefect is the planted defect of a Notes whose Get
// returns a value beside an error.
type notesGetValueBesideDefect struct {
	notesOverEmbedded
}

// Get is the override of the planted defect.
func (d *notesGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Note, r1 error) {
	r0, r1 = d.notesOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Note{SKU: "SKU-b"}
	}
	return r0, r1
}

// notesGetValueBeside returns the planted defect of a Notes whose Get
// returns a value beside an error over a fresh subject.
func notesGetValueBeside(h *NotesHandle) Notes {
	return &notesGetValueBesideDefect{notesOverEmbedded: h.Fresh()}
}

// The derived checks of Hooked:
//
// Hooked.Set/smoke: Set returns without a panic for derived inputs.
//
// Hooked.Set/smoke#generated: Set returns without a panic for derived
// inputs.
//
// Hooked.Set/cancel: Set returns an error for a cancelled context.
//
// Hooked.Set/cancel#generated: Set returns an error for a cancelled
// context.
//
// Hooked.Set/partition: Get does not return in one partition what Set
// wrote in another.
//
// Hooked.Set/partition#generated: Get does not return in one partition
// what Set wrote in another.
//
// Hooked.Get/smoke: Get returns without a panic for derived inputs.
//
// Hooked.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Hooked.Get/cancel: Get returns an error for a cancelled context.
//
// Hooked.Get/cancel#generated: Get returns an error for a cancelled
// context.

// The derived checks of Hooked that witness their laws at the shallow end:
//
// Hooked.Set/smoke
//
// Hooked.Get/smoke

// The unclaimed classifications of Hooked:
//
// The shape:multi-arg-writer of Hooked.Set has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Hooked writes.
type (
	HookedHandle = conform.Handle[Hooked]
	HookedCheck  = conform.Check[Hooked]
	HookedOption = conform.Option[Hooked]
)

// HookedChecks contains the name of each callable of Hooked, with the
// identities of its checks.
var HookedChecks = struct {
	Set hookedSet
	Get hookedGet
}{
	Set: hookedSet{"Hooked.Set"},
	Get: hookedGet{"Hooked.Get"},
}

// hookedSet is the name of Hooked.Set, with the identities of its checks.
type hookedSet struct {
	conform.Callable[Hooked]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c hookedSet) Smoke() conform.ID[Hooked] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Set returns
// without a panic for derived inputs.
func (c hookedSet) SmokeGenerated() conform.ID[Hooked] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Set returns an error for a
// cancelled context.
func (c hookedSet) Cancel() conform.ID[Hooked] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Set returns an
// error for a cancelled context.
func (c hookedSet) CancelGenerated() conform.ID[Hooked] {
	return c.Check("cancel#generated")
}

// Partition returns the identity of the check that Get does not return in
// one partition what Set wrote in another.
func (c hookedSet) Partition() conform.ID[Hooked] {
	return c.Check("partition")
}

// PartitionGenerated returns the identity of the check that Get does not
// return in one partition what Set wrote in another.
func (c hookedSet) PartitionGenerated() conform.ID[Hooked] {
	return c.Check("partition#generated")
}

// hookedGet is the name of Hooked.Get, with the identities of its checks.
type hookedGet struct {
	conform.Callable[Hooked]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c hookedGet) Smoke() conform.ID[Hooked] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c hookedGet) SmokeGenerated() conform.ID[Hooked] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c hookedGet) Cancel() conform.ID[Hooked] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c hookedGet) CancelGenerated() conform.ID[Hooked] {
	return c.Check("cancel#generated")
}

// HookedSuite returns the checks of Hooked, and the calls and the budgets
// of its callables.
func HookedSuite() conform.Suite[Hooked] {
	return conform.Suite[Hooked]{
		Declaration: "Hooked",
		Package:     "example.com/shop/drawn",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Hooked]{
			{
				ID:          HookedChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HookedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Set(ctx, string("scope-a"), int(42), func() int { return 0 })
					})
				},
				Defect: hookedSetPanics,
			},
			{
				ID:          HookedChecks.Set.SmokeGenerated(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *HookedHandle, c *prop.Case) {
					drawnScope := c.Draw(prop.Of[string](), "scope")
					drawnValue := c.Draw(prop.Of[int](), "value")
					fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
					drawnFn := func() func() int {
						var fnAt atomic.Int64
						return func() (r0 int) {
							i := int(fnAt.Add(1) - 1)
							if i >= len(fnResults) {
								return r0
							}
							return fnResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnScope, drawnValue, drawnFn())
					})
				},
				Defect: hookedSetPanics,
			},
			{
				ID:          HookedChecks.Set.Cancel(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *HookedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Set(ctx, string("scope-a"), int(42), func() int { return 0 })
					})
				},
				Defect: hookedSetIgnoresCancel,
			},
			{
				ID:          HookedChecks.Set.CancelGenerated(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *HookedHandle, c *prop.Case) {
					drawnScope := c.Draw(prop.Of[string](), "scope")
					drawnValue := c.Draw(prop.Of[int](), "value")
					fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
					drawnFn := func() func() int {
						var fnAt atomic.Int64
						return func() (r0 int) {
							i := int(fnAt.Add(1) - 1)
							if i >= len(fnResults) {
								return r0
							}
							return fnResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnScope, drawnValue, drawnFn())
					})
				},
				Defect: hookedSetIgnoresCancel,
			},
			{
				ID:          HookedChecks.Set.Partition(),
				Statement:   "Get does not return in one partition what Set wrote in another",
				Claim:       spec.Partition,
				Probe:       spec.KeepsPartitionsApart,
				Reason:      spec.CrossedPartitions,
				Law:         spec.OnePartitionDoesNotServeAnother,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *HookedHandle) {
					probe.KeepsPartitionsApart(h, func() func(ctx context.Context) (int, error) {
						calls := 0
						return func(ctx context.Context) (int, error) {
							calls++
							w := int(42)
							s0 := string("scope-a")
							s2 := func() int { return 0 }
							if calls > 1 {
								w = int(7)
								s0 = string("scope-b")
								s2 = func() int { return 0 }
							}
							err := h.Subject().Set(ctx, s0, w, s2)
							return w, err
						}
					}(), func(ctx context.Context) (int, error) {
						return h.Subject().Get(ctx, string("scope-a"), func() int { return 0 })
					})
				},
				Defect: hookedSetSucceedsAlways,
			},
			{
				ID:          HookedChecks.Set.PartitionGenerated(),
				Statement:   "Get does not return in one partition what Set wrote in another",
				Claim:       spec.Partition,
				Probe:       spec.KeepsPartitionsApart,
				Reason:      spec.CrossedPartitions,
				Law:         spec.OnePartitionDoesNotServeAnother,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *HookedHandle, c *prop.Case) {
					drawnValue := c.Draw(
						prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(2), prop.MaxSize(2)),
						"value",
					)
					drawnScope := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"scope",
					)
					fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
					drawnFn := func() func() int {
						var fnAt atomic.Int64
						return func() (r0 int) {
							i := int(fnAt.Add(1) - 1)
							if i >= len(fnResults) {
								return r0
							}
							return fnResults[i]
						}
					}
					probe.KeepsPartitionsApart(h, func() func(ctx context.Context) (int, error) {
						calls := 0
						return func(ctx context.Context) (int, error) {
							calls++
							w := drawnValue[0]
							s0 := drawnScope[0]
							s2 := drawnFn()
							if calls > 1 {
								w = drawnValue[1]
								s0 = drawnScope[1]
								s2 = func() int { return 0 }
							}
							err := h.Subject().Set(ctx, s0, w, s2)
							return w, err
						}
					}(), func(ctx context.Context) (int, error) {
						return h.Subject().Get(ctx, drawnScope[0], drawnFn())
					})
				},
				Defect: hookedSetSucceedsAlways,
			},
			{
				ID:          HookedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HookedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("scope-a"), func() int { return 0 })
						return err
					})
				},
				Defect: hookedGetPanics,
			},
			{
				ID:          HookedChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *HookedHandle, c *prop.Case) {
					drawnScope := c.Draw(prop.Of[string](), "scope")
					fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
					drawnFn := func() func() int {
						var fnAt atomic.Int64
						return func() (r0 int) {
							i := int(fnAt.Add(1) - 1)
							if i >= len(fnResults) {
								return r0
							}
							return fnResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnScope, drawnFn())
						return err
					})
				},
				Defect: hookedGetPanics,
			},
			{
				ID:          HookedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *HookedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("scope-a"), func() int { return 0 })
						return err
					})
				},
				Defect: hookedGetIgnoresCancel,
			},
			{
				ID:          HookedChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *HookedHandle, c *prop.Case) {
					drawnScope := c.Draw(prop.Of[string](), "scope")
					fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
					drawnFn := func() func() int {
						var fnAt atomic.Int64
						return func() (r0 int) {
							i := int(fnAt.Add(1) - 1)
							if i >= len(fnResults) {
								return r0
							}
							return fnResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnScope, drawnFn())
						return err
					})
				},
				Defect: hookedGetIgnoresCancel,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:multi-arg-writer",
				Callable:       "Hooked.Set",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Hooked]{
			{
				Calls:    hookedSetCalls,
				Over:     newHookedOver,
				Callable: HookedChecks.Set.Callable,
			},
			{
				Calls:    hookedGetCalls,
				Over:     newHookedOver,
				Callable: HookedChecks.Get.Callable,
			},
		},
	}
}

// RunHookedSuite applies the checks of Hooked, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunHookedSuite[T Hooked](t *testing.T, newSubject func() T, opts ...HookedOption) {
	t.Helper()
	RunHookedSuiteFromHandle(t, func(*HookedHandle) T { return newSubject() }, opts...)
}

// RunHookedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunHookedSuiteFromHandle[T Hooked](
	t *testing.T,
	newSubject func(h *HookedHandle) T,
	opts ...HookedOption,
) {
	t.Helper()
	conform.Run(t, HookedSuite(), hookedSubjects(newSubject), opts...)
}

// FuzzHookedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzHookedSuite[T Hooked](
	f *testing.F,
	callable conform.Callable[Hooked],
	newSubject func(h *HookedHandle) T,
	opts ...HookedOption,
) {
	f.Helper()
	conform.Fuzz(f, HookedSuite(), callable, hookedSubjects(newSubject), opts...)
}

// RunHookedBudgets applies the allocation ceilings of Hooked and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunHookedBudgets[T Hooked](
	t *testing.T,
	newSubject func(h *HookedHandle) T,
	opts ...HookedOption,
) {
	t.Helper()
	conform.RunBudgets(t, HookedSuite(), hookedSubjects(newSubject), opts...)
}

// BenchHookedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchHookedSuite[T Hooked](
	b conform.B,
	callable conform.Callable[Hooked],
	newSubject func(h *HookedHandle) T,
	opts ...HookedOption,
) {
	b.Helper()
	conform.Bench(b, HookedSuite(), callable, hookedSubjects(newSubject), opts...)
}

// hookedSubjects returns the factory of the subjects of Hooked that calls
// newSubject with the handle of each check.
func hookedSubjects[T Hooked](newSubject func(h *HookedHandle) T) func(h *HookedHandle) Hooked {
	return func(h *HookedHandle) Hooked { return newSubject(h) }
}

// hookedSetCalls generates the calls of Hooked.Set with the inputs that
// each case draws.
var hookedSetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hooked) error {
	scope := c.Draw(prop.Of[string](), "scope")
	value := c.Draw(prop.Of[int](), "value")
	fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
	var fnAt atomic.Int64
	fn := func() (r0 int) {
		i := int(fnAt.Add(1) - 1)
		if i >= len(fnResults) {
			return r0
		}
		return fnResults[i]
	}
	return func(ctx context.Context, subject Hooked) error {
		return subject.Set(ctx, scope, value, fn)
	}
})

// hookedGetCalls generates the calls of Hooked.Get with the inputs that
// each case draws.
var hookedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hooked) error {
	scope := c.Draw(prop.Of[string](), "scope")
	fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
	var fnAt atomic.Int64
	fn := func() (r0 int) {
		i := int(fnAt.Add(1) - 1)
		if i >= len(fnResults) {
			return r0
		}
		return fnResults[i]
	}
	return func(ctx context.Context, subject Hooked) error {
		_, err := subject.Get(ctx, scope, fn)
		return err
	}
})

// hookedOverEmbedded is Hooked under the name of the field that hookedOver
// embeds.
type hookedOverEmbedded = Hooked

// hookedOver is a subject of Hooked whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type hookedOver struct {
	hookedOverEmbedded
	spend func()
}

// Set calls spend and returns its zero results.
func (o hookedOver) Set(context.Context, string, int, func() int) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o hookedOver) Get(context.Context, string, func() int) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newHookedOver returns the planted defect of a ceiling of Hooked, whose
// callables call spend.
func newHookedOver(_ *HookedHandle, spend func()) Hooked {
	return hookedOver{spend: spend}
}

// hookedSetPanicsDefect is the planted defect of a Hooked whose Set
// panics.
type hookedSetPanicsDefect struct {
	hookedOverEmbedded
}

// Set is the override of the planted defect.
func (*hookedSetPanicsDefect) Set(p0 context.Context, p1 string, p2 int, p3 func() int) (r0 error) {
	panic("the planted defect of a Hooked whose Set panics")
}

// hookedSetPanics returns the planted defect of a Hooked whose Set panics
// over a fresh subject.
func hookedSetPanics(h *HookedHandle) Hooked {
	return &hookedSetPanicsDefect{hookedOverEmbedded: h.Fresh()}
}

// hookedSetIgnoresCancelDefect is the planted defect of a Hooked whose Set
// ignores the context.
type hookedSetIgnoresCancelDefect struct {
	hookedOverEmbedded
}

// Set is the override of the planted defect.
func (d *hookedSetIgnoresCancelDefect) Set(
	p0 context.Context,
	p1 string,
	p2 int,
	p3 func() int,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.hookedOverEmbedded.Set(p0, p1, p2, p3)
}

// hookedSetIgnoresCancel returns the planted defect of a Hooked whose Set
// ignores the context over a fresh subject.
func hookedSetIgnoresCancel(h *HookedHandle) Hooked {
	return &hookedSetIgnoresCancelDefect{hookedOverEmbedded: h.Fresh()}
}

// hookedSetSucceedsAlwaysDefect is the planted defect of a Hooked whose
// Set keeps no partition.
type hookedSetSucceedsAlwaysDefect struct {
	hookedOverEmbedded
}

// Set is the override of the planted defect.
func (*hookedSetSucceedsAlwaysDefect) Set(
	p0 context.Context,
	p1 string,
	p2 int,
	p3 func() int,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// hookedSetSucceedsAlways returns the planted defect of a Hooked whose Set
// keeps no partition over a fresh subject.
func hookedSetSucceedsAlways(h *HookedHandle) Hooked {
	return &hookedSetSucceedsAlwaysDefect{hookedOverEmbedded: h.Fresh()}
}

// hookedGetPanicsDefect is the planted defect of a Hooked whose Get
// panics.
type hookedGetPanicsDefect struct {
	hookedOverEmbedded
}

// Get is the override of the planted defect.
func (*hookedGetPanicsDefect) Get(p0 context.Context, p1 string, p2 func() int) (r0 int, r1 error) {
	panic("the planted defect of a Hooked whose Get panics")
}

// hookedGetPanics returns the planted defect of a Hooked whose Get panics
// over a fresh subject.
func hookedGetPanics(h *HookedHandle) Hooked {
	return &hookedGetPanicsDefect{hookedOverEmbedded: h.Fresh()}
}

// hookedGetIgnoresCancelDefect is the planted defect of a Hooked whose Get
// ignores the context.
type hookedGetIgnoresCancelDefect struct {
	hookedOverEmbedded
}

// Get is the override of the planted defect.
func (d *hookedGetIgnoresCancelDefect) Get(
	p0 context.Context,
	p1 string,
	p2 func() int,
) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.hookedOverEmbedded.Get(p0, p1, p2)
}

// hookedGetIgnoresCancel returns the planted defect of a Hooked whose Get
// ignores the context over a fresh subject.
func hookedGetIgnoresCancel(h *HookedHandle) Hooked {
	return &hookedGetIgnoresCancelDefect{hookedOverEmbedded: h.Fresh()}
}

// The derived checks of Lookups:
//
// Lookups.Has/smoke: Has returns without a panic for derived inputs.
//
// Lookups.Has/smoke#generated: Has returns without a panic for derived
// inputs.
//
// Lookups.Has/cancel: Has returns an error for a cancelled context.
//
// Lookups.Has/cancel#generated: Has returns an error for a cancelled
// context.
//
// Lookups.Has/zero-on-error: Has returns the zero value beside an error.
//
// Lookups.Has/zero-on-error#generated: Has returns the zero value beside
// an error.
//
// Lookups.IsPresent/smoke: IsPresent returns without a panic for derived
// inputs.
//
// Lookups.IsPresent/smoke#generated: IsPresent returns without a panic for
// derived inputs.
//
// Lookups.IsPresent/cancel: IsPresent returns an error for a cancelled
// context.
//
// Lookups.IsPresent/cancel#generated: IsPresent returns an error for a
// cancelled context.
//
// Lookups.IsPresent/zero-on-error: IsPresent returns the zero value beside
// an error.
//
// Lookups.IsPresent/zero-on-error#generated: IsPresent returns the zero
// value beside an error.
//
// Lookups.IsEnabled/smoke: IsEnabled returns without a panic for derived
// inputs.
//
// Lookups.IsEnabled/smoke#generated: IsEnabled returns without a panic for
// derived inputs.
//
// Lookups.IsEnabled/cancel: IsEnabled returns an error for a cancelled
// context.
//
// Lookups.IsEnabled/cancel#generated: IsEnabled returns an error for a
// cancelled context.
//
// Lookups.IsEnabled/zero-on-error: IsEnabled returns the zero value beside
// an error.
//
// Lookups.IsEnabled/zero-on-error#generated: IsEnabled returns the zero
// value beside an error.
//
// Lookups.State/smoke: State returns without a panic for derived inputs.
//
// Lookups.State/smoke#generated: State returns without a panic for derived
// inputs.
//
// Lookups.State/cancel: State returns an error for a cancelled context.
//
// Lookups.State/cancel#generated: State returns an error for a cancelled
// context.
//
// Lookups.State/zero-on-error: State returns the zero value beside an
// error.
//
// Lookups.State/zero-on-error#generated: State returns the zero value
// beside an error.
//
// Lookups.Toggles/smoke: Toggles returns without a panic for derived
// inputs.
//
// Lookups.Toggles/smoke#generated: Toggles returns without a panic for
// derived inputs.
//
// Lookups.Toggles/cancel: Toggles returns an error for a cancelled
// context.
//
// Lookups.Toggles/cancel#generated: Toggles returns an error for a
// cancelled context.
//
// Lookups.Toggles/zero-on-error: Toggles returns the zero value beside an
// error.
//
// Lookups.Toggles/zero-on-error#generated: Toggles returns the zero value
// beside an error.
//
// Lookups.Amount/smoke: Amount returns without a panic for derived inputs.
//
// Lookups.Amount/smoke#generated: Amount returns without a panic for
// derived inputs.
//
// Lookups.Amount/cancel: Amount returns an error for a cancelled context.
//
// Lookups.Amount/cancel#generated: Amount returns an error for a cancelled
// context.
//
// Lookups.Amount/zero-on-error: Amount returns the zero value beside an
// error.
//
// Lookups.Amount/zero-on-error#generated: Amount returns the zero value
// beside an error.
//
// Lookups.Ratio/smoke: Ratio returns without a panic for derived inputs.
//
// Lookups.Ratio/smoke#generated: Ratio returns without a panic for derived
// inputs.
//
// Lookups.Ratio/cancel: Ratio returns an error for a cancelled context.
//
// Lookups.Ratio/cancel#generated: Ratio returns an error for a cancelled
// context.
//
// Lookups.Ratio/zero-on-error: Ratio returns the zero value beside an
// error.
//
// Lookups.Ratio/zero-on-error#generated: Ratio returns the zero value
// beside an error.
//
// Lookups.Code/smoke: Code returns without a panic for derived inputs.
//
// Lookups.Code/smoke#generated: Code returns without a panic for derived
// inputs.
//
// Lookups.Code/cancel: Code returns an error for a cancelled context.
//
// Lookups.Code/cancel#generated: Code returns an error for a cancelled
// context.
//
// Lookups.Code/zero-on-error: Code returns the zero value beside an error.
//
// Lookups.Code/zero-on-error#generated: Code returns the zero value beside
// an error.

// The derived checks of Lookups that witness their laws at the shallow
// end:
//
// Lookups.Has/smoke
//
// Lookups.IsPresent/smoke
//
// Lookups.IsEnabled/smoke
//
// Lookups.State/smoke
//
// Lookups.Toggles/smoke
//
// Lookups.Amount/smoke
//
// Lookups.Ratio/smoke
//
// Lookups.Code/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Lookups writes.
type (
	LookupsHandle = conform.Handle[Lookups]
	LookupsCheck  = conform.Check[Lookups]
	LookupsOption = conform.Option[Lookups]
)

// LookupsChecks contains the name of each callable of Lookups, with the
// identities of its checks.
var LookupsChecks = struct {
	Has       lookupsHas
	IsPresent lookupsIsPresent
	IsEnabled lookupsIsEnabled
	State     lookupsState
	Toggles   lookupsToggles
	Amount    lookupsAmount
	Ratio     lookupsRatio
	Code      lookupsCode
}{
	Has:       lookupsHas{"Lookups.Has"},
	IsPresent: lookupsIsPresent{"Lookups.IsPresent"},
	IsEnabled: lookupsIsEnabled{"Lookups.IsEnabled"},
	State:     lookupsState{"Lookups.State"},
	Toggles:   lookupsToggles{"Lookups.Toggles"},
	Amount:    lookupsAmount{"Lookups.Amount"},
	Ratio:     lookupsRatio{"Lookups.Ratio"},
	Code:      lookupsCode{"Lookups.Code"},
}

// lookupsHas is the name of Lookups.Has, with the identities of its
// checks.
type lookupsHas struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that Has returns without a panic
// for derived inputs.
func (c lookupsHas) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Has returns
// without a panic for derived inputs.
func (c lookupsHas) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Has returns an error for a
// cancelled context.
func (c lookupsHas) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Has returns an
// error for a cancelled context.
func (c lookupsHas) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Has returns the zero
// value beside an error.
func (c lookupsHas) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Has returns
// the zero value beside an error.
func (c lookupsHas) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// lookupsIsPresent is the name of Lookups.IsPresent, with the identities
// of its checks.
type lookupsIsPresent struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that IsPresent returns without a
// panic for derived inputs.
func (c lookupsIsPresent) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that IsPresent returns
// without a panic for derived inputs.
func (c lookupsIsPresent) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that IsPresent returns an error
// for a cancelled context.
func (c lookupsIsPresent) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that IsPresent returns
// an error for a cancelled context.
func (c lookupsIsPresent) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that IsPresent returns the
// zero value beside an error.
func (c lookupsIsPresent) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that IsPresent
// returns the zero value beside an error.
func (c lookupsIsPresent) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// lookupsIsEnabled is the name of Lookups.IsEnabled, with the identities
// of its checks.
type lookupsIsEnabled struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that IsEnabled returns without a
// panic for derived inputs.
func (c lookupsIsEnabled) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that IsEnabled returns
// without a panic for derived inputs.
func (c lookupsIsEnabled) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that IsEnabled returns an error
// for a cancelled context.
func (c lookupsIsEnabled) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that IsEnabled returns
// an error for a cancelled context.
func (c lookupsIsEnabled) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that IsEnabled returns the
// zero value beside an error.
func (c lookupsIsEnabled) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that IsEnabled
// returns the zero value beside an error.
func (c lookupsIsEnabled) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// lookupsState is the name of Lookups.State, with the identities of its
// checks.
type lookupsState struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that State returns without a
// panic for derived inputs.
func (c lookupsState) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that State returns
// without a panic for derived inputs.
func (c lookupsState) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that State returns an error for
// a cancelled context.
func (c lookupsState) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that State returns an
// error for a cancelled context.
func (c lookupsState) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that State returns the
// zero value beside an error.
func (c lookupsState) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that State
// returns the zero value beside an error.
func (c lookupsState) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// lookupsToggles is the name of Lookups.Toggles, with the identities of
// its checks.
type lookupsToggles struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that Toggles returns without a
// panic for derived inputs.
func (c lookupsToggles) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Toggles returns
// without a panic for derived inputs.
func (c lookupsToggles) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Toggles returns an error
// for a cancelled context.
func (c lookupsToggles) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Toggles returns
// an error for a cancelled context.
func (c lookupsToggles) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Toggles returns the
// zero value beside an error.
func (c lookupsToggles) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Toggles
// returns the zero value beside an error.
func (c lookupsToggles) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// lookupsAmount is the name of Lookups.Amount, with the identities of its
// checks.
type lookupsAmount struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that Amount returns without a
// panic for derived inputs.
func (c lookupsAmount) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Amount returns
// without a panic for derived inputs.
func (c lookupsAmount) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Amount returns an error
// for a cancelled context.
func (c lookupsAmount) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Amount returns an
// error for a cancelled context.
func (c lookupsAmount) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Amount returns the
// zero value beside an error.
func (c lookupsAmount) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Amount
// returns the zero value beside an error.
func (c lookupsAmount) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// lookupsRatio is the name of Lookups.Ratio, with the identities of its
// checks.
type lookupsRatio struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that Ratio returns without a
// panic for derived inputs.
func (c lookupsRatio) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Ratio returns
// without a panic for derived inputs.
func (c lookupsRatio) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Ratio returns an error for
// a cancelled context.
func (c lookupsRatio) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Ratio returns an
// error for a cancelled context.
func (c lookupsRatio) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Ratio returns the
// zero value beside an error.
func (c lookupsRatio) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Ratio
// returns the zero value beside an error.
func (c lookupsRatio) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// lookupsCode is the name of Lookups.Code, with the identities of its
// checks.
type lookupsCode struct {
	conform.Callable[Lookups]
}

// Smoke returns the identity of the check that Code returns without a
// panic for derived inputs.
func (c lookupsCode) Smoke() conform.ID[Lookups] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Code returns
// without a panic for derived inputs.
func (c lookupsCode) SmokeGenerated() conform.ID[Lookups] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Code returns an error for
// a cancelled context.
func (c lookupsCode) Cancel() conform.ID[Lookups] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Code returns an
// error for a cancelled context.
func (c lookupsCode) CancelGenerated() conform.ID[Lookups] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Code returns the zero
// value beside an error.
func (c lookupsCode) ZeroOnError() conform.ID[Lookups] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Code returns
// the zero value beside an error.
func (c lookupsCode) ZeroOnErrorGenerated() conform.ID[Lookups] {
	return c.Check("zero-on-error#generated")
}

// LookupsSuite returns the checks of Lookups, and the calls and the
// budgets of its callables.
func LookupsSuite() conform.Suite[Lookups] {
	return conform.Suite[Lookups]{
		Declaration: "Lookups",
		Package:     "example.com/shop/drawn",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Lookups]{
			{
				ID:          LookupsChecks.Has.Smoke(),
				Statement:   "Has returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsHasPanics,
			},
			{
				ID:          LookupsChecks.Has.SmokeGenerated(),
				Statement:   "Has returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsHasPanics,
			},
			{
				ID:          LookupsChecks.Has.Cancel(),
				Statement:   "Has returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsHasIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Has.CancelGenerated(),
				Statement:   "Has returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsHasIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Has.ZeroOnError(),
				Statement:   "Has returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Found, error) {
						return h.Subject().Has(ctx, string("key-a"))
					}, Found{})
				},
				Defect: lookupsHasValueBeside,
			},
			{
				ID:          LookupsChecks.Has.ZeroOnErrorGenerated(),
				Statement:   "Has returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Found, error) {
						return h.Subject().Has(ctx, drawnKey)
					}, Found{})
				},
				Defect: lookupsHasValueBeside,
			},
			{
				ID:          LookupsChecks.IsPresent.Smoke(),
				Statement:   "IsPresent returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().IsPresent(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsIsPresentPanics,
			},
			{
				ID:          LookupsChecks.IsPresent.SmokeGenerated(),
				Statement:   "IsPresent returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().IsPresent(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsIsPresentPanics,
			},
			{
				ID:          LookupsChecks.IsPresent.Cancel(),
				Statement:   "IsPresent returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().IsPresent(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsIsPresentIgnoresCancel,
			},
			{
				ID:          LookupsChecks.IsPresent.CancelGenerated(),
				Statement:   "IsPresent returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().IsPresent(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsIsPresentIgnoresCancel,
			},
			{
				ID:          LookupsChecks.IsPresent.ZeroOnError(),
				Statement:   "IsPresent returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Present, error) {
						return h.Subject().IsPresent(ctx, string("key-a"))
					}, Present(false))
				},
				Defect: lookupsIsPresentValueBeside,
			},
			{
				ID:          LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
				Statement:   "IsPresent returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Present, error) {
						return h.Subject().IsPresent(ctx, drawnKey)
					}, Present(false))
				},
				Defect: lookupsIsPresentValueBeside,
			},
			{
				ID:          LookupsChecks.IsEnabled.Smoke(),
				Statement:   "IsEnabled returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().IsEnabled(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsIsEnabledPanics,
			},
			{
				ID:          LookupsChecks.IsEnabled.SmokeGenerated(),
				Statement:   "IsEnabled returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().IsEnabled(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsIsEnabledPanics,
			},
			{
				ID:          LookupsChecks.IsEnabled.Cancel(),
				Statement:   "IsEnabled returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().IsEnabled(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsIsEnabledIgnoresCancel,
			},
			{
				ID:          LookupsChecks.IsEnabled.CancelGenerated(),
				Statement:   "IsEnabled returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().IsEnabled(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsIsEnabledIgnoresCancel,
			},
			{
				ID:          LookupsChecks.IsEnabled.ZeroOnError(),
				Statement:   "IsEnabled returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Enabled, error) {
						return h.Subject().IsEnabled(ctx, string("key-a"))
					}, Enabled(false))
				},
				Defect: lookupsIsEnabledValueBeside,
			},
			{
				ID:          LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
				Statement:   "IsEnabled returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Enabled, error) {
						return h.Subject().IsEnabled(ctx, drawnKey)
					}, Enabled(false))
				},
				Defect: lookupsIsEnabledValueBeside,
			},
			{
				ID:          LookupsChecks.State.Smoke(),
				Statement:   "State returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().State(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsStatePanics,
			},
			{
				ID:          LookupsChecks.State.SmokeGenerated(),
				Statement:   "State returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().State(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsStatePanics,
			},
			{
				ID:          LookupsChecks.State.Cancel(),
				Statement:   "State returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().State(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsStateIgnoresCancel,
			},
			{
				ID:          LookupsChecks.State.CancelGenerated(),
				Statement:   "State returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().State(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsStateIgnoresCancel,
			},
			{
				ID:          LookupsChecks.State.ZeroOnError(),
				Statement:   "State returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Switch, error) {
						return h.Subject().State(ctx, string("key-a"))
					}, Switch{})
				},
				Defect: lookupsStateValueBeside,
			},
			{
				ID:          LookupsChecks.State.ZeroOnErrorGenerated(),
				Statement:   "State returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Switch, error) {
						return h.Subject().State(ctx, drawnKey)
					}, Switch{})
				},
				Defect: lookupsStateValueBeside,
			},
			{
				ID:          LookupsChecks.Toggles.Smoke(),
				Statement:   "Toggles returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Toggles(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsTogglesPanics,
			},
			{
				ID:          LookupsChecks.Toggles.SmokeGenerated(),
				Statement:   "Toggles returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Toggles(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsTogglesPanics,
			},
			{
				ID:          LookupsChecks.Toggles.Cancel(),
				Statement:   "Toggles returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Toggles(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsTogglesIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Toggles.CancelGenerated(),
				Statement:   "Toggles returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Toggles(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsTogglesIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Toggles.ZeroOnError(),
				Statement:   "Toggles returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]bool, error) {
						return h.Subject().Toggles(ctx, string("key-a"))
					}, []bool(nil))
				},
				Defect: lookupsTogglesValueBeside,
			},
			{
				ID:          LookupsChecks.Toggles.ZeroOnErrorGenerated(),
				Statement:   "Toggles returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) ([]bool, error) {
						return h.Subject().Toggles(ctx, drawnKey)
					}, []bool(nil))
				},
				Defect: lookupsTogglesValueBeside,
			},
			{
				ID:          LookupsChecks.Amount.Smoke(),
				Statement:   "Amount returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Amount(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsAmountPanics,
			},
			{
				ID:          LookupsChecks.Amount.SmokeGenerated(),
				Statement:   "Amount returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Amount(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsAmountPanics,
			},
			{
				ID:          LookupsChecks.Amount.Cancel(),
				Statement:   "Amount returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Amount(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsAmountIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Amount.CancelGenerated(),
				Statement:   "Amount returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Amount(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsAmountIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Amount.ZeroOnError(),
				Statement:   "Amount returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Amount, error) {
						return h.Subject().Amount(ctx, string("key-a"))
					}, Amount(0))
				},
				Defect: lookupsAmountValueBeside,
			},
			{
				ID:          LookupsChecks.Amount.ZeroOnErrorGenerated(),
				Statement:   "Amount returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Amount, error) {
						return h.Subject().Amount(ctx, drawnKey)
					}, Amount(0))
				},
				Defect: lookupsAmountValueBeside,
			},
			{
				ID:          LookupsChecks.Ratio.Smoke(),
				Statement:   "Ratio returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Ratio(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsRatioPanics,
			},
			{
				ID:          LookupsChecks.Ratio.SmokeGenerated(),
				Statement:   "Ratio returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Ratio(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsRatioPanics,
			},
			{
				ID:          LookupsChecks.Ratio.Cancel(),
				Statement:   "Ratio returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Ratio(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsRatioIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Ratio.CancelGenerated(),
				Statement:   "Ratio returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Ratio(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsRatioIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Ratio.ZeroOnError(),
				Statement:   "Ratio returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Ratio, error) {
						return h.Subject().Ratio(ctx, string("key-a"))
					}, Ratio(0))
				},
				Defect: lookupsRatioValueBeside,
			},
			{
				ID:          LookupsChecks.Ratio.ZeroOnErrorGenerated(),
				Statement:   "Ratio returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Ratio, error) {
						return h.Subject().Ratio(ctx, drawnKey)
					}, Ratio(0))
				},
				Defect: lookupsRatioValueBeside,
			},
			{
				ID:          LookupsChecks.Code.Smoke(),
				Statement:   "Code returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LookupsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Code(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsCodePanics,
			},
			{
				ID:          LookupsChecks.Code.SmokeGenerated(),
				Statement:   "Code returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Code(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsCodePanics,
			},
			{
				ID:          LookupsChecks.Code.Cancel(),
				Statement:   "Code returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Code(ctx, string("key-a"))
						return err
					})
				},
				Defect: lookupsCodeIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Code.CancelGenerated(),
				Statement:   "Code returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Code(ctx, drawnKey)
						return err
					})
				},
				Defect: lookupsCodeIgnoresCancel,
			},
			{
				ID:          LookupsChecks.Code.ZeroOnError(),
				Statement:   "Code returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LookupsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Code, error) {
						return h.Subject().Code(ctx, string("key-a"))
					}, Code(""))
				},
				Defect: lookupsCodeValueBeside,
			},
			{
				ID:          LookupsChecks.Code.ZeroOnErrorGenerated(),
				Statement:   "Code returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LookupsHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Code, error) {
						return h.Subject().Code(ctx, drawnKey)
					}, Code(""))
				},
				Defect: lookupsCodeValueBeside,
			},
		},
		Calls: []conform.Call[Lookups]{
			{
				Calls:    lookupsHasCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.Has.Callable,
			},
			{
				Calls:    lookupsIsPresentCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.IsPresent.Callable,
			},
			{
				Calls:    lookupsIsEnabledCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.IsEnabled.Callable,
			},
			{
				Calls:    lookupsStateCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.State.Callable,
			},
			{
				Calls:    lookupsTogglesCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.Toggles.Callable,
			},
			{
				Calls:    lookupsAmountCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.Amount.Callable,
			},
			{
				Calls:    lookupsRatioCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.Ratio.Callable,
			},
			{
				Calls:    lookupsCodeCalls,
				Over:     newLookupsOver,
				Callable: LookupsChecks.Code.Callable,
			},
		},
	}
}

// RunLookupsSuite applies the checks of Lookups, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLookupsSuite[T Lookups](t *testing.T, newSubject func() T, opts ...LookupsOption) {
	t.Helper()
	RunLookupsSuiteFromHandle(t, func(*LookupsHandle) T { return newSubject() }, opts...)
}

// RunLookupsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLookupsSuiteFromHandle[T Lookups](
	t *testing.T,
	newSubject func(h *LookupsHandle) T,
	opts ...LookupsOption,
) {
	t.Helper()
	conform.Run(t, LookupsSuite(), lookupsSubjects(newSubject), opts...)
}

// FuzzLookupsSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLookupsSuite[T Lookups](
	f *testing.F,
	callable conform.Callable[Lookups],
	newSubject func(h *LookupsHandle) T,
	opts ...LookupsOption,
) {
	f.Helper()
	conform.Fuzz(f, LookupsSuite(), callable, lookupsSubjects(newSubject), opts...)
}

// RunLookupsBudgets applies the allocation ceilings of Lookups and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLookupsBudgets[T Lookups](
	t *testing.T,
	newSubject func(h *LookupsHandle) T,
	opts ...LookupsOption,
) {
	t.Helper()
	conform.RunBudgets(t, LookupsSuite(), lookupsSubjects(newSubject), opts...)
}

// BenchLookupsSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchLookupsSuite[T Lookups](
	b conform.B,
	callable conform.Callable[Lookups],
	newSubject func(h *LookupsHandle) T,
	opts ...LookupsOption,
) {
	b.Helper()
	conform.Bench(b, LookupsSuite(), callable, lookupsSubjects(newSubject), opts...)
}

// lookupsSubjects returns the factory of the subjects of Lookups that
// calls newSubject with the handle of each check.
func lookupsSubjects[T Lookups](
	newSubject func(h *LookupsHandle) T,
) func(h *LookupsHandle) Lookups {
	return func(h *LookupsHandle) Lookups { return newSubject(h) }
}

// lookupsHasCalls generates the calls of Lookups.Has with the inputs that
// each case draws.
var lookupsHasCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.Has(ctx, key)
		return err
	}
})

// lookupsIsPresentCalls generates the calls of Lookups.IsPresent with the
// inputs that each case draws.
var lookupsIsPresentCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.IsPresent(ctx, key)
		return err
	}
})

// lookupsIsEnabledCalls generates the calls of Lookups.IsEnabled with the
// inputs that each case draws.
var lookupsIsEnabledCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.IsEnabled(ctx, key)
		return err
	}
})

// lookupsStateCalls generates the calls of Lookups.State with the inputs
// that each case draws.
var lookupsStateCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.State(ctx, key)
		return err
	}
})

// lookupsTogglesCalls generates the calls of Lookups.Toggles with the
// inputs that each case draws.
var lookupsTogglesCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.Toggles(ctx, key)
		return err
	}
})

// lookupsAmountCalls generates the calls of Lookups.Amount with the inputs
// that each case draws.
var lookupsAmountCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.Amount(ctx, key)
		return err
	}
})

// lookupsRatioCalls generates the calls of Lookups.Ratio with the inputs
// that each case draws.
var lookupsRatioCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.Ratio(ctx, key)
		return err
	}
})

// lookupsCodeCalls generates the calls of Lookups.Code with the inputs
// that each case draws.
var lookupsCodeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lookups) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lookups) error {
		_, err := subject.Code(ctx, key)
		return err
	}
})

// lookupsOverEmbedded is Lookups under the name of the field that
// lookupsOver embeds.
type lookupsOverEmbedded = Lookups

// lookupsOver is a subject of Lookups whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type lookupsOver struct {
	lookupsOverEmbedded
	spend func()
}

// Has calls spend and returns its zero results.
func (o lookupsOver) Has(context.Context, string) (r0 Found, r1 error) {
	o.spend()
	return r0, r1
}

// IsPresent calls spend and returns its zero results.
func (o lookupsOver) IsPresent(context.Context, string) (r0 Present, r1 error) {
	o.spend()
	return r0, r1
}

// IsEnabled calls spend and returns its zero results.
func (o lookupsOver) IsEnabled(context.Context, string) (r0 Enabled, r1 error) {
	o.spend()
	return r0, r1
}

// State calls spend and returns its zero results.
func (o lookupsOver) State(context.Context, string) (r0 Switch, r1 error) {
	o.spend()
	return r0, r1
}

// Toggles calls spend and returns its zero results.
func (o lookupsOver) Toggles(context.Context, string) (r0 []bool, r1 error) {
	o.spend()
	return r0, r1
}

// Amount calls spend and returns its zero results.
func (o lookupsOver) Amount(context.Context, string) (r0 Amount, r1 error) {
	o.spend()
	return r0, r1
}

// Ratio calls spend and returns its zero results.
func (o lookupsOver) Ratio(context.Context, string) (r0 Ratio, r1 error) {
	o.spend()
	return r0, r1
}

// Code calls spend and returns its zero results.
func (o lookupsOver) Code(context.Context, string) (r0 Code, r1 error) {
	o.spend()
	return r0, r1
}

// newLookupsOver returns the planted defect of a ceiling of Lookups, whose
// callables call spend.
func newLookupsOver(_ *LookupsHandle, spend func()) Lookups {
	return lookupsOver{spend: spend}
}

// lookupsHasPanicsDefect is the planted defect of a Lookups whose Has
// panics.
type lookupsHasPanicsDefect struct {
	lookupsOverEmbedded
}

// Has is the override of the planted defect.
func (*lookupsHasPanicsDefect) Has(p0 context.Context, p1 string) (r0 Found, r1 error) {
	panic("the planted defect of a Lookups whose Has panics")
}

// lookupsHasPanics returns the planted defect of a Lookups whose Has
// panics over a fresh subject.
func lookupsHasPanics(h *LookupsHandle) Lookups {
	return &lookupsHasPanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsHasIgnoresCancelDefect is the planted defect of a Lookups whose
// Has ignores the context.
type lookupsHasIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// Has is the override of the planted defect.
func (d *lookupsHasIgnoresCancelDefect) Has(p0 context.Context, p1 string) (r0 Found, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.Has(p0, p1)
}

// lookupsHasIgnoresCancel returns the planted defect of a Lookups whose
// Has ignores the context over a fresh subject.
func lookupsHasIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsHasIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsHasValueBesideDefect is the planted defect of a Lookups whose Has
// returns a value beside an error.
type lookupsHasValueBesideDefect struct {
	lookupsOverEmbedded
}

// Has is the override of the planted defect.
func (d *lookupsHasValueBesideDefect) Has(p0 context.Context, p1 string) (r0 Found, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.Has(p0, p1)
	if r1 != nil {
		r0 = Found{Found: true}
	}
	return r0, r1
}

// lookupsHasValueBeside returns the planted defect of a Lookups whose Has
// returns a value beside an error over a fresh subject.
func lookupsHasValueBeside(h *LookupsHandle) Lookups {
	return &lookupsHasValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsIsPresentPanicsDefect is the planted defect of a Lookups whose
// IsPresent panics.
type lookupsIsPresentPanicsDefect struct {
	lookupsOverEmbedded
}

// IsPresent is the override of the planted defect.
func (*lookupsIsPresentPanicsDefect) IsPresent(
	p0 context.Context,
	p1 string,
) (r0 Present, r1 error) {
	panic("the planted defect of a Lookups whose IsPresent panics")
}

// lookupsIsPresentPanics returns the planted defect of a Lookups whose
// IsPresent panics over a fresh subject.
func lookupsIsPresentPanics(h *LookupsHandle) Lookups {
	return &lookupsIsPresentPanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsIsPresentIgnoresCancelDefect is the planted defect of a Lookups
// whose IsPresent ignores the context.
type lookupsIsPresentIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// IsPresent is the override of the planted defect.
func (d *lookupsIsPresentIgnoresCancelDefect) IsPresent(
	p0 context.Context,
	p1 string,
) (r0 Present, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.IsPresent(p0, p1)
}

// lookupsIsPresentIgnoresCancel returns the planted defect of a Lookups
// whose IsPresent ignores the context over a fresh subject.
func lookupsIsPresentIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsIsPresentIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsIsPresentValueBesideDefect is the planted defect of a Lookups
// whose IsPresent returns a value beside an error.
type lookupsIsPresentValueBesideDefect struct {
	lookupsOverEmbedded
}

// IsPresent is the override of the planted defect.
func (d *lookupsIsPresentValueBesideDefect) IsPresent(
	p0 context.Context,
	p1 string,
) (r0 Present, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.IsPresent(p0, p1)
	if r1 != nil {
		r0 = Present(true)
	}
	return r0, r1
}

// lookupsIsPresentValueBeside returns the planted defect of a Lookups
// whose IsPresent returns a value beside an error over a fresh subject.
func lookupsIsPresentValueBeside(h *LookupsHandle) Lookups {
	return &lookupsIsPresentValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsIsEnabledPanicsDefect is the planted defect of a Lookups whose
// IsEnabled panics.
type lookupsIsEnabledPanicsDefect struct {
	lookupsOverEmbedded
}

// IsEnabled is the override of the planted defect.
func (*lookupsIsEnabledPanicsDefect) IsEnabled(
	p0 context.Context,
	p1 string,
) (r0 Enabled, r1 error) {
	panic("the planted defect of a Lookups whose IsEnabled panics")
}

// lookupsIsEnabledPanics returns the planted defect of a Lookups whose
// IsEnabled panics over a fresh subject.
func lookupsIsEnabledPanics(h *LookupsHandle) Lookups {
	return &lookupsIsEnabledPanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsIsEnabledIgnoresCancelDefect is the planted defect of a Lookups
// whose IsEnabled ignores the context.
type lookupsIsEnabledIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// IsEnabled is the override of the planted defect.
func (d *lookupsIsEnabledIgnoresCancelDefect) IsEnabled(
	p0 context.Context,
	p1 string,
) (r0 Enabled, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.IsEnabled(p0, p1)
}

// lookupsIsEnabledIgnoresCancel returns the planted defect of a Lookups
// whose IsEnabled ignores the context over a fresh subject.
func lookupsIsEnabledIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsIsEnabledIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsIsEnabledValueBesideDefect is the planted defect of a Lookups
// whose IsEnabled returns a value beside an error.
type lookupsIsEnabledValueBesideDefect struct {
	lookupsOverEmbedded
}

// IsEnabled is the override of the planted defect.
func (d *lookupsIsEnabledValueBesideDefect) IsEnabled(
	p0 context.Context,
	p1 string,
) (r0 Enabled, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.IsEnabled(p0, p1)
	if r1 != nil {
		r0 = Enabled(true)
	}
	return r0, r1
}

// lookupsIsEnabledValueBeside returns the planted defect of a Lookups
// whose IsEnabled returns a value beside an error over a fresh subject.
func lookupsIsEnabledValueBeside(h *LookupsHandle) Lookups {
	return &lookupsIsEnabledValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsStatePanicsDefect is the planted defect of a Lookups whose State
// panics.
type lookupsStatePanicsDefect struct {
	lookupsOverEmbedded
}

// State is the override of the planted defect.
func (*lookupsStatePanicsDefect) State(p0 context.Context, p1 string) (r0 Switch, r1 error) {
	panic("the planted defect of a Lookups whose State panics")
}

// lookupsStatePanics returns the planted defect of a Lookups whose State
// panics over a fresh subject.
func lookupsStatePanics(h *LookupsHandle) Lookups {
	return &lookupsStatePanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsStateIgnoresCancelDefect is the planted defect of a Lookups whose
// State ignores the context.
type lookupsStateIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// State is the override of the planted defect.
func (d *lookupsStateIgnoresCancelDefect) State(
	p0 context.Context,
	p1 string,
) (r0 Switch, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.State(p0, p1)
}

// lookupsStateIgnoresCancel returns the planted defect of a Lookups whose
// State ignores the context over a fresh subject.
func lookupsStateIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsStateIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsStateValueBesideDefect is the planted defect of a Lookups whose
// State returns a value beside an error.
type lookupsStateValueBesideDefect struct {
	lookupsOverEmbedded
}

// State is the override of the planted defect.
func (d *lookupsStateValueBesideDefect) State(p0 context.Context, p1 string) (r0 Switch, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.State(p0, p1)
	if r1 != nil {
		r0 = Switch{On: Flag(true)}
	}
	return r0, r1
}

// lookupsStateValueBeside returns the planted defect of a Lookups whose
// State returns a value beside an error over a fresh subject.
func lookupsStateValueBeside(h *LookupsHandle) Lookups {
	return &lookupsStateValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsTogglesPanicsDefect is the planted defect of a Lookups whose
// Toggles panics.
type lookupsTogglesPanicsDefect struct {
	lookupsOverEmbedded
}

// Toggles is the override of the planted defect.
func (*lookupsTogglesPanicsDefect) Toggles(p0 context.Context, p1 string) (r0 []bool, r1 error) {
	panic("the planted defect of a Lookups whose Toggles panics")
}

// lookupsTogglesPanics returns the planted defect of a Lookups whose
// Toggles panics over a fresh subject.
func lookupsTogglesPanics(h *LookupsHandle) Lookups {
	return &lookupsTogglesPanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsTogglesIgnoresCancelDefect is the planted defect of a Lookups
// whose Toggles ignores the context.
type lookupsTogglesIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// Toggles is the override of the planted defect.
func (d *lookupsTogglesIgnoresCancelDefect) Toggles(
	p0 context.Context,
	p1 string,
) (r0 []bool, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.Toggles(p0, p1)
}

// lookupsTogglesIgnoresCancel returns the planted defect of a Lookups
// whose Toggles ignores the context over a fresh subject.
func lookupsTogglesIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsTogglesIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsTogglesValueBesideDefect is the planted defect of a Lookups whose
// Toggles returns a value beside an error.
type lookupsTogglesValueBesideDefect struct {
	lookupsOverEmbedded
}

// Toggles is the override of the planted defect.
func (d *lookupsTogglesValueBesideDefect) Toggles(
	p0 context.Context,
	p1 string,
) (r0 []bool, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.Toggles(p0, p1)
	if r1 != nil {
		r0 = []bool{false}
	}
	return r0, r1
}

// lookupsTogglesValueBeside returns the planted defect of a Lookups whose
// Toggles returns a value beside an error over a fresh subject.
func lookupsTogglesValueBeside(h *LookupsHandle) Lookups {
	return &lookupsTogglesValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsAmountPanicsDefect is the planted defect of a Lookups whose
// Amount panics.
type lookupsAmountPanicsDefect struct {
	lookupsOverEmbedded
}

// Amount is the override of the planted defect.
func (*lookupsAmountPanicsDefect) Amount(p0 context.Context, p1 string) (r0 Amount, r1 error) {
	panic("the planted defect of a Lookups whose Amount panics")
}

// lookupsAmountPanics returns the planted defect of a Lookups whose Amount
// panics over a fresh subject.
func lookupsAmountPanics(h *LookupsHandle) Lookups {
	return &lookupsAmountPanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsAmountIgnoresCancelDefect is the planted defect of a Lookups
// whose Amount ignores the context.
type lookupsAmountIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// Amount is the override of the planted defect.
func (d *lookupsAmountIgnoresCancelDefect) Amount(
	p0 context.Context,
	p1 string,
) (r0 Amount, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.Amount(p0, p1)
}

// lookupsAmountIgnoresCancel returns the planted defect of a Lookups whose
// Amount ignores the context over a fresh subject.
func lookupsAmountIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsAmountIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsAmountValueBesideDefect is the planted defect of a Lookups whose
// Amount returns a value beside an error.
type lookupsAmountValueBesideDefect struct {
	lookupsOverEmbedded
}

// Amount is the override of the planted defect.
func (d *lookupsAmountValueBesideDefect) Amount(
	p0 context.Context,
	p1 string,
) (r0 Amount, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.Amount(p0, p1)
	if r1 != nil {
		r0 = Amount(5)
	}
	return r0, r1
}

// lookupsAmountValueBeside returns the planted defect of a Lookups whose
// Amount returns a value beside an error over a fresh subject.
func lookupsAmountValueBeside(h *LookupsHandle) Lookups {
	return &lookupsAmountValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsRatioPanicsDefect is the planted defect of a Lookups whose Ratio
// panics.
type lookupsRatioPanicsDefect struct {
	lookupsOverEmbedded
}

// Ratio is the override of the planted defect.
func (*lookupsRatioPanicsDefect) Ratio(p0 context.Context, p1 string) (r0 Ratio, r1 error) {
	panic("the planted defect of a Lookups whose Ratio panics")
}

// lookupsRatioPanics returns the planted defect of a Lookups whose Ratio
// panics over a fresh subject.
func lookupsRatioPanics(h *LookupsHandle) Lookups {
	return &lookupsRatioPanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsRatioIgnoresCancelDefect is the planted defect of a Lookups whose
// Ratio ignores the context.
type lookupsRatioIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// Ratio is the override of the planted defect.
func (d *lookupsRatioIgnoresCancelDefect) Ratio(
	p0 context.Context,
	p1 string,
) (r0 Ratio, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.Ratio(p0, p1)
}

// lookupsRatioIgnoresCancel returns the planted defect of a Lookups whose
// Ratio ignores the context over a fresh subject.
func lookupsRatioIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsRatioIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsRatioValueBesideDefect is the planted defect of a Lookups whose
// Ratio returns a value beside an error.
type lookupsRatioValueBesideDefect struct {
	lookupsOverEmbedded
}

// Ratio is the override of the planted defect.
func (d *lookupsRatioValueBesideDefect) Ratio(p0 context.Context, p1 string) (r0 Ratio, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.Ratio(p0, p1)
	if r1 != nil {
		r0 = Ratio(0.5)
	}
	return r0, r1
}

// lookupsRatioValueBeside returns the planted defect of a Lookups whose
// Ratio returns a value beside an error over a fresh subject.
func lookupsRatioValueBeside(h *LookupsHandle) Lookups {
	return &lookupsRatioValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsCodePanicsDefect is the planted defect of a Lookups whose Code
// panics.
type lookupsCodePanicsDefect struct {
	lookupsOverEmbedded
}

// Code is the override of the planted defect.
func (*lookupsCodePanicsDefect) Code(p0 context.Context, p1 string) (r0 Code, r1 error) {
	panic("the planted defect of a Lookups whose Code panics")
}

// lookupsCodePanics returns the planted defect of a Lookups whose Code
// panics over a fresh subject.
func lookupsCodePanics(h *LookupsHandle) Lookups {
	return &lookupsCodePanicsDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsCodeIgnoresCancelDefect is the planted defect of a Lookups whose
// Code ignores the context.
type lookupsCodeIgnoresCancelDefect struct {
	lookupsOverEmbedded
}

// Code is the override of the planted defect.
func (d *lookupsCodeIgnoresCancelDefect) Code(p0 context.Context, p1 string) (r0 Code, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lookupsOverEmbedded.Code(p0, p1)
}

// lookupsCodeIgnoresCancel returns the planted defect of a Lookups whose
// Code ignores the context over a fresh subject.
func lookupsCodeIgnoresCancel(h *LookupsHandle) Lookups {
	return &lookupsCodeIgnoresCancelDefect{lookupsOverEmbedded: h.Fresh()}
}

// lookupsCodeValueBesideDefect is the planted defect of a Lookups whose
// Code returns a value beside an error.
type lookupsCodeValueBesideDefect struct {
	lookupsOverEmbedded
}

// Code is the override of the planted defect.
func (d *lookupsCodeValueBesideDefect) Code(p0 context.Context, p1 string) (r0 Code, r1 error) {
	r0, r1 = d.lookupsOverEmbedded.Code(p0, p1)
	if r1 != nil {
		r0 = Code("code")
	}
	return r0, r1
}

// lookupsCodeValueBeside returns the planted defect of a Lookups whose
// Code returns a value beside an error over a fresh subject.
func lookupsCodeValueBeside(h *LookupsHandle) Lookups {
	return &lookupsCodeValueBesideDefect{lookupsOverEmbedded: h.Fresh()}
}

// The derived checks of Labels:
//
// Labels.Put/smoke: Put returns without a panic for derived inputs.
//
// Labels.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Labels.Put/cancel: Put returns an error for a cancelled context.
//
// Labels.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Labels.Get/smoke: Get returns without a panic for derived inputs.
//
// Labels.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Labels.Get/cancel: Get returns an error for a cancelled context.
//
// Labels.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Labels.Get/read-after-write: Get returns what Put stored before it.
//
// Labels.Get/read-after-write#generated: Get returns what Put stored
// before it.
//
// Labels.Get/reads-each-last-write: Get returns what Put wrote right
// before it, over any sequence of calls.

// The derived checks of Labels that witness their laws at the shallow end:
//
// Labels.Put/smoke
//
// Labels.Get/smoke
//
// Labels.Get/read-after-write
//
// Labels.Get/read-after-write#generated

// The unclaimed classifications of Labels:
//
// The shape:composite-writer of Labels.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Labels writes.
type (
	LabelsHandle = conform.Handle[Labels]
	LabelsCheck  = conform.Check[Labels]
	LabelsOption = conform.Option[Labels]
)

// LabelsChecks contains the name of each callable of Labels, with the
// identities of its checks.
var LabelsChecks = struct {
	Put labelsPut
	Get labelsGet
}{
	Put: labelsPut{"Labels.Put"},
	Get: labelsGet{"Labels.Get"},
}

// labelsPut is the name of Labels.Put, with the identities of its checks.
type labelsPut struct {
	conform.Callable[Labels]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c labelsPut) Smoke() conform.ID[Labels] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c labelsPut) SmokeGenerated() conform.ID[Labels] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c labelsPut) Cancel() conform.ID[Labels] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c labelsPut) CancelGenerated() conform.ID[Labels] {
	return c.Check("cancel#generated")
}

// labelsGet is the name of Labels.Get, with the identities of its checks.
type labelsGet struct {
	conform.Callable[Labels]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c labelsGet) Smoke() conform.ID[Labels] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c labelsGet) SmokeGenerated() conform.ID[Labels] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c labelsGet) Cancel() conform.ID[Labels] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c labelsGet) CancelGenerated() conform.ID[Labels] {
	return c.Check("cancel#generated")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c labelsGet) ReadAfterWrite() conform.ID[Labels] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that Get
// returns what Put stored before it.
func (c labelsGet) ReadAfterWriteGenerated() conform.ID[Labels] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Get returns
// what Put wrote right before it, over any sequence of calls.
func (c labelsGet) ReadsEachLastWrite() conform.ID[Labels] {
	return c.Check("reads-each-last-write")
}

// LabelsSuite returns the checks of Labels, and the calls and the budgets
// of its callables.
func LabelsSuite() conform.Suite[Labels] {
	return conform.Suite[Labels]{
		Declaration: "Labels",
		Package:     "example.com/shop/drawn",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Labels]{
			{
				ID:          LabelsChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LabelsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("name-a"), Label{Name: "Name-a"})
					})
				},
				Defect: labelsPutPanics,
			},
			{
				ID:          LabelsChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LabelsHandle, c *prop.Case) {
					drawnName := c.Draw(prop.Of[string](), "name")
					drawnV := c.Draw(prop.Of[Label](), "v")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnName, drawnV)
					})
				},
				Defect: labelsPutPanics,
			},
			{
				ID:          LabelsChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LabelsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("name-a"), Label{Name: "Name-a"})
					})
				},
				Defect: labelsPutIgnoresCancel,
			},
			{
				ID:          LabelsChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LabelsHandle, c *prop.Case) {
					drawnName := c.Draw(prop.Of[string](), "name")
					drawnV := c.Draw(prop.Of[Label](), "v")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnName, drawnV)
					})
				},
				Defect: labelsPutIgnoresCancel,
			},
			{
				ID:          LabelsChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LabelsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, func() string { return "" })
						return err
					})
				},
				Defect: labelsGetPanics,
			},
			{
				ID:          LabelsChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LabelsHandle, c *prop.Case) {
					nameResults := c.Draw(prop.List(prop.Of[string]()), "name")
					drawnName := func() func() string {
						var nameAt atomic.Int64
						return func() (r0 string) {
							i := int(nameAt.Add(1) - 1)
							if i >= len(nameResults) {
								return r0
							}
							return nameResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnName())
						return err
					})
				},
				Defect: labelsGetPanics,
			},
			{
				ID:          LabelsChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LabelsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, func() string { return "" })
						return err
					})
				},
				Defect: labelsGetIgnoresCancel,
			},
			{
				ID:          LabelsChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LabelsHandle, c *prop.Case) {
					nameResults := c.Draw(prop.List(prop.Of[string]()), "name")
					drawnName := func() func() string {
						var nameAt atomic.Int64
						return func() (r0 string) {
							i := int(nameAt.Add(1) - 1)
							if i >= len(nameResults) {
								return r0
							}
							return nameResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnName())
						return err
					})
				},
				Defect: labelsGetIgnoresCancel,
			},
			{
				ID:          LabelsChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LabelsHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Label, error) {
						return h.Subject().Get(ctx, func() string { return "" })
					}, func(ctx context.Context) (Label, error) {
						w := Label{Name: "Name-a"}
						err := h.Subject().Put(ctx, string("name-a"), w)
						return w, err
					})
				},
				Defect: labelsGetReadsAnotherValue,
			},
			{
				ID:          LabelsChecks.Get.ReadAfterWriteGenerated(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *LabelsHandle, c *prop.Case) {
					nameResults := c.Draw(prop.List(prop.Of[string]()), "name")
					drawnName := func() func() string {
						var nameAt atomic.Int64
						return func() (r0 string) {
							i := int(nameAt.Add(1) - 1)
							if i >= len(nameResults) {
								return r0
							}
							return nameResults[i]
						}
					}
					drawnName0 := c.Draw(prop.Of[string](), "name")
					drawnV := c.Draw(prop.Of[Label](), "v")
					probe.WriteThenRead(h, func(ctx context.Context) (Label, error) {
						return h.Subject().Get(ctx, drawnName())
					}, func(ctx context.Context) (Label, error) {
						w := drawnV
						err := h.Subject().Put(ctx, drawnName0, w)
						return w, err
					})
				},
				Defect: labelsGetReadsAnotherValue,
			},
			{
				ID:          LabelsChecks.Get.ReadsEachLastWrite(),
				Statement:   "Get returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Labels]{LabelsChecks.Get.ReadAfterWrite()},
				Property: func(h *LabelsHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						labelsActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							values, err := labelsGetCall(
								ctx,
								h.Subject(),
								[]any{func() string { return "" }},
							)
							return conform.OutputOf(err, values...)
						},
						func(args []any) Label {
							v, _ := args[1].(Label)
							return v
						},
					)
				},
				Defect: labelsGetReadsAStaleValueAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Labels.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Labels]{
			{
				Calls:    labelsPutCalls,
				Over:     newLabelsOver,
				Callable: LabelsChecks.Put.Callable,
			},
			{
				Calls:    labelsGetCalls,
				Over:     newLabelsOver,
				Callable: LabelsChecks.Get.Callable,
			},
		},
	}
}

// RunLabelsSuite applies the checks of Labels, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLabelsSuite[T Labels](t *testing.T, newSubject func() T, opts ...LabelsOption) {
	t.Helper()
	RunLabelsSuiteFromHandle(t, func(*LabelsHandle) T { return newSubject() }, opts...)
}

// RunLabelsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLabelsSuiteFromHandle[T Labels](
	t *testing.T,
	newSubject func(h *LabelsHandle) T,
	opts ...LabelsOption,
) {
	t.Helper()
	conform.Run(t, LabelsSuite(), labelsSubjects(newSubject), opts...)
}

// FuzzLabelsSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLabelsSuite[T Labels](
	f *testing.F,
	callable conform.Callable[Labels],
	newSubject func(h *LabelsHandle) T,
	opts ...LabelsOption,
) {
	f.Helper()
	conform.Fuzz(f, LabelsSuite(), callable, labelsSubjects(newSubject), opts...)
}

// RunLabelsBudgets applies the allocation ceilings of Labels and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLabelsBudgets[T Labels](
	t *testing.T,
	newSubject func(h *LabelsHandle) T,
	opts ...LabelsOption,
) {
	t.Helper()
	conform.RunBudgets(t, LabelsSuite(), labelsSubjects(newSubject), opts...)
}

// BenchLabelsSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLabelsSuite[T Labels](
	b conform.B,
	callable conform.Callable[Labels],
	newSubject func(h *LabelsHandle) T,
	opts ...LabelsOption,
) {
	b.Helper()
	conform.Bench(b, LabelsSuite(), callable, labelsSubjects(newSubject), opts...)
}

// labelsSubjects returns the factory of the subjects of Labels that calls
// newSubject with the handle of each check.
func labelsSubjects[T Labels](newSubject func(h *LabelsHandle) T) func(h *LabelsHandle) Labels {
	return func(h *LabelsHandle) Labels { return newSubject(h) }
}

// labelsPutCalls generates the calls of Labels.Put with the inputs that
// each case draws.
var labelsPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Labels) error {
	name := c.Draw(prop.Of[string](), "name")
	v := c.Draw(prop.Of[Label](), "v")
	return func(ctx context.Context, subject Labels) error {
		return subject.Put(ctx, name, v)
	}
})

// labelsGetCalls generates the calls of Labels.Get with the inputs that
// each case draws.
var labelsGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Labels) error {
	nameResults := c.Draw(prop.List(prop.Of[string]()), "name")
	var nameAt atomic.Int64
	name := func() (r0 string) {
		i := int(nameAt.Add(1) - 1)
		if i >= len(nameResults) {
			return r0
		}
		return nameResults[i]
	}
	return func(ctx context.Context, subject Labels) error {
		_, err := subject.Get(ctx, name)
		return err
	}
})

// labelsOverEmbedded is Labels under the name of the field that labelsOver
// embeds.
type labelsOverEmbedded = Labels

// labelsOver is a subject of Labels whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type labelsOver struct {
	labelsOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o labelsOver) Put(context.Context, string, Label) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o labelsOver) Get(context.Context, func() string) (r0 Label, r1 error) {
	o.spend()
	return r0, r1
}

// newLabelsOver returns the planted defect of a ceiling of Labels, whose
// callables call spend.
func newLabelsOver(_ *LabelsHandle, spend func()) Labels {
	return labelsOver{spend: spend}
}

// labelsPutPanicsDefect is the planted defect of a Labels whose Put
// panics.
type labelsPutPanicsDefect struct {
	labelsOverEmbedded
}

// Put is the override of the planted defect.
func (*labelsPutPanicsDefect) Put(p0 context.Context, p1 string, p2 Label) (r0 error) {
	panic("the planted defect of a Labels whose Put panics")
}

// labelsPutPanics returns the planted defect of a Labels whose Put panics
// over a fresh subject.
func labelsPutPanics(h *LabelsHandle) Labels {
	return &labelsPutPanicsDefect{labelsOverEmbedded: h.Fresh()}
}

// labelsPutIgnoresCancelDefect is the planted defect of a Labels whose Put
// ignores the context.
type labelsPutIgnoresCancelDefect struct {
	labelsOverEmbedded
}

// Put is the override of the planted defect.
func (d *labelsPutIgnoresCancelDefect) Put(p0 context.Context, p1 string, p2 Label) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.labelsOverEmbedded.Put(p0, p1, p2)
}

// labelsPutIgnoresCancel returns the planted defect of a Labels whose Put
// ignores the context over a fresh subject.
func labelsPutIgnoresCancel(h *LabelsHandle) Labels {
	return &labelsPutIgnoresCancelDefect{labelsOverEmbedded: h.Fresh()}
}

// labelsGetPanicsDefect is the planted defect of a Labels whose Get
// panics.
type labelsGetPanicsDefect struct {
	labelsOverEmbedded
}

// Get is the override of the planted defect.
func (*labelsGetPanicsDefect) Get(p0 context.Context, p1 func() string) (r0 Label, r1 error) {
	panic("the planted defect of a Labels whose Get panics")
}

// labelsGetPanics returns the planted defect of a Labels whose Get panics
// over a fresh subject.
func labelsGetPanics(h *LabelsHandle) Labels {
	return &labelsGetPanicsDefect{labelsOverEmbedded: h.Fresh()}
}

// labelsGetIgnoresCancelDefect is the planted defect of a Labels whose Get
// ignores the context.
type labelsGetIgnoresCancelDefect struct {
	labelsOverEmbedded
}

// Get is the override of the planted defect.
func (d *labelsGetIgnoresCancelDefect) Get(
	p0 context.Context,
	p1 func() string,
) (r0 Label, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.labelsOverEmbedded.Get(p0, p1)
}

// labelsGetIgnoresCancel returns the planted defect of a Labels whose Get
// ignores the context over a fresh subject.
func labelsGetIgnoresCancel(h *LabelsHandle) Labels {
	return &labelsGetIgnoresCancelDefect{labelsOverEmbedded: h.Fresh()}
}

// labelsGetReadsAnotherValueDefect is the planted defect of a Labels whose
// Get returns a value that Put did not write.
type labelsGetReadsAnotherValueDefect struct {
	labelsOverEmbedded
}

// Get is the override of the planted defect.
func (d *labelsGetReadsAnotherValueDefect) Get(
	p0 context.Context,
	p1 func() string,
) (r0 Label, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.labelsOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Label{Name: "Name-b"}
	}
	return r0, r1
}

// labelsGetReadsAnotherValue returns the planted defect of a Labels whose
// Get returns a value that Put did not write over a fresh subject.
func labelsGetReadsAnotherValue(h *LabelsHandle) Labels {
	return &labelsGetReadsAnotherValueDefect{labelsOverEmbedded: h.Fresh()}
}

// labelsGetReadsAStaleValueAfterThreeDefect is the planted defect of a
// Labels whose Get returns its first value again from its fourth call on.
type labelsGetReadsAStaleValueAfterThreeDefect struct {
	labelsOverEmbedded
	calls int
	first *Label
}

// Get is the override of the planted defect.
func (d *labelsGetReadsAStaleValueAfterThreeDefect) Get(
	p0 context.Context,
	p1 func() string,
) (r0 Label, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.labelsOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// labelsGetReadsAStaleValueAfterThree returns the planted defect of a
// Labels whose Get returns its first value again from its fourth call on
// over a fresh subject.
func labelsGetReadsAStaleValueAfterThree(h *LabelsHandle) Labels {
	return &labelsGetReadsAStaleValueAfterThreeDefect{labelsOverEmbedded: h.Fresh()}
}

// labelsActions returns the actions of the generated sequences of Labels
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func labelsActions(c *prop.Case, subject func() Labels) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Input: func(c *prop.Case, _ struct{}) any {
				name := c.Draw(prop.SampledFrom(keysString...), "put.name")
				v := c.Draw(prop.Of[Label](), "put.v")
				return []any{name, v}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Put", args, args[0])
				values, err := labelsPutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				nameResults := c.Draw(prop.List(prop.Of[string]()), "name")
				var nameAt atomic.Int64
				name := func() (r0 string) {
					i := int(nameAt.Add(1) - 1)
					if i >= len(nameResults) {
						return r0
					}
					return nameResults[i]
				}
				return []any{name}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args)
				values, err := labelsGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// labelsPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func labelsPutCall(ctx context.Context, subject Labels, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(Label)
	return nil, subject.Put(ctx, a0, a1)
}

// labelsGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func labelsGetCall(ctx context.Context, subject Labels, args []any) ([]any, error) {
	a0, _ := args[0].(func() string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}
