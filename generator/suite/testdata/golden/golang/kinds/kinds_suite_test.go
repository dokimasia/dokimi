package kinds_test

import (
	"testing"

	"example.com/shop/kinds"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// paramsSubjectEmbedded is kinds.Params under the name of the field that
// paramsSubject embeds.
type paramsSubjectEmbedded = kinds.Params

// paramsSubject is a value of Params that only TestParamsSuite builds.
type paramsSubject struct{ paramsSubjectEmbedded }

// TestParamsSuite checks the suite of Params and the two forms of its
// application.
func TestParamsSuite(t *testing.T) {
	t.Parallel()

	t.Run("ParamsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, kinds.ParamsSuite())
	})

	t.Run("RunParamsSuite", func(t *testing.T) {
		t.Parallel()
		var built *paramsSubject
		kinds.RunParamsSuite(
			t,
			func() *paramsSubject {
				built = &paramsSubject{}
				return built
			},
			conform.With(kinds.ParamsCheck{
				ID:        "Params/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.ParamsHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Params(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.ParamsHandle) kinds.Params { return &paramsSubject{} },
			}),
			conform.Excuse(
				kinds.ParamsChecks.Bool.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int8.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int16.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int32.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint8.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint16.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint32.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uintptr.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Float32.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Float64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Complex64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Complex128.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Array.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Func.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Map.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Slice.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.String.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunParamsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *paramsSubject
		kinds.RunParamsSuiteFromHandle(
			t,
			func(*kinds.ParamsHandle) *paramsSubject {
				built = &paramsSubject{}
				return built
			},
			conform.With(kinds.ParamsCheck{
				ID:        "Params/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.ParamsHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Params(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.ParamsHandle) kinds.Params { return &paramsSubject{} },
			}),
			conform.Excuse(
				kinds.ParamsChecks.Bool.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int8.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int16.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int32.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Int64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint8.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint16.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint32.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uint64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Uintptr.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Float32.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Float64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Complex64.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Complex128.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Array.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Func.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Map.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.Slice.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.ParamsChecks.String.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// paramsIdle returns the planted defect of a ceiling of Params with a
// spend of nothing, a subject that answers each call at once.
func paramsIdle(h *kinds.ParamsHandle) kinds.Params {
	return kinds.ParamsSuite().Calls[0].Over(h, func() {})
}

// TestParamsSuiteBudgets applies the allocation ceilings of the suite of
// Params to paramsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestParamsSuiteBudgets(t *testing.T) {
	kinds.RunParamsBudgets(
		t,
		paramsIdle,
		conform.Ceilings(kinds.ParamsChecks.Bool.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Int.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Int8.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Int16.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Int32.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Int64.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Uint.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Uint8.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			kinds.ParamsChecks.Uint16.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Uint32.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Uint64.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Uintptr.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Float32.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Float64.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Complex64.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Complex128.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(kinds.ParamsChecks.Array.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Chan.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Func.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Map.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Ptr.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.ParamsChecks.Slice.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			kinds.ParamsChecks.String.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.Struct.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			kinds.ParamsChecks.UnsafePointer.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteBool registers the fuzz target of Params.Bool on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteBool(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Bool.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteInt registers the fuzz target of Params.Int on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Int.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteInt8 registers the fuzz target of Params.Int8 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt8(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Int8.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteInt16 registers the fuzz target of Params.Int16 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt16(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Int16.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteInt32 registers the fuzz target of Params.Int32 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt32(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Int32.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteInt64 registers the fuzz target of Params.Int64 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt64(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Int64.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteUint registers the fuzz target of Params.Uint on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Uint.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteUint8 registers the fuzz target of Params.Uint8 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint8(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Uint8.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteUint16 registers the fuzz target of Params.Uint16 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint16(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Uint16.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteUint32 registers the fuzz target of Params.Uint32 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint32(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Uint32.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteUint64 registers the fuzz target of Params.Uint64 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint64(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Uint64.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteUintptr registers the fuzz target of Params.Uintptr on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUintptr(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Uintptr.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteFloat32 registers the fuzz target of Params.Float32 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteFloat32(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Float32.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteFloat64 registers the fuzz target of Params.Float64 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteFloat64(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Float64.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteComplex64 registers the fuzz target of Params.Complex64
// on paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteComplex64(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Complex64.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteComplex128 registers the fuzz target of Params.Complex128
// on paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteComplex128(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Complex128.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteArray registers the fuzz target of Params.Array on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteArray(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Array.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteChan registers the fuzz target of Params.Chan on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteChan(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Chan.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteFunc registers the fuzz target of Params.Func on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteFunc(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Func.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteInterface registers the fuzz target of Params.Interface
// on paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInterface(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Interface.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteMap registers the fuzz target of Params.Map on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteMap(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Map.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuitePtr registers the fuzz target of Params.Ptr on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuitePtr(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Ptr.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteSlice registers the fuzz target of Params.Slice on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteSlice(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Slice.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteString registers the fuzz target of Params.String on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteString(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.String.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteStruct registers the fuzz target of Params.Struct on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteStruct(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.Struct.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzParamsSuiteUnsafePointer registers the fuzz target of
// Params.UnsafePointer on paramsIdle, with every check of the suite
// excused.
func FuzzParamsSuiteUnsafePointer(f *testing.F) {
	kinds.FuzzParamsSuite(
		f,
		kinds.ParamsChecks.UnsafePointer.Callable,
		paramsIdle,
		conform.Excuse(
			kinds.ParamsChecks.Bool.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Int64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint8.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint16.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uint64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Uintptr.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float32.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Float64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex64.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Complex128.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Array.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Func.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Map.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.Slice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.ParamsChecks.String.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkParamsSuite measures each callable of the suite of Params on
// paramsIdle.
func BenchmarkParamsSuite(b *testing.B) {
	b.Run("Bool", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Bool.Callable, paramsIdle)
	})
	b.Run("Int", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Int.Callable, paramsIdle)
	})
	b.Run("Int8", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Int8.Callable, paramsIdle)
	})
	b.Run("Int16", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Int16.Callable, paramsIdle)
	})
	b.Run("Int32", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Int32.Callable, paramsIdle)
	})
	b.Run("Int64", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Int64.Callable, paramsIdle)
	})
	b.Run("Uint", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Uint.Callable, paramsIdle)
	})
	b.Run("Uint8", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Uint8.Callable, paramsIdle)
	})
	b.Run("Uint16", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Uint16.Callable, paramsIdle)
	})
	b.Run("Uint32", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Uint32.Callable, paramsIdle)
	})
	b.Run("Uint64", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Uint64.Callable, paramsIdle)
	})
	b.Run("Uintptr", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Uintptr.Callable, paramsIdle)
	})
	b.Run("Float32", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Float32.Callable, paramsIdle)
	})
	b.Run("Float64", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Float64.Callable, paramsIdle)
	})
	b.Run("Complex64", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Complex64.Callable, paramsIdle)
	})
	b.Run("Complex128", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Complex128.Callable, paramsIdle)
	})
	b.Run("Array", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Array.Callable, paramsIdle)
	})
	b.Run("Chan", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Chan.Callable, paramsIdle)
	})
	b.Run("Func", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Func.Callable, paramsIdle)
	})
	b.Run("Interface", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Interface.Callable, paramsIdle)
	})
	b.Run("Map", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Map.Callable, paramsIdle)
	})
	b.Run("Ptr", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Ptr.Callable, paramsIdle)
	})
	b.Run("Slice", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Slice.Callable, paramsIdle)
	})
	b.Run("String", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.String.Callable, paramsIdle)
	})
	b.Run("Struct", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.Struct.Callable, paramsIdle)
	})
	b.Run("UnsafePointer", func(b *testing.B) {
		kinds.BenchParamsSuite(b, kinds.ParamsChecks.UnsafePointer.Callable, paramsIdle)
	})
}

// everySubjectEmbedded is kinds.Every under the name of the field that
// everySubject embeds.
type everySubjectEmbedded = kinds.Every

// everySubject is a value of Every that only TestEverySuite builds.
type everySubject struct{ everySubjectEmbedded }

// TestEverySuite checks the suite of Every and the two forms of its
// application.
func TestEverySuite(t *testing.T) {
	t.Parallel()

	t.Run("EverySuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, kinds.EverySuite())
	})

	t.Run("RunEverySuite", func(t *testing.T) {
		t.Parallel()
		var built *everySubject
		kinds.RunEverySuite(
			t,
			func() *everySubject {
				built = &everySubject{}
				return built
			},
			conform.With(kinds.EveryCheck{
				ID:        "Every/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.EveryHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Every(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.EveryHandle) kinds.Every { return &everySubject{} },
			}),
			conform.Excuse(
				kinds.EveryChecks.Fields.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.EveryChecks.Elements.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunEverySuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *everySubject
		kinds.RunEverySuiteFromHandle(
			t,
			func(*kinds.EveryHandle) *everySubject {
				built = &everySubject{}
				return built
			},
			conform.With(kinds.EveryCheck{
				ID:        "Every/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.EveryHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Every(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.EveryHandle) kinds.Every { return &everySubject{} },
			}),
			conform.Excuse(
				kinds.EveryChecks.Fields.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.EveryChecks.Elements.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// everyIdle returns the planted defect of a ceiling of Every with a spend
// of nothing, a subject that answers each call at once.
func everyIdle(h *kinds.EveryHandle) kinds.Every {
	return kinds.EverySuite().Calls[0].Over(h, func() {})
}

// TestEverySuiteBudgets applies the allocation ceilings of the suite of
// Every to everyIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestEverySuiteBudgets(t *testing.T) {
	kinds.RunEveryBudgets(
		t,
		everyIdle,
		conform.Ceilings(kinds.EveryChecks.Fields.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			kinds.EveryChecks.Elements.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			kinds.EveryChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.EveryChecks.Elements.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzEverySuiteFields registers the fuzz target of Every.Fields on
// everyIdle, with every check of the suite excused.
func FuzzEverySuiteFields(f *testing.F) {
	kinds.FuzzEverySuite(
		f,
		kinds.EveryChecks.Fields.Callable,
		everyIdle,
		conform.Excuse(
			kinds.EveryChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.EveryChecks.Elements.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzEverySuiteElements registers the fuzz target of Every.Elements on
// everyIdle, with every check of the suite excused.
func FuzzEverySuiteElements(f *testing.F) {
	kinds.FuzzEverySuite(
		f,
		kinds.EveryChecks.Elements.Callable,
		everyIdle,
		conform.Excuse(
			kinds.EveryChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.EveryChecks.Elements.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzEverySuiteResults registers the fuzz target of Every.Results on
// everyIdle, with every check of the suite excused.
func FuzzEverySuiteResults(f *testing.F) {
	kinds.FuzzEverySuite(
		f,
		kinds.EveryChecks.Results.Callable,
		everyIdle,
		conform.Excuse(
			kinds.EveryChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.EveryChecks.Elements.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkEverySuite measures each callable of the suite of Every on
// everyIdle.
func BenchmarkEverySuite(b *testing.B) {
	b.Run("Fields", func(b *testing.B) {
		kinds.BenchEverySuite(b, kinds.EveryChecks.Fields.Callable, everyIdle)
	})
	b.Run("Elements", func(b *testing.B) {
		kinds.BenchEverySuite(b, kinds.EveryChecks.Elements.Callable, everyIdle)
	})
	b.Run("Results", func(b *testing.B) {
		kinds.BenchEverySuite(b, kinds.EveryChecks.Results.Callable, everyIdle)
	})
}

// placerSubjectEmbedded is kinds.Placer under the name of the field that
// placerSubject embeds.
type placerSubjectEmbedded = kinds.Placer

// placerSubject is a value of Placer that only TestPlacerSuite builds.
type placerSubject struct{ placerSubjectEmbedded }

// TestPlacerSuite checks the suite of Placer and the two forms of its
// application.
func TestPlacerSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlacerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, kinds.PlacerSuite())
	})

	t.Run("RunPlacerSuite", func(t *testing.T) {
		t.Parallel()
		var built *placerSubject
		kinds.RunPlacerSuite(
			t,
			func() *placerSubject {
				built = &placerSubject{}
				return built
			},
			conform.With(kinds.PlacerCheck{
				ID:        "Placer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.PlacerHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Placer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.PlacerHandle) kinds.Placer { return &placerSubject{} },
			}),
		)
	})

	t.Run("RunPlacerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *placerSubject
		kinds.RunPlacerSuiteFromHandle(
			t,
			func(*kinds.PlacerHandle) *placerSubject {
				built = &placerSubject{}
				return built
			},
			conform.With(kinds.PlacerCheck{
				ID:        "Placer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.PlacerHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Placer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.PlacerHandle) kinds.Placer { return &placerSubject{} },
			}),
		)
	})
}

// placerIdle returns the planted defect of a ceiling of Placer with a
// spend of nothing, a subject that answers each call at once.
func placerIdle(h *kinds.PlacerHandle) kinds.Placer {
	return kinds.PlacerSuite().Calls[0].Over(h, func() {})
}

// TestPlacerSuiteBudgets applies the allocation ceilings of the suite of
// Placer to placerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestPlacerSuiteBudgets(t *testing.T) {
	kinds.RunPlacerBudgets(
		t,
		placerIdle,
	)
}

// FuzzPlacerSuitePlace registers the fuzz target of Placer.Place on
// placerIdle, with every check of the suite excused.
func FuzzPlacerSuitePlace(f *testing.F) {
	kinds.FuzzPlacerSuite(
		f,
		kinds.PlacerChecks.Place.Callable,
		placerIdle,
	)
}

// BenchmarkPlacerSuite measures each callable of the suite of Placer on
// placerIdle.
func BenchmarkPlacerSuite(b *testing.B) {
	b.Run("Place", func(b *testing.B) {
		kinds.BenchPlacerSuite(b, kinds.PlacerChecks.Place.Callable, placerIdle)
	})
}

// routerSubjectEmbedded is kinds.Router under the name of the field that
// routerSubject embeds.
type routerSubjectEmbedded = kinds.Router

// routerSubject is a value of Router that only TestRouterSuite builds.
type routerSubject struct{ routerSubjectEmbedded }

// TestRouterSuite checks the suite of Router and the two forms of its
// application.
func TestRouterSuite(t *testing.T) {
	t.Parallel()

	t.Run("RouterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, kinds.RouterSuite())
	})

	t.Run("RunRouterSuite", func(t *testing.T) {
		t.Parallel()
		var built *routerSubject
		kinds.RunRouterSuite(
			t,
			func() *routerSubject {
				built = &routerSubject{}
				return built
			},
			conform.With(kinds.RouterCheck{
				ID:        "Router/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.RouterHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Router(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.RouterHandle) kinds.Router { return &routerSubject{} },
			}),
		)
	})

	t.Run("RunRouterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *routerSubject
		kinds.RunRouterSuiteFromHandle(
			t,
			func(*kinds.RouterHandle) *routerSubject {
				built = &routerSubject{}
				return built
			},
			conform.With(kinds.RouterCheck{
				ID:        "Router/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.RouterHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Router(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.RouterHandle) kinds.Router { return &routerSubject{} },
			}),
		)
	})
}

// routerIdle returns the planted defect of a ceiling of Router with a
// spend of nothing, a subject that answers each call at once.
func routerIdle(h *kinds.RouterHandle) kinds.Router {
	return kinds.RouterSuite().Calls[0].Over(h, func() {})
}

// TestRouterSuiteBudgets applies the allocation ceilings of the suite of
// Router to routerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRouterSuiteBudgets(t *testing.T) {
	kinds.RunRouterBudgets(
		t,
		routerIdle,
		conform.Ceilings(kinds.RouterChecks.Route.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.RouterChecks.Write.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.RouterChecks.Pass.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzRouterSuiteRoute registers the fuzz target of Router.Route on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteRoute(f *testing.F) {
	kinds.FuzzRouterSuite(
		f,
		kinds.RouterChecks.Route.Callable,
		routerIdle,
	)
}

// FuzzRouterSuiteRead registers the fuzz target of Router.Read on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteRead(f *testing.F) {
	kinds.FuzzRouterSuite(
		f,
		kinds.RouterChecks.Read.Callable,
		routerIdle,
	)
}

// FuzzRouterSuiteWrite registers the fuzz target of Router.Write on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteWrite(f *testing.F) {
	kinds.FuzzRouterSuite(
		f,
		kinds.RouterChecks.Write.Callable,
		routerIdle,
	)
}

// FuzzRouterSuitePass registers the fuzz target of Router.Pass on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuitePass(f *testing.F) {
	kinds.FuzzRouterSuite(
		f,
		kinds.RouterChecks.Pass.Callable,
		routerIdle,
	)
}

// FuzzRouterSuiteWatch registers the fuzz target of Router.Watch on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteWatch(f *testing.F) {
	kinds.FuzzRouterSuite(
		f,
		kinds.RouterChecks.Watch.Callable,
		routerIdle,
	)
}

// BenchmarkRouterSuite measures each callable of the suite of Router on
// routerIdle.
func BenchmarkRouterSuite(b *testing.B) {
	b.Run("Route", func(b *testing.B) {
		kinds.BenchRouterSuite(b, kinds.RouterChecks.Route.Callable, routerIdle)
	})
	b.Run("Read", func(b *testing.B) {
		kinds.BenchRouterSuite(b, kinds.RouterChecks.Read.Callable, routerIdle)
	})
	b.Run("Write", func(b *testing.B) {
		kinds.BenchRouterSuite(b, kinds.RouterChecks.Write.Callable, routerIdle)
	})
	b.Run("Pass", func(b *testing.B) {
		kinds.BenchRouterSuite(b, kinds.RouterChecks.Pass.Callable, routerIdle)
	})
	b.Run("Watch", func(b *testing.B) {
		kinds.BenchRouterSuite(b, kinds.RouterChecks.Watch.Callable, routerIdle)
	})
}

// drawerSubjectEmbedded is kinds.Drawer under the name of the field that
// drawerSubject embeds.
type drawerSubjectEmbedded = kinds.Drawer

// drawerSubject is a value of Drawer that only TestDrawerSuite builds.
type drawerSubject struct{ drawerSubjectEmbedded }

// TestDrawerSuite checks the suite of Drawer and the two forms of its
// application.
func TestDrawerSuite(t *testing.T) {
	t.Parallel()

	t.Run("DrawerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, kinds.DrawerSuite())
	})

	t.Run("RunDrawerSuite", func(t *testing.T) {
		t.Parallel()
		var built *drawerSubject
		kinds.RunDrawerSuite(
			t,
			func() *drawerSubject {
				built = &drawerSubject{}
				return built
			},
			conform.With(kinds.DrawerCheck{
				ID:        "Drawer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.DrawerHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Drawer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.DrawerHandle) kinds.Drawer { return &drawerSubject{} },
			}),
		)
	})

	t.Run("RunDrawerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *drawerSubject
		kinds.RunDrawerSuiteFromHandle(
			t,
			func(*kinds.DrawerHandle) *drawerSubject {
				built = &drawerSubject{}
				return built
			},
			conform.With(kinds.DrawerCheck{
				ID:        "Drawer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.DrawerHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Drawer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.DrawerHandle) kinds.Drawer { return &drawerSubject{} },
			}),
		)
	})
}

// drawerIdle returns the planted defect of a ceiling of Drawer with a
// spend of nothing, a subject that answers each call at once.
func drawerIdle(h *kinds.DrawerHandle) kinds.Drawer {
	return kinds.DrawerSuite().Calls[0].Over(h, func() {})
}

// TestDrawerSuiteBudgets applies the allocation ceilings of the suite of
// Drawer to drawerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestDrawerSuiteBudgets(t *testing.T) {
	kinds.RunDrawerBudgets(
		t,
		drawerIdle,
		conform.Ceilings(kinds.DrawerChecks.Draw.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.DrawerChecks.Count.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.DrawerChecks.Name.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.DrawerChecks.Seal.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.DrawerChecks.Test.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzDrawerSuiteDraw registers the fuzz target of Drawer.Draw on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteDraw(f *testing.F) {
	kinds.FuzzDrawerSuite(
		f,
		kinds.DrawerChecks.Draw.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteCount registers the fuzz target of Drawer.Count on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteCount(f *testing.F) {
	kinds.FuzzDrawerSuite(
		f,
		kinds.DrawerChecks.Count.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteName registers the fuzz target of Drawer.Name on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteName(f *testing.F) {
	kinds.FuzzDrawerSuite(
		f,
		kinds.DrawerChecks.Name.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuitePick registers the fuzz target of Drawer.Pick on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuitePick(f *testing.F) {
	kinds.FuzzDrawerSuite(
		f,
		kinds.DrawerChecks.Pick.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteSeal registers the fuzz target of Drawer.Seal on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteSeal(f *testing.F) {
	kinds.FuzzDrawerSuite(
		f,
		kinds.DrawerChecks.Seal.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteTest registers the fuzz target of Drawer.Test on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteTest(f *testing.F) {
	kinds.FuzzDrawerSuite(
		f,
		kinds.DrawerChecks.Test.Callable,
		drawerIdle,
	)
}

// BenchmarkDrawerSuite measures each callable of the suite of Drawer on
// drawerIdle.
func BenchmarkDrawerSuite(b *testing.B) {
	b.Run("Draw", func(b *testing.B) {
		kinds.BenchDrawerSuite(b, kinds.DrawerChecks.Draw.Callable, drawerIdle)
	})
	b.Run("Count", func(b *testing.B) {
		kinds.BenchDrawerSuite(b, kinds.DrawerChecks.Count.Callable, drawerIdle)
	})
	b.Run("Name", func(b *testing.B) {
		kinds.BenchDrawerSuite(b, kinds.DrawerChecks.Name.Callable, drawerIdle)
	})
	b.Run("Pick", func(b *testing.B) {
		kinds.BenchDrawerSuite(b, kinds.DrawerChecks.Pick.Callable, drawerIdle)
	})
	b.Run("Seal", func(b *testing.B) {
		kinds.BenchDrawerSuite(b, kinds.DrawerChecks.Seal.Callable, drawerIdle)
	})
	b.Run("Test", func(b *testing.B) {
		kinds.BenchDrawerSuite(b, kinds.DrawerChecks.Test.Callable, drawerIdle)
	})
}

// oddSubjectEmbedded is kinds.Odd under the name of the field that
// oddSubject embeds.
type oddSubjectEmbedded = kinds.Odd

// oddSubject is a value of Odd that only TestOddSuite builds.
type oddSubject struct{ oddSubjectEmbedded }

// TestOddSuite checks the suite of Odd and the two forms of its
// application.
func TestOddSuite(t *testing.T) {
	t.Parallel()

	t.Run("OddSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, kinds.OddSuite())
	})

	t.Run("RunOddSuite", func(t *testing.T) {
		t.Parallel()
		var built *oddSubject
		kinds.RunOddSuite(
			t,
			func() *oddSubject {
				built = &oddSubject{}
				return built
			},
			conform.With(kinds.OddCheck{
				ID:        "Odd/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.OddHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Odd(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.OddHandle) kinds.Odd { return &oddSubject{} },
			}),
			conform.Excuse(
				kinds.OddChecks.Keep.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunOddSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *oddSubject
		kinds.RunOddSuiteFromHandle(
			t,
			func(*kinds.OddHandle) *oddSubject {
				built = &oddSubject{}
				return built
			},
			conform.With(kinds.OddCheck{
				ID:        "Odd/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.OddHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Odd(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.OddHandle) kinds.Odd { return &oddSubject{} },
			}),
			conform.Excuse(
				kinds.OddChecks.Keep.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// oddIdle returns the planted defect of a ceiling of Odd with a spend of
// nothing, a subject that answers each call at once.
func oddIdle(h *kinds.OddHandle) kinds.Odd {
	return kinds.OddSuite().Calls[0].Over(h, func() {})
}

// TestOddSuiteBudgets applies the allocation ceilings of the suite of Odd
// to oddIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestOddSuiteBudgets(t *testing.T) {
	kinds.RunOddBudgets(
		t,
		oddIdle,
		conform.Ceilings(kinds.OddChecks.Loop.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			kinds.OddChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOddSuiteTag registers the fuzz target of Odd.Tag on oddIdle, with
// every check of the suite excused.
func FuzzOddSuiteTag(f *testing.F) {
	kinds.FuzzOddSuite(
		f,
		kinds.OddChecks.Tag.Callable,
		oddIdle,
		conform.Excuse(
			kinds.OddChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOddSuiteKeep registers the fuzz target of Odd.Keep on oddIdle, with
// every check of the suite excused.
func FuzzOddSuiteKeep(f *testing.F) {
	kinds.FuzzOddSuite(
		f,
		kinds.OddChecks.Keep.Callable,
		oddIdle,
		conform.Excuse(
			kinds.OddChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOddSuiteEmbed registers the fuzz target of Odd.Embed on oddIdle,
// with every check of the suite excused.
func FuzzOddSuiteEmbed(f *testing.F) {
	kinds.FuzzOddSuite(
		f,
		kinds.OddChecks.Embed.Callable,
		oddIdle,
		conform.Excuse(
			kinds.OddChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOddSuiteLoop registers the fuzz target of Odd.Loop on oddIdle, with
// every check of the suite excused.
func FuzzOddSuiteLoop(f *testing.F) {
	kinds.FuzzOddSuite(
		f,
		kinds.OddChecks.Loop.Callable,
		oddIdle,
		conform.Excuse(
			kinds.OddChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkOddSuite measures each callable of the suite of Odd on oddIdle.
func BenchmarkOddSuite(b *testing.B) {
	b.Run("Tag", func(b *testing.B) {
		kinds.BenchOddSuite(b, kinds.OddChecks.Tag.Callable, oddIdle)
	})
	b.Run("Keep", func(b *testing.B) {
		kinds.BenchOddSuite(b, kinds.OddChecks.Keep.Callable, oddIdle)
	})
	b.Run("Embed", func(b *testing.B) {
		kinds.BenchOddSuite(b, kinds.OddChecks.Embed.Callable, oddIdle)
	})
	b.Run("Loop", func(b *testing.B) {
		kinds.BenchOddSuite(b, kinds.OddChecks.Loop.Callable, oddIdle)
	})
}

// stdSubjectEmbedded is kinds.Std under the name of the field that
// stdSubject embeds.
type stdSubjectEmbedded = kinds.Std

// stdSubject is a value of Std that only TestStdSuite builds.
type stdSubject struct{ stdSubjectEmbedded }

// TestStdSuite checks the suite of Std and the two forms of its
// application.
func TestStdSuite(t *testing.T) {
	t.Parallel()

	t.Run("StdSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, kinds.StdSuite())
	})

	t.Run("RunStdSuite", func(t *testing.T) {
		t.Parallel()
		var built *stdSubject
		kinds.RunStdSuite(
			t,
			func() *stdSubject {
				built = &stdSubject{}
				return built
			},
			conform.With(kinds.StdCheck{
				ID:        "Std/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.StdHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Std(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.StdHandle) kinds.Std { return &stdSubject{} },
			}),
			conform.Excuse(
				kinds.StdChecks.Within.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.StdChecks.Within.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.StdChecks.Pair.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.StdChecks.Meet.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStdSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *stdSubject
		kinds.RunStdSuiteFromHandle(
			t,
			func(*kinds.StdHandle) *stdSubject {
				built = &stdSubject{}
				return built
			},
			conform.With(kinds.StdCheck{
				ID:        "Std/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *kinds.StdHandle) {
					got := h.Subject()
					assert.Equal(h, got, kinds.Std(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*kinds.StdHandle) kinds.Std { return &stdSubject{} },
			}),
			conform.Excuse(
				kinds.StdChecks.Within.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.StdChecks.Within.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.StdChecks.Pair.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				kinds.StdChecks.Meet.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// stdIdle returns the planted defect of a ceiling of Std with a spend of
// nothing, a subject that answers each call at once.
func stdIdle(h *kinds.StdHandle) kinds.Std {
	return kinds.StdSuite().Calls[0].Over(h, func() {})
}

// TestStdSuiteBudgets applies the allocation ceilings of the suite of Std
// to stdIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestStdSuiteBudgets(t *testing.T) {
	kinds.RunStdBudgets(
		t,
		stdIdle,
		conform.Ceilings(kinds.StdChecks.Copy.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Close.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Log.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Buffer.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Lock.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Fail.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Within.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Hold.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Pair.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(kinds.StdChecks.Meet.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteCopy registers the fuzz target of Std.Copy on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteCopy(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Copy.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteClose registers the fuzz target of Std.Close on stdIdle,
// with every check of the suite excused.
func FuzzStdSuiteClose(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Close.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteFrom registers the fuzz target of Std.From on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteFrom(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.From.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteAt registers the fuzz target of Std.At on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteAt(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.At.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteLog registers the fuzz target of Std.Log on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteLog(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Log.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteBuffer registers the fuzz target of Std.Buffer on stdIdle,
// with every check of the suite excused.
func FuzzStdSuiteBuffer(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Buffer.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteLock registers the fuzz target of Std.Lock on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteLock(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Lock.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteFail registers the fuzz target of Std.Fail on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteFail(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Fail.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteWithin registers the fuzz target of Std.Within on stdIdle,
// with every check of the suite excused.
func FuzzStdSuiteWithin(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Within.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteHold registers the fuzz target of Std.Hold on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteHold(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Hold.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuitePair registers the fuzz target of Std.Pair on stdIdle, with
// every check of the suite excused.
func FuzzStdSuitePair(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Pair.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStdSuiteMeet registers the fuzz target of Std.Meet on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteMeet(f *testing.F) {
	kinds.FuzzStdSuite(
		f,
		kinds.StdChecks.Meet.Callable,
		stdIdle,
		conform.Excuse(
			kinds.StdChecks.Within.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Within.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Pair.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			kinds.StdChecks.Meet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStdSuite measures each callable of the suite of Std on stdIdle.
func BenchmarkStdSuite(b *testing.B) {
	b.Run("Copy", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Copy.Callable, stdIdle)
	})
	b.Run("Close", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Close.Callable, stdIdle)
	})
	b.Run("From", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.From.Callable, stdIdle)
	})
	b.Run("At", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.At.Callable, stdIdle)
	})
	b.Run("Log", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Log.Callable, stdIdle)
	})
	b.Run("Buffer", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Buffer.Callable, stdIdle)
	})
	b.Run("Lock", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Lock.Callable, stdIdle)
	})
	b.Run("Fail", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Fail.Callable, stdIdle)
	})
	b.Run("Within", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Within.Callable, stdIdle)
	})
	b.Run("Hold", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Hold.Callable, stdIdle)
	})
	b.Run("Pair", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Pair.Callable, stdIdle)
	})
	b.Run("Meet", func(b *testing.B) {
		kinds.BenchStdSuite(b, kinds.StdChecks.Meet.Callable, stdIdle)
	})
}
