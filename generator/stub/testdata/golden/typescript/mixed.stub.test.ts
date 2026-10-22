import { Recorder, Standard, byIdentity } from "@dokimi/assert";
import { behaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, resetBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { SessionStub } from "./mixed.stub";

/** Calls prepare of double with the samples of its arguments. */
function sessionStubPrepare(double: SessionStub): Promise<void> {
  return double.prepare();
}

/** Calls commit of double with the samples of its arguments. */
function sessionStubCommit(double: SessionStub): Promise<void> {
  return double.commit();
}

/** Returns the samples of the arguments of old. An argument without a sample is undefined. */
function sessionStubOldArgs(): [string] {
  let a0 = undefined as unknown as string;
  a0 = "old-a";
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
  want0 = "old-a";
  return want0;
}

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * prepare alone.
 */
function sessionStubInstance(double: SessionStub) {
  return {
    stub: double.onPrepare,
    call: () => sessionStubPrepare(double),
    reset: () => double.resetCalls(),
  };
}

describe("SessionStub", () => {
  describe("prepare", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("runs the onRecord hook once for each call", async () => {
      const double = new SessionStub(undefined);
      let seen = 0;
      double.onPrepare.onRecord(() => {
        seen += 1;
      });
      await sessionStubPrepare(double);
      await sessionStubPrepare(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of prepare");
    });
  });

  describe("commit", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("runs the onRecord hook once for each call", async () => {
      const double = new SessionStub(undefined);
      let seen = 0;
      double.onCommit.onRecord(() => {
        seen += 1;
      });
      await sessionStubCommit(double);
      await sessionStubCommit(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of commit");
    });
  });

  describe("old", () => {
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
      it(c.name, () => c.run(new Standard()));
    }

    it("records the arguments of the call", async () => {
      const double = new SessionStub(undefined);
      sessionStubOld(double);
      expect(double.onOld.callCount).isEqualTo(1, "old records the call");
      const got = double.onOld.calls[0];
      const args = sessionStubOldArgs();
      expect(got?.key).isEqualTo(args[0], "the call records key");
    });

    it("runs the onRecord hook once for each call", async () => {
      const double = new SessionStub(undefined);
      let seen = 0;
      double.onOld.onRecord(() => {
        seen += 1;
      });
      sessionStubOld(double);
      sessionStubOld(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of old");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => sessionStubInstance(new SessionStub(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of prepare", async () => {
      const double = new SessionStub(undefined);
      await sessionStubPrepare(double);
      double.resetCalls();
      expect(double.onPrepare.callCount).isEqualTo(0, "resetCalls clears the calls of prepare");
    });

    it("clears the recorded calls of commit", async () => {
      const double = new SessionStub(undefined);
      await sessionStubCommit(double);
      double.resetCalls();
      expect(double.onCommit.callCount).isEqualTo(0, "resetCalls clears the calls of commit");
    });

    it("clears the recorded calls of old", async () => {
      const double = new SessionStub(undefined);
      sessionStubOld(double);
      double.resetCalls();
      expect(double.onOld.callCount).isEqualTo(0, "resetCalls clears the calls of old");
    });

    it("keeps the result of returns of old", async () => {
      const double = new SessionStub(undefined);
      const want = sessionStubOldPin();
      double.onOld.returns(want);
      double.resetCalls();
      expect(sessionStubOld(double)).isEqualTo(want, "resetCalls keeps the result of old");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => sessionStubInstance(SessionStub.newDelayed(seat, duration, clock)),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails prepare with its error", async () => {
      const want = new Error("prepare-failing");
      const double = SessionStub.newFailing(undefined, want);
      const outcome = await settle(() => sessionStubPrepare(double));
      const message = "prepare fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails commit with its error", async () => {
      const want = new Error("commit-failing");
      const double = SessionStub.newFailing(undefined, want);
      const outcome = await settle(() => sessionStubCommit(double));
      const message = "commit fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails old with its error", async () => {
      const want = new Error("old-failing");
      const double = SessionStub.newFailing(undefined, want);
      const outcome = await settle(() => sessionStubOld(double));
      const message = "old fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => sessionStubInstance(new SessionStub(seat, SessionStub.strict())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("refuses commit before prepare", async () => {
      const failures = new Recorder();
      const option = SessionStub.withSessionCommit(() => undefined as never);
      const double = new SessionStub(failures, SessionStub.strict(), option);
      await settle(() => sessionStubCommit(double));
      expect(failures.failed).isTrue("a strict double refuses commit before prepare");
    });
  });

  describe("delegateTo", () => {
    it("passes prepare to the implementation", async () => {
      const inner = new SessionStub(undefined);
      await settle(() => sessionStubPrepare(new SessionStub(undefined, SessionStub.delegateTo(inner))));
      expect(inner.onPrepare.callCount).isEqualTo(1, "prepare calls the implementation");
    });

    it("fails prepare with the error of the implementation", async () => {
      const want = new Error("prepare-delegate");
      const inner = new SessionStub(undefined);
      inner.onPrepare.faults(want, 1);
      const double = new SessionStub(undefined, SessionStub.delegateTo(inner));
      const outcome = await settle(() => sessionStubPrepare(double));
      const message = "prepare fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes commit to the implementation", async () => {
      const inner = new SessionStub(undefined);
      await settle(() => sessionStubCommit(new SessionStub(undefined, SessionStub.delegateTo(inner))));
      expect(inner.onCommit.callCount).isEqualTo(1, "commit calls the implementation");
    });

    it("fails commit with the error of the implementation", async () => {
      const want = new Error("commit-delegate");
      const inner = new SessionStub(undefined);
      inner.onCommit.faults(want, 1);
      const double = new SessionStub(undefined, SessionStub.delegateTo(inner));
      const outcome = await settle(() => sessionStubCommit(double));
      const message = "commit fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes old to the implementation", async () => {
      const inner = new SessionStub(undefined);
      await settle(() => sessionStubOld(new SessionStub(undefined, SessionStub.delegateTo(inner))));
      expect(inner.onOld.callCount).isEqualTo(1, "old calls the implementation");
    });

    it("fails old with the error of the implementation", async () => {
      const want = new Error("old-delegate");
      const inner = new SessionStub(undefined);
      inner.onOld.faults(want, 1);
      const double = new SessionStub(undefined, SessionStub.delegateTo(inner));
      const outcome = await settle(() => sessionStubOld(double));
      const message = "old fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => sessionStubInstance(new SessionStub(seat, SessionStub.withClock(clock))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      sessionStubInstance(new SessionStub(seat, SessionStub.withRandSource(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => sessionStubInstance(new SessionStub(seat, SessionStub.benchMode())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        sessionStubInstance(new SessionStub(seat, SessionStub.withClock(clock), SessionStub.withLatency(duration))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails prepare with its error", async () => {
      const want = new Error("prepare-failing");
      const double = new SessionStub(undefined, SessionStub.failing(want));
      const outcome = await settle(() => sessionStubPrepare(double));
      const message = "prepare fails with the error of SessionStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails commit with its error", async () => {
      const want = new Error("commit-failing");
      const double = new SessionStub(undefined, SessionStub.failing(want));
      const outcome = await settle(() => sessionStubCommit(double));
      const message = "commit fails with the error of SessionStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails old with its error", async () => {
      const want = new Error("old-failing");
      const double = new SessionStub(undefined, SessionStub.failing(want));
      const outcome = await settle(() => sessionStubOld(double));
      const message = "old fails with the error of SessionStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withSessionPrepare", () => {
    it("sets the function of prepare", async () => {
      let called = false;
      const option = SessionStub.withSessionPrepare(() => {
        called = true;
        return undefined as never;
      });
      await sessionStubPrepare(new SessionStub(undefined, option));
      expect(called).isTrue("withSessionPrepare sets the function of prepare");
    });
  });

  describe("withSessionCommit", () => {
    it("sets the function of commit", async () => {
      let called = false;
      const option = SessionStub.withSessionCommit(() => {
        called = true;
        return undefined as never;
      });
      await sessionStubCommit(new SessionStub(undefined, option));
      expect(called).isTrue("withSessionCommit sets the function of commit");
    });
  });

  describe("withSessionOld", () => {
    it("sets the function of old", async () => {
      let called = false;
      const option = SessionStub.withSessionOld(() => {
        called = true;
        return undefined as never;
      });
      sessionStubOld(new SessionStub(undefined, option));
      expect(called).isTrue("withSessionOld sets the function of old");
    });
  });
});

describe("SessionOldStub", () => {
  describe("returns", () => {
    it("sets the result of each call of old", async () => {
      const double = new SessionStub(undefined);
      const want = sessionStubOldPin();
      double.onOld.returns(want);
      expect(sessionStubOld(double)).isEqualTo(want, "old returns the result that returns sets");
    });

    it("configures old for strict mode", async () => {
      const failures = new Recorder();
      const double = new SessionStub(failures);
      double.onOld.strict();
      double.onOld.returns(undefined as never);
      await settle(() => sessionStubOld(double));
      expect(failures.failed).isFalse("strict mode accepts the call of old");
    });
  });

  describe("func", () => {
    it("passes the arguments of each call of old to the function", async () => {
      const double = new SessionStub(undefined);
      let got: unknown[] = [];
      double.onOld.func((...received: unknown[]) => {
        got = received;
        return undefined as never;
      });
      sessionStubOld(double);
      const args = sessionStubOldArgs();
      expect(got).isEqualTo([args[0]], "the function receives the arguments of the call");
    });
  });
});
