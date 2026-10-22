import { Controlled } from "@dokimi/assert";
import { describe, expect, it } from "vitest";
import { BoxStub, PointStub } from "./point.stub";

/** Assigns the sample of each property of PointStub to double. */
function pointStubProps(double: PointStub): void {
  double.x = 1.5;
  double.y = 1.5;
}

describe("PointStub", () => {
  describe("resetCalls", () => {
    const tests: { name: string; giveDouble: () => PointStub }[] = [
      {
        name: "keeps each property of a new double",
        giveDouble: () => new PointStub(undefined),
      },
      {
        name: "keeps each property of a strict double",
        giveDouble: () => new PointStub(undefined, PointStub.strict()),
      },
      {
        name: "keeps each property of a double with a clock",
        giveDouble: () => new PointStub(undefined, PointStub.withClock(new Controlled())),
      },
      {
        name: "keeps each property of a double with a random source",
        giveDouble: () => new PointStub(undefined, PointStub.withRandSource(() => 0)),
      },
      {
        name: "keeps each property of a double in bench mode",
        giveDouble: () => new PointStub(undefined, PointStub.benchMode()),
      },
      {
        name: "keeps each property of a double with a latency",
        giveDouble: () => new PointStub(undefined, PointStub.withLatency(1)),
      },
      {
        name: "keeps each property of a delayed double",
        giveDouble: () => PointStub.newDelayed(undefined, 1, new Controlled()),
      },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        const double = tt.giveDouble();
        pointStubProps(double);
        double.resetCalls();
        const want = new PointStub(undefined);
        pointStubProps(want);
        expect.soft(double.x).isEqualTo(want.x, "resetCalls keeps the property x");
        expect.soft(double.y).isEqualTo(want.y, "resetCalls keeps the property y");
      });
    }
  });
});

describe("BoxStub", () => {
  describe("resetCalls", () => {
    const tests: { name: string; giveDouble: () => BoxStub<unknown> }[] = [
      {
        name: "keeps each property of a new double",
        giveDouble: () => new BoxStub<unknown>(undefined),
      },
      {
        name: "keeps each property of a strict double",
        giveDouble: () => new BoxStub<unknown>(undefined, BoxStub.strict<unknown>()),
      },
      {
        name: "keeps each property of a double with a clock",
        giveDouble: () => new BoxStub<unknown>(undefined, BoxStub.withClock<unknown>(new Controlled())),
      },
      {
        name: "keeps each property of a double with a random source",
        giveDouble: () => new BoxStub<unknown>(undefined, BoxStub.withRandSource<unknown>(() => 0)),
      },
      {
        name: "keeps each property of a double in bench mode",
        giveDouble: () => new BoxStub<unknown>(undefined, BoxStub.benchMode<unknown>()),
      },
      {
        name: "keeps each property of a double with a latency",
        giveDouble: () => new BoxStub<unknown>(undefined, BoxStub.withLatency<unknown>(1)),
      },
      {
        name: "keeps each property of a delayed double",
        giveDouble: () => BoxStub.newDelayed<unknown>(undefined, 1, new Controlled()),
      },
    ];
    for (const tt of tests) {
      it(tt.name, () => {
        const double = tt.giveDouble();
        double.resetCalls();
        const want = new BoxStub<unknown>(undefined);
        expect(double.value).isEqualTo(want.value, "resetCalls keeps the property value");
      });
    }
  });
});
