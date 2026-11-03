# Copyright Dokimasia B.V. 2026
# SPDX-License-Identifier: Apache-2.0

# The settings of the repository of dokimi, which ergon init appends to the Makefile.

# The tests of the gate run their properties under the profile ci of assert, which derives the
# seed of each run from the contract of its property. The hooks and CI then run the same cases,
# so a gate that passes before a push passes in CI, and a failure in CI repeats on a workstation.
# A value of DOKIMI_ASSERT_PROP_PROFILE in the environment takes precedence, such as default for
# new random seeds.
test-go race-go: export DOKIMI_ASSERT_PROP_PROFILE ?= ci

# report-golden writes the run and the mutation record of the corpus fixture
# iface/mutation/deadcheck into REPORT_GOLDEN, which the spec of dokimi report reads. It keeps the
# attribute events of go test -json, which contain the suite records, as JSON Lines in run.jsonl.
# dokimi-mutate-go runs in the modes of every killer and of the coverage of each leaf test, and
# exits with 1 for the survivors that the fixture plants. The run takes the seeds of the profile
# ci, so each run writes the same seeds.
REPORT_GOLDEN := cmd/dokimi/testdata/report
REPORT_FIXTURE := ./iface/mutation/deadcheck/deadchecktest

.PHONY: report-golden
report-golden: export DOKIMI_ASSERT_PROP_PROFILE = ci
report-golden: ## Write the run and the mutation record of the fixture that the spec of dokimi report reads
	$(GO) -C corpus/go test -count=1 -json $(REPORT_FIXTURE) > $(REPORT_GOLDEN)/events.json
	grep '"Action":"attr"' $(REPORT_GOLDEN)/events.json > $(REPORT_GOLDEN)/run.jsonl
	rm $(REPORT_GOLDEN)/events.json
	$(DOKIMI_MUTATE_GO) -C corpus/go --all-killers --subtests --json $(REPORT_FIXTURE) \
		> $(REPORT_GOLDEN)/mutation.json || [ $$? -eq 1 ]
