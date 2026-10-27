---
rfc: 0012
title: The sample of a parameter in the generated calls
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-27
updated: 2026-10-27
discussion: https://github.com/dokimasia/dokimi/issues/36
supersedes: RFC-0006 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0012: The sample of a parameter in the generated calls

## Summary

A parameter with an authored sample, as `//dokimi:sample value=0x01` on the parameter states it, takes that sample in every call that the suite generates for its callable. The call generator draws every other parameter as before. The benchmarks, the allocation budgets and the planted defects of the ceilings of a callable take their calls from the call generator. So they measure the method in the part of its domain that its contract admits. A sample of the type of the parameter keeps its role as the starting value of the type, and the checks of a consumer draw their own inputs as before.

## Motivation

The call generator of a callable gives each parameter a value of the whole domain of its type, also where the parameter has an authored sample. `params` of the suite generator builds the value of each parameter from its type alone (`generator/suite/call.go:149` at dokimi 71136eb). The suite reads the sample of a type, in `startOf` and in `wholeStartOf` (`generator/suite/override.go:522` and `:614`), and not the sample of a parameter. eidos stamps the sample of a parameter under `gen.sample` on the parameter (`authored/sample.go` at eidos-core v0.10.0), and the suite generator does not read it.

So a method that refuses part of the domain of a parameter with the panic of its contract panics in the benchmark of its subject and in its allocation budget. `crypto.Hasher` of core is such a declaration. `HashTagged(r Role, data []byte)` takes a role below 0x80, and `CombineTagged(r Role, left, right Digest)` takes a role of 0x80 or above, so a leaf of a tree and an interior node never share a role. The benchmark of a hasher with dokimi v0.3.0 and conform v0.7.0 panics:

```text
panic: tagged: HashTagged takes a role below 0x80

goroutine 33 [running]:
example.com/tagged.SHA256.HashTagged(...)
	tagged.go:45
example.com/tagged/taggedtest.init.func1.1({0x0?, 0x0?}, {0xa6ba20?, 0xacd2e0?})
	taggedtest/tagged_suite.go:127 +0x37
go.dokimi.dev/conform.repeat[...](...)
	conform@v0.7.0/bench.go:123 +0xa7
go.dokimi.dev/conform.Bench[...](...)
	conform@v0.7.0/bench.go:71 +0x346
```

A sample of the type does not help. A sample on `Role` starts the generator of the type. Each case keeps the sample or draws any role. `HashTagged` and `CombineTagged` also take the two halves of the type, so one sample of `Role` serves one of them at most. `conform.Using` states one generator for each type of an application and has the same limit. Today the test of each hasher applies the suite of `crypto.Hasher` by hand. It writes its benchmarks and budgets again.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| A parameter with an authored sample takes the sample in every generated call of its callable | The suite generator |
| A parameter without a sample of its own takes a drawn value, as before | The suite generator |
| The sample of a type starts the generator of the type and fixes no parameter | The suite generator |
| A generated file refers only to declarations that it can import, also through the sample of a parameter | `Unimportable` of the suite generator |

### Components

| Component | Change |
|---|---|
| The suite generator | Reads the sample that a parameter states on itself, and reports `Unimportable` for a sample that the file of the suite cannot refer to |
| The template of the suites | Writes the sample of a parameter in place of its draw |
| The code `Unimportable` | A sample of a parameter that refers to a declaration that the file cannot import |
| The Go guide | A row for the sample of a parameter |

### The value of a parameter

The call generator gives a parameter its value by the first rule that applies:

1. A parameter in the role `context` takes the context of the case.
2. A parameter with a sample on itself takes the sample. The generator converts the sample to the type of the parameter, and the backend spells the conversion in the call generator of each case.
3. A parameter without a sample of its own takes a value by the rules of the values of parameters, as before.

The suite reads the sample of a parameter through the identity of the parameter in the node of its callable. It reads the sample that the parameter states on itself alone. The view of eidos falls back to the sample of the type, so the generator first reads whether the parameter states one.

A variadic parameter takes a sample of its element type and passes the sample as its only element. A sample on a parameter in the role `context` has no effect.

For the hasher of the motivation, the declaration states the sample of each role:

```go
type Hasher interface {
	HashTagged(
		//dokimi:sample value=0x01
		r Role,
		data []byte,
	) [32]byte

	CombineTagged(
		//dokimi:sample value=0x81
		r Role,
		left, right [32]byte,
	) [32]byte
}
```

The call generator of `HashTagged` takes the sample of `r` and draws `data`:

```go
var hasherHashTaggedCalls = prop.Composite(func(c *prop.Case) func(context.Context, tagged.Hasher) error {
	r := tagged.Role(1)
	data := c.Draw(prop.Of[[]byte](), "data")
	…
})
```

### The runs that the sample changes

| Run | The calls of a callable with a sample on a parameter |
|---|---|
| The benchmark and the parallel benchmark of a subject, and the benchmarks of the test of the suite | Take the sample in each call |
| The check of an allocation ceiling | Takes the sample in each case. The campaign of the property engine searches the other parameters alone |
| The planted defect of a ceiling | Takes the sample in each call |
| A check of a consumer | None. The check draws its own inputs |
| A fuzz target | None. The target draws the inputs of the property checks |
| A generated binding of a claim | Takes the sample in each call, because it draws the calls of the call generator |

`conform.Using` states the generator of a type for the values that a run draws. A parameter with a sample draws no value, so `Using` does not apply to it. A counterexample does not show the label of such a parameter, because the case draws no value for it.

### Failures

| Failure | Result | Recovery |
|---|---|---|
| The sample of a parameter refers to a declaration that the file of the suite cannot import | `Unimportable` at the declaration, and the suite records nothing for it | Write a sample that the file can spell, such as a value of an exported declaration |
| eidos refuses the text of the sample | The finding of eidos at the directive, as for the sample of a type | Correct the text |
| A sample on a parameter in the role `context` | No effect | Remove the sample |

### The contract in each language

| Part | Go | TypeScript |
|---|---|---|
| The sample of a parameter | `//dokimi:sample value=…` on the parameter | The same directive on the parameter |
| The value of the parameter in a generated call | The sample, spelled as a value of the type of the parameter | The same |

The TypeScript suites and their template tree do not exist yet. Their call generators follow the same rule when they do.

### Verification

- **The suite generator.** The specs of the plugin cover a parameter with a sample, a parameter whose type has a sample, a variadic parameter with a sample, a sample on a parameter in the role `context`, and a sample that the file cannot refer to. The golden files of a fixture pin the call generator that takes the sample.
- **The corpus.** A fixture of a hasher whose two methods refuse the two halves of a role. `go test -run '^$' -bench . -benchtime=100x` of the fixture ends without a panic, and the falsification of its checks passes.

### Migration

A declaration without a sample on a parameter generates the same files. A declaration with one generates call generators that take the sample. Its benchmarks and budgets measure other inputs from then on.

## Alternatives considered

### The sample of the type

The suite keeps the sample of the type as the source of a parameter, as today.

**Why not:** the sample of a type starts the generator of the type, and each case keeps the sample or draws any value. One type also serves parameters with different domains, such as the two halves of a role.

### A generator of each parameter

A param of a directive on the parameter refers to a generator of the kit, such as `values=leafRoles`. The call generator draws the parameter from it.

**Why not:** the sample already states a value that the contract admits, and eidos already stamps it on the parameter. A generator of each parameter adds a function of the kit for each refused domain. Reverses when a budget or a benchmark depends on the value of the narrowed parameter, so that one value measures too little.

### A generator of the application

The application passes the generator of each parameter, such as an option of `conform.Using` for one parameter.

**Why not:** the domain that the contract admits is a property of the declaration, so every subject of the suite would repeat the option.

### A recovered panic in the benchmarks

`Bench` recovers a panic of a call and counts the call as refused.

**Why not:** the benchmark would then hide a real panic of a subject, and it would measure the panic of each refused call.

## Drawbacks

- **One value of the parameter.** The benchmark and the budget of the callable measure one value of a parameter with a sample. An allocation that depends on that value goes outside the search of the campaign.
- **A generated binding sees one value.** A generated binding of a claim draws the calls of the call generator, so it takes the sample of the parameter too.
- **A counterexample without the parameter.** The case draws no value for a parameter with a sample, so a counterexample does not show it.

## Open questions

None.

## References

| What | Where |
|---|---|
| The request of the sample of a parameter | <https://github.com/dokimasia/dokimi/issues/36> |
| The annotator of the directive `sample` | <https://github.com/dokimasia/eidos/blob/main/eidos-core/authored/sample.go> |
