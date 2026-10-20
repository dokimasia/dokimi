// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the cas contract: a compare-and-swap
 * write guarded by a version.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * @packageDocumentation
 */

/**
 * Reports a write with a stale version.
 *
 * @remarks
 * The error is the contract's own sentinel, declared beside the shape so the
 * directive can name it.
 */
export const mismatch = new Error("cas: the write's version is stale");

/** Represents the payload that the roles of the contract exchange. */
export interface Value {
  /** Contains the body of the cell. */
  body: string;
  /** Contains the version that the compare-and-set guards the body by. */
  version: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=castest/
export interface Contract {
  /**
   * Writes the value where its version is current.
   *
   * @remarks
   * put is the cas contract's writer role, and hosts the directive that names
   * its partners.
   */
  //dokimi:contract cas role=writer version=version mismatch=mismatch
  put(v: Value): Promise<void>;

  /** Returns what the cell contains, which a compare-and-set is meaningless without. */
  get(): Promise<Value>;
}
