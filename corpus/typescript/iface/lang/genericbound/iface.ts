// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for a type parameter whose constraint
 * the generator cannot reason about.
 *
 * @remarks
 * An unconstrained parameter is the one bound whose type set is known without
 * reading anything: it admits every type. A named constraint is a reference
 * the derivation does not evaluate, so nothing tells it which types satisfy
 * the bound — and inventing one produces a companion that fails to compile
 * for a reason the author could not have predicted.
 *
 * The witness directive is the answer: the source names the type each
 * parameter is instantiated at, the witness annotator resolves and stamps
 * them, and every generator reads the stamp through its rules. This fixture
 * holds that path, as the generic fixture holds the derived one.
 *
 * @packageDocumentation
 */

/**
 * Represents the constraint of the keys.
 *
 * @remarks
 * The constraint is declared here rather than written inline, which is
 * exactly what makes it opaque to the generator.
 */
export type Ordered = number | string;

/**
 * Represents the value type.
 *
 * @remarks
 * A witness naming a type of this module exercises the import a bare name
 * needs, because the companion is in another directory.
 */
export interface Score {
  /** Counts the points of the score. */
  points: number;
}

/**
 * Represents an interface that is generic over an opaque constraint.
 *
 * @remarks
 * Without the witness its companion would be a note, because no derivation
 * can tell that `number` satisfies {@link Ordered}.
 */
//dokimi:stub
//dokimi:suite
//dokimi:witness K=number V=Score
//dokimi:out path=genericboundtest/
export interface Ranked<K extends Ordered, V> {
  /** Returns the value of the rank of the key. */
  rank(key: K): Promise<V>;
  /** Clears the ranks. */
  reset(): Promise<void>;
}
