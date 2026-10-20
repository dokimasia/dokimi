---
rfc: 0005
title: The package structure of dokimi
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

# RFC-0005: The package structure of dokimi

## Summary

dokimi is one Go module with one package for each plugin. Each generator is a package under `generator/`, and each annotator is a package under `annotator/`. The root package composes the plugins into a plan for Go and a plan for TypeScript, and `cmd/dokimi` runs the command line of eidos over that composition. A package under `internal/` contains one concept that two or more packages of the module use.

The file `plugin.go` of a plugin contains its code for every language, and a file `<lang>.go` contains its parts for one language, with the templates of that language in `templates/<lang>`. One generator then serves Go and TypeScript, and a new language adds files beside `plugin.go`. The doubles that the stub generator writes call the runtime of their language: dokimi-go or dokimi-typescript. The corpus of every language has the same structure.

## Motivation

A generator that serves one language hides the decisions that a second language exposes. The builder and the sentinel generator of eidos-reference serve Go alone. Their handlers call the Go satellite through `golang.PointerReceiver`, `golang.AssumedName` and the fact key `golang.SatisfiesErrorKey`. A TypeScript target would first have to remove those calls. The authoring guide of eidos states the rule that the calls break: a handler reads source-side projections through its match and emits neutral values through its emitter, and a target fact in a neutral declaration is a leak. A layout with one file for each language puts such calls into `golang.go`, and `plugin.go` then contains no code of one language.

Go and TypeScript are the two languages that eidos reads and writes in full:

| Satellite of eidos | Frontend | Rules | Backend |
|---|---|---|---|
| Go | Yes | Yes | Yes |
| TypeScript | Yes | Yes | Yes |
| Rust and Java | Yes | No | Yes |
| Protobuf | Yes | Yes | No |
| Kotlin and PHP | No | No | No |

Without rules, every parameter of a callable is an input and every return is a value under no error model. A generator of dokimi then has no facts to read in Rust or Java.

A layout by role makes each plugin easy to find, and depguard can then state the import rules in a few lines.

## Detailed design

### The repository

```text
dokimi/
├── cmd/dokimi/          the binary
├── dokimi.go            the root package: the brand, the plans and the composition
├── annotator/<name>/    one annotator, such as annotator/defaults
├── generator/<name>/    one generator, such as generator/builder
├── internal/<name>/     one concept that two or more packages of the module use
├── spec/                the conformance spec
├── testdata/            the tree of the composition and the files that it generates
└── corpus/
    ├── go/              the Go fixtures, in a module of their own
    └── typescript/      the TypeScript fixtures, in a package of their own
```

| Directory | Contents | Imports |
|---|---|---|
| `cmd/dokimi` | The function `main`, which runs the command line of eidos over the composition | The root package and the package `cli` of eidos |
| The root package | `Brand`, `Base`, `Plans` and `Compose` | The annotators and the generators, the packages `workspace` and `authored` of eidos-core, the annotators of the shape catalog of eidos, and the Go and TypeScript satellites |
| `annotator/<name>` | One annotator | The SDK facade of eidos, the package `naming` of eidos-lang and `internal/`. A file `<lang>.go` also imports the root package of the satellite of its language |
| `generator/<name>` | One generator | The SDK facade of eidos, the package `naming` of eidos-lang, the shape catalog of eidos and `internal/`. A file `<lang>.go` also imports the root package of the satellite of its language |
| `internal/<name>` | One concept that two or more packages use | The SDK facade of eidos. `internal/gentest` also imports the frontends and the rules of the satellites, and only test files import it |
| `spec` | The conformance spec and its Go package | The shape catalog of eidos, and the Go stdlib |
| `corpus/<lang>` | The fixtures of one language | The runtime and the assert library of the language |

depguard states the import rules of a plugin in `.ergon/local/.golangci.yml`:

- A plugin and its tests do not import a package of eidos-core. The import lint of eidos denies the kernel to its plugins in the same way.
- `plugin.go` does not import a satellite of eidos.
- A file `<lang>.go` imports the root package of its own satellite as its only package of a satellite.
- Each file of a plugin may import the package `naming` of eidos-lang, which cases a name and makes a name unique. It contains no syntax of a language.

The rules do not cover the root package and `cmd/dokimi`, because both compose the plugins.

No rule of depguard restricts the shape catalog of eidos, `go.dokimi.dev/eidos/plugin/shape`, so a generator may import it. Its specs contain the keys of the facts that the catalog stamps for each shape, mixin and contract. The package contains no syntax of a language.

### The composition

The root package composes the plugins in the way that the root package of eidos-reference composes its ensemble. `Base` registers these parts:

- The brand `dokimi`. A carrier then sets a directive as `//dokimi:stub`, the config file is `.dokimi.yaml` and the state directory is `.dokimi/`.
- The frontends and the rules of Go and TypeScript, and the targets of both satellites.
- The kernel annotators of the authored samples and witnesses, the annotator of the Go satellite, and the annotators of dokimi.
- The annotators of the shape catalog of eidos. The plugin `shape` declares the directives of the catalog and stamps the facts of each shape, mixin and contract. The plugin `shapecheck` validates the instances of the contracts. The stub generator reads the facts of the mixins `deprecated` and `order-after`.

The kernel directive `sample` declares the values that a check draws for a declaration: the sample, and with the param `alternate` a second value that differs from it. dokimi does not declare a directive of its own for these values. The directive applies to a struct, a class and an interface alike.

`Plans` returns two plans:

| Plan | Sources | Generators | Backend |
|---|---|---|---|
| `go` | The Go packages of the workspace, and the Go packages of the stores that they import | builder, stub, enum and sentinel | The Go backend |
| `typescript` | The TypeScript sources of the workspace, and the modules of the stores that they import | builder, stub, enum and sentinel | The TypeScript backend |

Its language alone scopes each plan. eidos then admits the packages of that language that a store provides, such as the standard library of Go, and a generator reads their declarations. A double of an interface that embeds `io.Closer` so has the method `Close`. A workspace narrows the sources of a plan in its config file, or disables a plan that it does not need. A plan that a pattern narrows still reads the packages of the stores of its language. Both plans share one instance of each generator.

### A plugin

Every generator and every annotator has this layout:

```text
generator/builder/
├── doc.go               the documentation of the package
├── plugin.go            the plugin and the model that it derives for every language
├── plugin_test.go       the spec of plugin.go
├── golang.go            the parts of the plugin for Go
├── golang_test.go       the spec of golang.go
├── typescript.go        the parts of the plugin for TypeScript
├── typescript_test.go   the spec of typescript.go
├── codes.go             the codes of its findings, where it reports any
├── codes_test.go        the spec of codes.go
├── templates/
│   ├── golang/          the file templates of the target golang
│   └── typescript/      the file templates of the target typescript
└── testdata/
    ├── golang/          the Go fixtures of the specs
    └── typescript/      the TypeScript fixtures of the specs
```

- `plugin.go` declares the identity, the directives, the options and the outputs of the plugin, and its handlers. A handler derives the model of its output from its match and the facts of eidos, and records the model through the emitter. Where a step differs by language, the handler calls the part of the language of its subject. `plugin.go` contains no word and no syntax of a language.
- `<lang>.go` declares the part of one language. The file embeds `templates/<lang>`, and the part declares the target with that tree through `sdk.TemplatesOf` and with its helpers. The part also takes each step that the facts of eidos leave to the language, such as the types that the language writes as sets and the name of a companion function.
- `New` in `plugin.go` builds the plugin from the parts of the languages in one list. A handler runs with the part of the language of its match. For a language without a part, the plugin does not record a model. A new language adds its `<lang>.go`, its `<lang>_test.go` and its `templates/<lang>`, and one entry of that list.
- Each family of the plugin claims its files through a file template of each target. The template writes the model in the conventions of its language, such as `NewAddress()` in Go and `AddressBuilder.create()` in TypeScript. The conventions cover the form of the API, the words of its names and the idioms of the language. The stubgen plugin of eidos-reference serves both targets in this way. The Build of a composition refuses a plan whose target a generator does not serve.
- `plugin_test.go` runs `plugintest.RunPluginSuite` over the fixtures of each language, which lints every template against the funcmap of each target, and checks the records that the handlers place. `<lang>_test.go` checks the steps of its language over the fixtures in `testdata/<lang>`.
- A plugin without a step of one language has no `<lang>.go` for it. An annotator has no templates.

annotator/defaults stamps the declared default of each field and of each type alias or defined type. The builder reads the default of each member, and the suite generator reads the default of a type. The annotator reads the facts of eidos alone, so it has `plugin.go` and no `<lang>.go`. Its directive rule reads `//dokimi:default <value>` in both languages. Its tag rule reads the tag of the key that the option `tag-key` sets, `default` unless a config sets another, through the rules of the language. The option `no-tags` turns the tag rule off, and the rules of TypeScript state no tags.

annotator/fault stamps the params `retry` and `partition` of the directive `fault` on a method, which the stub generator reads. It also reads the facts of eidos alone.

### Shared packages

`internal/gentest` builds a fixture of the SDK over a tree of Go and TypeScript sources under the brand. It also loads the scripted language of the SDK, which no plugin of dokimi serves, so a spec shows that a plugin leaves the declarations of another language alone. The annotators of a test stamp the fixture before the plugin under test runs. Its helpers read the records and the findings of a phase call, and check the template trees and the codes that a plugin declares. The specs of every plugin use it.

`internal/lang` runs the handlers of a plugin with the part of the language of each subject. The interface of the parts of a plugin embeds `lang.Part`, whose methods return the language, the target and the options of the target. `lang.Register` declares the target of each part on the builder of the plugin. `lang.Served` wraps a handler, so the handler runs with the part of the language of its match and leaves a subject of another language alone.

The SDK of eidos supplies the other helpers that more than one plugin needs. A handler that serves two kinds of match takes `sdk.SubjectMatch`. The call record of a double declares the type parameters that `rules.Mentioned` returns. The templates write type arguments through `typeargs`, and a TypeScript construction spells its class through `spellvalue`.

### The generators of Go and TypeScript

| Generator | Subject in Go | Subject in TypeScript | Output | Checks assert through |
|---|---|---|---|---|
| stub | An interface, or a struct with methods | An interface, or a class | A typed double that calls the runtime, and its test | assert-go and `@dokimi/assert` |
| builder | A struct, or an alias of one | A class | A builder with a setter for each field, and its test | assert-go and `@dokimi/assert` |
| enum | A defined type with its constants | An enum | The parse function, the text, the validity and the list of the variants, and their test | assert-go and `@dokimi/assert` |
| sentinel | The error values and the error types of a package | The error classes of a module | A test of the errors of each package, with the planted defects of its checks | assert-go and `@dokimi/assert` |

The rules of eidos supply the facts of both languages:

- An error type is a declaration whose rules state an error contract. In Go it is a struct with the method `Error`. In TypeScript it is a class that extends `Error`, directly or through other classes.
- An error value is a variable whose name the rules of its language accept, such as `ErrNotFound` in Go. The rules of TypeScript do not accept such a name, so the sentinel generator checks the error classes of a TypeScript module alone.

### The runtimes

| Repository | Package | Contents |
|---|---|---|
| `dokimasia/dokimi-go` | The module `go.dokimi.dev/dokimi-go`, with the packages `stub`, `fault`, `spec`, `conform` and `probe` | The double runtime: the record of each call, the order of results, strict mode and the expectations of calls. The fault strategies. The constants of the vendored conformance spec, the checks with their test handle, falsification and suite records, and the probes |
| `dokimasia/dokimi-typescript` | The package `@dokimi/runtime`, with the modules `stub/*`, `fault`, `spec`, `conform` and `probe`. Each module is a subpath export, such as `@dokimi/runtime/stub/answer` | The same contents in TypeScript |

The runtime contains the code that every generated double and every generated check would otherwise repeat. Each double reads the time from the clock of the assert library of its language, and it takes its random values from a source that the test passes, such as the source of a property case. RFC-0002 states the API of the conformance parts.

### The corpus

The corpus of every language has the same structure:

```text
corpus/<lang>/
└── <family>/[<group>/]<fixture>/
    ├── iface.go                 the declarations of the fixture, with their directives
    └── <fixture>test/           the directory that the fixture routes its files to
        ├── inmemory.go          the subject of the suite
        └── inmemory_test.go     the test that runs the suite against the subject
```

- The families are `enum`, `errors`, `iface`, `integration` and `struct`. The family `iface` has the groups `composite`, `contract`, `detector`, `lang` and `mixin`, and the other families have no group.
- A fixture has the same path, the same declarations and the same files in every language. The tree shows the names of Go, and TypeScript names the same files `iface.ts`, `inmemory.ts` and `inmemory.test.ts`.
- A fixture that a language cannot state, such as a Go channel in TypeScript, has the file `absent.md` in its directory of that language. The file states the reason.
- dokimi writes the generated files of each fixture into its routed directory. The generation gate fails while a committed file differs from a fresh run.

`corpus/go` is the module `go.dokimi.dev/dokimi/corpus/go`, with the hand-written fixtures of the corpus of testkit. `corpus/typescript` is the package `@dokimi/corpus`, with the same fixtures in TypeScript and a `tsconfig.json`.

### Tests

| Level | Check | What it proves |
|---|---|---|
| A plugin | `plugintest.RunPluginSuite`, and the records of its fixtures | Determinism, parallel dispatch, the lint of every template of both targets, and the records of each subject in each language |
| The composition | `workspacetest` over `testdata/tree`, with the goldens of both plans in `testdata/want` | The rendered files of both targets, byte for byte |
| The generated files | The toolchain adapters of the Go and the TypeScript satellites | The Go output type-checks and its tests pass. The TypeScript output type-checks and its tests pass under vitest |
| The binary | `acceptancetest` of eidos-cli | The exit statuses, the idempotence of a second run and the compile of the generated output |
| The corpus | `dokimi run --check` over each directory of `corpus/` | The committed generated files equal a fresh run |

### Requests

| Project | Request | Why |
|---|---|---|
| assert-typescript | A release of `@dokimi/assert` on npm | The generated TypeScript tests import it, and the toolchain check installs it |

### Invariants

| Invariant | Enforced by |
|---|---|
| A plugin and its tests do not import a package of eidos-core | depguard |
| `plugin.go` does not import a satellite of eidos, and `<lang>.go` imports the root package of its own satellite alone | depguard |
| Every generator declares a template tree for golang and for typescript | The Build of the composition |
| Each family of a generator claims its files through a template of each target | The template lint of plugintest |
| Every fixture of a corpus has the same path in the corpus of every other language, or the file `absent.md` | A test of the root package, which walks every corpus |
| The goldens of both plans equal a fresh run | workspacetest |
| A generated test passes in its language | The toolchain adapters |

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| `plugin.go` imports a satellite | depguard fails the lint | Move the code of that language into `<lang>.go` |
| A generator lacks the tree of a target | The Build of the composition fails and names the plugin and the target | Add the tree |
| A fixture of one corpus is missing in another | The test of the corpus fails and names the fixture and the language | Add the fixture, or `absent.md` with the reason |
| The goldens differ from a fresh run | workspacetest fails with the difference | Review the difference, then update the goldens with `-update` |
| A generated TypeScript test cannot import `@dokimi/assert` | The toolchain check fails | Release `@dokimi/assert` |
| A toolchain is missing on a machine | The check skips there with its reason, and CI requires it | Install the toolchain |

## Alternatives considered

### Wrappers of the reference generators of eidos

The builder and the sentinel generator of dokimi embed the generators of eidos-reference and replace their assertions.

**Why not:** the reference generators serve Go alone. Their checks call the package `testing` inline, so a wrapper has no helper to replace. eidos-reference is the canary of the ring of eidos, and its plugins change with the needs of eidos.

### One target first

Each generator serves Go. A TypeScript tree follows it.

**Why not:** the handlers of a generator for Go alone call the Go satellite, as the reference builder does. The TypeScript tree then has to undo those calls. A generator that serves both targets has neutral handlers from its first line.

### One file for the code of every language

A plugin has one file with its handlers, and the templates contain every difference between the languages.

**Why not:** a step that the facts of eidos leave to a language, such as the receiver of a Go method, is then code in the handlers, next to the generic code. A new language would edit the handlers of every plugin, and nothing would show which line belongs to which language.

### The runtimes in this repository

The runtimes are the directories `runtimes/<language>` of the dokimi repository, and each release of dokimi releases them.

**Why not:** a consumer installs the runtime from the registry of its language, so the runtime needs the baseline, the gate and the release flow of that language in a repository of its own.

### The composition under internal/

The composition is the package `internal/compose`.

**Why not:** `cmd/dokimi` and the tests of the composition import it, and a binary of a consumer that adds plugins composes over `Base`. An internal package hides `Base` from such a binary. The root package of eidos-reference is the same entry point for its ensemble.

### Packages by language

The directories `golang/` and `typescript/` contain the generators of each language.

**Why not:** a split by language would split each generator into two packages, and both packages would repeat its handlers.

## Drawbacks

- **Two template trees for every generator.** Each generator keeps a tree for golang and one for typescript, and a change of its output changes both.
- **A Node toolchain in CI.** The checks of the generated TypeScript need Node, tsc and vitest on every runner of CI.
- **The TypeScript checks wait on a release.** Until `@dokimi/assert` is on npm, the generated TypeScript tests cannot install it.
- **A tree for each new language.** A language that eidos reads with rules needs a template tree in every generator.

## Open questions

None.

## References

| What | Where |
|---|---|
| The authoring surface of eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/06b-authoring.md> |
| Testing and conformance in eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/13-testing-and-conformance.md> |
| Workspaces and plans in eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/08-workspace-and-plans.md> |
| stubgen, a generator of eidos-reference for Go and TypeScript | <https://github.com/dokimasia/eidos/tree/main/eidos-reference/stubgen> |
| The composition of eidos-reference | <https://github.com/dokimasia/eidos/blob/main/eidos-reference/reference.go> |
| The shared helpers of plugins in eidos | <https://github.com/dokimasia/eidos/blob/main/docs/rfc/0069-constructions-module-directives-and-the-shared-helpers-of-plugins.md> |
