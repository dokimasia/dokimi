// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the ttl mixin, which states that stored
 * data stops being readable once its lifetime elapses.
 *
 * @remarks
 * The directive names everything the claim needs: the lifetime, the pair of
 * callables it governs, and the sentinel a lapsed read reports. It is the one
 * classification using all three reference kinds at once — a verbatim value,
 * a pair of sibling callables, and a module-level constant.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/**
 * Reports a read past the lifetime.
 *
 * @remarks
 * Module-level because the directive names it: a sentinel is a constant
 * rather than a callable, which is why it resolves through its own scope.
 */
export const expired = new Error("ttl: entry expired");

/** Represents the payload that the store keeps. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=ttltest/
export interface Mixed {
  /** Stores a value, which starts its lifetime. */
  put(v: Value): Promise<void>;

  /** Returns the value under the key while the lifetime lasts, and rejects with the sentinel after. */
  //dokimi:mixin ttl duration=1m put=put read=read not-found=expired
  read(key: string): Promise<Value>;
}
