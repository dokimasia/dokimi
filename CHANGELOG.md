# go.dokimi.dev/dokimi

## 0.11.0

### Minor Changes

- 8c423be: Add `//dokimi:action final` for a change of the consumer that the run takes once after the faults end.
- 154c04f: Add the param `unknown` to `//dokimi:system` for the errors of a call of a component that leave the outcome of the call unknown.

### Patch Changes

- 363d72a: Import the package of each name of the length of an array in a generated Go suite.
- a0bb7c2: Leave out each entailment of a check that a gate refuses from a generated Go suite.
- a0bb7c2: Leave out the generated binding of a sample binding that a gate refuses from a generated Go suite.
- 363d72a: Require `go.dokimi.dev/eidos/core` v0.16.0 and the releases of the other eidos modules that depend on it.

## 0.10.0

### Minor Changes

- d74c1ef: Let a `rules.Licence` of a declaration name a classification that a callable of the declaration must have.
- c2eb6f1: Check that a child error value of a family matches its parent and opens with its message, in place of the two checks that refused it.
- 48a0255: Add the params `mean-steps` and `max-steps` to `//dokimi:sim` and `//dokimi:system` for the steps of a case.
- 20889c0: Accept the faults `reset`, `slow-link` and `disk-replacement` in the params `tolerate` and `require` of a simulation.
- cf24c4e: Add the param `dependency` to `//dokimi:perform func=`, so the faults of the interface that the function calls strike its effect.
- cf24c4e: **Breaking:** Require a function under `//dokimi:perform func=` to return `sim.Outcome[E, F]`.
- b89679f: Add the param `end` to `//dokimi:perform func=` for the event that ends a wait at an outage after the wait has delivered an event.
- b89679f: Add the param `ends` to `//dokimi:perform` for an effect that ends the waits of the effects that it lists.
- 48a0255: Accept a function of `//dokimi:perform func=` that takes a pointer to its effect.
- cf24c4e: Accept a system and the defect of an invariant of the form `func(c *prop.Case) func(w *sim.World[E, F])`, which take their inputs from the case.
- 20889c0: Generate the simulation of a system of components from a function under `//dokimi:system` whose world is `*sim.World[sim.None, sim.None]`, with the modes `tasks` and `blocking` in its param `tolerate`.
- e9a9822: Add the param `falsify` to the directive `suite` for the function of the options that the spec of the suite passes to `conform.Falsify` and `conform.FalsifyAllocs`.
- e9a9822: Generate the world of a declaration with the mixin `durable`, and the checks of `linearizable-under-faults` and `completes-after-recovery` over it.
- d74c1ef: Add the rules of the claims `linearizable-under-faults` and `completes-after-recovery`, with their laws and defect forms.
- e9a9822: Add the key that a call of the mixin `creates` returns to the pool of its type in the world of a declaration.

### Patch Changes

- e9a9822: Fix the doc comments of the planted defects of durable claims to state the clause and the door in separate sentences.
- 9f75b0a: Fix `dokimi run` to update `dokimi.go.lock` in a clone without `.dokimi/`.
- 7686b19: Fix `dokimi conform` to show the last output of a test binary that ends in the middle of a test.
- 7ac433a: Require `go.dokimi.dev/assert` v0.19.0.
- c26cbdf: Require `go.dokimi.dev/assert` v0.20.0 and `go.dokimi.dev/spec` v0.5.1.
- 7ac433a: Require `go.dokimi.dev/eidos/core` v0.14.0 and the releases of the other eidos modules that depend on it.
- cf24c4e: Require `go.dokimi.dev/eidos/core` v0.15.0 and `go.dokimi.dev/eidos/sdk` v0.11.0, and the releases of the other eidos modules that depend on them.
- c797659: Require `go.dokimi.dev/eidos/lang/go` v0.8.2.
- 1ff4198: Fix the strict run of a tree in which a module imports the S3 client of the AWS SDK.
- c26cbdf: Require `go.stealthscale.io/cli` v0.8.1 and `go.yaml.in/yaml/v3` v3.0.6.
- cf24c4e: Fix the simulation to bind a performer whose outcome has a name of its own.
- 20889c0: Report `DOKIMISIM-0003` at the directive system of a machine interface that states a param of its simulation.
- cb10807: Fix the sim generator to leave the types of test files out of the events and the effects of a machine.
- e9a9822: Fix the suite to start a struct once from its own sources where a transparent alias spells it.
- 530c71d: Fix the suite of a callable whose input contains `proto.Message`.
- 1db2c45: Fix the check of `linearizable-under-faults` to search its history within the budget of the other history checks.

## 0.9.1

### Patch Changes

- c2a7f3d: Require `go.dokimi.dev/eidos/core` v0.13.1 and the releases of the other eidos modules that depend on it.
- c2a7f3d: Require `go.dokimi.dev/spec` v0.4.1.

## 0.9.0

### Minor Changes

- d9b82f1: Add the claim `returns-on-error`, which checks that a failed call returns the field of its request.
- c4b00ef: Add the claims `serializable` and `snapshot-isolation`, which check the histories of concurrent transactions.
- d9b82f1: Add the directive `refuse`, which withdraws a derived claim of a callable with a reason.
- 3b93926: Refuse the claims `deterministic` and `predicate` of a callable with the mixin `time-aware` or `side-effect`.
- b08c44f: Add `--mutation` to `dokimi report`, which reports the dead and the subsumed checks of the records of mutation runs.
- d9b82f1: Take the inputs of each derived call of a callable from the builder of its mixin `sample`.
- e41c924: **Breaking:** Remove the header comment of a generated Go suite.
- e41c924: **Breaking:** Name a generated binding as `StoreChecks.Get.ZeroOnError().Generated()` in place of its identity method `ZeroOnErrorGenerated`.
- e41c924: Excuse the checks of a suite in its generated Go tests with one `conform.ExcuseSuite`.
- e41c924: Leave the probe, the reason code and the law out of a generated derived check, which conform takes from its claim.
- e41c924: Build the suite of a declaration once in each generated Go test file that applies its planted defect of a ceiling.
- 903079f: **Breaking:** Remove the constructors `New<Double>Delayed` and `New<Double>Failing` of a generated Go double.
- 903079f: Check the arguments, the hook `OnRecord`, the reset and the fixed results of each method of a generated Go double through `stub.Behaviour` and `stub.ResetBehaviour`.

### Patch Changes

- 9e701f9: Leave out `StrictKanon` for an enumeration that kanon encodes as a `kanon.Validator`.

## 0.8.0

### Minor Changes

- 26a0797: Add the claim `causal-history` to the spec, with its probe and defect form.
- 554bc2a: Add the claim `computes-once` to the spec, with its law, reason, probe and defect form.
- 554bc2a: Add the claim `linearizable` to the spec, with its law, reason, probe and defect form.
- 554bc2a: Add the claim `one-holder-at-a-time` to the spec, with its probe and defect form.
- 26a0797: Add the claim `pins-the-snapshot-under-writes` to the spec, with its probe and defect form.
- 554bc2a: Add the claim `reads-agree-under-concurrency` to the spec, with its law, reason, probe and defect form.
- 554bc2a: Add the five claims of the rules of a runtime of the generator `sim` to the spec, with their laws, reasons, probes and defect forms.
- 554bc2a: Add the seven claims of deliveries, iterations, totals and failed calls to the spec, with their laws, reasons, probes, gates and defect forms.
- e766b0b: Add the claims `synced-writes-survive`, `no-torn-record`, `corruption-is-reported` and `full-medium-is-reported` of the mixin `durable` to the spec.
- e766b0b: Add the gate `reference-role` to the spec, which the generator `fake` judges.
- 24d69ed: Add the generator `fake`, which writes an in-memory fake of an interface under the directive `fake` and a test that applies the suite of the interface to it.
- 5365049: Add the generator `sim`, which writes the deterministic simulation of a system of sans-IO machines and its tests from the directives on their declarations.
- 378b01a: Add the partner mode `paired`, which binds the one read of a declaration that pairs with the licensed write.
- e766b0b: Add the field `points` to the check record of the suite schema, the points of a medium that a check enumerated.
- e766b0b: Add the reference forms of the contracts `persister`, `cache` and `appender` to the spec, with `spec.References` and `spec.ReferenceOf`.
- 554bc2a: Add the stance `gap` of the mixin `returns-on-error` to the spec.
- f9e5c64: Export `suite.Suited` and `suite.NamesOf` for a generator that writes a test of the suite of a declaration.
- f9e5c64: Generate the checks of the mixin `durable` over `probe.Durable`, with planted defects that open over a medium.
- dd9cc74: Judge each fake that the key `fakes` of an expectation file lists, as `dokimi conform` judges a subject without an excusal.
- 378b01a: Judge the classifications, the keys and the results of the callables of a declaration under the directive `fake` by the reference forms of the spec.
- 88a1009: License a claim of the scope declaration through a mixin of one of the callables of the declaration.
- 2ff1bac: Write the marker `StrictKanon` of `kanon.Strict` beside the generated `MarshalText` and `UnmarshalText` of a Go enumeration.
- 114b6f3: Accept the faults `torn-write`, `corruption` and `full-disk` in the params `tolerate` and `require` of `//dokimi:sim`.
- c7e65c0: Write the check `causal-history` and its planted defect `racy-causal-read` in the Go suite of a read with the mixin `causal`.
- eb9d0db: Write the check `computes-once` and its planted defect `computes-per-call` in the Go suite of a callable with the contract `single-flight`.
- eb9d0db: Write the check `linearizable` and its planted defect `racy-read-write` in the Go suite of a concurrent declaration with a reference.
- eb9d0db: Write the check `one-holder-at-a-time` and its planted defect `check-then-grant` in the Go suite of a callable with the contract `lease`.
- c7e65c0: Write the check `pins-the-snapshot-under-writes` and its planted defect `racy-pin` in the Go suite of a read with the mixin `point-in-time`.
- eb9d0db: Write the check `reads-agree-under-concurrency` and its planted defect `racy-read` in the Go suite of a callable with the mixin `concurrent-readers`.
- eb9d0db: Write the checks of deliveries, iterations, totals and failed calls in the Go suites.

### Patch Changes

- 88a1009: Fix the derivation to bind the claims of an input of a defined function type.
- eb9d0db: Fix the generated binding to keep the sample of a generated function whose results are generated functions.
- 8ab433b: Fix `dokimi conform` to write each test of the corpus that fails, with its output.
- eb9d0db: Name the partners of a check in the reason and the remedy of each gate that the suite generator judges.
- eb9d0db: Fix the order of the declarations under the directive `suite` to their identity order, whatever the spelling of the directive.

## 0.7.0

### Minor Changes

- deab4c1: Add the seven claims of the machinery `call-that-fails-on-demand` to the spec, with their laws, reasons, probes and defect forms.
- 8e6e7cf: Add the 19 claims of the generated sequences of a declaration, with their laws, probes, reasons, gates and defect forms.
- 8e6e7cf: Add the claims of the eight classifications without one, and remove the stance `gap` from the unclaimed table.
- 4302828: Add the command `dokimi conform`, which judges the corpus of a language by the expected outcomes of its fixtures.
- 027e24b: Add the command `dokimi report`, which prints the suite records of a run of `go test -json`.
- 6812b3b: Add the package `spec` with the laws, claims, gates, probes, reason codes, defect forms, machineries and stances of the conformance spec.
- edd824f: Add the door `faults` of the directive `check`, which lists the machinery `call-that-fails-on-demand` in the check.
- deab4c1: **Breaking:** Add the required field `gate` to each refusal of the suite record and to `spec.RefusalRecord`.
- 2a3df20: Add the gate `partner-arguments` to the conformance spec.
- 3b020ed: Add the gate `returns-unaided`, which withdraws the claims of a fresh subject from a callable with the mixin `waits`.
- deab4c1: Add the gates `input-samples`, `classification-params`, `interface-handle` and `importable-samples`, so every refusal has a gate of the spec.
- 4538348: Add the inputs `list`, `valid`, `outside` and `name` of the probes of enumerations to the conformance spec.
- e2c0508: Add `spec.Key` and the key vectors of `spec/keys.yaml` for the key form of a check identity.
- a18d071: Add `dokimi.go.lock` with a line for each check, ceiling, refusal and unclaimed classification of the suites of the plan go.
- 8d892c5: Add the option `isolate` to the generator `suite`, which passes `conform.Isolate` in the conformance test of each subject.
- edd824f: Add the param `machineries` of the directive `suite`, and `DOKIMISUITE-0013` for a value that names no machinery of a consumer.
- 1357933: Add the param `normalize` to the directive `suite` to map the values of each call of a generated sequence.
- ef3d7eb: Add `spec.SuiteRecord`, `spec.ValidateSuiteRecord` and `schemas/suite.schema.json` to the conformance spec.
- 8e6e7cf: Add the licence of every declaration, the source `declaration`, the default of a param and four input kinds of a probe to the spec.
- 6caeec8: Bind a claim of a whole declaration once, on its first callable that passes the gates of the claim.
- 6caeec8: Bind a claim of generated sequences once, without a generated binding and without the gate `input-samples`.
- 6caeec8: Bind a partner that a probe of generated sequences calls by its name without a binding of its arguments.
- 1357933: Bind the probes of the classifications without a claim in the suite generator.
- 3b020ed: Add the derived checks and their planted defects to the Go suites of the generator `suite`.
- 1357933: Drain the `iter.Seq` or `iter.Seq2` that a callable returns into a list for the probes of a drain.
- 027e24b: Fail `dokimi report` for a suite record whose field `ended` lists the problems that stopped its application.
- 6caeec8: Pair a batch read with the writer of the elements of its list.
- 027e24b: Print each problem of the field `ended` on a line of the report of its suite.
- cb1a9a5: Add `DOKIMIBUILDER-0007` at a declaration whose builder refers to a declaration that the file of the builder cannot import.
- 3485cf2: Add `DOKIMISUITE-0012`, which `dokimi run --check` reports for each check of the committed lock that the run no longer generates.
- ef3d7eb: State `checks[].counterexample` of the suite record as the list of entries of the counterexample of a run of assert-spec.
- b21bf98: Add the package, the generator, the refusals and the unclaimed classifications to each Go suite of the generator `suite`.
- 1357933: Write the actions, the call functions and the reference of a declaration whose claims run generated sequences.
- 1357933: Write the checks of generated sequences with their planted defects, and refuse them by the gates `countable-use`, `takes-a-replica` and `graph-of-states`.
- edd824f: **Breaking:** Write the gate of each refusal into the generated suite of Go, which needs `conform.Refusal.Gate`.
- 8d892c5: Add the header of a generated suite, which lists its derived checks, its refused claims and its unclaimed classifications.
- 453d5eb: Add `spec/naming.json` with the Go and TypeScript spelling of each identifier of the conformance spec.

### Patch Changes

- 8d892c5: Add ADR-0018, which keeps the tests of enumerations and errors as fixed templates.
- 8d892c5: Fix the planted defect of the claim `cancel` with the new defect form `ignores-cancel`, which succeeds for a cancelled context and calls the subject for any other.
- 3b020ed: Fix the planted defect of the claim `if-match` with the new defect form `ignores-the-predicate`, which applies its first write and accepts every later write.
- 8d892c5: Replace the defect form `forgets-writes` with `reads-another-value`, which breaks the read of a read-after-write law instead of its write.
- 3b020ed: Fix the planted defect of the claim `tx` with the new defect form `settles-twice`, whose rollback settles a committed transaction.
- 8d892c5: Fix the planted defect of the claim `writes-follow-reads` with the new defect form `stamps-the-past`.
- 8d892c5: Fix the partner mode `derived` to bind the callable that takes the value whose field the param names.
- 3b020ed: Fix `//dokimi:excuse` to accept the identity of a derived check, such as `Store.Get/zero-on-error`.
- 8d892c5: Fix the defect form `fails-always` to fail each call where the call of the subject succeeds.
- edd824f: Fix the note of a needed machinery to spell the param `machineries` as a list, such as `machineries=[faults]`.
- 4538348: Fix the claim `partition` to read its partner in the mode `reader`.
- 3b020ed: Fix the claim `total` to need the gate `failure`.
- 8d892c5: Fix each planted defect of a callable with a context and an error to return the error of a cancelled context first.
- edd824f: Fix the defect form `keeps-failed-writes` to return its zero results only until a call of a dependency fails.
- 3efcbb7: Fix a strict run that refuses an interface with two embedded interfaces of one bare name.
- cb1a9a5: Fix the tests of `builder` to leave out each default, sample, member and witness that only the package of its declaration can refer to.
- aca5bc3: Fix the tests of `stub` to leave out each sample and witness that only the package of its declaration can refer to.
- aca5bc3: Fix `stub` to write the double of an unexported interface with unexported constructors and options and without checks.
- 7ffbb81: Fix the error of `dokimi report` for a file that does not exist to name the file once.
- 6caeec8: Fix the pairing of a callable without a value with the named partner that reads what it writes, such as the `observe` of `accumulates`.
- 3b020ed: Fix the partner mode `role` to pair a callable without results as the write.
- 6caeec8: Fix the panic of the suite generator on a read of `writes-follow-reads` whose write returns no value of its type.
- 71b3508: Fix `dokimi conform` to report a group of fixtures that is a file on Windows.
- aca5bc3: Fix `stub` to report `DOKIMISTUB-0003` at an unexported subject whose double the directive `out` moves into another package.
- 3b020ed: Require `go.dokimi.dev/assert` v0.12.0.
- c238979: Warn about a directive whose brand is near the brand dokimi, through the eidos modules at 482b0bb7.
- 1738f06: Require the eidos modules at 8d191b8e, which bound the zero value of a cycle of aliases.
- 8d892c5: Fix the five test modules of the generators and the command to require `go.dokimi.dev/assert` v0.12.0 again.
- 8d892c5: Fix the defect forms `value-beside` and `values-beside` to return a sample beside the error of a failed call of the subject alone.
- 1357933: Return an empty stream in place of a nil stream from the planted defect of a callable that returns a stream.
- 8d892c5: Fix the write of a derived check to return the value that the write returns of the type that it writes.
- 3b020ed: Fix the second write of a derived check to take the changed field from the samples of the field.
- 8d892c5: Fix the derivation to take two function literals of one spelling as two samples of an input.
- 3b020ed: Fix the written value of a pair to skip an input that the read takes by its name.
- 8d892c5: Fix the edges of the conformance spec to the edge from `deterministic` to `smoke`, which its probe proves.
- ef3d7eb: Fix the vacuous conditions of five probes in the conformance spec to name each guard that ends the probe.
- 8d892c5: Fix the arguments of a partner to give a function parameter, and each parameter of a read whose write takes no input, its own sample.
- cb1a9a5: Fix `builder` to write the builder of an unexported struct with unexported constructors and without checks.
- 3b020ed: Fix the check of the claim `if-match` to write a first value and judge a second value that differs in the field of the param `field`.

## 0.6.0

### Minor Changes

- 33b06f2: Add the param `field` of `//dokimi:stub`, the name template of the field that configures each method of a double.
- 33b06f2: Report `DOKIMISTUB-0007` at a declaration whose param `field` is not the name template of a field.

### Patch Changes

- f76be57: Quote the empty name template of an option of enum in its error, as the other errors of a name template do.
- 5994f41: Fix a run over a module whose dependency carries directives of dokimi, through the eidos modules at 83dc1f97.

## 0.5.0

### Minor Changes

- 6f2a905: **Breaking:** Pass each call of a Go double to `stub.Answer` by value, as go.dokimi.dev/stub v0.2.0 requires.

### Patch Changes

- 082136a: Fix the double of a Go struct to implement the exported methods of the struct alone.
- 082136a: Fix the generated checks of a Go double to leave out the unexported methods of the double.
- 082136a: Report `DOKIMISTUB-0003` at a Go interface with an unexported method that its double in another package cannot implement.

## 0.4.0

### Minor Changes

- b586ce9: Add the param `around` of the directive `check` to run each run of the check inside an environment, such as a `testing/synctest` bubble.
- b586ce9: Leave each callable with the mixin `waits` out of the generated benchmarks and out of the allocation ceilings of the test of a suite.
- b586ce9: Write `Waits: true` into the call of each callable with the mixin `waits` in a generated suite.
- b586ce9: Pass the sample that a parameter states on itself in each generated call of its callable.
- a5e2ee8: Release the binaries of the command `dokimi` and its Homebrew cask `dokimasia/tap/dokimi`.
- b586ce9: Report `DOKIMISUITE-0007` at a budget of a callable with the mixin `waits`.
- b586ce9: Report `DOKIMISUITE-0009` at a declaration with a parameter whose sample the file of its suite cannot refer to.
- b586ce9: Report `DOKIMISUITE-0005` at a check whose environment does not take a `*testing.T` and a `func(*testing.T)`, or returns a result.

### Patch Changes

- 71136eb: Require go.dokimi.dev/assert v0.10.0.
- 71136eb: Require eidos-plugin-shape v0.1.0.

## 0.3.0

### Minor Changes

- 50ef79f: Add the params `reference`, `second`, `restart` and `crash` of the directive `check` to give a planted defect its doors.
- 50ef79f: Add the error `DOKIMISUITE-0011` for a subject whose function of the options its tests cannot call.
- 50ef79f: Add the param `options` of the directive `check` to state the options of the property run of a check.
- 50ef79f: Add the param `options` of the directive `subject` to pass the options of an application to each test of a subject.
- 76749dd: Add the Warning `DOKIMISUITE-0010` for a type in the inputs of a suite without a generator.
- 76749dd: Draw nil or `conform.ErrGenerated` for an error, `context.Background()` for a context, nil for `unsafe.Pointer`, and a discarding logger or writer for `*log.Logger`, `*slog.Logger` and `*bufio.Writer`.
- 76749dd: Draw one of the exported implementations of each interface with an unexported method in the inputs of a suite.
- 76749dd: Return a closed channel or the sink of the test from a generated function after its drawn results.
- 76749dd: Generate a closed stream for each channel that a subject receives from in the inputs of a suite, and the sink of the test for each send channel.
- 76749dd: Generate an implementation of each interface of the standard library in the inputs of a suite, and a `*conform.Bytes` for each interface of `io` that it implements.
- 9d599ed: **Breaking:** Keep the plural s of an acronym with the acronym in generated names, such as the field `IDs` of a stub for a parameter `ids` in place of `Ids`.
- 76749dd: Return the address of each generated implementation of an interface, so its value is comparable.
- 76749dd: Start each struct without a readable field from its zero value, also a struct of the standard library and a generic struct.

### Patch Changes

- 44fbf9a: Require go.dokimi.dev/assert v0.9.0.
- 76749dd: Require eidos-core v0.10.0, eidos-sdk v0.7.1, eidos-lang v0.3.3, eidos-lang-go v0.6.0, eidos-lang-typescript v0.4.3, eidos-plugin-shape v0.0.10 and eidos-cli v0.3.2.
- 9d599ed: Require eidos-lang v0.4.0, eidos-lang-go v0.6.1 and eidos-lang-typescript v0.4.4.

## 0.2.2

### Patch Changes

- c379e06: Require eidos-lang-go v0.5.1.

## 0.2.1

### Patch Changes

- 62a5f68: Fix the suite to generate a recursive function type or interface through `prop.Recursive`, with at most 10 base values in one value.
- 62a5f68: Fix the suite to leave out each sample, default, member and declared error that its file cannot refer to.
- e4425ef: Require go.dokimi.dev/assert v0.8.0, eidos-core v0.9.0, eidos-sdk v0.7.0, eidos-lang-go v0.5.0 and the other eidos releases of 2026-10-26.
- 62a5f68: Fix the suite to start a struct without a field that `prop` reads from its zero value.
- 62a5f68: Fix the suite to take the builder of a struct of the package of its declaration alone.

## 0.2.0

### Minor Changes

- 30c1768: Add the param `remote` of the directive `sentinel`, whose error types have no check of the message of the zero value and of the prefix, and the finding `UnknownRemote`.
- e30629f: **Breaking:** Declare `Bench<X>SuiteParallel` only in the suite of a declaration with a concurrent callable.
- e6175e9: Draw a value of `any` in a suite, and draw `*big.Int`, `*big.Rat` and `*time.Location` as the pointers that `prop` reads.
- b066e19: Generate the surface of an unexported Go enumeration with unexported functions and error, as `parseState`, and without generated checks.
- 40ab37e: **Breaking:** Name the allocation test, the benchmarks and the fuzz targets of a subject after the struct and its declaration, as `BenchmarkStoreConformsToStore`.
- e1f4eea: Compare each identity of a check of a suite with its method in the level `<X>Checks` of the generated `Test<X>Suite`.
- b066e19: **Breaking:** Return an error that wraps the error of the parse function from the generated `MarshalText` of a value that the validity method refuses.
- 40ab37e: Let a struct take one directive `subject` for each declaration under the directive `suite`, and report `RepeatedSubject` at a second directive of one declaration.
- 40ab37e: **Breaking:** Report `BadFactory` at a subject whose factory its package does not export, in place of a generated test inside the package.
- e6175e9: Add the finding `Unimportable` of `suite` at a declaration with a callable whose signature the file of its suite cannot import.
- 1533d9e: Add the param `counted` of the directive `check` to run a check that counts allocations with the allocation ceilings of its suite.
- e30629f: Add `Test<X>SuiteBudgets`, `Fuzz<X>Suite<Callable>`, `Benchmark<X>Suite` and `Benchmark<X>SuiteParallel` to the generated test of the suite of an interface or a struct.
- b638e24: **Breaking:** Write the double of a declaration under the directive `suite` only where the directive `stub` asks for it.
- e30629f: Write the conformance test of a subject into its external test package when the factory of the directive `subject` is declared there.

### Patch Changes

- e1f4eea: Build the factory of each form of the application of a suite through `<x>Subjects`, so the checks of `Run<X>Suite` cover it.
- 30c1768: Fix `sentinel` to leave the check of the prefix of a re-exported error value of another package to the package that declares it.
- 5ad4cd7: Fix the suite to draw a generated function or an implementation for each value of a function type or an interface at every depth of an input.
- 00d434d: Fix `stub` to write one double of a declaration under the directive `suite` that a method of another double returns.
- e30629f: Fix `suite` to embed the subject type through an alias, so an interface with a method of its own name gets a suite that compiles.
- 077dc00: Fix the generated doc comments of an enumeration to describe a failed parse as a text that does not match a declared variant.
- 5ad4cd7: Fix the suite to leave each type that its file cannot import to a registration of the test process.
- b638e24: Fix `stub` to report `Unimplementable` at a method whose signature the file of its double cannot import.
- d068e26: Require the eidos modules at 01e3506c, which give each record of a unit the data of its own call.
- c56c3cc: Require the eidos modules at 86dbe80e and assert v0.7.0, which complete the member walk of an embedded type of another module.
- 53d5881: Require the eidos modules at e4ae6a9a, which resolve a factory of the external test package by its name alone.
- 80207ae: Require eidos-core v0.8.0, eidos-sdk v0.6.0 and the other eidos releases of 2026-10-25.
- e6175e9: Fix the directive `subject` to take a factory of the external test package of the directory of its directive `out`.

## 0.1.0

### Minor Changes

- ebc38cf: Add `Compose` that builds the workspace of dokimi with the plans of Go and TypeScript.
- daca2cd: Add `defaults.Value`, which reads the text of a declared default as a value of its type.
- aac310d: Add `RunXBudgets`, `BenchXSuite` and `BenchXSuiteParallel` to each suite of `suite`.
- 40c9a27: Add the annotator `defaults` that stamps the declared default of each field and type.
- 778a836: Add the annotator `fault` that stamps the faults of the double of a method.
- aac310d: Add `Test<Subject>Allocs`, `Benchmark<Subject>` and `Benchmark<Subject>Parallel` to the test of each subject.
- aac310d: Add the call generator of each callable to each suite of `suite`.
- e684f0f: Add the command `dokimi` that runs the plans of Go and TypeScript over a workspace.
- aac310d: Add the directive `budget` with the ceilings `allocs`, `bytes`, `p99` and `mean` and the param `warmup` of a callable.
- f5d0b91: Add the directive `excuse`, which excuses a check of a suite from the conformance test of a subject.
- d844ebb: Add the facts `builder.name` and `builder.seed` of the builder of each declaration.
- 920583c: Add the falsification test `TestXSuite` of each suite.
- 7587cf3: Add `FuzzXSuite` to each suite of `suite`, and a fuzz target of each callable to the test of each subject.
- 44e6594: Add the generator `builder` of a fixture builder and its checks for each struct and class.
- 0f5116c: Add the generator `enum` of the text surface and its checks for each enumeration.
- cf6be8b: Add the generator `sentinel` of the checks of the error values and the error types of a package.
- d80668a: Add the generator `stub` of a recording double and its test for each interface, struct and class.
- 53e9bb8: Add the generator `suite` of the conformance suite of each Go declaration and the conformance test of each subject.
- f5d0b91: Add the levels `RunXSuite` and `RunXSuiteFromHandle` to the generated test of each suite.
- bc95344: Add the name templates `parse-function`, `values-function`, `parse-error` and `valid-method` to the options of `enum`.
- aac310d: Add the param `cases` of the directive `subject` for the number of cases of each property run of a subject.
- aac310d: Add the starting value of each struct and defined type to the calls of a suite.
- cda41be: Answer an unconfigured call of a double with a new value from the builder of each result type.
- 7587cf3: Apply the suite of a function to the function in the level `<Function>` of its test, and declare `Fuzz<Function>`.
- b82e408: Assert the generated TypeScript tests through the assertions of `@dokimi/assert` on the `expect` of vitest.
- 920583c: Bind a function of a consumer to a callable of a suite as a check through the directive `check`.
- 7587cf3: Bind a function of the handle and a `*prop.Case` under the directive `check` of `suite` as a property check.
- bc95344: Check the parse function, the function that lists the variants and the error that a package declares in the generated test of an enumeration.
- b82e408: Check the replacement of each member in the generated test of a builder.
- f5d0b91: Compare the subject of a check with the value of its factory by identity in the generated test of each suite.
- cc4ae7d: Import the runtime of a double from the modules `go.dokimi.dev/stub` and `go.dokimi.dev/fault`.
- d89e811: Pass the context of the test to each context parameter in the generated test of a Go double.
- 437f741: Quote the refused text in the error of the parse function of a Go enumeration.
- cd5e451: Render the text of a Go enumeration of integers in the form of stringer.
- f5d0b91: Report `BadCheck` at a check of a function with the segment `builds-the-subject` or `receives-the-function`.
- bc95344: Report `Conflict` of `enum` at a declaration that has the name of a member of the surface and another form.
- b82e408: Structure the generated test of each double, builder, enumeration and error with a level for each function and a case for each behaviour.
- d89e811: Take any implementation of the methods of a Go struct in the option `DelegateTo` of its double.
- f5d0b91: Write the conformance test of a subject whose factory is not exported inside the package of the subject.
- bc95344: Write the function that lists the variants of an enumeration beside a render that the package declares.
- f5d0b91: Write the tests of the suite generator into one family with the word `suite`, as `<file>_suite_test.go`.

### Patch Changes

- d89e811: Alias the type of the instance of each double in its generated Go test.
- 7d7beff: Build the cause of each error type with `errors.New` in the generated test of the errors.
- 437f741: Check the variants of a Go enumeration in one case list in `IsValid`.
- 7d7beff: Compare the message of each error without a conversion to string in the generated test of the errors.
- b82e408: Copy the lists and maps of a Go builder with `slices.Clone`, `maps.Clone` and `maps.Copy`.
- 66904b0: Fix `builder` to refer to the constructors of a builder through emitted references.
- bc95344: Fix `enum` to leave out the parse function that a package declares beside a generated render.
- ddbb9da: Fix `enum` to refer to the generated members of an enumeration through emitted references.
- bc95344: Fix `enum` to wrap the error of the parse function that a package declares instead of declaring it again.
- c74c8d2: Fix `stub` to refer to the declarations of a double and to the builders of its results through emitted references.
- 4174b1a: Fix `stub` to write the double of each interface and struct under the directive `suite` without the directive `stub`.
- b82e408: Rewrite the doc comments of the generated doubles, builders, enumerations and errors.
- 4174b1a: Stamp `suite.double` on each interface and struct under the directive `suite`, and provide the capability `suite` in the annotator role.
- b82e408: Write zero literals in place of `*new(T)` in the generated test of a Go builder.
