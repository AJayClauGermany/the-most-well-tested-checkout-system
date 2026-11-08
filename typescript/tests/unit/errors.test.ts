// Unit tests for the error codes and the CheckoutError class.

import { describe, expect, it } from "vitest";

import { CheckoutError, ERROR_CODES } from "../../src/index.ts";

describe("ERROR_CODES", () => {
  it("lists exactly the six codes from rules §9, spelled the same way", () => {
    // toEqual compares the contents of two arrays (toBe would check that
    // they are the very same array object, which they are not).
    expect(ERROR_CODES).toEqual([
      "EMPTY_BASKET",
      "UNKNOWN_PRODUCT",
      "INVALID_QUANTITY",
      "UNKNOWN_SHIPPING_METHOD",
      "UNKNOWN_COUPON",
      "COUPON_MINIMUM_NOT_MET",
    ]);
  });
});

describe("CheckoutError", () => {
  it("carries the error code", () => {
    const error = new CheckoutError("UNKNOWN_COUPON");

    expect(error.code).toBe("UNKNOWN_COUPON");
  });

  it("is a normal Error, so it can be thrown and caught", () => {
    const error = new CheckoutError("EMPTY_BASKET", "The basket has no items");

    expect(error).toBeInstanceOf(Error);
    expect(error.message).toBe("The basket has no items");
    expect(() => {
      throw error;
    }).toThrow(CheckoutError);
  });
});
