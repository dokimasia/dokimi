import { type Seat, Standard, byIdentity } from "@dokimi/assert";
import { type Subject, behaviour, resetBehaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import type { Cursor } from "./handles";
import { type CursorCloseCall, type CursorNextCall, CursorStub, type ProducerOpenCall, ProducerStub } from "./handles.stub";

/** Calls open of double with the samples of its arguments. */
function producerStubOpen(double: ProducerStub): Promise<Cursor> {
  return double.open();
}

/** Returns the subject of behaviour and resetBehaviour for open, on a new double that reports to seat. */
function producerStubOpenSubject(seat: Seat): Subject<ProducerOpenCall, Cursor | undefined> {
  const double = new ProducerStub(seat);
  return {
    stub: double.onOpen,
    call: () => producerStubOpen(double),
    reset: () => double.resetCalls(),
    result: async () => await producerStubOpen(double),
    override: (seen) => {
      double.onOpen.func(() => {
        seen({} as ProducerOpenCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => producerStubOpen(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call open alone.
 */
function producerStubInstance(double: ProducerStub) {
  return {
    stub: double.onOpen,
    call: () => producerStubOpen(double),
  };
}

describe("ProducerStub", () => {
  describe("open", () => {
    for (const c of behaviour("Producer.open", producerStubOpenSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("open", () => {
      for (const c of resetBehaviour("Producer.open", producerStubOpenSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => producerStubInstance(new ProducerStub(seat, ProducerStub.strict())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes open to the implementation", async () => {
      const inner = new ProducerStub(undefined);
      await settle(() => producerStubOpen(new ProducerStub(undefined, ProducerStub.delegateTo(inner))));
      expect(inner.onOpen.callCount).isEqualTo(1, "open calls the implementation");
    });

    it("fails open with the error of the implementation", async () => {
      const want = new Error("open-delegate");
      const inner = new ProducerStub(undefined);
      inner.onOpen.faults(want, 1);
      const double = new ProducerStub(undefined, ProducerStub.delegateTo(inner));
      const outcome = await settle(() => producerStubOpen(double));
      const message = "open fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => producerStubInstance(new ProducerStub(seat, ProducerStub.withClock(clock))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      producerStubInstance(new ProducerStub(seat, ProducerStub.withRandSource(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => producerStubInstance(new ProducerStub(seat, ProducerStub.benchMode())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        producerStubInstance(new ProducerStub(seat, ProducerStub.withClock(clock), ProducerStub.withLatency(duration))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails open with its error", async () => {
      const want = new Error("open-failing");
      const double = new ProducerStub(undefined, ProducerStub.failing(want));
      const outcome = await settle(() => producerStubOpen(double));
      const message = "open fails with the error of ProducerStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withProducerOpen", () => {
    it("sets the function of open", async () => {
      let called = false;
      const option = ProducerStub.withProducerOpen(() => {
        called = true;
        return undefined as never;
      });
      await producerStubOpen(new ProducerStub(undefined, option));
      expect(called).isTrue("withProducerOpen sets the function of open");
    });
  });
});

/** Calls next of double with the samples of its arguments. */
function cursorStubNext(double: CursorStub): Promise<string> {
  return double.next();
}

/** Returns the sample of the result of next. */
function cursorStubNextPin(): string {
  let want0: string;
  want0 = "next-a";
  return want0;
}

/** Returns the subject of behaviour and resetBehaviour for next, on a new double that reports to seat. */
function cursorStubNextSubject(seat: Seat): Subject<CursorNextCall, string | undefined> {
  const double = new CursorStub(seat);
  return {
    stub: double.onNext,
    call: () => cursorStubNext(double),
    reset: () => double.resetCalls(),
    result: async () => await cursorStubNext(double),
    returns: () => {
      const want = cursorStubNextPin();
      double.onNext.returns(want);
      return want;
    },
    override: (seen) => {
      double.onNext.func(() => {
        seen({} as CursorNextCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => cursorStubNext(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Calls close of double with the samples of its arguments. */
function cursorStubClose(double: CursorStub): void {
  return double.close();
}

/** Returns the subject of behaviour and resetBehaviour for close, on a new double that reports to seat. */
function cursorStubCloseSubject(seat: Seat): Subject<CursorCloseCall, undefined> {
  const double = new CursorStub(seat);
  return {
    stub: double.onClose,
    call: () => cursorStubClose(double),
    reset: () => double.resetCalls(),
    override: (seen) => {
      double.onClose.func(() => {
        seen({} as CursorCloseCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => cursorStubClose(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call next alone.
 */
function cursorStubInstance(double: CursorStub) {
  return {
    stub: double.onNext,
    call: () => cursorStubNext(double),
  };
}

describe("CursorStub", () => {
  describe("next", () => {
    for (const c of behaviour("Cursor.next", cursorStubNextSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("close", () => {
    for (const c of behaviour("Cursor.close", cursorStubCloseSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("next", () => {
      for (const c of resetBehaviour("Cursor.next", cursorStubNextSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("close", () => {
      for (const c of resetBehaviour("Cursor.close", cursorStubCloseSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => cursorStubInstance(new CursorStub(seat, CursorStub.strict())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes next to the implementation", async () => {
      const inner = new CursorStub(undefined);
      await settle(() => cursorStubNext(new CursorStub(undefined, CursorStub.delegateTo(inner))));
      expect(inner.onNext.callCount).isEqualTo(1, "next calls the implementation");
    });

    it("fails next with the error of the implementation", async () => {
      const want = new Error("next-delegate");
      const inner = new CursorStub(undefined);
      inner.onNext.faults(want, 1);
      const double = new CursorStub(undefined, CursorStub.delegateTo(inner));
      const outcome = await settle(() => cursorStubNext(double));
      const message = "next fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes close to the implementation", async () => {
      const inner = new CursorStub(undefined);
      await settle(() => cursorStubClose(new CursorStub(undefined, CursorStub.delegateTo(inner))));
      expect(inner.onClose.callCount).isEqualTo(1, "close calls the implementation");
    });

    it("fails close with the error of the implementation", async () => {
      const want = new Error("close-delegate");
      const inner = new CursorStub(undefined);
      inner.onClose.faults(want, 1);
      const double = new CursorStub(undefined, CursorStub.delegateTo(inner));
      const outcome = await settle(() => cursorStubClose(double));
      const message = "close fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => cursorStubInstance(new CursorStub(seat, CursorStub.withClock(clock))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      cursorStubInstance(new CursorStub(seat, CursorStub.withRandSource(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => cursorStubInstance(new CursorStub(seat, CursorStub.benchMode())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        cursorStubInstance(new CursorStub(seat, CursorStub.withClock(clock), CursorStub.withLatency(duration))),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails next with its error", async () => {
      const want = new Error("next-failing");
      const double = new CursorStub(undefined, CursorStub.failing(want));
      const outcome = await settle(() => cursorStubNext(double));
      const message = "next fails with the error of CursorStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails close with its error", async () => {
      const want = new Error("close-failing");
      const double = new CursorStub(undefined, CursorStub.failing(want));
      const outcome = await settle(() => cursorStubClose(double));
      const message = "close fails with the error of CursorStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withCursorNext", () => {
    it("sets the function of next", async () => {
      let called = false;
      const option = CursorStub.withCursorNext(() => {
        called = true;
        return undefined as never;
      });
      await cursorStubNext(new CursorStub(undefined, option));
      expect(called).isTrue("withCursorNext sets the function of next");
    });
  });

  describe("withCursorClose", () => {
    it("sets the function of close", async () => {
      let called = false;
      const option = CursorStub.withCursorClose(() => {
        called = true;
        return undefined as never;
      });
      cursorStubClose(new CursorStub(undefined, option));
      expect(called).isTrue("withCursorClose sets the function of close");
    });
  });
});
