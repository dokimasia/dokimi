import type { Item, Store } from "./store";

/** Represents a store in memory, the subject of the suite of Store. */
//dokimi:subject Store new=newMemory
//dokimi:excuse Store.size/smoke reason="size reads no input"
export class Memory implements Store {
  /** Contains the name of the store. */
  readonly zone = "memory";
  readonly #items = new Map<string, Item>();

  /** Returns the item of an id. */
  get(id: string, signal: AbortSignal): Promise<Item> {
    signal.throwIfAborted();
    const item = this.#items.get(id);
    if (item === undefined) {
      return Promise.reject(new Error("memory: not found"));
    }
    return Promise.resolve(item);
  }

  /** Keeps an item. */
  put(item: Item, signal?: AbortSignal): Promise<void> {
    signal?.throwIfAborted();
    this.#items.set(item.id, item);
    return Promise.resolve();
  }

  /** Returns the count of the items. */
  size(): number {
    return this.#items.size;
  }
}

/** Returns a new store in memory. */
export function newMemory(): Memory {
  return new Memory();
}
