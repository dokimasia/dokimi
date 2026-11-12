import { byIdentity, check } from "@dokimi/assert";
import { falsify } from "@dokimi/runtime/conform/falsify";
import { expect, it } from "vitest";
import type { Store } from "./store";
import { type StoreHandle, runStoreSuite, runStoreSuiteFromHandle, storeChecks, storeSuite } from "./store.suite";

falsify(storeSuite());

/**
 * The subject that the factory of the last check of an application built.
 */
let storeBuilt: Store | undefined;

runStoreSuite(
  () => {
    storeBuilt = {} as Store;
    return storeBuilt;
  },
  {
    checks: [
      {
        id: "Store/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: StoreHandle) => {
          check.isEqualTo<Store | undefined>(h, h.subject(), storeBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Store,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);

runStoreSuiteFromHandle(
  () => {
    storeBuilt = {} as Store;
    return storeBuilt;
  },
  {
    checks: [
      {
        id: "Store/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: StoreHandle) => {
          check.isEqualTo<Store | undefined>(h, h.subject(), storeBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Store,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);

it("storeChecks", () => {
  expect(storeChecks.get.keepsTheItem()).isEqualTo(
    "Store.get/keeps-the-item",
    "the identity of the check joins the callable and the segment of the check",
  );
  expect(storeChecks.get.missesTheUnknown()).isEqualTo(
    "Store.get/misses-the-unknown",
    "the identity of the check joins the callable and the segment of the check",
  );
});
