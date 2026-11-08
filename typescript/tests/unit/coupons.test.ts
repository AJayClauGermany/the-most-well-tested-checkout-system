// Unit tests for coupons (rules §4.4).

import { describe, expect, it } from "vitest";

import { CheckoutError, couponDiscount, normalizeCouponCode } from "../../src/index.ts";

describe("normalizeCouponCode", () => {
  it.each([
    ["SAVE10", "SAVE10"],
    ["save10", "SAVE10"],
    ["  save10 ", "SAVE10"],
    ["\tWelcome5\n", "WELCOME5"],
    ["nope", "NOPE"],
  ])("turns %j into %j", (input, expected) => {
    expect(normalizeCouponCode(input)).toBe(expected);
  });

  it.each([undefined, null, "", "   ", "\t\n"])("treats %j as no coupon", (input) => {
    expect(normalizeCouponCode(input)).toBeUndefined();
  });
});

describe("couponDiscount", () => {
  it("is 0 without a coupon", () => {
    expect(couponDiscount(undefined, 10000)).toBe(0);
    expect(couponDiscount(null, 10000)).toBe(0);
    expect(couponDiscount("", 10000)).toBe(0);
    expect(couponDiscount("  ", 10000)).toBe(0);
  });

  describe("SAVE10", () => {
    it("takes 10% off, rounded half up", () => {
      expect(couponDiscount("SAVE10", 500)).toBe(50);
      expect(couponDiscount("SAVE10", 225)).toBe(23); // 22.5 → 23
      expect(couponDiscount("SAVE10", 224)).toBe(22); // 22.4 → 22
      expect(couponDiscount("SAVE10", 1798)).toBe(180); // 179.8 → 180
    });

    it("has no minimum spend", () => {
      expect(couponDiscount("SAVE10", 1)).toBe(0); // 0.1 → 0
      expect(couponDiscount("SAVE10", 5)).toBe(1); // 0.5 → 1
    });

    it("ignores case and surrounding spaces", () => {
      expect(couponDiscount(" save10 ", 500)).toBe(50);
    });
  });

  describe("WELCOME5", () => {
    it("takes exactly $5.00 off", () => {
      expect(couponDiscount("WELCOME5", 2500)).toBe(500);
      expect(couponDiscount("WELCOME5", 99999)).toBe(500);
    });

    it("works at exactly the $25.00 minimum and fails one cent below it", () => {
      expect(couponDiscount("WELCOME5", 2500)).toBe(500);
      expect(() => couponDiscount("WELCOME5", 2499)).toThrow(
        expect.objectContaining({ code: "COUPON_MINIMUM_NOT_MET" }),
      );
    });

    it("ignores case and surrounding spaces", () => {
      expect(couponDiscount("welcome5", 2500)).toBe(500);
    });
  });

  it("rejects a code it does not know with UNKNOWN_COUPON", () => {
    expect(() => couponDiscount("FREESTUFF", 10000)).toThrow(CheckoutError);
    expect(() => couponDiscount("FREESTUFF", 10000)).toThrow(
      expect.objectContaining({ code: "UNKNOWN_COUPON" }),
    );
  });

  it("explains what went wrong in the error message", () => {
    expect(() => couponDiscount("freestuff", 10000)).toThrow("Unknown coupon: FREESTUFF");
    expect(() => couponDiscount("WELCOME5", 2499)).toThrow(
      "WELCOME5 needs a spend of at least 2500 cents, got 2499",
    );
  });

  it("does not accept a code that merely starts or ends with a real code", () => {
    expect(() => couponDiscount("SAVE100", 10000)).toThrow(
      expect.objectContaining({ code: "UNKNOWN_COUPON" }),
    );
    expect(() => couponDiscount("WELCOME50", 10000)).toThrow(
      expect.objectContaining({ code: "UNKNOWN_COUPON" }),
    );
  });
});
