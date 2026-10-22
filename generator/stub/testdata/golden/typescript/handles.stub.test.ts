import { Recorder, Standard, byIdentity } from "@dokimi/assert";
import { behaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, resetBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import type { Cursor } from "./handles";
import { CursorStub, ProducerStub } from "./handles.stub";

/** Calls open of double with the samples of its arguments. */
function producerStubOpen(double: ProducerStub): Promise<Cursor> {
  return double.open();
}

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * open alone.
 */
function producerStubInstance(double: ProducerStub) {
  return {
    stub: double.onOpen,
    call: () => producerStubOpen(double),
    reset: () => double.resetCalls(),
  };
}

describe("ProducerStub", () => {
  describe("open", () => {
    const checks = behaviour("Producer.open", (seat) => {
      const double = new ProducerStub(seat);
      return {
        stub: double.onOpen,
        call: () => producerStubOpen(double),
        override: (mark) => {
          double.onOpen.func(() => {
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
      const double = new ProducerStub(undefined);
      let seen = 0;
      double.onOpen.onRecord(() => {
        seen += 1;
      });
      await producerStubOpen(double);
      await producerStubOpen(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of open");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => producerStubInstance(new ProducerStub(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of open", async () => {
      const double = new ProducerStub(undefined);
      await producerStubOpen(double);
      double.resetCalls();
      expect(double.onOpen.callCount).isEqualTo(0, "resetCalls clears the calls of open");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => producerStubInstance(ProducerStub.newDelayed(seat, duration, clock)),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails open with its error", async () => {
      const want = new Error("open-failing");
      const double = ProducerStub.newFailing(undefined, want);
      const outcome = await settle(() => producerStubOpen(double));
      const message = "open fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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

describe("ProducerOpenStub", () => {
  describe("returns", () => {
    it("configures open for strict mode", async () => {
      const failures = new Recorder();
      const double = new ProducerStub(failures);
      double.onOpen.strict();
      double.onOpen.returns(undefined as never);
      await settle(() => producerStubOpen(double));
      expect(failures.failed).isFalse("strict mode accepts the call of open");
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

/** Calls close of double with the samples of its arguments. */
function cursorStubClose(double: CursorStub): void {
  return double.close();
}

/**
 * Returns the instance of double that the checks of the constructors and the
 * options drive. The options apply to the whole double, so the checks call
 * next alone.
 */
function cursorStubInstance(double: CursorStub) {
  return {
    stub: double.onNext,
    call: () => cursorStubNext(double),
    reset: () => double.resetCalls(),
  };
}

describe("CursorStub", () => {
  describe("next", () => {
    const checks = behaviour("Cursor.next", (seat) => {
      const double = new CursorStub(seat);
      return {
        stub: double.onNext,
        call: () => cursorStubNext(double),
        override: (mark) => {
          double.onNext.func(() => {
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
      const double = new CursorStub(undefined);
      let seen = 0;
      double.onNext.onRecord(() => {
        seen += 1;
      });
      await cursorStubNext(double);
      await cursorStubNext(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of next");
    });
  });

  describe("close", () => {
    const checks = behaviour("Cursor.close", (seat) => {
      const double = new CursorStub(seat);
      return {
        stub: double.onClose,
        call: () => cursorStubClose(double),
        override: (mark) => {
          double.onClose.func(() => {
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
      const double = new CursorStub(undefined);
      let seen = 0;
      double.onClose.onRecord(() => {
        seen += 1;
      });
      cursorStubClose(double);
      cursorStubClose(double);
      expect(seen).isEqualTo(2, "onRecord runs once for each call of close");
    });
  });

  describe("resetCalls", () => {
    const checks = resetBehaviour((seat) => cursorStubInstance(new CursorStub(seat)));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }

    it("clears the recorded calls of next", async () => {
      const double = new CursorStub(undefined);
      await cursorStubNext(double);
      double.resetCalls();
      expect(double.onNext.callCount).isEqualTo(0, "resetCalls clears the calls of next");
    });

    it("keeps the result of returns of next", async () => {
      const double = new CursorStub(undefined);
      const want = cursorStubNextPin();
      double.onNext.returns(want);
      double.resetCalls();
      expect(await cursorStubNext(double)).isEqualTo(want, "resetCalls keeps the result of next");
    });

    it("clears the recorded calls of close", async () => {
      const double = new CursorStub(undefined);
      cursorStubClose(double);
      double.resetCalls();
      expect(double.onClose.callCount).isEqualTo(0, "resetCalls clears the calls of close");
    });
  });

  describe("newDelayed", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) => cursorStubInstance(CursorStub.newDelayed(seat, duration, clock)),
      true,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("newFailing", () => {
    it("fails next with its error", async () => {
      const want = new Error("next-failing");
      const double = CursorStub.newFailing(undefined, want);
      const outcome = await settle(() => cursorStubNext(double));
      const message = "next fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails close with its error", async () => {
      const want = new Error("close-failing");
      const double = CursorStub.newFailing(undefined, want);
      const outcome = await settle(() => cursorStubClose(double));
      const message = "close fails with the error of newFailing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
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

describe("CursorNextStub", () => {
  describe("returns", () => {
    it("sets the result of each call of next", async () => {
      const double = new CursorStub(undefined);
      const want = cursorStubNextPin();
      double.onNext.returns(want);
      expect(await cursorStubNext(double)).isEqualTo(want, "next returns the result that returns sets");
    });

    it("configures next for strict mode", async () => {
      const failures = new Recorder();
      const double = new CursorStub(failures);
      double.onNext.strict();
      double.onNext.returns(undefined as never);
      await settle(() => cursorStubNext(double));
      expect(failures.failed).isFalse("strict mode accepts the call of next");
    });
  });
});
