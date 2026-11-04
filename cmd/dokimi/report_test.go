// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"
	"go.dokimi.dev/spec"
)

// The names of the file of go test -json and of the file of a mutation record
// in the workspace of a case.
const (
	runFile      = "run.json"
	mutationFile = "mutation.json"
)

// deadCheckReport are the lines of the attribution of the corpus fixture
// iface/mutation/deadcheck, from the run in run.jsonl and the mutation record
// that make report-golden writes into testdata/report. Close/smoke kills no mutant that
// it covers. The checks of smoke of Get and Put kill the constructor that
// returns nil, which Get/cancel kills too. The spec of Len alone covers Len.
const deadCheckReport = "  mutants     25   covered by checks   15 killed   10 survived\n" +
	"  dead         1   Store.Close/smoke: covered 5 mutants, killed none\n" +
	"  subsumed     4   Store.Get/smoke: every kill also by Store.Get/cancel\n" +
	"                   Store.Get/smoke#generated: every kill also by Store.Get/cancel\n" +
	"                   Store.Put/smoke: every kill also by Store.Get/cancel\n" +
	"                   Store.Put/smoke#generated: every kill also by Store.Get/cancel\n" +
	"  unattributed 2 survivors in iface/mutation/deadcheck/deadchecktest/inmemory.go:76-77, covered by no check\n"

// eventsOf returns the events of go test -json that write the suite record of
// the declaration Store in the application with one check of the result.
func eventsOf(application spec.Application, result spec.Result) string {
	data, _ := json.Marshal(spec.SuiteRecord{
		Spec:        "1.0.0",
		Suite:       spec.RecordSuite{Language: "go", Package: "example.com/store", Declaration: "Store"},
		Generator:   spec.RecordTool{Name: "dokimi", Version: "0.1.0"},
		Runtime:     spec.RecordTool{Name: "go.dokimi.dev/conform", Version: "v0.9.0"},
		Seed:        "0x5f3a1c9e27d4b860",
		Application: application,
		Checks:      []spec.CheckRecord{{Identity: "Store.Read/keeps", Key: "store.read/keeps", Result: result}},
		Refusals:    []spec.RefusalRecord{},
		Unclaimed:   []spec.UnclaimedRecord{},
	})
	return attr(string(data))
}

// attr returns the attribute event of go test -json of the value under the
// key dokimi.suite.1.
func attr(value string) string {
	line, _ := json.Marshal(map[string]string{
		"Action": "attr", "Package": "example.com/store/storetest", "Test": "TestStore", "Key": "dokimi.suite.1",
		"Value": value,
	})
	return string(line) + "\n"
}

// runReport runs the test binary as the binary with the arguments of report
// in the directory dir, and returns its standard output, its standard error
// and its exit status.
func runReport(t *testing.T, dir string, args ...string) (stdout, stderr string, status int) {
	t.Helper()
	binary, err := os.Executable()
	assert.NoError(t, err, "the path of the test binary")
	cmd := exec.CommandContext(t.Context(), binary, append([]string{reportName}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), childVar+"=1")
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	err = cmd.Run()
	if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
		return out.String(), errOut.String(), exit.ExitCode()
	}
	assert.NoError(t, err, "the binary runs")
	return out.String(), errOut.String(), 0
}

func TestReport(t *testing.T) {
	t.Parallel()

	t.Run("reportRun", func(t *testing.T) {
		t.Parallel()

		t.Run("writes the report of each suite of the file and exits 0", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{runFile: files.Text(eventsOf(spec.Falsification, spec.Proven))})

			stdout, stderr, status := runReport(t, dir, runFile)

			assert.Equal(t, status, 0, "a run without a failed proof exits 0")
			assert.Empty(t, stderr, "the binary writes no error")
			assert.Contains(
				t,
				stdout,
				"example.com/store  Store\n  claims       1   proven 1",
				"the report has the suite",
			)
			assert.NotContains(t, stdout, "  mutants", "a run without --mutation has no line of mutants")
		})

		t.Run("writes one JSON document of every suite under --json", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{runFile: files.Text(eventsOf(spec.OnSubject, spec.Pass))})

			stdout, _, status := runReport(t, dir, "--json", runFile)

			assert.Equal(t, status, 0, "the run exits 0")
			var doc struct {
				Suites []struct {
					Declaration string `json:"declaration"`
				} `json:"suites"`
			}
			assert.NoError(t, json.Unmarshal([]byte(stdout), &doc), "the output is one JSON document")
			assert.Length(t, doc.Suites, 1, "the document has the suite")
			assert.Equal(t, doc.Suites[0].Declaration, "Store", "the suite has its declaration")
		})

		tests := []struct {
			name       string
			giveEvents string
			wantError  string
		}{
			{
				name:       "exits 1 for a check whose falsification result is not proven",
				giveEvents: eventsOf(spec.Falsification, spec.WrongReason),
				wantError:  "Store.Read/keeps is wrong-reason in falsification",
			},
			{
				name:       "exits 1 for a record that does not validate",
				giveEvents: attr(`{"spec": "1.0.0"}`),
				wantError:  "run.json: example.com/store/storetest TestStore dokimi.suite.1: spec: the suite record does not validate",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				dir := files.Workspace(t, files.Tree{runFile: files.Text(tt.giveEvents)})

				_, stderr, status := runReport(t, dir, runFile)

				assert.Equal(t, status, 1, "the run exits 1")
				assert.Contains(t, stderr, tt.wantError, "the error names the failure")
			})
		}

		t.Run("fails when the standard output does not take the JSON document", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				runFile: files.Text(eventsOf(spec.OnSubject, spec.Pass)), "out": files.Text(""),
			})
			readOnly, err := os.Open(filepath.Join(dir, "out"))
			assert.NoError(t, err, "the output file opens")
			t.Cleanup(func() { assert.NoError(t, readOnly.Close(), "the output file closes") })
			binary, err := os.Executable()
			assert.NoError(t, err, "the path of the test binary")
			cmd := exec.CommandContext(t.Context(), binary, reportName, "--json", runFile)
			cmd.Dir, cmd.Env, cmd.Stdout = dir, append(os.Environ(), childVar+"=1"), readOnly
			var errOut bytes.Buffer
			cmd.Stderr = &errOut

			assert.HasError(t, cmd.Run(), "the run fails")
			assert.Contains(t, errOut.String(), "write the report", "the error names the report")
		})

		t.Run("fails for a file that does not exist", func(t *testing.T) {
			t.Parallel()
			path := filepath.Join("absent", runFile)

			_, stderr, status := runReport(t, t.TempDir(), path)

			assert.NotEqual(t, status, 0, "the run fails")
			assert.Contains(t, stderr, "read the run: open "+path+": ", "the error names the file once")
		})

		t.Run("attributes the mutants of each record of --mutation to the checks", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				runFile: files.Text(eventsOf(spec.OnSubject, spec.Pass)),
				mutationFile: files.Text(`{"record": "dokimi-mutate", "version": 2, "target": {"language": "go", ` +
					`"name": "example.com/store/storetest"}, "killers": "all", "coverage": "leaves", "mutants": [` +
					`{"key": "0000000000000001", "file": "store.go", "start": {"line": 4, "column": 2}, ` +
					`"end": {"line": 4, "column": 9}, "verdict": "survived", "coveredBy": ["TestStore/Read/keeps"]}]}`),
			})

			stdout, stderr, status := runReport(t, dir, "--mutation", mutationFile, runFile)

			assert.Equal(t, status, 1, "a dead check fails the run")
			assert.Contains(t, stdout, "  mutants      1   covered by checks   0 killed   1 survived\n"+
				"  dead         1   Store.Read/keeps: covered 1 mutant, killed none\n",
				"the report attributes the survivor to the check")
			assert.Contains(t, stderr, "Store.Read/keeps is dead: it covered 1 mutant and killed none",
				"the error names the dead check")
		})

		t.Run("exits 0 for a mutation record whose checks killed each mutant that they covered", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				runFile: files.Text(eventsOf(spec.OnSubject, spec.Pass)),
				mutationFile: files.Text(`{"record": "dokimi-mutate", "version": 2, "target": {"language": "go", ` +
					`"name": "example.com/store/storetest"}, "killers": "all", "coverage": "leaves", "mutants": [` +
					`{"key": "0000000000000001", "file": "store.go", "start": {"line": 4, "column": 2}, ` +
					`"end": {"line": 4, "column": 9}, "verdict": "killed", "tests": ["TestStore/Read/keeps"], ` +
					`"allKillers": true, "coveredBy": ["TestStore/Read/keeps"]}]}`),
			})

			stdout, stderr, status := runReport(t, dir, "--mutation", mutationFile, runFile)

			assert.Equal(t, status, 0, "a run without a dead check exits 0")
			assert.Empty(t, stderr, "the binary writes no error")
			assert.Contains(t, stdout, "  mutants      1   covered by checks   1 killed   0 survived\n",
				"the report attributes the kill to the check")
		})

		t.Run("exits 1 for a mutation record of another kind", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{
				runFile:      files.Text(eventsOf(spec.OnSubject, spec.Pass)),
				mutationFile: files.Text(`{"record": "stryker", "version": 2}`),
			})

			_, stderr, status := runReport(t, dir, "--mutation", mutationFile, runFile)

			assert.Equal(t, status, 1, "the run exits 1")
			assert.Contains(t, stderr, `mutation.json: record 1: the record "stryker" of version 2 is no record `+
				`dokimi-mutate of version 2`, "the error names the file and the record")
		})

		t.Run("fails for a mutation record that does not exist", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, files.Tree{runFile: files.Text(eventsOf(spec.OnSubject, spec.Pass))})

			_, stderr, status := runReport(t, dir, "--mutation", mutationFile, runFile)

			assert.NotEqual(t, status, 0, "the run fails")
			assert.Contains(t, stderr, "read the mutation record: open "+mutationFile+": ", "the error names the file")
		})

		t.Run("reports the dead, the subsumed and the unattributed of the corpus fixture of a dead check",
			func(t *testing.T) {
				t.Parallel()
				golden, err := filepath.Abs(filepath.Join("testdata", "report"))
				assert.NoError(t, err, "the directory of the golden files has an absolute path")

				stdout, stderr, status := runReport(t, golden, "--mutation", mutationFile, "run.jsonl")

				assert.Equal(t, status, 1, "the dead check of the fixture fails the run")
				assert.Contains(t, stdout, deadCheckReport, "the report has the findings of the fixture")
				assert.Contains(t, stderr, "Store.Close/smoke is dead: it covered 5 mutants and killed none",
					"the error names the dead check")
			})
	})
}
