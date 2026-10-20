// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture that arms the redeliver role of the
 * publisher: the at-least-once bound proven through an actual redelivery
 * rather than on the single publish.
 *
 * @remarks
 * republish re-offers a message the broker already delivered, so each
 * subscriber sees the duplicate the mode permits — which is exactly what
 * tells at-least-once apart from exactly-once. Its unarmed siblings exercise
 * the role's optional omission, and their headers say so.
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
//dokimi:out path=publisherredelivertest/
export interface Contract {
  /** Publishes the value, in the publish role of the publisher contract. */
  //dokimi:contract publisher role=publish mode=at-least-once
  publish(v: Value): Promise<void>;

  /**
   * Re-offers a message to every current subscriber, in the redeliver role of
   * the publisher contract.
   *
   * @remarks
   * The duplicate is what this mode permits.
   */
  //dokimi:contract publisher role=redeliver
  republish(v: Value): Promise<void>;

  /** Returns the published values, in the subscribe role of the publisher contract. */
  //dokimi:contract publisher role=subscribe
  subscribe(): AsyncIterable<Value>;
}
