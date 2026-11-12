/** Represents an item that a store keeps. */
export interface Item {
  /** Contains the identifier of the item. */
  id: string;
  /** Contains the count of the item. */
  count: number;
  /** Contains an optional label. */
  label?: string;
}

/**
 * Represents a store under the directive suite.
 *
 * @remarks
 * get takes a signal, put takes an optional signal, and size is synchronous
 * and takes none.
 */
//dokimi:suite
export interface Store {
  /** Contains the name of the store, which a planted defect reads from its fresh subject. */
  readonly zone: string;
  /** Contains the region of the store, which a planted defect leaves out. */
  readonly region?: string;
  /** Returns the item of an id. */
  get(id: string, signal: AbortSignal): Promise<Item>;
  /** Keeps an item. */
  put(item: Item, signal?: AbortSignal): Promise<void>;
  /** Returns the count of the items. */
  size(): number;
}
