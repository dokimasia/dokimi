// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for AUTHORED VALUES: a //dokimi:sample
 * names the function that a generated check writes a type with.
 *
 * @remarks
 * The derived value would break the run here — a tag the subject validates
 * against a vocabulary its structure does not state.
 *
 * The directive is load-bearing here, not decorative. apply refuses a tag
 * outside the vocabulary, a refused first call leaves the idempotent claim
 * vacuous, and the census fails a run whose claim never engaged — so without
 * the authored sample the generated suite cannot be green against a
 * validating subject.
 *
 * Level carries both value directives at once, so the generated fixture also
 * pins their precedence: the sample outranks the default, because the
 * fixture is what checks write and that is the sample directive's half of the
 * composition — defaults seed a constructor.
 *
 * @packageDocumentation
 */

/**
 * Reports a tag outside the vocabulary of the subject.
 *
 * @remarks
 * The refusal makes a derived "test-apply" value unable to engage the
 * idempotent claim.
 */
export const unknownTag = new Error("sampled: unknown tag");

/**
 * Represents a word of the validated vocabulary.
 *
 * @remarks
 * A derived sample satisfies the compiler and fails the vocabulary, which is
 * the case the sample directive exists for.
 */
//dokimi:sample value=newTestTag()
export type Tag = string;

/** Returns a tag that the vocabulary admits, as the function that the sample directive names. */
export function newTestTag(): Tag {
  return "release";
}

/**
 * Represents a level with both directives of different values.
 *
 * @remarks
 * The generated fixture proves which one wins.
 */
//dokimi:sample value=newTestLevel()
//dokimi:default 3
export type Level = number;

/** Returns the authored level, which the fixture must draw over the declared default. */
export function newTestLevel(): Level {
  return 7;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=sampledtest/
export interface Labeller {
  /**
   * Applies the tag at the level.
   *
   * @remarks
   * apply is idempotent for a known tag and refuses an unknown one with
   * {@link unknownTag}. The refusal is the point: a refused first call leaves
   * the idempotent claim vacuous, so only an authored value the vocabulary
   * admits lets the claim engage.
   */
  //dokimi:mixin idempotent
  apply(tag: Tag, level: Level): Promise<void>;
}
