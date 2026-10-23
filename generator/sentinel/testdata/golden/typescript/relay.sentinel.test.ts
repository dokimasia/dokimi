import { describe, expect, it } from "vitest";
import { CallError, TimeoutError } from "./relay";

describe("CallError", () => {
  it("matches no other error class of its module", () => {
    const got = new CallError("message-a", 1.5);
    const others = [new TimeoutError("op-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });
});

function timeoutErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new TimeoutError("op-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });

  it("opens with the prefix of its module", () => {
    const got = new TimeoutError("op-a");
    expect(got.message).hasPrefix("relay: ", "the message opens with the prefix of the module");
  });

  it("contains each text member", () => {
    const got = new TimeoutError("op-a");
    expect(got.message).doesContain(got.op, "the message contains op");
  });
}

describe("TimeoutError", () => {
  it("matches no other error class of its module", () => {
    const got = new TimeoutError("op-a");
    const others = [new CallError("message-a", 1.5)];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", timeoutErrorMessageCases);
});
