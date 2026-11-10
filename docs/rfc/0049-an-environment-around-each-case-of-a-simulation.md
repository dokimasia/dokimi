---
rfc: 0049
title: An environment of the consumer around each case of a simulation
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-10
updated: 2026-11-10
discussion: https://github.com/dokimasia/dokimi/issues/90
supersedes: RFC-0025 and RFC-0036 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0049: An environment of the consumer around each case of a simulation

## Summary

We propose that a consumer can run each case of a simulation inside an environment that it chooses:

- `//dokimi:around` marks the function `func(t *testing.T, run func(*testing.T))` that runs each case, such as a function that seeds `crypto/rand` with `testing/cryptotest.SetGlobalRandom`.
- The generated simulation sets `Around` to that function. A system of components runs `conform.Bubble` inside it.
- `sim.Simulation` gains `Serial`. `Falsify` proves each check of a simulation with `Serial` set in a subtest without `t.Parallel()`. The generated simulation sets `Serial` with an around.
- The generated `TestSimulation` and `TestSimulationFalsify` of such a simulation do not call `t.Parallel()`.

## Motivation

A machine that reads `crypto/rand` makes `sim.Test` fail as flaky at its first replayed case, because the replay reads other values:

- The generated simulation of a machine interface has no `Around`. `generator/sim/templates/golang/sim.tmpl:275` at ba748e1 writes `Around: conform.Bubble` for a system of components alone, and no directive sets an `Around`.
- The RawNode of etcd-io/raft v3.7.0 takes each randomized election timeout from `crypto/rand` (`raft.go:2054`, through the `globalRand` of `raft.go:97-104`). Its `Config` takes no source of randomness.
- `cryptotest.SetGlobalRandom` sets a deterministic source for `crypto/rand` during a test, and a test may call it again to reset the stream. It affects the whole process, so it panics in a parallel test and in a test with a parallel ancestor.
- The generated `TestSimulation` and `TestSimulationFalsify` call `t.Parallel()`, and `Falsify` runs each check in a parallel subtest (`sim/run.go:263-271` of dokimi-go at sim v0.12.0).

`prop.Around` runs each case, each shrink candidate and each replay inside its wrap, one at a time on the default single worker. A wrap that seeds `crypto/rand` gives each replay the random values of its case. sharder simulates the members of its Raft groups on the RawNode before it builds them, and its design seeds `crypto/rand` in the `Around` of the simulation.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The directive around | `go.dokimi.dev/dokimi`, generator sim | Marks the function that runs each case of the simulation |
| The serial falsification | `go.dokimi.dev/sim`, `Simulation.Serial` and `Falsify` | Proves each check in a subtest without `t.Parallel()` |
| The generated tests | `go.dokimi.dev/dokimi`, generator sim | Omit `t.Parallel()` for a simulation with an around |

### The directive

```go
// around seeds crypto/rand for each case, so a replay takes the same election
// timeouts.
//
//dokimi:around
func around(t *testing.T, run func(*testing.T)) {
	cryptotest.SetGlobalRandom(t, 1)
	run(t)
}
```

The function is in the package of the simulation, as the spec and the fingerprint are. A simulation has one around. The generator reports `BadFunction` at a second around and at a function of another signature.

The generated simulation of a machine interface writes the function into `Around` and sets `Serial`:

```go
		Around: around,
		Serial: true,
```

The simulation of a system of components runs `conform.Bubble` inside the around, so the seed of the around covers each bubble:

```go
		Around: func(t *testing.T, run func(*testing.T)) {
			around(t, func(t *testing.T) { conform.Bubble(t, run) })
		},
		Serial: true,
```

### The serial falsification

```go
type Simulation[E, F, S any] struct {
	// The fields of a simulation, as before.

	// Serial makes Falsify prove the checks one at a time, in subtests
	// without t.Parallel, for an Around that changes the state of the
	// process, such as one that calls cryptotest.SetGlobalRandom.
	Serial bool
}
```

`Falsify` runs each check in a subtest, as before. Without `Serial` the subtest calls `t.Parallel()`, and with `Serial` it runs to its end before the next one starts. `Test` and `Prove` do not change, because they run their properties on the test that they take. A property runs one case at a time on its default single worker.

### The generated tests

`TestSimulation` and `TestSimulationFalsify` of a simulation with an around do not call `t.Parallel()`, so no test of the simulation has a parallel ancestor. `FuzzSimulation` is unchanged, because the test of each input of a fuzz target does not run in parallel.

### Failure handling

| Failure | Behaviour |
|---|---|
| A function under `around` with another signature | `BadFunction` at the function |
| A second function under `around` | `BadFunction` at the second function, in the order of the identities |
| An around that seeds `crypto/rand` in a simulation without `Serial`, such as one written by hand | `cryptotest.SetGlobalRandom` panics in the first parallel subtest of `Falsify` |

### Migration

- A simulation without an around is unchanged, and its tests run in parallel as before.
- A simulation written by hand that seeds the process in its `Around` sets `Serial`, and its tests do not call `t.Parallel()`.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| A function of the consumer as the around | in | The consumer knows which state of the process its machines read, such as `crypto/rand`, the environment or the working directory | none known |
| An around on a system of components | in | A component can read `crypto/rand` too, and the around then runs outside the bubble of each case | none known |
| Serial tests for every around | in | `Falsify` cannot tell from a function whether it changes the state of the process, and each change of state that `testing` guards panics in a parallel test | An around that a consumer marks as safe for parallel tests |
| A directive that states only a seed | out | A seed covers `crypto/rand` alone, while an around also covers `t.Setenv` and `t.Chdir` | none known |

## Alternatives considered

### A seed in the simulation

`sim.Simulation` could take a seed and call `cryptotest.SetGlobalRandom` itself. sim would then decide which state of the process a case seeds, and a machine that reads the environment would need a second field.

### A source of randomness in the configuration of the machine

The consumer could pass a deterministic source to the machine. The `Config` of etcd-io/raft takes no such source, and its tests pin each timeout through a setter that only the tests of its package can call.

## Drawbacks

- The proofs of `Falsify` of a simulation with an around run one at a time. Each check adds its proof to the time of the test.
- The tests of such a simulation run outside the parallel tests of the package.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 90, a machine simulation cannot seed crypto/rand before each case | <https://github.com/dokimasia/dokimi/issues/90> |
| `testing/cryptotest.SetGlobalRandom` | <https://pkg.go.dev/testing/cryptotest#SetGlobalRandom> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
| RFC-0036, components that run goroutines | 0036-components-that-run-goroutines.md |
