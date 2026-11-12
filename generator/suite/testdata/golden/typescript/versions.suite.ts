import { prop } from "@dokimi/assert";
import { compatV1 } from "@dokimi/runtime/conform/compat";
import type { Handle } from "@dokimi/runtime/conform/handle";
import type { Options } from "@dokimi/runtime/conform/options";
import { complete } from "@dokimi/runtime/conform/output";
import { run } from "@dokimi/runtime/conform/run";
import { type Callable, type Check, type Id, type Suite, generated } from "@dokimi/runtime/conform/suite";
import { readVersionsNeverFall } from "@dokimi/runtime/probe/read-versions-never-fall";
import type { Actions } from "@dokimi/runtime/probe/sequence";
import { survives } from "@dokimi/runtime/probe/survives";
import { underCancelled } from "@dokimi/runtime/probe/under-cancelled";
import type { Rec, Versions } from "./versions";

/**
 * The handle of a check of the suite of Versions.
 */
export type VersionsHandle = Handle<Versions>;

/**
 * A check of the suite of Versions.
 */
export type VersionsCheck = Check<Versions>;

/**
 * The options of an application of the suite of Versions.
 */
export type VersionsOptions = Options<Versions>;

/**
 * The name of each callable of Versions, with the identity of each of its
 * checks.
 */
export const versionsChecks = {
  version: {
    callable: "Versions.version" as Callable<Versions>,
    smoke: (): Id<Versions> => "Versions.version/smoke",
    cancel: (): Id<Versions> => "Versions.version/cancel",
    readsNeverGoBack: (): Id<Versions> => "Versions.version/reads-never-go-back",
  },
  put: {
    callable: "Versions.put" as Callable<Versions>,
    smoke: (): Id<Versions> => "Versions.put/smoke",
    cancel: (): Id<Versions> => "Versions.put/cancel",
  },
  replace: {
    callable: "Versions.replace" as Callable<Versions>,
    smoke: (): Id<Versions> => "Versions.replace/smoke",
    cancel: (): Id<Versions> => "Versions.replace/cancel",
  },
};

/**
 * The generator of the values of string that the generated bindings draw.
 */
const versionsKeyShape = prop.ofShape<string>('{"shape":"string"}');

/**
 * The generator of the values of Rec that the generated bindings draw.
 */
const versionsRShape = prop.ofShape<Rec>('{"definitions":{"Rec":{"fields":[["key",{"shape":"string"}],["rev",{"shape":"float","width":64}]],"shape":"record"}},"name":"Rec","shape":"ref"}');

/**
 * Returns the checks of Versions, the claims that the generator refuses
 * for its callables, and the classifications of its callables without a
 * check.
 */
export function versionsSuite(): Suite<Versions> {
  return {
    declaration: "Versions",
    package: "versions",
    generator: { name: "dokimi", version: "0.1.0" },
    checks: [
      {
        id: versionsChecks.version.smoke(),
        statement: "version returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, (signal) => h.subject().version("key-a", signal));
        },
        defect: versionsVersionPanics,
      },
      {
        id: generated(versionsChecks.version.smoke()),
        statement: "version returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnKey = c.draw(versionsKeyShape, "key");
          await survives(h, (signal) => h.subject().version(drawnKey, signal));
        },
        defect: versionsVersionPanics,
      },
      {
        id: versionsChecks.version.cancel(),
        statement: "version returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["nothing-beyond-a-caller"],
        apply: async (h) => {
          await underCancelled(h, (signal) => h.subject().version("key-a", signal));
        },
        defect: versionsVersionIgnoresCancel,
      },
      {
        id: generated(versionsChecks.version.cancel()),
        statement: "version returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnKey = c.draw(versionsKeyShape, "key");
          await underCancelled(h, (signal) => h.subject().version(drawnKey, signal));
        },
        defect: versionsVersionIgnoresCancel,
      },
      {
        id: versionsChecks.version.readsNeverGoBack(),
        statement: "version never returns a version below one that it returned before, over any sequence of calls",
        claim: "reads-never-go-back",
        machineries: ["generated-sequences"],
        property: async (h, c) => {
          await readVersionsNeverFall(h, c, versionsActions, "version", (v: Rec) => v.rev);
        },
        defect: versionsVersionReadsThePastAfterThree,
      },
      {
        id: versionsChecks.put.smoke(),
        statement: "put returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, (signal) => h.subject().put({ key: "key-a", rev: 1.5 }, signal));
        },
        defect: versionsPutPanics,
      },
      {
        id: generated(versionsChecks.put.smoke()),
        statement: "put returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnR = c.draw(versionsRShape, "r");
          await survives(h, (signal) => h.subject().put(drawnR, signal));
        },
        defect: versionsPutPanics,
      },
      {
        id: versionsChecks.put.cancel(),
        statement: "put returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["nothing-beyond-a-caller"],
        apply: async (h) => {
          await underCancelled(h, (signal) => h.subject().put({ key: "key-a", rev: 1.5 }, signal));
        },
        defect: versionsPutIgnoresCancel,
      },
      {
        id: generated(versionsChecks.put.cancel()),
        statement: "put returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnR = c.draw(versionsRShape, "r");
          await underCancelled(h, (signal) => h.subject().put(drawnR, signal));
        },
        defect: versionsPutIgnoresCancel,
      },
      {
        id: versionsChecks.replace.smoke(),
        statement: "replace returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, (signal) => h.subject().replace({ key: "key-a", rev: 1.5 }, signal));
        },
        defect: versionsReplacePanics,
      },
      {
        id: generated(versionsChecks.replace.smoke()),
        statement: "replace returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnR = c.draw(versionsRShape, "r");
          await survives(h, (signal) => h.subject().replace(drawnR, signal));
        },
        defect: versionsReplacePanics,
      },
      {
        id: versionsChecks.replace.cancel(),
        statement: "replace returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["nothing-beyond-a-caller"],
        apply: async (h) => {
          await underCancelled(h, (signal) => h.subject().replace({ key: "key-a", rev: 1.5 }, signal));
        },
        defect: versionsReplaceIgnoresCancel,
      },
      {
        id: generated(versionsChecks.replace.cancel()),
        statement: "replace returns an error for a cancelled context",
        claim: "cancel",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnR = c.draw(versionsRShape, "r");
          await underCancelled(h, (signal) => h.subject().replace(drawnR, signal));
        },
        defect: versionsReplaceIgnoresCancel,
      },
    ],
    refusals: [
      {
        claim: "zero-on-error",
        callable: "Versions.version",
        gate: "returns-beside-error",
        why: "version returns no value when it fails",
        remedy: "remove the classification, because a failed call of version returns no value",
      },
      {
        claim: "monotonic-reads",
        callable: "Versions.version",
        gate: "one-partner",
        why: "two callables of the declaration can be the partner write of version",
        remedy: "name the partner of version with the param of the classification",
      },
    ],
    unclaimed: [
      {
        classification: "shape:writer",
        callable: "Versions.put",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
      {
        classification: "shape:writer",
        callable: "Versions.replace",
        stance: "silent",
        note: "A write is an input of the other checks, and its result is an error alone.",
      },
    ],
  };
}

/**
 * Registers the application of the checks of Versions, and of the checks
 * that options add, to a fresh subject of newSubject for each check.
 */
export function runVersionsSuite(newSubject: () => Versions | Promise<Versions>, options?: VersionsOptions): void {
  runVersionsSuiteFromHandle(() => newSubject(), options);
}

/**
 * Registers the same application to a subject that newSubject builds from
 * the handle of each check.
 */
export function runVersionsSuiteFromHandle(newSubject: (h: VersionsHandle) => Versions | Promise<Versions>, options?: VersionsOptions): void {
  run(versionsSuite(), newSubject, options);
}

/**
 * Returns the planted defect of a Versions whose version panics over a
 * fresh subject.
 */
async function versionsVersionPanics(h: VersionsHandle): Promise<Versions> {
  const s = await h.fresh();
  return {
    version: () => {
      throw new TypeError("the planted defect of a Versions whose version panics");
    },
    put: s.put.bind(s),
    replace: s.replace.bind(s),
  };
}

/**
 * Returns the planted defect of a Versions whose version ignores the
 * context over a fresh subject.
 */
async function versionsVersionIgnoresCancel(h: VersionsHandle): Promise<Versions> {
  const s = await h.fresh();
  return {
    version: async (p0, p1) => {
      if (p1?.aborted) {
        return { key: "key-b", rev: 1.5 };
      }
      return s.version(p0, p1);
    },
    put: s.put.bind(s),
    replace: s.replace.bind(s),
  };
}

/**
 * Returns the planted defect of a Versions whose version returns its first
 * value again from its fourth call on over a fresh subject.
 */
async function versionsVersionReadsThePastAfterThree(h: VersionsHandle): Promise<Versions> {
  const s = await h.fresh();
  let calls = 0;
  let first: { readonly value: Rec } | undefined;
  return {
    version: async (p0, p1) => {
      calls += 1;
      const r0 = await s.version(p0, p1);
      first ??= { value: r0 };
      return calls > 3 ? first.value : r0;
    },
    put: s.put.bind(s),
    replace: s.replace.bind(s),
  };
}

/**
 * Returns the planted defect of a Versions whose put panics over a fresh
 * subject.
 */
async function versionsPutPanics(h: VersionsHandle): Promise<Versions> {
  const s = await h.fresh();
  return {
    version: s.version.bind(s),
    put: () => {
      throw new TypeError("the planted defect of a Versions whose put panics");
    },
    replace: s.replace.bind(s),
  };
}

/**
 * Returns the planted defect of a Versions whose put ignores the context
 * over a fresh subject.
 */
async function versionsPutIgnoresCancel(h: VersionsHandle): Promise<Versions> {
  const s = await h.fresh();
  return {
    version: s.version.bind(s),
    put: async (p0, p1) => {
      if (p1?.aborted) {
        return;
      }
      return s.put(p0, p1);
    },
    replace: s.replace.bind(s),
  };
}

/**
 * Returns the planted defect of a Versions whose replace panics over a
 * fresh subject.
 */
async function versionsReplacePanics(h: VersionsHandle): Promise<Versions> {
  const s = await h.fresh();
  return {
    version: s.version.bind(s),
    put: s.put.bind(s),
    replace: () => {
      throw new TypeError("the planted defect of a Versions whose replace panics");
    },
  };
}

/**
 * Returns the planted defect of a Versions whose replace ignores the
 * context over a fresh subject.
 */
async function versionsReplaceIgnoresCancel(h: VersionsHandle): Promise<Versions> {
  const s = await h.fresh();
  return {
    version: s.version.bind(s),
    put: s.put.bind(s),
    replace: async (p0, p1) => {
      if (p1?.aborted) {
        return;
      }
      return s.replace(p0, p1);
    },
  };
}

/**
 * The actions of the generated sequences of Versions over the subject that
 * subject returns. The actions of a case take their keys from the pools
 * that it draws.
 */
const versionsActions: Actions<Versions> = (c, subject) => {
  const keysString = c.draw(prop.list(versionsKeyShape, { unique: true, minSize: 1, maxSize: 3 }), "keys.string");
  return [
    {
      name: "version",
      input: (step) => [
        step.draw(prop.sampledFrom(...keysString), "version.key"),
      ],
      run: async (step, client, input) => {
        const args = input as readonly unknown[];
        const call = step.history().invoke(client, "version", args, args[0]);
        await versionsVersionCall(step.signal, subject(), args).then(
          (values) => {
            complete(step.signal, call, undefined, ...values);
          },
          (error: unknown) => {
            complete(step.signal, call, error);
          },
        );
      },
    },
    {
      name: "put",
      input: (step) => [
        Object.assign(step.draw(versionsRShape, "put.r"), { key: step.draw(prop.sampledFrom(...keysString), "put.r.key"), }),
      ],
      run: async (step, client, input) => {
        const args = input as readonly unknown[];
        const call = step.history().invoke(client, "put", args, (args[0] as Rec).key);
        await versionsPutCall(step.signal, subject(), args).then(
          (values) => {
            complete(step.signal, call, undefined, ...values);
          },
          (error: unknown) => {
            complete(step.signal, call, error);
          },
        );
      },
    },
    {
      name: "replace",
      input: (step) => [
        Object.assign(step.draw(versionsRShape, "replace.r"), { key: step.draw(prop.sampledFrom(...keysString), "replace.r.key"), }),
      ],
      run: async (step, client, input) => {
        const args = input as readonly unknown[];
        const call = step.history().invoke(client, "replace", args, (args[0] as Rec).key);
        await versionsReplaceCall(step.signal, subject(), args).then(
          (values) => {
            complete(step.signal, call, undefined, ...values);
          },
          (error: unknown) => {
            complete(step.signal, call, error);
          },
        );
      },
    },
  ];
};

/**
 * Calls version on subject with the arguments args, and resolves with the
 * values of the call.
 */
async function versionsVersionCall(signal: AbortSignal, subject: Versions, args: readonly unknown[]): Promise<readonly unknown[]> {
  return [await subject.version(args[0] as string, signal)];
}

/**
 * Calls put on subject with the arguments args, and resolves with the
 * values of the call.
 */
async function versionsPutCall(signal: AbortSignal, subject: Versions, args: readonly unknown[]): Promise<readonly unknown[]> {
  await subject.put(args[0] as Rec, signal);
  return [];
}

/**
 * Calls replace on subject with the arguments args, and resolves with the
 * values of the call.
 */
async function versionsReplaceCall(signal: AbortSignal, subject: Versions, args: readonly unknown[]): Promise<readonly unknown[]> {
  await subject.replace(args[0] as Rec, signal);
  return [];
}

void compatV1;
