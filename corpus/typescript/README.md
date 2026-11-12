<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# The TypeScript corpus

The TypeScript corpus is the TypeScript form of `corpus/go`. It has the same axes, the same fixture directories and the same directives. Each fixture declares its subject in `iface.ts`, and the generated files are beside it or in the directory that its `//dokimi:out` names.

## Run the corpus

Install the packages, generate the files and type-check them, then run the gate of the corpus:

```sh
pnpm --dir corpus/typescript install
dokimi -C corpus/typescript run
pnpm --dir corpus/typescript run types
dokimi conform --lang typescript
```

`dokimi conform --lang typescript` checks that the generated files are current, and runs the tests of the fixtures at the level of `overlay.yaml` through `pnpm exec vitest run --reporter=json`. It then judges each fixture against its expected outcomes under `corpus/expect`. `overlay.yaml` declares each difference of TypeScript from the corpus of Go. The gate fails on a difference without an entry.

The script `test` measures the coverage of the doubles, the builders and the surfaces of the enumerations, and fails below 100%. `vitest.config.ts` runs `vitest.setup.ts` before each test file, and the setup file registers the assertions of `@dokimi/assert` on the `expect` of vitest, which the generated tests call. The configuration in `tsconfig.json` turns on every strict check, `noUnusedLocals` and `noUnusedParameters` included, so the generated code compiles in the strictest project that consumes it.

## The translation of a Go fixture

| Go | TypeScript |
|---|---|
| A method with a `context.Context` | An asynchronous method whose result is a `Promise`. Its last parameter is `signal: AbortSignal`, and a rest parameter follows the signal |
| A method without a `context.Context` | A synchronous method |
| The results `(T, error)` | The result `T`. A failure rejects the promise or throws |
| The result `error` alone | The result `void` |
| The results `(A, B, error)` | The tuple `[A, B]` |
| The results `(T, bool)` and the result `*T` | The result `T \| undefined` |
| `[]T`, `map[K]V` and `map[K]struct{}` | `T[]`, `Map<K, V>` and `Set<K>` |
| `[]byte`, `time.Time` and `time.Duration` | `Uint8Array`, `Date` and a `number` of milliseconds |
| `iter.Seq2[T, error]` and `<-chan T` | `AsyncIterable<T>` |
| A struct | An interface, or a class under `//dokimi:builder` |
| An error variable such as `ErrNotFound` | A constant such as `notFound` |
| An error type | A class that extends `Error` |
| A defined type over a builtin | An alias |
| An exported name such as `Get` | The same name in camel case, such as `get`, in the declaration and in each directive param |

The fixtures document their declarations in TSDoc. The directives are line comments directly above the declaration that they address, after its documentation.

## Fixtures without a TypeScript form

- `enum/split` declares the members of an enumeration in another file than its type. A TypeScript enum declares its members in one declaration.
- `errors/split` spreads the errors of one package over three files. A TypeScript module is one file.
- `iface/lang/namedreturns` names the results of its methods. A TypeScript method returns one value without a name.
- `TagConfig` and `MixedConfig` of `struct/config` declare defaults in struct tags, which TypeScript does not have.
