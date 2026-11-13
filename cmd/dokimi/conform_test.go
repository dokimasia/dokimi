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
	"strconv"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"
	"go.dokimi.dev/eidos/cli"
	"go.dokimi.dev/spec"
	gocli "go.stealthscale.io/cli"
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

// The test binary stands in for Vitest where its first argument is vitestArg,
// and a run of conform calls it in place of Vitest where the environment has
// vitestReportVar. The stand-in writes the file of vitestReportVar to the path
// of its --outputFile, writes its filters to the file of vitestArgsVar, one
// on each line, and exits with the status of vitestStatusVar.
const (
	vitestArg       = "-dokimi.test.vitest"
	vitestReportVar = "DOKIMI_TEST_VITEST_REPORT"
	vitestArgsVar   = "DOKIMI_TEST_VITEST_ARGS"
	vitestStatusVar = "DOKIMI_TEST_VITEST_STATUS"
)

// standInVitest runs the stand-in of Vitest with the arguments that follow
// vitestArg and returns its exit status.
func standInVitest(args []string) int {
	report, err := os.ReadFile(os.Getenv(vitestReportVar))
	var filters []string
	for _, a := range args {
		out, isOutput := strings.CutPrefix(a, "--outputFile=")
		switch {
		case isOutput && err == nil:
			err = os.WriteFile(out, report, 0o600)
		case !isOutput && a != "run" && a != "--reporter=json":
			filters = append(filters, a)
		}
	}
	if f := os.Getenv(vitestArgsVar); f != "" && err == nil {
		err = os.WriteFile(f, []byte(strings.Join(filters, "\n")), 0o600)
	}
	status, _ := strconv.Atoi(os.Getenv(vitestStatusVar))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return status
}

// The files of a corpus of TypeScript: the overlay of TypeScript, the report
// of the run of the stand-in of Vitest, and the file of its filters.
const (
	tsOverlay = "corpus/typescript/overlay.yaml"
	tsReport  = "report.json"
	tsFilters = "filters.txt"
)

// typescriptCorpusOf returns the tree of a repository whose corpus of
// TypeScript has the fixture enum/x with the expected outcomes, and the report
// that the stand-in of Vitest writes.
func typescriptCorpusOf(report, expected string) files.Tree {
	return files.Tree{
		"corpus/typescript/package.json": files.Text(
			"{\n  \"name\": \"corpus\",\n  \"private\": true,\n  \"type\": \"module\"\n}\n"),
		"corpus/typescript/tsconfig.json": files.Text("{}\n"),
		"corpus/typescript/.dokimi.yaml":  files.Text("version: 1\n"),
		tsOverlay:                         files.Text("language: typescript\nlevel: 4\ndiverge: []\n"),
		"corpus/typescript/enum/x/x.ts": files.Text(
			"/** Returns one. */\nexport function one(): number {\n  return 1;\n}\n"),
		corpusExpected: files.Text(expected),
		tsReport:       files.Text(report),
	}
}

// vitestReportOf returns the JSON report of Vitest with one test of the file
// of the fixture enum/x, whose meta has the records under the key
// dokimi.suite, and with the status of the test and its messages.
func vitestReportOf(status string, messages []string, records ...spec.SuiteRecord) string {
	data, _ := json.Marshal(map[string]any{"testResults": []any{map[string]any{
		"name": "/corpus/typescript/enum/x/x.suite.test.ts", "status": status, "message": "",
		"assertionResults": []any{map[string]any{
			"fullName": "xSuite suite records", "status": status, "failureMessages": messages,
			"meta": map[string]any{"dokimi.suite": records},
		}},
	}}})
	return string(data)
}

// typescriptRecording returns the falsification record of the TypeScript
// suite of X in the module enum/x/x with the check x.get/cancel, with the
// result and the reason code of its probe.
func typescriptRecording(result spec.Result, reason spec.ReasonID) spec.SuiteRecord {
	return spec.SuiteRecord{
		Spec:        "1.0.0",
		Suite:       spec.RecordSuite{Language: "typescript", Package: "enum/x/x", Declaration: "X"},
		Generator:   spec.RecordTool{Name: "dokimi", Version: "0.1.0"},
		Runtime:     spec.RecordTool{Name: "@dokimi/runtime", Version: "0.1.0"},
		Seed:        "0x0000000000000001",
		Application: spec.Falsification,
		Checks: []spec.CheckRecord{{
			Identity: "X.get/cancel", Key: "x.get/cancel", Claim: spec.Cancel, Law: spec.CancelledContextIsReported,
			Result: result, Reason: reason,
		}},
		Refusals:  []spec.RefusalRecord{},
		Unclaimed: []spec.UnclaimedRecord{},
	}
}

// tempDirEnv returns the variables of the environment that set the directory
// of the temporary files of a child process to dir. os.TempDir reads TMPDIR on
// Unix and TMP on Windows.
func tempDirEnv(dir string) []string {
	return []string{"TMPDIR=" + dir, "TMP=" + dir}
}

// vitestEnv returns the environment of a run of conform in dir that calls the
// stand-in of Vitest with the report of the tree and the status.
func vitestEnv(dir string, status int) []string {
	return []string{
		vitestReportVar + "=" + filepath.Join(dir, tsReport),
		vitestArgsVar + "=" + filepath.Join(dir, tsFilters),
		vitestStatusVar + "=" + strconv.Itoa(status),
	}
}

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

		levels := []struct {
			name      string
			giveLevel int
			want      string
		}{
			{
				name: "counts a fixture at the level of the overlay", giveLevel: 3,
				want: "go: level 3, 1 fixtures, 1 conformant, 0 declared differences\n",
			},
			{
				name: "leaves out a fixture above the level of the overlay", giveLevel: 2,
				want: "go: level 2, 0 fixtures, 0 conformant, 0 declared differences\n",
			},
		}
		for _, tt := range levels {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				tree := corpusOf(passing, noOutcomes)
				tree[corpusOverlay] = files.Text(fmt.Sprintf("language: go\nlevel: %d\ndiverge: []\n", tt.giveLevel))
				dir := files.Workspace(t, tree)

				stdout, stderr, status := runConform(t, dir, standalone, "--lang", "go")

				assert.Equal(t, status, 0, "the corpus conforms: "+stderr)
				assert.Contains(t, stdout, tt.want, "the summary counts the fixtures at the level and below it")
			})
		}

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
			assert.Contains(t, stderr, "the tests of the corpus failed: go test: exit status 1",
				"the gate names the failure")
			assert.Contains(t, stderr, " TestX\n", "the gate names the test that failed")
			assert.Contains(t, stderr, "the fixture fails", "the gate writes the output of the test that failed")
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

			_, stderr, status := runConform(t, t.TempDir(), standalone, "--lang", "rust")

			assert.Equal(t, status, 64, "a language without a binding is a usage error")
			assert.Contains(t, stderr, "has no binding of rust", "the error names the language")
		})

		t.Run("refuses an update of the expected outcomes from the corpus of TypeScript", func(t *testing.T) {
			t.Parallel()

			_, stderr, status := runConform(t, t.TempDir(), standalone, "--lang", "typescript", "--update")

			assert.Equal(t, status, 64, "an update from another corpus than the corpus of Go is a usage error")
			assert.Contains(t, stderr, "from the corpus of go alone", "the error names the corpus of the update")
		})

		t.Run("passes a corpus of TypeScript whose outcomes meet their expectations", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, typescriptCorpusOf(vitestReportOf("passed", []string{}), noOutcomes))

			stdout, stderr, status := runConform(t, dir, vitestEnv(dir, 0), "--lang", "typescript")

			assert.Equal(t, status, 0, "the corpus conforms: "+stderr)
			assert.Contains(t, stdout, "typescript: level 4, 1 fixtures, 1 conformant, 0 declared differences\n",
				"the summary counts the fixtures of TypeScript")
		})

		t.Run("fails a planted difference of TypeScript under the fixture of the directory of its module",
			func(t *testing.T) {
				t.Parallel()
				report := vitestReportOf(
					"passed",
					[]string{},
					typescriptRecording(spec.WrongReason, spec.ValueBesideError),
				)
				dir := files.Workspace(t, typescriptCorpusOf(report, cancelled))

				stdout, _, status := runConform(t, dir, vitestEnv(dir, 0), "--lang", "typescript")

				assert.Equal(t, status, 1, "a finding fails the gate")
				assert.Contains(t, stdout, "enum/x\n"+
					"  generation     FAIL  x.get/cancel: want check of cancel, got absent\n"+
					"  falsification  FAIL  x.get/cancel: want proven by ignored-cancel, got wrong-reason by value-beside-error\n"+
					"typescript: level 4, 1 fixtures, 0 conformant, 0 declared differences\n",
					"the gate writes each finding of the module enum/x/x under the fixture enum/x")
			})

		t.Run("declares a difference of TypeScript that an entry of its overlay states", func(t *testing.T) {
			t.Parallel()
			tree := typescriptCorpusOf(vitestReportOf("passed", []string{}), cancelled)
			tree[tsOverlay] = files.Text("language: typescript\nlevel: 4\ndiverge:\n" +
				"  - fixture: enum/x\n    key: x.get/cancel\n    expect: check\n    actual: absent\n" +
				"    why: The function one takes no signal\n")
			dir := files.Workspace(t, tree)

			stdout, stderr, status := runConform(t, dir, vitestEnv(dir, 0), "--lang", "typescript")

			assert.Equal(t, status, 0, "the entry declares the check that the lock lacks: "+stdout+stderr)
			assert.Contains(t, stdout, "typescript: level 4, 1 fixtures, 1 conformant, 1 declared differences\n",
				"the summary counts the declared difference")
		})

		t.Run("runs the tests of the fixtures at the level of the overlay alone", func(t *testing.T) {
			t.Parallel()
			tree := typescriptCorpusOf(vitestReportOf("passed", []string{}), noOutcomes)
			tree[tsOverlay] = files.Text("language: typescript\nlevel: 3\ndiverge: []\n")
			tree["corpus/typescript/enum/y/y.ts"] = files.Text("/** Returns two. */\nexport const two = 2;\n")
			tree["corpus/expect/enum/y.yaml"] = files.Text(
				conformHeader + "fixture: enum/y\nlevel: 4\nsubjects: []\ncontrols: []\ndeclarations: {}\n")
			dir := files.Workspace(t, tree)

			_, stderr, status := runConform(t, dir, vitestEnv(dir, 0), "--lang", "typescript")

			assert.Equal(t, status, 0, "the corpus conforms at level 3: "+stderr)
			files.HasContent(t, filepath.Join(dir, tsFilters), "enum/x/",
				"Vitest runs the test files of enum/x and leaves out enum/y at level 4")
		})

		t.Run("fails a corpus of TypeScript whose tests fail", func(t *testing.T) {
			t.Parallel()
			report := vitestReportOf("failed", []string{"AssertionFailed: the fixture fails"})
			dir := files.Workspace(t, typescriptCorpusOf(report, noOutcomes))

			_, stderr, status := runConform(t, dir, vitestEnv(dir, 1), "--lang", "typescript")

			assert.Equal(t, status, 1, "a failed test fails the gate")
			assert.Contains(t, stderr, "the tests of the corpus failed: vitest: exit status 1",
				"the gate names the failure")
			assert.Contains(t, stderr, "xSuite suite records\nAssertionFailed: the fixture fails",
				"the gate writes the test that failed with its message")
		})

		t.Run("returns the error of a report of Vitest that the run does not write", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, typescriptCorpusOf(vitestReportOf("passed", []string{}), noOutcomes))
			env := vitestEnv(dir, 0)
			env[0] = vitestReportVar + "=" + filepath.Join(dir, "missing.json")

			_, stderr, status := runConform(t, dir, env, "--lang", "typescript")

			assert.Equal(t, status, 1, "the gate fails")
			assert.Contains(t, stderr, "read the report of Vitest", "the error names the step")
		})

		t.Run("removes the directory of the report of Vitest after the run", func(t *testing.T) {
			t.Parallel()
			tree := typescriptCorpusOf(vitestReportOf("passed", []string{}), noOutcomes)
			tree["tmp"] = files.Dir()
			dir := files.Workspace(t, tree)
			env := append(vitestEnv(dir, 0), tempDirEnv(filepath.Join(dir, "tmp"))...)

			_, stderr, status := runConform(t, dir, env, "--lang", "typescript")

			assert.Equal(t, status, 0, "the corpus conforms: "+stderr)
			files.Equal(t, os.DirFS(filepath.Join(dir, "tmp")), files.Tree{},
				"the run removes the directory of the report of Vitest")
		})

		t.Run("returns the error of a directory of the report of Vitest that the run cannot make", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, typescriptCorpusOf(vitestReportOf("passed", []string{}), noOutcomes))
			env := append(vitestEnv(dir, 0), tempDirEnv(filepath.Join(dir, tsReport))...)

			_, stderr, status := runConform(t, dir, env, "--lang", "typescript")

			assert.Equal(t, status, 1, "the gate fails")
			assert.Contains(t, stderr, "make the directory of the report of Vitest", "the error names the step")
		})

		t.Run("refuses a working directory that does not exist", func(t *testing.T) {
			t.Parallel()

			_, stderr, status := runConform(t, t.TempDir(), standalone, "--lang", "go", "-C", "missing")

			assert.Equal(t, status, 64, "a working directory that does not exist is a usage error")
			assert.Contains(t, stderr, "missing", "the error names the directory")
		})
	})
}
