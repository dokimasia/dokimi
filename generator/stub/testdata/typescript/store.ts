/** Represents an item that a store keeps. */
export interface Item {
  /** Contains the identifier of the item. */
  id: string;
}

/** Reports an id without an item. */
export const notFound = new Error("not found");

/**
 * Represents a store with a member of each form of a callable and a property.
 *
 * @remarks
 * The store extends the method close from {@link Closer}.
 */
//dokimi:stub
export interface Store extends Closer {
  /** Contains the zone of the store, which the double declares as a field. */
  readonly zone: string;
  /** Contains an optional label. */
  label?: string;
  /**
   * Returns the item of an id.
   *
   * @remarks
   * The double recovers on the third attempt and selects the calls of a
   * partition fault by the id. Its stub has a helper of the error value
   * {@link notFound}.
   */
  //dokimi:fault errors=[notFound] retry=3 partition=id
  get(id: string): Promise<Item>;
  /** Keeps an item, with an optional parameter and no result. */
  put(item: Item, note?: string): void;
  /** Returns the count of the items. */
  size(): number;
  /** Writes a line from a format and its variadic arguments. */
  log(format: string, ...args: string[]): void;
}

/** Represents the release of the resources of a store. */
export interface Closer {
  /** Releases the resources, and returns a promise of no value. */
  close(): Promise<void>;
}

/** Represents a repository with a type parameter. */
//dokimi:stub
export interface Repo<T> {
  /** Returns a value of the type parameter. */
  find(id: string): T;
}

/** Represents a wrapper with a type parameter whose name contains a dollar sign. */
//dokimi:stub
export interface Wrap<T$> {
  /** Keeps an array of the type parameter. */
  put(items: T$[]): void;
}

/**
 * Represents a table that is generic over its keys and values.
 *
 * @remarks
 * get mentions both type parameters, and size mentions neither.
 */
//dokimi:stub
export interface Table<K, V> {
  /** Returns the value under the key. */
  get(key: K): Promise<V>;
  /** Returns the count of the keys. */
  size(): number;
}
