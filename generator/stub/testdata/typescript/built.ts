import type { Invoice } from "./billing";

/** Represents a value with a builder. */
//dokimi:builder
export class Account {
  /** Identifies the account. */
  id: string = "";
  /** Contains the email address of the account. */
  email: string = "";
}

/** Represents a value without a builder. */
export class Note {
  /** Contains the text of the note. */
  text: string = "";
}

/** Represents a store that returns values of types with a builder, and of types without one. */
//dokimi:stub
export interface Store {
  /** Returns a promise of an account. */
  get(id: string): Promise<Account>;
  /** Returns an account, or undefined where the store has none. */
  find(id: string): Account | undefined;
  /** Returns a value of a type without a builder. */
  note(id: string): Note;
  /** Returns a value of a type whose builder is in another module. */
  bill(id: string): Invoice;
}
