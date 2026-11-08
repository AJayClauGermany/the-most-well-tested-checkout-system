// Unit tests for the two kinds of promotion (rules §4.2).

import { describe, expect, it } from "vitest";

import { promotionDiscount } from "../../src/index.ts";
import type { Promotion } from "../../src/index.ts";

const buyTwoGetOne: Promotion = { type: "buy_x_get_y_free", buy: 2, free: 1 };
const tenPercentOffTen: Promotion = { type: "bulk_discount", min_quantity: 10, percent: 10 };

describe("promotionDiscount", () => {
  it("is 0 when the product has no promotion", () => {
    expect(promotionDiscount(null, 50, 899)).toBe(0);
  });

  describe("buy X, get Y free", () => {
    // Coffee is buy 2 get 1 free at $8.99.
    it.each([
      { quantity: 1, free: 0 },
      { quantity: 2, free: 0 },
      { quantity: 3, free: 1 },
      { quantity: 4, free: 1 },
      { quantity: 5, free: 1 },
      { quantity: 6, free: 2 },
      { quantity: 99, free: 33 },
    ])("$quantity coffees → $free free", ({ quantity, free }) => {
      expect(promotionDiscount(buyTwoGetOne, quantity, 899)).toBe(free * 899);
    });

    it("uses the promotion's own numbers, not 2 and 1", () => {
      // Buy 1 get 1 free: every second item is free.
      expect(promotionDiscount({ type: "buy_x_get_y_free", buy: 1, free: 1 }, 4, 100)).toBe(200);
      expect(promotionDiscount({ type: "buy_x_get_y_free", buy: 1, free: 1 }, 5, 100)).toBe(200);
      // Buy 3 get 2 free: groups of 5, 2 of each group free.
      expect(promotionDiscount({ type: "buy_x_get_y_free", buy: 3, free: 2 }, 4, 100)).toBe(0);
      expect(promotionDiscount({ type: "buy_x_get_y_free", buy: 3, free: 2 }, 5, 100)).toBe(200);
      expect(promotionDiscount({ type: "buy_x_get_y_free", buy: 3, free: 2 }, 10, 100)).toBe(400);
    });

    it("multiplies by the unit price", () => {
      expect(promotionDiscount(buyTwoGetOne, 3, 100)).toBe(100);
      expect(promotionDiscount(buyTwoGetOne, 3, 1)).toBe(1);
    });
  });

  describe("bulk discount", () => {
    // Notebooks are 10% off from 10 at $3.49.
    it.each([
      { quantity: 9, discount: 0, why: "one short of the minimum" },
      { quantity: 10, discount: 349, why: "exactly the minimum: 10% of $34.90" },
      { quantity: 11, discount: 384, why: "383.9 cents rounds up" },
      { quantity: 15, discount: 524, why: "exactly 523.5 cents rounds UP" },
      { quantity: 99, discount: 3455, why: "10% of $345.51 is 3455.1" },
    ])("$quantity notebooks → $discount cents off ($why)", ({ quantity, discount }) => {
      expect(promotionDiscount(tenPercentOffTen, quantity, 349)).toBe(discount);
    });

    it("uses the promotion's own minimum and percentage", () => {
      const quarterOffFromThree: Promotion = {
        type: "bulk_discount",
        min_quantity: 3,
        percent: 25,
      };
      expect(promotionDiscount(quarterOffFromThree, 2, 100)).toBe(0);
      expect(promotionDiscount(quarterOffFromThree, 3, 100)).toBe(75);
      expect(promotionDiscount(quarterOffFromThree, 3, 101)).toBe(76); // 75.75 → 76
    });
  });
});
