import { prop } from "@dokimi/assert";
import { compatV1 } from "@dokimi/runtime/conform/compat";
import type { Handle } from "@dokimi/runtime/conform/handle";
import type { Options } from "@dokimi/runtime/conform/options";
import { run } from "@dokimi/runtime/conform/run";
import { type Callable, type Check, type Id, type Suite, generated } from "@dokimi/runtime/conform/suite";
import { survives } from "@dokimi/runtime/probe/survives";
import { underCancelled } from "@dokimi/runtime/probe/under-cancelled";
import { forgets, keepsTheItem, missesTheUnknown, remembers } from "./checks";
import type { Item, Store } from "./store";

/**
 * The handle of a check of the suite of Store.
 */
export type StoreHandle = Handle<Store>;

/**
 * A check of the suite of Store.
 */
export type StoreCheck = Check<Store>;

/**
 * The options of an application of the suite of Store.
 */
export type StoreOptions = Options<Store>;

/**
 * The name of each callable of Store, with the identity of each of its
 * checks.
 */
export const storeChecks = {
  get: {
    callable: "Store.get" as Callable<Store>,
    keepsTheItem: (): Id<Store> => "Store.get/keeps-the-item",
    missesTheUnknown: (): Id<Store> => "Store.get/misses-the-unknown",
    smoke: (): Id<Store> => "Store.get/smoke",
    cancel: (): Id<Store> => "Store.get/cancel",
  },
  put: {
    callable: "Store.put" as Callable<Store>,
    smoke: (): Id<Store> => "Store.put/smoke",
    cancel: (): Id<Store> => "Store.put/cancel",
  },
  size: {
    callable: "Store.size" as Callable<Store>,
    smoke: (): Id<Store> => "Store.size/smoke",
  },
};

/**
 * The generator of the values of string that the generated bindings draw.
 */
const storeIdShape = prop.ofShape<string>('{"shape":"string"}');

/**
 * The generator of the values of Item that the generated bindings draw.
 */
const storeItemShape = prop.ofShape<Item>('{"definitions":{"Item":{"fields":[["id",{"shape":"string"}],["count",{"shape":"float","width":64}],["label",{"of":{"shape":"string"},"shape":"optional"}]],"shape":"record"}},"name":"Item","shape":"ref"}');

/**
 * Returns the checks of Store, the claims that the generator refuses for
 * its callables, and the classifications of its callables without a check.
 */
export function storeSuite(): Suite<Store> {
  return {
    declaration: "Store",
    package: "store",
    generator: { name: "dokimi", version: "0.1.0" },
    checks: [
      {
        id: storeChecks.get.smoke(),
        statement: "get returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, (signal) => h.subject().get("id-a", signal));
        },
        defect: storeGetPanics,
      },
      {
        id: generated(storeChecks.get.smoke()),
        statement: "get returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnId = c.draw(storeIdShape, "id");
          await survives(h, (signal) => h.subject().get(drawnId, signal));
        },
        defect: storeGetPanics,
      },
      {
        id: storeChecks.get.cancel(),
        statement: "get returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["nothing-beyond-a-caller"],
        apply: async (h) => {
          await underCancelled(h, (signal) => h.subject().get("id-a", signal));
        },
        defect: storeGetIgnoresCancel,
      },
      {
        id: generated(storeChecks.get.cancel()),
        statement: "get returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnId = c.draw(storeIdShape, "id");
          await underCancelled(h, (signal) => h.subject().get(drawnId, signal));
        },
        defect: storeGetIgnoresCancel,
      },
      {
        id: storeChecks.put.smoke(),
        statement: "put returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, (signal) => h.subject().put({ id: "id-a", count: 1.5 }, signal));
        },
        defect: storePutPanics,
      },
      {
        id: generated(storeChecks.put.smoke()),
        statement: "put returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnItem = c.draw(storeItemShape, "item");
          await survives(h, (signal) => h.subject().put(drawnItem, signal));
        },
        defect: storePutPanics,
      },
      {
        id: storeChecks.put.cancel(),
        statement: "put returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["nothing-beyond-a-caller"],
        apply: async (h) => {
          await underCancelled(h, (signal) => h.subject().put({ id: "id-a", count: 1.5 }, signal));
        },
        defect: storePutIgnoresCancel,
      },
      {
        id: generated(storeChecks.put.cancel()),
        statement: "put returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnItem = c.draw(storeItemShape, "item");
          await underCancelled(h, (signal) => h.subject().put(drawnItem, signal));
        },
        defect: storePutIgnoresCancel,
      },
      {
        id: storeChecks.size.smoke(),
        statement: "size returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().size());
        },
        defect: storeSizePanics,
      },
      {
        id: storeChecks.get.keepsTheItem(),
        statement: "get returns the item that put stored",
        apply: keepsTheItem,
        defect: forgets,
      },
      {
        id: storeChecks.get.missesTheUnknown(),
        statement: "get rejects an id that put did not store",
        property: missesTheUnknown,
        defect: remembers,
      },
    ],
    refusals: [
      {
        claim: "zero-on-error",
        callable: "Store.get",
        gate: "returns-beside-error",
        why: "get returns no value when it fails",
        remedy: "remove the classification, because a failed call of get returns no value",
      },
      {
        claim: "zero-on-error",
        callable: "Store.size",
        gate: "returns-beside-error",
        why: "size returns no value when it fails",
        remedy: "remove the classification, because a failed call of size returns no value",
      },
    ],
    unclaimed: [
      {
        classification: "shape:writer",
        callable: "Store.put",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
    ],
  };
}

/**
 * Registers the application of the checks of Store, and of the checks that
 * options add, to a fresh subject of newSubject for each check.
 */
export function runStoreSuite(newSubject: () => Store | Promise<Store>, options?: StoreOptions): void {
  runStoreSuiteFromHandle(() => newSubject(), options);
}

/**
 * Registers the same application to a subject that newSubject builds from
 * the handle of each check.
 */
export function runStoreSuiteFromHandle(newSubject: (h: StoreHandle) => Store | Promise<Store>, options?: StoreOptions): void {
  run(storeSuite(), newSubject, options);
}

/**
 * Returns the planted defect of a Store whose get panics over a fresh
 * subject.
 */
async function storeGetPanics(h: StoreHandle): Promise<Store> {
  const s = await h.fresh();
  return {
    get zone() {
      return s.zone;
    },
    get: () => {
      throw new TypeError("the planted defect of a Store whose get panics");
    },
    put: s.put.bind(s),
    size: s.size.bind(s),
  };
}

/**
 * Returns the planted defect of a Store whose get ignores the context over
 * a fresh subject.
 */
async function storeGetIgnoresCancel(h: StoreHandle): Promise<Store> {
  const s = await h.fresh();
  return {
    get zone() {
      return s.zone;
    },
    get: async (p0, p1) => {
      if (p1?.aborted) {
        return { id: "id-b", count: 1.5 };
      }
      return s.get(p0, p1);
    },
    put: s.put.bind(s),
    size: s.size.bind(s),
  };
}

/**
 * Returns the planted defect of a Store whose put panics over a fresh
 * subject.
 */
async function storePutPanics(h: StoreHandle): Promise<Store> {
  const s = await h.fresh();
  return {
    get zone() {
      return s.zone;
    },
    get: s.get.bind(s),
    put: () => {
      throw new TypeError("the planted defect of a Store whose put panics");
    },
    size: s.size.bind(s),
  };
}

/**
 * Returns the planted defect of a Store whose put ignores the context over
 * a fresh subject.
 */
async function storePutIgnoresCancel(h: StoreHandle): Promise<Store> {
  const s = await h.fresh();
  return {
    get zone() {
      return s.zone;
    },
    get: s.get.bind(s),
    put: async (p0, p1) => {
      if (p1?.aborted) {
        return;
      }
      return s.put(p0, p1);
    },
    size: s.size.bind(s),
  };
}

/**
 * Returns the planted defect of a Store whose size panics over a fresh
 * subject.
 */
async function storeSizePanics(h: StoreHandle): Promise<Store> {
  const s = await h.fresh();
  return {
    get zone() {
      return s.zone;
    },
    get: s.get.bind(s),
    put: s.put.bind(s),
    size: () => {
      throw new TypeError("the planted defect of a Store whose size panics");
    },
  };
}

void compatV1;
