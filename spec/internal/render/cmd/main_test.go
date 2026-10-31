// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/files"

	"go.dokimi.dev/dokimi/spec/internal/render"
)

// childVar is the variable of the environment that makes the test binary run
// main in place of its tests. The cases run the test binary again with it,
// so go test merges the coverage of main into the profile.
const childVar = "DOKIMI_RENDER_MAIN"

// specDir is the directory of the spec. It contains the YAML sources and the
// files that a render of them writes.
const specDir = "../../.."

func TestMain(m *testing.M) {
	if os.Getenv(childVar) != "" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// The command renders the spec in its working directory, and exits with
// status 1 and the error where the spec does not render.
func TestCmd(t *testing.T) {
	t.Parallel()

	t.Run("main", func(t *testing.T) {
		t.Parallel()

		t.Run("renders the committed files of the spec from its sources", func(t *testing.T) {
			t.Parallel()
			sources := files.Tree{}
			for _, name := range append(render.Files(), "keys") {
				sources[name+".yaml"] = files.Text(files.Read(t, filepath.Join(specDir, name+".yaml")))
			}
			dir := files.Workspace(t, sources)

			out, status := run(t, dir)

			assert.Equal(t, status, 0, "the render passes: "+out)
			for _, name := range append(render.Files(), "keys", "naming") {
				want := files.Read(t, filepath.Join(specDir, name+".json"))
				files.HasContent(t, filepath.Join(dir, name+".json"), want,
					"the render writes the committed "+name+".json")
			}
			want := files.Read(t, filepath.Join(specDir, "names.gen.go"))
			files.HasContent(t, filepath.Join(dir, "names.gen.go"), want,
				"the render writes the committed names.gen.go")
		})

		t.Run("exits with status 1 where the spec does not render", func(t *testing.T) {
			t.Parallel()

			out, status := run(t, t.TempDir())

			assert.Equal(t, status, 1, "a directory without the sources of the spec fails the render")
			expect.HasPrefix(t, out, "render: ", "the command writes the error of the render")
		})
	})
}

// run runs main in the directory dir through the test binary, and returns
// its output and its exit status.
func run(t *testing.T, dir string) (string, int) {
	t.Helper()

	binary, err := os.Executable()
	assert.NoError(t, err, "the path of the test binary")
	cmd := exec.CommandContext(t.Context(), binary)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), childVar+"=1")
	out, err := cmd.CombinedOutput()
	if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
		return string(out), exit.ExitCode()
	}
	assert.NoError(t, err, "the test binary runs")
	return string(out), 0
}
