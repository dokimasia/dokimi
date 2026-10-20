// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for type parameters on the interface
 * itself, so every emitted type has to carry them through.
 *
 * @remarks
 * This axis varies the type system rather than the classification, so these
 * break generators independently of any directive.
 *
 * @packageDocumentation
 */

/**
 * Represents a store that is generic over key and value.
 *
 * @remarks
 * A generator that drops the type parameters produces code that does not
 * compile, which is the failure this fixture exists to catch.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=generictest/
export interface Store<K, V> {
  /** Returns the value under the key. */
  get(key: K): Promise<V>;
  /** Stores the value under the key. */
  put(key: K, value: V): Promise<void>;
}
