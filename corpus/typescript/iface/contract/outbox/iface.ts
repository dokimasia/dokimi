// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the outbox contract: a
 * transactional outbox and its reader.
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
//dokimi:out path=outboxtest/
export interface Contract {
  /** Appends the value to the outbox, in the append role of the outbox contract. */
  //dokimi:contract outbox role=append
  append(v: Value): Promise<void>;

  /** Returns the appended values, in the subscribe role of the outbox contract. */
  //dokimi:contract outbox role=subscribe
  subscribe(): AsyncIterable<Value>;
}
