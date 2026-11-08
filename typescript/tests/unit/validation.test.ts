// Unit tests for basket validation (rules §4.1, rows 1 to 4) and the order errors are
// reported in. The check order matters: a basket can be wrong in several ways at once.

import { describe, expect, it } from "vitest";

import {
  CheckoutError,
  isValidQuantity,
  MAX_QUANTITY,
  MIN_QUANTITY,
  validateItems,
} from "../../src/index.ts";
import type { Basket, ErrorCode } from "../../src/index.ts";
import { tinyCatalog } from "../support/fixtures.ts";

/** Asserts that validating `basket` fails with exactly this code. */
function expectError(basket: Basket, code: ErrorCode): void {
  expect(() => validateItems(basket, tinyCatalog)).toThrow(CheckoutError);
  expect(() => validateItems(basket, tinyCatalog)).toThrow(expect.objectContaining({ code }));
}

describe("isValidQuantity", () => {
  it("has the limits from the rules", () => {
    expect(MIN_QUANTITY).toBe(1);
    expect(MAX_QUANTITY).toBe(99);
  });

  it.each([1, 2, 50, 98, 99])("accepts %i", (quantity) => {
    expect(isValidQuantity(quantity)).toBe(true);
  });

  it.each([0, -1, 100, 101, 1000])("rejects %i (outside 1 to 99)", (quantity) => {
    expect(isValidQuantity(quantity)).toBe(false);
  });

  it.each([1.5, 0.999, 99.0001, Number.NaN, Number.POSITIVE_INFINITY])(
    "rejects %d (not a whole number)",
    (quantity) => {
      expect(isValidQuantity(quantity)).toBe(false);
    },
  );
});

describe("validateItems", () => {
  it("accepts a normal basket (and returns nothing)", () => {
    const basket: Basket = {
      items: [
        { sku: "APPLE", quantity: 1 },
        { sku: "COFFEE", quantity: 99 },
      ],
    };

    expect(validateItems(basket, tinyCatalog)).toBeUndefined();
  });

  it("row 1: an empty basket", () => {
    expectError({ items: [] }, "EMPTY_BASKET");
  });

  it("row 2: a sku that is not in the catalog", () => {
    expectError({ items: [{ sku: "BANANA", quantity: 1 }] }, "UNKNOWN_PRODUCT");
  });

  it("row 2: skus are case sensitive", () => {
    expectError({ items: [{ sku: "apple", quantity: 1 }] }, "UNKNOWN_PRODUCT");
  });

  it("row 2: an unknown product on any line, not only the first", () => {
    expectError(
      {
        items: [
          { sku: "APPLE", quantity: 1 },
          { sku: "PEN", quantity: 1 },
          { sku: "BANANA", quantity: 1 },
        ],
      },
      "UNKNOWN_PRODUCT",
    );
  });

  it.each([0, -1, 100, 1.5])("row 3: a quantity of %d", (quantity) => {
    expectError({ items: [{ sku: "APPLE", quantity }] }, "INVALID_QUANTITY");
  });

  it("row 3: a bad quantity on a later line", () => {
    expectError(
      {
        items: [
          { sku: "APPLE", quantity: 1 },
          { sku: "PEN", quantity: 0 },
        ],
      },
      "INVALID_QUANTITY",
    );
  });

  it("row 4: lines of one sku that add up to more than 99", () => {
    expectError(
      {
        items: [
          { sku: "APPLE", quantity: 50 },
          { sku: "APPLE", quantity: 50 },
        ],
      },
      "INVALID_QUANTITY",
    );
  });

  it("row 4: lines of one sku that add up to exactly 99 are fine", () => {
    const basket: Basket = {
      items: [
        { sku: "APPLE", quantity: 50 },
        { sku: "PEN", quantity: 99 },
        { sku: "APPLE", quantity: 49 },
      ],
    };

    expect(() => validateItems(basket, tinyCatalog)).not.toThrow();
  });

  it("row 4: 99 of one sku and 99 of another is fine (the limit is per sku)", () => {
    const basket: Basket = {
      items: [
        { sku: "APPLE", quantity: 99 },
        { sku: "PEN", quantity: 99 },
      ],
    };

    expect(() => validateItems(basket, tinyCatalog)).not.toThrow();
  });

  describe("check order", () => {
    it("an unknown product beats a bad quantity, even on a later line", () => {
      expectError(
        {
          items: [
            { sku: "APPLE", quantity: 0 },
            { sku: "BANANA", quantity: 1 },
          ],
        },
        "UNKNOWN_PRODUCT",
      );
    });

    it("a bad single quantity beats a bad merged quantity", () => {
      // Both rows 3 and 4 apply; either way the code is INVALID_QUANTITY.
      expectError(
        {
          items: [
            { sku: "APPLE", quantity: 100 },
            { sku: "APPLE", quantity: 1 },
          ],
        },
        "INVALID_QUANTITY",
      );
    });

    it("ignores the coupon and shipping method: those are checked later, by checkout()", () => {
      const basket: Basket = {
        items: [{ sku: "APPLE", quantity: 1 }],
        coupon: "NOPE",
        shipping: "overnight",
      };

      expect(() => validateItems(basket, tinyCatalog)).not.toThrow();
    });
  });
});
