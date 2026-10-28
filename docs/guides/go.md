<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# How to generate the test code of a Go package

dokimi reads directives in the doc comments of your Go code and writes test helpers, tests and
conformance suites beside the code. Write each directive as a line `//dokimi:<name>` in the doc
comment of its declaration, or of the package clause for a directive of the package.

## Choose a generator

| If you want to | Put this on | Directive | dokimi writes |
|---|---|---|---|
| Build test values of a struct with many fields | The struct | `builder` | `New<X>Builder` with a setter of each field, `Mutate`, `Clone` and `Build`, and its spec |
| Replace a dependency in a test | The interface or the struct | `stub` | `<X>Stub`, which records each call, returns the configured results and injects faults and latency, and its spec |
| Prove that each implementation of an interface keeps its contract | The interface, the struct or the function | `suite` | The checks, `Run<X>Suite`, the fuzz targets, the budgets and the benchmarks, the test of the planted defects, and the double of `stub` |
| Run a suite against one implementation | The struct of the implementation | `subject <X> new=<factory>` | `Test<Impl>ConformsTo<X>`, the benchmarks, and the allocation test where a budget caps allocations |
| Parse and print the constants of a defined type | The defined type | `enum` | `Parse<X>`, `String`, `IsValid`, `<X>Values`, `MarshalText` and `UnmarshalText`, and their spec |
| Keep the error values of a package distinct | The package clause | `sentinel` | The spec of each error value and each error type |

Leave a declaration without a directive when no test needs it. A generated helper is code that
you review and keep.

- Use `builder` on structs whose fields a test sets. A builder writes the fields directly, so it
  bypasses a constructor that enforces an invariant.
- Use `stub` on the dependencies of the code under test, not on the code under test itself.
- Use `suite` on an interface whose contract each implementation must keep, such as a store in
  memory and a store on disk, or on a callable that needs fuzz targets and budgets.

## Shape the generated code

Each of these directives changes what a generator writes, and none of them generates code itself.

| If you want to | Write |
|---|---|
| Move the generated files into a package of their own | `//dokimi:out path=storetest/` on the package clause or on the declaration |
| Start each builder from a declared value | `default:"1"` on a field, or `//dokimi:default "new"` on a defined type |
| Start each builder from a function | `//dokimi:builder defaults=DefaultConfig` |
| Give a type the value that the generated specs use | `//dokimi:sample value=DefaultNote()` |
| Pass one value of a parameter in each generated call of its callable, such as the part of a domain that its contract admits | `//dokimi:sample value=0x01` on the line before the parameter |
| Give a double helpers that fail a call with an error value of the package | `//dokimi:fault errors=[ErrNotFound]` on the method |
| Name the field that configures each method of a double, such as for a subject with the methods `Node` and `OnNode` | `//dokimi:stub field={{.Method}}Stub`, which gives `Node` the field `NodeStub` in place of `OnNode` |
| State a property that a suite checks | `//dokimi:mixin not-found sentinel=ErrNotFound`, `//dokimi:mixin concurrent` and the other mixins of the shape catalog, on the method |
| Leave a callable that waits for another call out of each benchmark and budget, such as a lock whose `Acquire` waits for `Release` | `//dokimi:mixin waits` on the method or the function |
| Add a check of your own to a suite | `//dokimi:check Store.Get defect=forgetful` on a function in the package of the generated suite, whose doc comment opens with "<Name> checks that" |
| Give the planted defect of a check a restart, a crash, a reference or a second subject | `restart=forget`, `crash=tear`, `reference=` or `second=` beside `defect`. Each value is a function of the package of the check |
| Require a share of the outcomes of a property check in each application | `options=machineOptions` beside `defect`, where `machineOptions` takes the handle and returns `[]prop.Option`, such as `prop.Require` |
| Run each run of a check inside a `testing/synctest` bubble, also on its planted defect | `around=bubble` beside `defect`, where `bubble` takes a `*testing.T` and a `func(*testing.T)` and calls `synctest.Test` |
| Pass the doors and other options of an application to the test of a subject | `//dokimi:subject store.Store new=New options=StoreOptions`, where `StoreOptions` returns the list of options |
| Excuse a check that an implementation does not pass | `//dokimi:excuse Store.Get/keeps-the-item reason="…"` beside `subject` |
| Cap the allocations, the bytes or the latency of a call | `//dokimi:budget allocs=1 bytes=96 p99=2us mean=1us warmup=500` on the method or the function |
| Take the text of each constant from its trailing comment | `//dokimi:enum text=comment` |
| Check that no error message overlaps a neighbouring package | `//dokimi:sentinel prefix=kv no-overlap-with=[example.com/shop/neighbour]` |

## Examples

Each example is a fixture of the corpus in `corpus/go`, with the generated files beside it.

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
//dokimi:stub
//dokimi:out path=readertest/
type Reader interface {
	//dokimi:mixin not-found sentinel=ErrNotFound
	//dokimi:fault errors=[ErrNotFound] retry=3 partition=key
	Get(ctx context.Context, key string) (Value, error)
}
```

```go
source := readertest.NewReaderStub(t)
source.OnGet.Returns(want, nil).RetrySchedule(errUnavailable)
got, err := retrying.New(source, 3).Get(t.Context(), want.Key)
assert.NoError(t, err, "the third attempt succeeds")
assert.Equal(t, source.OnGet.CallCount(), 3, "Get calls the source once for each attempt")
```

`FaultNotFound` arms `ErrNotFound` for the next call, because the directive `fault` lists it.

### Check each implementation of an interface

```go
//dokimi:suite
//dokimi:out path=causalchaintest/
type Log interface {
	Append(ctx context.Context, e Entry) error
	Replay(ctx context.Context) ([]Entry, error)
}
```

Bind a check of your own in the package of the suite. Its planted defect is the double of `Log`
with one method replaced:

```go
// KeepsTheBody checks that Replay returns the body that Append admitted.
//
//dokimi:check Log.Replay defect=losesBodies
func KeepsTheBody(h *LogHandle) {
	s := h.Subject()
	assert.NoError(h, s.Append(h.Context(), causalchain.RootEntry()), "Append admits the root entry")
	got, err := s.Replay(h.Context())
	assert.NoError(h, err, "Replay reads the log")
	assert.Equal(h, got[0].Body, causalchain.RootEntry().Body, h.Statement())
}

func losesBodies(h *LogHandle) causalchain.Log {
	var entries []causalchain.Entry
	return NewLogStub(h,
		WithLogAppend(func(_ context.Context, e causalchain.Entry) error {
			e.Body = ""
			entries = append(entries, e)
			return nil
		}),
		WithLogReplay(func(context.Context) ([]causalchain.Entry, error) {
			return slices.Clone(entries), nil
		}))
}
```

Apply the suite to an implementation in its own test:

```go
func TestLog(t *testing.T) {
	t.Parallel()

	causalchaintest.RunLogSuite(t, memory.New)
}
```

Or let dokimi write that test, and excuse each check that the implementation does not pass:

```go
//dokimi:subject causalchain.Log new=New
//dokimi:excuse Log.Replay/keeps-the-body reason="a replica refuses each append"
type Log struct {
	source causalchain.Log
}
```

The suite starts a struct without exported fields from its zero value. Pass `conform.Using` to an
application to state other values of such a type, such as a resolution that only a registry
builds. The generator of the option replaces the zero value in each run of the application:

```go
func TestApplier(t *testing.T) {
	t.Parallel()

	registry := resolved.NewRegistry(3)
	resolvedtest.RunApplierSuite(t, memory.New,
		conform.Using[resolved.Applier](prop.Just(registry.Resolve())))
}
```

The suite gives each type in the inputs of a callable a value. It passes a closed channel of drawn
elements for a receive channel and for a channel of both directions, so each receive returns. It
passes the sink of the test for a send channel. The sink receives and discards each value until
the test ends. It passes a `*conform.Bytes` over drawn bytes for an interface of the package `io`.
dokimi reports the warning `DOKIMISUITE-0010` at the declaration of the suite for each type
without a generator. An interface with an unexported method is such a type when the standard
library declares it or when its package does not export an implementation. Register a generator of
such a type with `prop.Register` in the test process and suppress the warning with
`//dokimi:diag off=DOKIMISUITE-0010`.

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
s, err := status.ParseStatus("Published") // errors.Is(err, ErrUnknownStatus) for any other text
s.String()                                // "Published"
status.StatusValues()                     // Draft, Published and Archived
```

An enumeration takes no directive `out`, because Go declares methods in the package of their type
alone.

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

## Run the generators

1. Add `.dokimi.yaml` with `version: 1` to the root of the module. Set the policy `golang.width`
   to the width at which your lint wraps lines, such as `"120"`.
2. Run `dokimi run ./...` from the root of the module. The run writes `<file>_<generator>.go` and
   `<file>_<generator>_test.go` beside the source, or into the package of `out`.
3. Run `go test ./...`. The generated specs check the generated code, and a suite falsifies each of
   its checks on a planted defect.
4. Commit the generated files with the directives that produce them.
5. Run `dokimi run --check --strict ./...` in CI. It fails on a file that a run would change and on
   every warning.

If a run reports a finding, run `dokimi explain` on the file or the declaration to see what the
run recorded about it. Run `dokimi doctor` to check the config and the directives. Run
`dokimi prune` to remove each generated file whose directive you removed.

Never edit a generated file. A run refuses to write over a file that was edited after its stamp.
Change the directive or the source instead, and run again.
