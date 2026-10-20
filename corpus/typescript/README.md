<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# The TypeScript corpus

The TypeScript corpus is the TypeScript form of `corpus/go`. It has the same axes, the same fixture directories and the same directives. Each fixture declares its subject in `iface.ts`, and the generated files are beside it or in the directory that its `//dokimi:out` names.

## Run the corpus

Generate the files, then type-check them and run their tests:

```sh
dokimi -C corpus/typescript run
npm --prefix corpus/typescript install
npm --prefix corpus/typescript run types
npm --prefix corpus/typescript test
```

The script `test` measures the coverage of the doubles, the builders and the surfaces of the enumerations, and fails below 100%. The configuration in `tsconfig.json` turns on every strict check, `noUnusedLocals` and `noUnusedParameters` included, so the generated code compiles in the strictest project that consumes it.

## The translation of a Go fixture

| Go | TypeScript |
|---|---|
| A method with a `context.Context` | An asynchronous method, whose result is a `Promise` |
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
