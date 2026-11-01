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
	"path/filepath"
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

// conformLang is the language whose binding the gate judges.
const conformLang = "go"

// conformOptions are the options of the command conform.
type conformOptions struct {
	// Lang is the language of the corpus under corpus/<language>.
	Lang string `help:"The language of the corpus, such as go." required:""`

	// Update rewrites the declarations of the expected outcomes.
	Update bool `help:"Rewrite the declarations of the expected outcomes from the lock and the records of the run."`
}

// conformRun judges the corpus of a language under the directory corpus of
// the working directory. It checks that the generated files of the corpus are
// up to date with dokimi run --check, runs go test -count=1 -json in the
// corpus, and judges the lock of the plan go and the suite records of the run
// by the expected outcomes under corpus/expect and by the overlay of the
// language. It writes each finding and a summary to the standard output.
// Under --update it rewrites the declarations of each file of expected
// outcomes from the run in place of the judgement.
//
// conformRun exits with the status of the generation when the generated files
// are out of date, and with the status 1 when a test fails, a record does not
// validate, or a finding remains.
func conformRun(ctx context.Context, run *gocli.Run[cli.Global, conformOptions]) error {
	lang := run.Config().Lang
	if lang != conformLang {
		return &cli.UsageError{Err: fmt.Errorf("the gate judges the corpus of go and has no binding of %s", lang)}
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
	mod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return fmt.Errorf("read the module of the corpus: %w", err)
	}
	env := *gocli.OS()
	env.Args, env.Stdout, env.Stderr = []string{"-C", dir, "run", "--check", "--no-color"}, run.Stdout(), run.Stderr()
	if status := cli.Run(ctx, cli.New(help, dokimi.Compose), dokimi.Compose, &env); status != cli.StatusOK {
		return cli.ExitError{Code: status}
	}
	var events bytes.Buffer
	test := exec.CommandContext(ctx, "go", "test", "-count=1", "-json", "./...")
	test.Dir, test.Stdout, test.Stderr = dir, &events, run.Stderr()
	var errs []error
	if failed := test.Run(); failed != nil {
		errs = append(errs, fmt.Errorf("the tests of the corpus failed: %w", failed))
	}
	records, err := report.Read(&events)
	errs = append(errs, err)
	// A corpus without a suite has no lock.
	lock, err := os.ReadFile(filepath.Join(dir, "dokimi.go.lock"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("read the lock: %w", err)
	}
	outcomes := corpus.Outcomes(modfile.ModulePath(mod), fixtures, strings.Split(string(lock), "\n"), records)
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
