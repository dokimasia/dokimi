<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# dokimi

[![CI](https://github.com/dokimasia/dokimi/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/dokimasia/dokimi/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/dokimasia/dokimi)](https://github.com/dokimasia/dokimi/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/go.dokimi.dev/dokimi.svg)](https://pkg.go.dev/go.dokimi.dev/dokimi)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

dokimi generates the tests and the test helpers of your code from directives in its doc comments.
Mark a declaration with a directive, such as `//dokimi:builder` on a struct, and run `dokimi run`.
dokimi reads the code through [eidos](https://github.com/dokimasia/eidos) and writes the generated
files beside it.

A test that cannot fail looks the same as a test that works, because both pass. dokimi pairs each
check of a conformance suite with a planted defect. A planted defect is an implementation that is
correct except for the one rule that its check verifies. The test of the suite applies each check
to its planted defect and fails when the check passes on it.

## Install

Each [release](https://github.com/dokimasia/dokimi/releases/latest) has archives for Linux, macOS
and Windows on amd64 and arm64, deb, rpm and apk packages for Linux, and `checksums.txt` with the
SHA-256 digest of each file. cosign signs `checksums.txt` without a key, and each archive and
package has an SPDX SBOM. Download the archive for your system, check its digest, and put `dokimi`
on your `PATH`. `gh attestation verify <file> --repo dokimasia/dokimi` checks the build provenance
of a file.

With Homebrew, install the cask of the tap:

```sh
brew install --cask dokimasia/tap/dokimi
```

The cask installs the completions of bash, zsh and fish. The binary is not notarized, so on macOS
the cask removes the quarantine attribute.

To install the command from source, you need Go 1.27.3 or later and a C compiler, because dokimi
parses TypeScript with the C grammar of tree-sitter:

```sh
go install go.dokimi.dev/dokimi/cmd/dokimi@latest
```

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

`TestStoreSuite` applies the check to `forgetful`. It passes only when the check fails there at the
assertion whose message is `h.Statement()`. `TestStoreConformsToStore` then applies the check to
the implementation in memory.

Commit the generated files with the directives, and add `.dokimi/` to `.gitignore`. A run keeps
its state in that directory. In CI, fail the build when a generated file is out of date:

```sh
dokimi run --check --strict ./...
```

`--check` does not write a file and reports each file that a run would change. `--strict` reports
each warning as an error. Never edit a generated file, because a run refuses to write over a file
that was edited after its stamp. Change the directive or the source instead, and run again.

## What dokimi generates

| Directive | Put it on | dokimi writes |
|---|---|---|
| `suite` | An interface, a struct or a function | The suite, `Run<X>Suite`, the fuzz targets, the budgets, the benchmarks and the test that applies each check to its planted defect |
| `check` | A function in the package of the suite | The binding of the function and its planted defect to a callable of the suite |
| `subject` | The struct of an implementation | `Test<Impl>ConformsTo<X>`, the benchmarks, the fuzz targets, and the allocation test where a budget caps allocations |
| `stub` | An interface or a struct | `<X>Stub`, which records each call, returns the configured results and injects faults and latency, and its spec |
| `builder` | A struct | `New<X>Builder` with a setter of each field, `Mutate`, `Clone` and `Build`, and its spec |
| `enum` | A defined type with constants | `Parse<X>`, `String`, `IsValid`, `<X>Values`, `MarshalText` and `UnmarshalText`, and their spec |
| `sentinel` | The package clause | The spec of each error value and each error type |

The [Go guide](docs/guides/go.md) shows an example of each generator and lists the directives that
change what a generator writes, such as `out`, `budget`, `fault`, `mixin` and `excuse`.

## Languages

dokimi runs one plan for each language. The plan of Go reads the Go packages of the module and
runs every generator. The plan of TypeScript reads the TypeScript sources and runs `builder`,
`stub`, `enum` and `sentinel`. Its output imports the packages `@dokimi/assert` and
`@dokimi/runtime`. npm does not publish either package.

The key `sources` of a plan in `.dokimi.yaml` replaces the packages that the plan reads, and
`enabled: false` turns the plan off:

```yaml
version: 1
plans:
  go:
    sources: {lang: golang, packages: [./internal/...]}
  typescript:
    enabled: false
```

`dokimi schema` prints the JSON Schema of `.dokimi.yaml`.

## The runtime of the generated code

The generated Go code imports the modules of [dokimi-go](https://github.com/dokimasia/dokimi-go)
and the assertions of [assert-go](https://github.com/dokimasia/assert-go):

| Module | Imported by |
|---|---|
| [`go.dokimi.dev/conform`](https://pkg.go.dev/go.dokimi.dev/conform) | Each suite and the test of each suite |
| [`go.dokimi.dev/stub`](https://pkg.go.dev/go.dokimi.dev/stub) | Each double |
| [`go.dokimi.dev/fault`](https://pkg.go.dev/go.dokimi.dev/fault) | The fault strategies of each double |
| [`go.dokimi.dev/assert`](https://pkg.go.dev/go.dokimi.dev/assert) | Each spec and each check |

Each generated suite refers to the function `conform.CompatV1`. A major version of conform renames
that function, so a suite that dokimi wrote for another major version does not compile.

## Command line

```text
dokimi <command> [options]
```

| Command | What it does |
|---|---|
| `run` | Generates the output of each workspace and writes it |
| `plan` | Prints the composition of each workspace without running it |
| `explain` | Prints what the last run recorded about a file, a declaration, a fact or a finding |
| `prune` | Removes the generated files that no plan produces |
| `doctor` | Checks the config, the sealed state and the directives of each workspace |
| `watch` | Runs each workspace again after each change, until an interrupt |
| `version` | Prints the versions of the binary, of eidos and of the plugins of each workspace |
| `schema` | Prints the JSON Schema of the config file |
| `completion` | Writes the completion script of bash, zsh, fish or PowerShell |

Each option also reads an environment variable, such as `DOKIMI_STRICT` for `--strict`.
`dokimi <command> --help` lists the options of a command.

## Documentation

| Directory | Contents |
|---|---|
| [docs/guides](docs/guides) | How to use dokimi in a project of one language |
| [docs/rfc](docs/rfc) | The proposals and the argument behind them |
| [docs/adr](docs/adr) | The decisions and what they cost |

The corpus in [corpus/go](corpus/go) contains the fixtures of the generators, with the generated
files beside each fixture. The examples of the Go guide come from it.

## Development

The repository uses [ergon](https://github.com/dokimasia/ergon) for its checks and its releases.
[CONTRIBUTING.md](CONTRIBUTING.md) describes how to propose a change. The project follows the
[Contributor Covenant 2.1](CODE_OF_CONDUCT.md).

```sh
make check   # the gate: lint, tests, race detector and vulnerability scan of every module
make test    # the tests of every module
make help    # every target
```

## Security

Report a vulnerability by email to security@dokimi.dev, or through a private security advisory of
the repository. Do not open a public issue. [SECURITY.md](SECURITY.md) describes what a report
contains.

## License

Dokimasia B.V. licenses dokimi under the Apache License 2.0. See [LICENSE](LICENSE) and
[NOTICE](NOTICE).
