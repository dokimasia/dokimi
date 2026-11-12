// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"go.dokimi.dev/eidos/cli"
	gocli "go.stealthscale.io/cli"
	"golang.org/x/mod/modfile"

	"go.dokimi.dev/dokimi"
	"go.dokimi.dev/dokimi/internal/corpus"
	"go.dokimi.dev/dokimi/internal/report"
)

// The name and the help text of the command conform.
const (
	conformName = "conform"
	conformHelp = "Judges the corpus of a language by the expected outcomes of its fixtures."
)

// conformLang is the language whose corpus the expected outcomes come from.
const conformLang = "go"

// conformOptions are the options of the command conform.
type conformOptions struct {
	// Lang is the language of the corpus under corpus/<language>.
	Lang string `help:"The language of the corpus, such as go or typescript." required:""`

	// Update rewrites the declarations of the expected outcomes.
	Update bool `help:"Rewrite the declarations of the expected outcomes from the lock and the records of the run."`
}

// binding is the binding of the gate to the corpus of one language under
// corpus/<language>. tests runs the tests of the fixtures in the directory of
// the corpus and returns their report, beside the error of a run whose tests
// failed. module returns the module path of the corpus, and lock is the file
// of the lock of the plan of the language. fixtureOf returns the package of
// the fixture of a package of a line of the lock or of a record.
type binding struct {
	tests     func(ctx context.Context, dir string, fixtures []string, stderr io.Writer) ([]byte, error)
	module    func(dir string) (string, error)
	fixtureOf func(pkg string) string
	lock      string
}

// bindings are the bindings of the gate by language. A package of Go is a
// directory, and a module of TypeScript is a file of the directory of its
// package.
var bindings = map[string]binding{
	conformLang: {
		tests: goTests, module: goModule, fixtureOf: func(pkg string) string { return pkg },
		lock: "dokimi." + dokimi.GoPlan + ".lock",
	},
	dokimi.TypeScriptPlan: {
		tests: vitestTests, module: func(string) (string, error) { return "", nil }, fixtureOf: path.Dir,
		lock: "dokimi." + dokimi.TypeScriptPlan + ".lock",
	},
}

// vitest is the command that runs Vitest in the corpus of TypeScript.
var vitest = []string{"pnpm", "exec", "vitest"}

// conformRun judges the corpus of a language under the directory corpus of
// the working directory. It checks that the generated files of the corpus are
// up to date with dokimi run --check, runs the tests of the fixtures at the
// level of the overlay of the language and below it, and judges the lock of
// the plan of the language and the suite records of the run by the expected
// outcomes under corpus/expect and by the overlay. The tests of Go run with
// go test -count=1 -json in the corpus, and the tests of TypeScript with
// Vitest and its JSON report. It writes each finding and a summary to the
// standard output. Under --update it rewrites the declarations of each file
// of expected outcomes from the run of the corpus of Go in place of the
// judgement.
//
// conformRun exits with the status of the generation when the generated files
// are out of date, and with the status 1 when a test fails, a record does not
// validate, or a finding remains.
func conformRun(ctx context.Context, run *gocli.Run[cli.Global, conformOptions]) error {
	lang := run.Config().Lang
	b, bound := bindings[lang]
	if !bound {
		return &cli.UsageError{Err: fmt.Errorf("the gate has no binding of %s", lang)}
	}
	if run.Config().Update && lang != conformLang {
		return &cli.UsageError{Err: fmt.Errorf("the gate updates the expected outcomes from the corpus of %s alone",
			conformLang)}
	}
	stdio, err := cli.IOOf(run)
	if err != nil {
		return fmt.Errorf("read the working directory: %w", err)
	}
	expect, dir := filepath.Join(stdio.Dir, "corpus", "expect"), filepath.Join(stdio.Dir, "corpus", lang)
	expected, err := corpus.ReadExpectations(os.DirFS(expect))
	if err != nil {
		return fmt.Errorf("read the expected outcomes: %w", err)
	}
	if len(expected) == 0 {
		return fmt.Errorf("read the expected outcomes: %s has no file of expected outcomes", expect)
	}
	fixtures, err := corpus.Fixtures(os.DirFS(dir), expected)
	if err != nil {
		return fmt.Errorf("list the fixtures of the corpus: %w", err)
	}
	overlay, err := corpus.ReadOverlay(os.DirFS(dir), "overlay.yaml")
	if err != nil {
		return fmt.Errorf("read the overlay: %w", err)
	}
	module, err := b.module(dir)
	if err != nil {
		return fmt.Errorf("read the module of the corpus: %w", err)
	}
	env := *gocli.OS()
	env.Args, env.Stdout, env.Stderr = []string{"-C", dir, "run", "--check", "--no-color"}, run.Stdout(), run.Stderr()
	if status := cli.Run(ctx, cli.New(help, dokimi.Compose), dokimi.Compose, &env); status != cli.StatusOK {
		return cli.ExitError{Code: status}
	}
	var judged []string
	for i := range expected {
		if expected[i].Level <= overlay.Level {
			judged = append(judged, expected[i].Fixture)
		}
	}
	events, failed := b.tests(ctx, dir, judged, run.Stderr())
	var errs []error
	if failed != nil {
		errs = append(errs, fmt.Errorf("the tests of the corpus failed: %w", failed))
		// Read reports an event that does not decode, so the error of Failed
		// repeats it.
		failures, _ := report.Failed(bytes.NewReader(events))
		for _, f := range failures {
			fmt.Fprintf(run.Stderr(), "%s\n%s", strings.TrimSpace(f.Package+" "+f.Test), f.Output)
		}
	}
	records, err := report.Read(bytes.NewReader(events))
	errs = append(errs, err)
	// A corpus without a suite has no lock.
	lock, err := os.ReadFile(filepath.Join(dir, b.lock))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("read the lock: %w", err)
	}
	lines := strings.Split(string(lock), "\n")
	for i, text := range lines {
		if pkg, rest, tabbed := strings.Cut(text, "\t"); tabbed && !strings.HasPrefix(text, "#") {
			lines[i] = b.fixtureOf(pkg) + "\t" + rest
		}
	}
	for i := range records {
		r := &records[i]
		r.Suite.Package = b.fixtureOf(r.Suite.Package)
		if at := strings.LastIndex(r.Subject, "."); at >= 0 {
			r.Subject = b.fixtureOf(r.Subject[:at]) + r.Subject[at:]
		}
	}
	outcomes := corpus.Outcomes(module, fixtures, lines, records)
	if run.Config().Update {
		updated := corpus.Update(expected, outcomes)
		for i := range updated {
			if err := corpus.WriteExpectation(expect, &updated[i]); err != nil {
				return fmt.Errorf("write the expected outcomes: %w", err)
			}
		}
		fmt.Fprintf(run.Stdout(), "%s: wrote %d files of expected outcomes\n", lang, len(updated))
	} else {
		found, declared := corpus.Judge(expected, overlay, outcomes)
		printJudgement(run.Stdout(), lang, expected, overlay.Level, found, declared)
		if len(found) > 0 {
			errs = append(errs, fmt.Errorf("the corpus of %s has %d findings", lang, len(found)))
		}
	}
	if err := errors.Join(errs...); err != nil {
		fmt.Fprintln(run.Stderr(), err)
		return cli.ExitError{Code: cli.StatusFailed}
	}
	return nil
}

// goTests runs go test -count=1 -json over every package of the corpus of Go
// in dir and returns its events. The level of the overlay judges the fixtures,
// and every test of the corpus runs.
func goTests(ctx context.Context, dir string, _ []string, stderr io.Writer) ([]byte, error) {
	var events bytes.Buffer
	test := exec.CommandContext(ctx, "go", "test", "-count=1", "-json", "./...")
	test.Dir, test.Stdout, test.Stderr = dir, &events, stderr
	if err := test.Run(); err != nil {
		return events.Bytes(), fmt.Errorf("go test: %w", err)
	}
	return events.Bytes(), nil
}

// goModule returns the module path of the go.mod of the corpus in dir.
func goModule(dir string) (string, error) {
	mod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	return modfile.ModulePath(mod), err
}

// vitestTests runs Vitest in the corpus of TypeScript in dir over the test
// files of fixtures, and returns its JSON report. A filter of Vitest selects
// the files whose path contains it, so each filter is the path of a fixture
// with a slash, which iface/lang/genericbound does not contain for
// iface/lang/generic. Vitest writes its own output to stderr.
func vitestTests(ctx context.Context, dir string, fixtures []string, stderr io.Writer) ([]byte, error) {
	out, err := os.MkdirTemp("", "dokimi-conform-")
	if err != nil {
		return nil, fmt.Errorf("make the directory of the report of Vitest: %w", err)
	}
	defer func() { _ = os.RemoveAll(out) }()
	file := filepath.Join(out, "vitest.json")
	args := append(slices.Clone(vitest[1:]), "run", "--reporter=json", "--outputFile="+file)
	for _, f := range fixtures {
		args = append(args, f+"/")
	}
	test := exec.CommandContext(ctx, vitest[0], args...)
	test.Dir, test.Stdout, test.Stderr = dir, stderr, stderr
	var failed error
	if ran := test.Run(); ran != nil {
		failed = fmt.Errorf("vitest: %w", ran)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Join(failed, fmt.Errorf("read the report of Vitest: %w", err))
	}
	return data, failed
}

// printJudgement writes the findings of each fixture under the path of the
// fixture, and the summary of the corpus of lang at the level: the fixtures
// that the gate judged, those without a finding, and the differences that the
// overlay declared.
func printJudgement(
	w io.Writer, lang string, expected []corpus.Expectation, level int, found []corpus.Finding, declared int,
) {
	failed := map[string]bool{}
	for i := range found {
		f := &found[i]
		if !failed[f.Fixture] {
			fmt.Fprintln(w, f.Fixture)
		}
		failed[f.Fixture] = true
		key := ""
		if f.Key != "" {
			key = f.Key + ": "
		}
		fmt.Fprintf(w, "  %-13s  FAIL  %swant %s, got %s\n", f.Phase, key, f.Want, f.Got)
	}
	judged, conformant := 0, 0
	for i := range expected {
		if expected[i].Level <= level {
			judged++
			if !failed[expected[i].Fixture] {
				conformant++
			}
		}
	}
	fmt.Fprintf(w, "%s: level %d, %d fixtures, %d conformant, %d declared differences\n",
		lang, level, judged, conformant, declared)
}
