// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the if-match contract: a write
 * conditional on a predicate.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * The predicate has two spellings and this fixture uses the callable one:
 * match plays the role match of the contract. `pred=` is the other form and
 * carries an expression — `pred=version==expected` — which the resolver
 * deliberately leaves verbatim, because there is no callable in it.
 *
 * A conformance check has to call the predicate, so only the callable form is
 * reachable from this tier. The expression form is a declaration the model
 * tier can act on and this one cannot.
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
//dokimi:out path=ifmatchtest/
export interface Contract {
  /** Writes the value, in the writer role of the if-match contract. */
  //dokimi:contract if-match role=writer field=body
  put(v: Value): Promise<void>;

  /**
   * Reports whether the write of the value may land.
   *
   * @remarks
   * match is the predicate the write is conditional on. It answers about the
   * value put takes, because a predicate over anything else is one the write
   * cannot be conditional on.
   */
  //dokimi:contract if-match role=match
  match(v: Value): Promise<boolean>;
}
