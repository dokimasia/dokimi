// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the integration fixture for the default results of a double,
 * which the builders of the result types construct.
 *
 * @remarks
 * Ledger has a method for each row of the table of default results: a value
 * of a class with a builder, a promise of one, a value or undefined, an
 * instance of a generic class with a builder, a class without a builder and
 * a class whose builder is in another module. The test of the double of
 * Ledger checks the default results of each method.
 *
 * The builders and the double are in one directory, and the double imports
 * each builder from the module of this file.
 *
 * @packageDocumentation
 */
//dokimi:out path=builttest/

import type { Invoice } from "./billing/billing";

/** Represents an account, whose builder {@link accountDefaults} seeds. */
//dokimi:builder
export class Account {
  /** Identifies the account. */
  id: string = "";

  /** Contains the email address of the owner of the account. */
  email: string = "";
}

/** Returns the account that each builder of Account starts from. */
export function accountDefaults(): Account {
  const account = new Account();
  account.id = "acct-1";
  account.email = "owner@example.com";
  return account;
}

/** Represents a page of items, a generic class whose builder builds each instance. */
//dokimi:builder
export class Page<T> {
  /** Contains the items of the page. */
  items: T[] = [];

  /** Contains the cursor of the next page. */
  next: string = "";
}

/** Represents a note, a class without a builder. */
export class Note {
  /** Contains the text of the note. */
  text: string = "";
}

/** Represents a ledger that returns a value of each type of the table of default results. */
//dokimi:stub
export interface Ledger {
  /** Returns a promise of an account. */
  get(id: string): Promise<Account>;

  /** Returns a promise of an account, or of undefined where the ledger has none. */
  find(id: string): Promise<Account | undefined>;

  /** Returns an account, or undefined where the ledger has none. */
  lookup(id: string): Account | undefined;

  /** Returns a promise of a page of accounts. */
  list(): Promise<Page<Account>>;

  /** Returns a note, whose class has no builder, so an unconfigured call returns undefined. */
  note(id: string): Note;

  /**
   * Returns an invoice, whose builder is in another module.
   *
   * @remarks
   * The double refers to a builder of its own module alone, so an unconfigured
   * call returns undefined.
   */
  bill(id: string): Invoice;
}
