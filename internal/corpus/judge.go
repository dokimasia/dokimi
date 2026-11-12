// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package corpus

import (
	"maps"
	"slices"
	"strings"

	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/internal/report"
)

// The kinds of a key in the generation of a fixture. The status of a line of
// the lock is its kind, and the stance of an unclaimed classification is the
// status of its line.
const (
	kindCheck     = "check"
	kindAuthored  = "authored"
	kindCeiling   = "ceiling"
	kindRefused   = "refused"
	kindUnclaimed = "unclaimed"
	kindAbsent    = "absent"
)

// lockFields is the number of fields of a line of the lock: the package, the
// key, the claim, the status and the partners.
const lockFields = 5

// noPartners is the partners field of a line of the lock without partners.
const noPartners = "-"

// none is the actual outcome of a finding whose expected outcome is missing.
const none = "none"

// state is the state of one key in the generation of a fixture: its kind, the
// claim and the partners of a check, the gate of a refusal, and the stance of
// an unclaimed classification.
type state struct {
	kind     string
	claim    string
	partners string
	gate     spec.GateID
	stance   spec.StanceID
}

// String returns s in the words of a finding, such as "check of cancel" or
// "refused by failure".
func (s *state) String() string {
	switch s.kind {
	case kindCheck:
		if s.partners == noPartners {
			return "check of " + s.claim
		}
		return "check of " + s.claim + " with " + s.partners
	case kindRefused:
		if s.gate == "" {
			return kindRefused
		}
		return "refused by " + string(s.gate)
	case kindUnclaimed:
		return "unclaimed as " + string(s.stance)
	}
	return s.kind
}

// difference is one difference of the generation of a fixture from its
// expected outcome: the key, the expected state and the actual state. The
// key of an unclaimed classification is the key of its callable, a space and
// the classification, and a missing fixture has no key.
type difference struct {
	key       string
	want, got state
}

// declares reports whether e declares the difference d of fixture. An entry
// with the stance blocked and no key declares the missing fixture. An entry
// with the actual outcome absent and the stance blocked declares a missing
// key, such as a ceiling that the runtime of the language cannot check.
func (e *Entry) declares(fixture string, d *difference) bool {
	if e.Fixture != fixture || e.Key != d.key {
		return false
	}
	if d.key == "" {
		return e.Stance == spec.Blocked
	}
	blocked := d.got.kind == kindAbsent && e.Stance == spec.Blocked
	return e.Expect == d.want.kind && e.Actual == d.got.kind && e.Gate == d.got.gate &&
		(e.Stance == d.got.stance || blocked)
}

// Judge returns the findings of each fixture of expected whose level is at
// most the level of overlay, sorted by fixture and then by phase:
//
//   - generation: a key of the lock of the fixture that is missing, added or
//     changed, and a refusal with another gate in the records, unless an
//     entry of overlay declares the difference.
//   - falsification: a check of a falsification record that the planted
//     defect does not prove, and a derived check with another reason code,
//     another law, or no falsification record.
//   - subject: a subject without a record, and a check of a subject that
//     does not pass or that its excusal does not excuse.
//   - controls: a control without a record, and a check that fails on a
//     control.
//   - fakes: a fake without a record, and a check of a fake that does not
//     pass. A fake has no excusal.
//   - declaration: an entry of overlay that declares no difference.
//
// A falsification, subject, control or fake record whose field ended lists
// problems fails its phase. Judge also reports each outcome without an
// expectation, and each fixture of the level without an outcome that no
// entry declares as blocked. declared is the number of the entries of
// overlay that declare a difference.
func Judge(expected []Expectation, overlay Overlay, outcomes []Outcome) (found []Finding, declared int) {
	byFixture := make(map[string]*Outcome, len(outcomes))
	for i := range outcomes {
		byFixture[outcomes[i].Fixture] = &outcomes[i]
	}
	matched := make([]bool, len(overlay.Diverge))
	declare := func(fixture string, d *difference) bool {
		i := slices.IndexFunc(overlay.Diverge, func(e Entry) bool { return e.declares(fixture, d) })
		if i >= 0 {
			matched[i] = true
			declared++
		}
		return i >= 0
	}
	levels := make(map[string]int, len(expected))
	for i := range expected {
		e := &expected[i]
		levels[e.Fixture] = e.Level
		if e.Level > overlay.Level {
			continue
		}
		o, ok := byFixture[e.Fixture]
		if !ok {
			if !declare(e.Fixture, &difference{}) {
				found = append(found, Finding{
					Fixture: e.Fixture, Phase: PhaseGeneration, Want: "the fixture", Got: none,
				})
			}
			continue
		}
		differences := generation(e, o)
		changed := map[string]bool{}
		for j := range differences {
			d := &differences[j]
			changed[d.key] = true
			if !declare(e.Fixture, d) {
				found = append(found, Finding{
					Fixture: e.Fixture, Phase: PhaseGeneration, Key: d.key, Want: d.want.String(), Got: d.got.String(),
				})
			}
		}
		found = slices.Concat(found, falsification(e, o, changed), subjects(e, o, changed), controls(e, o),
			fakes(e, o, changed))
	}
	for i := range outcomes {
		if _, ok := levels[outcomes[i].Fixture]; !ok {
			found = append(found, Finding{
				Fixture: outcomes[i].Fixture, Phase: PhaseGeneration, Want: "a file of expected outcomes", Got: none,
			})
		}
	}
	for i := range overlay.Diverge {
		// A fixture without an expectation has the level 0, which every
		// overlay judges.
		e := &overlay.Diverge[i]
		if matched[i] || levels[e.Fixture] > overlay.Level {
			continue
		}
		actual := state{kind: e.Actual, gate: e.Gate, stance: e.Stance}
		want := "the declared difference from " + e.Expect + " to " + actual.String()
		if e.Key == "" {
			want = "a missing fixture that the entry declares as " + string(e.Stance)
		}
		found = append(found, Finding{Fixture: e.Fixture, Phase: PhaseDeclaration, Key: e.Key, Want: want, Got: none})
	}
	// The findings of one fixture are in the order of the phases.
	slices.SortStableFunc(found, func(a, b Finding) int { return strings.Compare(a.Fixture, b.Fixture) })
	return found, declared
}

// generation returns the differences of the lock of o from the declarations
// of e, sorted by key. A refusal has the gate of its refusal in the records
// of o. An expected key that the lock lacks has the state of an unclaimed
// classification of its callable that e does not expect, and that
// classification has no difference of its own.
func generation(e *Expectation, o *Outcome) []difference {
	wantKeys, wantUnclaimed := map[string]state{}, map[string]state{}
	for _, d := range e.Declarations {
		for key, c := range d.Checks {
			wantKeys[key] = state{kind: kindCheck, claim: string(c.Claim), partners: partnersField(c.Partners)}
		}
		for _, key := range d.Authored {
			wantKeys[key] = state{kind: kindAuthored}
		}
		for _, key := range d.Ceilings {
			wantKeys[key] = state{kind: kindCeiling}
		}
		for key, gate := range d.Refusals {
			wantKeys[key] = state{kind: kindRefused, gate: gate}
		}
		for _, u := range d.Unclaimed {
			wantUnclaimed[u.Callable+" "+u.Classification] = state{kind: kindUnclaimed, stance: u.Stance}
		}
	}
	gates := gatesOf(o.Records)
	gotKeys, gotUnclaimed := map[string]state{}, map[string]state{}
	for _, text := range o.Lock {
		f := fieldsOf(text)
		key, claim, status := f[1], f[2], f[3]
		switch status {
		case kindCheck:
			gotKeys[key] = state{kind: kindCheck, claim: claim, partners: f[4]}
		case kindAuthored, kindCeiling:
			gotKeys[key] = state{kind: status}
		case kindRefused:
			gotKeys[key] = state{kind: kindRefused, gate: gates[refusalKey(key, claim)]}
		default:
			gotUnclaimed[key+" "+claim] = state{kind: kindUnclaimed, stance: spec.StanceID(status)}
		}
	}
	explained := map[string]bool{}
	unclaimedOf := func(key string) state {
		callable, _, _ := strings.Cut(key, "/")
		for _, u := range slices.Sorted(maps.Keys(gotUnclaimed)) {
			if _, expected := wantUnclaimed[u]; strings.HasPrefix(u, callable+" ") && !expected && !explained[u] {
				explained[u] = true
				return gotUnclaimed[u]
			}
		}
		return state{kind: kindAbsent}
	}
	var out []difference
	for _, key := range unionOf(wantKeys, gotKeys) {
		// A key of the union is in one of the maps at least, so an absent
		// actual state has an expected one.
		want, got := stateOf(wantKeys, key), stateOf(gotKeys, key)
		if got.kind == kindAbsent {
			got = unclaimedOf(key)
		}
		if want != got {
			out = append(out, difference{key: key, want: want, got: got})
		}
	}
	for _, key := range unionOf(wantUnclaimed, gotUnclaimed) {
		if want, got := stateOf(wantUnclaimed, key), stateOf(gotUnclaimed, key); want != got && !explained[key] {
			out = append(out, difference{key: key, want: want, got: got})
		}
	}
	return out
}

// stateOf returns the state of key in states, and the state absent for a
// key that states lacks.
func stateOf(states map[string]state, key string) state {
	if s, ok := states[key]; ok {
		return s
	}
	return state{kind: kindAbsent}
}

// unionOf returns the keys of a and of b, sorted, each once.
func unionOf(a, b map[string]state) []string {
	keys := slices.AppendSeq(slices.Collect(maps.Keys(a)), maps.Keys(b))
	slices.Sort(keys)
	return slices.Compact(keys)
}

// falsification returns the findings of the falsification records of o. A
// check that a falsification excuses or skips has no result to judge. A
// derived check without a falsification record is a finding unless its key
// is changed, which the phase generation judges.
func falsification(e *Expectation, o *Outcome, changed map[string]bool) []Finding {
	checks := checksOf(e)
	judged := map[string]bool{}
	var found []Finding
	for i := range o.Records {
		r := &o.Records[i]
		if r.Application != spec.Falsification {
			continue
		}
		if r.Ended != "" {
			found = append(found, ended(e.Fixture, PhaseFalsification, r))
			continue
		}
		for j := range r.Checks {
			c := &r.Checks[j]
			want, expected := checks[c.Key]
			finding := Finding{Fixture: e.Fixture, Phase: PhaseFalsification, Key: c.Key, Got: resultOf(c)}
			switch {
			case c.Result == spec.Excused || c.Result == spec.Skipped:
				continue
			case !expected && c.Result != spec.Proven:
				finding.Want = string(spec.Proven)
			case expected && (c.Result != spec.Proven || c.Reason != want.Reason):
				finding.Want = string(spec.Proven) + " by " + string(want.Reason)
			case expected && c.Law != want.Law:
				finding.Want, finding.Got = "the law "+string(want.Law), "the law "+string(c.Law)
			default:
				judged[c.Key] = true
				continue
			}
			judged[c.Key] = true
			found = append(found, finding)
		}
	}
	for _, key := range slices.Sorted(maps.Keys(checks)) {
		if !judged[key] && !changed[key] {
			found = append(found, Finding{
				Fixture: e.Fixture, Phase: PhaseFalsification, Key: key,
				Want: string(spec.Proven) + " by " + string(checks[key].Reason), Got: none,
			})
		}
	}
	return found
}

// subjects returns the findings of the subject records of each subject of e
// in o, as [applied] judges them.
func subjects(e *Expectation, o *Outcome, changed map[string]bool) []Finding {
	var found []Finding
	for _, s := range e.Subjects {
		found = slices.Concat(found, applied(e, o, changed, PhaseSubject, &s))
	}
	return found
}

// fakes returns the findings of the subject records of each fake of e in o,
// as [applied] judges them for a subject without an excusal.
func fakes(e *Expectation, o *Outcome, changed map[string]bool) []Finding {
	var found []Finding
	for _, fake := range e.Fakes {
		found = slices.Concat(found, applied(e, o, changed, PhaseFakes, &Subject{Type: fake}))
	}
	return found
}

// applied returns the findings of phase of the subject records of s in o.
// Each derived and each authored check of the declarations of the records of
// s has a record on s, unless its key is changed, which the phase generation
// judges. A subject without a check builds no value, so a subject of a
// fixture without a check needs no record. A ceiling of latency runs in a
// benchmark, which go test does not run, so the phase judges a ceiling only
// where it has a record.
func applied(e *Expectation, o *Outcome, changed map[string]bool, phase Phase, s *Subject) []Finding {
	var found []Finding
	applies := appliedKeys(e)
	records := applications(o, s.Type)
	if len(records) == 0 && len(applies) > 0 {
		return []Finding{{Fixture: e.Fixture, Phase: phase, Key: s.Type, Want: "a subject record", Got: none}}
	}
	recorded, declarations := map[string]bool{}, map[string]bool{}
	for _, r := range records {
		declarations[spec.Key(r.Suite.Declaration)] = true
		if r.Ended != "" {
			found = append(found, ended(e.Fixture, phase, r))
			continue
		}
		for i := range r.Checks {
			c := &r.Checks[i]
			recorded[c.Key] = true
			// The excusal of a check also excuses its generated binding.
			excused := slices.ContainsFunc(s.Excused, func(key string) bool {
				return key == c.Key || derive.GeneratedOf(key) == c.Key
			})
			want := spec.Pass
			if excused {
				want = spec.Excused
			}
			if c.Result != want {
				found = append(found, Finding{
					Fixture: e.Fixture, Phase: phase, Key: c.Key, Want: string(want),
					Got: resultOf(c) + " on " + s.Type,
				})
			}
		}
	}
	for _, key := range applies {
		if declarations[declarationOf(key)] && !recorded[key] && !changed[key] {
			found = append(found, Finding{
				Fixture: e.Fixture, Phase: phase, Key: key, Want: string(spec.Pass), Got: "no record on " + s.Type,
			})
		}
	}
	return found
}

// controls returns the findings of the subject records of each control of e
// in o.
func controls(e *Expectation, o *Outcome) []Finding {
	var found []Finding
	for _, control := range e.Controls {
		records := applications(o, control)
		if len(records) == 0 {
			found = append(found, Finding{
				Fixture: e.Fixture, Phase: PhaseControls, Key: control, Want: "a subject record", Got: none,
			})
		}
		for _, r := range records {
			if r.Ended != "" {
				found = append(found, ended(e.Fixture, PhaseControls, r))
				continue
			}
			for i := range r.Checks {
				if c := &r.Checks[i]; c.Result == spec.Fail {
					found = append(found, Finding{
						Fixture: e.Fixture, Phase: PhaseControls, Key: c.Key, Want: string(spec.Pass),
						Got: resultOf(c) + " on " + control,
					})
				}
			}
		}
	}
	return found
}

// applications returns the subject records of o whose subject is typ. A
// falsification record has no subject.
func applications(o *Outcome, typ string) []*report.Record {
	var out []*report.Record
	for i := range o.Records {
		if r := &o.Records[i]; r.Subject == typ {
			out = append(out, r)
		}
	}
	return out
}

// ended returns the finding of phase of the record r, whose application ended
// before its checks.
func ended(fixture string, phase Phase, r *report.Record) Finding {
	return Finding{
		Fixture: fixture, Phase: phase, Key: spec.Key(r.Suite.Declaration), Want: "every check",
		Got: "ended in " + r.Test + ": " + strings.ReplaceAll(r.Ended, "\n", "; "),
	}
}

// resultOf returns the result of c in the words of a finding, with the reason
// code that its probe reported, the other check of its result, and why.
func resultOf(c *spec.CheckRecord) string {
	s := string(c.Result)
	if c.Reason != "" {
		s += " by " + string(c.Reason)
	}
	if c.Other != "" {
		s += " beside " + c.Other
	}
	if c.Why != "" {
		s += ": " + c.Why
	}
	return s
}

// checksOf returns the derived checks of every declaration of e by their
// keys.
func checksOf(e *Expectation) map[string]Check {
	checks := map[string]Check{}
	for _, d := range e.Declarations {
		maps.Copy(checks, d.Checks)
	}
	return checks
}

// appliedKeys returns the keys of the derived and the authored checks of
// every declaration of e, sorted.
func appliedKeys(e *Expectation) []string {
	var keys []string
	for _, d := range e.Declarations {
		keys = slices.Concat(keys, slices.Collect(maps.Keys(d.Checks)), d.Authored)
	}
	slices.Sort(keys)
	return keys
}

// gatesOf returns the gate of each refusal of records by [refusalKey]. Each
// record of a suite has the refusals of the suite.
func gatesOf(records []report.Record) map[string]spec.GateID {
	gates := map[string]spec.GateID{}
	for i := range records {
		for _, r := range records[i].Refusals {
			gates[spec.Key(r.Callable)+"/"+string(r.Claim)] = r.Gate
		}
	}
	return gates
}

// refusalKey returns the key of the callable of the check key, a slash and
// the claim. A record names a refusal by its callable and its claim.
func refusalKey(key, claim string) string {
	callable, _, _ := strings.Cut(key, "/")
	return callable + "/" + claim
}

// fieldsOf returns the fields of a line of the lock. A line with fewer fields
// has the empty string in the place of each field that it lacks.
func fieldsOf(text string) [lockFields]string {
	var f [lockFields]string
	copy(f[:], strings.Split(text, "\t"))
	return f
}

// declarationOf returns the key of the declaration of a key of the lock: the
// key up to its first dot or slash, such as mixed for mixed.read/cancel.
func declarationOf(key string) string {
	callable, _, _ := strings.Cut(key, "/")
	declaration, _, _ := strings.Cut(callable, ".")
	return declaration
}

// partnersField returns the partners field of a line of the lock: each
// partner as its name and the key of its callable, in the order of the names,
// or a dash without partners.
func partnersField(partners map[string]string) string {
	if len(partners) == 0 {
		return noPartners
	}
	parts := make([]string, 0, len(partners))
	for _, name := range slices.Sorted(maps.Keys(partners)) {
		parts = append(parts, name+"="+partners[name])
	}
	return strings.Join(parts, ",")
}
