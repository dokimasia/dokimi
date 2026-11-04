// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package report_test

import (
	"encoding/json"
	"strings"

	"go.dokimi.dev/spec"
)

// record returns the JSON of a suite record of the declaration Store in the
// application, with the checks.
func record(application spec.Application, checks ...spec.CheckRecord) string {
	r := spec.SuiteRecord{
		Spec:        "1.0.0",
		Suite:       spec.RecordSuite{Language: "go", Package: "example.com/store", Declaration: "Store"},
		Generator:   spec.RecordTool{Name: "dokimi", Version: "0.1.0"},
		Runtime:     spec.RecordTool{Name: "go.dokimi.dev/conform", Version: "v0.9.0"},
		Seed:        "0x5f3a1c9e27d4b860",
		Application: application,
		Checks:      checks,
		Refusals: []spec.RefusalRecord{{
			Claim: spec.ZeroOnError, Callable: "Store.Scan", Gate: spec.Failure, Why: "Scan cannot fail",
			Remedy: "return an error",
		}},
		Unclaimed: []spec.UnclaimedRecord{{
			Classification: "mixin:concurrent", Callable: "Store.Put", Stance: spec.Needs,
			Note: "needs concurrent callers",
		}},
	}
	data, _ := json.Marshal(r)
	return string(data)
}

// check returns the record of the check of the identity with the result.
func check(identity string, result spec.Result) spec.CheckRecord {
	return spec.CheckRecord{Identity: identity, Key: strings.ToLower(identity), Result: result}
}

// events returns the events of go test -json that write each record under the
// test TestStore, one attribute event for each record, split into parts of at
// most size bytes.
func events(size int, records ...string) string {
	var b strings.Builder
	b.WriteString(`{"Action":"start","Package":"example.com/store/storetest"}` + "\n")
	b.WriteString("a line that is no event\n")
	for i, r := range records {
		key := "dokimi.suite." + string(rune('1'+i))
		for len(r) > 0 {
			part := r[:min(size, len(r))]
			r = r[len(part):]
			line, _ := json.Marshal(map[string]string{
				"Action": "attr", "Package": "example.com/store/storetest", "Test": "TestStore", "Key": key,
				"Value": part,
			})
			b.Write(line)
			b.WriteByte('\n')
		}
	}
	b.WriteString(`{"Action":"attr","Package":"example.com/store/storetest","Test":"TestStore",` +
		`"Key":"dokimi.assert.1","Value":"{}"}` + "\n")
	return b.String()
}
