// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// statusFailed is the status of a test file and of a test of a Vitest report
// that failed.
const statusFailed = "failed"

// vitestReport is the part of the JSON report of Vitest that Read and Failed
// read.
type vitestReport struct {
	TestResults []vitestFile `json:"testResults"`
}

// vitestFile is a test file of a Vitest report: its path, its status, the
// message of the error of the file, and its tests.
type vitestFile struct {
	Name             string       `json:"name"`
	Status           string       `json:"status"`
	Message          string       `json:"message"`
	AssertionResults []vitestTest `json:"assertionResults"`
}

// vitestTest is a test of a Vitest report: its name with the names of its
// blocks, its status, the messages of its failures, and the suite records in
// its meta.
type vitestTest struct {
	FullName        string   `json:"fullName"`
	Status          string   `json:"status"`
	FailureMessages []string `json:"failureMessages"`
	Meta            struct {
		Records []json.RawMessage `json:"dokimi.suite"`
	} `json:"meta"`
}

// vitestOf decodes the first JSON value of r. It returns the Vitest report
// when the value is an object with the key testResults. Otherwise it returns
// nil and a reader of every byte of r from the start, for the events of
// go test -json. It returns an error for a report whose test files do not
// decode.
func vitestOf(r io.Reader) (*vitestReport, io.Reader, error) {
	var head bytes.Buffer
	var first struct {
		TestResults json.RawMessage `json:"testResults"`
	}
	// The events of go test -json start with an event or with a line that is
	// no JSON value. A value that does not decode sets no field, because the
	// decoder reads the whole value before it sets one.
	_ = json.NewDecoder(io.TeeReader(r, &head)).Decode(&first)
	if first.TestResults == nil {
		return nil, io.MultiReader(&head, r), nil
	}
	var v vitestReport
	if err := json.Unmarshal(first.TestResults, &v.TestResults); err != nil {
		return nil, nil, fmt.Errorf("decode the Vitest report: %w", err)
	}
	return &v, nil, nil
}

// records returns the suite records in the meta of each test of the report,
// in the order of the report. It returns an error that names the file, the
// test and the index of each record that does not validate, beside the
// records that validate.
func (v *vitestReport) records() ([]Record, error) {
	var records []Record
	var errs []error
	for _, f := range v.TestResults {
		for _, t := range f.AssertionResults {
			for i, data := range t.Meta.Records {
				r, err := decodeRecord(f.Name, t.FullName, fmt.Sprintf("%s[%d]", kind, i), data)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				records = append(records, r)
			}
		}
	}
	return records, errors.Join(errs...)
}

// failures returns each test of the report that failed, with its failure
// messages, and each file that failed without a test that failed, with the
// message of its error.
func (v *vitestReport) failures() []Failure {
	var out []Failure
	for _, f := range v.TestResults {
		failed := false
		for _, t := range f.AssertionResults {
			if t.Status == statusFailed {
				failed = true
				out = append(out, Failure{Package: f.Name, Test: t.FullName, Output: lines(t.FailureMessages...)})
			}
		}
		if f.Status == statusFailed && !failed {
			out = append(out, Failure{Package: f.Name, Output: lines(f.Message)})
		}
	}
	return out
}

// lines returns each message that is not empty, followed by a newline.
func lines(messages ...string) string {
	var b strings.Builder
	for _, m := range messages {
		if m != "" {
			b.WriteString(m)
			b.WriteByte('\n')
		}
	}
	return b.String()
}
