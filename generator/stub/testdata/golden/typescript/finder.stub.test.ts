import { Recorder, Standard, byIdentity } from "@dokimi/assert";
import { behaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, resetBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { FinderStub } from "./finder.stub";

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

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * findById alone.
 */
function finderStubInstance(double: FinderStub) {
  return {
    stub: double.onFindById,
    call: () => finderStubFindById(double),
    reset: () => double.resetCalls(),
  };
}

describe("FinderStub", () => {
  describe("findById", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new FinderStub(undefined);
      finderStubFindById(double);
      expect(double.onFindById.callCount).isEqualTo(1, "findById records the call");
      const got = double.onFindById.calls[0];
      const args = finderStubFindByIdArgs();
      expect(got?.id).isEqualTo(args[0], "the call records id");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new FinderStub(undefined);
      let seen = 0;
      double.onFindById.onRecord(() => {
        seen += 1;
      });
      finderStubFindById(double);
      finderStubFindById(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of findById");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => finderStubInstance(new FinderStub(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of findById", async () => {
      const double = new FinderStub(undefined);
      finderStubFindById(double);
      double.resetCalls();
      expect(double.onFindById.callCount).isEqualTo(0, "resetCalls clears the calls of findById");
    });

    it("keeps the result of returns of findById", async () => {
      const double = new FinderStub(undefined);
      const want = finderStubFindByIdPin();
      double.onFindById.returns(want);
      double.resetCalls();
      expect(finderStubFindById(double)).isEqualTo(want, "resetCalls keeps the result of findById");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => finderStubInstance(FinderStub.newDelayed(seat, duration, clock)),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails findById with its error", async () => {
      const want = new Error("findById-failing");
      const double = FinderStub.newFailing(undefined, want);
      const outcome = await settle(() => finderStubFindById(double));
      const message = "findById fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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

describe("FinderFindByIdStub", () => {
  describe("returns", () => {
    it("sets the result of each call of findById", async () => {
      const double = new FinderStub(undefined);
      const want = finderStubFindByIdPin();
      double.onFindById.returns(want);
      expect(finderStubFindById(double)).isEqualTo(want, "findById returns the result that returns sets");
    });

    it("configures findById for strict mode", async () => {
      const failures = new Recorder();
      const double = new FinderStub(failures);
      double.onFindById.strict();
      double.onFindById.returns(undefined as never);
      await settle(() => finderStubFindById(double));
      expect(failures.failed).isFalse("strict mode accepts the call of findById");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of findById to the function", async () => {
      const double = new FinderStub(undefined);
      let got: unknown[] = [];
      double.onFindById.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      finderStubFindById(double);
      const args = finderStubFindByIdArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });
});
