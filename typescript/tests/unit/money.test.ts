// Unit tests for rounding (rules §7). Rounding is where checkout systems classically go
// wrong, so it gets its own small function and its own tests.

import { describe, expect, it } from "vitest";

import { roundHalfUp } from "../../src/index.ts";

describe("roundHalfUp(amount, percent)", () => {
  // The table from rules §7, as a parameterised test: one row, one test.
  it.each([
    { amount: 224, percent: 10, exact: "22.4", expected: 22 },
    { amount: 225, percent: 10, exact: "22.5", expected: 23 },
    { amount: 226, percent: 10, exact: "22.6", expected: 23 },
  ])("rounds $exact cents to $expected", ({ amount, percent, expected }) => {
    expect(roundHalfUp(amount, percent)).toBe(expected);
  });

  it("rounds an exact half UP, where banker's rounding would go to the even number", () => {
    // 10% of $52.35 is exactly 523.5 cents. Math.round gives 524 too, but a language's
    // "round half to even" would give 524 as well here - so also check 22.5 → 23 (above),
    // where half-to-even gives 22. Both together pin the rule down.
    expect(roundHalfUp(5235, 10)).toBe(524);
  });

  it("does not round when the result is already whole", () => {
    expect(roundHalfUp(500, 10)).toBe(50);
    expect(roundHalfUp(50, 8)).toBe(4);
  });

  it("rounds 8% of $1.20 (9.6 cents) up to 10 and 8% of $2.40 (19.2 cents) down to 19", () => {
    expect(roundHalfUp(120, 8)).toBe(10);
    expect(roundHalfUp(240, 8)).toBe(19);
  });

  it("gives 0 for an amount of 0 and the whole amount for 100%", () => {
    expect(roundHalfUp(0, 10)).toBe(0);
    expect(roundHalfUp(899, 100)).toBe(899);
  });

  it("gives 0 for 0%", () => {
    expect(roundHalfUp(899, 0)).toBe(0);
  });

  it("stays exact for amounts where floating point maths drifts", () => {
    // 0.1 + 0.2 !== 0.3 in JavaScript. 10% of 3 cents is 0.3 - it must round to 0, not 1.
    expect(roundHalfUp(3, 10)).toBe(0);
    // 8% of 1050 = 84.0 exactly; a float implementation might land on 83.99999.
    expect(roundHalfUp(1050, 8)).toBe(84);
  });
});
