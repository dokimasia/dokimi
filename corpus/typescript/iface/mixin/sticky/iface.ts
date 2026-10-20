// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the sticky mixin, which states that
 * once a key resolves, it keeps resolving to the same value.
 *
 * @remarks
 * The interface carries a write beside the read because the claim relates
 * them: a read alone has nothing to be judged against.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

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
//dokimi:out path=stickytest/
export interface Mixed {
  /** Stores the value, as the write the read is judged against. */
  store(v: Value): Promise<void>;

  /** Returns the value under the key, which keeps the first value it resolved to. */
  //dokimi:mixin sticky key=key
  get(key: string): Promise<Value>;
}
