// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the exactly-once bound of the
 * publisher contract: each subscriber receives a published message exactly
 * once — neither loss nor duplication.
 *
 * @remarks
 * The mode rides the directive, and the bound law it selects counts each
 * subscriber's copies of a published message against it. The redeliver role
 * is armed here through replay, whose dedupe is the exactly-once claim
 * itself: the law re-offers the published message and the duplicate must
 * never reach a subscriber — which is what tells this mode apart from the
 * at-least-once fixture that redelivers without asking.
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
//dokimi:out path=publisherexactlyoncetest/
export interface Contract {
  /** Publishes the value, in the publish role of the publisher contract. */
  //dokimi:contract publisher role=publish mode=exactly-once
  publish(v: Value): Promise<void>;

  /**
   * Re-offers a message that the broker may already have delivered, in the
   * redeliver role of the publisher contract.
   *
   * @remarks
   * The exactly-once claim is that the duplicate never reaches a subscriber.
   */
  //dokimi:contract publisher role=redeliver
  replay(v: Value): Promise<void>;

  /** Returns the published values, in the subscribe role of the publisher contract. */
  //dokimi:contract publisher role=subscribe
  subscribe(): AsyncIterable<Value>;
}
