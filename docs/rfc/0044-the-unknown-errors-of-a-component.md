---
rfc: 0044
title: The errors of a component that leave the outcome of a call unknown
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-09
updated: 2026-11-09
discussion: https://github.com/dokimasia/dokimi/issues/83
supersedes: RFC-0036 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0044: The errors of a component that leave the outcome of a call unknown

## Summary

We propose that a system of components states the errors of a call of a component that leave the outcome of the call unknown, as a machine states them for a reply:

- The directive system of a system of components takes the param `unknown`, such as `unknown=[volume.ErrAmbiguous]`.
- The generated simulation writes those errors into `Anatomy.Unknown`.
- `World.Invoke` completes a call as unknown for an error that `Anatomy.Unknown` reports, beside the four causes that it completes as unknown now.

## Motivation

`World.Invoke` completes a call of a component as unknown only for an error that wraps `context.Canceled`, `context.DeadlineExceeded`, `medium.ErrCrashed` or `host.ErrReset`. It completes the call as failed for any other error (`sim/component.go:63-65` and `sim/component.go:165-172` of dokimi-go at sim v0.10.2). `history.Call.Fail` records that the call took no effect, so a later read that finds the effect of the call fails the check of the history.

theke returns an error that wraps `volume.ErrAmbiguous` for a write that its volume admitted and whose outcome the call does not know, such as a write whose sync failed. A crash of the disk can then keep the write or lose it. When the disk keeps it, a later read returns the write. The simulation of a pool of theke records such a write as failed. The check of its history then rejects the read that returns the write.

A machine already states such errors in the param `unknown` of its directive reply. The generated anatomy reports them through `Anatomy.Unknown`. A system of components has no such param. Its simulation has the zero `Anatomy`.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The completion of a call of a component | `go.dokimi.dev/sim`, `World.Invoke` | Completes a call as unknown for an error that `Anatomy.Unknown` reports |
| The param `unknown` of the directive system | `go.dokimi.dev/dokimi`, generator sim | States the errors of a call of a component that leave its outcome unknown |
| The anatomy of a system of components | The template of the generator sim | An `Anatomy` with the field `Unknown` alone, where the param states errors |

### The completion of a call

| Error of the call | Completion of the call |
|---|---|
| None | OK with the output of the call |
| An error that wraps `context.Canceled`, `context.DeadlineExceeded`, `medium.ErrCrashed` or `host.ErrReset` | Unknown with the error |
| An error for which `Anatomy.Unknown` reports true | Unknown with the error |
| Any other error | Failed with the error |

The doc comment of `Anatomy.Unknown` states that the field covers the error of a reply of a machine and the error of a call of a component.

### The directive

```go
//dokimi:system tolerate=[crash:1, tasks] unknown=[volume.ErrAmbiguous]
func pool(w *sim.World[sim.None, sim.None]) { ... }
```

The param `unknown` is a list of package variables, as the param `unknown` of the directive reply is. A variable of another package has the name of its package, such as `volume.ErrAmbiguous`. The generated simulation of the system then has this anatomy:

```go
Anatomy: sim.Anatomy[sim.None, sim.None]{
	Unknown: func(err error) bool {
		return errors.Is(err, volume.ErrAmbiguous)
	},
},
```

The directive system of a machine interface refuses the param with `BadParam`, because the directive reply of the machine states such errors.

### Failure handling

| Failure | Behaviour |
|---|---|
| An error that a cause of `Invoke` and `Anatomy.Unknown` both cover | Unknown, as for either |
| A name in the param that is no package variable | The directive reports the error of its resolution |
| The param `unknown` on the system of a machine interface | `BadParam` at the function of the system |

### Migration

- A simulation without the param completes its calls as before.
- A consumer that classified such an error as failed states the error in the param, and its stored cases can stop replaying the same history.
- The corpus does not change. Its systems of components state no such error.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| The field `Unknown` of `Anatomy` | in | It already classifies the errors of a reply of a machine, so one field states the errors of every call whose outcome is unknown | A consumer whose machines and components need different classifications of one error |
| A field of `Simulation` for the errors of components | out | It would state the classification of `Anatomy.Unknown` a second time | none known |
| An interface of the error, such as a method `Unknown() bool` | out | The errors of theke are sentinels of the module of theke, and such a sentinel does not implement an interface of sim | A consumer whose errors are types that it declares |

## Alternatives considered

### A field of `Simulation`

`Simulation` could take the errors of a call of a component in a separate field. The field would state the classification that `Anatomy.Unknown` states for a machine, so a reader would find the errors of an unknown outcome in two places.

### An interface of the error

`Invoke` could complete a call as unknown for an error with a method that reports an unknown outcome. `volume.ErrAmbiguous` is a sentinel of theke without such a method. theke would have to wrap each such error in a type that it declares.

### A cancel in the error of the component

A component could wrap its error in `context.Canceled` to make `Invoke` complete the call as unknown. The history would then state a cancel that the call never had.

## Drawbacks

- The simulation of a system of components, whose anatomy was the zero `Anatomy`, has an anatomy with one field where the param states errors.
- The template of the generator writes the field `Unknown` in two places: the anatomy of a machine and the anatomy of a system of components.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 83, a call of a component cannot end as unknown for an error of the component | <https://github.com/dokimasia/dokimi/issues/83> |
| RFC-0036, a simulation of components that run goroutines | 0036-components-that-run-goroutines.md |
