// Unit tests for tax (rules §4.6).

import { describe, expect, it } from "vitest";

import { TAX_PERCENT, taxOn } from "../../src/index.ts";

describe("taxOn", () => {
  it("is 8 percent", () => {
    expect(TAX_PERCENT).toBe(8);
  });

  it.each([
    { amount: 50, tax: 4, why: "8% of $0.50 is exactly 4 cents" },
    { amount: 120, tax: 10, why: "9.6 cents rounds up" },
    { amount: 240, tax: 19, why: "19.2 cents rounds down" },
    { amount: 625, tax: 50, why: "exactly 50 cents" },
    { amount: 1706, tax: 136, why: "136.48 rounds down" },
    { amount: 11880, tax: 950, why: "950.4 rounds down" },
    { amount: 0, tax: 0, why: "no spend, no tax" },
  ])("taxes $amount cents as $tax cents ($why)", ({ amount, tax }) => {
    expect(taxOn(amount)).toBe(tax);
  });

  it("rounds plain fractions the usual way: 150.48 → 150, 151.04 → 151", () => {
    // 8% of a whole number of cents can never be exactly half a cent (x × 8 never ends in 50),
    // so the half-up rule is pinned down by the 10% coupon tests; here we check plain rounding.
    expect(taxOn(1875)).toBe(150); // exactly 150.00
    expect(taxOn(1881)).toBe(150); // 150.48
    expect(taxOn(1888)).toBe(151); // 151.04
  });
});
