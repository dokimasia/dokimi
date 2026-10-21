// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package suite generates the conformance suite of each declaration under the
// directive suite, and the conformance test of each subject under the
// directive subject.
//
// The suite of a declaration contains the names of its callables with the
// identities of their checks, the function that returns its checks, and the
// applications that run the checks on a subject and on a control. The runtime
// conform of dokimi-go declares the types of these parts:
//
//	//dokimi:suite
//	//dokimi:out path=storetest/
//	type Store interface {
//		Get(ctx context.Context, id string) (Item, error)
//	}
//
// A consumer applies the suite to a subject, and adds checks of its own:
//
//	storetest.RunStore(t, memory.New, conform.With(keepsTheItem()))
//
// The directive subject on a struct of the consumer generates that test in
// the package of the struct. Its param new names the factory of the
// subjects:
//
//	//dokimi:subject store.Store new=New
//	type Store struct{ … }
//
// The suite of a declaration has no checks of its own, so an application runs
// the checks that the consumer adds.
//
// # Kinds of declaration
//
// The subjects of the suite of an interface have the type of the interface,
// at the witnesses of its type parameters. The suite of a struct declares the
// interface of the exported methods of the struct, such as StoreMethods, and
// its subjects have that type. The suite of a function declares the type of
// the function, such as SlugFunc, and its applications take the function
// itself.
//
// # Languages
//
// The file golang.go contains the part of Go, which declares the target of Go
// with its template tree. suite leaves a declaration of a language without a
// part alone.
//
// # Diagnostics
//
// suite reports [Uninstantiable] at a declaration whose type parameters its
// suite cannot instantiate, [NoSuite] at a subject whose declaration has no
// directive suite, and [BadFactory] at a subject whose factory an application
// cannot call.
//
// # Dependency position
//
// The package imports the packages internal/lang and generator/stub of the
// module, the root package of the Go satellite in golang.go, the package
// naming of eidos-lang, the root, diag, directive, emit, node, plugin, rules
// and symbol packages of the SDK facade of eidos, and the Go stdlib.
package suite
