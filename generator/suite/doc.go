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
// A function that takes a *prop.Case after the handle is a property check.
// The runtime runs it once for each case of the property engine, and the fuzz
// target of its callable runs it under the fuzzer:
//
//	// KeepsTheLastPut checks that Get returns the session of the last Put of
//	// its ID.
//	//
//	//dokimi:check Sessions.Get defect=keepsTheFirstPut
//	func KeepsTheLastPut(h *SessionsHandle, c *prop.Case) { … }
//
// The param counted=true marks a check that counts the allocations of the
// whole process, such as with assert.MaxAllocs. Such a check runs with the
// allocation ceilings in TestStoreAllocs and TestStoreSuiteAllocs, which do
// not run in parallel:
//
//	// CreditsInPlace checks that Credit allocates nothing once the account
//	// has a balance.
//	//
//	//dokimi:check Ledger.Credit defect=copiesBalances counted=true
//	func CreditsInPlace(h *LedgerHandle) { … }
//
// # Subjects
//
// A consumer applies the suite to a subject, and adds checks of its own:
//
//	storetest.RunStoreSuite(t, memory.New, conform.With(keepsTheItem()))
//
// The directive subject on a struct of the consumer generates that test in
// the external test package of the struct. Its param new contains the name of
// the factory of the subjects, an exported function of the package of the
// struct or a function of its external test package:
//
//	//dokimi:subject store.Store new=New
//	type Store struct{ … }
//
// The test of a subject also declares a fuzz target of each callable, such as
// FuzzStoreGet, which calls the fuzz form FuzzStoreSuite of the suite. The
// test of the suite of a function applies the suite to the function itself,
// and declares its fuzz target, such as FuzzSlug.
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
// # Budgets
//
// The directive budget states the ceilings of a callable: the most heap
// allocations and bytes of one call, and the highest p99 and mean latency of
// one call. On a method or a function under the directive suite, the budget
// applies to every subject of the declaration:
//
//	// Get returns the item with the identity id.
//	//
//	//dokimi:budget allocs=1 bytes=96 p99=2us mean=1us
//	Get(ctx context.Context, id string) (Item, error)
//
// On a subject, the budget names its callable and applies to that subject
// alone:
//
//	//dokimi:subject store.Store new=New
//	//dokimi:budget Store.Get p99=3ms mean=1ms
//	type Store struct{ … }
//
// The runtime derives a check with a planted defect from each ceiling, such
// as Store.Get/allocation-budget. The names of the callables return the
// identity of each ceiling of the budget of a declaration, such as
// StoreChecks.Get.AllocationBudget(), and a subject excuses a ceiling of its
// declaration with the directive excuse. The param warmup of a budget states
// the number of iterations that each benchmark of the callable runs before it
// measures, as in //dokimi:budget allocs=0 warmup=20000. A budget may state a
// warm-up without a ceiling.
//
// # Calls and benchmarks
//
// The suite declares the call generator of each callable, which draws the
// inputs of each call through the property engine of the assert library. A
// context takes the context of the call, and a variadic parameter takes a
// list. A pointer takes the address of a value, or a pointer that can be nil
// under the mixin nil-safe. The members of each enumeration that the inputs
// take at some depth replace the generator of the enumeration. The suite also
// declares the planted defect of the ceilings: a subject whose callables spend
// the cost that exceeds a ceiling.
//
// The test of a subject declares a benchmark of each callable, such as the
// level Get of BenchmarkStore. It declares a parallel benchmark, such as
// BenchmarkStoreParallel, of each callable with the mixin concurrent or
// concurrent-readers, and a test of the allocation ceilings, such as
// TestStoreAllocs, where a budget states one or a check counts allocations.
// The test of a suite with an allocation ceiling or a check that counts
// allocations declares TestStoreSuiteAllocs, which falsifies both. The test of the suite of a function declares the same parts for
// the function, such as BenchmarkSlug and TestSlugAllocs.
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
// without a statement or with the wrong signature, of a check of a function
// with the segment of a check that the spec of the suite adds, and of a check
// with the segment of a ceiling. It reports [UnknownExcusal] at an excusal
// whose identity does not match a check of the suite, and [BadBudget] at a
// budget that no check of a suite applies.
//
// # Dependency position
//
// The package imports the packages annotator/defaults, annotator/fault,
// generator/builder, generator/stub, internal/lang and internal/weave of the
// module, the root package of the Go satellite in golang.go and call.go, the
// package naming of eidos-lang, the package shape of the shape plugin of
// eidos, the root, diag, directive, emit, meta, node, plugin, rules, store and
// symbol packages of the SDK facade of eidos, and the Go stdlib.
package suite
