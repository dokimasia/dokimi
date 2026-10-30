// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package report reads the suite records of a run from the events of
// go test -json, and prints what each suite claimed, proved, refused and
// found vacuous.
package report

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"go.dokimi.dev/dokimi/spec"
)

// keyPrefix starts the key of each attribute event of a suite record, which
// ends with the number of the record in its test.
const keyPrefix = "dokimi.suite."

// maxLine is the longest line of an event that Read reads.
const maxLine = 1 << 20

// Record is one suite record of a run, with the package and the test of
// go test that wrote it.
type Record struct {
	Package string `json:"package"`
	Test    string `json:"test"`
	spec.SuiteRecord
}

// event is the part of an event of go test -json that Read reads. An event
// of the action attr alone has a key.
type event struct {
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Key     string `json:"Key"`
	Value   string `json:"Value"`
}

// Read reads the events of go test -json from r, and returns the suite
// records in the order of their first events. It joins the values of the
// attribute events of one key of one test into one record. It skips a line
// that is no JSON object, such as the output of a build that failed.
//
// Read returns an error that names the line for an event that does not
// decode, and an error that names the package, the test and the key of each
// record that does not validate against the schema of the suite record. It
// returns the records that validate beside the second error.
func Read(r io.Reader) ([]Record, error) {
	type part struct{ pkg, test, key string }
	var order []part
	values := map[part]*strings.Builder{}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, maxLine)
	for line := 1; scanner.Scan(); line++ {
		text := scanner.Bytes()
		if len(text) == 0 || text[0] != '{' {
			continue
		}
		var e event
		if err := json.Unmarshal(text, &e); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		if !strings.HasPrefix(e.Key, keyPrefix) {
			continue
		}
		p := part{e.Package, e.Test, e.Key}
		b, seen := values[p]
		if !seen {
			b = &strings.Builder{}
			values[p] = b
			order = append(order, p)
		}
		b.WriteString(e.Value)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read the events: %w", err)
	}
	var records []Record
	var errs []error
	for _, p := range order {
		data := []byte(values[p].String())
		if err := spec.ValidateSuiteRecord(data); err != nil {
			errs = append(errs, fmt.Errorf("%s %s %s: %w", p.pkg, p.test, p.key, err))
			continue
		}
		r := Record{Package: p.pkg, Test: p.test}
		_ = json.Unmarshal(data, &r.SuiteRecord)
		records = append(records, r)
	}
	return records, errors.Join(errs...)
}

// Suite contains the records of one declaration.
type Suite struct {
	// Package is the package of the declaration.
	Package string `json:"package"`

	Declaration string `json:"declaration"`

	// Records are the records of the declaration in the order of the run.
	Records []Record `json:"records"`

	// Vacuous are the identities of the checks that each subject application
	// that judged them found vacuous. An application that excused or skipped
	// a check, or left it unchecked, did not judge it.
	Vacuous []string `json:"vacuous"`
}

// Suites groups the records by the package and the name of their
// declaration, in the order of the first record of each declaration.
func Suites(records []Record) []Suite {
	var suites []Suite
	for i := range records {
		r := &records[i]
		at := slices.IndexFunc(suites, func(s Suite) bool {
			return s.Package == r.Suite.Package && s.Declaration == r.Suite.Declaration
		})
		if at < 0 {
			suites = append(suites, Suite{Package: r.Suite.Package, Declaration: r.Suite.Declaration})
			at = len(suites) - 1
		}
		suites[at].Records = append(suites[at].Records, *r)
	}
	for i := range suites {
		suites[i].Vacuous = vacuous(suites[i].Records)
	}
	return suites
}

// Failures returns each check of a falsification record of the suite whose
// result is not proven, and each problem that ended an application of the
// suite before its checks. A check that an excusal covers or that skipped, and
// an allocation ceiling that the run did not check, has no result to judge.
func (s *Suite) Failures() []string {
	var found []string
	for _, c := range verdicts(s.Records, spec.Falsification) {
		switch c.Result {
		case spec.Proven, spec.Excused, spec.Skipped, spec.Unchecked:
		default:
			found = append(found, fmt.Sprintf("%s is %s in falsification", c.Identity, c.Result))
		}
	}
	for _, e := range endings(s.Records) {
		found = append(found, "the application in "+e.test+" ended before its checks: "+e.problem)
	}
	return found
}

// ending is a problem that ended an application before its checks, with the
// test of the application.
type ending struct{ test, problem string }

// endings returns each problem of the field ended of the records, in their
// order. Both records of one application state the same problems, so endings
// returns the problems of one test once.
func endings(records []Record) []ending {
	var found []ending
	seen := map[string]bool{}
	for i := range records {
		r := &records[i]
		key := r.Test + "\n" + r.Ended
		if seen[key] {
			continue
		}
		seen[key] = true
		for problem := range strings.Lines(r.Ended) {
			found = append(found, ending{test: r.Test, problem: strings.TrimSuffix(problem, "\n")})
		}
	}
	return found
}

// verdicts returns the check records of the records of the application.
func verdicts(records []Record, application spec.Application) []*spec.CheckRecord {
	var found []*spec.CheckRecord
	for i := range records {
		if records[i].Application != application {
			continue
		}
		for j := range records[i].Checks {
			found = append(found, &records[i].Checks[j])
		}
	}
	return found
}

// vacuous returns the identities of the checks that each subject application
// of the records that judged them found vacuous, in the order of their first
// verdicts.
func vacuous(records []Record) []string {
	var order []string
	judged, empty := map[string]int{}, map[string]int{}
	for _, c := range verdicts(records, spec.OnSubject) {
		if c.Result == spec.Excused || c.Result == spec.Skipped || c.Result == spec.Unchecked {
			continue
		}
		if judged[c.Identity] == 0 {
			order = append(order, c.Identity)
		}
		judged[c.Identity]++
		if c.Result == spec.Vacuous {
			empty[c.Identity]++
		}
	}
	return slices.DeleteFunc(order, func(id string) bool { return empty[id] < judged[id] })
}

// claimed are the results of falsification that Print always writes.
var claimed = []spec.Result{spec.Proven, spec.WrongReason, spec.Passed}

// judgedOn are the results on a subject that Print always writes.
var judgedOn = []spec.Result{spec.Pass, spec.Vacuous, spec.Fail, spec.Excused, spec.Skipped}

// Print writes one block for each suite to w: the counts of the results of
// falsification and on the subject, the refusals, the unclaimed
// classifications, the vacuous checks, the problems that ended an application
// before its checks, and the seed and the versions of the first record. A
// count line writes each other result that occurs after the results that it
// always writes.
func Print(w io.Writer, suites []Suite) {
	for i := range suites {
		s := &suites[i]
		first := &s.Records[0]
		fmt.Fprintf(w, "%s  %s\n", s.Package, s.Declaration)
		results(w, "claims", s.Records, spec.Falsification, claimed)
		results(w, "subject", s.Records, spec.OnSubject, judgedOn)
		refused := make([]string, len(first.Refusals))
		for i, r := range first.Refusals {
			refused[i] = fmt.Sprintf("%s/%s: %s", r.Callable, r.Claim, r.Why)
		}
		entries(w, "refused", refused)
		unclaimed := make([]string, len(first.Unclaimed))
		for i, u := range first.Unclaimed {
			unclaimed[i] = fmt.Sprintf("%s on %s: %s", u.Classification, u.Callable, u.Note)
		}
		entries(w, "unclaimed", unclaimed)
		for _, id := range s.Vacuous {
			fmt.Fprintf(w, "  %-10s       %s\n", "vacuous", id)
		}
		for _, e := range endings(s.Records) {
			fmt.Fprintf(w, "  %-10s       %s: %s\n", "ended", e.test, e.problem)
		}
		fmt.Fprintf(w, "seed %s   spec %s   %s %s\n", first.Seed, first.Spec, first.Runtime.Name, first.Runtime.Version)
	}
}

// results writes the line label with the number of the checks of the records
// of the application, the count of each result of always, and the count of
// each other result that occurs, as "proven 11".
func results(w io.Writer, label string, records []Record, application spec.Application, always []spec.Result) {
	total, by := 0, map[spec.Result]int{}
	for _, c := range verdicts(records, application) {
		total++
		by[c.Result]++
	}
	parts := make([]string, 0, len(spec.Results))
	for _, result := range always {
		parts = append(parts, fmt.Sprintf("%s %d", result, by[result]))
	}
	for _, result := range spec.Results {
		if by[result] > 0 && !slices.Contains(always, result) {
			parts = append(parts, fmt.Sprintf("%s %d", result, by[result]))
		}
	}
	fmt.Fprintf(w, "  %-10s %3d   %s\n", label, total, strings.Join(parts, "   "))
}

// entries writes the line label with the number of entries and the first
// entry, and each other entry on a line of its own.
func entries(w io.Writer, label string, list []string) {
	if len(list) == 0 {
		fmt.Fprintf(w, "  %-10s %3d\n", label, 0)
		return
	}
	fmt.Fprintf(w, "  %-10s %3d   %s\n", label, len(list), list[0])
	for _, e := range list[1:] {
		fmt.Fprintf(w, "  %-10s       %s\n", "", e)
	}
}
