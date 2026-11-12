import { Recorder, type Seat, Standard, byIdentity, check } from "@dokimi/assert";
import { type Subject, behaviour, resetBehaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { type SessionCommitCall, type SessionOldCall, type SessionPrepareCall, SessionStub } from "./mixed.stub";

/** Calls prepare of double with the samples of its arguments. */
function sessionStubPrepare(double: SessionStub): Promise<void> {
  return double.prepare();
}

/** Returns the subject of behaviour and resetBehaviour for prepare, on a new double that reports to seat. */
function sessionStubPrepareSubject(seat: Seat): Subject<SessionPrepareCall, undefined> {
  const double = new SessionStub(seat);
  return {
    stub: double.onPrepare,
    call: () => sessionStubPrepare(double),
    reset: () => double.resetCalls(),
    override: (seen) => {
      double.onPrepare.func(() => {
        seen({} as SessionPrepareCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => sessionStubPrepare(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Calls commit of double with the samples of its arguments. */
function sessionStubCommit(double: SessionStub): Promise<void> {
  return double.commit();
}

/** Returns the subject of behaviour and resetBehaviour for commit, on a new double that reports to seat. */
function sessionStubCommitSubject(seat: Seat): Subject<SessionCommitCall, undefined> {
  const double = new SessionStub(seat);
  return {
    stub: double.onCommit,
    call: () => sessionStubCommit(double),
    reset: () => double.resetCalls(),
    override: (seen) => {
      double.onCommit.func(() => {
        seen({} as SessionCommitCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => sessionStubCommit(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
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

/** Returns the subject of behaviour and resetBehaviour for old, on a new double that reports to seat. */
function sessionStubOldSubject(seat: Seat): Subject<SessionOldCall, string | undefined> {
  const double = new SessionStub(seat);
  return {
    stub: double.onOld,
    call: () => sessionStubOld(double),
    reset: () => double.resetCalls(),
    arguments: (s, c) => {
      const args = sessionStubOldArgs();
      check.isEqualTo(s, c.key, args[0], "the call records key");
    },
    result: async () => sessionStubOld(double),
    returns: () => {
      const want = sessionStubOldPin();
      double.onOld.returns(want);
      return want;
    },
    override: (seen) => {
      double.onOld.func((...received: unknown[]) => {
        seen({ key: received[0] } as SessionOldCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => sessionStubOld(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call prepare alone.
 */
function sessionStubInstance(double: SessionStub) {
  return {
    stub: double.onPrepare,
    call: () => sessionStubPrepare(double),
  };
}

describe("SessionStub", () => {
  describe("prepare", () => {
    for (const c of behaviour("Session.prepare", sessionStubPrepareSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("commit", () => {
    for (const c of behaviour("Session.commit", sessionStubCommitSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("old", () => {
    for (const c of behaviour("Session.old", sessionStubOldSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("prepare", () => {
      for (const c of resetBehaviour("Session.prepare", sessionStubPrepareSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("commit", () => {
      for (const c of resetBehaviour("Session.commit", sessionStubCommitSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("old", () => {
      for (const c of resetBehaviour("Session.old", sessionStubOldSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
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
