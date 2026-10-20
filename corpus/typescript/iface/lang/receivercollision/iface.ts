// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for parameters whose identifiers are the
 * ones a generated method binds to its locals.
 *
 * @remarks
 * A TypeScript method has no receiver to collide with, but a method of the
 * double binds the call it records to `call` and the function that the test
 * sets to `fn`. A parameter of either name declares the identifier twice in
 * one body and does not compile. The generator cannot pick the locals before
 * it knows the parameters, which is why nothing catches this by inspection:
 * the double is correct for every interface whose author happened not to use
 * those names.
 *
 * This axis varies the type system rather than the classification, so these
 * break generators independently of any directive.
 *
 * No suite directive, and that is the fixture's second finding. The double
 * derives a recorded field per parameter name AND type, so the call of put
 * and the key of get are separate values; the suite derives one fixture field
 * per role, and a role is the parameter name alone.
 *
 * @packageDocumentation
 */

/** Represents the payload, named so that a parameter of it reads naturally as call. */
export interface Call {
  /** Identifies the call. */
  id: string;
}

/**
 * Represents a store that declares the collision on every shape with a
 * parameter.
 *
 * @remarks
 * Locals bound before the parameters are known fail here rather than in one
 * arm of the template that a fixture happened to miss.
 */
//dokimi:stub
//dokimi:out path=receivercollisiontest/
export interface Store {
  /** Stores the call, as the writer form: one parameter, named for its type. */
  put(call: Call): Promise<void>;

  /** Returns the call under the key fn, as the reader form with the collision on the key. */
  get(fn: string): Promise<Call>;

  /** Marks the call as used, as the form that returns nothing. */
  touch(call: Call): Promise<void>;
}
