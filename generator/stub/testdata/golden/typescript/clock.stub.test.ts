import { type Seat, Standard, byIdentity, check } from "@dokimi/assert";
import { type Subject, behaviour, resetBehaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { type ClockEveryCall, type ClockNowCall, type ClockStopperCall, ClockStub } from "./clock.stub";

/** Calls now of double with the samples of its arguments. */
function clockStubNow(double: ClockStub): number {
  return double.now();
}

/** Returns the sample of the result of now. */
function clockStubNowPin(): number {
  let want0: number;
  want0 = 1.5;
  return want0;
}

/** Returns the subject of behaviour and resetBehaviour for now, on a new double that reports to seat. */
function clockStubNowSubject(seat: Seat): Subject<ClockNowCall, number | undefined> {
  const double = new ClockStub(seat);
  return {
    stub: double.onNow,
    call: () => clockStubNow(double),
    reset: () => double.resetCalls(),
    result: async () => clockStubNow(double),
    returns: () => {
      const want = clockStubNowPin();
      double.onNow.returns(want);
      return want;
    },
    override: (seen) => {
      double.onNow.func(() => {
        seen({} as ClockNowCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => clockStubNow(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Returns the samples of the arguments of every. An argument without a sample is undefined. */
function clockStubEveryArgs(): [() => void] {
  let a0 = undefined as unknown as () => void;
  return [a0];
}

/** Calls every of double with the samples of its arguments. */
function clockStubEvery(double: ClockStub): void {
  const args = clockStubEveryArgs();
  return double.every(args[0]);
}

/** Returns the subject of behaviour and resetBehaviour for every, on a new double that reports to seat. */
function clockStubEverySubject(seat: Seat): Subject<ClockEveryCall, undefined> {
  const double = new ClockStub(seat);
  return {
    stub: double.onEvery,
    call: () => clockStubEvery(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = clockStubEveryArgs();
      check.isEqualTo(s, c.fn, args[0], "the call records fn");
    },
    override: (seen) => {
      double.onEvery.func((...received: unknown[]) => {
        seen({ fn: received[0] } as ClockEveryCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => clockStubEvery(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Calls stopper of double with the samples of its arguments. */
function clockStubStopper(double: ClockStub): () => void {
  return double.stopper();
}

/** Returns the subject of behaviour and resetBehaviour for stopper, on a new double that reports to seat. */
function clockStubStopperSubject(seat: Seat): Subject<ClockStopperCall, (() => void) | undefined> {
  const double = new ClockStub(seat);
  return {
    stub: double.onStopper,
    call: () => clockStubStopper(double),
    reset: () => double.resetCalls(),
    result: async () => clockStubStopper(double),
    override: (seen) => {
      double.onStopper.func(() => {
        seen({} as ClockStopperCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => clockStubStopper(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call now alone.
 */
function clockStubInstance(double: ClockStub) {
  return {
    stub: double.onNow,
    call: () => clockStubNow(double),
  };
}

describe("ClockStub", () => {
  describe("now", () => {
    for (const c of behaviour("Clock.now", clockStubNowSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("every", () => {
    for (const c of behaviour("Clock.every", clockStubEverySubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("stopper", () => {
    for (const c of behaviour("Clock.stopper", clockStubStopperSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("now", () => {
      for (const c of resetBehaviour("Clock.now", clockStubNowSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("every", () => {
      for (const c of resetBehaviour("Clock.every", clockStubEverySubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("stopper", () => {
      for (const c of resetBehaviour("Clock.stopper", clockStubStopperSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => clockStubInstance(new ClockStub(seat, ClockStub.strict())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes now to the implementation", async () => {
      const inner = new ClockStub(undefined);
      await settle(() => clockStubNow(new ClockStub(undefined, ClockStub.delegateTo(inner))));
      expect(inner.onNow.callCount).isEqualTo(1, "now calls the implementation");
    });

    it("fails now with the error of the implementation", async () => {
      const want = new Error("now-delegate");
      const inner = new ClockStub(undefined);
      inner.onNow.faults(want, 1);
      const double = new ClockStub(undefined, ClockStub.delegateTo(inner));
      const outcome = await settle(() => clockStubNow(double));
      const message = "now fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes every to the implementation", async () => {
      const inner = new ClockStub(undefined);
      await settle(() => clockStubEvery(new ClockStub(undefined, ClockStub.delegateTo(inner))));
      expect(inner.onEvery.callCount).isEqualTo(1, "every calls the implementation");
    });

    it("fails every with the error of the implementation", async () => {
      const want = new Error("every-delegate");
      const inner = new ClockStub(undefined);
      inner.onEvery.faults(want, 1);
      const double = new ClockStub(undefined, ClockStub.delegateTo(inner));
      const outcome = await settle(() => clockStubEvery(double));
      const message = "every fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes stopper to the implementation", async () => {
      const inner = new ClockStub(undefined);
      await settle(() => clockStubStopper(new ClockStub(undefined, ClockStub.delegateTo(inner))));
      expect(inner.onStopper.callCount).isEqualTo(1, "stopper calls the implementation");
    });

    it("fails stopper with the error of the implementation", async () => {
      const want = new Error("stopper-delegate");
      const inner = new ClockStub(undefined);
      inner.onStopper.faults(want, 1);
      const double = new ClockStub(undefined, ClockStub.delegateTo(inner));
      const outcome = await settle(() => clockStubStopper(double));
      const message = "stopper fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => clockStubInstance(new ClockStub(seat, ClockStub.withClock(clock))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      clockStubInstance(new ClockStub(seat, ClockStub.withRandSource(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => clockStubInstance(new ClockStub(seat, ClockStub.benchMode())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        clockStubInstance(new ClockStub(seat, ClockStub.withClock(clock), ClockStub.withLatency(duration))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails now with its error", async () => {
      const want = new Error("now-failing");
      const double = new ClockStub(undefined, ClockStub.failing(want));
      const outcome = await settle(() => clockStubNow(double));
      const message = "now fails with the error of ClockStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails every with its error", async () => {
      const want = new Error("every-failing");
      const double = new ClockStub(undefined, ClockStub.failing(want));
      const outcome = await settle(() => clockStubEvery(double));
      const message = "every fails with the error of ClockStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails stopper with its error", async () => {
      const want = new Error("stopper-failing");
      const double = new ClockStub(undefined, ClockStub.failing(want));
      const outcome = await settle(() => clockStubStopper(double));
      const message = "stopper fails with the error of ClockStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClockNow", () => {
    it("sets the function of now", async () => {
      let called = false;
      const option = ClockStub.withClockNow(() => {
        called = true;
        return undefined as never;
      });
      clockStubNow(new ClockStub(undefined, option));
      expect(called).isTrue("withClockNow sets the function of now");
    });
  });

  describe("withClockEvery", () => {
    it("sets the function of every", async () => {
      let called = false;
      const option = ClockStub.withClockEvery(() => {
        called = true;
        return undefined as never;
      });
      clockStubEvery(new ClockStub(undefined, option));
      expect(called).isTrue("withClockEvery sets the function of every");
    });
  });

  describe("withClockStopper", () => {
    it("sets the function of stopper", async () => {
      let called = false;
      const option = ClockStub.withClockStopper(() => {
        called = true;
        return undefined as never;
      });
      clockStubStopper(new ClockStub(undefined, option));
      expect(called).isTrue("withClockStopper sets the function of stopper");
    });
  });
});
