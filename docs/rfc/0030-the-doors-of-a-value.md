---
rfc: 0030
title: The doors of a value
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-01
updated: 2026-11-01
discussion: https://github.com/dokimasia/dokimi/issues/49
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0030: The doors of a value

## Summary

A check receives only subjects from the doors of its handle: a reference implementation, a second subject, a restart and a crash. A check of a store also needs values that only the code that builds the store can supply, such as the client of the backend of the store. conform gets the door of a value of any type. `conform.Supply` returns the door of a value of a type `V` from the function that opens it for a subject. An application passes the door through the option `conform.Supplies`, a check passes the door of its planted defect through the field `Values` of `conform.Doors`, and `conform.Open[V]` returns the value for a subject. The directive `check` gets the param `values`, which lists the functions of the doors of values of its planted defect.

## Motivation

`conform.Doors` declares the doors of a handle as functions that return a subject of the type `S` of the suite (`conform/handle.go:58-75` at 494250b). A check cannot receive a value of another type, also when only the code that builds the subject can supply it.

core designs a capability of `blob.Store` that signs the requests through which a client without the credentials of the store reads or creates one object (thesm-os/core#109). Each check of the suites of `blob.GetSigner` and `blob.CreateSigner` sends a signed request to the backend of its subject, as a client without credentials does. The check then compares the response with the state of the store. The test of the store supplies a client of the backend, the time of the clock of the store and a wait on the clock of the backend:

- The store in memory serves its requests from an httptest server, and a store over Amazon S3 needs a client with the roots of the system.
- A check chooses each expiry from the time of the clock, because the store refuses a lifetime below one second or above 7 days.
- The wait lets a request pass its expiry. The test of the store in memory advances its fake clock, and the test of a store over S3 waits real time.

The suite in core declares the three as one interface, `blobtest.Backend`, with the methods `Do`, `Time` and `Elapse`. Each check of these suites has a planted defect over the store in memory, and the defect serves its requests from its own httptest server. So the check must give the defect a backend of its own, as a check gives a planted defect that restarts a restart of its own (dokimasia/dokimi-go#22, dokimasia/dokimi#49).

The handle supplies every other door, and the suite generator writes the doors of each planted defect into its check. No layer above conform and the suite generator sees a handle, so the change belongs in these two.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| A handle on a planted defect opens the value of the door of its check, and the value of a door of the options of the falsification where the check has no door of that type | `Doors.over` and `Open` of conform |
| A handle on a subject of an application opens the values of the doors of the options of the application alone | `newHandle` of conform |
| A door that no option and no check supplies ends the check with a fault that wraps `ErrOption` and names the remedy | `Open` |
| A check supplies one door for each type of value | The problems of an application in conform, and `BadCheck` of the suite generator |
| A door function of a check takes the handle and the planted defect, and returns one value | The annotate step of the suite generator, under `BadCheck` |

### Components

| Component | Change |
|---|---|
| `conform.Supplied[S]` | A new type: the door of a value of one type, with the function that opens it |
| `conform.Supply` | A new function that returns the door of a value |
| `conform.Supplies` | A new option that supplies the doors of values to an application |
| `conform.Open` | A new function that returns the value of a door for a subject |
| `conform.Doors[S]` | The field `Values`, the doors of values of a planted defect |
| The problems of an application | Two doors of one type in the doors of a check |
| The directive `check` | The param `values` |
| The codes of the suite generator | Doors of values in the meaning of `BadCheck` |
| The template of the suites | The field `Values` of the doors of each check with the param |

### The doors of a value

A door of a value is identified by the type of its value. A suite declares the type once, such as the interface `Backend` of the kit of a store:

```go
// Supplied is the door of a value of one type: the function that opens the
// value for a subject. An application passes it through the option
// [Supplies], and a check passes it to its planted defect through the field
// Values of [Doors]. [Supply] returns it.
type Supplied[S any] struct {
	of   reflect.Type
	open func(h *Handle[S], s S) any
}

// Supply returns the door of a value of the type V, which open opens for a
// subject s, such as the client of the backend of s.
func Supply[S, V any](open func(h *Handle[S], s S) V) Supplied[S]

// Supplies supplies the doors of values of the subjects of an application. A
// door replaces the door of its type that an earlier option supplies.
func Supplies[S any](doors ...Supplied[S]) Option[S]

// Open returns the value of the type V for the subject s, from the door of
// that type of the planted defect of the check, or from the option [Supplies]
// of the application. Without the door it ends the check with a fault that
// wraps [ErrOption], and it returns the zero value to a seat that takes the
// fault and goes on.
func Open[V, S any](h *Handle[S], s S) V
```

`Doors` gains the doors of values:

```go
type Doors[S any] struct {
	// The fields Reference, Second, Restart, Crash and Faults as before.

	// Values are the doors of values of the planted defect, one for each type,
	// each from [Supply]. [Open] returns the value of the door of its type.
	Values []Supplied[S]
}
```

A check opens a door for the subject that it calls:

```go
// SignGetServesTheObject checks that the backend serves the body that Put
// stored for a request that SignGet signed.
func SignGetServesTheObject(h *GetSignerHandle) {
	s := h.Subject()
	…
	response, err := conform.Open[Backend](h, s).Do(request)
	…
}
```

`Open` calls the function of the door each time with the handle and the subject. A function that returns the backend of the store `s` returns one backend for two calls with one store. `Open` calls the function as a factory builds a subject, so a read of the clock of the handle in the function is a read by a subject, and a backend over the clock of the handle can call `Handle.Advance`. Falsification resolves the door of a type on a handle of a planted defect in the order of the doors of a subject: the door of the check, then the door of the options of the falsification.

The fault of a handle on a planted defect contains the door of the check:

```text
conform: the handle lacks the door: GetSigner.SignGet/serves-the-object opens a value of the type blobtest.Backend
on its planted defect, and the check has no door of the type. Give the check the door:
	conform.Doors[blob.GetSigner]{Values: []conform.Supplied[blob.GetSigner]{conform.Supply(func(h *conform.Handle[blob.GetSigner], s blob.GetSigner) blobtest.Backend { … })}}
```

On a subject of an application, the fault contains the option instead: `… without the option conform.Supplies. Pass the option to the application: conform.Supplies(conform.Supply(…))`.

### The directive `check`

The directive `check` takes the param `values`, a list of functions of the package of the check, as the value of `defect` is. The param has the name of the field of `conform.Doors`, as the param `restart` has. Each function takes the handle of a check and the planted defect, and returns the value of its door:

```go
// SignGetServesTheObject checks that the backend serves the body that Put
// stored for a request that SignGet signed.
//
//dokimi:check blob.GetSigner.SignGet defect=servesAnotherObject values=[backendOf]
func SignGetServesTheObject(h *GetSignerHandle) { … }

// backendOf returns the backend of the store in memory that the planted
// defect d wraps.
func backendOf(_ *GetSignerHandle, d blob.GetSigner) Backend { return d.(anotherObject).Store.Backend() }
```

The annotate step stamps the identities of the functions as the fact `suite.values` of the check. It reports `BadCheck` at the directive in each of these cases:

- A function takes another number of parameters than two.
- A function returns another number of values than one.
- A function returns a value of the type of another function of the list.

The generator writes the doors into the check of the suite:

```go
{
	ID:        GetSignerChecks.SignGet.ServesTheObject(),
	Statement: "the backend serves the body that Put stored for a request that SignGet signed",
	Apply:     SignGetServesTheObject,
	Defect:    servesAnotherObject,
	Doors: conform.Doors[blob.GetSigner]{
		Values: []conform.Supplied[blob.GetSigner]{conform.Supply(backendOf)},
	},
},
```

A subject passes the doors of its application through the function of its param `options`, so the directive `subject` does not change:

```go
// StoreOptions returns the door of the backend of a store in memory.
func StoreOptions() []blobtest.GetSignerOption {
	return []blobtest.GetSignerOption{conform.Supplies(conform.Supply(backendOf))}
}
```

### The contract in each language

| Part | Go | TypeScript |
|---|---|---|
| The door of a value | `conform.Supply(open)` | `supply(Type, open)` |
| The doors of an application | The option `Supplies(doors...)` | The field `supplies` of `Options<S>` |
| The doors of a planted defect | The field `Values` of `Doors[S]` | The field `values` of `Doors<S>` |
| The value of a door | `conform.Open[V](h, s)` | `handle.open(Type, s)` |
| The param of the directive | `values` of `check` | The same param |

A TypeScript value has no type at run time, so its door takes the class or the token of the value in place of the type. The TypeScript runtime and its template tree do not exist yet. Their parts take these names when they do.

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A check opens a door that its planted defect lacks, in the spec of its suite | The falsification of the check ends with a fault that wraps `ErrOption` and names the door of the check | Give the check the door, such as `values=[backendOf]` |
| A check opens a door that a subject lacks | The check ends with a fault that wraps `ErrOption` and names the option | Pass the option, or excuse the check |
| The doors of a check supply two doors of one type | The application reports the problem and runs no check | Remove one door |
| A function of `values` has another signature, or returns the type of another function of the list | `BadCheck` at the directive. The check is not bound | Correct the function |
| Two doors of a suite return values of one type | A check cannot open both | Declare a type for each, such as `type ReplicaClient struct{ *http.Client }` |
| Another check calls `Open` on a planted defect that lacks the door | The isolation gives the check no verdict on the defect, as for every other door | Give the defect the door |

### Verification

- **conform.** The spec of the handle covers `Open` with the door of a check, with the door of the options on a defect whose check lacks it, on a subject of an application, and without a door on a defect and on an application, with the message of each fault. The spec of the options covers `Supplies`, a door that replaces an earlier one and the problem of two doors of one type in a check. The spec of a check covers the doors of values of a defect in the isolation of another check.
- **The suite generator.** The specs of the plugin cover the param `values`, each case of `BadCheck` of a door of a value, and the field `Values` of a check. The golden files of the fixture tree pin the generated check.
- **The corpus.** The fixture `iface/door/signer` applies a suite with a check of a signed request to a signer in memory, whose backend refuses a request that another signer signed. The signer passes the door of its backend through `options`, and the check gives its planted defect a backend of its own.

### Migration

The change adds a type, three functions, a field, a problem and a param, so no existing call breaks. A check without the param has no field `Values`, so the generated suites of every other check do not change. A generated suite with the field requires a release of conform with the doors of values.

## Alternatives considered

### A door with a name

Each door is a variable of the suite, such as `var BackendDoor = conform.NewDoor[S, Backend]("backend")`. `conform.Supply(BackendDoor, open)` returns the option of an application, and `conform.Open(h, BackendDoor, s)` returns the value in a check, as dokimasia/dokimi-go#22 proposes.

**Why not:** the directive `check` would have to name the door and its function, and the params of a directive have no type for a pair. The generator would find the door of a function by the type of its value, as the type identifies the door here. The name adds a variable and a lookup to each door. The cost of the type as the identity is a type for each of two doors of one type, such as the client of a primary and the client of a replica.

### A harness as the subject

The suite applies to a harness that contains the store and the backend, as the `Harness` of the Go CDK supplies `HTTPClient` beside the driver (`blob/drivertest/drivertest.go:60` at v0.47.0).

**Why not:** the type of the subjects is the contract of the declaration. Every derived check of the suite would apply to the harness and not to the store, and every planted defect would wrap the harness.

### A method of the subject

The store in memory declares `Backend() Backend`, and a check asserts the method on its subject.

**Why not:** the runtime would assert an undeclared interface on a value of the consumer's type, and a misspelled method would end the check with the fault of a missing door without a word about the spelling. A store over S3 would declare a method that only its test calls.

### A value of the context of the check

An application stores the value in the context of each check, through `context.WithValue`, and a check reads it from `h.Context()`.

**Why not:** a check has one context for all its subjects, so a check with two subjects would read one backend for both. A check would read a missing value as nil and fail on it, without the remedy of the fault of a door.

## Drawbacks

- **A sixth door in two places.** An application supplies the door of a value through an option, and a check supplies the door of its planted defect through `Doors`, as for the other five doors. A reader of a suite looks in both places to learn which value a handle opens.
- **A type for each door.** A suite with doors of values of one type declares a type for each door, because the type identifies the door.
- **A backend in each defect.** A planted defect over the store in memory needs a backend that serves its own requests, so each such defect builds a second backend.
- **A release of conform first.** A generated suite with the field `Values` compiles only against a release of conform with the doors of values.
- **One more param.** The directive `check` grows from ten params to eleven.

## Open questions

None.

## Unresolved and future work

- No derived claim of the conformance spec opens a door of a value. A claim that needs one is not proposed here.
- The TypeScript runtime and the TypeScript template tree of the suite generator do not exist. The table of the contract fixes their names.

## References

| What | Where |
|---|---|
| The request of a door of a value in conform | <https://github.com/dokimasia/dokimi-go/issues/22> |
| The request of the param of the directive `check` | <https://github.com/dokimasia/dokimi/issues/49> |
| The signer of a store in core | <https://github.com/thesm-os/core/issues/109> |
| The doors of conform | <https://github.com/dokimasia/dokimi-go/blob/main/conform/handle.go> |
| The harness of the Go CDK | <https://github.com/google/go-cloud/blob/v0.47.0/blob/drivertest/drivertest.go#L60> |
