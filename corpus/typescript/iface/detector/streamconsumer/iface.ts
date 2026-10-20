// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the streamconsumer shape: a stream
 * in rather than out.
 *
 * @remarks
 * It is not quite the mirror of a stream reader, and the asymmetry is worth
 * knowing: eidos recognises a produced stream as an async iterable return,
 * but a consumed stream as an *interface* parameter. A method taking an async
 * iterable is a method taking a library type, and no amount of naming makes
 * the detector see a stream in it.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the element type that the fixture ingests. */
export interface Value {
  /** Contains the key of the element. */
  key: string;
  /** Contains the body of the element. */
  body: string;
}

/**
 * Represents the stream that {@link StreamConsumer.ingest} consumes.
 *
 * @remarks
 * It has to be an interface: the detector tests whether the parameter is one.
 *
 * Both fixture values are declared rather than derived, because an interface
 * has no literal form and neither half could be worked out. Without the
 * sample directive every check needing a Source is declined, which is what
 * the named pair below is here to close.
 */
//dokimi:stub
//dokimi:out path=streamconsumertest/
export interface Source {
  /** Returns the next value, or undefined once the source is exhausted. */
  next(): Promise<Value | undefined>;
}

/**
 * Represents the subject of the fixture.
 *
 * @remarks
 * The model tier is unarmed and the reason is on the one method: a stream is
 * consumed by reading it, and the fixture holds one value where driving two
 * sides needs a fresh stream each. The suite tier's checks over ingest are
 * the ones a stored value can answer, and they exist because Source declares
 * its pair — which is what this fixture is here to show.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=streamconsumertest/
export interface StreamConsumer {
  /**
   * Reads the source to its end, and returns the count of the values it read.
   *
   * @remarks
   * ingest takes one parameter and returns one value. Both counts are
   * load-bearing: a second parameter or a second value takes the method out
   * of the shape.
   */
  ingest(src: Source): Promise<number>;
}

/** Returns the seed of a Source, which the harness finds by name. */
export function sourceDefaults(): Source {
  return new ListSource([]);
}

/**
 * Returns the second value of a Source.
 *
 * @remarks
 * The fixture directive names this function, because nothing can derive a
 * second value that differs from the first.
 */
export function sourceAlternate(): Source {
  return new ListSource([{ key: "k", body: "b" }]);
}

/** Walks a fixed list, and then reports exhaustion. */
class ListSource implements Source {
  readonly #items: Value[];

  constructor(items: Value[]) {
    this.#items = items;
  }

  next(): Promise<Value | undefined> {
    return Promise.resolve(this.#items.shift());
  }
}
