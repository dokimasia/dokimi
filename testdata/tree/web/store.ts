/** Represents an item that a store keeps. */
export interface Item {
  /** Contains the identifier of the item. */
  id: string;
}

/** Represents the store of the items of a customer. */
//dokimi:stub
export interface Store {
  /** Contains the zone of the store. */
  readonly zone: string;
  /** Returns the item of an id. */
  get(id: string): Promise<Item>;
  /** Keeps an item with an optional note. */
  put(item: Item, note?: string): void;
  /** Writes a line of the store. */
  log(format: string, ...args: string[]): void;
}
