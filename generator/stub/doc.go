// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package stub generates a typed double of each interface, struct or class
// under the directive stub. It also generates a test of each double.
//
// The plugin derives a model of each double from the facts of eidos and
// records it. The model contains the methods of the member walk of the
// declaration, the call and the results of each method, whether a fault can
// fail it, and the properties that the double declares as fields. The
// templates of each target write the double against the stub runtime of
// their language:
//
//	store := NewStoreStub(t)
//	store.OnGet.Returns(item, nil)
//
//	const store = new StoreStub(seat);
//	store.onGet.returns(item);
//
// The test of a double with a method has a level for each method and each
// option of the double, and for ResetCalls, NewStoreStubDelayed and
// NewStoreStubFailing. These levels run the behaviour suites of the runtime
// against the double. Their other cases check the results of each method with
// a sample of each argument and of each result. The test of the stub of a
// method has a level for each method of the stub, such as Returns and Func.
// The test of a double without a method builds the double under each option,
// and checks that a reset keeps the sample of each property.
//
// # Default results
//
// An unconfigured call returns a new value from the builder of each result
// whose type has one, and the zero value of each other result. The builder
// generator stamps the names of the builder of a type as facts. The records
// of a double and of its test refer to the builder through a bare reference,
// which the layout of eidos qualifies for the file of each record. The layout
// resolves a bare reference in the package of its unit alone, so a type of
// another package than the subject has no default result. A TypeScript
// module is a package of its own.
//
// # Languages
//
// The file golang.go contains the part of Go. The file typescript.go
// contains the part of TypeScript. Each part declares the target of its
// language with its template tree and its helpers. It reports whether the
// language types a value by its members, which subjects the file of a double
// can refer to and which interfaces only constrain type parameters. It sets
// the receiver of a method of a double, and returns the type that the option
// DelegateTo wraps. For a Go struct the part of Go returns no type, and the
// template writes the interface of the methods of the double. Each part
// returns the names of the member that configures a method, of the fields of
// a call and of the fields of the results. It returns the type that a builder
// builds for a result, and the declaration that a default result calls. stub
// leaves a declaration of a language without a part alone.
//
// # Diagnostics
//
// stub reports [Constraint] at an interface that only constrains type
// parameters, [IncompleteMembers] at a declaration with a gap in its member
// walk, [Unimplementable] at a declaration with a member that a double
// cannot implement or that the file of its double cannot refer to,
// [Collision] at a declaration whose double would have two members of one
// name, and [Empty] at a declaration without a member. stub does not record a
// declaration that it reports. stub also reports [Incomparable] at a
// declaration with a partition fault that cannot compare the arguments of
// its parameter, and records the double without the helpers of the
// partition.
//
// # Dependency position
//
// The package imports the annotator fault, the generator builder and the
// package internal/lang of the module, the root packages of the Go and the
// TypeScript satellites in golang.go and typescript.go, the package naming of
// eidos-lang, the shape catalog of eidos, the root, diag, directive, emit,
// meta, node, plugin, rules and symbol packages of the SDK facade of eidos,
// and the Go stdlib.
package stub
