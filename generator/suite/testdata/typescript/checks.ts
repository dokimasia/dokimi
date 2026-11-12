import { check, prop } from "@dokimi/assert";
import type { Store } from "./store";
import type { StoreHandle } from "./store.suite";

/** Checks that get returns the item that put stored. */
//dokimi:check Store.get defect=forgets
export async function keepsTheItem(h: StoreHandle): Promise<void> {
  const s = h.subject();
  await s.put({ id: "a", count: 1 }, h.signal);
  const got = await s.get("a", h.signal).catch(() => undefined);
  check.isEqualTo(h, got?.id, "a", h.statement());
}

/** Checks that get rejects an id that put did not store. */
//dokimi:check Store.get defect=remembers
export async function missesTheUnknown(h: StoreHandle, c: prop.Case): Promise<void> {
  const id = c.draw(prop.string(), "id");
  const found = await h
    .subject()
    .get(id, h.signal)
    .then(
      () => true,
      () => false,
    );
  check.isFalse(h, found, h.statement());
}

/** Returns a store whose get forgets each item. */
export function forgets(_h: StoreHandle): Store {
  return {
    zone: "forgetful",
    get: () => Promise.reject(new Error("forgotten")),
    put: async (_item, signal) => {
      signal?.throwIfAborted();
    },
    size: () => 0,
  };
}

/** Returns a store whose get returns an item for each id. */
export function remembers(_h: StoreHandle): Store {
  return {
    zone: "remembering",
    get: async (id, signal) => {
      signal.throwIfAborted();
      return { id, count: 0 };
    },
    put: async (_item, signal) => {
      signal?.throwIfAborted();
    },
    size: () => 0,
  };
}
