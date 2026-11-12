import { compatV1 } from "@dokimi/runtime/conform/compat";
import { GeneratedError } from "@dokimi/runtime/conform/errors";
import type { Handle } from "@dokimi/runtime/conform/handle";
import type { Options } from "@dokimi/runtime/conform/options";
import { run } from "@dokimi/runtime/conform/run";
import type { Callable, Check, Id, Suite } from "@dokimi/runtime/conform/suite";
import { reportsNoError } from "@dokimi/runtime/probe/reports-no-error";
import { survives } from "@dokimi/runtime/probe/survives";
import { underCancelled } from "@dokimi/runtime/probe/under-cancelled";
import type { Resettable, Single } from "./single";

/**
 * The handle of a check of the suite of Single.
 */
export type SingleHandle = Handle<Single>;

/**
 * A check of the suite of Single.
 */
export type SingleCheck = Check<Single>;

/**
 * The options of an application of the suite of Single.
 */
export type SingleOptions = Options<Single>;

/**
 * The name of each callable of Single, with the identity of each of its
 * checks.
 */
export const singleChecks = {
  apply: {
    callable: "Single.apply" as Callable<Single>,
    smoke: (): Id<Single> => "Single.apply/smoke",
    cancel: (): Id<Single> => "Single.apply/cancel",
  },
};

/**
 * Returns the checks of Single, the claims that the generator refuses for
 * its callables, and the classifications of its callables without a check.
 */
export function singleSuite(): Suite<Single> {
  return {
    declaration: "Single",
    package: "single",
    generator: { name: "dokimi", version: "0.1.0" },
    checks: [
      {
        id: singleChecks.apply.smoke(),
        statement: "apply returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, (signal) => h.subject().apply(signal));
        },
        defect: singleApplyPanics,
      },
      {
        id: singleChecks.apply.cancel(),
        statement: "apply returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["nothing-beyond-a-caller"],
        apply: async (h) => {
          await underCancelled(h, (signal) => h.subject().apply(signal));
        },
        defect: singleApplyIgnoresCancel,
      },
    ],
    refusals: [],
    unclaimed: [
      {
        classification: "shape:lifecycle",
        callable: "Single.apply",
        stance: "silent",
        note: "A lifecycle call returns nothing that a check could compare.",
      },
    ],
  };
}

/**
 * Registers the application of the checks of Single, and of the checks
 * that options add, to a fresh subject of newSubject for each check.
 */
export function runSingleSuite(newSubject: () => Single | Promise<Single>, options?: SingleOptions): void {
  runSingleSuiteFromHandle(() => newSubject(), options);
}

/**
 * Registers the same application to a subject that newSubject builds from
 * the handle of each check.
 */
export function runSingleSuiteFromHandle(newSubject: (h: SingleHandle) => Single | Promise<Single>, options?: SingleOptions): void {
  run(singleSuite(), newSubject, options);
}

/**
 * Returns the planted defect of a Single whose apply panics.
 */
function singleApplyPanics(_h: SingleHandle): Single {
  return {
    apply: () => {
      throw new TypeError("the planted defect of a Single whose apply panics");
    },
  };
}

/**
 * Returns the planted defect of a Single whose apply ignores the context
 * over a fresh subject.
 */
async function singleApplyIgnoresCancel(h: SingleHandle): Promise<Single> {
  const s = await h.fresh();
  return {
    apply: async (p0) => {
      if (p0?.aborted) {
        return;
      }
      return s.apply(p0);
    },
  };
}

void compatV1;

/**
 * The handle of a check of the suite of Resettable.
 */
export type ResettableHandle = Handle<Resettable>;

/**
 * A check of the suite of Resettable.
 */
export type ResettableCheck = Check<Resettable>;

/**
 * The options of an application of the suite of Resettable.
 */
export type ResettableOptions = Options<Resettable>;

/**
 * The name of each callable of Resettable, with the identity of each of
 * its checks.
 */
export const resettableChecks = {
  reset: {
    callable: "Resettable.reset" as Callable<Resettable>,
    smoke: (): Id<Resettable> => "Resettable.reset/smoke",
    freshReportsNoPoison: (): Id<Resettable> => "Resettable.reset/fresh-reports-no-poison",
  },
};

/**
 * Returns the checks of Resettable, the claims that the generator refuses
 * for its callables, and the classifications of its callables without a
 * check.
 */
export function resettableSuite(): Suite<Resettable> {
  return {
    declaration: "Resettable",
    package: "single",
    generator: { name: "dokimi", version: "0.1.0" },
    checks: [
      {
        id: resettableChecks.reset.smoke(),
        statement: "reset returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject().reset());
        },
        defect: resettableResetPanics,
      },
      {
        id: resettableChecks.reset.freshReportsNoPoison(),
        statement: "reset reports no error for a fresh receiver",
        claim: "fresh-reports-no-poison",
        machineries: ["nothing-beyond-a-caller"],
        apply: async (h) => {
          await reportsNoError(h, () => h.subject().reset());
        },
        defect: resettableResetReportsPoison,
      },
    ],
    refusals: [],
    unclaimed: [],
  };
}

/**
 * Registers the application of the checks of Resettable, and of the checks
 * that options add, to a fresh subject of newSubject for each check.
 */
export function runResettableSuite(newSubject: () => Resettable | Promise<Resettable>, options?: ResettableOptions): void {
  runResettableSuiteFromHandle(() => newSubject(), options);
}

/**
 * Registers the same application to a subject that newSubject builds from
 * the handle of each check.
 */
export function runResettableSuiteFromHandle(newSubject: (h: ResettableHandle) => Resettable | Promise<Resettable>, options?: ResettableOptions): void {
  run(resettableSuite(), newSubject, options);
}

/**
 * Returns the planted defect of a Resettable whose reset panics.
 */
function resettableResetPanics(_h: ResettableHandle): Resettable {
  return {
    reset: () => {
      throw new TypeError("the planted defect of a Resettable whose reset panics");
    },
  };
}

/**
 * Returns the planted defect of a Resettable whose reset reports an error
 * for a fresh receiver.
 */
function resettableResetReportsPoison(_h: ResettableHandle): Resettable {
  return {
    reset: () => {
      throw new GeneratedError("the planted defect of a Resettable whose reset reports an error for a fresh receiver");
    },
  };
}

void compatV1;
