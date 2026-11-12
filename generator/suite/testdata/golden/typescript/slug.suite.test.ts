import { byIdentity, check } from "@dokimi/assert";
import { falsify } from "@dokimi/runtime/conform/falsify";
import { slug } from "./slug";
import { type SlugFunc, type SlugHandle, runSlugSuite, runSlugSuiteFromHandle, slugSuite } from "./slug.suite";

falsify(slugSuite());

runSlugSuite(slug, {
  checks: [
    {
      id: "slug/receives-the-function",
      statement: "the subject of a check is the function of the application",
      apply: async (h: SlugHandle) => {
        check.isEqualTo(h, h.subject(), slug, h.statement(), byIdentity());
      },
      defect: (): SlugFunc => slug.bind(undefined),
    },
  ],
  excuseSuite: "the test applies the check of its option alone",
});

runSlugSuiteFromHandle(() => slug, {
  checks: [
    {
      id: "slug/builds-the-subject",
      statement: "the subject of a check is the function that the factory returns",
      apply: async (h: SlugHandle) => {
        check.isEqualTo(h, h.subject(), slug, h.statement(), byIdentity());
      },
      defect: (): SlugFunc => slug.bind(undefined),
    },
  ],
  excuseSuite: "the test applies the check of its option alone",
});
