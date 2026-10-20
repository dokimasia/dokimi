// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the producer arm of the cursor
 * contract.
 *
 * @remarks
 * The directive sits on the method that returns a fresh cursor, and its
 * next-method= and close-method= name members of the handle it returns — the
 * produced type declares no directives of its own, so the producing method is
 * the only place its contract can be licensed. The standalone arm keeps its
 * own fixture in ../cursor.
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
export const closed = new Error("cursoropen: closed");

/** Represents the payload that the produced cursor drains. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/**
 * Represents the handle that open returns.
 *
 * @remarks
 * The members that the directive's next-method= and close-method= name are
 * declared here, on the cursor, because draining and releasing are the
 * handle's operations, not the producer's.
 */
export interface Cursor {
  /**
   * Returns the next value, or undefined when the cursor is exhausted.
   *
   * @remarks
   * After close, next reports {@link closed}.
   */
  next(): Promise<Value | undefined>;

  /** Releases the cursor, where a second close changes nothing. */
  close(): Promise<void>;
}

/**
 * Represents the subject of the fixture.
 *
 * @remarks
 * No //dokimi:model: a producer-only interface has no method that maps to a
 * sequence action — the produced cursor's laws instantiate at the CURSOR's
 * type and lift onto the producer.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=cursoropentest/
export interface Contract {
  /**
   * Returns a fresh cursor, in the open role of the cursor contract.
   *
   * @remarks
   * open hosts the directive that names the handle's members and the
   * sentinel they share.
   */
  //dokimi:contract cursor role=open next-method=next close-method=close sentinel=closed
  open(): Promise<Cursor>;
}
