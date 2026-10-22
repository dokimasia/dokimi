import { Recorder, Standard, byIdentity } from "@dokimi/assert";
import { behaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, resetBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { ClockStub } from "./clock.stub";

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

/** Calls stopper of double with the samples of its arguments. */
function clockStubStopper(double: ClockStub): () => void {
  return double.stopper();
}

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * now alone.
 */
function clockStubInstance(double: ClockStub) {
  return {
    stub: double.onNow,
    call: () => clockStubNow(double),
    reset: () => double.resetCalls(),
  };
}

describe("ClockStub", () => {
  describe("now", () => {
    const checks = behaviour("Clock.now", (seat) => {
      const double = new ClockStub(seat);
      return {
        stub: double.onNow,
        call: () => clockStubNow(double),
        override: (mark) => {
          double.onNow.func(() => {
            mark();
            return undefined as never;
          });
        },
      };
    });
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("runs the onRecord hook once for each call", async () => {
      const double = new ClockStub(undefined);
      let seen = 0;
      double.onNow.onRecord(() => {
        seen += 1;
      });
      clockStubNow(double);
      clockStubNow(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of now");
    });
  });

  describe("every", () => {
    const checks = behaviour("Clock.every", (seat) => {
      const double = new ClockStub(seat);
      return {
        stub: double.onEvery,
        call: () => clockStubEvery(double),
        override: (mark) => {
          double.onEvery.func(() => {
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
      const double = new ClockStub(undefined);
      clockStubEvery(double);
      expect(double.onEvery.callCount).isEqualTo(1, "every records the call");
      const got = double.onEvery.calls[0];
      const args = clockStubEveryArgs();
      expect(got?.fn).isEqualTo(args[0], "the call records fn");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new ClockStub(undefined);
      let seen = 0;
      double.onEvery.onRecord(() => {
        seen += 1;
      });
      clockStubEvery(double);
      clockStubEvery(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of every");
    });
  });

  describe("stopper", () => {
    const checks = behaviour("Clock.stopper", (seat) => {
      const double = new ClockStub(seat);
      return {
        stub: double.onStopper,
        call: () => clockStubStopper(double),
        override: (mark) => {
          double.onStopper.func(() => {
            mark();
            return undefined as never;
          });
        },
      };
    });
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("runs the onRecord hook once for each call", async () => {
      const double = new ClockStub(undefined);
      let seen = 0;
      double.onStopper.onRecord(() => {
        seen += 1;
      });
      clockStubStopper(double);
      clockStubStopper(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of stopper");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => clockStubInstance(new ClockStub(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of now", async () => {
      const double = new ClockStub(undefined);
      clockStubNow(double);
      double.resetCalls();
      expect(double.onNow.callCount).isEqualTo(0, "resetCalls clears the calls of now");
    });

    it("keeps the result of returns of now", async () => {
      const double = new ClockStub(undefined);
      const want = clockStubNowPin();
      double.onNow.returns(want);
      double.resetCalls();
      expect(clockStubNow(double)).isEqualTo(want, "resetCalls keeps the result of now");
    });

    it("clears the recorded calls of every", async () => {
      const double = new ClockStub(undefined);
      clockStubEvery(double);
      double.resetCalls();
      expect(double.onEvery.callCount).isEqualTo(0, "resetCalls clears the calls of every");
    });

    it("clears the recorded calls of stopper", async () => {
      const double = new ClockStub(undefined);
      clockStubStopper(double);
      double.resetCalls();
      expect(double.onStopper.callCount).isEqualTo(0, "resetCalls clears the calls of stopper");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => clockStubInstance(ClockStub.newDelayed(seat, duration, clock)),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails now with its error", async () => {
      const want = new Error("now-failing");
      const double = ClockStub.newFailing(undefined, want);
      const outcome = await settle(() => clockStubNow(double));
      const message = "now fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails every with its error", async () => {
      const want = new Error("every-failing");
      const double = ClockStub.newFailing(undefined, want);
      const outcome = await settle(() => clockStubEvery(double));
      const message = "every fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails stopper with its error", async () => {
      const want = new Error("stopper-failing");
      const double = ClockStub.newFailing(undefined, want);
      const outcome = await settle(() => clockStubStopper(double));
      const message = "stopper fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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

describe("ClockNowStub", () => {
  describe("returns", () => {
    it("sets the result of each call of now", async () => {
      const double = new ClockStub(undefined);
      const want = clockStubNowPin();
      double.onNow.returns(want);
      expect(clockStubNow(double)).isEqualTo(want, "now returns the result that returns sets");
    });

    it("configures now for strict mode", async () => {
      const failures = new Recorder();
      const double = new ClockStub(failures);
      double.onNow.strict();
      double.onNow.returns(undefined as never);
      await settle(() => clockStubNow(double));
      expect(failures.failed).isFalse("strict mode accepts the call of now");
    });
  });
});

describe("ClockEveryStub", () => {
  describe("func", () => {
    it("passes the arguments of each call of every to the function", async () => {
      const double = new ClockStub(undefined);
      let got: unknown[] = [];
      double.onEvery.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      clockStubEvery(double);
      const args = clockStubEveryArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });
});

describe("ClockStopperStub", () => {
  describe("returns", () => {
    it("configures stopper for strict mode", async () => {
      const failures = new Recorder();
      const double = new ClockStub(failures);
      double.onStopper.strict();
      double.onStopper.returns(undefined as never);
      await settle(() => clockStubStopper(double));
      expect(failures.failed).isFalse("strict mode accepts the call of stopper");
    });
  });
});
