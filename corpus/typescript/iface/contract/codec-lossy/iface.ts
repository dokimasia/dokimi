// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the lossy form of the codec
 * contract: a forward transform that discards information, paired with an
 * inverse that recovers only what survived.
 *
 * @remarks
 * `fidelity=lossy` selects the weaker roundtrip law: not the identity, but
 * agreement on the second pass — encode, decode, encode again, and the two
 * encodings match because everything the transform was going to lose is
 * already gone. The exact form lives in the sibling `codec` fixture; both
 * exist because the parameter selects between two laws, and a corpus stating
 * only one leaves the other's selection untested.
 *
 * The model tier rides the twin floor and that is not a floor here, for the
 * reason the lossless codec beside it gives: a codec holds no state for an
 * oracle to model, and the claim is about the round trip.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=codeclossytest/
export interface Contract {
  /**
   * Returns the encoding of the input, in the forward role of the codec
   * contract.
   *
   * @remarks
   * The directive sets the fidelity of the round trip.
   */
  //dokimi:contract codec role=forward fidelity=lossy
  encode(input: string): Promise<string>;

  /** Returns the decoding of the input, in the inverse role of the codec contract. */
  //dokimi:contract codec role=inverse
  decode(input: string): Promise<string>;
}
