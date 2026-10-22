package inputs

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// pipelineSubject is a value of Pipeline that only TestPipelineSuite
// builds.
type pipelineSubject struct{ Pipeline }

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
