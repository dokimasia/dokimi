// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package corpus

import (
	"strings"

	"go.dokimi.dev/dokimi/spec"
)

// Update returns expected with the declarations of each fixture rewritten
// from its outcome in outcomes. The checks, the authored checks, the
// ceilings, the refusals, the partners and the unclaimed classifications come
// from the lines of the lock, the gate of each refusal from the records, and
// the law and the reason code of each check from its claim in the spec.
// Update keeps the level, the subjects and the controls of each expectation,
// and the declarations of a fixture without an outcome.
func Update(expected []Expectation, outcomes []Outcome) []Expectation {
	byFixture := make(map[string]*Outcome, len(outcomes))
	for i := range outcomes {
		byFixture[outcomes[i].Fixture] = &outcomes[i]
	}
	out := make([]Expectation, len(expected))
	for i, e := range expected {
		if o, ok := byFixture[e.Fixture]; ok {
			e.Declarations = declarationsOf(o)
		}
		out[i] = e
	}
	return out
}

// declarationsOf returns the declarations of the lines of the lock of o by
// the keys of their names.
func declarationsOf(o *Outcome) map[string]Declaration {
	gates := gatesOf(o.Records)
	declarations := map[string]Declaration{}
	for _, text := range o.Lock {
		f := fieldsOf(text)
		key, claim, status := f[1], f[2], f[3]
		name := declarationOf(key)
		d, ok := declarations[name]
		if !ok {
			d = Declaration{Checks: map[string]Check{}, Refusals: map[string]spec.GateID{}}
		}
		switch status {
		case kindCheck:
			c, _ := spec.ClaimOf(spec.ClaimID(claim))
			d.Checks[key] = Check{
				Claim: spec.ClaimID(claim), Law: c.Law, Reason: c.Defect.Reason, Partners: partnersOf(f[4]),
			}
		case kindAuthored:
			d.Authored = append(d.Authored, key)
		case kindCeiling:
			d.Ceilings = append(d.Ceilings, key)
		case kindRefused:
			d.Refusals[key] = gates[refusalKey(key, claim)]
		default:
			d.Unclaimed = append(d.Unclaimed, Unclaimed{
				Classification: claim, Callable: key, Stance: spec.StanceID(status),
			})
		}
		declarations[name] = d
	}
	return declarations
}

// partnersOf returns the partners of the partners field of a line of the
// lock, by their names, and nil for a dash.
func partnersOf(field string) map[string]string {
	if field == noPartners {
		return nil
	}
	partners := map[string]string{}
	for part := range strings.SplitSeq(field, ",") {
		name, key, _ := strings.Cut(part, "=")
		partners[name] = key
	}
	return partners
}
