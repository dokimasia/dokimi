---
rfc: 0046
title: A change of the consumer in the final phase
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-09
updated: 2026-11-09
discussion: https://github.com/dokimasia/dokimi/issues/85
supersedes: RFC-0025 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0046: A change of the consumer in the final phase

## Summary

We propose that the directive action takes the word `final` alone, for a change of the consumer that the run takes once at the start of the final phase:

- `//dokimi:action final` marks a function of the form of an action, `func(c *prop.Case) func(w *sim.World[E, F])`.
- The generated simulation writes the action with `Drain` and without `Call`.
- `//dokimi:action call final` keeps its meaning: a call that the run makes once after the faults end.

## Motivation

The run of sim takes each available action with `Drain` once at the start of the final phase, after the recovery ends each fault in force (`sim/run.go:791-799` of dokimi-go at sim v0.11.0). It takes an action that is not a call there too, with an input that has no client (`sim/run.go:679-688`).

The directive action has the positional params `call` and `final` (`generator/sim/plugin.go:603-619` at 711a71ab). Its first word must be `call`, so the generator refuses `//dokimi:action final` with `EID-0009`. The directive takes no keyed param, so `final=final` fails with `EID-0006`. A consumer cannot declare a change of the final phase.

A write transaction of the simulation of sharder that is still open in the final phase blocks the acquisition of the fence of its shard. The controller then cannot add the shard to a host, and the convergence fails at "each transition has ended". A deployment ends its transactions once its faults end. The simulation models that with a final change that commits each open transaction. Such a change is no call of a client.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The directive action | `go.dokimi.dev/dokimi`, generator sim | Takes `final` as its first word |
| The action of the model | `go.dokimi.dev/dokimi`, generator sim | Has `Final` without `Call` for a final change |
| The run of sim | `go.dokimi.dev/sim` | No change, because it already takes an action with `Drain` that is not a call |

### The forms of the directive

| Directive | Function | Meaning |
|---|---|---|
| `//dokimi:action` | `func(c *prop.Case) func(w *sim.World[E, F])` | A change that the case chooses in the active phase |
| `//dokimi:action call` | `func(c *prop.Case) func(w *sim.World[E, F], client int)` | A call of a client in the active phase |
| `//dokimi:action final` | `func(c *prop.Case) func(w *sim.World[E, F])` | A change that the run takes once at the start of the final phase |
| `//dokimi:action call final` | `func(c *prop.Case) func(w *sim.World[E, F], client int)` | A call that the run makes once at the start of the final phase |

The first word of the directive is `call` or `final`, and the second word is `final` after `call`. A directive whose first word is `final` takes no second word. The generator reports `BadParam` at a function whose directive has a word after `final`, because the schema of a directive states the choices of each word alone.

The generated simulation writes a final change as an action with `Drain` and the run of a change:

```go
{
	Name:  "commitWrites",
	Drain: true,
	Input: func(c *prop.Case) any { return commitWrites(c) },
	Run: func(w *world, _ int, in any) {
		run, _ := in.(func(*world))
		run(w)
	},
},
```

### Failure handling

| Failure | Behaviour |
|---|---|
| `//dokimi:action final final` | `BadParam` at the function |
| `//dokimi:action final call` | The directive refuses the second word, as it refuses each word other than `final` there |
| `//dokimi:action final` on a function of the form of a call | `BadFunction` at the function, as for each change of another form |

### Migration

- A simulation without the new form generates as before.
- The corpus does not change, because none of its simulations states the form.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| `final` as the only word of the directive action | in | The run already takes a final action that is not a call, and one directive marks each action of a simulation | none known |
| A keyed param `final` | out | The directive takes positional words alone, so a final call would need a second spelling | none known |
| A second directive, such as `//dokimi:settle` | out | A final change has the form of an action, so a second directive would repeat its form and its binding | A final change with another form than an action |

## Alternatives considered

### A keyed param

`//dokimi:action final=true` would mark a final change. The directive takes no keyed param. A final call would need the keyed form too, so one meaning would have two spellings.

### A second directive

A directive `//dokimi:settle` could mark a change of the final phase. It would bind a function of the form of an action, and a reader would find the actions of a simulation under two directives.

### An end of the transactions in sim

The run could end each open transaction of the consumer at the start of the final phase. sim does not know the transactions of an application, and the consumer decides how an application ends them.

## Drawbacks

- The first positional param has the key `call` and also takes the word `final`, so the error of a wrong first word names the param `call`.
- The generator checks the word after `final` itself, beside the schema of the directive.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 85, the word `final` alone in the directive action | <https://github.com/dokimasia/dokimi/issues/85> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
