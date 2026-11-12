import { byIdentity, check } from "@dokimi/assert";
import { falsify } from "@dokimi/runtime/conform/falsify";
import type { Versions } from "./versions";
import { type VersionsHandle, runVersionsSuite, runVersionsSuiteFromHandle, versionsSuite } from "./versions.suite";

falsify(versionsSuite());

/**
 * The subject that the factory of the last check of an application built.
 */
let versionsBuilt: Versions | undefined;

runVersionsSuite(
  () => {
    versionsBuilt = {} as Versions;
    return versionsBuilt;
  },
  {
    checks: [
      {
        id: "Versions/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: VersionsHandle) => {
          check.isEqualTo<Versions | undefined>(h, h.subject(), versionsBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Versions,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);

runVersionsSuiteFromHandle(
  () => {
    versionsBuilt = {} as Versions;
    return versionsBuilt;
  },
  {
    checks: [
      {
        id: "Versions/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: VersionsHandle) => {
          check.isEqualTo<Versions | undefined>(h, h.subject(), versionsBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Versions,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);
