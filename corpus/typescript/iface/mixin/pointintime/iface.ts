// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the pointintime mixin, which states
 * that two reads of one key agree even when a write lands between them.
 *
 * @remarks
 * The interface carries a write and a read because the claim relates them:
 * a read alone has nothing to be consistent with, and the generated subtest
 * would pass by having nothing to check.
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
//dokimi:out path=pointintimetest/
export interface Mixed {
  /** Stores the value, as the write the read is judged against. */
  store(v: Value): Promise<void>;

  /**
   * Returns the value under the key.
   *
   * @remarks
   * Stronger than cacheable, which permits the second read to observe a
   * concurrent write.
   *
   * write= names the store the two reads straddle. Without it the directive
   * states the property and no derivation can find the partner, because which
   * write matters is the author's choice rather than a fact of the signature.
   */
  //dokimi:mixin point-in-time write=store
  get(key: string): Promise<Value>;
}
