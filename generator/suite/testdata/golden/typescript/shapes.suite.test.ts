import { byIdentity, check } from "@dokimi/assert";
import { falsify } from "@dokimi/runtime/conform/falsify";
import type { Sink } from "./shapes";
import { type SinkHandle, runSinkSuite, runSinkSuiteFromHandle, sinkSuite } from "./shapes.suite";

falsify(sinkSuite());

/**
 * The subject that the factory of the last check of an application built.
 */
let sinkBuilt: Sink | undefined;

runSinkSuite(
  () => {
    sinkBuilt = {} as Sink;
    return sinkBuilt;
  },
  {
    checks: [
      {
        id: "Sink/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: SinkHandle) => {
          check.isEqualTo<Sink | undefined>(h, h.subject(), sinkBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Sink,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);

runSinkSuiteFromHandle(
  () => {
    sinkBuilt = {} as Sink;
    return sinkBuilt;
  },
  {
    checks: [
      {
        id: "Sink/builds-the-subject",
        statement: "the subject of a check is the value that the factory builds",
        apply: async (h: SinkHandle) => {
          check.isEqualTo<Sink | undefined>(h, h.subject(), sinkBuilt, h.statement(), byIdentity());
        },
        defect: () => ({}) as Sink,
      },
    ],
    excuseSuite: "the test applies the check of its option alone",
  },
);
