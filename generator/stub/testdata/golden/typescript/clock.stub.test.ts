import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { behaviour, settle } from '@dokimi/runtime/stub/behaviour';
import { doubleBehaviour } from '@dokimi/runtime/stub/double-behaviour';
import { describe } from 'vitest';
import { ClockStub } from './clock.stub';

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

/** Returns the samples of the arguments of every, or undefined. */
function clockStubEveryArgs(): [()=>void] {
  let a0 = undefined as unknown as ()=>void;
  return [a0];
}

/** Calls every of double with the samples of its arguments. */
function clockStubEvery(double: ClockStub): void {
  const args = clockStubEveryArgs();
  return double.every(args[0]);
}

/** Calls stopper of double with the samples of its arguments. */
function clockStubStopper(double: ClockStub): ()=>void {
  return double.stopper();
}

describe("ClockStub.now", () => {
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
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("returns the result of returns", async ({ seat }) => {
    const double = new ClockStub(seat);
    const want = clockStubNowPin();
    double.onNow.returns(want);
    check.equal(seat, clockStubNow(double), want, "now returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new ClockStub(seat);
    clockStubNow(double);
    double.onNow.assertCalledOnce(seat, "now records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new ClockStub(seat);
    let seen = 0;
    double.onNow.onRecord(() => {
      seen += 1;
    });
    clockStubNow(double);
    clockStubNow(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of now");
  });

  test("sets the function of withClockNow at construction", async ({ seat }) => {
    let called = false;
    const option = ClockStub.withClockNow(() => {
      called = true;
      return undefined as never;
    });
    clockStubNow(new ClockStub(seat, option));
    check.isTrue(seat, called, "withClockNow sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("now-failing");
    const double = ClockStub.newFailing(seat, want);
    const outcome = await settle(() => clockStubNow(double));
    const message = "now fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new ClockStub(seat);
    const want = clockStubNowPin();
    double.onNow.returns(want);
    clockStubNow(double);
    double.resetCalls();
    check.equal(seat, clockStubNow(double), want, "a reset keeps the result of returns");
  });
});

describe("ClockStub.every", () => {
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
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new ClockStub(seat);
    clockStubEvery(double);
    const got = double.onEvery.assertCalledOnce(seat, "every records the call");
    const args = clockStubEveryArgs();
    check.equal(seat, got?.fn, args[0], "the call records fn");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new ClockStub(seat);
    let seen = 0;
    double.onEvery.onRecord(() => {
      seen += 1;
    });
    clockStubEvery(double);
    clockStubEvery(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of every");
  });

  test("sets the function of withClockEvery at construction", async ({ seat }) => {
    let called = false;
    const option = ClockStub.withClockEvery(() => {
      called = true;
      return undefined as never;
    });
    clockStubEvery(new ClockStub(seat, option));
    check.isTrue(seat, called, "withClockEvery sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("every-failing");
    const double = ClockStub.newFailing(seat, want);
    const outcome = await settle(() => clockStubEvery(double));
    const message = "every fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

describe("ClockStub.stopper", () => {
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
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("records the arguments of the call", async ({ seat }) => {
    const double = new ClockStub(seat);
    clockStubStopper(double);
    double.onStopper.assertCalledOnce(seat, "stopper records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new ClockStub(seat);
    let seen = 0;
    double.onStopper.onRecord(() => {
      seen += 1;
    });
    clockStubStopper(double);
    clockStubStopper(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of stopper");
  });

  test("sets the function of withClockStopper at construction", async ({ seat }) => {
    let called = false;
    const option = ClockStub.withClockStopper(() => {
      called = true;
      return undefined as never;
    });
    clockStubStopper(new ClockStub(seat, option));
    check.isTrue(seat, called, "withClockStopper sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("stopper-failing");
    const double = ClockStub.newFailing(seat, want);
    const outcome = await settle(() => clockStubStopper(double));
    const message = "stopper fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps recording across a reset", async ({ seat }) => {
    const double = new ClockStub(seat);
    double.onStopper.returns(undefined as never);
    clockStubStopper(double);
    double.resetCalls();
    clockStubStopper(double);
    double.onStopper.assertCalledOnce(seat, "stopper records the call after the reset");
  });
});

describe("ClockStub", () => {
  const instance = (double: ClockStub) => ({
    stub: double.onNow,
    call: () => clockStubNow(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new ClockStub(seat)),
    withClock: (seat, clock) => instance(new ClockStub(seat, ClockStub.withClock(clock))),
    withLatency: (seat, duration, clock) => instance(ClockStub.newDelayed(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new ClockStub(seat, ClockStub.withRandSource(random))),
    benchMode: (seat) => instance(new ClockStub(seat, ClockStub.benchMode())),
    strict: (seat) => instance(new ClockStub(seat, ClockStub.strict())),
    async: false,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes now to the implementation of delegateTo", async ({ seat }) => {
    const inner = new ClockStub(seat);
    await settle(() => clockStubNow(new ClockStub(seat, ClockStub.delegateTo(inner))));
    inner.onNow.assertCalledOnce(seat, "now reaches the implementation");
  });

  test("fails now with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("now-delegate");
    const inner = new ClockStub(seat);
    inner.onNow.faults(want, 1);
    const double = new ClockStub(seat, ClockStub.delegateTo(inner));
    const outcome = await settle(() => clockStubNow(double));
    const message = "now fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes every to the implementation of delegateTo", async ({ seat }) => {
    const inner = new ClockStub(seat);
    await settle(() => clockStubEvery(new ClockStub(seat, ClockStub.delegateTo(inner))));
    inner.onEvery.assertCalledOnce(seat, "every reaches the implementation");
  });

  test("fails every with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("every-delegate");
    const inner = new ClockStub(seat);
    inner.onEvery.faults(want, 1);
    const double = new ClockStub(seat, ClockStub.delegateTo(inner));
    const outcome = await settle(() => clockStubEvery(double));
    const message = "every fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes stopper to the implementation of delegateTo", async ({ seat }) => {
    const inner = new ClockStub(seat);
    await settle(() => clockStubStopper(new ClockStub(seat, ClockStub.delegateTo(inner))));
    inner.onStopper.assertCalledOnce(seat, "stopper reaches the implementation");
  });

  test("fails stopper with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("stopper-delegate");
    const inner = new ClockStub(seat);
    inner.onStopper.faults(want, 1);
    const double = new ClockStub(seat, ClockStub.delegateTo(inner));
    const outcome = await settle(() => clockStubStopper(double));
    const message = "stopper fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});
