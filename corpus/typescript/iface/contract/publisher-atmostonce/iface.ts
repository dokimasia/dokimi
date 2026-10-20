// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the at-most-once bound of the
 * publisher contract: each subscriber receives a published message zero or
 * one times — loss permitted, duplicates forbidden.
 *
 * @remarks
 * The mode rides the directive, and the bound law it selects counts each
 * subscriber's copies of a published message against it. The redeliver role
 * stays undeclared here, so the bound is exercised on the single publish and
 * the header says the arm went unarmed — `publisher-redeliver` and the
 * exactly-once sibling are the fixtures that declare the role.
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
//dokimi:out path=publisheratmostoncetest/
export interface Contract {
  /** Publishes the value, in the publish role of the publisher contract. */
  //dokimi:contract publisher role=publish mode=at-most-once
  publish(v: Value): Promise<void>;

  /** Returns the published values, in the subscribe role of the publisher contract. */
  //dokimi:contract publisher role=subscribe
  subscribe(): AsyncIterable<Value>;
}
