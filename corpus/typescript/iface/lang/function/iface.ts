// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for free functions.
 *
 * @remarks
 * eidos annotates module-level functions through the same three axes it
 * applies to interface methods — its walker has an OnFunction hook alongside
 * OnMethod. Every other fixture in this corpus declares an interface, so a
 * generator or a gate that reads only methods would look complete against all
 * of them and still miss every classification a consumer put on a plain
 * function.
 *
 * The functions here carry a body rather than being declarations, because a
 * free function has nowhere else to live: an interface method is a signature
 * and a function is not.
 *
 * Nothing is generated from them, and that is the gap this fixture records
 * rather than closes. Every generator of dokimi starts from a declaration
 * with members, so the directives below reach the annotation census and stop
 * there. The census is what counts them: a classification reachable only
 * outside an interface would read as uncovered without these, which is what
 * the module is for. No //dokimi:out appears here to promise otherwise.
 *
 * @packageDocumentation
 */

/** Reports a key that get misses. */
export const notFound = new Error("function: not found");

/** Represents the payload that these functions exchange. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/**
 * Returns the value under the key.
 *
 * @remarks
 * get is a reader by signature alone, classified with no directive — the
 * same inference that applies to a method of the same shape.
 *
 * @throws {@link notFound} for the empty key.
 */
export function get(key: string): Promise<Value> {
  if (key === "") {
    return Promise.reject(notFound);
  }
  return Promise.resolve({ key, body: "" });
}

/**
 * Stores the value.
 *
 * @remarks
 * put carries a mixin, proving the directive axes reach a function and not
 * only a method.
 */
//dokimi:mixin idempotent
export function put(_v: Value): Promise<void> {
  return Promise.resolve();
}

/**
 * Acquires the lease of the key, in the acquire role of the lease contract on
 * a free function.
 *
 * @remarks
 * A contract binds the functions of a module as it binds the methods of an
 * interface.
 */
//dokimi:contract lease role=acquire
export function acquire(_key: string): Promise<void> {
  return Promise.resolve();
}

/** Releases the lease of the key, in the release role of the lease contract. */
//dokimi:contract lease role=release
export function release(_key: string): Promise<void> {
  return Promise.resolve();
}
