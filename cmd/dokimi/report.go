// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.dokimi.dev/eidos/cli"
	gocli "go.stealthscale.io/cli"

	"go.dokimi.dev/dokimi/internal/report"
)

// The name and the help text of the command report.
const (
	reportName = "report"
	reportHelp = "Prints what each suite of a run of go test -json claimed, proved, refused and found vacuous."
)

// reportOptions are the options of the command report.
type reportOptions struct {
	// Files are the files of go test -json.
	Files []gocli.InputFile `arg:"" help:"The files of go test -json, or - for the standard input."`

	// JSON writes one JSON document of every suite in place of the text.
	JSON bool `help:"Write one JSON document of every suite."`
}

// reportDocument is the JSON document of the command report.
type reportDocument struct {
	Suites []report.Suite `json:"suites"`
}

// reportRun reads the suite records of each file, and writes the report of
// every suite to the standard output. It writes each record that does not
// validate, each check whose falsification result is not proven and each
// problem that ended an application to the standard error. It exits with the
// status 1 for any of them.
func reportRun(_ context.Context, run *gocli.Run[cli.Global, reportOptions]) error {
	var records []report.Record
	var errs []error
	for _, f := range run.Config().Files {
		r, err := run.Open(f)
		if err != nil {
			return fmt.Errorf("read the run: %w", err)
		}
		read, err := report.Read(r)
		// A file that the command only reads loses nothing on a close that fails.
		_ = r.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", f, err))
		}
		records = append(records, read...)
	}
	suites := report.Suites(records)
	if run.Config().JSON {
		if err := json.NewEncoder(run.Stdout()).Encode(reportDocument{Suites: suites}); err != nil {
			return fmt.Errorf("write the report: %w", err)
		}
	} else {
		report.Print(run.Stdout(), suites)
	}
	for i := range suites {
		for _, f := range suites[i].Failures() {
			errs = append(errs, errors.New(f))
		}
	}
	if err := errors.Join(errs...); err != nil {
		fmt.Fprintln(run.Stderr(), err)
		return cli.ExitError{Code: cli.StatusFailed}
	}
	return nil
}
