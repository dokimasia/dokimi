// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the pagination contract: a reader
 * whose protocol cursors through pages.
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
 * Represents where a walk resumes.
 *
 * @remarks
 * The empty cursor is both where a walk starts and where it ends.
 */
export type Cursor = string;

/**
 * Represents what one call returns: the entries of the page and the cursor
 * that a caller passes back to go on.
 *
 * @remarks
 * A record rather than a spread of results, because the cursor is the member
 * the directive names and a member has to live on something.
 */
export interface Page {
  /** Lists the entries of the page. */
  items: Value[];
  /** Contains the cursor where the walk resumes. */
  next: Cursor;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=paginationtest/
export interface Contract {
  /**
   * Returns the page at the cursor, in the reader role of the pagination
   * contract.
   *
   * @remarks
   * fetch hosts the directive that names its partners. It returns one page
   * and where the walk resumes — the shape both walk claims need, because a
   * keyed read has no cursor to resume from.
   */
  //dokimi:contract pagination role=reader cursor=next
  fetch(cur: Cursor): Promise<Page>;

  /**
   * Stores an entry that the walk visits.
   *
   * @remarks
   * Without a writer the sequences seed nothing and every walk holds over
   * empty pages.
   */
  put(v: Value): Promise<void>;
}
