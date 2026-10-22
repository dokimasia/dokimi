import { describe, expect, it } from "vitest";
import { ClosedError, MissingError, NotFoundError, WrappedError } from "./store";

function closedErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new ClosedError();
    expect(got.message).isNotEqualTo("", "the error has a message");
  });

  it("opens with the prefix of its module", () => {
    const got = new ClosedError();
    expect(got.message).hasPrefix("kv: ", "the message opens with the prefix of the module");
  });
}

describe("ClosedError", () => {
  it("matches no other error class of its module", () => {
    const got = new ClosedError();
    const others = [new MissingError("key-a"), new NotFoundError("key-a"), new WrappedError("op-a", "cause-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", closedErrorMessageCases);
});

function missingErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new MissingError("key-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });

  it("opens with the prefix of its module", () => {
    const got = new MissingError("key-a");
    expect(got.message).hasPrefix("kv: ", "the message opens with the prefix of the module");
  });
}

describe("MissingError", () => {
  it("matches no other error class of its module", () => {
    const got = new MissingError("key-a");
    const others = [new ClosedError(), new NotFoundError("key-a"), new WrappedError("op-a", "cause-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", missingErrorMessageCases);
});

function notFoundErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new NotFoundError("key-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });

  it("opens with the prefix of its module", () => {
    const got = new NotFoundError("key-a");
    expect(got.message).hasPrefix("kv: ", "the message opens with the prefix of the module");
  });

  it("contains each text member", () => {
    const got = new NotFoundError("key-a");
    expect(got.message).doesContain(got.key, "the message contains key");
  });
}

describe("NotFoundError", () => {
  it("matches no other error class of its module", () => {
    const got = new NotFoundError("key-a");
    const others = [new ClosedError(), new MissingError("key-a"), new WrappedError("op-a", "cause-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", notFoundErrorMessageCases);
});

function wrappedErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new WrappedError("op-a", "cause-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });

  it("opens with the prefix of its module", () => {
    const got = new WrappedError("op-a", "cause-a");
    expect(got.message).hasPrefix("kv: ", "the message opens with the prefix of the module");
  });

  it("contains each text member", () => {
    const got = new WrappedError("op-a", "cause-a");
    expect(got.message).doesContain(got.op, "the message contains op");
  });
}

function wrappedErrorCauseCases(): void {
  it("is the cause of the construction", () => {
    const cause = new Error("cause");
    const got = new WrappedError("op-a", cause);
    expect(got.cause).isEqualTo(cause, "the error keeps its cause");
  });
}

describe("WrappedError", () => {
  it("matches no other error class of its module", () => {
    const got = new WrappedError("op-a", "cause-a");
    const others = [new ClosedError(), new MissingError("key-a"), new NotFoundError("key-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", wrappedErrorMessageCases);
  describe("cause", wrappedErrorCauseCases);
});
