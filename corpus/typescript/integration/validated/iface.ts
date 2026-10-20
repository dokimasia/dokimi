// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the integration fixture for a type whose values a derived sample
 * cannot guess.
 *
 * @remarks
 * The other axes vary one thing about a declaration and check what gets
 * stamped. This one varies nothing about the classification and everything
 * about the subject: a real implementation, with real validation, held to the
 * generated contract. It is where "the generator emitted something" and "the
 * something it emitted is true of a working implementation" stop being the
 * same claim.
 *
 * Account is the point. Its fields accept only some strings — an identifier
 * shaped like a UUID, an address with an `@` — and the sample a generator
 * derives from a field's name and type is `test-id` and `test-email`, which
 * the store refuses. A harness built on the derived value would report a
 * conformance violation against a correct implementation, which is the
 * failure mode that gets a suite switched off.
 *
 * {@link accountDefaults} is the answer, and it is deliberately an ordinary
 * function rather than a directive: the validation rules live in the author's
 * head, and the convention that names them is the one the builder generator
 * already reads. One function, two readers — builder seeds the builder with
 * it, and the suite's fixture takes it over anything it could derive.
 *
 * A fixture holds two values of every type it draws, and a companion answers
 * only the first. Left there, the second is derived — so the pair would be
 * one account the store accepts and one it rejects, which is the same failure
 * a field further in. `//dokimi:fixture alternate=` is the second half, and
 * it is a directive rather than a second convention on purpose: nothing but
 * the harness reads it, so a name found by looking would revert to the
 * derived value the moment someone renamed the function, with nothing
 * failing.
 *
 * Routing is declared once for the module rather than on each declaration:
 * the builder, the double and the harness all belong beside each other, and
 * a per-declaration directive is the same statement written three times.
 *
 * @packageDocumentation
 */
//dokimi:out path=validatedtest/

/**
 * Represents the payload, with fields that admit only some values.
 *
 * @remarks
 * The builder exists so a check that needs a *variant* of a valid account —
 * one field changed, the rest still acceptable — can say so in one line
 * instead of restating every field and hoping it stays valid.
 */
//dokimi:builder
export class Account {
  /** Identifies the account, and the store rejects it unless it is shaped like a UUID. */
  id: string = "";

  /** Contains the email address, which the store rejects unless it has an `@`. */
  email: string = "";
}

/**
 * Returns an Account the store accepts.
 *
 * @remarks
 * Hand-written, because nothing about the declaration says what `id` admits.
 * A generator deriving `test-id` is not wrong to try — it is right for the
 * many types with no rules — and this is how a type with rules says so once,
 * for every generator that needs a value of it.
 */
export function accountDefaults(): Account {
  const account = new Account();
  account.id = "6ba7b810-9dad-11d1-80b4-00c04fd430c8";
  account.email = "conformance@example.com";
  return account;
}

/**
 * Returns a second Account the store accepts, which differs from
 * {@link accountDefaults} in both fields.
 *
 * @remarks
 * Named rather than found: the directive on Account points here. Spelled
 * `Second` rather than `Alternate` deliberately — nothing resolves this by
 * name, so a name that looks like a convention would suggest one exists.
 */
export function accountSecond(): Account {
  const account = new Account();
  account.id = "7c9e6679-7425-40de-944b-e07fc1f90ae7";
  account.email = "second@example.com";
  return account;
}

/** Represents the subject that the generated harness holds an implementation to. */
export interface Store {
  /** Stores the account, and refuses an account that does not validate. */
  put(a: Account): Promise<void>;

  /**
   * Returns the account under the id.
   *
   * @remarks
   * get reports the empty account alongside every rejection, which is the
   * property the reader's own check is about.
   */
  get(id: string): Promise<Account>;
}
