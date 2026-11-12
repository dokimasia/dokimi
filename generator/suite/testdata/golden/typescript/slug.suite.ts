import { prop } from "@dokimi/assert";
import { compatV1 } from "@dokimi/runtime/conform/compat";
import type { Handle } from "@dokimi/runtime/conform/handle";
import type { Options } from "@dokimi/runtime/conform/options";
import { run } from "@dokimi/runtime/conform/run";
import { type Callable, type Check, type Id, type Suite, generated } from "@dokimi/runtime/conform/suite";
import { repeatAgrees } from "@dokimi/runtime/probe/repeat-agrees";
import { survives } from "@dokimi/runtime/probe/survives";

/**
 * The handle of a check of the suite of slug.
 */
export type SlugHandle = Handle<SlugFunc>;

/**
 * A check of the suite of slug.
 */
export type SlugCheck = Check<SlugFunc>;

/**
 * The options of an application of the suite of slug.
 */
export type SlugOptions = Options<SlugFunc>;

/**
 * The type of the function slug. Each subject of the suite of slug has it.
 */
export type SlugFunc = (title: string) => string;

/**
 * The name of each callable of slug, with the identity of each of its
 * checks.
 */
export const slugChecks = {
  callable: "slug" as Callable<SlugFunc>,
  smoke: (): Id<SlugFunc> => "slug/smoke",
  pure: (): Id<SlugFunc> => "slug/pure",
};

/**
 * The generator of the values of string that the generated bindings draw.
 */
const slugTitleShape = prop.ofShape<string>('{"shape":"string"}');

/**
 * Returns the checks of slug, the claims that the generator refuses for
 * its callables, and the classifications of its callables without a check.
 */
export function slugSuite(): Suite<SlugFunc> {
  return {
    declaration: "slug",
    package: "slug",
    generator: { name: "dokimi", version: "0.1.0" },
    checks: [
      {
        id: slugChecks.smoke(),
        statement: "slug returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await survives(h, () => h.subject()("title-a"));
        },
        defect: slugPanics,
      },
      {
        id: generated(slugChecks.smoke()),
        statement: "slug returns without a panic for derived inputs",
        claim: "smoke",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnTitle = c.draw(slugTitleShape, "title");
          await survives(h, () => h.subject()(drawnTitle));
        },
        defect: slugPanics,
      },
      {
        id: slugChecks.pure(),
        statement: "slug returns the same value on a repeat",
        claim: "pure",
        machineries: ["nothing-beyond-a-caller"],
        shallow: true,
        apply: async (h) => {
          await repeatAgrees(h, () => h.subject()("title-a"));
        },
        defect: slugDrifts,
      },
      {
        id: generated(slugChecks.pure()),
        statement: "slug returns the same value on a repeat",
        claim: "pure",
        machineries: ["generated-inputs", "nothing-beyond-a-caller"],
        property: async (h, c) => {
          const drawnTitle = c.draw(slugTitleShape, "title");
          await repeatAgrees(h, () => h.subject()(drawnTitle));
        },
        defect: slugDrifts,
      },
    ],
    refusals: [
      {
        claim: "zero-on-error",
        callable: "slug",
        gate: "returns-beside-error",
        why: "slug returns no value when it fails",
        remedy: "remove the classification, because a failed call of slug returns no value",
      },
    ],
    unclaimed: [],
  };
}

/**
 * Registers the application of the checks of slug, and of the checks that
 * options add, to subject.
 */
export function runSlugSuite(subject: SlugFunc, options?: SlugOptions): void {
  runSlugSuiteFromHandle(() => subject, options);
}

/**
 * Registers the same application to a subject that newSubject builds from
 * the handle of each check.
 */
export function runSlugSuiteFromHandle(newSubject: (h: SlugHandle) => SlugFunc | Promise<SlugFunc>, options?: SlugOptions): void {
  run(slugSuite(), newSubject, options);
}

/**
 * Returns the planted defect of a slug whose slug panics.
 */
function slugPanics(_h: SlugHandle): SlugFunc {
  return () => {
      throw new TypeError("the planted defect of a slug whose slug panics");
    };
}

/**
 * Returns the planted defect of a slug whose slug returns another value on
 * the repeat.
 */
function slugDrifts(_h: SlugHandle): SlugFunc {
  let calls = 0;
  return () => {
      calls += 1;
      return calls === 1 ? "sample-a" : "sample-b";
    };
}

void compatV1;
