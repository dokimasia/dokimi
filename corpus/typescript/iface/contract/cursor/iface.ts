// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the cursor contract: a cursor that
 * next drains and close releases.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * @packageDocumentation
 */

/**
 * Reports a call of next after close.
 *
 * @remarks
 * The error is the contract's own sentinel, declared beside the shape so the
 * directive can name it.
 */
export const closed = new Error("cursor: closed");

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
//dokimi:out path=cursortest/
export interface Contract {
  /** Returns the next value, or undefined once the cursor is exhausted, in the next role of the cursor contract. */
  //dokimi:contract cursor role=next sentinel=closed
  next(): Promise<Value | undefined>;

  /** Releases the cursor, in the close role of the cursor contract. */
  //dokimi:contract cursor role=close
  close(): Promise<void>;
}
