// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the persister contract: a write
 * paired with the read that observes it.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * @packageDocumentation
 */

/**
 * Reports a store whose medium has failed.
 *
 * @remarks
 * A check can ask for that state by identity, and whichever harness knows how
 * to produce it answers.
 */
export const mediumGone = new Error("persister: the medium is gone");

/**
 * Represents the payload that the roles of the contract exchange.
 *
 * @remarks
 * Annotated for the builder as well, which nothing else in the corpus does:
 * the builder axis and the model axis had never met, so a class that is both
 * a value type and a buildable one was a combination nothing exercised. They
 * compose — the builder, the double, the harness and the model tier share
 * this type without colliding over a name.
 *
 * The routing is repeated here because a declaration's //dokimi:out scopes to
 * that declaration, and the interface below carries its own. Without this
 * line the builder still emits — beside the source, which is where an output
 * with no routing goes, and not where the generated module is. Declaring it
 * once in the header of the module covers everything beneath it and is the
 * other way round; the validated integration fixture does that.
 */
//dokimi:out path=persistertest/
//dokimi:builder
export class Value {
  /**
   * Identifies the record.
   *
   * @remarks
   * Roled, which is what opens a pool: a role says what a value IS, and the
   * pool it opens carries a third member hostile to that role.
   */
  //dokimi:default "test-key"
  key: string = "";

  /** Contains the body of the record. */
  //dokimi:default "test-body"
  body: string = "";
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=persistertest/
export interface Contract {
  /**
   * Files the value, and returns the record as filed, in the writer role of
   * the persister contract.
   *
   * @remarks
   * The record AS FILED is what separates a persister from any other keyed
   * write. The retrievability claim reads that answer's key back, and a store
   * that files under a key of its own choosing is found under that one rather
   * than the one the caller supplied. A writer answering nothing leaves the
   * claim nothing to read back but the key the caller already had, which is
   * the write-observable row beside it.
   */
  //dokimi:contract persister role=writer
  put(v: Value): Promise<Value>;

  /** Returns the record under the key, in the reader role of the persister contract. */
  //dokimi:contract persister role=reader
  get(key: string): Promise<Value>;
}
