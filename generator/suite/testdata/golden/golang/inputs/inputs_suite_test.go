package inputs

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// pipelineSubjectEmbedded is Pipeline under the name of the field that
// pipelineSubject embeds.
type pipelineSubjectEmbedded = Pipeline

// pipelineSubject is a value of Pipeline that only TestPipelineSuite
// builds.
type pipelineSubject struct{ pipelineSubjectEmbedded }

// TestPipelineSuite checks the suite of Pipeline and the two forms of its
// application.
func TestPipelineSuite(t *testing.T) {
	t.Parallel()

	t.Run("PipelineSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PipelineSuite())
	})

	t.Run("RunPipelineSuite", func(t *testing.T) {
		t.Parallel()
		var built *pipelineSubject
		RunPipelineSuite(
			t,
			func() *pipelineSubject {
				built = &pipelineSubject{}
				return built
			},
			conform.With(PipelineCheck{
				ID:        "Pipeline/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PipelineHandle) {
					got := h.Subject()
					assert.Equal(h, got, Pipeline(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PipelineHandle) Pipeline { return &pipelineSubject{} },
			}),
		)
	})

	t.Run("RunPipelineSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *pipelineSubject
		RunPipelineSuiteFromHandle(
			t,
			func(*PipelineHandle) *pipelineSubject {
				built = &pipelineSubject{}
				return built
			},
			conform.With(PipelineCheck{
				ID:        "Pipeline/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PipelineHandle) {
					got := h.Subject()
					assert.Equal(h, got, Pipeline(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PipelineHandle) Pipeline { return &pipelineSubject{} },
			}),
		)
	})
}

// pipelineIdle returns the planted defect of a ceiling of Pipeline with a
// spend of nothing, a subject that answers each call at once.
func pipelineIdle(h *PipelineHandle) Pipeline {
	return PipelineSuite().Calls[0].Over(h, func() {})
}

// TestPipelineSuiteBudgets applies the allocation ceilings of the suite of
// Pipeline to pipelineIdle, and a ceiling of 0 to each callable without a
// budget whose calls build no input of their own.
func TestPipelineSuiteBudgets(t *testing.T) {
	RunPipelineBudgets(
		t,
		pipelineIdle,
		conform.Ceilings(PipelineChecks.Each.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Handle.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Do.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Notify.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Apply.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Spin.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Copy.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Configure.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Subscribe.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PipelineChecks.Label.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzPipelineSuitePump registers the fuzz target of Pipeline.Pump on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuitePump(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Pump.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteMerge registers the fuzz target of Pipeline.Merge on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteMerge(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Merge.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteEach registers the fuzz target of Pipeline.Each on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteEach(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Each.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteHandle registers the fuzz target of Pipeline.Handle on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteHandle(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Handle.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteDo registers the fuzz target of Pipeline.Do on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteDo(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Do.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteNotify registers the fuzz target of Pipeline.Notify on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteNotify(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Notify.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteApply registers the fuzz target of Pipeline.Apply on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteApply(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Apply.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteFold registers the fuzz target of Pipeline.Fold on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteFold(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Fold.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuitePut registers the fuzz target of Pipeline.Put on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuitePut(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Put.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteCount registers the fuzz target of Pipeline.Count on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteCount(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Count.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteSpin registers the fuzz target of Pipeline.Spin on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteSpin(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Spin.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteCopy registers the fuzz target of Pipeline.Copy on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteCopy(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Copy.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteConfigure registers the fuzz target of
// Pipeline.Configure on pipelineIdle, with every check of the suite
// excused.
func FuzzPipelineSuiteConfigure(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Configure.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteSubscribe registers the fuzz target of
// Pipeline.Subscribe on pipelineIdle, with every check of the suite
// excused.
func FuzzPipelineSuiteSubscribe(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Subscribe.Callable,
		pipelineIdle,
	)
}

// FuzzPipelineSuiteLabel registers the fuzz target of Pipeline.Label on
// pipelineIdle, with every check of the suite excused.
func FuzzPipelineSuiteLabel(f *testing.F) {
	FuzzPipelineSuite(
		f,
		PipelineChecks.Label.Callable,
		pipelineIdle,
	)
}

// BenchmarkPipelineSuite measures each callable of the suite of Pipeline
// on pipelineIdle.
func BenchmarkPipelineSuite(b *testing.B) {
	b.Run("Pump", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Pump.Callable, pipelineIdle)
	})
	b.Run("Merge", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Merge.Callable, pipelineIdle)
	})
	b.Run("Each", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Each.Callable, pipelineIdle)
	})
	b.Run("Handle", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Handle.Callable, pipelineIdle)
	})
	b.Run("Do", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Do.Callable, pipelineIdle)
	})
	b.Run("Notify", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Notify.Callable, pipelineIdle)
	})
	b.Run("Apply", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Apply.Callable, pipelineIdle)
	})
	b.Run("Fold", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Fold.Callable, pipelineIdle)
	})
	b.Run("Put", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Put.Callable, pipelineIdle)
	})
	b.Run("Count", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Count.Callable, pipelineIdle)
	})
	b.Run("Spin", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Spin.Callable, pipelineIdle)
	})
	b.Run("Copy", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Copy.Callable, pipelineIdle)
	})
	b.Run("Configure", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Configure.Callable, pipelineIdle)
	})
	b.Run("Subscribe", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Subscribe.Callable, pipelineIdle)
	})
	b.Run("Label", func(b *testing.B) {
		BenchPipelineSuite(b, PipelineChecks.Label.Callable, pipelineIdle)
	})
}
