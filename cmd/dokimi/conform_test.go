// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"
	"go.dokimi.dev/eidos/cli"
	gocli "go.stealthscale.io/cli"

	"go.dokimi.dev/dokimi/spec"
)

// The files of a corpus of Go without a suite: the fixture enum/x with its
// test, its expected outcomes, and the overlay of Go.
const (
	corpusMod      = "corpus/go/go.mod"
	corpusOverlay  = "corpus/go/overlay.yaml"
	corpusTest     = "corpus/go/enum/x/x_test.go"
	corpusExpected = "corpus/expect/enum/x.yaml"
	corpusLock     = "corpus/go/dokimi.go.lock"
)

// passing is the test of the fixture enum/x that passes and writes no record.
const passing = "package x_test\n\nimport \"testing\"\n\nfunc TestX(t *testing.T) {}\n"

// conformHeader is the comment lines that open a file of expected outcomes.
const conformHeader = "# Copyright Dokimasia B.V. 2026\n# SPDX-License-Identifier: Apache-2.0\n\n"

// noOutcomes is the file of the expected outcomes of the fixture enum/x
// without a declaration, as conform --update writes it.
const noOutcomes = conformHeader + `fixture: enum/x
level: 3
subjects: []
controls: []
declarations: {}
`

// cancelled is the file of the expected outcomes of the fixture enum/x with
// the check x.get/cancel.
const cancelled = conformHeader + `fixture: enum/x
level: 3
subjects: []
controls: []
declarations:
  x:
    checks:
      x.get/cancel:
        claim: cancel
        law: cancelled-context-is-reported
        reason: ignored-cancel
`

// corpusOf returns the tree of a repository whose corpus of Go has the
// fixture enum/x with the test and the expected outcomes.
func corpusOf(test, expected string) files.Tree {
	return files.Tree{
		corpusMod:                files.Text("module example.com/corpus\n\ngo 1.27.2\n"),
		"corpus/go/.dokimi.yaml": files.Text("version: 1\n"),
		corpusOverlay:            files.Text("language: go\nlevel: 4\ndiverge: []\n"),
		"corpus/go/enum/x/x.go":  files.Text("// Package x is the fixture of the corpus.\npackage x\n"),
		corpusTest:               files.Text(test),
		corpusExpected:           files.Text(expected),
	}
}

// recording returns the test of the fixture enum/x that writes the
// falsification record of the check x.get/cancel with the result and the
// reason code of its probe.
func recording(result spec.Result, reason spec.ReasonID) string {
	data, _ := json.Marshal(spec.SuiteRecord{
		Spec:        "1.0.0",
		Suite:       spec.RecordSuite{Language: "go", Package: "example.com/corpus/enum/x", Declaration: "X"},
		Generator:   spec.RecordTool{Name: "dokimi", Version: "0.1.0"},
		Runtime:     spec.RecordTool{Name: "go.dokimi.dev/conform", Version: "v0.9.0"},
		Seed:        "0x0000000000000001",
		Application: spec.Falsification,
		Checks: []spec.CheckRecord{{
			Identity: "X.Get/cancel", Key: "x.get/cancel", Claim: spec.Cancel, Law: spec.CancelledContextIsReported,
			Result: result, Reason: reason,
		}},
		Refusals:  []spec.RefusalRecord{},
		Unclaimed: []spec.UnclaimedRecord{},
	})
	return fmt.Sprintf("package x_test\n\nimport \"testing\"\n\nfunc TestX(t *testing.T) {\n\tt.Attr(%q, %q)\n}\n",
		"dokimi.suite.1", data)
}

// instrumentedVar is the variable of the environment that a mutation run sets
// for the test binary that it instrumented. The test binaries of a corpus
// are not instrumented, and they count no allocations under the variable.
const instrumentedVar = "DOKIMI_MUTATE_INSTRUMENTED"

// runConform runs the test binary as the binary with the arguments of
// conform and --no-color in the directory dir, with env beside the
// environment of the test without instrumentedVar, and returns its standard
// output, its standard error and its exit status.
func runConform(t *testing.T, dir string, env []string, args ...string) (stdout, stderr string, status int) {
	t.Helper()
	binary, err := os.Executable()
	assert.NoError(t, err, "the path of the test binary")
	cmd := exec.CommandContext(t.Context(), binary, append([]string{conformName, "--no-color"}, args...)...)
	cmd.Dir = dir
	environ := slices.DeleteFunc(os.Environ(), func(v string) bool { return strings.HasPrefix(v, instrumentedVar+"=") })
	cmd.Env = append(append(environ, childVar+"=1"), env...)
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	err = cmd.Run()
	if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
		return out.String(), errOut.String(), exit.ExitCode()
	}
	assert.NoError(t, err, "the binary runs")
	return out.String(), errOut.String(), 0
}

// standalone is the environment of a corpus outside every workspace of Go.
var standalone = []string{"GOWORK=off"}

func TestConform(t *testing.T) {
	t.Parallel()

	t.Run("conformRun", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the corpus of Go of the repository", func(t *testing.T) {
			t.Parallel()
			root, err := filepath.Abs(filepath.Join("..", ".."))
			assert.NoError(t, err, "the root of the repository")

			stdout, stderr, status := runConform(t, root, nil, "--lang", "go")

			assert.Equal(t, status, 0, "the corpus of Go conforms: "+stdout+stderr)
			assert.Contains(t, stdout, "go: level 4, ", "the summary names the language and its level")
		})

		t.Run("passes a corpus whose outcomes meet their expectations", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, corpusOf(passing, noOutcomes))

			stdout, stderr, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 0, "the corpus conforms: "+stderr)
			assert.Contains(t, stdout, "go: level 4, 1 fixtures, 1 conformant, 0 declared differences\n",
				"the summary counts the fixtures")
		})

		t.Run("fails a planted difference with the fixture and the two reason codes", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, corpusOf(recording(spec.WrongReason, spec.ValueBesideError), cancelled))

			stdout, _, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "a finding fails the gate")
			assert.Contains(t, stdout, "enum/x\n"+
				"  generation     FAIL  x.get/cancel: want check of cancel, got absent\n"+
				"  falsification  FAIL  x.get/cancel: want proven by ignored-cancel, got wrong-reason by value-beside-error\n"+
				"go: level 4, 1 fixtures, 0 conformant, 0 declared differences\n",
				"the gate writes each finding under its fixture")
		})

		t.Run("fails a corpus whose test fails", func(t *testing.T) {
			t.Parallel()
			failing := "package x_test\n\nimport \"testing\"\n\nfunc TestX(t *testing.T) { t.Error(\"the fixture fails\") }\n"
			dir := files.Workspace(t, corpusOf(failing, noOutcomes))

			_, stderr, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "a failed test fails the gate")
			assert.Contains(t, stderr, "the tests of the corpus failed: exit status 1", "the gate names the failure")
		})

		t.Run("fails a corpus with a record that does not validate", func(t *testing.T) {
			t.Parallel()
			invalid := "package x_test\n\nimport \"testing\"\n\nfunc TestX(t *testing.T) { t.Attr(\"dokimi.suite.1\", \"{}\") }\n"
			dir := files.Workspace(t, corpusOf(invalid, noOutcomes))

			_, stderr, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "an invalid record fails the gate")
			assert.Contains(t, stderr, "dokimi.suite.1", "the gate names the record")
		})

		t.Run("fails a corpus whose generated files are out of date", func(t *testing.T) {
			t.Parallel()
			tree := corpusOf(passing, noOutcomes)
			tree["corpus/go/enum/x/x.go"] = files.Text("// Package x is the fixture of the corpus.\npackage x\n\n" +
				"// Level is a level.\n//\n//dokimi:enum\ntype Level int\n\n// The levels.\nconst (\n\tLow Level = iota\n\tHigh\n)\n")
			dir := files.Workspace(t, tree)

			_, _, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "the generation of the corpus fails its check")
			files.Absent(t, filepath.Join(dir, "corpus", "go", "enum", "x", "level_enum.go"),
				"the check writes no generated file")
		})

		t.Run("rewrites the declarations of the expected outcomes under --update", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, corpusOf(passing, cancelled))

			stdout, stderr, status := runConform(t, dir, standalone, "--lang", "go", "--update")

			assert.Equal(t, status, 0, "the update succeeds: "+stderr)
			assert.Contains(t, stdout, "go: wrote 1 files of expected outcomes\n", "the update counts the files")
			files.HasContent(t, filepath.Join(dir, filepath.FromSlash(corpusExpected)), noOutcomes,
				"the declarations are those of the run")
		})

		t.Run("returns the error of a file of expected outcomes that it cannot write", func(t *testing.T) {
			t.Parallel()
			tree := corpusOf(passing, cancelled)
			tree[corpusExpected] = files.Text(cancelled).WithMode(0o444)
			dir := files.Workspace(t, tree)

			_, stderr, status := runConform(t, dir, standalone, "--lang", "go", "--update")

			assert.Equal(t, status, 1, "the update fails")
			assert.Contains(t, stderr, "write the expected outcomes", "the error names the step")
		})

		t.Run("fails a fixture of the corpus without expected outcomes", func(t *testing.T) {
			t.Parallel()
			tree := corpusOf(passing, noOutcomes)
			tree["corpus/go/enum/y/y.go"] = files.Text(
				"// Package y is a fixture without expected outcomes.\npackage y\n",
			)
			dir := files.Workspace(t, tree)

			stdout, _, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "a fixture without expected outcomes fails the gate")
			assert.Contains(t, stdout, "enum/y\n  generation     FAIL  want a file of expected outcomes, got none\n",
				"the gate names the fixture")
		})

		t.Run("returns the error of a directory of expected outcomes that does not exist", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, corpusOf(passing, noOutcomes))
			assert.NoError(t, os.RemoveAll(filepath.Join(dir, "corpus", "expect")), "the directory is removed")

			_, stderr, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "the gate fails")
			assert.Contains(t, stderr, "read the expected outcomes: corpus:", "the error is the error of the walk")
		})

		t.Run("returns the error of a group of fixtures that is not a directory", func(t *testing.T) {
			t.Parallel()
			tree := corpusOf(passing, noOutcomes)
			delete(tree, "corpus/go/enum/x/x.go")
			delete(tree, corpusTest)
			tree["corpus/go/enum"] = files.Text("a file in the place of the group\n")
			dir := files.Workspace(t, tree)

			_, stderr, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "the gate fails")
			assert.Contains(t, stderr, "list the fixtures of the corpus", "the error names the step")
		})

		t.Run("returns the error of a working directory that does not resolve", func(t *testing.T) {
			t.Parallel()
			var out, errOut bytes.Buffer
			run, err := gocli.NewRun(
				&cli.Global{Chdir: filepath.Join(t.TempDir(), "missing")},
				&conformOptions{Lang: conformLang},
				&gocli.Env{Stdout: &out, Stderr: &errOut},
			)
			assert.NoError(t, err, "the run has its streams")

			err = conformRun(t.Context(), run)

			usage := assert.ErrorAs[*cli.UsageError](t, err,
				"a working directory that does not resolve is a usage error")
			assert.Contains(t, usage.Error(), "find the working directory", "the usage error names its fault")
			assert.Contains(t, err.Error(), "read the working directory", "the error names the step")
		})

		t.Run("returns the error of a lock that it cannot read", func(t *testing.T) {
			t.Parallel()
			tree := corpusOf(passing, noOutcomes)
			tree[corpusLock] = files.Dir()
			dir := files.Workspace(t, tree)

			_, stderr, status := runConform(t, dir, standalone, "--lang", "go")

			assert.Equal(t, status, 1, "the gate fails")
			assert.Contains(t, stderr, "read the lock", "the error names the step")
		})

		tests := []struct {
			name       string
			giveRemove string
			want       string
		}{
			{
				name:       "returns the error of a corpus without expected outcomes",
				giveRemove: corpusExpected,
				want:       "has no file of expected outcomes",
			},
			{
				name:       "returns the error of a corpus without an overlay",
				giveRemove: corpusOverlay,
				want:       "read the overlay",
			},
			{
				name:       "returns the error of a corpus without a module",
				giveRemove: corpusMod,
				want:       "read the module of the corpus",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				dir := files.Workspace(t, corpusOf(passing, noOutcomes))
				removed := os.Remove(filepath.Join(dir, filepath.FromSlash(tt.giveRemove)))
				assert.NoError(t, removed, "the file is removed")

				_, stderr, status := runConform(t, dir, standalone, "--lang", "go")

				assert.Equal(t, status, 1, "the gate fails")
				assert.Contains(t, stderr, tt.want, "the error names the step")
			})
		}

		t.Run("refuses a language without a binding", func(t *testing.T) {
			t.Parallel()

			_, stderr, status := runConform(t, t.TempDir(), standalone, "--lang", "typescript")

			assert.Equal(t, status, 64, "a language without a binding is a usage error")
			assert.Contains(t, stderr, "has no binding of typescript", "the error names the language")
		})

		t.Run("refuses a working directory that does not exist", func(t *testing.T) {
			t.Parallel()

			_, stderr, status := runConform(t, t.TempDir(), standalone, "--lang", "go", "-C", "missing")

			assert.Equal(t, status, 64, "a working directory that does not exist is a usage error")
			assert.Contains(t, stderr, "missing", "the error names the directory")
		})
	})
}
