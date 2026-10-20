/** Represents the postal address of a customer. */
//dokimi:builder
export class Address {
  /** Contains the street and the house number. */
  //dokimi:default "Dam 1"
  street: string = "";
  /** Contains the city of the address. */
  city: string = "";
  /** Contains the floor of the address, where it has one. */
  floor?: number;
  /** Lists the lines of the address. */
  lines: string[] = [];
  /** Contains the labels of the address. */
  labels: Map<string, string> = new Map();
  /** Contains the tags of the address. */
  tags: Set<string> = new Set();
}

/** Represents the account of a customer. */
//dokimi:builder
export class Account {
  /** Contains the owner of the account. */
  //dokimi:default "owner"
  owner: string = "";
  /** Contains the plan of the account. */
  plan: string = "";
}

/** Returns the account of a new customer. */
export function accountDefaults(): Account {
  const account = new Account();
  account.owner = "seed";
  account.plan = "basic";
  return account;
}

/** Represents a page of results. */
//dokimi:builder
export class Page<T> {
  /** Lists the results on the page. */
  items: T[] = [];
}

/** Counts the visits of a page. */
//dokimi:builder
export class Counter {
  /** Contains the number of visits, with the type of its initializer. */
  count = 0;
}

/** Represents the tier of an account, whose text is its value. */
//dokimi:enum
export enum Tier {
  Free = "free",
  Pro = "pro",
}
