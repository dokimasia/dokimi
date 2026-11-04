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

	"go.dokimi.dev/dokimi/rules/internal/render"
)

// childVar is the variable of the environment that makes the test binary run
// main in place of its tests. The cases run the test binary again with it,
// so go test merges the coverage of main into the profile.
const childVar = "DOKIMI_RENDER_RULES_MAIN"

// rulesDir is the directory of the rules. It contains the YAML sources and the
// files that a render of them writes.
const rulesDir = "../../.."

func TestMain(m *testing.M) {
	if os.Getenv(childVar) != "" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// The command renders the rules in its working directory, and exits with
// status 1 and the error where the rules do not render.
func TestCmd(t *testing.T) {
	t.Parallel()

	t.Run("main", func(t *testing.T) {
		t.Parallel()

		t.Run("renders the committed files of the rules from their sources", func(t *testing.T) {
			t.Parallel()
			sources := files.Tree{}
			for _, name := range render.Files() {
				sources[name+".yaml"] = files.Text(files.Read(t, filepath.Join(rulesDir, name+".yaml")))
			}
			dir := files.Workspace(t, sources)

			out, status := run(t, dir)

			assert.Equal(t, status, 0, "the render passes: "+out)
			for _, name := range render.Files() {
				for _, file := range []string{name + ".json", filepath.Join("schemas", name+".schema.json")} {
					files.HasContent(t, filepath.Join(dir, file), files.Read(t, filepath.Join(rulesDir, file)),
						"the render writes the committed "+file)
				}
			}
		})

		t.Run("exits with status 1 where the rules do not render", func(t *testing.T) {
			t.Parallel()

			out, status := run(t, t.TempDir())

			assert.Equal(t, status, 1, "a directory without the sources of the rules fails the render")
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
