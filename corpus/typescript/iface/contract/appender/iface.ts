// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the appender contract: an
 * append-only log whose writes only ever add.
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
//dokimi:out path=appendertest/
export interface Contract {
  /**
   * Appends the value, and returns its offset in the log.
   *
   * @remarks
   * run is the appender contract's fn role, and hosts the directive that
   * names its partners. The offset it answers is the contract's observable:
   * AUTO-APPENDER-MONOTONIC-OFFSETS holds successive answers strictly
   * increasing.
   */
  //dokimi:contract appender role=fn
  run(v: Value): Promise<number>;
}
