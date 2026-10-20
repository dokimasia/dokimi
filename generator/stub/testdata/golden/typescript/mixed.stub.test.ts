import { Recorder, byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { behaviour, settle } from '@dokimi/runtime/stub/behaviour';
import { doubleBehaviour } from '@dokimi/runtime/stub/double-behaviour';
import { describe } from 'vitest';
import { SessionStub } from './mixed.stub';

/** Calls prepare of double with the samples of its arguments. */
function sessionStubPrepare(double: SessionStub): Promise<void> {
  return double.prepare();
}

/** Calls commit of double with the samples of its arguments. */
function sessionStubCommit(double: SessionStub): Promise<void> {
  return double.commit();
}

/** Returns the samples of the arguments of old, or undefined. */
function sessionStubOldArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = 'old-a';
  return [a0];
}

/** Calls old of double with the samples of its arguments. */
function sessionStubOld(double: SessionStub): string {
  const args = sessionStubOldArgs();
  return double.old(args[0]);
}

/** Returns the sample of the result of old. */
function sessionStubOldPin(): string {
  let want0: string;
  want0 = 'old-a';
  return want0;
}

describe("SessionStub.prepare", () => {
  const checks = behaviour("Session.prepare", (seat) => {
    const double = new SessionStub(seat);
    return {
      stub: double.onPrepare,
      call: () => sessionStubPrepare(double),
      override: (mark) => {
        double.onPrepare.func(() => {
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
    const double = new SessionStub(seat);
    await sessionStubPrepare(double);
    double.onPrepare.assertCalledOnce(seat, "prepare records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new SessionStub(seat);
    let seen = 0;
    double.onPrepare.onRecord(() => {
      seen += 1;
    });
    await sessionStubPrepare(double);
    await sessionStubPrepare(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of prepare");
  });

  test("sets the function of withSessionPrepare at construction", async ({ seat }) => {
    let called = false;
    const option = SessionStub.withSessionPrepare(() => {
      called = true;
      return undefined as never;
    });
    await sessionStubPrepare(new SessionStub(seat, option));
    check.isTrue(seat, called, "withSessionPrepare sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("prepare-failing");
    const double = SessionStub.newFailing(seat, want);
    const outcome = await settle(() => sessionStubPrepare(double));
    const message = "prepare fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});

describe("SessionStub.commit", () => {
  const checks = behaviour("Session.commit", (seat) => {
    const double = new SessionStub(seat);
    return {
      stub: double.onCommit,
      call: () => sessionStubCommit(double),
      override: (mark) => {
        double.onCommit.func(() => {
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
    const double = new SessionStub(seat);
    await sessionStubCommit(double);
    double.onCommit.assertCalledOnce(seat, "commit records the call");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new SessionStub(seat);
    let seen = 0;
    double.onCommit.onRecord(() => {
      seen += 1;
    });
    await sessionStubCommit(double);
    await sessionStubCommit(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of commit");
  });

  test("sets the function of withSessionCommit at construction", async ({ seat }) => {
    let called = false;
    const option = SessionStub.withSessionCommit(() => {
      called = true;
      return undefined as never;
    });
    await sessionStubCommit(new SessionStub(seat, option));
    check.isTrue(seat, called, "withSessionCommit sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("commit-failing");
    const double = SessionStub.newFailing(seat, want);
    const outcome = await settle(() => sessionStubCommit(double));
    const message = "commit fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("refuses commit before prepare in strict mode", async ({ seat }) => {
    const failures = new Recorder();
    const double = new SessionStub(failures, SessionStub.strict());
    await settle(() => sessionStubCommit(double));
    check.isTrue(seat, failures.failed, "commit before prepare is refused");
  });
});

describe("SessionStub.old", () => {
  const checks = behaviour("Session.old", (seat) => {
    const double = new SessionStub(seat);
    return {
      stub: double.onOld,
      call: () => sessionStubOld(double),
      override: (mark) => {
        double.onOld.func(() => {
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
    const double = new SessionStub(seat);
    const want = sessionStubOldPin();
    double.onOld.returns(want);
    check.equal(seat, sessionStubOld(double), want, "old returns the result of returns");
  });

  test("records the arguments of the call", async ({ seat }) => {
    const double = new SessionStub(seat);
    sessionStubOld(double);
    const got = double.onOld.assertCalledOnce(seat, "old records the call");
    const args = sessionStubOldArgs();
    check.equal(seat, got?.key, args[0], "the call records key");
  });

  test("runs the onRecord hook for each call", async ({ seat }) => {
    const double = new SessionStub(seat);
    let seen = 0;
    double.onOld.onRecord(() => {
      seen += 1;
    });
    sessionStubOld(double);
    sessionStubOld(double);
    check.equal(seat, seen, 2, "onRecord runs once for each call of old");
  });

  test("sets the function of withSessionOld at construction", async ({ seat }) => {
    let called = false;
    const option = SessionStub.withSessionOld(() => {
      called = true;
      return undefined as never;
    });
    sessionStubOld(new SessionStub(seat, option));
    check.isTrue(seat, called, "withSessionOld sets the function");
  });

  test("fails with the error of the failing option", async ({ seat }) => {
    const want = new Error("old-failing");
    const double = SessionStub.newFailing(seat, want);
    const outcome = await settle(() => sessionStubOld(double));
    const message = "old fails with the error of the option";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("keeps the result of returns across a reset", async ({ seat }) => {
    const double = new SessionStub(seat);
    const want = sessionStubOldPin();
    double.onOld.returns(want);
    sessionStubOld(double);
    double.resetCalls();
    check.equal(seat, sessionStubOld(double), want, "a reset keeps the result of returns");
  });
});

describe("SessionStub", () => {
  const instance = (double: SessionStub) => ({
    stub: double.onPrepare,
    call: () => sessionStubPrepare(double),
    reset: () => double.resetCalls(),
  });
  const checks = doubleBehaviour({
    create: (seat) => instance(new SessionStub(seat)),
    withClock: (seat, clock) => instance(new SessionStub(seat, SessionStub.withClock(clock))),
    withLatency: (seat, duration, clock) => instance(SessionStub.newDelayed(seat, duration, clock)),
    withRandSource: (seat, random) =>
      instance(new SessionStub(seat, SessionStub.withRandSource(random))),
    benchMode: (seat) => instance(new SessionStub(seat, SessionStub.benchMode())),
    strict: (seat) => instance(new SessionStub(seat, SessionStub.strict())),
    async: true,
  });
  for (const c of checks) {
    test(c.name, ({ seat }) => c.run(seat));
  }

  test("passes prepare to the implementation of delegateTo", async ({ seat }) => {
    const inner = new SessionStub(seat);
    await settle(() => sessionStubPrepare(new SessionStub(seat, SessionStub.delegateTo(inner))));
    inner.onPrepare.assertCalledOnce(seat, "prepare reaches the implementation");
  });

  test("fails prepare with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("prepare-delegate");
    const inner = new SessionStub(seat);
    inner.onPrepare.faults(want, 1);
    const double = new SessionStub(seat, SessionStub.delegateTo(inner));
    const outcome = await settle(() => sessionStubPrepare(double));
    const message = "prepare fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes commit to the implementation of delegateTo", async ({ seat }) => {
    const inner = new SessionStub(seat);
    await settle(() => sessionStubCommit(new SessionStub(seat, SessionStub.delegateTo(inner))));
    inner.onCommit.assertCalledOnce(seat, "commit reaches the implementation");
  });

  test("fails commit with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("commit-delegate");
    const inner = new SessionStub(seat);
    inner.onCommit.faults(want, 1);
    const double = new SessionStub(seat, SessionStub.delegateTo(inner));
    const outcome = await settle(() => sessionStubCommit(double));
    const message = "commit fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });

  test("passes old to the implementation of delegateTo", async ({ seat }) => {
    const inner = new SessionStub(seat);
    await settle(() => sessionStubOld(new SessionStub(seat, SessionStub.delegateTo(inner))));
    inner.onOld.assertCalledOnce(seat, "old reaches the implementation");
  });

  test("fails old with the error of the implementation of delegateTo", async ({ seat }) => {
    const want = new Error("old-delegate");
    const inner = new SessionStub(seat);
    inner.onOld.faults(want, 1);
    const double = new SessionStub(seat, SessionStub.delegateTo(inner));
    const outcome = await settle(() => sessionStubOld(double));
    const message = "old fails with the error of the implementation";
    check.equal(seat, outcome.threw && outcome.error, want, message, byIdentity());
  });
});
