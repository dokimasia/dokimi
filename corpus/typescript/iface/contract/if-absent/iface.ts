// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the if-absent contract: a write that
 * only lands when the key is absent.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * @packageDocumentation
 */

/**
 * Reports a duplicate write.
 *
 * @remarks
 * The error is the contract's own conflict sentinel, declared beside the
 * shape so the directive can name it.
 */
export const exists = new Error("ifabsent: exists");

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
//dokimi:out path=ifabsenttest/
export interface Contract {
  /**
   * Writes the value where its key is absent.
   *
   * @remarks
   * put is the if-absent contract's writer role, and hosts the directive that
   * names its partners. The conflict key is what turns "refused" into
   * "refused for being present": without it, a subject failing on a missing
   * map passes the duplicate check.
   */
  //dokimi:contract if-absent role=writer conflict=exists
  put(v: Value): Promise<void>;
}
