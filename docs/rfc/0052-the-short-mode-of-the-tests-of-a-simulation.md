---
rfc: 0052
title: The short mode of the tests of a simulation
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-12
updated: 2026-11-12
discussion: https://github.com/dokimasia/dokimi/issues/95
supersedes: RFC-0025 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0052: The short mode of the tests of a simulation

## Summary

We propose that the generated tests of a simulation shorten their run under `go test -short`:

- `TestSimulation` calls `sim.TestShort`. It runs 20 random runs without the shares of the faults that the simulation requires, and runs each scenario once without its crash points.
- `TestSimulationFalsify` skips.
- `FuzzSimulation` does not change.

Without `-short`, the three tests run as they run at dokimi v0.13.0.

## Motivation

The generated tests of a simulation do not check `testing.Short`:

- The template of the tests calls `sim.Test` and `sim.Falsify` without a check of the mode (`generator/sim/templates/golang/sim.test.tmpl` at v0.13.0).
- `sim.Test` and `sim.Falsify` do not read the mode either (`sim/run.go:250-283` of dokimi-go at sim v0.14.0).

A repository on ergon runs its test stage and its race stage with `-short`. The job Longtest Go of its `nightly.yml` runs the long tests. In sharder the two simulations take 753 s of the 946 s of the race stage.

We measured the parts of both simulations on a copy of the working tree of sharder, with sim v0.14.0 and assert v0.22.0. The runs used `go test -short` without `-race`, Go 1.27.4, 16 cores and the seeds of the ci profile:

| Part | `membertest` | `machinetest` |
|---|---|---|
| `TestSimulation` | 31.5 s | 31.6 s |
| The random runs of `TestSimulation`, with the shares | 13.5 s for 200 runs | 23.9 s for 200 runs |
| The scenarios of `TestSimulation`, with their crash points | 18.0 s for 7 scenarios | 7.6 s for 4 scenarios |
| The proofs of `TestSimulationFalsify` | 5.5 s for 14 proofs | 46.3 s for 24 proofs |
| `FuzzSimulation` | 0.00 s | 0.00 s |
| 100 random runs without the shares | 6.6 s | 13.3 s |
| 10 random runs without the shares | 0.6 s | 1.2 s |

The time of the scenarios is the time of `TestSimulation` minus the time of its random runs at the same seed. The proofs of `machinetest` run in parallel subtests, so the table lists the sum of their times. The longest proof took 14.7 s.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| `TestShort` | `go.dokimi.dev/sim` | Runs the short form of `Test` |
| The generated tests | `go.dokimi.dev/dokimi`, generator sim | Call `TestShort` and skip the proofs in short mode |

### The short form of `Test`

```go
// TestShort runs the short form of Test on tb for the short mode of go test:
// 20 random runs of s without the shares of the faults that s requires, and
// then each scenario of s once without its crash points. It replays the 16th
// passing case from its choices, as Test does.
func TestShort[E, F, S any](tb assert.TB, s Simulation[E, F, S])
```

- The random runs have the contract and the store of the random runs of `Test`. A case that a run of `Test` stored runs again in `TestShort`, and a failing case of `TestShort` is stored for the next run of `Test`.
- The 20 random runs of `TestShort` are a fifth of the 100 runs of `Test`. rapid v1.3.0 also runs a fifth of its checks in short mode. The determinism check replays every 16th passing case, so 20 runs keep one replay.
- The random runs do not check the shares. `prop.Require` decides a share after the stated number of cases, and again each time the number doubles, up to eight times the number. With the shares, a short run of 20 cases can take 160 runs.
- `TestShort` does not run the crash points. They run a scenario of n steps 2n more times, and they took 18.0 s of the 31.5 s of `TestSimulation` of `membertest`.

### The generated tests

```go
// TestSimulation runs the simulation of Machine over random runs and the
// crash points of each scenario. In short mode it runs 20 random runs and
// each scenario once.
func TestSimulation(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		sim.TestShort(t, membertest.MachineSimulation())
		return
	}
	sim.Test(t, membertest.MachineSimulation())
}

// TestSimulationFalsify proves each rule of a runtime and each invariant
// of the simulation of Machine on its planted defect. It skips in short
// mode.
func TestSimulationFalsify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping the proofs of the planted defects in short mode")
	}
	t.Parallel()

	sim.Falsify(t, membertest.MachineSimulation())
}
```

- `TestSimulationFalsify` skips before `t.Parallel()`, as the section Skipping of the package `testing` shows. Each proof runs at a seed from the key of its check, so its outcome changes only with the code of the simulation.
- `FuzzSimulation` does not change. Without `-fuzz`, `go test` runs the seed corpus of a fuzz target alone. For a simulation that corpus is the stored cases, which took 0.00 s in both packages.
- A simulation with an around keeps its tests out of the parallel tests of its package.

### The cost in short mode

We regenerated the copy of sharder with this design. The race stage of ergon runs `go test -race -short -count=1 -p=1`, and we ran both packages with that command on 16 cores:

| Package | At dokimi v0.13.0, from the issue | With the short mode |
|---|---|---|
| `internal/machine/machinetest` | 320.9 s | 18.7 s |
| `store/quorum/internal/member/membertest` | 431.9 s | 11.8 s |

Without `-race`, `TestSimulation` took 2.7 s in `machinetest` and 1.1 s in `membertest`.

### `testing.Short` in the generated tests

The generated tests read `testing.Short`, and `sim` does not. A test that runs `sim.Test` as its subject then runs the same way under `go test` and `go test -short`. The specs of `sim` are such tests, and ergon runs them with `-short`. The fixtures of the corpus that prove that a planted defect fails `sim.Test` are such tests too.

### Failure handling

| Case | Behaviour |
|---|---|
| A random run of `TestShort` fails | `TestShort` fails with the record of the random runs of `Test`, and stores the case for the next run of `Test` |
| A scenario of `TestShort` fails | `TestShort` fails with the record of the scenario, as `Test` fails at a scenario without a crash |
| The cases of the full run miss a share | The short run passes, and the full run fails |
| A crash point or a proof fails | The short run passes, and the full run fails |

### Migration

- A regeneration writes the new tests of each simulation. Until then the generated tests run in full in short mode.
- A test of a simulation that a consumer writes by hand calls `sim.TestShort` in short mode to take the short form.
- A run without `-short`, such as the job Longtest Go of `nightly.yml`, runs every test in full.

### Verdicts

| Capability | Verdict | Reason | What would reverse it |
|---|---|---|---|
| A short form of `TestSimulation` | in | The random runs, the determinism check and the scenarios check the system in short mode too | none known |
| 20 random runs | in | They are a fifth of the 100 runs of `Test`, as rapid runs a fifth of its checks, and they keep one replay of the determinism check | A simulation whose 20 runs take longer than the short stage of its repository allows |
| The shares in the short form | out | A share describes the cases of a long run, and with the shares a short run of 20 cases can take 160 runs | none known |
| The crash points in the short form | out | They took 18.0 s of the 31.5 s of `TestSimulation` of `membertest` | none known |
| The proofs in short mode | out | They check the oracles at fixed seeds, and they took 46.3 s of `machinetest` | none known |
| A read of the mode inside `sim` | out | A test that runs `sim.Test` as its subject would change with `-short` | none known |
| A count of short runs for each simulation | out | With 20 runs the race stage of the two packages of sharder took 18.7 s and 11.8 s | A simulation whose 20 runs take longer than the short stage of its repository allows |
| A short mode of the TypeScript tests | out | Vitest has no short mode, and the TypeScript runtime has no simulation | A short mode of the TypeScript stages of ergon |

## Alternatives considered

### `sim` reads the short mode

`Test` and `Falsify` could call `testing.Short` themselves, as `checkTB` of rapid does. The specs of `sim` run with `-short` in the test stage of dokimi-go. A spec of the shares or of the crash points would then check the short form there. A fixture of the corpus that proves that a planted defect fails `sim.Test` would pass under `go test` and could fail under `go test -short`.

### A field of the simulation

`Simulation` could take a field `Short` that `Test` and `Falsify` read. The generated tests would copy the simulation and set the field from `testing.Short`. A reader of the call `sim.Test(t, s)` could then not tell the form of the run from the call.

### 100 random runs without the shares

The short form could keep the 100 runs of `Test` and drop the shares alone. The 100 runs took 6.6 s and 13.3 s, five times the time of 20 runs.

### Skip `TestSimulation` in short mode

The short stage would then not check the simulation at all. A change that breaks a scenario or the determinism of the machines would pass the CI of the change.

## Drawbacks

- A defect that only the shares, the crash points or the proofs find passes the CI of a change. The next run without `-short` fails.
- The generated `TestSimulation` makes one of two calls of the module, by the mode.
- `sim` gains one exported function.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 95, the generated tests of a simulation ignore `-short` | <https://github.com/dokimasia/dokimi/issues/95> |
| The flag `-short` of `go test`: "Tell long-running tests to shorten their run time." | <https://pkg.go.dev/cmd/go#hdr-Testing_flags> |
| The section Skipping of the package `testing` | <https://pkg.go.dev/testing#hdr-Skipping> |
| rapid v1.3.0, a fifth of the checks in short mode | <https://github.com/flyingmutant/rapid/blob/v1.3.0/engine.go#L274-L277> |
| RFC-0025, a deterministic simulation of sans-IO machines | [RFC-0025](0025-a-simulation-of-sans-io-machines.md) |
| RFC-0049, an environment of the consumer around each case of a simulation | [RFC-0049](0049-an-environment-around-each-case-of-a-simulation.md) |
