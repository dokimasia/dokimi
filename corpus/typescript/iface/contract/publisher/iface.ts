// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the publisher contract: a publish
 * paired with the subscribe that delivers it.
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
//dokimi:out path=publishertest/
export interface Contract {
  /** Publishes the value, in the publish role of the publisher contract. */
  //dokimi:contract publisher role=publish
  publish(v: Value): Promise<void>;

  /** Returns the published values, in the subscribe role of the publisher contract. */
  //dokimi:contract publisher role=subscribe
  subscribe(): AsyncIterable<Value>;
}
