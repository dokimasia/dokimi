---
rfc: 0047
title: The planted defect of a convergence
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-10
updated: 2026-11-10
discussion: https://github.com/dokimasia/dokimi/issues/87
supersedes: RFC-0025 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0047: The planted defect of a convergence

## Summary

We propose that a convergence has a planted defect, as an invariant has one:

- `sim.Convergence` gains `Defect`, the system of its planted defect.
- `Falsify` proves each convergence, and `Prove` proves a convergence on its defect.
- `//dokimi:converges` takes the required param `defect`, as `//dokimi:invariant` does.
- `Prove` fails for an invariant or a convergence without a planted defect.
- `Prove` fails for a name that two checks of the simulation share, because it finds a check by its name.

## Motivation

`//dokimi:invariant` requires the param `defect`, and `TestSimulationFalsify` proves each invariant on the system of its defect. `//dokimi:converges` takes the param `within` alone (`generator/sim/plugin.go:628-648` at 8c423be2). `Falsify` collects the identities of the rules and the statements of the invariants, and `Prove` proves each of them (`sim/run.go:262-305` of dokimi-go at sim v0.11.0). No generated test proves that a convergence can fail.

A spec of the consumer cannot run a convergence on a broken system either. The functions under the directives are unexported in the package of the simulation, and a spec is a black-box test in the package with the suffix `_test`.

The simulator of a pool of theke checks a liveness property of `MaxCompactionLag`: once the faults end and the lag has passed, the pool reports 0 overdue eligible bytes. The testing standard of theke requires a proof that each check fails against a system that breaks it. The convergence is the one check of the simulator without such a proof.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The convergence | `go.dokimi.dev/sim`, `Convergence` | Gains `Defect` |
| The falsification | `go.dokimi.dev/sim`, `Falsify` and `Simulation.Prove` | Proves each convergence on its defect. Fails for a check without a defect, and for a name that two checks share |
| The directive converges | `go.dokimi.dev/dokimi`, generator sim | Takes the required param `defect` |

### The convergence

```go
// Convergence is a liveness property that the run checks once the faults
// have ended and Within has passed, with the system of its planted defect.
type Convergence[E, F any] struct {
	Check     func(h *Check, w *World[E, F])
	Defect    System[E, F]
	Statement string
	Within    time.Duration
}
```

### The proof of a convergence

`Prove` looks a check up by its name: a rule by its identity, then an invariant and then a convergence by its statement. It proves a convergence on a simulation with the system `Defect` and that convergence alone, without the spec, the rules and the invariants, as it proves an invariant alone.

- Each random run of the defect takes the steps of the active phase and ends the faults. It then settles past `Within` and checks the convergence.
- The runs must fail at the seed of the check. Their first failure must be the failure of the assertion whose message is the statement of the convergence.
- The proof does not run the scenarios, because a scenario ends without the final phase.

`Falsify` adds a parallel test for each convergence, named after its statement. `Prove` fails with `sim: the check "<name>" has no planted defect` for an invariant or a convergence whose `Defect` is nil.

### The directive

```go
// Once the faults end, the backup has the value of the primary.
//
//dokimi:converges within=settled defect=stale
func converged(h *sim.Check, w *world) { ... }
```

The param `defect` names a function of the form of a system, `func(w *sim.World[E, F])` or `func(c *prop.Case) func(w *sim.World[E, F])`, as the param `defect` of an invariant does. The generated simulation writes it into `Defect`:

```go
Converges: []sim.Convergence[E, F]{
	{
		Statement: "Once the faults end, the backup has the value of the primary",
		Check:     converged,
		Defect:    func(*prop.Case) func(w *world) { return stale },
		Within:    settled(),
	},
},
```

The generator reports `BadFunction` at a convergence whose defect has another form, as at an invariant.

### The names of the checks

`Falsify` names each test after its check, and `Prove` finds the check by that name: a rule by its identity, an invariant or a convergence by its statement. Both tests of a name that two checks share would prove the first of the two checks. `Prove` fails with `sim: the simulation has two checks named "<name>"` for such a name. The failure covers a generated simulation and a simulation written by hand.

### Failure handling

| Failure | Behaviour |
|---|---|
| `//dokimi:converges` without `defect` | The validation of the directive reports the missing required param, and the generator leaves the convergence out |
| A defect that is no system | `BadFunction` at the convergence |
| A name that two checks share | `Prove` fails with a message that contains the name |
| A convergence whose defect passes the convergence | `TestSimulationFalsify` fails at the seed of the convergence |
| A `Convergence` without `Defect` in a simulation written by hand | `Prove` fails with a message that contains the check |

### Migration

- Each `//dokimi:converges` states `defect`.
- A simulation written by hand states `Defect` for each convergence, or `Falsify` fails at it.
- A simulation with two checks of one name gives one of them another name, or `Falsify` fails at both.
- The corpus fixtures `sim/mirrored` and `sim/replicated` gain a planted defect of their convergence.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| A required param `defect` | in | Each check of a simulation then has a proof, as each invariant has one | A convergence that no system can break |
| An optional param `defect` | out | A convergence without a defect would pass `Falsify` without a proof | none known |
| A proof by the scenarios | out | A scenario ends without the final phase, so it never checks a convergence | A scenario that runs the final phase |
| A planted machine of the form of a rule | out | A convergence states a property of the whole system, and a planted machine breaks one claim of one machine | none known |
| A refusal of a second statement in the generator | out | `Prove` fails for a name that two checks share, in every simulation | A consumer who needs the position of the second check before the tests run |

## Alternatives considered

### An optional param

The param `defect` could be optional, and `Falsify` could prove the convergences that have one. A consumer would then not see which convergence lacks a proof, and the testing standard of theke would still fail at the simulator.

### A proof in the specs of the consumer

The consumer could write a spec that runs the convergence against a broken system. The functions of the simulation are unexported, and the black-box spec cannot call them.

## Drawbacks

- Each simulation with a convergence needs a defect system. The defect of `sim/mirrored` has 6 lines, and the defect of `sim/replicated` has 24.
- `TestSimulationFalsify` runs one property more for each convergence.
- The change breaks each `//dokimi:converges` without `defect`. It also breaks each simulation written by hand whose convergence has no `Defect`, or whose two checks share a name.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 87, a convergence takes no planted defect | <https://github.com/dokimasia/dokimi/issues/87> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
