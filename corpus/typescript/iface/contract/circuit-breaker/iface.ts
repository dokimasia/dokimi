// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the circuit-breaker contract: a call
 * that trips open after repeated failure.
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

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=circuitbreakertest/
export interface Contract {
  /**
   * Runs the guarded call for the key.
   *
   * @remarks
   * run is the circuit-breaker contract's fn role, and hosts the directive
   * that names its partners.
   */
  //dokimi:contract circuit-breaker role=fn
  run(key: string): Promise<void>;
}
