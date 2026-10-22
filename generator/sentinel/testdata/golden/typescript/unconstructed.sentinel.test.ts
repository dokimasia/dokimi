import { describe, expect, it } from "vitest";
import { ChainError, KeptError, LostError } from "./unconstructed";

function chainErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new ChainError("cause-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });
}

function chainErrorCauseCases(): void {
  it("is the cause of the construction", () => {
    const cause = new Error("cause");
    const got = new ChainError(cause);
    expect(got.cause).isEqualTo(cause, "the error keeps its cause");
  });
}

describe("ChainError", () => {
  it("matches no other error class of its module", () => {
    const got = new ChainError("cause-a");
    const others = [new KeptError(), new LostError("message-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", chainErrorMessageCases);
  describe("cause", chainErrorCauseCases);
});

function keptErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new KeptError();
    expect(got.message).isNotEqualTo("", "the error has a message");
  });
}

describe("KeptError", () => {
  it("matches no other error class of its module", () => {
    const got = new KeptError();
    const others = [new ChainError("cause-a"), new LostError("message-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", keptErrorMessageCases);
});

function lostErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new LostError("message-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });
}

describe("LostError", () => {
  it("matches no other error class of its module", () => {
    const got = new LostError("message-a");
    const others = [new ChainError("cause-a"), new KeptError()];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", lostErrorMessageCases);
});
