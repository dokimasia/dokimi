# go.dokimi.dev/dokimi

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
