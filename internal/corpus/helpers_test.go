// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package corpus_test

import (
	"slices"

	"go.dokimi.dev/dokimi/internal/corpus"
	"go.dokimi.dev/dokimi/internal/report"
	"go.dokimi.dev/dokimi/spec"
)

// The fixture of the specs, its subject, its control and its fake.
const (
	ttl     = "iface/mixin/ttl"
	subject = "iface/mixin/ttl/ttltest.InMemory"
	control = "iface/mixin/ttl/deadline.Store"
	fake    = "iface/mixin/ttl/ttltest.MixedFake"
)

// expectation returns the expected outcome of the fixture ttl: two derived
// checks, an authored check, a ceiling, a refusal and an unclaimed
// classification of the declaration Mixed.
func expectation() corpus.Expectation {
	return corpus.Expectation{
		Fixture:  ttl,
		Level:    2,
		Subjects: []corpus.Subject{{Type: subject}},
		Controls: []string{control},
		Declarations: map[string]corpus.Declaration{"mixed": {
			Checks: map[string]corpus.Check{
				"mixed.read/cancel": {
					Claim: spec.Cancel, Law: spec.CancelledContextIsReported, Reason: spec.IgnoredCancel,
				},
				"mixed.read/expires": {
					Claim: spec.Expires, Law: spec.EntryLapsesOnTheClock, Reason: spec.NotExpired,
					Partners: map[string]string{"write": "mixed.put"},
				},
			},
			Authored:  []string{"mixed.read/keeps-the-item"},
			Ceilings:  []string{"mixed.read/latency-budget"},
			Refusals:  map[string]spec.GateID{"mixed.scan/zero-on-error": spec.InputSamples},
			Unclaimed: []corpus.Unclaimed{{Classification: "shape:writer", Callable: "mixed.put", Stance: spec.Silent}},
		}},
	}
}

// outcome returns the outcome of the fixture ttl that meets [expectation]:
// the lines of its lock, its falsification record, and a record of each
// application to the subject and to the control.
func outcome() corpus.Outcome {
	proven := []spec.CheckRecord{
		{
			Key: "mixed.read/cancel", Claim: spec.Cancel, Law: spec.CancelledContextIsReported,
			Result: spec.Proven, Reason: spec.IgnoredCancel,
		},
		{
			Key: "mixed.read/expires", Claim: spec.Expires, Law: spec.EntryLapsesOnTheClock,
			Result: spec.Proven, Reason: spec.NotExpired,
		},
		{Key: "mixed.read/latency-budget", Result: spec.Proven},
	}
	passed := []spec.CheckRecord{
		{Key: "mixed.read/cancel", Result: spec.Pass},
		{Key: "mixed.read/expires", Result: spec.Pass},
		{Key: "mixed.read/keeps-the-item", Result: spec.Pass},
		{Key: "mixed.read/latency-budget", Result: spec.Pass},
	}
	return corpus.Outcome{
		Fixture: ttl,
		Lock: []string{
			ttl + "\tmixed.put\tshape:writer\tsilent\t-",
			ttl + "\tmixed.read/cancel\tcancel\tcheck\t-",
			ttl + "\tmixed.read/expires\texpires\tcheck\twrite=mixed.put",
			ttl + "\tmixed.read/keeps-the-item\tkeeps-the-item\tauthored\t-",
			ttl + "\tmixed.read/latency-budget\tlatency-budget\tceiling\t-",
			ttl + "\tmixed.scan/zero-on-error\tzero-on-error\trefused\t-",
		},
		Records: []report.Record{
			record(spec.Falsification, "", proven...),
			record(spec.OnSubject, subject, passed...),
			record(spec.OnSubject, control, passed...),
		},
	}
}

// record returns a record of the declaration Mixed of the fixture ttl in the
// application, on the subject typ, with the checks and the refusal of the
// suite.
func record(application spec.Application, typ string, checks ...spec.CheckRecord) report.Record {
	return report.Record{
		Package:     ttl + "/ttltest",
		Test:        "TestMixedSuite",
		Suite:       spec.RecordSuite{Language: "go", Package: ttl, Declaration: "Mixed"},
		Application: application,
		Subject:     typ,
		Checks:      slices.Clone(checks),
		Refusals: []spec.RefusalRecord{{
			Claim: spec.ZeroOnError, Callable: "Mixed.Scan", Gate: spec.InputSamples,
			Why:    "eidos derives no sample of the parameter visit of Scan",
			Remedy: "author a sample of the parameter visit with the directive sample",
		}},
	}
}
