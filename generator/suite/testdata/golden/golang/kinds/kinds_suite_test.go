package kinds

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// paramsSubjectEmbedded is Params under the name of the field that
// paramsSubject embeds.
type paramsSubjectEmbedded = Params

// paramsSubject is a value of Params that only TestParamsSuite builds.
type paramsSubject struct{ paramsSubjectEmbedded }

// TestParamsSuite checks the suite of Params and the two forms of its
// application.
func TestParamsSuite(t *testing.T) {
	t.Parallel()

	t.Run("ParamsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ParamsSuite())
	})

	t.Run("RunParamsSuite", func(t *testing.T) {
		t.Parallel()
		var built *paramsSubject
		RunParamsSuite(
			t,
			func() *paramsSubject {
				built = &paramsSubject{}
				return built
			},
			conform.With(ParamsCheck{
				ID:        "Params/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ParamsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Params(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ParamsHandle) Params { return &paramsSubject{} },
			}),
		)
	})

	t.Run("RunParamsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *paramsSubject
		RunParamsSuiteFromHandle(
			t,
			func(*ParamsHandle) *paramsSubject {
				built = &paramsSubject{}
				return built
			},
			conform.With(ParamsCheck{
				ID:        "Params/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ParamsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Params(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ParamsHandle) Params { return &paramsSubject{} },
			}),
		)
	})
}

// paramsIdle returns the planted defect of a ceiling of Params with a
// spend of nothing, a subject that answers each call at once.
func paramsIdle(h *ParamsHandle) Params {
	return ParamsSuite().Calls[0].Over(h, func() {})
}

// TestParamsSuiteBudgets applies the allocation ceilings of the suite of
// Params to paramsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestParamsSuiteBudgets(t *testing.T) {
	RunParamsBudgets(
		t,
		paramsIdle,
		conform.Ceilings(ParamsChecks.Bool.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Int.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Int8.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Int16.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Int32.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Int64.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Uint.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Uint8.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Uint16.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Uint32.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Uint64.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Uintptr.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Float32.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Float64.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Complex64.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Complex128.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Array.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Chan.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Func.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Map.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Ptr.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Slice.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.String.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ParamsChecks.Struct.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			ParamsChecks.UnsafePointer.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
	)
}

// FuzzParamsSuiteBool registers the fuzz target of Params.Bool on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteBool(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Bool.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteInt registers the fuzz target of Params.Int on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Int.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteInt8 registers the fuzz target of Params.Int8 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt8(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Int8.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteInt16 registers the fuzz target of Params.Int16 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt16(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Int16.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteInt32 registers the fuzz target of Params.Int32 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt32(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Int32.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteInt64 registers the fuzz target of Params.Int64 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInt64(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Int64.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteUint registers the fuzz target of Params.Uint on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Uint.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteUint8 registers the fuzz target of Params.Uint8 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint8(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Uint8.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteUint16 registers the fuzz target of Params.Uint16 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint16(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Uint16.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteUint32 registers the fuzz target of Params.Uint32 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint32(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Uint32.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteUint64 registers the fuzz target of Params.Uint64 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUint64(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Uint64.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteUintptr registers the fuzz target of Params.Uintptr on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteUintptr(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Uintptr.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteFloat32 registers the fuzz target of Params.Float32 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteFloat32(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Float32.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteFloat64 registers the fuzz target of Params.Float64 on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteFloat64(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Float64.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteComplex64 registers the fuzz target of Params.Complex64
// on paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteComplex64(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Complex64.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteComplex128 registers the fuzz target of Params.Complex128
// on paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteComplex128(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Complex128.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteArray registers the fuzz target of Params.Array on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteArray(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Array.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteChan registers the fuzz target of Params.Chan on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteChan(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Chan.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteFunc registers the fuzz target of Params.Func on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteFunc(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Func.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteInterface registers the fuzz target of Params.Interface
// on paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteInterface(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Interface.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteMap registers the fuzz target of Params.Map on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteMap(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Map.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuitePtr registers the fuzz target of Params.Ptr on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuitePtr(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Ptr.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteSlice registers the fuzz target of Params.Slice on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteSlice(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Slice.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteString registers the fuzz target of Params.String on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteString(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.String.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteStruct registers the fuzz target of Params.Struct on
// paramsIdle, with every check of the suite excused.
func FuzzParamsSuiteStruct(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.Struct.Callable,
		paramsIdle,
	)
}

// FuzzParamsSuiteUnsafePointer registers the fuzz target of
// Params.UnsafePointer on paramsIdle, with every check of the suite
// excused.
func FuzzParamsSuiteUnsafePointer(f *testing.F) {
	FuzzParamsSuite(
		f,
		ParamsChecks.UnsafePointer.Callable,
		paramsIdle,
	)
}

// BenchmarkParamsSuite measures each callable of the suite of Params on
// paramsIdle.
func BenchmarkParamsSuite(b *testing.B) {
	b.Run("Bool", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Bool.Callable, paramsIdle)
	})
	b.Run("Int", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Int.Callable, paramsIdle)
	})
	b.Run("Int8", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Int8.Callable, paramsIdle)
	})
	b.Run("Int16", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Int16.Callable, paramsIdle)
	})
	b.Run("Int32", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Int32.Callable, paramsIdle)
	})
	b.Run("Int64", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Int64.Callable, paramsIdle)
	})
	b.Run("Uint", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Uint.Callable, paramsIdle)
	})
	b.Run("Uint8", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Uint8.Callable, paramsIdle)
	})
	b.Run("Uint16", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Uint16.Callable, paramsIdle)
	})
	b.Run("Uint32", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Uint32.Callable, paramsIdle)
	})
	b.Run("Uint64", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Uint64.Callable, paramsIdle)
	})
	b.Run("Uintptr", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Uintptr.Callable, paramsIdle)
	})
	b.Run("Float32", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Float32.Callable, paramsIdle)
	})
	b.Run("Float64", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Float64.Callable, paramsIdle)
	})
	b.Run("Complex64", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Complex64.Callable, paramsIdle)
	})
	b.Run("Complex128", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Complex128.Callable, paramsIdle)
	})
	b.Run("Array", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Array.Callable, paramsIdle)
	})
	b.Run("Chan", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Chan.Callable, paramsIdle)
	})
	b.Run("Func", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Func.Callable, paramsIdle)
	})
	b.Run("Interface", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Interface.Callable, paramsIdle)
	})
	b.Run("Map", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Map.Callable, paramsIdle)
	})
	b.Run("Ptr", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Ptr.Callable, paramsIdle)
	})
	b.Run("Slice", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Slice.Callable, paramsIdle)
	})
	b.Run("String", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.String.Callable, paramsIdle)
	})
	b.Run("Struct", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.Struct.Callable, paramsIdle)
	})
	b.Run("UnsafePointer", func(b *testing.B) {
		BenchParamsSuite(b, ParamsChecks.UnsafePointer.Callable, paramsIdle)
	})
}

// everySubjectEmbedded is Every under the name of the field that
// everySubject embeds.
type everySubjectEmbedded = Every

// everySubject is a value of Every that only TestEverySuite builds.
type everySubject struct{ everySubjectEmbedded }

// TestEverySuite checks the suite of Every and the two forms of its
// application.
func TestEverySuite(t *testing.T) {
	t.Parallel()

	t.Run("EverySuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, EverySuite())
	})

	t.Run("RunEverySuite", func(t *testing.T) {
		t.Parallel()
		var built *everySubject
		RunEverySuite(
			t,
			func() *everySubject {
				built = &everySubject{}
				return built
			},
			conform.With(EveryCheck{
				ID:        "Every/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *EveryHandle) {
					got := h.Subject()
					assert.Equal(h, got, Every(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*EveryHandle) Every { return &everySubject{} },
			}),
		)
	})

	t.Run("RunEverySuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *everySubject
		RunEverySuiteFromHandle(
			t,
			func(*EveryHandle) *everySubject {
				built = &everySubject{}
				return built
			},
			conform.With(EveryCheck{
				ID:        "Every/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *EveryHandle) {
					got := h.Subject()
					assert.Equal(h, got, Every(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*EveryHandle) Every { return &everySubject{} },
			}),
		)
	})
}

// everyIdle returns the planted defect of a ceiling of Every with a spend
// of nothing, a subject that answers each call at once.
func everyIdle(h *EveryHandle) Every {
	return EverySuite().Calls[0].Over(h, func() {})
}

// TestEverySuiteBudgets applies the allocation ceilings of the suite of
// Every to everyIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestEverySuiteBudgets(t *testing.T) {
	RunEveryBudgets(
		t,
		everyIdle,
		conform.Ceilings(EveryChecks.Fields.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(EveryChecks.Elements.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzEverySuiteFields registers the fuzz target of Every.Fields on
// everyIdle, with every check of the suite excused.
func FuzzEverySuiteFields(f *testing.F) {
	FuzzEverySuite(
		f,
		EveryChecks.Fields.Callable,
		everyIdle,
	)
}

// FuzzEverySuiteElements registers the fuzz target of Every.Elements on
// everyIdle, with every check of the suite excused.
func FuzzEverySuiteElements(f *testing.F) {
	FuzzEverySuite(
		f,
		EveryChecks.Elements.Callable,
		everyIdle,
	)
}

// FuzzEverySuiteResults registers the fuzz target of Every.Results on
// everyIdle, with every check of the suite excused.
func FuzzEverySuiteResults(f *testing.F) {
	FuzzEverySuite(
		f,
		EveryChecks.Results.Callable,
		everyIdle,
	)
}

// BenchmarkEverySuite measures each callable of the suite of Every on
// everyIdle.
func BenchmarkEverySuite(b *testing.B) {
	b.Run("Fields", func(b *testing.B) {
		BenchEverySuite(b, EveryChecks.Fields.Callable, everyIdle)
	})
	b.Run("Elements", func(b *testing.B) {
		BenchEverySuite(b, EveryChecks.Elements.Callable, everyIdle)
	})
	b.Run("Results", func(b *testing.B) {
		BenchEverySuite(b, EveryChecks.Results.Callable, everyIdle)
	})
}

// placerSubjectEmbedded is Placer under the name of the field that
// placerSubject embeds.
type placerSubjectEmbedded = Placer

// placerSubject is a value of Placer that only TestPlacerSuite builds.
type placerSubject struct{ placerSubjectEmbedded }

// TestPlacerSuite checks the suite of Placer and the two forms of its
// application.
func TestPlacerSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlacerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PlacerSuite())
	})

	t.Run("RunPlacerSuite", func(t *testing.T) {
		t.Parallel()
		var built *placerSubject
		RunPlacerSuite(
			t,
			func() *placerSubject {
				built = &placerSubject{}
				return built
			},
			conform.With(PlacerCheck{
				ID:        "Placer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlacerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Placer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlacerHandle) Placer { return &placerSubject{} },
			}),
		)
	})

	t.Run("RunPlacerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *placerSubject
		RunPlacerSuiteFromHandle(
			t,
			func(*PlacerHandle) *placerSubject {
				built = &placerSubject{}
				return built
			},
			conform.With(PlacerCheck{
				ID:        "Placer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlacerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Placer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlacerHandle) Placer { return &placerSubject{} },
			}),
		)
	})
}

// placerIdle returns the planted defect of a ceiling of Placer with a
// spend of nothing, a subject that answers each call at once.
func placerIdle(h *PlacerHandle) Placer {
	return PlacerSuite().Calls[0].Over(h, func() {})
}

// TestPlacerSuiteBudgets applies the allocation ceilings of the suite of
// Placer to placerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestPlacerSuiteBudgets(t *testing.T) {
	RunPlacerBudgets(
		t,
		placerIdle,
	)
}

// FuzzPlacerSuitePlace registers the fuzz target of Placer.Place on
// placerIdle, with every check of the suite excused.
func FuzzPlacerSuitePlace(f *testing.F) {
	FuzzPlacerSuite(
		f,
		PlacerChecks.Place.Callable,
		placerIdle,
	)
}

// BenchmarkPlacerSuite measures each callable of the suite of Placer on
// placerIdle.
func BenchmarkPlacerSuite(b *testing.B) {
	b.Run("Place", func(b *testing.B) {
		BenchPlacerSuite(b, PlacerChecks.Place.Callable, placerIdle)
	})
}

// routerSubjectEmbedded is Router under the name of the field that
// routerSubject embeds.
type routerSubjectEmbedded = Router

// routerSubject is a value of Router that only TestRouterSuite builds.
type routerSubject struct{ routerSubjectEmbedded }

// TestRouterSuite checks the suite of Router and the two forms of its
// application.
func TestRouterSuite(t *testing.T) {
	t.Parallel()

	t.Run("RouterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, RouterSuite())
	})

	t.Run("RunRouterSuite", func(t *testing.T) {
		t.Parallel()
		var built *routerSubject
		RunRouterSuite(
			t,
			func() *routerSubject {
				built = &routerSubject{}
				return built
			},
			conform.With(RouterCheck{
				ID:        "Router/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *RouterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Router(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*RouterHandle) Router { return &routerSubject{} },
			}),
		)
	})

	t.Run("RunRouterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *routerSubject
		RunRouterSuiteFromHandle(
			t,
			func(*RouterHandle) *routerSubject {
				built = &routerSubject{}
				return built
			},
			conform.With(RouterCheck{
				ID:        "Router/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *RouterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Router(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*RouterHandle) Router { return &routerSubject{} },
			}),
		)
	})
}

// routerIdle returns the planted defect of a ceiling of Router with a
// spend of nothing, a subject that answers each call at once.
func routerIdle(h *RouterHandle) Router {
	return RouterSuite().Calls[0].Over(h, func() {})
}

// TestRouterSuiteBudgets applies the allocation ceilings of the suite of
// Router to routerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRouterSuiteBudgets(t *testing.T) {
	RunRouterBudgets(
		t,
		routerIdle,
		conform.Ceilings(RouterChecks.Route.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(RouterChecks.Write.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(RouterChecks.Pass.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzRouterSuiteRoute registers the fuzz target of Router.Route on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteRoute(f *testing.F) {
	FuzzRouterSuite(
		f,
		RouterChecks.Route.Callable,
		routerIdle,
	)
}

// FuzzRouterSuiteRead registers the fuzz target of Router.Read on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteRead(f *testing.F) {
	FuzzRouterSuite(
		f,
		RouterChecks.Read.Callable,
		routerIdle,
	)
}

// FuzzRouterSuiteWrite registers the fuzz target of Router.Write on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteWrite(f *testing.F) {
	FuzzRouterSuite(
		f,
		RouterChecks.Write.Callable,
		routerIdle,
	)
}

// FuzzRouterSuitePass registers the fuzz target of Router.Pass on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuitePass(f *testing.F) {
	FuzzRouterSuite(
		f,
		RouterChecks.Pass.Callable,
		routerIdle,
	)
}

// FuzzRouterSuiteWatch registers the fuzz target of Router.Watch on
// routerIdle, with every check of the suite excused.
func FuzzRouterSuiteWatch(f *testing.F) {
	FuzzRouterSuite(
		f,
		RouterChecks.Watch.Callable,
		routerIdle,
	)
}

// BenchmarkRouterSuite measures each callable of the suite of Router on
// routerIdle.
func BenchmarkRouterSuite(b *testing.B) {
	b.Run("Route", func(b *testing.B) {
		BenchRouterSuite(b, RouterChecks.Route.Callable, routerIdle)
	})
	b.Run("Read", func(b *testing.B) {
		BenchRouterSuite(b, RouterChecks.Read.Callable, routerIdle)
	})
	b.Run("Write", func(b *testing.B) {
		BenchRouterSuite(b, RouterChecks.Write.Callable, routerIdle)
	})
	b.Run("Pass", func(b *testing.B) {
		BenchRouterSuite(b, RouterChecks.Pass.Callable, routerIdle)
	})
	b.Run("Watch", func(b *testing.B) {
		BenchRouterSuite(b, RouterChecks.Watch.Callable, routerIdle)
	})
}

// drawerSubjectEmbedded is Drawer under the name of the field that
// drawerSubject embeds.
type drawerSubjectEmbedded = Drawer

// drawerSubject is a value of Drawer that only TestDrawerSuite builds.
type drawerSubject struct{ drawerSubjectEmbedded }

// TestDrawerSuite checks the suite of Drawer and the two forms of its
// application.
func TestDrawerSuite(t *testing.T) {
	t.Parallel()

	t.Run("DrawerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, DrawerSuite())
	})

	t.Run("RunDrawerSuite", func(t *testing.T) {
		t.Parallel()
		var built *drawerSubject
		RunDrawerSuite(
			t,
			func() *drawerSubject {
				built = &drawerSubject{}
				return built
			},
			conform.With(DrawerCheck{
				ID:        "Drawer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *DrawerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Drawer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*DrawerHandle) Drawer { return &drawerSubject{} },
			}),
		)
	})

	t.Run("RunDrawerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *drawerSubject
		RunDrawerSuiteFromHandle(
			t,
			func(*DrawerHandle) *drawerSubject {
				built = &drawerSubject{}
				return built
			},
			conform.With(DrawerCheck{
				ID:        "Drawer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *DrawerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Drawer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*DrawerHandle) Drawer { return &drawerSubject{} },
			}),
		)
	})
}

// drawerIdle returns the planted defect of a ceiling of Drawer with a
// spend of nothing, a subject that answers each call at once.
func drawerIdle(h *DrawerHandle) Drawer {
	return DrawerSuite().Calls[0].Over(h, func() {})
}

// TestDrawerSuiteBudgets applies the allocation ceilings of the suite of
// Drawer to drawerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestDrawerSuiteBudgets(t *testing.T) {
	RunDrawerBudgets(
		t,
		drawerIdle,
		conform.Ceilings(DrawerChecks.Draw.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(DrawerChecks.Count.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(DrawerChecks.Name.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(DrawerChecks.Seal.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(DrawerChecks.Test.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzDrawerSuiteDraw registers the fuzz target of Drawer.Draw on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteDraw(f *testing.F) {
	FuzzDrawerSuite(
		f,
		DrawerChecks.Draw.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteCount registers the fuzz target of Drawer.Count on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteCount(f *testing.F) {
	FuzzDrawerSuite(
		f,
		DrawerChecks.Count.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteName registers the fuzz target of Drawer.Name on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteName(f *testing.F) {
	FuzzDrawerSuite(
		f,
		DrawerChecks.Name.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuitePick registers the fuzz target of Drawer.Pick on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuitePick(f *testing.F) {
	FuzzDrawerSuite(
		f,
		DrawerChecks.Pick.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteSeal registers the fuzz target of Drawer.Seal on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteSeal(f *testing.F) {
	FuzzDrawerSuite(
		f,
		DrawerChecks.Seal.Callable,
		drawerIdle,
	)
}

// FuzzDrawerSuiteTest registers the fuzz target of Drawer.Test on
// drawerIdle, with every check of the suite excused.
func FuzzDrawerSuiteTest(f *testing.F) {
	FuzzDrawerSuite(
		f,
		DrawerChecks.Test.Callable,
		drawerIdle,
	)
}

// BenchmarkDrawerSuite measures each callable of the suite of Drawer on
// drawerIdle.
func BenchmarkDrawerSuite(b *testing.B) {
	b.Run("Draw", func(b *testing.B) {
		BenchDrawerSuite(b, DrawerChecks.Draw.Callable, drawerIdle)
	})
	b.Run("Count", func(b *testing.B) {
		BenchDrawerSuite(b, DrawerChecks.Count.Callable, drawerIdle)
	})
	b.Run("Name", func(b *testing.B) {
		BenchDrawerSuite(b, DrawerChecks.Name.Callable, drawerIdle)
	})
	b.Run("Pick", func(b *testing.B) {
		BenchDrawerSuite(b, DrawerChecks.Pick.Callable, drawerIdle)
	})
	b.Run("Seal", func(b *testing.B) {
		BenchDrawerSuite(b, DrawerChecks.Seal.Callable, drawerIdle)
	})
	b.Run("Test", func(b *testing.B) {
		BenchDrawerSuite(b, DrawerChecks.Test.Callable, drawerIdle)
	})
}

// oddSubjectEmbedded is Odd under the name of the field that oddSubject
// embeds.
type oddSubjectEmbedded = Odd

// oddSubject is a value of Odd that only TestOddSuite builds.
type oddSubject struct{ oddSubjectEmbedded }

// TestOddSuite checks the suite of Odd and the two forms of its
// application.
func TestOddSuite(t *testing.T) {
	t.Parallel()

	t.Run("OddSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, OddSuite())
	})

	t.Run("RunOddSuite", func(t *testing.T) {
		t.Parallel()
		var built *oddSubject
		RunOddSuite(
			t,
			func() *oddSubject {
				built = &oddSubject{}
				return built
			},
			conform.With(OddCheck{
				ID:        "Odd/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *OddHandle) {
					got := h.Subject()
					assert.Equal(h, got, Odd(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*OddHandle) Odd { return &oddSubject{} },
			}),
		)
	})

	t.Run("RunOddSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *oddSubject
		RunOddSuiteFromHandle(
			t,
			func(*OddHandle) *oddSubject {
				built = &oddSubject{}
				return built
			},
			conform.With(OddCheck{
				ID:        "Odd/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *OddHandle) {
					got := h.Subject()
					assert.Equal(h, got, Odd(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*OddHandle) Odd { return &oddSubject{} },
			}),
		)
	})
}

// oddIdle returns the planted defect of a ceiling of Odd with a spend of
// nothing, a subject that answers each call at once.
func oddIdle(h *OddHandle) Odd {
	return OddSuite().Calls[0].Over(h, func() {})
}

// TestOddSuiteBudgets applies the allocation ceilings of the suite of Odd
// to oddIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestOddSuiteBudgets(t *testing.T) {
	RunOddBudgets(
		t,
		oddIdle,
		conform.Ceilings(OddChecks.Loop.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzOddSuiteTag registers the fuzz target of Odd.Tag on oddIdle, with
// every check of the suite excused.
func FuzzOddSuiteTag(f *testing.F) {
	FuzzOddSuite(
		f,
		OddChecks.Tag.Callable,
		oddIdle,
	)
}

// FuzzOddSuiteKeep registers the fuzz target of Odd.Keep on oddIdle, with
// every check of the suite excused.
func FuzzOddSuiteKeep(f *testing.F) {
	FuzzOddSuite(
		f,
		OddChecks.Keep.Callable,
		oddIdle,
	)
}

// FuzzOddSuiteEmbed registers the fuzz target of Odd.Embed on oddIdle,
// with every check of the suite excused.
func FuzzOddSuiteEmbed(f *testing.F) {
	FuzzOddSuite(
		f,
		OddChecks.Embed.Callable,
		oddIdle,
	)
}

// FuzzOddSuiteLoop registers the fuzz target of Odd.Loop on oddIdle, with
// every check of the suite excused.
func FuzzOddSuiteLoop(f *testing.F) {
	FuzzOddSuite(
		f,
		OddChecks.Loop.Callable,
		oddIdle,
	)
}

// BenchmarkOddSuite measures each callable of the suite of Odd on oddIdle.
func BenchmarkOddSuite(b *testing.B) {
	b.Run("Tag", func(b *testing.B) {
		BenchOddSuite(b, OddChecks.Tag.Callable, oddIdle)
	})
	b.Run("Keep", func(b *testing.B) {
		BenchOddSuite(b, OddChecks.Keep.Callable, oddIdle)
	})
	b.Run("Embed", func(b *testing.B) {
		BenchOddSuite(b, OddChecks.Embed.Callable, oddIdle)
	})
	b.Run("Loop", func(b *testing.B) {
		BenchOddSuite(b, OddChecks.Loop.Callable, oddIdle)
	})
}

// stdSubjectEmbedded is Std under the name of the field that stdSubject
// embeds.
type stdSubjectEmbedded = Std

// stdSubject is a value of Std that only TestStdSuite builds.
type stdSubject struct{ stdSubjectEmbedded }

// TestStdSuite checks the suite of Std and the two forms of its
// application.
func TestStdSuite(t *testing.T) {
	t.Parallel()

	t.Run("StdSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, StdSuite())
	})

	t.Run("RunStdSuite", func(t *testing.T) {
		t.Parallel()
		var built *stdSubject
		RunStdSuite(
			t,
			func() *stdSubject {
				built = &stdSubject{}
				return built
			},
			conform.With(StdCheck{
				ID:        "Std/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *StdHandle) {
					got := h.Subject()
					assert.Equal(h, got, Std(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*StdHandle) Std { return &stdSubject{} },
			}),
		)
	})

	t.Run("RunStdSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *stdSubject
		RunStdSuiteFromHandle(
			t,
			func(*StdHandle) *stdSubject {
				built = &stdSubject{}
				return built
			},
			conform.With(StdCheck{
				ID:        "Std/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *StdHandle) {
					got := h.Subject()
					assert.Equal(h, got, Std(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*StdHandle) Std { return &stdSubject{} },
			}),
		)
	})
}

// stdIdle returns the planted defect of a ceiling of Std with a spend of
// nothing, a subject that answers each call at once.
func stdIdle(h *StdHandle) Std {
	return StdSuite().Calls[0].Over(h, func() {})
}

// TestStdSuiteBudgets applies the allocation ceilings of the suite of Std
// to stdIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestStdSuiteBudgets(t *testing.T) {
	RunStdBudgets(
		t,
		stdIdle,
		conform.Ceilings(StdChecks.Copy.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Close.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Log.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Buffer.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Lock.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Fail.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Within.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Hold.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Pair.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StdChecks.Meet.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzStdSuiteCopy registers the fuzz target of Std.Copy on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteCopy(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Copy.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteClose registers the fuzz target of Std.Close on stdIdle,
// with every check of the suite excused.
func FuzzStdSuiteClose(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Close.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteFrom registers the fuzz target of Std.From on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteFrom(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.From.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteAt registers the fuzz target of Std.At on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteAt(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.At.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteLog registers the fuzz target of Std.Log on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteLog(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Log.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteBuffer registers the fuzz target of Std.Buffer on stdIdle,
// with every check of the suite excused.
func FuzzStdSuiteBuffer(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Buffer.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteLock registers the fuzz target of Std.Lock on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteLock(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Lock.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteFail registers the fuzz target of Std.Fail on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteFail(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Fail.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteWithin registers the fuzz target of Std.Within on stdIdle,
// with every check of the suite excused.
func FuzzStdSuiteWithin(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Within.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteHold registers the fuzz target of Std.Hold on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteHold(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Hold.Callable,
		stdIdle,
	)
}

// FuzzStdSuitePair registers the fuzz target of Std.Pair on stdIdle, with
// every check of the suite excused.
func FuzzStdSuitePair(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Pair.Callable,
		stdIdle,
	)
}

// FuzzStdSuiteMeet registers the fuzz target of Std.Meet on stdIdle, with
// every check of the suite excused.
func FuzzStdSuiteMeet(f *testing.F) {
	FuzzStdSuite(
		f,
		StdChecks.Meet.Callable,
		stdIdle,
	)
}

// BenchmarkStdSuite measures each callable of the suite of Std on stdIdle.
func BenchmarkStdSuite(b *testing.B) {
	b.Run("Copy", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Copy.Callable, stdIdle)
	})
	b.Run("Close", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Close.Callable, stdIdle)
	})
	b.Run("From", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.From.Callable, stdIdle)
	})
	b.Run("At", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.At.Callable, stdIdle)
	})
	b.Run("Log", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Log.Callable, stdIdle)
	})
	b.Run("Buffer", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Buffer.Callable, stdIdle)
	})
	b.Run("Lock", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Lock.Callable, stdIdle)
	})
	b.Run("Fail", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Fail.Callable, stdIdle)
	})
	b.Run("Within", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Within.Callable, stdIdle)
	})
	b.Run("Hold", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Hold.Callable, stdIdle)
	})
	b.Run("Pair", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Pair.Callable, stdIdle)
	})
	b.Run("Meet", func(b *testing.B) {
		BenchStdSuite(b, StdChecks.Meet.Callable, stdIdle)
	})
}
