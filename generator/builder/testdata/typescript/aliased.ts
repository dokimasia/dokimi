import { Account } from "./catalog";

/**
 * Represents an alias of the class {@link Account} of another module.
 *
 * @remarks
 * TypeScript has no value of an alias, so the builder constructs Account.
 */
//dokimi:builder
export type Mirror = Account;
