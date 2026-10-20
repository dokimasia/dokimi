// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the defaultonerror mixin, which states
 * that a failed read hands back the type's declared default rather than an
 * undefined value.
 *
 * @remarks
 * The result is a named type carrying //dokimi:default, because that is
 * where the vocabulary looks for "the type's default": a claim about the
 * error path's expected answer needs the declaration to state that answer,
 * or there is nothing to hold the failed read against. The mixin REPLACES
 * the zero-on-error claim the reader shape would otherwise license — the two
 * contradict on one method, and the authored statement wins.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/**
 * Represents what the store returns for a key.
 *
 * @remarks
 * The directive declares the resting value a failed read owes — the mixin's
 * whole claim — and {@link defaultBody} is the same value as a constant, so
 * the subject and the seats never respell it.
 */
//dokimi:default "fallback-body"
export type Body = string;

/**
 * Contains the declared default as a value.
 *
 * @remarks
 * The directive on {@link Body} declares the same spelling for the
 * derivation. Code uses this constant.
 */
export const defaultBody: Body = "fallback-body";

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=defaultonerrortest/
export interface Mixed {
  /** Stores the body under the key, as the write the read is judged against. */
  store(key: string, body: Body): Promise<void>;

  /**
   * Returns the body under the key.
   *
   * @remarks
   * A failed read returns the declared default rather than whatever the
   * lookup left behind.
   */
  //dokimi:mixin default-on-error
  get(key: string): Promise<Body>;
}
