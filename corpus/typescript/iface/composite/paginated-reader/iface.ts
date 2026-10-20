// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the fixture that stacks the reader detector with the pagination
 * contract on the same method.
 *
 * @remarks
 * It earns a composite fixture because the contract changes what the
 * detector generates rather than adding to it. A bare reader's suite asserts
 * one call and one result. Once the method is pagination's reader role the
 * same signature has to be driven as a walk — call, take the cursor, call
 * again — and the assertion becomes that the second page is not the first,
 * which does not exist for an unpaginated reader.
 *
 * A generator that reads the axes independently emits a single-call reader
 * subtest and a pagination subtest that never runs, which passes.
 *
 * @packageDocumentation
 */

/** Represents the payload that the pages carry. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/**
 * Represents what one call returns.
 *
 * @remarks
 * next is a number offset rather than an opaque token, so the derivation is
 * held to the member's type rather than to a shape one fixture happens to
 * have.
 */
export interface Page {
  /** Lists the entries of the page. */
  items: Value[];
  /** Contains the offset where the walk resumes. */
  next: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=paginatedreadertest/
export interface PaginatedReader {
  /**
   * Returns the page at the cursor.
   *
   * @remarks
   * fetch is a reader by signature and pagination's reader role by directive.
   * The cursor parameter is what the contract keys on.
   */
  //dokimi:contract pagination role=reader cursor=next
  fetch(cursor: number): Promise<Page>;

  /**
   * Stores an entry that the walk visits, so a page has something to follow
   * it.
   *
   * @remarks
   * An operation nobody can call through the interface is one no check can
   * reach.
   */
  add(v: Value): Promise<void>;
}
