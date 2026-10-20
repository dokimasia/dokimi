// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the validates mixin, which states that
 * invalid input is rejected rather than stored.
 *
 * @remarks
 * The interface carries whatever methods the law needs to be stateable. A
 * mixin whose law spans two calls cannot be hosted by a single method: there
 * would be nothing to compare against, and the generated subtest would pass
 * by having nothing to check.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents what validate accepts or refuses. */
export interface Payload {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/**
 * Contains the value that the directive names: a payload without a key,
 * which validate refuses.
 *
 * @remarks
 * Declared rather than derived. A drawn Payload carries a key and a body and
 * is perfectly valid, so a check over one passes whatever the store did — it
 * stays green against a store that screens nothing at all.
 */
export const invalid: Payload = { key: "", body: "unkeyed" };

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=validatestest/
export interface Mixed {
  /**
   * Stores the payload, and rejects what validate refuses.
   *
   * @remarks
   * The fn parameter resolves to a sibling, so validate has to exist here for
   * the reference to bind.
   */
  //dokimi:mixin validates fn=validate invalid=invalid
  store(v: Payload): Promise<void>;

  /** Throws for a payload that it refuses, as the predicate that fn names. */
  validate(v: Payload): void;

  /** Returns the payload under the key, which proves that a rejected value was not stored. */
  read(key: string): Promise<Payload>;
}
