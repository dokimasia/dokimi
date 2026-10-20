// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the updater contract: an update
 * paired with the read that observes it.
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
//dokimi:out path=updatertest/
export interface Contract {
  /** Updates the value, in the writer role of the updater contract. */
  //dokimi:contract updater role=writer
  put(v: Value): Promise<void>;

  /** Returns the value under the key, in the reader role of the updater contract. */
  //dokimi:contract updater role=reader
  get(key: string): Promise<Value>;
}
