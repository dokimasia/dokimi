import { type Seat, Standard, byIdentity } from "@dokimi/assert";
import { type Subject, behaviour, resetBehaviour, settle } from "@dokimi/runtime/stub/behaviour";
import { benchModeBehaviour, clockBehaviour, latencyBehaviour, randSourceBehaviour, strictBehaviour } from "@dokimi/runtime/stub/double-behaviour";
import { describe, expect, it } from "vitest";
import { type FeedNamesCall, type FeedRowsCall, FeedStub } from "./stream.stub";

/** Calls names of double with the samples of its arguments. */
function feedStubNames(double: FeedStub): Iterable<string> {
  return double.names();
}

/** Returns the subject of behaviour and resetBehaviour for names, on a new double that reports to seat. */
function feedStubNamesSubject(seat: Seat): Subject<FeedNamesCall, Iterable<string> | undefined> {
  const double = new FeedStub(seat);
  return {
    stub: double.onNames,
    call: () => feedStubNames(double),
    reset: () => double.resetCalls(),
    result: async () => feedStubNames(double),
    override: (seen) => {
      double.onNames.func(() => {
        seen({} as FeedNamesCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => feedStubNames(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/** Calls rows of double with the samples of its arguments. */
function feedStubRows(double: FeedStub): AsyncIterable<number> {
  return double.rows();
}

/** Returns the subject of behaviour and resetBehaviour for rows, on a new double that reports to seat. */
function feedStubRowsSubject(seat: Seat): Subject<FeedRowsCall, AsyncIterable<number> | undefined> {
  const double = new FeedStub(seat);
  return {
    stub: double.onRows,
    call: () => feedStubRows(double),
    reset: () => double.resetCalls(),
    result: async () => feedStubRows(double),
    override: (seen) => {
      double.onRows.func(() => {
        seen({} as FeedRowsCall);
        return undefined as never;
      });
    },
    fails: async () => {
      const outcome = await settle(() => feedStubRows(double));
      if (outcome.threw) {
        return outcome.error;
      }
      return undefined;
    },
  };
}

/**
 * Returns the instance of double that the checks of the options drive. The
 * options apply to the whole double, so the checks call names alone.
 */
function feedStubInstance(double: FeedStub) {
  return {
    stub: double.onNames,
    call: () => feedStubNames(double),
  };
}

describe("FeedStub", () => {
  describe("names", () => {
    for (const c of behaviour("Feed.names", feedStubNamesSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("rows", () => {
    for (const c of behaviour("Feed.rows", feedStubRowsSubject)) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("resetCalls", () => {
    describe("names", () => {
      for (const c of resetBehaviour("Feed.names", feedStubNamesSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });

    describe("rows", () => {
      for (const c of resetBehaviour("Feed.rows", feedStubRowsSubject)) {
        it(c.name, () => c.run(new Standard()));
      }
    });
  });

  describe("strict", () => {
    const checks = strictBehaviour((seat) => feedStubInstance(new FeedStub(seat, FeedStub.strict())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("delegateTo", () => {
    it("passes names to the implementation", async () => {
      const inner = new FeedStub(undefined);
      await settle(() => feedStubNames(new FeedStub(undefined, FeedStub.delegateTo(inner))));
      expect(inner.onNames.callCount).isEqualTo(1, "names calls the implementation");
    });

    it("fails names with the error of the implementation", async () => {
      const want = new Error("names-delegate");
      const inner = new FeedStub(undefined);
      inner.onNames.faults(want, 1);
      const double = new FeedStub(undefined, FeedStub.delegateTo(inner));
      const outcome = await settle(() => feedStubNames(double));
      const message = "names fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("passes rows to the implementation", async () => {
      const inner = new FeedStub(undefined);
      await settle(() => feedStubRows(new FeedStub(undefined, FeedStub.delegateTo(inner))));
      expect(inner.onRows.callCount).isEqualTo(1, "rows calls the implementation");
    });

    it("fails rows with the error of the implementation", async () => {
      const want = new Error("rows-delegate");
      const inner = new FeedStub(undefined);
      inner.onRows.faults(want, 1);
      const double = new FeedStub(undefined, FeedStub.delegateTo(inner));
      const outcome = await settle(() => feedStubRows(double));
      const message = "rows fails with the error of the implementation";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withClock", () => {
    const checks = clockBehaviour(
      (seat, clock) => feedStubInstance(new FeedStub(seat, FeedStub.withClock(clock))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withRandSource", () => {
    const checks = randSourceBehaviour((seat, random) =>
      feedStubInstance(new FeedStub(seat, FeedStub.withRandSource(random))),
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("benchMode", () => {
    const checks = benchModeBehaviour((seat) => feedStubInstance(new FeedStub(seat, FeedStub.benchMode())));
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("withLatency", () => {
    const checks = latencyBehaviour(
      (seat, duration, clock) =>
        feedStubInstance(new FeedStub(seat, FeedStub.withClock(clock), FeedStub.withLatency(duration))),
      false,
    );
    for (const c of checks) {
      it(c.name, () => c.run(new Standard()));
    }
  });

  describe("failing", () => {
    it("fails names with its error", async () => {
      const want = new Error("names-failing");
      const double = new FeedStub(undefined, FeedStub.failing(want));
      const outcome = await settle(() => feedStubNames(double));
      const message = "names fails with the error of FeedStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });

    it("fails rows with its error", async () => {
      const want = new Error("rows-failing");
      const double = new FeedStub(undefined, FeedStub.failing(want));
      const outcome = await settle(() => feedStubRows(double));
      const message = "rows fails with the error of FeedStub.failing";
      expect(outcome.threw && outcome.error).isEqualTo(want, message, byIdentity());
    });
  });

  describe("withFeedNames", () => {
    it("sets the function of names", async () => {
      let called = false;
      const option = FeedStub.withFeedNames(() => {
        called = true;
        return undefined as never;
      });
      feedStubNames(new FeedStub(undefined, option));
      expect(called).isTrue("withFeedNames sets the function of names");
    });
  });

  describe("withFeedRows", () => {
    it("sets the function of rows", async () => {
      let called = false;
      const option = FeedStub.withFeedRows(() => {
        called = true;
        return undefined as never;
      });
      feedStubRows(new FeedStub(undefined, option));
      expect(called).isTrue("withFeedRows sets the function of rows");
    });
  });
});

describe("FeedRowsStub", () => {
  describe("yields", () => {
    it("yields each item", async () => {
      const double = new FeedStub(undefined);
      double.onRows.yields(undefined as never, undefined as never);
      let seen = 0;
      for await (const _ of feedStubRows(double)) {
        seen += 1;
      }
      expect(seen).isEqualTo(2, "rows yields each item");
    });

    it("stops yielding when the caller breaks", async () => {
      const double = new FeedStub(undefined);
      double.onRows.yields(undefined as never, undefined as never, undefined as never);
      let seen = 0;
      for await (const _ of feedStubRows(double)) {
        seen += 1;
        break;
      }
      expect(seen).isEqualTo(1, "rows stops when the caller breaks");
    });
  });
});
