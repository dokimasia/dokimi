---
rfc: 0003
title: The suite generator
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-19
updated: 2026-10-19
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0003: The suite generator

## Summary

The suite generator writes the checks of each declaration that the directive `suite` marks. For each callable it reads the facts that the rules of eidos derive and the classifications of the shape catalog. It binds every claim of the conformance spec whose licence the callable has and whose gates pass. It lists every other licensed claim as a refusal, and every classification without a claim as unclaimed with a stance. It then writes a suite in the language of the declaration, which contains:

- a header
- the checks with their probe calls
- the entry points for the subjects of the consumer
- a planted defect for each check
- a falsification test and a structural test
- a lock file of check identities

The derivation is one function of the binary for every language.

## Motivation

The derivation produces the checks of a declaration. A declaration that the source states in two languages must get the same checks in both, or the claims of dokimi mean different things in different languages. So the derivation reads only facts that eidos derives in the same form for every language. These facts are the roles of parameters and returns, the error model, the classifications and the samples.

testkit derives the same kind of checks for Go, and its derivation reads Go directly. Its 19 gates test Go signature shapes, such as a `(value, error)` result, a `(value, bool)` result and a lone pointer result. The rows bind Go generic probes, and the derivation imports the Go SDK of an earlier eidos. A second language would need a second derivation. The rules of eidos now assign roles and error models in each language, so the gates can read those roles instead.

## Detailed design

### The directive

`dokimi:suite` marks one of three kinds of declaration:

- **An interface.** A Go or TypeScript interface, a Rust trait, a Python `Protocol` or abstract class, or a protobuf service.
- **A concrete type with methods.** A Go or Rust struct, or a class of TypeScript, Python, Java or Kotlin. The generator writes an interface of the method set, and the double and the planted defects implement it. In Go, TypeScript and Python a type satisfies an interface by its structure, so the concrete type satisfies the generated interface as it is. In Rust, Java, Kotlin, PHP and C#, the generator also writes an implementation of the interface that calls the type.
- **A function.** A check calls the function through a function value, and its planted defect is another function value.

The directive has no parameters. The kernel directive `out` routes the files of the suite, as it routes the files of every generator.

An interface or a concrete type with `suite` also gets the typed double of the stub generator, because the planted defects of its suite are built from that double. The suite generator requires the capability `stub`, and the stub generator provides it. A function does not need a double.

### The facts of a callable

The generator builds one value of facts for each callable of the declaration. Each part of it comes from eidos:

| Fact | Source in eidos |
|---|---|
| The parameters with the roles `input` and `context`, and the returns with the roles `value`, `ok-bool`, `stream` and `error` | `rules.Callable` from `Bound.CallableOf` |
| The error model: `none`, `last-return`, `result-type`, `thrown` or `raised` | `rules.Callable.Errors` |
| The canonical shape of each type, such as an integer, a list or the optional form | `Bound.TypeOf` |
| Two distinct samples of a type, or the reason why eidos refuses them | `Bound.SamplesOf` |
| The zero value of a type | `Bound.ZeroValue` |
| The shape, the mixins and the contract memberships, with their parameters | `shape.Of` |

```go
// Facts are the facts that eidos derives for one callable of a declaration.
type Facts struct {
	// Declaration and Callable are the names that the source declares.
	Declaration, Callable string

	// Signature is the normalized view of the callable.
	Signature rules.Callable

	// Class is the classification of the callable.
	Class shape.Classification

	// Inputs pairs each parameter in the role input with what eidos derives
	// for its type. Results does the same for each return in the role value.
	Inputs, Results []Value
}

// Value is one parameter or return with what eidos derives for its type.
type Value struct {
	Name    string
	Shape   rules.TypeShape
	Samples [2]rules.Sample
	Zero    emit.Value
	HasZero bool
}
```

### The claims table

The claims table is the file `claims.yaml` of the conformance spec, and the suite generator binds its claims whose generator is `suite`. The package `spec` decodes each claim into a `Claim`. A claim binds to a callable that has one of its licences and passes all its gates. The claim then becomes one check:

```go
// Claim is one claim of the conformance spec.
type Claim struct {
	// Segment identifies the claim, and it is the last part of the identity
	// of each check of the claim, such as zero-on-error.
	Segment Segment

	// Generator is the generator that binds the claim: suite, enum or
	// sentinel.
	Generator string

	// Licences are the classifications that license the claim. A callable
	// needs one of them. A claim of enum or sentinel has none.
	Licences []Licence

	// Unless is a mixin that withdraws the claim from a callable that has it.
	Unless shape.Mixin

	// Gates are the conditions that the facts of the callable must meet.
	Gates []GateID

	// Law is the law that the check witnesses.
	Law LawID

	// Machinery is what the run supplies before the check exercises its
	// claim.
	Machinery MachineryID

	// Statement is the sentence of the check. The generator replaces
	// {declaration} and {callable} with the names.
	Statement string

	// Probe is the probe that the check calls.
	Probe ProbeID

	// Inputs are the bindings of the inputs of the probe, in the order of the
	// probe.
	Inputs []Binding

	// Defect is the planted defect of the check.
	Defect Defect

	// Entails is the segment of a shallower claim that this claim implies,
	// and empty when the claim implies none.
	Entails Segment
}

// Licence is one classification that licenses a claim. Exactly one of
// Callable, Context, Shape, Mixin and Contract is set.
type Licence struct {
	// Callable licenses the claim on every callable.
	Callable bool

	// Context licenses the claim on every callable with a parameter in the
	// role context.
	Context bool

	Shape    shape.Shape
	Mixin    shape.Mixin
	Contract shape.Contract

	// Role is the role of the callable in Contract.
	Role shape.Role

	// With is a second classification that the callable must also have, for
	// a claim of two classifications together.
	With *Licence
}

// Gate is one gate of the conformance spec. A licensed claim whose gate
// fails is a refusal with the reason and the remedy of the gate.
type Gate struct {
	ID GateID

	// Passes is the condition of the gate, in words.
	Passes string

	Why    string
	Remedy string
}

// Binding is one input of a probe and the source of its value.
type Binding struct {
	// Input is the name of the input in the spec.
	Input string

	// From is the source of the value of the input.
	From Source
}

// Source is the fact of the subject from which a check takes the value of
// an input. The subject is a callable, an enumeration or the errors of a
// package.
type Source struct {
	Kind SourceKind

	// Param is the parameter of the classification for a partner, a sentinel,
	// a value or a duration.
	Param string

	// Role is the role of a partner in the contract of the callable.
	Role shape.Role

	// Sample is the index of a sample, 0 or 1.
	Sample int
}

// SourceKind is the kind of fact of a source.
type SourceKind uint8

const (
	// FromCallable is a call of the callable itself with sample inputs.
	FromCallable SourceKind = iota + 1

	// FromPartner is a call of the callable in a parameter of the
	// classification, such as the write of read-after-write.
	FromPartner

	// FromRole is a call of the callable with the role Role in the same
	// contract instance, such as the backing of a cache.
	FromRole

	// FromSample is a sample of the type of a result.
	FromSample

	// FromZero is the zero value of the type of a result.
	FromZero

	// FromParam is the value of a parameter of the classification, such as
	// the lifetime of a ttl.
	FromParam

	// FromClock is the controlled clock of the run.
	FromClock

	// FromMembers are the members of an enumeration with their names.
	FromMembers

	// FromName is the function that returns the name of a member of an
	// enumeration, and FromParse the function that parses a name into a
	// member.
	FromName
	FromParse

	// FromErrorValues are the error values of a package with their names.
	FromErrorValues

	// FromErrorTypes are the values of the error types of a package with
	// their names.
	FromErrorTypes

	// FromCauses are the functions of the error types of a package with a
	// cause. Each makes an error of its type with a given cause.
	FromCauses
)

// Defect is the planted defect of a claim.
type Defect struct {
	// Form is the identifier of the template of the defect in every template
	// tree, such as value-beside.
	Form DefectForm

	// Clause describes the defect. The generator replaces {declaration} and
	// {callable} with the names.
	Clause string

	// Reason is the reason code that the check must report against the
	// defect.
	Reason ReasonID
}
```

`shape.Shape`, `shape.Mixin`, `shape.Contract` and `shape.Role` are the types of the shape catalog. The package `spec` decodes each classification of a claim through the shape catalog. A test of the package fails on a name that the catalog lacks, and on a name of a law, gate, probe, reason code or machinery that the spec lacks. Each template tree has one template for each defect form of the claims, and a test of the template tree fails on a form without one.

### The gates

A gate reads the facts that the rules of each language decide. The derivation implements the condition of each gate of the spec once:

```go
// conditions contains the condition of each gate of the spec over the facts
// of a callable. A test fails on a gate of the spec without a condition.
var conditions map[spec.GateID]func(f *Facts) bool
```

The conditions of the gates:

| Gate | Passes when |
|---|---|
| `failure` | The error model of the callable is not `none` |
| `cancellation` | A parameter has the role `context` |
| `one-value` | The callable returns exactly one value |
| `two-values` | The callable returns exactly two values |
| `no-value` | The callable returns no value |
| `presence-flag` | A return has the role `ok-bool` |
| `optional-value` | The value result has the optional form |
| `list-value` | The value result is a list, or a return has the role `stream` |
| `ordered-value` | The value result is an integer, a float or a string |
| `value-samples` | eidos derives two distinct samples of each value result |
| `optional-input` | A parameter in the role `input` has the optional form |

Every claim also needs two samples of each input parameter, because a check calls the callable with sample inputs. The first sample of a type is its authored sample, else the default that the annotator `defaults` stamps on the type, else the sample that eidos derives. When eidos refuses a sample, each licensed claim on the callable refuses with the refusal reason of eidos. A claim with the licence `Callable` does not refuse when a gate fails, because no classification of the author licenses it. It is silent instead.

dokimi has no gate for equality. A probe compares two values with the equality of assert-spec, which compares any two values by structure and compares functions by identity.

The rules of each language assign the roles and forms that a gate reads. These constructs pass each such gate in each language:

| Gate | Go | TypeScript | Rust | Python | Protobuf | Java | Kotlin | PHP |
|---|---|---|---|---|---|---|---|---|
| `failure` | A last return of type `error` | Every callable | A `Result` return | Every callable | Every rpc under the gRPC binding of its target | Every method | Every function | Every method |
| `cancellation` | A `context.Context` parameter | An `AbortSignal` parameter | Never, because a caller cancels a future by dropping it | Never | Every rpc under the gRPC binding of its target | Never | Never, because a coroutine is cancelled through its job | Never |
| `presence-flag` | A `bool` as the second of two returns | Never | Never | Never | Never | Never | Never | Never |
| `optional-value` | A pointer result | A result of `T \| undefined` or `T \| null` | A result of `Option<T>` | A result of `T \| None` | Never | A result of `Optional<T>` | A result of `T?` | A result of `?T` |
| `list-value` | A slice, an array, or an `iter.Seq` result | An array, or an `AsyncIterable` result | A `Vec` or an iterator result | A `list` or an iterator result | A stream response | A `List` or a `Stream` result | A `List` or a `Flow` result | An `array` or an iterable result |

In eidos, the rules of Go and TypeScript assign the roles and forms of their columns. Protobuf gives every rpc the error model `none`, so the gRPC binding of a service must assign the roles and forms of its column. Each other language needs rules that assign the roles and forms of its column.

### The derivation

```go
// Derive returns the checks of one declaration, with its refusals, its
// unclaimed classifications and its entailment edges. tier and runtime are
// the machineries that the suite tier and the runtime of the target supply.
func Derive(decl []Facts, table []spec.Claim, tier, runtime []spec.MachineryID) Suite

// Suite is the derivation of one declaration.
type Suite struct {
	Checks    []Check
	Refusals  []Refusal
	Unclaimed []Unclaimed
	Edges     []Edge
}

// Check is one claim bound to one callable.
type Check struct {
	Identity string // Store.Read/zero-on-error
	Key      string // store.read/zero-on-error
	Claim    *spec.Claim
	Facts    *Facts
	Inputs   []Bound

	// Shallow is true when the law of the claim needs a deeper machinery
	// than the claim.
	Shallow bool
}

// Refusal is a licensed claim whose gate failed on a callable.
type Refusal struct {
	Claim       Segment
	Callable    string
	Why, Remedy string
}

// Unclaimed is a classification of a callable without a bound or a refused
// claim.
type Unclaimed struct {
	Classification string
	Callable       string
	Stance         spec.StanceID
	Note           string
}

// Edge is a pair of checks on one callable where the deeper check implies
// the shallower one.
type Edge struct {
	Deeper, Shallower string
}
```

`Derive` visits the callables in the order of the declaration and the claims in the order of `claims.yaml`:

1. For each claim with a licence of the callable and without its `Unless` mixin, it tests the gates. When every gate passes, it binds the inputs. When a binding needs a parameter or a role that the classification does not supply, the claim refuses with the name of the missing parameter. Each claim that binds becomes a check.
2. A claim whose machinery the tier or the runtime does not supply does not become a check. Its classification is unclaimed with the stance `needs`, unless another claim of the classification binds.
3. Each classification of the callable without a bound or a refused claim is unclaimed. Its entry in the section `unclaimed` of the conformance spec has its stance and note.
4. Each check whose claim entails another claim that bound on the same callable adds an edge.

The identity of a check is `<declaration>.<callable>/<segment>`. A claim that binds twice on one callable adds `#<id>`, where `<id>` is the parameter `id` of the contract instance or the index of the mixin instance. The key of a check follows the key grammar of the conformance spec.

The tests of `Derive` build facts as values and run without eidos or a language, because `Derive` is a pure function.

### The claims

The claims of the suite generator in the conformance spec restate the claims that testkit derives for Go over the facts of eidos. Every claim needs the machinery `nothing-beyond-a-caller`, except the six claims marked with a clock, which need `a-clock-the-run-advances`.

| Claim | Licence | Gates | Probe | Reason code | Law |
|---|---|---|---|---|---|
| `smoke` | Every callable | — | `survives` | `panicked` | `no-panic-on-any-input` |
| `cancel` | A parameter in the role context | cancellation, failure | `under-cancelled` | `ignored-cancel` | `a-cancelled-context-is-reported` |
| `zero-on-error` | Shapes reader, aggregator, answering-writer and stream-consumer, without the mixin default-on-error | failure, one-value, value-samples | `zero-beside-error` | `value-beside-error` | `no-value-beside-an-error` |
| `zeroes-on-error` | Shapes multi-reader and multi-aggregator | failure, two-values, value-samples | `zero-pair-beside-error` | `value-beside-error` | `no-value-beside-an-error` |
| `zero-on-miss` | Shape reader-with-bool | presence-flag, one-value, value-samples | `zero-beside-false` | `value-beside-false` | `no-value-beside-a-false-flag` |
| `zeroes-on-miss` | Shape lookup | presence-flag, two-values, value-samples | `zero-pair-beside-false` | `value-beside-false` | `no-value-beside-a-false-flag` |
| `nil-on-miss` | Shape pointer-reader | one-value, optional-value | `nil-on-miss` | `non-nil-miss` | `a-miss-returns-the-zero` |
| `miss-is-zero` | Shape reader-no-error | one-value, value-samples | `miss-is-zero` | `non-zero-miss` | `a-miss-returns-the-zero` |
| `deterministic` | Shape computation | one-value, value-samples | `repeat-agrees` | `unstable-repeat` | `a-repeat-returns-alike` |
| `predicate` | Shape predicate | one-value | `repeat-agrees` | `unstable-repeat` | `a-repeat-returns-alike` |
| `repeat-survives` | Shape void-lifecycle | no-value | `survives-repeat` | `panicked-repeat` | `a-repeated-call-survives` |
| `miss` | Mixin not-found | failure, one-value | `reports-sentinel` | `wrong-sentinel` | `a-miss-reports-the-declared-sentinel` |
| `idempotent` | Mixin idempotent | cancellation, failure | `repeat-is-quiet` | `refused-repeat` | `a-repeat-is-accepted` |
| `cacheable` | Mixin cacheable | one-value | `repeat-agrees` | `unstable-repeat` | `a-repeat-returns-alike` |
| `pure` | Mixin pure | one-value | `repeat-agrees` | `unstable-repeat` | `a-repeat-returns-alike` |
| `read-after-write` | Mixin read-after-write | cancellation, one-value | `write-then-read` | `unobserved-write` | `a-read-returns-the-write-before-it` |
| `causal` | Mixin causal | cancellation, one-value | `write-then-read` | `unobserved-write` | `a-read-observes-its-causal-past` |
| `read-your-writes` | Mixin read-your-writes | cancellation, one-value | `write-then-read` | `unobserved-write` | `a-read-observes-its-own-writes` |
| `monotonic-writes` | Mixin monotonic-writes | cancellation, one-value | `write-then-read` | `unobserved-write` | `writes-apply-in-issue-order` |
| `point-in-time` | Mixin point-in-time | cancellation, one-value | `pinned-across-write` | `unpinned-read` | `a-read-keeps-its-snapshot` |
| `no-duplicates` | Mixin no-duplicates | cancellation, list-value | `no-element-twice` | `repeated-element` | `no-element-is-drained-twice` |
| `stable-order` | Mixin stable-order | cancellation, list-value | `drains-in-one-order` | `reordered-drain` | `two-drains-agree-on-order` |
| `permutation` | Mixin permutation | cancellation, list-value | `gains-what-went-in` | `miscounted-drain` | `a-drain-gains-exactly-what-went-in` |
| `over-match` | Mixin over-match | cancellation, list-value | `keeps-what-went-in` | `shrunk-drain` | `a-drain-keeps-what-went-in` |
| `monotonic` | Mixin monotonic | one-value, ordered-value, value-samples | `never-decreases` | `regressed` | `results-never-decrease` |
| `monotonic-reads` | Mixin monotonic-reads | cancellation, one-value, value-samples | `reads-never-regress` | `regressed` | `reads-never-go-backwards` |
| `writes-follow-reads` | Mixin writes-follow-reads | cancellation, one-value, value-samples | `writes-advance` | `unordered-write` | `a-write-follows-the-read-before-it` |
| `xss-safe` | Mixin xss-safe | cancellation, one-value | `escapes-the-unsafe` | `unescaped-output` | `the-unsafe-value-does-not-survive` |
| `injection-safe` | Mixin injection-safe | cancellation, failure | `write-then-read` | `unobserved-write` | `the-unsafe-value-is-kept-as-data` |
| `total` | Mixin total | — | `accepts-the-edge` | `refused-the-edge` | `the-call-accepts-its-declared-domain` |
| `validates` | Mixin validates | cancellation, failure | `refuses-the-invalid` | `unscreened-write` | `the-declared-invalid-value-is-refused` |
| `nil-safe` | Mixin nil-safe | optional-input | `survives` | `panicked` | `a-nil-argument-is-accepted` |
| `lifecycle-after-close` | Mixin lifecycle-after-close | cancellation, failure | `reports-sealed` | `alive-after-close` | `work-after-close-is-refused` |
| `order-after` | Mixin order-after | cancellation, failure | `refuses-early` | `accepted-early` | `the-call-refuses-before-its-predecessor` |
| `ready-after` | Mixin order-after | cancellation, failure | `works-after` | `never-ready` | `the-call-works-after-its-predecessor` |
| `settles` | Mixin eventually | cancellation, failure | `settles-then-observes` | `unsettled-effect` | `an-effect-arrives-once-settled` |
| `within` (clock) | Mixin timeout | cancellation, failure | `returns-within` | `overran-budget` | `a-call-returns-within-its-budget` |
| `windowed` (clock) | Mixin windowed | cancellation, failure | `counts-within-window` | `misread-window` | `a-count-covers-a-bounded-window` |
| `scheduled` (clock) | Mixin scheduled | cancellation, failure | `runs-on-schedule` | `off-schedule` | `a-task-runs-once-its-instant-passes` |
| `time-aware` (clock) | Mixin time-aware | cancellation, one-value | `moves-with-the-clock` | `unmoved-by-the-clock` | `a-result-depends-on-the-clock` |
| `steady` (clock) | Mixin time-aware | cancellation, one-value | `steady-while-the-clock-stops` | `drifted-without-clock` | `a-result-is-stable-while-the-clock-stops` |
| `expires` (clock) | Mixin ttl | cancellation, one-value | `expires-after` | `not-expired` | `an-entry-lapses-on-the-clock` |
| `delete-removes` | Mixin delete-removes | cancellation, failure | `removes-on-delete` | `readable-after-delete` | `a-delete-makes-the-key-unreadable` |
| `side-effect` | Mixin side-effect | cancellation, failure | `moves-observation` | `unmoved-observation` | `the-call-moves-the-named-observation` |
| `poisonable` | Mixin poisonable | failure | `remains-poisoned` | `healthy-after-induce` | `an-induced-failure-latches` |
| `tamper-evident` | Mixin tamper-evident | cancellation, failure | `detects-tamper` | `silent-tamper` | `a-tamper-is-detected` |
| `retry-succeeds-tx` | Mixin retry-succeeds with contract tx | cancellation, failure | `converges-after` | `never-converged` | `the-call-converges-once-its-scope-is-open` |
| `retry-succeeds` | Mixin retry-succeeds | cancellation, failure | `succeeds-within` | `never-converged` | `the-call-converges-within-the-declared-attempts` |
| `hooks` | Mixin hooks | cancellation, failure | `runs-the-hook` | `hook-not-run` | `a-registered-hook-runs` |
| `accumulates` | Mixin accumulates | cancellation, failure | `each-call-compounds` | `uncompounded` | `repeated-calls-compound` |
| `partition` | Mixin partition | cancellation, failure | `keeps-partitions-apart` | `crossed-partitions` | `one-partition-does-not-serve-another` |
| `persister` | Contract persister, role writer | cancellation, one-value | `write-then-read` | `unobserved-write` | `a-read-returns-the-write-before-it` |
| `chain-causal` | Contract chain, role append, with mixin causal | cancellation, failure | `replays-in-stamp-order` | `unordered-history` | `a-history-replays-in-its-stamped-order` |
| `chain` | Contract chain, role append | cancellation, failure | `replay-grows` | `replay-not-grown` | `an-admitted-entry-grows-the-replay` |
| `appender` | Contract appender, role fn | one-value, ordered-value, value-samples | `repeat-ascends` | `not-ascending` | `an-append-returns-above-the-last` |
| `codec` | Contract codec, role forward | cancellation, one-value | `round-trips` | `broken-round-trip` | `an-encoding-survives-its-round-trip` |
| `if-absent` | Contract if-absent, role writer | cancellation, failure | `refuses-the-second-write` | `unguarded-write` | `a-write-is-refused-where-its-key-exists` |
| `lease` | Contract lease, role acquire | cancellation, failure | `refuses-the-second-write` | `unguarded-write` | `a-granted-lease-is-not-granted-twice` |
| `cas` | Contract cas, role writer | cancellation, failure | `refuses-the-second-write` | `unguarded-write` | `a-write-is-refused-on-a-stale-version` |
| `if-match` | Contract if-match, role writer | cancellation, failure | `refuses-what-the-predicate-rejects` | `unmatched-write` | `a-write-is-refused-where-the-predicate-rejects-it` |
| `cursor` | Contract cursor, role next | cancellation, failure | `reports-sealed` | `alive-after-close` | `work-after-close-is-refused` |
| `cursor-open` | Contract cursor, role open | cancellation, failure | `seals-after-close` | `alive-after-close` | `work-after-close-is-refused` |
| `rate-limit` | Contract rate-limit, role fn | cancellation, failure | `refuses-past-the-burst` | `past-the-burst` | `the-burst-is-spent-then-refused` |
| `leader-election` | Contract leader-election, role campaign | cancellation, failure | `resigns-leadership` | `still-leading` | `resigned-leadership-is-released` |
| `transaction` | Contract transaction, role fn | cancellation, failure | `rolls-back-on-failure` | `not-rolled-back` | `a-failed-scope-leaves-nothing` |
| `cache` | Contract cache, role cache | cancellation, one-value | `agrees-with-the-backing` | `cache-disagrees` | `the-cache-agrees-with-its-backing` |
| `tx` | Contract tx, role begin | cancellation, failure | `seals-after-close` | `alive-after-close` | `a-settled-transaction-is-not-settled-again` |
| `pool` | Contract pool, role get | cancellation, one-value | `takes-back-what-it-lent` | `unreturned-loan` | `a-borrowed-value-is-taken-back` |
| `pagination` | Contract pagination, role reader | cancellation, one-value | `continues-from-the-cursor` | `repeated-page` | `a-walk-resumes-where-it-left-off` |
| `saga` | Contract saga, role step | cancellation, failure | `compensates-the-step` | `uncompensated-step` | `a-compensation-accepts-its-step` |
| `batch-writer` | Contract batch-writer, role writer | cancellation, failure | `refused-leaves-nothing` | `partial-write` | `a-refused-write-leaves-nothing` |
| `updater` | Contract updater, role writer | cancellation, failure | `second-write-replaces` | `unreplaced-write` | `a-second-write-replaces-the-first` |
| `upserter` | Contract upserter, role writer | cancellation, failure | `second-write-replaces` | `unreplaced-write` | `a-repeated-write-keeps-one-key` |

The table has 73 claims, 59 probes, 54 reason codes and 64 laws. The claim `repeat-survives` entails `smoke`.

### The stances of unclaimed classifications

Each classification without a claim has a stance and a note in the section `unclaimed` of the conformance spec:

| Stance | Classifications |
|---|---|
| `silent` | The shapes writer, composite-writer, multi-arg-writer, mutator, deleter, closer and lifecycle. A write is an input of the other checks, and a call that returns nothing has no result to judge. The documentary mixins deprecated, errors and scope, and the mixins integration-only and sample |
| `needs` | The mixins concurrent, concurrent-readers, atomic, serializable and snapshot-isolation, and the contract single-flight, which need concurrent callers. The mixins associative, commutative, bounded, conservative, crdt-merge and stream-reflects-mutations, and the contracts outbox, publisher, watcher and workflow, which need generated sequences. The contract circuit-breaker, which needs a call that fails on demand |
| `blocked` | The mixin default-on-error, because it declares no default value for a check to compare |
| `gap` | The shapes batch-reader, stream-reader and poison-accessor, and the mixins indexed, leak-free, sticky and wrapped-via |

### The generated suite

The generator writes three families of files for each declaration. The template tree of each target spells them:

| Family | Per | Contents |
|---|---|---|
| The suite | Source file | The header, the compatibility reference, the checks and the entry points |
| The proofs | Source file, as a test | The planted defects, the falsification test and the structural test |
| The lock | Package | The keys of every check of the package, one per line, sorted |

The header lists:

- the checks of each callable, with their statements
- the checks whose law needs a deeper machinery, which witness the law at its shallow end
- the refusals, each with its reason and remedy
- the unclaimed classifications, each with its stance and note

The suite has two entry points for the consumer:

- `Run` applies every check to a fresh subject from a factory of the consumer.
- `Control` applies every check to a correct implementation that differs from the subject. Every check that the control does not excuse must pass on it.

Both accept an excusal of a check by its identity, with a reason, and a factory that pairs a fresh subject with a controlled clock. A check that needs a clock fails, naming the missing option, when the consumer did not supply that factory.

Each planted defect is the typed double of the declaration, with the callables that the defect form names replaced. The double is not a correct implementation of the declaration. A planted defect is correct for the calls that its check makes, and wrong at the claim of that check. Each defect form has one template in each template tree.

The falsification test runs in the package of the generated suite and does not need a subject of the consumer:

1. It checks that every check has a planted defect, and that every planted defect belongs to a check.
2. It runs each check against its planted defect, and it fails unless the check fails with the reason code of the claim.
3. It runs each edge in both directions. The deeper check must fail on the defect of the shallower check, and the shallower check must pass on the defect of the deeper check.

The structural test checks that the lock lists every check and nothing else, that every identity is unique and that every check has its law and its machinery.

### The lock

The lock is a source file of the target language with the keys of the checks of a package, one per line, sorted. A regeneration that loses a check removes one line from the lock. The generation gate fails while the committed lock differs from a fresh run, so a reviewer sees a lost check as a deleted line in a small file. A removal is legitimate when that deleted line is in the same change as the edit of the declaration that caused it.

### Invariants

| Invariant | Enforced by |
|---|---|
| A claim binds only on a callable with one of its licences, and only when all its gates pass | `Derive`, and its tests over facts built as values |
| A licensed claim that does not bind is a refusal with a reason and a remedy | `Derive` |
| Every classification of a callable is read by a claim or has a stance | `Derive`, and the test of the package spec that walks every classification of the shape catalog |
| The derivation reads no fact of a single language | depguard: `internal/derive` imports no package of a language satellite |
| Every check comes with a planted defect, and every defect belongs to a check | The falsification test of each generated suite |
| A check fails on its planted defect with the reason code of its claim | The falsification test of each generated suite |
| A lost check changes the committed lock | The lock and the generation gate |
| Every claim refers to an existing law, gate, probe, reason code, machinery and classification | The test of the package spec |
| Every gate of the spec has a condition in the derivation | A test of the derivation |

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| The language of the declaration has no rules in eidos | eidos reports the absent rules and derives no samples, so no claim binds on a callable with inputs | The rules of the language in eidos |
| eidos refuses a sample of an input parameter | Every licensed claim on the callable refuses with the reason of eidos | Author a sample with the directive `sample` |
| A gate of a licensed claim fails | The header lists the refusal with its reason and remedy | Change the declaration, or accept the refusal |
| A binding reads a parameter that the classification does not set | The claim refuses and names the parameter | Set the parameter in the directive |
| Two callables of the declaration have one key | Generation fails and names both callables | Rename one callable |
| The target lacks the template of a defect form | The build of the composition fails | Add the template to the template tree |
| A gate of the spec has no condition in the derivation | A test of the derivation fails and names the gate | Implement the condition |
| A planted defect passes its check | The falsification test fails and names the check | Fix the probe, the defect form or the claim |
| The committed lock differs from a fresh run | The generation gate fails | Run `dokimi run`, review the lock and commit |

## Alternatives considered

### The signature shapes of Go

The gates read the shapes of Go signatures, as the gates of testkit do. Such shapes are a `(value, error)` result, a `(value, bool)` result and a lone pointer result.

**Why not:** these shapes exist only in Go. A `(value, bool)` result has no counterpart in TypeScript, Rust or Python, so each other language would need gates of its own. A gate that reads the roles and forms of eidos works in every language whose rules assign them.

### Planted defects that wrap the subject of the consumer

A planted defect delegates every call to the subject of the consumer. It breaks only the callable of the check.

**Why not:** falsification then needs the subject of the consumer, so it runs in the tests of the consumer. A partner of the check can behave differently in that subject. The check would then fail falsification for a reason outside the check. A defect built from the double runs in the package of the generated suite. Its falsification depends on nothing outside dokimi.

### Planted defects by mutation of the subject

A mutation tool such as dokimi-mutate-go changes the code of the subject. Falsification expects each check to fail on some mutant.

**Why not:** a mutant does not target a claim. A check that fails on a mutant does not show that it fails for the reason of its claim. Mutation also needs the source of the subject and a mutation engine for each language. dokimi-mutate-go remains a separate measure of the code of the consumer.

### Checks with argued evidence

A check may come without a planted defect when its author argues that it can fail. testkit allows this with `Because`.

**Why not:** dokimi exists to remove checks that no planted defect proves able to fail. The structural test covers the wiring that testkit argues for. That wiring is not a claim.

### A gate for equality

The gates include the equality of each language. testkit requires a result type that compares with `==` for each claim that compares results.

**Why not:** eidos does not derive the equality of a type in each language. The gate would need a rule for each language in the derivation. assert-spec fixes one structural equality for every language, and each probe compares through it.

### The lock as a separate text file

The lock is a text file beside the generated files. testkit proposed such a file as `checks.lock`.

**Why not:** eidos writes the files of a plan through the backend of its target. Each backend spells a source file of its language. A lock in the language of the target is one more file family, and its structural test reads the lock as code.

## Drawbacks

- **Templates for every defect form.** testkit uses 36 defect forms. With eight targets, that is 288 defect templates, and a new defect form needs a template in every tree.
- **Every runtime implements 59 probes.** The probes are host code, so each runtime writes and tests all 59.
- **A double for every suite.** A declaration with `suite` also gets the typed double. For the corpus fixture of testkit with one method, the suite and its proofs have 392 lines and the double with its test has 465.
- **Structural equality.** Probes compare results by structure. A type that its language compares by identity, such as a class in TypeScript, compares by structure in a check.
- **Claims that need parameters.** A partner claim refuses until its directive sets the partner, so a mixin without its parameters adds refusals to the header.

## Open questions

None.

## References

| What | Where |
|---|---|
| The claims table of testkit | <https://github.com/thesm-os/testkit/blob/main/generator/suite/claims.go> |
| The shape catalog of eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/12-shape-catalog.md> |
| The callable that the rules of eidos derive | <https://github.com/dokimasia/eidos/blob/main/eidos-core/rules/callable.go> |
| The equality of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/assertions.yaml> |
