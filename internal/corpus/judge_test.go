// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package corpus_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/corpus"
)

// The positions of the lines of the lock of [outcome].
const (
	lineCancel  = 1
	lineExpires = 2
)

// The positions of the records of [outcome].
const (
	falsified = 0
	onSubject = 1
	onControl = 2
)

// refuseCancel turns the check mixed.read/cancel of o into a refusal by the
// gate in its lock and in its records.
func refuseCancel(o *corpus.Outcome, gate spec.GateID) {
	o.Lock[lineCancel] = ttl + "\tmixed.read/cancel\tcancel\trefused\t-"
	for i := range o.Records {
		o.Records[i].Refusals = append(o.Records[i].Refusals, spec.RefusalRecord{
			Claim: spec.Cancel, Callable: "Mixed.Read", Gate: gate, Why: "Read takes no context",
			Remedy: "add a context parameter to Read",
		})
	}
}

// leaveCancelUnclaimed turns the check mixed.read/cancel of o into an
// unclaimed classification of mixed.read with the stance gap.
func leaveCancelUnclaimed(o *corpus.Outcome) {
	o.Lock[lineCancel] = ttl + "\tmixed.read\tmixin:cancellable\tgap\t-"
}

func TestJudge(t *testing.T) {
	t.Parallel()

	t.Run("Judge", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name         string
			give         func(e *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay)
			want         []corpus.Finding
			wantDeclared int
		}{
			{
				name: "returns no finding for an outcome that meets its expectation",
				give: func(*corpus.Expectation, *corpus.Outcome, *corpus.Overlay) {},
			},
			{
				name: "reports a check that the lock lacks",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Lock = slices.Delete(o.Lock, lineCancel, lineCancel+1)
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
					Want: "check of cancel", Got: "absent",
				}},
			},
			{
				name: "reports a check that the lock adds",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Lock = append(o.Lock, ttl+"\tmixed.get/cancel\tcancel\tcheck\t-")
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.get/cancel",
					Want: "absent", Got: "check of cancel",
				}},
			},
			{
				name: "reports a check with another partner",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Lock[lineExpires] = ttl + "\tmixed.read/expires\texpires\tcheck\twrite=mixed.set"
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/expires",
					Want: "check of expires with write=mixed.put", Got: "check of expires with write=mixed.set",
				}},
			},
			{
				name: "reports a refusal with another gate in the records",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					for i := range o.Records {
						o.Records[i].Refusals[0].Gate = spec.Failure
					}
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.scan/zero-on-error",
					Want: "refused by input-samples", Got: "refused by failure",
				}},
			},
			{
				name: "reports a refusal without a record of its gate",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					for i := range o.Records {
						o.Records[i].Refusals = nil
					}
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.scan/zero-on-error",
					Want: "refused by input-samples", Got: "refused",
				}},
			},
			{
				name: "reports an unclaimed classification with another stance",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Lock[0] = ttl + "\tmixed.put\tshape:writer\tgap\t-"
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.put shape:writer",
					Want: "unclaimed as silent", Got: "unclaimed as gap",
				}},
			},
			{
				name: "reports an unclaimed classification that the expectation lacks",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Lock = append(o.Lock, ttl+"\tmixed.scan\tshape:reader\tsilent\t-")
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.scan shape:reader",
					Want: "absent", Got: "unclaimed as silent",
				}},
			},
			{
				name: "reports a check that the lock leaves unclaimed as one difference",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					leaveCancelUnclaimed(o)
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
					Want: "check of cancel", Got: "unclaimed as gap",
				}},
			},
			{
				name: "declares a refusal that an entry of the overlay states",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					refuseCancel(o, spec.Cancellation)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "check", Actual: "refused",
						Gate: spec.Cancellation, Why: "a call of the language takes no context",
					}}
				},
				wantDeclared: 1,
			},
			{
				name: "declares an unclaimed check that an entry of the overlay states",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					leaveCancelUnclaimed(o)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "check", Actual: "unclaimed",
						Stance: spec.Gap, Why: "the runtime of the language cannot cancel a call",
					}}
				},
				wantDeclared: 1,
			},
			{
				name: "reports an entry whose outcome meets the expectation",
				give: func(_ *corpus.Expectation, _ *corpus.Outcome, overlay *corpus.Overlay) {
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "check", Actual: "refused",
						Gate: spec.Cancellation, Why: "a call of the language takes no context",
					}}
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseDeclaration, Key: "mixed.read/cancel",
					Want: "the declared difference from check to refused by cancellation", Got: "none",
				}},
			},
			{
				name: "reports an entry that declares another difference than the run",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					refuseCancel(o, spec.Failure)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "check", Actual: "refused",
						Gate: spec.Cancellation, Why: "a call of the language takes no context",
					}}
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "refused by failure",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseDeclaration, Key: "mixed.read/cancel",
						Want: "the declared difference from check to refused by cancellation", Got: "none",
					},
				},
			},
			{
				name: "declares a refusal of a check without a record",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					refuseCancel(o, spec.Cancellation)
					for i := range o.Records {
						o.Records[i].Checks = slices.DeleteFunc(o.Records[i].Checks, func(c spec.CheckRecord) bool {
							return c.Key == "mixed.read/cancel"
						})
					}
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "check", Actual: "refused",
						Gate: spec.Cancellation, Why: "a call of the language takes no context",
					}}
				},
				wantDeclared: 1,
			},
			{
				name: "reports a difference whose entry has another expected outcome",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					refuseCancel(o, spec.Cancellation)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "authored", Actual: "refused",
						Gate: spec.Cancellation, Why: "a call of the language takes no context",
					}}
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "refused by cancellation",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseDeclaration, Key: "mixed.read/cancel",
						Want: "the declared difference from authored to refused by cancellation", Got: "none",
					},
				},
			},
			{
				name: "reports a difference whose entry has another actual outcome",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					refuseCancel(o, spec.Cancellation)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "check", Actual: "absent",
						Gate: spec.Cancellation, Why: "a call of the language takes no context",
					}}
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "refused by cancellation",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseDeclaration, Key: "mixed.read/cancel",
						Want: "the declared difference from check to absent", Got: "none",
					},
				},
			},
			{
				name: "reports a difference whose entry has another stance",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					leaveCancelUnclaimed(o)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/cancel", Expect: "check", Actual: "unclaimed",
						Stance: spec.Needs, Why: "the runtime of the language cannot cancel a call",
					}}
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "unclaimed as gap",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseDeclaration, Key: "mixed.read/cancel",
						Want: "the declared difference from check to unclaimed as needs", Got: "none",
					},
				},
			},
			{
				name: "judges a fixture at the level of the overlay",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					overlay.Level = 2
					o.Lock = slices.Delete(o.Lock, lineCancel, lineCancel+1)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/expires", Expect: "check", Actual: "absent",
						Why: "the language has no clock",
					}}
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "absent",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseDeclaration, Key: "mixed.read/expires",
						Want: "the declared difference from check to absent", Got: "none",
					},
				},
			},
			{
				name: "reports a check absent beside a classification of another callable",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Lock[lineCancel] = ttl + "\tmixed.put\tmixin:cancellable\tgap\t-"
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "absent",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.put mixin:cancellable",
						Want: "absent", Got: "unclaimed as gap",
					},
				},
			},
			{
				name: "reports a check absent beside an expected classification of its callable",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					mixed := e.Declarations["mixed"]
					mixed.Unclaimed = append(mixed.Unclaimed, corpus.Unclaimed{
						Classification: "mixin:sorted", Callable: "mixed.read", Stance: spec.Silent,
					})
					e.Declarations["mixed"] = mixed
					o.Lock[lineCancel] = ttl + "\tmixed.read\tmixin:sorted\tsilent\t-"
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
					Want: "check of cancel", Got: "absent",
				}},
			},
			{
				name: "explains one absent check with one classification of its callable",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					leaveCancelUnclaimed(o)
					o.Lock = slices.Delete(o.Lock, lineExpires, lineExpires+1)
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "unclaimed as gap",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/expires",
						Want: "check of expires with write=mixed.put", Got: "absent",
					},
				},
			},
			{
				name: "accepts a proven check that the expectation does not list",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[falsified].Checks = append(o.Records[falsified].Checks, spec.CheckRecord{
						Key: "mixed.get/cancel", Claim: spec.Cancel, Law: spec.CancelledContextIsReported,
						Result: spec.Proven, Reason: spec.IgnoredCancel,
					})
				},
			},
			{
				name: "judges the checks of the declarations of the records of a subject",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Declarations["other"] = corpus.Declaration{Checks: map[string]corpus.Check{
						"other.get/cancel": {
							Claim: spec.Cancel, Law: spec.CancelledContextIsReported, Reason: spec.IgnoredCancel,
						},
					}}
					o.Lock = append(o.Lock, ttl+"\tother.get/cancel\tcancel\tcheck\t-")
					o.Records[falsified].Checks = append(o.Records[falsified].Checks, spec.CheckRecord{
						Key: "other.get/cancel", Claim: spec.Cancel, Law: spec.CancelledContextIsReported,
						Result: spec.Proven, Reason: spec.IgnoredCancel,
					})
				},
			},
			{
				name: "reports a difference whose entry has another key",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					refuseCancel(o, spec.Cancellation)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/expires", Expect: "check", Actual: "refused",
						Gate: spec.Cancellation, Why: "a call of the language takes no context",
					}}
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "refused by cancellation",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseDeclaration, Key: "mixed.read/expires",
						Want: "the declared difference from check to refused by cancellation", Got: "none",
					},
				},
			},
			{
				name: "reports a difference whose entry has another fixture",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					refuseCancel(o, spec.Cancellation)
					overlay.Diverge = []corpus.Entry{{
						Fixture: "iface/mixin/ttlperwrite", Key: "mixed.read/cancel",
						Expect: "check", Actual: "refused", Gate: spec.Cancellation,
						Why: "a call of the language takes no context",
					}}
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseGeneration, Key: "mixed.read/cancel",
						Want: "check of cancel", Got: "refused by cancellation",
					},
					{
						Fixture: "iface/mixin/ttlperwrite", Phase: corpus.PhaseDeclaration, Key: "mixed.read/cancel",
						Want: "the declared difference from check to refused by cancellation", Got: "none",
					},
				},
			},
			{
				name: "reports an entry of a fixture without an expectation",
				give: func(_ *corpus.Expectation, _ *corpus.Outcome, overlay *corpus.Overlay) {
					overlay.Diverge = []corpus.Entry{{
						Fixture: "enum/priority", Key: "priority/parses", Expect: "check", Actual: "absent",
						Why: "the language has no parse function",
					}}
				},
				want: []corpus.Finding{{
					Fixture: "enum/priority", Phase: corpus.PhaseDeclaration, Key: "priority/parses",
					Want: "the declared difference from check to absent", Got: "none",
				}},
			},
			{
				name: "skips a fixture above the level of the overlay",
				give: func(e *corpus.Expectation, o *corpus.Outcome, overlay *corpus.Overlay) {
					e.Level, overlay.Level = 3, 2
					o.Lock = slices.Delete(o.Lock, lineCancel, lineCancel+1)
					overlay.Diverge = []corpus.Entry{{
						Fixture: ttl, Key: "mixed.read/expires", Expect: "check", Actual: "absent",
						Why: "the language has no clock",
					}}
				},
			},
			{
				name: "reports a check that its planted defect proves with another reason code",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					c := &o.Records[falsified].Checks[0]
					c.Result, c.Reason = spec.WrongReason, spec.ValueBesideError
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/cancel",
					Want: "proven by ignored-cancel", Got: "wrong-reason by value-beside-error",
				}},
			},
			{
				name: "reports a proven check with another reason code",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[falsified].Checks[0].Reason = spec.WrongSentinel
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/cancel",
					Want: "proven by ignored-cancel", Got: "proven by wrong-sentinel",
				}},
			},
			{
				name: "reports a proven check with another law",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[falsified].Checks[0].Law = spec.NoValueBesideAnError
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/cancel",
					Want: "the law cancelled-context-is-reported", Got: "the law no-value-beside-an-error",
				}},
			},
			{
				name: "reports a check whose defect fails another check",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					c := &o.Records[falsified].Checks[0]
					c.Result, c.Other = spec.NotIsolated, "Mixed.Read/expires"
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/cancel",
					Want: "proven by ignored-cancel", Got: "not-isolated by ignored-cancel beside Mixed.Read/expires",
				}},
			},
			{
				name: "reports a derived check without a falsification",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[falsified].Checks = o.Records[falsified].Checks[:1]
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/expires",
					Want: "proven by not-expired", Got: "none",
				}},
			},
			{
				name: "reports another check whose defect passes it",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[falsified].Checks = append(o.Records[falsified].Checks, spec.CheckRecord{
						Key: "mixed.read/latency-budget", Result: spec.Passed,
					})
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/latency-budget",
					Want: "proven", Got: "passed",
				}},
			},
			{
				name: "skips a check that a falsification excuses or skips",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records = append(o.Records, record(spec.Falsification, "",
						spec.CheckRecord{Key: "mixed.read/cancel", Result: spec.Excused, Why: "the test excuses it"},
						spec.CheckRecord{Key: "mixed.read/expires", Result: spec.Skipped, Why: "the clock stops"},
					))
				},
			},
			{
				name: "reports a falsification that ended before its checks",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[falsified].Ended = "DOKIMI_SEED is not a number\nthe option Cases asks for 0 cases"
					o.Records[falsified].Checks = nil
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed", Want: "every check",
						Got: "ended in TestMixedSuite: DOKIMI_SEED is not a number; the option Cases asks for 0 cases",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/cancel",
						Want: "proven by ignored-cancel", Got: "none",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseFalsification, Key: "mixed.read/expires",
						Want: "proven by not-expired", Got: "none",
					},
				},
			},
			{
				name: "reports a subject without a record",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records = slices.Delete(o.Records, onSubject, onSubject+1)
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseSubject, Key: subject, Want: "a subject record", Got: "none",
				}},
			},
			{
				name: "reports a check that fails on the subject",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[onSubject].Checks[0].Result = spec.Fail
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed.read/cancel",
					Want: "pass", Got: "fail on " + subject,
				}},
			},
			{
				name: "reports a check that the subject excuses without an excusal of the fixture",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					c := &o.Records[onSubject].Checks[2]
					c.Result, c.Why = spec.Excused, "the subject keeps no item"
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed.read/keeps-the-item",
					Want: "pass", Got: "excused: the subject keeps no item on " + subject,
				}},
			},
			{
				name: "accepts a check that the excusal of the subject excuses",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Subjects[0].Excused = []string{"mixed.read/keeps-the-item"}
					c := &o.Records[onSubject].Checks[2]
					c.Result, c.Why = spec.Excused, "the subject keeps no item"
				},
			},
			{
				name: "accepts the generated binding of a check that the excusal of the subject excuses",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Subjects[0].Excused = []string{"mixed.read/cancel"}
					r := &o.Records[onSubject]
					r.Checks[0].Result, r.Checks[0].Why = spec.Excused, "the subject cannot cancel"
					r.Checks = append(r.Checks, spec.CheckRecord{
						Key: "mixed.read/cancel#generated", Result: spec.Excused, Why: "the subject cannot cancel",
					})
				},
			},
			{
				name: "reports an excusal of a check that passes on the subject",
				give: func(e *corpus.Expectation, _ *corpus.Outcome, _ *corpus.Overlay) {
					e.Subjects[0].Excused = []string{"mixed.read/keeps-the-item"}
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed.read/keeps-the-item",
					Want: "excused", Got: "pass on " + subject,
				}},
			},
			{
				name: "reports an authored check without a record on the subject",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[onSubject].Checks = slices.Delete(o.Records[onSubject].Checks, 2, 3)
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed.read/keeps-the-item",
					Want: "pass", Got: "no record on " + subject,
				}},
			},
			{
				name: "accepts a ceiling without a record on the subject",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[onSubject].Checks = o.Records[onSubject].Checks[:3]
				},
			},
			{
				name: "accepts a subject without a record of a fixture without a check",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					mixed := e.Declarations["mixed"]
					mixed.Checks, mixed.Authored = map[string]corpus.Check{}, nil
					e.Declarations["mixed"] = mixed
					o.Lock = slices.Delete(o.Lock, lineCancel, lineExpires+2)
					o.Records = o.Records[onControl:]
				},
			},
			{
				name: "reports an application to the subject that ended before its checks",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[onSubject].Ended = "the check Mixed.Read/cancel has no statement"
					o.Records[onSubject].Checks = nil
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed", Want: "every check",
						Got: "ended in TestMixedSuite: the check Mixed.Read/cancel has no statement",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed.read/cancel",
						Want: "pass", Got: "no record on " + subject,
					},
					{
						Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed.read/expires",
						Want: "pass", Got: "no record on " + subject,
					},
					{
						Fixture: ttl, Phase: corpus.PhaseSubject, Key: "mixed.read/keeps-the-item",
						Want: "pass", Got: "no record on " + subject,
					},
				},
			},
			{
				name: "reports a control without a record",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records = o.Records[:onControl]
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseControls, Key: control, Want: "a subject record", Got: "none",
				}},
			},
			{
				name: "reports a check that fails on a control",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[onControl].Checks[0].Result = spec.Fail
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseControls, Key: "mixed.read/cancel",
					Want: "pass", Got: "fail on " + control,
				}},
			},
			{
				name: "accepts a check that is vacuous on a control",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[onControl].Checks[0].Result = spec.Vacuous
				},
			},
			{
				name: "reports an application to a control that ended before its checks",
				give: func(_ *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					o.Records[onControl].Ended = "the seed is not a number"
					o.Records[onControl].Checks = nil
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseControls, Key: "mixed", Want: "every check",
					Got: "ended in TestMixedSuite: the seed is not a number",
				}},
			},
			{
				name: "accepts a fake whose record passes each check of the declaration",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Fakes = []string{fake}
					o.Records = append(o.Records, record(spec.OnSubject, fake, o.Records[onSubject].Checks...))
				},
			},
			{
				name: "reports a fake without a record",
				give: func(e *corpus.Expectation, _ *corpus.Outcome, _ *corpus.Overlay) {
					e.Fakes = []string{fake}
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFakes, Key: fake, Want: "a subject record", Got: "none",
				}},
			},
			{
				name: "reports a check that is vacuous on a fake",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Fakes = []string{fake}
					r := record(spec.OnSubject, fake, o.Records[onSubject].Checks...)
					r.Checks[0].Result = spec.Vacuous
					o.Records = append(o.Records, r)
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFakes, Key: "mixed.read/cancel",
					Want: "pass", Got: "vacuous on " + fake,
				}},
			},
			{
				name: "reports a check that a fake excuses, because a fake has no excusal",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Fakes = []string{fake}
					e.Subjects[0].Excused = []string{"mixed.read/keeps-the-item"}
					r := record(spec.OnSubject, fake, o.Records[onSubject].Checks...)
					r.Checks[2].Result, r.Checks[2].Why = spec.Excused, "the fake keeps no item"
					o.Records = append(o.Records, r)
					c := &o.Records[onSubject].Checks[2]
					c.Result, c.Why = spec.Excused, "the subject keeps no item"
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFakes, Key: "mixed.read/keeps-the-item",
					Want: "pass", Got: "excused: the fake keeps no item on " + fake,
				}},
			},
			{
				name: "reports a check of the declaration without a record on a fake",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Fakes = []string{fake}
					o.Records = append(o.Records, record(spec.OnSubject, fake, o.Records[onSubject].Checks[1:]...))
				},
				want: []corpus.Finding{{
					Fixture: ttl, Phase: corpus.PhaseFakes, Key: "mixed.read/cancel",
					Want: "pass", Got: "no record on " + fake,
				}},
			},
			{
				name: "reports an application to a fake that ended before its checks",
				give: func(e *corpus.Expectation, o *corpus.Outcome, _ *corpus.Overlay) {
					e.Fakes = []string{fake}
					r := record(spec.OnSubject, fake, o.Records[onSubject].Checks...)
					r.Ended, r.Checks = "the seed is not a number", nil
					o.Records = append(o.Records, r)
				},
				want: []corpus.Finding{
					{
						Fixture: ttl, Phase: corpus.PhaseFakes, Key: "mixed", Want: "every check",
						Got: "ended in TestMixedSuite: the seed is not a number",
					},
					{
						Fixture: ttl, Phase: corpus.PhaseFakes, Key: "mixed.read/cancel",
						Want: "pass", Got: "no record on " + fake,
					},
					{
						Fixture: ttl, Phase: corpus.PhaseFakes, Key: "mixed.read/expires",
						Want: "pass", Got: "no record on " + fake,
					},
					{
						Fixture: ttl, Phase: corpus.PhaseFakes, Key: "mixed.read/keeps-the-item",
						Want: "pass", Got: "no record on " + fake,
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				e, o, overlay := expectation(), outcome(), corpus.Overlay{Language: "go", Level: 4}
				tt.give(&e, &o, &overlay)

				got, declared := corpus.Judge([]corpus.Expectation{e}, overlay, []corpus.Outcome{o})

				assert.Equal(t, got, tt.want, "the findings are the differences that the overlay does not declare")
				assert.Equal(t, declared, tt.wantDeclared, "each entry that matches a difference declares it")
			})
		}

		t.Run("reports a fixture of the level without an outcome", func(t *testing.T) {
			t.Parallel()

			got, _ := corpus.Judge([]corpus.Expectation{expectation()}, corpus.Overlay{Level: 4}, nil)

			assert.Equal(t, got, []corpus.Finding{{
				Fixture: ttl, Phase: corpus.PhaseGeneration, Want: "the fixture", Got: "none",
			}}, "a fixture that the corpus lacks fails the phase generation")
		})

		t.Run("declares a fixture without an outcome that an entry blocks", func(t *testing.T) {
			t.Parallel()
			overlay := corpus.Overlay{Level: 4, Diverge: []corpus.Entry{{
				Fixture: ttl, Stance: spec.Blocked, Why: "the language has no clock",
			}}}

			got, declared := corpus.Judge([]corpus.Expectation{expectation()}, overlay, nil)

			assert.Empty(t, got, "the blocked entry declares the missing fixture")
			assert.Equal(t, declared, 1, "the blocked entry declares one difference")
		})

		t.Run("reports a fixture without an outcome whose entry is not blocked", func(t *testing.T) {
			t.Parallel()
			overlay := corpus.Overlay{Level: 4, Diverge: []corpus.Entry{{
				Fixture: ttl, Stance: spec.Gap, Why: "the language has no clock",
			}}}

			got, declared := corpus.Judge([]corpus.Expectation{expectation()}, overlay, nil)

			assert.Equal(t, got, []corpus.Finding{
				{Fixture: ttl, Phase: corpus.PhaseGeneration, Want: "the fixture", Got: "none"},
				{
					Fixture: ttl, Phase: corpus.PhaseDeclaration,
					Want: "a missing fixture that the entry declares as gap", Got: "none",
				},
			}, "only the stance blocked declares a missing fixture")
			assert.Equal(t, declared, 0, "the entry declares no difference")
		})

		t.Run("reports an outcome without an expectation before a later fixture", func(t *testing.T) {
			t.Parallel()
			later := expectation()
			later.Fixture = "struct/plain"
			stray := outcome()
			stray.Fixture = "enum/priority"

			got, _ := corpus.Judge([]corpus.Expectation{later}, corpus.Overlay{Level: 4}, []corpus.Outcome{stray})

			assert.Equal(t, got, []corpus.Finding{
				{
					Fixture: "enum/priority", Phase: corpus.PhaseGeneration,
					Want: "a file of expected outcomes", Got: "none",
				},
				{Fixture: "struct/plain", Phase: corpus.PhaseGeneration, Want: "the fixture", Got: "none"},
			}, "the findings are in the order of their fixtures")
		})
	})
}
