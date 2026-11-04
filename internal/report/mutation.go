// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"go.dokimi.dev/spec"
)

// The kind and the version of the records of mutate-spec that ReadMutations
// reads.
const (
	mutationKind    = "dokimi-mutate"
	mutationVersion = 2
)

// The modes of a mutation run that dead and subsumed checks need: every test
// that killed a mutant, and the coverage of each leaf test.
const (
	killersAll     = "all"
	coverageLeaves = "leaves"
)

// The verdicts of a mutant that Attribute attributes.
const (
	verdictKilled   = "killed"
	verdictSurvived = "survived"
)

// Mutation is the record of a mutation run of mutate-spec, with the fields
// that Attribute reads.
type Mutation struct {
	// Target is the package whose code the run changed.
	Target MutationTarget `json:"target"`

	// Killers is all for a run that names every test that failed on a mutant,
	// and first or empty for a run that names the first.
	Killers string `json:"killers,omitempty"`

	// Coverage is leaves for a run that names the leaf tests that executed the
	// site of a mutant, such as the test of a check, and tests or empty for a
	// run that names the top-level tests.
	Coverage string `json:"coverage,omitempty"`

	Mutants []Mutant `json:"mutants"`
}

// MutationTarget is the target of a mutation run: its language and its name,
// such as the import path of a Go package.
type MutationTarget struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

// Mutant is one mutant of a mutation record. A record names a test of the
// target by its name, such as TestStore/Get/smoke, and a test of another
// package by the package, a colon, a space and the name.
type Mutant struct {
	Key string `json:"key"`

	// File is the path of the file of the mutant's site, relative to the root
	// of the record.
	File string `json:"file"`

	Verdict string `json:"verdict"`

	// Tests are the tests that failed on the mutant.
	Tests []string `json:"tests,omitempty"`

	// CoveredBy are the tests whose run alone executed the mutant's site.
	CoveredBy []string `json:"coveredBy,omitempty"`

	Start Position `json:"start"`
	End   Position `json:"end"`

	// AllKillers is true when Tests names every test that failed on the mutant.
	AllKillers bool `json:"allKillers,omitempty"`
}

// Position is a line and a column of a file, both counted from 1.
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// ReadMutations reads the records of mutation runs from r: one record, or one
// record after another, as dokimi-mutate-go writes them with --json. It returns
// an error with the number of the record for a record of another kind or
// version, and for a record that does not decode.
func ReadMutations(r io.Reader) ([]Mutation, error) {
	decoder := json.NewDecoder(r)
	var found []Mutation
	for n := 1; ; n++ {
		var record struct {
			Kind string `json:"record"`
			Mutation
			Version int `json:"version"`
		}
		err := decoder.Decode(&record)
		if errors.Is(err, io.EOF) {
			return found, nil
		}
		if err != nil {
			return nil, fmt.Errorf("record %d: %w", n, err)
		}
		if record.Kind != mutationKind || record.Version != mutationVersion {
			return nil, fmt.Errorf("record %d: the record %q of version %d is no record %s of version %d",
				n, record.Kind, record.Version, mutationKind, mutationVersion)
		}
		found = append(found, record.Mutation)
	}
}

// Attribution is what the checks of a suite covered and killed of the mutants
// of mutation records.
type Attribution struct {
	// Dead are the checks that covered a mutant and killed none.
	Dead []DeadCheck `json:"dead"`

	// Subsumed are the checks whose every kill another check of the suite also
	// made, and which killed fewer mutants than that check.
	Subsumed []SubsumedCheck `json:"subsumed"`

	// Unattributed are the survivors that the checks did not cover, by file, of
	// the targets whose mutants a check of the suite covered.
	Unattributed []Survivors `json:"unattributed"`

	// Lacking are the targets whose records lack a mode that dead and subsumed
	// checks need, of the targets whose mutants a check of the suite covered.
	Lacking []Lack `json:"lacking"`

	// Covered is the number of the mutants that a check of the suite covered.
	// Killed of them were killed, and Survived survived.
	Covered  int `json:"covered"`
	Killed   int `json:"killed"`
	Survived int `json:"survived"`
}

// DeadCheck is a check that covered mutants and killed none.
type DeadCheck struct {
	Identity string `json:"identity"`

	// Covered is the number of the mutants that the check covered.
	Covered int `json:"covered"`
}

// SubsumedCheck is a check whose every kill the check By also made.
type SubsumedCheck struct {
	Identity string `json:"identity"`
	By       string `json:"by"`
}

// Survivors are the survivors of a file of a target that the checks did not
// cover: their number, and the first line and the last line of their sites.
type Survivors struct {
	Target string `json:"target"`
	File   string `json:"file"`
	First  int    `json:"first"`
	Last   int    `json:"last"`
	Count  int    `json:"count"`
}

// Lack is a target whose record lacks the modes of a mutation run that dead
// and subsumed checks need, such as killers all.
type Lack struct {
	Target string   `json:"target"`
	Modes  []string `json:"modes"`
}

// check is a check of a suite that a name of a test of a mutation record
// joins.
type check struct {
	identity string
	suite    int
}

// tally is what one check covered and killed of the mutants whose records name
// every killer and the coverage of each leaf test.
type tally struct {
	killed  map[string]bool
	covered int
}

// Attribute sets the attribution of the mutants of the records to the checks
// of each suite. It joins a name of a test of a record with the checks of the
// subject records of the suites. In a run of the coverage of each leaf test,
// the name of a check's test is the name of the test of its record, a slash and
// the identity of the check without its declaration. In any other run, the
// name of the top-level test of a record joins every check of the record.
//
// A killed mutant whose record names every killer, and a survivor, count
// toward dead and subsumed checks in a record of both modes. A kill whose run
// ended before every killer failed counts as killed alone.
func Attribute(suites []Suite, mutations []Mutation) {
	for i := range suites {
		suites[i].Mutation = &Attribution{}
	}
	tallies := map[check]*tally{}
	for i := range mutations {
		m := &mutations[i]
		complete := m.Killers == killersAll && m.Coverage == coverageLeaves
		names := joins(suites, m.Target.Name, m.Coverage == coverageLeaves)
		reached := map[int]bool{}
		var uncovered []*Mutant
		for j := range m.Mutants {
			mu := &m.Mutants[j]
			if mu.Verdict != verdictKilled && mu.Verdict != verdictSurvived {
				continue
			}
			covering := owners(names, mu.CoveredBy)
			if len(covering) == 0 {
				if mu.Verdict == verdictSurvived {
					uncovered = append(uncovered, mu)
				}
				continue
			}
			count(suites, covering, mu.Verdict, reached)
			if !complete || (mu.Verdict == verdictKilled && !mu.AllKillers) {
				continue
			}
			for _, c := range covering {
				tallyOf(tallies, c).covered++
			}
			id := m.Target.Name + "\x00" + mu.Key
			for _, c := range owners(names, mu.Tests) {
				tallyOf(tallies, c).killed[id] = true
			}
		}
		lack := lacking(m.Killers, m.Coverage)
		for s := range suites {
			if !reached[s] {
				continue
			}
			a := suites[s].Mutation
			a.Unattributed = append(a.Unattributed, survivors(m.Target.Name, uncovered)...)
			if len(lack) > 0 {
				a.Lacking = append(a.Lacking, Lack{Target: m.Target.Name, Modes: lack})
			}
		}
	}
	judge(suites, tallies)
}

// joins returns the checks of the subject records of the suites by the names of
// their tests in a record of the target: the name of each check's test for a
// run of leaves, and the name of the top-level test of each record for any
// other run. A test of another package has the package, a colon and a space
// before its name.
func joins(suites []Suite, target string, leaves bool) map[string][]check {
	names := map[string][]check{}
	for i := range suites {
		declaration := suites[i].Declaration
		for j := range suites[i].Records {
			r := &suites[i].Records[j]
			if r.Application != spec.OnSubject {
				continue
			}
			prefix := ""
			if r.Package != target {
				prefix = r.Package + ": "
			}
			top, _, _ := strings.Cut(r.Test, "/")
			for k := range r.Checks {
				identity := r.Checks[k].Identity
				name := prefix + top
				if leaves {
					name = prefix + r.Test + "/" + subtest(declaration, identity)
				}
				names[name] = append(names[name], check{identity: identity, suite: i})
			}
		}
	}
	return names
}

// subtest returns the name of the test of the check with the identity in the
// suite of the declaration: the identity without the declaration, such as
// Get/keeps-the-value for Store.Get/keeps-the-value. The identity of a check of
// a function, such as Slug/keeps-the-case, is its own name.
func subtest(declaration, identity string) string {
	if rest, cut := strings.CutPrefix(identity, declaration+"."); cut {
		return rest
	}
	return identity
}

// owners returns the checks that the names join, each once, in the order of the
// names.
func owners(names map[string][]check, tests []string) []check {
	var found []check
	for _, name := range tests {
		for _, c := range names[name] {
			if !slices.Contains(found, c) {
				found = append(found, c)
			}
		}
	}
	return found
}

// count counts a mutant of the verdict once in each suite of the checks that
// covered it, and marks each such suite as reached.
func count(suites []Suite, covering []check, verdict string, reached map[int]bool) {
	counted := map[int]bool{}
	for _, c := range covering {
		reached[c.suite] = true
		if counted[c.suite] {
			continue
		}
		counted[c.suite] = true
		a := suites[c.suite].Mutation
		a.Covered++
		if verdict == verdictKilled {
			a.Killed++
		} else {
			a.Survived++
		}
	}
}

// tallyOf returns the tally of c, and adds a new tally of c to tallies when
// tallies lacks one.
func tallyOf(tallies map[check]*tally, c check) *tally {
	t, ok := tallies[c]
	if !ok {
		t = &tally{killed: map[string]bool{}}
		tallies[c] = t
	}
	return t
}

// lacking returns the modes that a record of the killers and the coverage
// lacks for dead and subsumed checks.
func lacking(killers, coverage string) []string {
	var modes []string
	if killers != killersAll {
		modes = append(modes, "killers "+killersAll)
	}
	if coverage != coverageLeaves {
		modes = append(modes, "coverage "+coverageLeaves)
	}
	return modes
}

// survivors returns the survivors of the target by file, in the order of the
// files.
func survivors(target string, uncovered []*Mutant) []Survivors {
	var found []Survivors
	for _, mu := range uncovered {
		at := slices.IndexFunc(found, func(s Survivors) bool { return s.File == mu.File })
		if at < 0 {
			found = append(found, Survivors{Target: target, File: mu.File, First: mu.Start.Line, Last: mu.End.Line})
			at = len(found) - 1
		}
		s := &found[at]
		s.First, s.Last = min(s.First, mu.Start.Line), max(s.Last, mu.End.Line)
		s.Count++
	}
	slices.SortFunc(found, func(a, b Survivors) int { return strings.Compare(a.File, b.File) })
	return found
}

// judge sets the dead and the subsumed checks of each suite from the tallies,
// in the order of the checks of its subject records.
func judge(suites []Suite, tallies map[check]*tally) {
	for i := range suites {
		order := identities(&suites[i])
		a := suites[i].Mutation
		for _, id := range order {
			t := tallies[check{identity: id, suite: i}]
			switch {
			case t == nil:
			case len(t.killed) == 0:
				a.Dead = append(a.Dead, DeadCheck{Identity: id, Covered: t.covered})
			default:
				if by, ok := subsumer(i, order, t, tallies); ok {
					a.Subsumed = append(a.Subsumed, SubsumedCheck{Identity: id, By: by})
				}
			}
		}
	}
}

// identities returns the identity of each check of the subject records of the
// suite once, in the order of the records.
func identities(s *Suite) []string {
	var order []string
	for _, c := range verdicts(s.Records, spec.OnSubject) {
		if !slices.Contains(order, c.Identity) {
			order = append(order, c.Identity)
		}
	}
	return order
}

// subsumer returns the first check of the order in the suite that made every
// kill of the tally t and killed more mutants. The check of t killed no more
// mutants than itself, so it is never its own subsumer.
func subsumer(suite int, order []string, t *tally, tallies map[check]*tally) (string, bool) {
	for _, other := range order {
		o := tallies[check{identity: other, suite: suite}]
		if o == nil || len(o.killed) <= len(t.killed) {
			continue
		}
		if within(t.killed, o.killed) {
			return other, true
		}
	}
	return "", false
}

// within reports whether every mutant of a is a mutant of b.
func within(a, b map[string]bool) bool {
	for id := range a {
		if !b[id] {
			return false
		}
	}
	return true
}
