import type { prop } from "@dokimi/assert";
import type { Store } from "./store";
import type { StoreHandle } from "./store.suite";

/** Represents a case of this module, which the property engine does not pass. */
export interface Case {
  /** Contains the seed of the case. */
  readonly seed: number;
}

/** Checks that get returns the item that put stored. */
//dokimi:check Store.get defect=forgets
export async function takesAnotherCase(h: StoreHandle, _c: Case): Promise<void> {
  await h.subject().get("a", h.signal);
}

/** Checks that get returns the item that put stored. */
//dokimi:check Store.get defect=forgets
export async function takesAnOption(h: StoreHandle, _o: prop.Option): Promise<void> {
  await h.subject().get("a", h.signal);
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
