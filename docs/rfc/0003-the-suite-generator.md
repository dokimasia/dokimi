---
rfc: 0003
title: The suite generator
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-19
updated: 2026-10-22
discussion: none
supersedes: none
superseded-by: ADR-0001, ADR-0002 and ADR-0006 in part
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
- the entry points for the subjects and the checks of the consumer
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

The generator writes three families of files. The template tree of each target spells them:

| Family | Per | Contents |
|---|---|---|
| The suite | Source file of the declaration | The header, the compatibility reference, the checks with their planted defects, and the application |
| The tests | Source file of a declaration or of a subject, as a test | The proofs of the suite of a declaration: the falsification test, the test of the two forms of its application and the structural test. The conformance test of a subject that the directive `subject` declares: the application of the suite to the subject |
| The lock | Package of the suite | The keys of every check of the package, one per line, sorted |

eidos gives a plugin one family of test files, so the proofs and the conformance tests share it. The family has the word `suite` and tests its package from outside, so a Go file of the family is `<file>_suite_test.go` in the external test package of its directory. The proofs of `Store` are the spec of `<file>_suite.go` in the package of its suite. The conformance test of a subject is in the directory of the subject. Where the factory of a subject is not exported, the test of the subject is in `<file>_suite_internal_test.go` inside the package of the subject, because only that package can call the factory.

The header lists:

- the checks of each callable, with their statements
- the checks whose law needs a deeper machinery, which witness the law at its shallow end
- the refusals, each with its reason and remedy
- the unclaimed classifications, each with its stance and note

The suite has one application for the consumer. It applies every check to a fresh subject from a factory of the consumer. Its two forms differ in the factory alone: one form takes a constructor without arguments, and the other takes a factory that builds the subject from the handle of the check. Both forms accept the checks of the consumer, an excusal of a check by its identity with a reason, and the doors of the machineries that only the consumer can build.

The corpus applies each suite to a control through the same application. A control is a correct implementation that differs from the subject, so a check that fails on the control has a defect of its own.

Each planted defect is the typed double of the declaration, with the callables that the defect form names replaced. The double is not a correct implementation of the declaration. A planted defect is wrong at the claim of its check and correct for the calls of every other check of the suite. Each defect form has one template in each template tree.

The falsification test runs in the package of the generated suite and does not need a subject of the consumer:

1. It checks that every check has a planted defect, and that every planted defect belongs to a check.
2. It runs each check against its planted defect. A derived check must fail with the reason code of its claim. A check of the consumer must fail at the assertion whose message is its statement.
3. It runs every other check of the suite against that defect, and each of them must pass. The checks that share the defect are exempt, and so is each deeper check of an edge of the check. A check that ends with a fault on the defect, such as a check that calls a door without its option, has no verdict on it. The failure lists the defect and each other check that fails on it, with its first failing assertion or the value of its panic.
4. It runs each edge in both directions. The deeper check must fail on the defect of the shallower check, and the shallower check must pass on the defect of the deeper check.

The test of the two forms of the application applies a check of its own through each form. The factory of the test builds a value of a type that only the test declares, and the check compares the subject of the check with that value by identity. The test excuses each check of the suite, so the check of the test runs alone. The form `Run` of the suite of a function takes the function itself, so its check compares the subject with the function. The identity of a check of the consumer of a function has no callable, so the generator refuses the segments `builds-the-subject` and `receives-the-function` of the checks of the test on it.

In Go the proofs are the test of the suite, with a level for the suite and one for each form of its application. The level of the suite is one call of the runtime:

```go
// storeSubject is a value of Store that only TestStoreSuite builds.
type storeSubject struct{ store.Store }

// TestStoreSuite checks the suite of Store and the two forms of its
// application.
func TestStoreSuite(t *testing.T) {
	t.Parallel()

	t.Run("StoreSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, StoreSuite())
	})

	t.Run("RunStoreSuite", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		RunStoreSuite(
			t,
			func() *storeSubject {
				built = &storeSubject{}
				return built
			},
			conform.With(StoreCheck{
				ID:        "Store/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *StoreHandle) {
					got := h.Subject()
					assert.Equal(h, got, store.Store(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*StoreHandle) store.Store { return &storeSubject{} },
			}),
			conform.Excuse(
				StoreChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStoreSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		RunStoreSuiteFromHandle(
			t,
			func(*StoreHandle) *storeSubject {
				built = &storeSubject{}
				return built
			},
			// The same check and the same excusals as RunStoreSuite.
		)
	})
}
```

The structural test checks that the lock lists every check and nothing else, that every identity is unique and that every check has its law and its machinery.

### The contract of a suite

A consumer codes against these parts of a suite. The runtime of each language declares their types, and the generated file declares their names for one declaration:

| Part | Go | TypeScript |
|---|---|---|
| The names of the callables and of the checks | `StoreChecks.Read`, which embeds `conform.Callable[store.Store]` and has a method for each check, such as `StoreChecks.Read.ZeroOnError()` | `storeChecks.read`, a `Callable<Store>` with a function for each check, such as `storeChecks.read.zeroOnError()` |
| The checks | `StoreSuite()`, which returns a `conform.Suite[store.Store]` | `storeSuite()`, which returns a `Suite<Store>` |
| The application to a subject | `RunStoreSuite(t, newSubject, opts...)` and `RunStoreSuiteFromHandle(t, newSubject, opts...)` | `runStoreSuite(newSubject, options)` and `runStoreSuiteFromHandle(newSubject, options)` |
| The options of an application | `conform.With`, `conform.Excuse`, `conform.Reference`, `conform.Second` and `conform.Restartable` | The fields `checks`, `excuse`, `reference`, `second` and `restart` of `Options<S>` |
| The names of the types | `StoreHandle`, `StoreCheck` and `StoreOption` | `StoreHandle`, `StoreCheck` and `StoreOptions` |

The contract does not depend on the claims that the derivation binds. A declaration without a bound claim has the same parts, and its suite runs the checks of the consumer alone.

Each kind of declaration has its subject type:

- An interface is its own subject type.
- For a concrete type such as the struct `Store`, the generator writes the interface `StoreMethods` of its exported methods, and that interface is the subject type.
- For a function such as `Slug`, the generator writes the type `SlugFunc` of the function. An application takes the function itself, as in `RunSlugSuite(t, slug.Slug)`.

#### The application

`RunStoreSuite` takes a constructor without arguments. `RunStoreSuiteFromHandle` takes a factory that builds each subject from the handle of its check, such as on its clock or in its temporary directory. The words `StoreSuite` in each name refer to the suite that `StoreSuite` returns. The suffix `FromHandle` refers to the argument of the factory. Both functions are generic over the type of the subject, so a constructor such as `NewInMemory`, which returns `*InMemory`, does not need a conversion:

```go
func TestInMemory(t *testing.T) {
	t.Parallel()

	storetest.RunStoreSuite(t, storetest.NewInMemory,
		conform.Excuse(storetest.StoreChecks.Read.Cancel(), "Read reads memory and never blocks"))
}

func TestCache(t *testing.T) {
	t.Parallel()

	storetest.RunStoreSuiteFromHandle(t, func(h *storetest.StoreHandle) *ttl.Cache { return ttl.New(h.Clock()) })
}
```

A TypeScript test calls the generated function inside a `describe` block:

```ts
describe('InMemory', () => {
  runStoreSuite(() => new InMemory(), {
    excuse: { [storeChecks.read.cancel()]: 'read reads memory and never blocks' },
  });
});
```

`RunStoreSuite` calls `RunStoreSuiteFromHandle` with a factory that does not read the handle. `RunStoreSuiteFromHandle` runs every check as a test of its own through `conform.Run` (RFC-0002). The name of the test is the identity of the check without the name of the declaration, such as `Read/zero-on-error`, because the name of the test of the subject contains the declaration. `runStoreSuite` registers each check with vitest in a `describe` block of the declaration.

The subject type of its suite types each identity: `conform.ID[store.Store]` in Go and `Id<Store>` in TypeScript. An option takes its subject type from the identity, and the identity of another suite fails the compile.

#### The handle

A check takes a handle. The handle implements the seat of the assert library: `assert.TB`, `assert.Reporter` and `assert.Clocked` in Go, and `Seat` in TypeScript. Its Go method `Unwrap` returns the seat of the check, so each assertion of the check writes its call record to that seat. The handle also builds the subjects of the check and supplies its machinery:

| Method | Returns | Machinery |
|---|---|---|
| `Subject()` | The subject of the check, built once | Every check |
| `New()` | A fresh subject from the same source | `generated-sequences` and `concurrent-callers` |
| `Reference()` | A fresh reference implementation | `a-reference-implementation` |
| `Second()` | A fresh second subject | `a-second-subject` |
| `Restart(s)` | The subject `s` after a restart of its process | `a-killable-process-or-breakable-medium` |
| `Clock()` and `Advance(d)` | The controlled clock of the check, and its advance | `a-clock-the-run-advances` |
| `Statement()` | The statement of the check | Every check |

Under falsification `Subject()` and `New()` return the planted defect of the check. The consumer supplies the reference, the second subject and the restart through the options. A check that calls a door without its option ends with a fault that wraps `conform.ErrOption`. The message of the fault contains the option to pass.

A check that advances the clock ends with a fault that wraps `conform.ErrClock` when no factory of its subjects read the clock of the handle. The message of the fault contains the call of `RunStoreSuiteFromHandle` that builds the subject on the clock.

#### The checks of a consumer

A consumer declares a check with the directive `check` on a function of the package of the generated suite. The suite generator binds it as it binds a claim:

```go
// KeepsTheBody checks that Replay returns the body that Append stored.
//
//dokimi:check Log.Replay defect=losesBodies
func KeepsTheBody(h *LogHandle) {
	s := h.Subject()
	assert.NoError(h, s.Append(h.Context(), causalchain.RootEntry()), "Append admits the root entry")
	got, err := s.Replay(h.Context())
	assert.NoError(h, err, "Replay reads the log")
	assert.Length(h, got, 1, "Replay returns the root entry")
	assert.Equal(h, got[0].Body, causalchain.RootEntry().Body, h.Statement())
}
```

- The segment is the name of the function in kebab case, and the identity is `Log.Replay/keeps-the-body`. The statement is the first sentence of the doc comment after the words "checks that".
- `defect` names a function of the package that returns the planted defect of the check from its handle.
- `Log` is the name of a declaration of the plan with the directive `suite`. Where two such declarations have the name, the directive qualifies it with the name of its package, as in `causalchain.Log.Replay`.
- The check is part of the suite. The header lists it, the lock contains its key, the generated falsification test applies it to its defect, and `LogChecks.Replay.KeepsTheBody()` returns its identity.
- The generator reports each of these mistakes at the position of the directive:
  - a name that no declaration with the directive `suite` has, or that two of them have
  - an unknown callable
  - a segment that a derived check of the callable has
  - a defect without a function, or with another signature
  - a function without a doc comment
  - the segment `builds-the-subject` or `receives-the-function` on a check of a function, because a check of the test of the two forms has each of them

The suite generator binds a check in its two steps:

1. Its annotate step runs on each function with the directive `check`. It resolves the declaration, the callable and the defect. It stamps facts on the function: the identity of the declaration, the name of the callable, the segment, the statement and the identity of the defect.
2. The generate step of a declaration lists the functions with the fact of a declaration through `sdk.Subjects` of eidos, and reads the facts of each function through `sdk.FactOf`. It binds each function whose fact contains the identity of the declaration.

`sdk.Subjects` records a read of the set of functions with the key of the fact. eidos then runs the generate step of each suite again when a function gains or loses the fact. A change of any other declaration does not run it again.

A check in a test file cannot be bound, because the generated suite cannot refer to a test file. `conform.With` adds such a check to one application. The application falsifies the check against every check that it does not excuse before it applies the check to the subject.

The decisive assertion of a check passes `h.Statement()` as its message. Falsification passes when the first failure record on the planted defect has the statement as its contract. A failure of another assertion, such as a guard, fails the falsification. Falsification also fails when another check of the suite fails on the defect.

A run fails when two of its checks have one identity, and the failure message contains both identities. The same failure occurs when a release derives a check with the identity of a check of the consumer.

The usual planted defect is the verified fake of the declaration with one method replaced. The double of the stub generator forwards each call to the fake through `DelegateTo`, and the defect replaces one method:

```go
func losesBodies(h *LogHandle) causalchain.Log {
	fake := NewLogFake()
	return NewLogStub(h, LogStubDelegateTo(fake),
		WithLogAppend(func(ctx context.Context, e causalchain.Entry) error {
			e.Body = ""
			return fake.Append(ctx, e)
		}))
}
```

#### The subjects of a consumer

A consumer declares a subject with the directive `subject` on its type, and an excusal with the directive `excuse`:

```go
// Store is the Log on PostgreSQL.
//
//dokimi:subject causalchain.Log new=NewTestStore
//dokimi:excuse Log.Replay/cancel reason="the test server finishes a cancelled query"
type Store struct{ … }
```

The suite generator writes a conformance test of the subject into the external test package of the subject:

```go
func TestStoreConformsToLog(t *testing.T) {
	t.Parallel()

	causalchaintest.RunLogSuite(
		t,
		pgstore.NewTestStore,
		conform.Excuse(
			causalchaintest.LogChecks.Replay.Cancel(),
			"the test server finishes a cancelled query",
		),
	)
}
```

A factory without parameters selects `RunLogSuite`, and a factory whose one parameter is the handle selects `RunLogSuiteFromHandle`. The test refers to the function as a declaration that the plan emits, which the layout of eidos qualifies by the package of its file. The generator reports an excusal whose identity matches no check of the suite at the position of the directive, with the identities of the checks of the callable. A subject whose wiring a directive cannot state calls `RunLogSuite` or `RunLogSuiteFromHandle` in a test of its own.

#### Excusals

An excusal contains the identity of a check and a reason. The check reports `excused` with the reason, and never `pass`. A run fails when the identity of an excusal does not match a check of the run. So an excusal fails once a release no longer derives its check.

#### Failure messages

A failure of the contract contains the code that fixes it, and the command that runs the check again with the seed of the run:

```text
--- FAIL: TestStoreConformsToLog/Append/expires
    Log.Append/expires advances a clock that the subject never read. Build the subject on the clock of the handle:
        RunLogSuiteFromHandle(t, func(h *conform.Handle[causalchain.Log]) causalchain.Log { … h.Clock() … })
    Run it again:
        DOKIMI_SEED=0x5f3a1c9e27d4b860 go test -run 'TestStoreConformsToLog/Append/expires' ./pgstore
```

The message does not contain an excusal. A test that the directive `subject` generates takes the directive `excuse`, and a test of the consumer takes the option `conform.Excuse`. Both tests make the same calls of the runtime.

#### The generated file

The Go file of `Store` declares these names:

```go
// The names that a test of Store writes.
type (
	StoreHandle = conform.Handle[store.Store]
	StoreCheck  = conform.Check[store.Store]
	StoreOption = conform.Option[store.Store]
)

// StoreChecks contains the name of each callable of Store, with the
// identities of its checks.
var StoreChecks = struct {
	Read  storeRead
	Write storeWrite
}{Read: storeRead{"Store.Read"}, Write: storeWrite{"Store.Write"}}

// storeRead is the name of Store.Read, with the identities of its checks.
type storeRead struct{ conform.Callable[store.Store] }

// ZeroOnError returns the identity of the check that Read returns the zero
// value beside an error.
func (c storeRead) ZeroOnError() conform.ID[store.Store] { return c.Check("zero-on-error") }

// StoreSuite returns the checks of Store.
func StoreSuite() conform.Suite[store.Store]

// RunStoreSuite applies the checks of Store, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunStoreSuite[T store.Store](t *testing.T, newSubject func() T, opts ...StoreOption)

// RunStoreSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStoreSuiteFromHandle[T store.Store](t *testing.T, newSubject func(h *StoreHandle) T, opts ...StoreOption)
```

`StoreSuite` takes the identity of each check from `StoreChecks`, so a run of the suite calls each function of an identity. The proofs of `Store` declare `TestStoreSuite`. The TypeScript file declares the same names in camel case. The function `storeChecks.read.zeroOnError()` returns the identity of a check.

### The lock

The lock is a source file of the target language with the keys of the checks of a package, one per line, sorted. A regeneration that loses a check removes one line from the lock. The generation gate fails while the committed lock differs from a fresh run, so a reviewer sees a lost check as a deleted line in a small file. A removal is legitimate when that deleted line is in the same change as the edit of the declaration that caused it.

### A worked example

A team declares a store of sessions, adds a check of its own and applies the suite to two implementations. One implementation reads the clock of the handle, and the other excuses the check that needs that clock.

| File | Written by | Contents |
|---|---|---|
| `session/session.go` | The team | The declaration with its directives |
| `session/sessiontest/session_suite.go` | dokimi | The names, the checks with their planted defects, and the application |
| `session/sessiontest/session_suite_test.go` | dokimi | The falsification test |
| `session/sessiontest/session_stub.go` | dokimi | The double, from which the planted defects are built |
| `session/sessiontest/checks.go` | The team | The check of the team with its planted defect |
| `session/memory/memory.go` | The team | The store in memory |
| `session/memory/memory_test.go` | The team | The application on the clock of the handle, with a check of the store in memory |
| `session/postgres/store.go` | The team | The store on PostgreSQL, with the directives `subject` and `excuse` |
| `session/postgres/store_suite_test.go` | dokimi | The application to the store on PostgreSQL |

#### The declaration

```go
// Package session keeps the sessions of signed-in users.
package session

// ErrNotFound is the error of a read of a session that the store does not
// have.
var ErrNotFound = errors.New("session: not found")

// Session is the session of one signed-in user.
type Session struct {
	ID    string
	User  string
	Roles []string
}

// Store keeps the sessions of signed-in users.
//
//dokimi:suite
//dokimi:out path=sessiontest/
type Store interface {
	// Put stores s under s.ID and starts its lifetime.
	Put(ctx context.Context, s Session) error

	// Get returns the session with the ID id until its lifetime ends.
	//dokimi:mixin not-found sentinel=ErrNotFound
	//dokimi:mixin read-after-write write=Put
	//dokimi:mixin ttl duration=30m put=Put read=Get not-found=ErrNotFound
	Get(ctx context.Context, id string) (Session, error)

	// Delete removes the session with the ID id.
	//dokimi:mixin delete-removes read=Get sentinel=ErrNotFound
	Delete(ctx context.Context, id string) error
}
```

The derivation binds eleven claims, and the header of `session_suite.go` lists them:

| Identity | Licence | Machinery |
|---|---|---|
| `Store.Put/smoke` and `Store.Put/cancel` | Every callable, and the parameter in the role context | `nothing-beyond-a-caller` |
| `Store.Get/smoke` and `Store.Get/cancel` | Every callable, and the parameter in the role context | `nothing-beyond-a-caller` |
| `Store.Get/zero-on-error` | The shape reader | `nothing-beyond-a-caller` |
| `Store.Get/miss` | The mixin not-found | `nothing-beyond-a-caller` |
| `Store.Get/read-after-write` | The mixin read-after-write | `nothing-beyond-a-caller` |
| `Store.Get/expires` | The mixin ttl | `a-clock-the-run-advances` |
| `Store.Delete/smoke` and `Store.Delete/cancel` | Every callable, and the parameter in the role context | `nothing-beyond-a-caller` |
| `Store.Delete/delete-removes` | The mixin delete-removes | `nothing-beyond-a-caller` |

The header also lists the shape writer of `Put` and the shape deleter of `Delete` as unclaimed classifications with the stance `silent`.

#### The check of the team

No claim of the spec states that `Get` returns a copy of the roles, so the team adds a check of it. The check is in the package of the suite, so every implementation runs it:

```go
package sessiontest

import (
	"context"
	"math"

	"go.dokimi.dev/assert"

	"example.com/shop/session"
	"example.com/shop/session/memory"
)

// ReturnsACopyOfTheRoles checks that a change of the roles that Get returns
// leaves the session in the store unchanged.
//
//dokimi:check Store.Get defect=sharesTheRoles
func ReturnsACopyOfTheRoles(h *StoreHandle) {
	s := h.Subject()
	put := session.Session{ID: "s-1", User: "u-1", Roles: []string{"admin"}}
	assert.NoError(h, s.Put(h.Context(), put), "Put stores the session")
	first, err := s.Get(h.Context(), "s-1")
	assert.NoError(h, err, "Get reads the session")
	assert.Length(h, first.Roles, 1, "Get returns the role of the session")
	first.Roles[0] = "guest"
	second, err := s.Get(h.Context(), "s-1")
	assert.NoError(h, err, "Get reads the session again")
	assert.Equal(h, second.Roles, []string{"admin"}, h.Statement())
}

// sharesTheRoles returns a Store whose Get returns one slice of roles for an
// ID until the next Put of the ID.
func sharesTheRoles(h *StoreHandle) session.Store {
	fake := memory.New(h.Clock(), math.MaxInt)
	shared := map[string][]string{}
	return NewStoreStub(h, StoreStubDelegateTo(fake),
		WithStorePut(func(ctx context.Context, s session.Session) error {
			delete(shared, s.ID)
			return fake.Put(ctx, s)
		}),
		WithStoreGet(func(ctx context.Context, id string) (session.Session, error) {
			s, err := fake.Get(ctx, id)
			if err != nil {
				return s, err
			}
			if roles, ok := shared[id]; ok {
				s.Roles = roles
			}
			shared[id] = s.Roles
			return s, nil
		}))
}
```

The planted defect passes each guard of the check and fails its last assertion, whose message is the statement. It forwards every other call to the store in memory, so it passes every other check of the suite. The suite generator adds the check to the suite:

```go
// ReturnsACopyOfTheRoles returns the identity of the check that a change of
// the roles that Get returns leaves the session in the store unchanged.
func (c storeGet) ReturnsACopyOfTheRoles() conform.ID[session.Store] {
	return c.Check("returns-a-copy-of-the-roles")
}

// StoreSuite returns the checks of Store.
func StoreSuite() conform.Suite[session.Store] {
	return conform.Suite[session.Store]{Declaration: "Store", Checks: []conform.Check[session.Store]{
		// The eleven derived checks.
		{
			ID:        StoreChecks.Get.ReturnsACopyOfTheRoles(),
			Statement: "a change of the roles that Get returns leaves the session in the store unchanged",
			Apply:     ReturnsACopyOfTheRoles,
			Defect:    sharesTheRoles,
		},
	}}
}
```

#### The store in memory

`memory.New(clock, limit)` returns a store that reads the time from `clock` and keeps at most `limit` sessions. The test of the team builds each subject on the clock of the handle, because the check `Store.Get/expires` advances that clock. It adds a check of the limit, which only the store in memory has:

```go
package memory_test

// TestStore applies the suite of Store to the store in memory, with the check
// of its limit.
func TestStore(t *testing.T) {
	t.Parallel()

	sessiontest.RunStoreSuiteFromHandle(t, func(h *sessiontest.StoreHandle) *memory.Store {
		return memory.New(h.Clock(), 2)
	}, conform.With(evictsTheFirstSession()))
}

// evictsTheFirstSession returns the check that Put evicts the session that it
// stored first once the store has two sessions. Its planted defect is the
// store in memory without a limit, which keeps every session.
func evictsTheFirstSession() sessiontest.StoreCheck {
	return sessiontest.StoreCheck{
		ID:        sessiontest.StoreChecks.Put.Check("evicts-the-first-session"),
		Statement: "Put evicts the session that it stored first once the store has two sessions",
		Apply: func(h *sessiontest.StoreHandle) {
			s := h.Subject()
			for _, id := range []string{"s-1", "s-2", "s-3"} {
				assert.NoError(h, s.Put(h.Context(), session.Session{ID: id, User: "u-1"}), "Put stores the session")
			}
			_, err := s.Get(h.Context(), "s-1")
			assert.ErrorIs(h, err, session.ErrNotFound, h.Statement())
		},
		Defect: func(h *sessiontest.StoreHandle) session.Store {
			return memory.New(h.Clock(), math.MaxInt)
		},
	}
}
```

The application falsifies `Store.Put/evicts-the-first-session` on its planted defect before it runs the check on the store. A test that built the store on the clock of the system through `RunStoreSuite` would fail `Store.Get/expires` with the message of [Failure messages](#failure-messages).

#### The store on PostgreSQL

PostgreSQL computes the end of a lifetime with `now()`, so no test can advance the clock of the store. The team declares the subject and the excusal on the type:

```go
// Store is the Store on PostgreSQL.
//
//dokimi:subject session.Store new=NewTestStore
//dokimi:excuse Store.Get/expires reason="PostgreSQL computes the expiry with now(), which a test cannot advance"
type Store struct{ … }

// NewTestStore returns a Store on a schema of its own in the database of the
// environment variable SESSION_TEST_DATABASE.
func NewTestStore() *Store { … }
```

The suite generator writes the test of the subject:

```go
// Code generated by dokimi. DO NOT EDIT.
// dokimi:plugin suite
// dokimi:source session/postgres/store.go

package postgres_test

import (
	"testing"

	"example.com/shop/session/postgres"
	"example.com/shop/session/sessiontest"
	"go.dokimi.dev/conform"
)

// TestStoreConformsToStore applies the checks of the suite of Store to Store.
func TestStoreConformsToStore(t *testing.T) {
	t.Parallel()

	sessiontest.RunStoreSuite(
		t,
		postgres.NewTestStore,
		conform.Excuse(
			sessiontest.StoreChecks.Get.Expires(),
			"PostgreSQL computes the expiry with now(), which a test cannot advance",
		),
	)
}
```

#### The output

Each check runs as a subtest whose name is the identity of the check without the name of the declaration. The level of the suite passes when every check fails on its planted defect:

```text
$ go test -v -run 'TestStoreSuite/StoreSuite' ./session/sessiontest
…
--- PASS: TestStoreSuite (0.01s)
    --- PASS: TestStoreSuite/StoreSuite (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Put/smoke (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Put/cancel (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/smoke (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/cancel (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/zero-on-error (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/miss (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/read-after-write (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/expires (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/returns-a-copy-of-the-roles (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Delete/smoke (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Delete/cancel (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Delete/delete-removes (0.00s)
```

The test of the store on PostgreSQL reports the excused check as skipped, with the line of the application and the reason:

```text
$ go test -v -run 'TestStoreConformsToStore/Get/expires' ./session/postgres
=== RUN   TestStoreConformsToStore
=== PAUSE TestStoreConformsToStore
=== CONT  TestStoreConformsToStore
=== RUN   TestStoreConformsToStore/Get/expires
=== PAUSE TestStoreConformsToStore/Get/expires
=== CONT  TestStoreConformsToStore/Get/expires
    store_suite_test.go:19: Store.Get/expires is excused: PostgreSQL computes the expiry with now(), which a test cannot advance
--- PASS: TestStoreConformsToStore (0.00s)
    --- SKIP: TestStoreConformsToStore/Get/expires (0.00s)
PASS
ok  	example.com/shop/session/postgres	0.214s
```

A change of the store on PostgreSQL that caches each session and returns the roles of the cache fails the check of the team at its statement:

```text
$ go test ./session/postgres
--- FAIL: TestStoreConformsToStore (0.38s)
    --- FAIL: TestStoreConformsToStore/Get/returns-a-copy-of-the-roles (0.01s)
        checks.go:26: a change of the roles that Get returns leaves the session in the store unchanged: (-want +got)
            	[0]: -"admin" +"guest"
FAIL
FAIL	example.com/shop/session/postgres	0.402s
FAIL
```

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
| A check of a consumer fails on its planted defect at the assertion of its statement before it runs on a subject | The generated falsification test, and the application for a check that `conform.With` adds |
| A planted defect passes every other check of its suite, except the checks that share it and the deeper checks of its edges | The generated falsification test, and the application for a check that `conform.With` adds |
| The subject of a check is the value that the factory of each form of the application builds | The generated test of the two forms of the application |
| Every excusal names a check of its run | The generator for the directive `excuse`, and the application |
| The checks of a run have distinct identities | The generator for the directive `check`, and the application |
| A function under the directive `check` binds to one callable of one declaration with the directive `suite` | The annotate step of the suite generator |
| A check that advances the clock runs on a subject that read the clock of its handle | The handle |
| A check that calls a door has the option of the door | The handle |

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
| The planted defect of a check of a consumer passes the check, or fails an assertion of another contract than its statement | The falsification of the check fails and names the contract | Fix the check or its defect, or pass the statement to the decisive assertion |
| A planted defect fails another check of its suite | The falsification of its check fails, and lists the defect with each other check that fails on it | Change the defect so that it passes the other checks, or give each check that fails on it the same defect |
| The directive `check` names a declaration without the directive `suite`, an unknown callable, a taken segment or a missing defect | Generation reports it at the position of the directive | Correct the directive |
| The directive `check` names a declaration that two declarations with the directive `suite` have | Generation reports it at the position of the directive, with the package of each declaration | Qualify the name with the name of its package |
| A check of a function has the segment `builds-the-subject` or `receives-the-function` | Generation reports it at the position of the directive, with the identity of the check of the test of the two forms | Rename the function of the check |
| The directive `excuse` names no check of the suite | Generation reports it at the position of the directive, with the identities of the checks of the callable | Correct or remove the excusal |
| An excusal of `conform.Excuse` names no check of the run | The run fails and names the identity | Remove the excusal |
| Two checks of a run have one identity | The run fails and names both checks | Rename the check of the consumer |
| A check advances the clock of a subject that never read it | The check ends with a fault that wraps `conform.ErrClock`, and the message contains the call of `RunStoreSuiteFromHandle` that builds the subject on the clock | Build the subject on the clock of the handle |
| A check calls a door without its option | The check ends with a fault that wraps `conform.ErrOption` and names the option | Pass the option |

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

### A separate falsification of the checks of a consumer

The generated file declares `FalsifyStore`, which a consumer calls beside `RunStoreSuite` with the same checks.

**Why not:** a consumer who does not call it runs checks that no planted defect has failed. An application that falsifies each check of the consumer first cannot run an unproven check.

### Reason codes of the consumer

A check of a consumer declares a reason code of its own, and its body reports the code before each assertion.

**Why not:** the assert library records the assertion of each failure, and no assertion records a code of the consumer. Each body would need a call of the runtime beside its assertions, and a missing call would turn each falsification into a wrong reason.

### The assertion of the first failure as the reason

A check of a consumer declares an assertion of assert-spec, such as `equal`, and falsification compares it with the assertion of the first failure record.

**Why not:** a guard of the same assertion fails first on many defects, so its failure would count as proof. The consumer would also restate in a field what the decisive assertion already states.

### Checks of a consumer registered at run time alone

A consumer passes each of its checks to `conform.With`, and no directive declares one.

**Why not:** the lock and the header then miss the checks of the consumer, so a lost check does not show in review. A misspelled callable or defect fails only when the test runs.

### A factory of every subject that takes the handle

`RunStoreSuite` takes a factory with the handle as its parameter for every subject. The conformance tests of the Go CDK take one factory in this way: `drivertest.RunConformanceTests(t, newHarness, asTests)` of `gocloud.dev/blob` takes a `HarnessMaker`, which receives the context and the test.

**Why not:** most constructors do not need the handle, and each test would wrap its constructor in a function literal. `RunStoreSuiteFromHandle` serves the subjects that need the handle.

### The names `RunStore` and `RunStoreOn`

The application of the suite of `Store` is `RunStore`, and its form with the handle is `RunStoreOn`.

**Why not:** a reader takes `RunStore` for an instruction to run the store, because the name does not mention the suite. From the word `On`, a reader cannot tell what the factory receives. Each conformance helper of Go has the name of what it tests: `fstest.TestFS`, `iotest.TestReader`, `nettest.TestConn` and `drivertest.RunConformanceTests`. In `RunStoreSuiteFromHandle`, the words `StoreSuite` refer to the suite that `StoreSuite` returns, and `FromHandle` refers to the argument of the factory. The suffix `WithHandle` would collide with `conform.With`, which adds checks.

### An application to a control

The generated file declares `ControlStore` and `ControlStoreFromHandle`, which apply the checks to a control. The suite record of such a run has the application `control`, so a reader of the record blames a failure on the check and not on the implementation.

**Why not:** the runtime runs a control exactly as it runs a subject, so the two functions differ from the application only in the record. The controls of dokimi are in its corpus. The corpus fails on every failed check, so the record adds nothing that the corpus does not enforce. Every suite would declare two functions that no consumer calls. The verdict reverses when the service needs to tell the failures of a control of a consumer from the failures of a subject.

## Drawbacks

- **Templates for every defect form.** testkit uses 36 defect forms. With eight targets, that is 288 defect templates, and a new defect form needs a template in every tree.
- **Every runtime implements 59 probes.** The probes are host code, so each runtime writes and tests all 59.
- **A double for every suite.** A declaration with `suite` also gets the typed double. For the corpus fixture of testkit with one method, the suite and its proofs have 392 lines and the double with its test has 465.
- **Structural equality.** Probes compare results by structure. A type that its language compares by identity, such as a class in TypeScript, compares by structure in a check.
- **Claims that need parameters.** A partner claim refuses until its directive sets the partner, so a mixin without its parameters adds refusals to the header.
- **A check of a consumer runs twice.** It runs on its planted defect and then on the subject.
- **Falsification grows with the square of the checks.** Each check of a suite runs on the planted defect of every other check, so a suite of 40 checks makes 1,600 applications.
- **A guard can reuse the statement.** A guard that passes `h.Statement()` as its message turns its own failure into proof. eidos does not model the bodies of functions, so the generator cannot detect it.
- **Bound checks are outside test files.** The functions of the directive `check` are in the package of the generated suite, beside the doubles. A production package that imports that package compiles test code.
- **Long names.** The longest declaration of the corpus gives `RunLeasedIdempotentWriterSuiteFromHandle`, a name of 40 characters.

## Open questions

None.

## References

| What | Where |
|---|---|
| The claims table of testkit | <https://github.com/thesm-os/testkit/blob/main/generator/suite/claims.go> |
| The shape catalog of eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/12-shape-catalog.md> |
| The callable that the rules of eidos derive | <https://github.com/dokimasia/eidos/blob/main/eidos-core/rules/callable.go> |
| The equality of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/assertions.yaml> |
| The reads of the directives of a plan in eidos, `Reader.ByDirective` among them | <https://github.com/dokimasia/eidos/blob/main/eidos-core/store/reader.go> |
| The subjects of a key and the instances of a directive on a subject in eidos, `Subjects` and `SubjectMatch.Directives` | <https://github.com/dokimasia/eidos/blob/main/eidos-core/match.go> |
| The conformance helpers of the Go standard library | <https://pkg.go.dev/testing/fstest#TestFS> and <https://pkg.go.dev/testing/iotest#TestReader> |
| `nettest.TestConn` | <https://pkg.go.dev/golang.org/x/net/nettest#TestConn> |
| The conformance tests of the blob drivers of the Go CDK | <https://pkg.go.dev/gocloud.dev/blob/drivertest> |
