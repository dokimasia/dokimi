package inputs

import (
	"context"
	"testing"

	"go.dokimi.dev/assert"
)

// TestPipelineSuiteInputs checks that each method of each implementation
// of an interface of the suite of Pipeline calls the function of its
// field.
func TestPipelineSuiteInputs(t *testing.T) {
	t.Parallel()

	t.Run("pipelineFeedInput", func(t *testing.T) {
		t.Parallel()

		t.Run("Next", func(t *testing.T) {
			t.Parallel()
			called := false
			in := pipelineFeedInput{
				onNext: func(a0 context.Context) (r0 Item, r1 bool, r2 error) {
					called = true
					return r0, r1, r2
				},
			}
			var a0 context.Context
			in.Next(a0)
			assert.True(t, called, "Next calls the function of the field onNext")
		})

		t.Run("Get", func(t *testing.T) {
			t.Parallel()
			called := false
			in := pipelineFeedInput{
				onGet: func(a0 context.Context, a1 string) (r0 Item, r1 error) {
					called = true
					return r0, r1
				},
			}
			var a0 context.Context
			var a1 string
			in.Get(a0, a1)
			assert.True(t, called, "Get calls the function of the field onGet")
		})

		t.Run("Seek", func(t *testing.T) {
			t.Parallel()
			called := false
			in := pipelineFeedInput{
				onSeek: func(a0 ...string) {
					called = true
				},
			}
			var a0 []string
			in.Seek(a0...)
			assert.True(t, called, "Seek calls the function of the field onSeek")
		})

		t.Run("Close", func(t *testing.T) {
			t.Parallel()
			called := false
			in := pipelineFeedInput{
				onClose: func() (r0 error) {
					called = true
					return r0
				},
			}
			in.Close()
			assert.True(t, called, "Close calls the function of the field onClose")
		})
	})

	t.Run("pipelineSourceInput", func(t *testing.T) {
		t.Parallel()

		t.Run("Read", func(t *testing.T) {
			t.Parallel()
			called := false
			in := pipelineSourceInput{
				onRead: func() (r0 int) {
					called = true
					return r0
				},
			}
			in.Read()
			assert.True(t, called, "Read calls the function of the field onRead")
		})
	})

	t.Run("pipelineSourceInput0", func(t *testing.T) {
		t.Parallel()

		t.Run("Read", func(t *testing.T) {
			t.Parallel()
			called := false
			in := pipelineSourceInput0{
				onRead: func() (r0 string) {
					called = true
					return r0
				},
			}
			in.Read()
			assert.True(t, called, "Read calls the function of the field onRead")
		})
	})
}
