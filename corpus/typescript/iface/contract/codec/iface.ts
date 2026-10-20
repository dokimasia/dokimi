// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the codec contract: a forward
 * transform and the inverse that undoes it.
 *
 * @remarks
 * `fidelity` says which equality the pair claims. This fixture declares the
 * default `exact` explicitly rather than omitting it, so the fixture states
 * what it means instead of relying on a reader knowing the default.
 *
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * The model tier rides the twin floor and that is not a floor here. A codec
 * holds nothing, so there is no state for an oracle to model; the claim is
 * the round trip, which AUTO-CODEC-ROUNDTRIP states, and two instances
 * agreeing is exactly what determinism means for a pure pair.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=codectest/
export interface Contract {
  /**
   * Returns the encoding of the input, in the forward role of the codec
   * contract.
   *
   * @remarks
   * The directive sets the fidelity of the round trip.
   */
  //dokimi:contract codec role=forward fidelity=exact
  encode(input: string): Promise<string>;

  /** Returns the decoding of the input, in the inverse role of the codec contract. */
  //dokimi:contract codec role=inverse
  decode(input: string): Promise<string>;
}
