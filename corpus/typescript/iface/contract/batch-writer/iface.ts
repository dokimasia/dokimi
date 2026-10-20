// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the batch-writer contract: a write
 * that takes many values at once.
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
  /**
   * Contains the key that files the value.
   *
   * @remarks
   * This subject refuses an empty key, which is what gives `mode=atomic` a
   * failure to be about.
   *
   * Roled, and that is the point: a role opens a pool, and a pool is how a
   * consumer says which inputs their implementation accepts. The adversarial
   * arm draws the empty string among others, and a subject that narrows what
   * it takes has to say so or be red against inputs its own author ruled out.
   */
  //dokimi:default "test-key"
  key: string;

  /** Contains the body of the value. */
  //dokimi:default "test-body"
  body: string;
}

/**
 * Contains the value that this writer is declared to turn down: a body with
 * nothing to file it under.
 *
 * @remarks
 * Declared rather than drawn, because a derived value is built to be
 * ACCEPTED — it is what every other claim writes — so a check using one never
 * reaches the failure `mode=atomic` is about. Only the author knows what
 * their writer refuses, which is what refused= is for.
 */
export const unkeyed: Value = { key: "", body: "test-body" };

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=batchwritertest/
export interface Contract {
  /**
   * Writes the value, in the writer role of the batch-writer contract.
   *
   * @remarks
   * `mode=atomic` is a claim about state AFTER a failure, so it needs
   * something to observe that state with — which is what the reader role is
   * for. Without it the only statable claim is that a good write succeeds,
   * which holds for a store that also keeps half a refused one.
   */
  //dokimi:contract batch-writer role=writer mode=atomic refused=unkeyed
  put(v: Value): Promise<void>;

  /**
   * Returns the value under the key, in the reader role of the batch-writer
   * contract.
   *
   * @remarks
   * The reader is the observation the mode's claim is read through. A role
   * rather than a param, so the pairing is recorded on both members — a
   * consumer holding this method can find the contract it confirms, which a
   * one-directional `read=` could not answer.
   */
  //dokimi:contract batch-writer role=reader
  get(key: string): Promise<Value>;
}
