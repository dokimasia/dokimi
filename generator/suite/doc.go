// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package suite generates the conformance suite of each declaration under the
// directive suite, and the conformance test of each subject under the
// directive subject.
//
// The suite of a declaration contains the names of its callables with the
// identities of their checks, the function that returns its checks, and the
// two forms of the application that runs the checks on a subject. The spec of
// the file of the suite falsifies each check on its planted defect. It also
// checks that each form passes the value of its factory to a check as the
// subject. The module go.dokimi.dev/conform declares the types of these parts:
//
//	//dokimi:suite
//	//dokimi:out path=storetest/
//	type Store interface {
//		Get(ctx context.Context, id string) (Item, error)
//	}
//
// # Checks of a consumer
//
// The directive check binds a function of the package of the generated suite
// to a callable of the suite as a check. Its param defect contains the name of
// the function that returns the planted defect of the check. The doc comment
// of the function opens with the name of the function and the words "checks
// that", and the rest of its first sentence is the statement of the check:
//
//	// KeepsTheItem checks that Get returns the item that Put stored.
//	//
//	//dokimi:check Store.Get defect=forgetful
//	func KeepsTheItem(h *StoreHandle) { … }
//
// Where two declarations under the directive suite have the name Store, the
// directive writes the name of the package before it, as in store.Store.Get.
// The annotate half of suite stamps the facts of each check that it binds, and
// the generate half reads them into the suite of the declaration.
//
// # Subjects
//
// A consumer applies the suite to a subject, and adds checks of its own:
//
//	storetest.RunStoreSuite(t, memory.New, conform.With(keepsTheItem()))
//
// The directive subject on a struct of the consumer generates that test in
// the external test package of the struct. Where the factory is not exported,
// the test is inside the package of the struct. Its param new contains the
// name of the factory of the subjects:
//
//	//dokimi:subject store.Store new=New
//	type Store struct{ … }
//
// # Excusals
//
// The directive excuse on the struct of a subject excuses a check of the
// suite from the conformance test of the subject. Its positional param is the
// identity of the check, and its param reason states why the subject does not
// pass the check. The test passes each excusal to the application through
// conform.Excuse, and the application skips the check with the reason:
//
//	//dokimi:subject store.Store new=New
//	//dokimi:excuse Store.Get/keeps-the-item reason="the store drops each item"
//	type Store struct{ … }
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
// cannot call. It reports [UnknownCallable] at the directive of a check whose
// callable it cannot find. It reports [BadCheck] at the directive of a check
// without a statement or with the wrong signature, and of a check of a
// function with the segment of a check that the spec of the suite adds. It
// reports [UnknownExcusal] at an excusal whose identity does not match a check
// of the suite.
//
// # Dependency position
//
// The package imports the packages internal/lang and generator/stub of the
// module, the root package of the Go satellite in golang.go, the package
// naming of eidos-lang, the root, diag, directive, emit, meta, node, plugin,
// rules and symbol packages of the SDK facade of eidos, and the Go stdlib.
package suite
