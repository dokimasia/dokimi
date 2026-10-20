// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the injectionsafe mixin, which states
 * that a value carrying a control sequence is stored and returned as data
 * rather than interpreted.
 *
 * @remarks
 * The store-and-load pair is the claim's own shape: a hostile value goes in
 * under a key, and what the key answers afterwards is the data that went in
 * — never an interpretation of it. AUTO-INJECTION-SAFE probes exactly that
 * round trip with adversarial payloads on both slots.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/**
 * Contains the value that the directive names: a fragment that is syntax to
 * a SQL interpreter downstream, which a store must keep as the data it is.
 *
 * @remarks
 * Declared rather than derived. Which separator or quote is dangerous
 * depends on an interpreter only the author knows, and a drawn string is
 * "test-value" — nothing in it survives differently from anything else, so a
 * check over one passes a store that sanitises everything it is given.
 */
export const unsafe = "'; DROP TABLE users; --";

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=injectionsafetest/
export interface Mixed {
  /** Stores the value as data, under the key as data. */
  //dokimi:mixin injection-safe read=load unsafe=unsafe
  store(key: string, value: string): Promise<void>;

  /** Returns what store kept under the key, uninterpreted. */
  load(key: string): Promise<string>;
}
