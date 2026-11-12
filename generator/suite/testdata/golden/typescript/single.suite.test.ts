import { byIdentity, check } from "@dokimi/assert";
import { falsify } from "@dokimi/runtime/conform/falsify";
import type { Resettable, Single } from "./single";
import { type ResettableHandle, type SingleHandle, resettableSuite, runResettableSuite, runResettableSuiteFromHandle, runSingleSuite, runSingleSuiteFromHandle, singleSuite } from "./single.suite";

falsify(singleSuite());

/**
 * The subject that the factory of the last check of an application built.
 */
let singleBuilt: Single | undefined;

runSingleSuite(
  () => {
    singleBuilt = {} as Single;
    return singleBuilt;
  },
  {
    checks: [
      {
        id: "Single/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: SingleHandle) => {
          check.isEqualTo<Single | undefined>(h, h.subject(), singleBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Single,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);

runSingleSuiteFromHandle(
  () => {
    singleBuilt = {} as Single;
    return singleBuilt;
  },
  {
    checks: [
      {
        id: "Single/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: SingleHandle) => {
          check.isEqualTo<Single | undefined>(h, h.subject(), singleBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Single,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);

falsify(resettableSuite());

/**
 * The subject that the factory of the last check of an application built.
 */
let resettableBuilt: Resettable | undefined;

runResettableSuite(
  () => {
    resettableBuilt = {} as Resettable;
    return resettableBuilt;
  },
  {
    checks: [
      {
        id: "Resettable/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: ResettableHandle) => {
          check.isEqualTo<Resettable | undefined>(h, h.subject(), resettableBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Resettable,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);

runResettableSuiteFromHandle(
  () => {
    resettableBuilt = {} as Resettable;
    return resettableBuilt;
  },
  {
    checks: [
      {
        id: "Resettable/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: ResettableHandle) => {
          check.isEqualTo<Resettable | undefined>(h, h.subject(), resettableBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Resettable,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);
