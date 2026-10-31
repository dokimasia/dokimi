package inputs_test

import (
	"testing"

	"example.com/shop/inputs"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// pipelineSubjectEmbedded is inputs.Pipeline under the name of the field
// that pipelineSubject embeds.
type pipelineSubjectEmbedded = inputs.Pipeline

// pipelineSubject is a value of Pipeline that only TestPipelineSuite
// builds.
type pipelineSubject struct{ pipelineSubjectEmbedded }

// TestPipelineSuite checks the suite of Pipeline and the two forms of its
// application.
func TestPipelineSuite(t *testing.T) {
	t.Parallel()

	t.Run("PipelineSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, inputs.PipelineSuite())
	})

	t.Run("RunPipelineSuite", func(t *testing.T) {
		t.Parallel()
		var built *pipelineSubject
		inputs.RunPipelineSuite(
			t,
			func() *pipelineSubject {
				built = &pipelineSubject{}
				return built
			},
			conform.With(inputs.PipelineCheck{
				ID:        "Pipeline/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *inputs.PipelineHandle) {
					got := h.Subject()
					assert.Equal(h, got, inputs.Pipeline(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*inputs.PipelineHandle) inputs.Pipeline { return &pipelineSubject{} },
			}),
			conform.Excuse(
				inputs.PipelineChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Do.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Notify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Apply.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Label.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPipelineSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *pipelineSubject
		inputs.RunPipelineSuiteFromHandle(
			t,
			func(*inputs.PipelineHandle) *pipelineSubject {
				built = &pipelineSubject{}
				return built
			},
			conform.With(inputs.PipelineCheck{
				ID:        "Pipeline/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *inputs.PipelineHandle) {
					got := h.Subject()
					assert.Equal(h, got, inputs.Pipeline(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*inputs.PipelineHandle) inputs.Pipeline { return &pipelineSubject{} },
			}),
			conform.Excuse(
				inputs.PipelineChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Do.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Notify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Apply.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				inputs.PipelineChecks.Label.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// pipelineIdle returns the planted defect of a ceiling of Pipeline with a
// spend of nothing, a subject that answers each call at once.
func pipelineIdle(h *inputs.PipelineHandle) inputs.Pipeline {
	return inputs.PipelineSuite().Calls[0].Over(h, func() {})
}

// TestPipelineSuiteBudgets applies the allocation ceilings of the suite of
// Pipeline to pipelineIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestPipelineSuiteBudgets(t *testing.T) {
	inputs.RunPipelineBudgets(
		t,
		pipelineIdle,
		conform.Ceilings(
			inputs.PipelineChecks.Each.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Handle.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(inputs.PipelineChecks.Do.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			inputs.PipelineChecks.Notify.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Apply.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Spin.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Copy.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Configure.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Subscribe.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			inputs.PipelineChecks.Label.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuitePump registers the fuzz target of Pipeline.Pump on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuitePump(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Pump.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteMerge registers the fuzz target of Pipeline.Merge on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteMerge(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Merge.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteEach registers the fuzz target of Pipeline.Each on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteEach(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Each.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteHandle registers the fuzz target of Pipeline.Handle on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteHandle(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Handle.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteDo registers the fuzz target of Pipeline.Do on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteDo(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Do.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteNotify registers the fuzz target of Pipeline.Notify on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteNotify(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Notify.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteApply registers the fuzz target of Pipeline.Apply on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteApply(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Apply.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteFold registers the fuzz target of Pipeline.Fold on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteFold(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Fold.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuitePut registers the fuzz target of Pipeline.Put on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuitePut(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Put.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteCount registers the fuzz target of Pipeline.Count on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteCount(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Count.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteSpin registers the fuzz target of Pipeline.Spin on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteSpin(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Spin.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteCopy registers the fuzz target of Pipeline.Copy on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteCopy(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Copy.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteConfigure registers the fuzz target of
// Pipeline.Configure on pipelineIdle, with every check of the suite
// excused.
func FuzzPipelineSuiteConfigure(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Configure.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteSubscribe registers the fuzz target of
// Pipeline.Subscribe on pipelineIdle, with every check of the suite
// excused.
func FuzzPipelineSuiteSubscribe(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Subscribe.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPipelineSuiteLabel registers the fuzz target of Pipeline.Label on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteLabel(f *testing.F) {
	inputs.FuzzPipelineSuite(
		f,
		inputs.PipelineChecks.Label.Callable,
		pipelineIdle,
		conform.Excuse(
			inputs.PipelineChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Notify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Apply.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			inputs.PipelineChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPipelineSuite measures each callable of the suite of Pipeline
// on pipelineIdle.
func BenchmarkPipelineSuite(b *testing.B) {
	b.Run("Pump", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Pump.Callable, pipelineIdle)
	})
	b.Run("Merge", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Merge.Callable, pipelineIdle)
	})
	b.Run("Each", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Each.Callable, pipelineIdle)
	})
	b.Run("Handle", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Handle.Callable, pipelineIdle)
	})
	b.Run("Do", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Do.Callable, pipelineIdle)
	})
	b.Run("Notify", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Notify.Callable, pipelineIdle)
	})
	b.Run("Apply", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Apply.Callable, pipelineIdle)
	})
	b.Run("Fold", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Fold.Callable, pipelineIdle)
	})
	b.Run("Put", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Put.Callable, pipelineIdle)
	})
	b.Run("Count", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Count.Callable, pipelineIdle)
	})
	b.Run("Spin", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Spin.Callable, pipelineIdle)
	})
	b.Run("Copy", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Copy.Callable, pipelineIdle)
	})
	b.Run("Configure", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Configure.Callable, pipelineIdle)
	})
	b.Run("Subscribe", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Subscribe.Callable, pipelineIdle)
	})
	b.Run("Label", func(b *testing.B) {
		inputs.BenchPipelineSuite(b, inputs.PipelineChecks.Label.Callable, pipelineIdle)
	})
}
