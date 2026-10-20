import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { behaviour, settle } from '@dokimi/runtime/stub/behaviour';
import { doubleBehaviour } from '@dokimi/runtime/stub/double-behaviour';
import { describe } from 'vitest';
import { FinderStub } from './finder.stub';

/** Returns the samples of the arguments of findById, or undefined. */
function finderStubFindByIdArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'findById-a';
  return [a0];
}

/** Calls findById of double with the samples of its arguments. */
function finderStubFindById(double: FinderStub): string {
  const args = finderStubFindByIdArgs();
  return double.findById(args[0]);
}

/** Returns the sample of the result of findById. */
function finderStubFindByIdPin(): string {
  let want0: string;
  want0 = 'findById-a';
  return want0;
}

describe("FinderStub.findById", () => {
  const checks = behaviour("Finder.findById", (seat) => {
    const double = new FinderStub(seat);
    return {
      stub: double.onFindById,
      call: () => finderStubFindById(double),
      override: (mark) => {
        double.onFindById.func(() => {
          mark();
          return undefined as never;
        });
      },
    };
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new FinderStub(seat);
    const want = finderStubFindByIdPin();
    double.onFindById.returns(want);
    check.equal(seat, finderStubFindById(double), want, "findById returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new FinderStub(seat);
    finderStubFindById(double);
    const got = double.onFindById.assertCalledOnce(seat, "findById records the call");
    const args = finderStubFindByIdArgs();
    check.equal(seat, got?.id, args[0], "the call records id");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new FinderStub(seat);
    let seen = 0;
    double.onFindById.onRecord(() => {
      seen += 1;
    });
    finderStubFindById(double);
    finderStubFindById(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of findById");
  });

  test("sets the function of withFinderFindById at construction", async ({ seat }) => {
    let called = false;
    const option = FinderStub.withFinderFindById(() => {
      called = true;
      return undefined as never;
    });
    finderStubFindById(new FinderStub(seat, option));
    check.isTrue(seat, called, "withFinderFindById sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("findById-failing");
    const double = FinderStub.newFailing(seat, want);
    const outcome = await settle(() => finderStubFindById(double));
    const message = "findById fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new FinderStub(seat);
    const want = finderStubFindByIdPin();
    double.onFindById.returns(want);
    finderStubFindById(double);
    double.resetCalls();
    check.equal(seat, finderStubFindById(double), want, "a reset keeps the result of returns");
  });
});

describe("FinderStub", () => {
  const instance = (double: FinderStub) => ({
    stub: double.onFindById,
    call: () => finderStubFindById(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new FinderStub(seat)),
    withClock: (seat, clock) => instance(new FinderStub(seat, FinderStub.withClock(clock))),
    withLatency: (seat, duration, clock) => instance(FinderStub.newDelayed(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new FinderStub(seat, FinderStub.withRandSource(random))),
    benchMode: (seat) => instance(new FinderStub(seat, FinderStub.benchMode())),
    strict: (seat) => instance(new FinderStub(seat, FinderStub.strict())),
    async: false,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes findById to the implementation of delegateTo", async ({ seat }) => {
    const inner = new FinderStub(seat);
    await settle(() => finderStubFindById(new FinderStub(seat, FinderStub.delegateTo(inner))));
    inner.onFindById.assertCalledOnce(seat, "findById reaches the implementation");
  });

  test("fails findById with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("findById-delegate");
    const inner = new FinderStub(seat);
    inner.onFindById.faults(want, 1);
    const double = new FinderStub(seat, FinderStub.delegateTo(inner));
    const outcome = await settle(() => finderStubFindById(double));
    const message = "findById fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});
