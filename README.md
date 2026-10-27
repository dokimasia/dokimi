<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# dokimi

[![CI](https://github.com/dokimasia/dokimi/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/dokimasia/dokimi/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/dokimasia/dokimi)](https://github.com/dokimasia/dokimi/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/go.dokimi.dev/dokimi.svg)](https://pkg.go.dev/go.dokimi.dev/dokimi)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/dokimasia/dokimi/badge)](https://scorecard.dev/viewer/?uri=github.com/dokimasia/dokimi)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

dokimi generates the tests and the test helpers of your Go code from directives in its doc
comments. It proves that each check of a generated suite can fail.

For an interface marked `//dokimi:suite`, dokimi writes a conformance suite that every
implementation of the interface must pass. A struct marked `//dokimi:builder` gets a fixture
builder with its spec. The generated files are plain Go that `go test` runs.

A test that cannot fail looks the same as a test that works, because both pass. dokimi pairs each
check of a suite with a planted defect. A planted defect is an implementation that is correct
except for the one rule that its check verifies. `go test` applies each check to its planted defect
and fails when the check passes on it.

## Features

- Conformance suites of interfaces, structs and functions, with fuzz targets, allocation budgets
  and benchmarks for each method or function.
- Checks that must each fail on a planted defect, so a check that cannot fail fails the test of
  its suite.
- Doubles that record each call, return the results that a test configures, and inject faults and
  latency.
- Fixture builders, the parse and print functions of enumerations, and the specs of error values.
- Generated code that imports the runtime modules and the assertions alone, and runs under
  `go test` without dokimi.
- A gate for CI that fails the build on a generated file that is out of date.
- Releases for Linux, macOS and Windows with signed checksums, SPDX SBOMs and build provenance.

## Install

With Homebrew on macOS or Linux:

```sh
brew install --cask dokimasia/tap/dokimi
```

The cask installs the completions of bash, zsh and fish. The binary is not notarized, so on macOS
the cask removes the quarantine attribute.

With Go 1.27.3 or later and a C compiler:

```sh
go install go.dokimi.dev/dokimi/cmd/dokimi@latest
```

The build needs the C compiler because dokimi parses TypeScript with the C grammar of tree-sitter.

Each [release](https://github.com/dokimasia/dokimi/releases/latest) also has archives for Linux,
macOS and Windows on amd64 and arm64, and deb, rpm and apk packages for Linux. `checksums.txt`
contains the SHA-256 digest of each file, and cosign signs it without a key. Each archive and
package has an SPDX SBOM. `gh attestation verify <file> --repo dokimasia/dokimi` checks the build
provenance of a file.

`dokimi version` prints the versions of the binary, of eidos and of each plugin.

## Quick start

This example generates the conformance suite of an interface `Store` with one check, and applies
the suite to an implementation in memory.

1. Add `.dokimi.yaml` to the root of the Go module:

   ```yaml
   version: 1
   ```

2. Mark the interface in `store/store.go`. `suite` writes its suite, `stub` writes its double, and
   `out` puts both into the package `storetest`:

   ```go
   // Store keeps items by key.
   //
   //dokimi:suite
   //dokimi:stub
   //dokimi:out path=storetest/
   type Store interface {
   	Put(ctx context.Context, key, item string) error
   	Get(ctx context.Context, key string) (string, error)
   }
   ```

3. Write a check in `store/storetest/checks.go`. The doc comment of a check opens with
   `<Name> checks that`, and `defect` names the function that returns its planted defect. Here
   the planted defect is the double of `Store` with a Get that returns an empty item:

   ```go
   // ReturnsThePutItem checks that Get returns the item that Put stored.
   //
   //dokimi:check Store.Get defect=forgetful
   func ReturnsThePutItem(h *StoreHandle) {
   	s := h.Subject()
   	assert.NoError(h, s.Put(h.Context(), "k", "v"), "Put stores the item")
   	got, err := s.Get(h.Context(), "k")
   	assert.NoError(h, err, "Get finds the item")
   	assert.Equal(h, got, "v", h.Statement())
   }

   // forgetful returns a store whose Get returns an empty item for each key.
   func forgetful(h *StoreHandle) store.Store {
   	return NewStoreStub(h, WithStoreGet(func(context.Context, string) (string, error) {
   		return "", nil
   	}))
   }
   ```

4. Mark each implementation with the interface that it implements and its constructor, such as
   the type `Store` of `memory/memory.go`:

   ```go
   // Store keeps items in a map.
   //
   //dokimi:subject store.Store new=New
   type Store struct {
   	items map[string]string
   	mu    sync.Mutex
   }
   ```

5. Generate the code and run the tests:

   ```sh
   dokimi run ./...
   go mod tidy
   go test ./...
   ```

`dokimi run` writes the suite and the double into `store/storetest`, and the test
`TestStoreConformsToStore` into `memory/memory_suite_test.go`. The generated code imports
`go.dokimi.dev/assert`, `go.dokimi.dev/conform` and `go.dokimi.dev/stub`, and `go mod tidy` adds
them to `go.mod`.

`go test -v ./...` shows both runs of the check:

```text
--- PASS: TestStoreConformsToStore (0.00s)
    --- PASS: TestStoreConformsToStore/Get/returns-the-put-item (0.00s)
--- PASS: TestStoreSuite (0.00s)
    --- PASS: TestStoreSuite/StoreSuite (0.00s)
        --- PASS: TestStoreSuite/StoreSuite/Get/returns-the-put-item (0.00s)
```

`TestStoreSuite` applies the check to `forgetful`. It passes only when the check fails there at the
assertion whose message is `h.Statement()`. `TestStoreConformsToStore` applies the check to the
implementation in memory.

Commit the generated files with the directives, and add `.dokimi/` to `.gitignore`.

## Generators

| Directive | Put it on | dokimi writes |
|---|---|---|
| `suite` | An interface, a struct or a function | The suite, `Run<X>Suite`, the fuzz targets, the budgets, the benchmarks and the test that applies each check to its planted defect |
| `check` | A function in the package of the suite | The binding of the function and its planted defect to a method or a function of the suite |
| `subject` | The struct of an implementation | `Test<Impl>ConformsTo<X>`, the benchmarks, the fuzz targets, and the allocation test where a budget caps allocations |
| `stub` | An interface or a struct | `<X>Stub`, which records each call, returns the configured results and injects faults and latency, and its spec |
| `builder` | A struct | `New<X>Builder` with a setter of each field, `Mutate`, `Clone` and `Build`, and its spec |
| `enum` | A defined type with constants | `Parse<X>`, `String`, `IsValid`, `<X>Values`, `MarshalText` and `UnmarshalText`, and their spec |
| `sentinel` | The package clause | The spec of each error value and each error type |

dokimi writes `<file>_<generator>.go` and `<file>_<generator>_test.go` beside the source, or into
the package that `//dokimi:out` names. The [Go guide](docs/guides/go.md) lists the directives that
change what a generator writes, such as `budget`, `fault`, `mixin` and `excuse`.

### Build a test value

```go
//dokimi:out path=plaintest/
package plain

//dokimi:builder
type Item struct {
	ID    string
	Tags  []string
	Count int
}
```

```go
item := plaintest.NewItemBuilder().WithID("a").AppendTags("new").WithCount(2).Build()
```

### Replace a dependency

```go
store := storetest.NewStoreStub(t)
store.OnGet.Returns("v", nil)
store.OnPut.FaultsOnce(errBusy)
```

The double records each call. `FaultsOnce` fails the next call of Put with `errBusy`.

### Parse an enumeration

```go
//dokimi:enum
type Status int

const (
	Draft Status = iota + 1
	Published
	Archived
)
```

```go
s, err := ParseStatus("Published") // errors.Is(err, ErrUnknownStatus) for any other text
s.String()                         // "Published"
StatusValues()                     // Draft, Published and Archived
```

### Check the error values of a package

```go
//dokimi:sentinel prefix=kv
package store

var (
	ErrNotFound = errors.New("kv: not found")
	ErrConflict = errors.New("kv: conflict")
)
```

The generated spec fails where a message does not open with `kv:`, where two messages are the same,
and where one error value matches another.

## Use dokimi in CI

Run the gate of generation before the tests:

```yaml
- uses: actions/setup-go@v7
  with:
    go-version-file: go.mod
- run: go install go.dokimi.dev/dokimi/cmd/dokimi@latest
- run: dokimi run --check --strict ./...
- run: go test ./...
```

`--check` does not write a file and reports each file that a run would change as an error.
`--strict` reports each warning as an error. A generated file that is out of date fails the step:

```text
store/storetest/store_suite.go: error EID-0064: store/storetest/store_suite.go is out of date: plan "go" would have updated it (close)
```

## Configuration

`.dokimi.yaml` at the root of the module configures a run:

```yaml
version: 1
policies:
  golang.width: "120"
plans:
  go:
    sources: {lang: golang, packages: [./internal/...]}
  typescript:
    enabled: false
```

- `golang.width` is the width at which the generated Go wraps its lines. Set it to the width of
  your linter.
- The key `sources` of a plan replaces the packages that the plan reads, and `enabled: false` turns
  the plan off.
- The key `options` sets the options of each generator, such as `options.enum.parse-function`.

`dokimi schema` prints the JSON Schema of the file, and `dokimi doctor` checks the file and the
directives.

## Commands

| Command | What it does |
|---|---|
| `dokimi run [PATTERNS...]` | Generates the files of the directives in the packages of the patterns |
| `dokimi watch` | Runs again after each change, until an interrupt |
| `dokimi explain` | Prints what the last run recorded about a file, a declaration, a fact or a finding |
| `dokimi doctor` | Checks the config, the sealed state and the directives |
| `dokimi prune` | Removes the generated files that no plan produces |
| `dokimi plan` | Prints the composition of each workspace without running it |
| `dokimi schema` | Prints the JSON Schema of `.dokimi.yaml` |
| `dokimi version` | Prints the versions of the binary, of eidos and of the plugins |
| `dokimi completion SHELL` | Writes the completion script of bash, zsh, fish or PowerShell |

Each option also reads an environment variable, such as `DOKIMI_STRICT` for `--strict`.
`dokimi <command> --help` lists the options of a command.

## How it works

dokimi reads your code through the frontends of [eidos](https://github.com/dokimasia/eidos) and runs
one plan for each language. Each generator turns the directives of a declaration into files, and
the backend of the language writes them.

- Each generated file opens with `// Code generated by dokimi. DO NOT EDIT.` and ends with a
  provenance stamp. A run refuses to write over a file that was edited after its stamp. Change the
  directive or the source instead, and run again.
- A run records its state in `.dokimi/`, and the next run regenerates only what changed. A run
  without that state writes the same bytes.
- The test of a suite calls `conform.Falsify`, which applies each check to its planted defect. It
  accepts a failure only at the assertion whose message is the statement of the check, so a check
  that fails for another reason also fails the test.
- `dokimi prune` removes the generated files whose directive you removed.

The generated Go code imports these modules:

| Module | Imported by |
|---|---|
| [`go.dokimi.dev/conform`](https://pkg.go.dev/go.dokimi.dev/conform) | Each suite and the test of each suite |
| [`go.dokimi.dev/stub`](https://pkg.go.dev/go.dokimi.dev/stub) | Each double |
| [`go.dokimi.dev/fault`](https://pkg.go.dev/go.dokimi.dev/fault) | The fault strategies of each double |
| [`go.dokimi.dev/assert`](https://pkg.go.dev/go.dokimi.dev/assert) | Each spec and each check |

[dokimi-go](https://github.com/dokimasia/dokimi-go) contains the first three, and
[assert-go](https://github.com/dokimasia/assert-go) contains the assertions. Each generated suite
refers to the function `conform.CompatV1`. A major version of conform renames that function, so a
suite that dokimi wrote for another major version does not compile.

## Comparison

- **Mock generators.** gomock and mockery generate mocks of interfaces. dokimi generates the double
  of an interface or a struct with fault strategies, such as a share of failing calls or a window
  of time, and a latency. The spec of the double checks the double itself.
- **stringer.** stringer generates the method `String` of a type with constants. The enumeration of
  dokimi also has `Parse<X>`, `IsValid`, `<X>Values`, `MarshalText` and `UnmarshalText`, and their
  spec.
- **Property-based testing libraries.** rapid runs the properties that you write and shrinks each
  failure. The property checks of a suite run on the property engine of assert-go. The suite also
  applies each check to each implementation and proves that each check fails on its planted defect.
- **Conformance helpers.** `fstest.TestFS` and `iotest.TestReader` of the standard library check one
  interface each. dokimi writes the suite of each interface that you mark.

## FAQ

### Does a test need dokimi to run?

No. The generated files import the runtime modules and the assertions alone, so `go test` runs
them on a machine without the binary.

### How does a planted defect differ from a mutant?

A mutation tool such as [dokimi-mutate-go](https://github.com/dokimasia/mutate-go) changes an
expression of your code and reports whether any test failed. The planted defect of a check breaks
the one rule that the check verifies. Its run shows that this check fails for the reason that it
states, and it costs one run of that check.

### What if an implementation does not pass a check?

Its test fails with the message of the assertion that failed. When the difference is intended,
excuse the check with its reason beside the directive `subject`:

```go
//dokimi:subject store.Store new=New
//dokimi:excuse Store.Get/returns-the-put-item reason="a replica returns the item of its last sync"
```

The test of the subject then skips the check and reports the reason.

### What if a generated file needs a change?

Change the directive or the source, and run `dokimi run` again. A run refuses to write over a
generated file that was edited after its stamp, and `dokimi run --overwrite-drift` writes over it.

### Does dokimi generate TypeScript?

The plan of TypeScript runs `builder`, `stub`, `enum` and `sentinel` on the TypeScript sources of a
workspace. Its output imports the packages `@dokimi/assert` and `@dokimi/runtime`. npm does not
publish either package.

## Status

dokimi is before its 1.0 release. A minor release can change the generated code, and the
[changelog](CHANGELOG.md) starts each such entry with **Breaking:**. The binary builds with Go
1.27.3 or later, and the generated code needs Go 1.27.2 or later. The runtime modules have
versions of their own.

## Documentation

| Directory | Contents |
|---|---|
| [docs/guides](docs/guides) | How to use dokimi in a project of one language |
| [docs/rfc](docs/rfc) | The proposals and their arguments |
| [docs/adr](docs/adr) | The decisions and what they cost |
| [corpus/go](corpus/go) | The fixtures of the generators, with the generated files beside each fixture |

## Get help and contribute

- To report a defect, open an issue with the
  [bug report form](https://github.com/dokimasia/dokimi/issues/new?template=bug.yml).
- To propose a capability, open an issue with the
  [feature request form](https://github.com/dokimasia/dokimi/issues/new?template=feature.yml).
- To change the code, read [CONTRIBUTING.md](CONTRIBUTING.md). The repository uses
  [ergon](https://github.com/dokimasia/ergon) for its checks and its releases, and `make check`
  runs the gate of every module.

Everyone who takes part follows the [Contributor Covenant 2.1](CODE_OF_CONDUCT.md).

## Security

Report a vulnerability by email to security@dokimi.dev, or through a private security advisory of
the repository. Do not open a public issue. [SECURITY.md](SECURITY.md) describes what a report
contains.

## License

Dokimasia B.V. maintains dokimi and licenses it under the Apache License 2.0. See
[LICENSE](LICENSE) and [NOTICE](NOTICE).
