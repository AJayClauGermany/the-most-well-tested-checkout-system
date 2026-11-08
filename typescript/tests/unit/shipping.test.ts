// Unit tests for shipping (rules §4.5).

import { describe, expect, it } from "vitest";

import { CheckoutError, parseShippingMethod, shippingCost } from "../../src/index.ts";

describe("parseShippingMethod", () => {
  it("defaults to standard when the basket says nothing", () => {
    expect(parseShippingMethod(undefined)).toBe("standard");
  });

  it("accepts the two known methods", () => {
    expect(parseShippingMethod("standard")).toBe("standard");
    expect(parseShippingMethod("express")).toBe("express");
  });

  it.each(["overnight", "Express", "STANDARD", " standard", "express ", ""])(
    "rejects %j with UNKNOWN_SHIPPING_METHOD",
    (method) => {
      expect(() => parseShippingMethod(method)).toThrow(CheckoutError);
      expect(() => parseShippingMethod(method)).toThrow(
        expect.objectContaining({ code: "UNKNOWN_SHIPPING_METHOD" }),
      );
    },
  );
});

describe("parseShippingMethod error message", () => {
  it("names the method it did not recognise", () => {
    expect(() => parseShippingMethod("overnight")).toThrow("Unknown shipping method: overnight");
  });
});

describe("shippingCost", () => {
  describe("standard", () => {
    it("costs $4.99 below $50.00", () => {
      expect(shippingCost("standard", 0)).toBe(499);
      expect(shippingCost("standard", 4999)).toBe(499);
    });

    it("is free from exactly $50.00", () => {
      expect(shippingCost("standard", 5000)).toBe(0);
      expect(shippingCost("standard", 5001)).toBe(0);
      expect(shippingCost("standard", 100000)).toBe(0);
    });
  });

  describe("express", () => {
    it("always costs $9.99", () => {
      expect(shippingCost("express", 0)).toBe(999);
      expect(shippingCost("express", 4999)).toBe(999);
      expect(shippingCost("express", 5000)).toBe(999);
      expect(shippingCost("express", 100000)).toBe(999);
    });
  });
});
