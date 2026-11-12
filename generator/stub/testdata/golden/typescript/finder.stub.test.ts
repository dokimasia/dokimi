import { type Seat, Standard, byIdentity, check } from "@dokimi/assert";
import { type Subject, behaviour, resetBehaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { type FinderFindByIdCall, FinderStub } from "./finder.stub";

/** Returns the samples of the arguments of findById. An argument without a sample is undefined. */
function finderStubFindByIdArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "findById-a";
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
  want0 = "findById-a";
  return want0;
}

/** Returns the subject of behaviour and resetBehaviour for findById, on a new double that reports to seat. */
function finderStubFindByIdSubject(seat: Seat): Subject<FinderFindByIdCall, string | undefined> {
  const double = new FinderStub(seat);
  return {
    stub: double.onFindById,
    call: () => finderStubFindById(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = finderStubFindByIdArgs();
      check.isEqualTo(s, c.id, args[0], "the call records id");
    },
    result: async () => finderStubFindById(double),
    returns: () => {
      const want = finderStubFindByIdPin();
      double.onFindById.returns(want);
      return want;
    },
    override: (seen) => {
      double.onFindById.func((...received: unknown[]) => {
        seen({ id: received[0] } as FinderFindByIdCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => finderStubFindById(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call findById alone.
 */
function finderStubInstance(double: FinderStub) {
  return {
    stub: double.onFindById,
    call: () => finderStubFindById(double),
  };
}

describe("FinderStub", () => {
  describe("findById", () => {
    for (const c of behaviour("Finder.findById", finderStubFindByIdSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("findById", () => {
      for (const c of resetBehaviour("Finder.findById", finderStubFindByIdSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => finderStubInstance(new FinderStub(seat, FinderStub.strict())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes findById to the implementation", async () => {
      const inner = new FinderStub(undefined);
      await settle(() => finderStubFindById(new FinderStub(undefined, FinderStub.delegateTo(inner))));
      expect(inner.onFindById.callCount).isEqualTo(1, "findById calls the implementation");
    });

    it("fails findById with the error of the implementation", async () => {
      const want = new Error("findById-delegate");
      const inner = new FinderStub(undefined);
      inner.onFindById.faults(want, 1);
      const double = new FinderStub(undefined, FinderStub.delegateTo(inner));
      const outcome = await settle(() => finderStubFindById(double));
      const message = "findById fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => finderStubInstance(new FinderStub(seat, FinderStub.withClock(clock))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      finderStubInstance(new FinderStub(seat, FinderStub.withRandSource(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => finderStubInstance(new FinderStub(seat, FinderStub.benchMode())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        finderStubInstance(new FinderStub(seat, FinderStub.withClock(clock), FinderStub.withLatency(duration))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails findById with its error", async () => {
      const want = new Error("findById-failing");
      const double = new FinderStub(undefined, FinderStub.failing(want));
      const outcome = await settle(() => finderStubFindById(double));
      const message = "findById fails with the error of FinderStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withFinderFindById", () => {
    it("sets the function of findById", async () => {
      let called = false;
      const option = FinderStub.withFinderFindById(() => {
        called = true;
        return undefined as never;
      });
      finderStubFindById(new FinderStub(undefined, option));
      expect(called).isTrue("withFinderFindById sets the function of findById");
    });
  });
});
