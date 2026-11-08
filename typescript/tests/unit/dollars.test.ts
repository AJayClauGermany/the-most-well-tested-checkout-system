// Test code deserves tests too: the BDD steps and the e2e tests rely on dollarsToCents,
// so a bug there would hide bugs in the checkout.

import { describe, expect, it } from "vitest";

import { dollarsToCents } from "../support/dollars.ts";

describe("dollarsToCents", () => {
  // A parameterised test: the same check runs once per row of the table.
  it.each([
    ["$0.00", 0],
    ["$0.50", 50],
    ["$12.39", 1239],
    ["$5", 500],
    ["$5.5", 550],
    ["  $8.99 ", 899],
    ["0.04", 4],
  ])("turns %s into %i cents", (text, cents) => {
    expect(dollarsToCents(text)).toBe(cents);
  });

  it("refuses text that is not an amount", () => {
    expect(() => dollarsToCents("ten dollars")).toThrow("not a dollar amount");
    expect(() => dollarsToCents("$1.234")).toThrow("not a dollar amount");
  });
});
