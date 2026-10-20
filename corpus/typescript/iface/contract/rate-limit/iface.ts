// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the rate-limit contract: a call
 * bounded by a rate and a burst.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * @packageDocumentation
 */

/** Represents the payload that the roles of the contract exchange. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/**
 * Reports a call refused because the burst is spent.
 *
 * @remarks
 * In the module of the interface rather than of the subject, because
 * limited= names it and the directive resolves through the scope of this
 * module. A sentinel a consumer cannot reach is one their implementation
 * cannot report, and the claim is about identity: any error would let a
 * subject that fails for its own reasons pass as a limiter.
 */
export const limited = new Error("ratelimit: burst spent");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=ratelimittest/
export interface Contract {
  /**
   * Runs the limited call for the key.
   *
   * @remarks
   * run is the rate-limit contract's fn role, and hosts the directive that
   * names its partners.
   */
  //dokimi:contract rate-limit role=fn rate=100 burst=10 limited=limited
  run(key: string): Promise<void>;
}
