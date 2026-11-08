// Unit tests of checkout() as a whole, with the tiny made-up catalog from tests/support/.
// The 60 shared spec cases (tests/integration/spec-cases.test.ts) cover the standard catalog;
// these tests check the pieces fit together and the things the totals alone cannot show,
// such as the receipt lines.

import { describe, expect, it } from "vitest";

import { CheckoutError, checkout } from "../../src/index.ts";
import type { Basket, ErrorCode, Receipt } from "../../src/index.ts";
import { tinyCatalog } from "../support/fixtures.ts";

function expectError(basket: Basket, code: ErrorCode): void {
  expect(() => checkout(basket, tinyCatalog)).toThrow(CheckoutError);
  expect(() => checkout(basket, tinyCatalog)).toThrow(expect.objectContaining({ code }));
}

describe("checkout", () => {
  it("prices one apple: 50 + $4.99 shipping + 4 tax = $5.53", () => {
    const receipt = checkout({ items: [{ sku: "APPLE", quantity: 1 }] }, tinyCatalog);

    const expected: Receipt = {
      lines: [
        {
          sku: "APPLE",
          name: "Apple",
          quantity: 1,
          unit_price: 50,
          line_total: 50,
          line_discount: 0,
        },
      ],
      subtotal: 50,
      promotion_discount: 0,
      coupon_discount: 0,
      shipping: 499,
      tax: 4,
      total: 553,
    };
    expect(receipt).toEqual(expected);
  });

  it("writes one receipt line per distinct sku, in order of first appearance, with merged quantities", () => {
    const receipt = checkout(
      {
        items: [
          { sku: "PEN", quantity: 1 },
          { sku: "APPLE", quantity: 2 },
          { sku: "PEN", quantity: 3 },
        ],
      },
      tinyCatalog,
    );

    expect(receipt.lines.map((line) => [line.sku, line.quantity])).toEqual([
      ["PEN", 4],
      ["APPLE", 2],
    ]);
  });

  it("works the totals out from the lines", () => {
    // 3 coffees ($26.97, one free: -$8.99) and 10 notebooks ($34.90, 10% off: -$3.49)
    const receipt = checkout(
      {
        items: [
          { sku: "COFFEE", quantity: 3 },
          { sku: "NOTEBOOK", quantity: 10 },
        ],
      },
      tinyCatalog,
    );

    expect(receipt.subtotal).toBe(2697 + 3490);
    expect(receipt.promotion_discount).toBe(899 + 349);
    expect(receipt.coupon_discount).toBe(0);
    expect(receipt.shipping).toBe(499); // $49.39 after promotions: not free
    expect(receipt.tax).toBe(395); // 8% of 4939 = 395.12
    expect(receipt.total).toBe(4939 + 499 + 395);
  });

  it("takes the coupon off after promotions and before shipping and tax", () => {
    // TV $100.00 with SAVE10: $90.00 → free shipping, tax 8% of $90.00 = $7.20
    const receipt = checkout(
      { items: [{ sku: "TV", quantity: 1 }], coupon: "SAVE10" },
      tinyCatalog,
    );

    expect(receipt).toMatchObject({
      subtotal: 10000,
      promotion_discount: 0,
      coupon_discount: 1000,
      shipping: 0,
      tax: 720,
      total: 9720,
    });
  });

  it("charges express shipping on top, untaxed", () => {
    const receipt = checkout(
      { items: [{ sku: "TV", quantity: 1 }], shipping: "express" },
      tinyCatalog,
    );

    expect(receipt).toMatchObject({ shipping: 999, tax: 800, total: 10000 + 999 + 800 });
  });

  it("does not change the basket it is given", () => {
    const basket: Basket = {
      items: [
        { sku: "APPLE", quantity: 1 },
        { sku: "APPLE", quantity: 1 },
      ],
      coupon: " save10 ",
    };

    checkout(basket, tinyCatalog);

    expect(basket).toEqual({
      items: [
        { sku: "APPLE", quantity: 1 },
        { sku: "APPLE", quantity: 1 },
      ],
      coupon: " save10 ",
    });
  });

  describe("reports the first failing row of rules §4.1", () => {
    it("empty basket before anything else", () => {
      expectError({ items: [], coupon: "NOPE", shipping: "overnight" }, "EMPTY_BASKET");
    });

    it("unknown product before bad quantity, shipping and coupon", () => {
      expectError(
        {
          items: [
            { sku: "APPLE", quantity: 0 },
            { sku: "BANANA", quantity: 1 },
          ],
          coupon: "NOPE",
          shipping: "x",
        },
        "UNKNOWN_PRODUCT",
      );
    });

    it("bad quantity before shipping and coupon", () => {
      expectError(
        { items: [{ sku: "APPLE", quantity: 0 }], coupon: "NOPE", shipping: "x" },
        "INVALID_QUANTITY",
      );
    });

    it("merged quantity over 99 before shipping and coupon", () => {
      expectError(
        {
          items: [
            { sku: "APPLE", quantity: 60 },
            { sku: "APPLE", quantity: 40 },
          ],
          coupon: "NOPE",
          shipping: "x",
        },
        "INVALID_QUANTITY",
      );
    });

    it("bad shipping before bad coupon", () => {
      expectError(
        { items: [{ sku: "APPLE", quantity: 1 }], coupon: "NOPE", shipping: "x" },
        "UNKNOWN_SHIPPING_METHOD",
      );
    });

    it("bad shipping before an unmet coupon minimum", () => {
      expectError(
        { items: [{ sku: "APPLE", quantity: 1 }], coupon: "WELCOME5", shipping: "x" },
        "UNKNOWN_SHIPPING_METHOD",
      );
    });

    it("unknown coupon", () => {
      expectError({ items: [{ sku: "APPLE", quantity: 1 }], coupon: "NOPE" }, "UNKNOWN_COUPON");
    });

    it("coupon minimum not met, judged on the price after promotions", () => {
      // 3 coffees are $26.97 but $17.98 after the free one: below WELCOME5's $25.00.
      expectError(
        { items: [{ sku: "COFFEE", quantity: 3 }], coupon: "WELCOME5" },
        "COUPON_MINIMUM_NOT_MET",
      );
    });
  });
});
