// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the watcher contract: a watch woken
 * by a trigger.
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
 * Represents the handle that watch returns.
 *
 * @remarks
 * The members that the directive's next= and stop= name are declared here,
 * on the subscription, because reading the next event and ending the watch
 * are the handle's operations, not the store's.
 */
export interface Subscription {
  /**
   * Returns the next change within the wait, or undefined where none arrived.
   *
   * @param timeout - The wait, in milliseconds.
   */
  next(timeout: number): Value | undefined;

  /** Ends the subscription. */
  stop(): void;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=watchertest/
export interface Contract {
  /**
   * Watches the key, in the watch role of the watcher contract.
   *
   * @remarks
   * The directive names the two members of the subscription that watch
   * returns.
   */
  //dokimi:contract watcher role=watch next=next stop=stop
  watch(key: string): Promise<Subscription>;

  /**
   * Publishes a change of the key, in the trigger role of the watcher
   * contract.
   *
   * @remarks
   * trigger carries the value the change publishes, because a watch that
   * fires is only half the claim — what arrived has to be what was written.
   */
  //dokimi:contract watcher role=trigger
  trigger(key: string, v: Value): Promise<void>;
}
