// Unit tests for merging basket items into receipt lines and pricing them (rules §4.2).

import { describe, expect, it } from "vitest";

import { mergeItems, priceLine } from "../../src/index.ts";
import { APPLE, COFFEE, NOTEBOOK } from "../support/fixtures.ts";

describe("mergeItems", () => {
  it("keeps items with different skus apart, in basket order", () => {
    const merged = mergeItems([
      { sku: "PEN", quantity: 2 },
      { sku: "APPLE", quantity: 1 },
    ]);

    expect(merged).toEqual([
      { sku: "PEN", quantity: 2 },
      { sku: "APPLE", quantity: 1 },
    ]);
  });

  it("adds up the quantities of the same sku, keeping the place of its first appearance", () => {
    // The example from rules §4.2: [APPLE x1, PEN x2, APPLE x3] gives [APPLE x4, PEN x2].
    const merged = mergeItems([
      { sku: "APPLE", quantity: 1 },
      { sku: "PEN", quantity: 2 },
      { sku: "APPLE", quantity: 3 },
    ]);

    expect(merged).toEqual([
      { sku: "APPLE", quantity: 4 },
      { sku: "PEN", quantity: 2 },
    ]);
  });

  it("merges three lines of one sku into one", () => {
    const merged = mergeItems([
      { sku: "COFFEE", quantity: 1 },
      { sku: "COFFEE", quantity: 1 },
      { sku: "COFFEE", quantity: 1 },
    ]);

    expect(merged).toEqual([{ sku: "COFFEE", quantity: 3 }]);
  });

  it("gives an empty list for an empty list", () => {
    expect(mergeItems([])).toEqual([]);
  });

  it("does not change the basket's own items", () => {
    const items = [
      { sku: "APPLE", quantity: 1 },
      { sku: "APPLE", quantity: 1 },
    ];

    mergeItems(items);

    expect(items).toEqual([
      { sku: "APPLE", quantity: 1 },
      { sku: "APPLE", quantity: 1 },
    ]);
  });
});

describe("priceLine", () => {
  it("copies sku, name and unit price from the product and multiplies for the line total", () => {
    expect(priceLine(APPLE, 3)).toEqual({
      sku: "APPLE",
      name: "Apple",
      quantity: 3,
      unit_price: 50,
      line_total: 150,
      line_discount: 0,
    });
  });

  it("includes the promotion discount", () => {
    expect(priceLine(COFFEE, 3)).toMatchObject({ line_total: 2697, line_discount: 899 });
    expect(priceLine(NOTEBOOK, 10)).toMatchObject({ line_total: 3490, line_discount: 349 });
  });

  it("gives no discount below the promotion's threshold", () => {
    expect(priceLine(COFFEE, 2).line_discount).toBe(0);
    expect(priceLine(NOTEBOOK, 9).line_discount).toBe(0);
  });
});
