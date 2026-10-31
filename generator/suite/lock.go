// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"

	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// LockTag is the tag of the family of the lock, and LockName the name of the
// file of the lock of a plan whose layout does not rename it.
const (
	LockTag  = "lock"
	LockName = "dokimi.lock"
)

// lockWord is the word of the family of the lock, and lockTemplate the
// template that writes its file.
const (
	lockWord     = "lock"
	lockTemplate = "lock.tmpl"
)

// The statuses of the lines of the lock. The stance of an unclaimed
// classification is the status of its line.
const (
	statusCheck    = "check"
	statusAuthored = "authored"
	statusCeiling  = "ceiling"
	statusRefused  = "refused"
)

// noPartners is the partners field of a line without partners.
const noPartners = "-"

// lockFields is the number of fields of a line of the lock: the package, the
// key, the claim, the status and the partners.
const lockFields = 5

// Lock is the record of one suite in the lock of its plan: one line for each
// derived check, each check of the consumer, each ceiling of a budget, each
// refusal and each unclaimed classification of the suite, sorted, each line
// once.
type Lock struct {
	Lines []string
}

// lockOf returns the record of the suite s in the lock. A derived check that
// the suite refuses after the derivation, because the file of the suite
// cannot write it, has the line of its refusal alone.
func lockOf(s *Suite) Lock {
	var lines []string
	add := func(key, claim, status, partners string) {
		lines = append(lines, strings.Join([]string{s.Package, key, claim, status, partners}, "\t"))
	}
	refused := map[string]bool{}
	for _, r := range s.Derivation.Refusals {
		refused[r.Key] = true
		add(r.Key, string(r.Claim), statusRefused, noPartners)
	}
	for i := range s.Derivation.Checks {
		if c := &s.Derivation.Checks[i]; !refused[c.Key] {
			add(c.Key, string(c.Claim.ID), statusCheck, partnersOf(c))
		}
	}
	for i := range s.Callables {
		c := &s.Callables[i]
		for j := range s.Checks {
			if check := &s.Checks[j]; check.Field == c.Field {
				add(spec.Key(c.Name+"/"+check.Segment), check.Segment, statusAuthored, noPartners)
			}
		}
		for _, segment := range s.Calls[i].Budget.segments() {
			add(spec.Key(c.Name+"/"+segment), segment, statusCeiling, noPartners)
		}
	}
	for _, u := range s.Derivation.Unclaimed {
		add(u.Key, u.Classification, string(u.Stance), noPartners)
	}
	slices.Sort(lines)
	return Lock{Lines: slices.Compact(lines)}
}

// partnersOf returns the partners field of the line of the check c: each
// partner as its name and the key of its callable, in the order of the names.
func partnersOf(c *derive.Check) string {
	if len(c.Pairs) == 0 {
		return noPartners
	}
	names := slices.Sorted(maps.Keys(c.Pairs))
	parts := make([]string, len(names))
	for i, name := range names {
		parts[i] = name + "=" + c.Pairs[name].Partner.Key()
	}
	return strings.Join(parts, ",")
}

// lockLines returns the lines of the records of a file of the lock, sorted,
// each line once. records is the field Records of the view of the file, and
// the field Data of each of its elements is the record of one suite.
func lockLines(records any) []string {
	v := reflect.ValueOf(records)
	var lines []string
	for i := range v.Len() {
		lock, _ := reflect.TypeAssert[Lock](v.Index(i).FieldByName("Data"))
		lines = append(lines, lock.Lines...)
	}
	slices.Sort(lines)
	return slices.Compact(lines)
}

// lost returns [LostCheck] for a removed line of the lock whose status is
// check, authored or ceiling. It returns false for every other line: a line
// of a refusal or of an unclaimed classification, whose diff shows in review,
// and a line of the header.
func lost(r plugin.RemovedLine) (diag.Diag, bool) {
	fields := strings.Split(r.Text, "\t")
	if len(fields) != lockFields || !slices.Contains([]string{statusCheck, statusAuthored, statusCeiling}, fields[3]) {
		return diag.Diag{}, false
	}
	return diag.Diag{
		Code: LostCheck,
		Msg: fmt.Sprintf("%s (%s) of %s is no longer generated as the lock states. Commit the regenerated lock "+
			"to approve the change", fields[1], fields[2], fields[0]),
	}, true
}
