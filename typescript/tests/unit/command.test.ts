// Unit tests for the command line program's logic (rules §8), run in-process.
// The real program (src/cli.ts) is tested from the outside in tests/e2e/.

import { existsSync } from "node:fs";
import { dirname, join } from "node:path";

import { describe, expect, it } from "vitest";

import { loadStandardCatalog, runCheckoutCommand, STANDARD_CATALOG_PATH } from "../../src/index.ts";
import type { Catalog } from "../../src/index.ts";
import { tinyCatalog } from "../support/fixtures.ts";

describe("loadStandardCatalog", () => {
  it("points at spec/catalog.json at the root of the repository", () => {
    expect(STANDARD_CATALOG_PATH.endsWith(join("spec", "catalog.json"))).toBe(true);
    // The repository root is the folder that also holds the rules.
    expect(
      existsSync(join(dirname(dirname(STANDARD_CATALOG_PATH)), "spec", "checkout-rules.md")),
    ).toBe(true);
  });

  it("loads the six standard products", () => {
    const catalog = loadStandardCatalog();

    expect(catalog.products.map((product) => product.sku)).toEqual([
      "APPLE",
      "BREAD",
      "COFFEE",
      "NOTEBOOK",
      "PEN",
      "HEADPHONES",
    ]);
  });
});

describe("runCheckoutCommand", () => {
  it("prints the receipt as JSON and exits with 0 for a valid basket", () => {
    const result = runCheckoutCommand('{"items":[{"sku":"APPLE","quantity":1}]}', tinyCatalog);

    expect(result.exitCode).toBe(0);
    expect(JSON.parse(result.output)).toEqual({
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
    });
  });

  it('prints {"error": CODE} and exits with 1 for an invalid basket', () => {
    const result = runCheckoutCommand('{"items":[]}', tinyCatalog);

    expect(result.exitCode).toBe(1);
    expect(JSON.parse(result.output)).toEqual({ error: "EMPTY_BASKET" });
  });

  it("reports every kind of error the same way", () => {
    const result = runCheckoutCommand(
      '{"items":[{"sku":"APPLE","quantity":1}],"coupon":"NOPE"}',
      tinyCatalog,
    );

    expect(result).toEqual({ output: '{"error":"UNKNOWN_COUPON"}', exitCode: 1 });
  });

  it("prints compact JSON on a single line (one receipt, one line)", () => {
    const result = runCheckoutCommand('{"items":[{"sku":"APPLE","quantity":1}]}', tinyCatalog);

    expect(result.output).not.toContain("\n");
  });

  it("accepts surrounding whitespace, as a shell pipeline adds a newline", () => {
    const result = runCheckoutCommand('  {"items":[{"sku":"APPLE","quantity":1}]}\n', tinyCatalog);

    expect(result.exitCode).toBe(0);
  });

  it("never hides a bug behind an error code: other errors pass straight through", () => {
    // A broken catalog makes checkout() throw a TypeError, not a CheckoutError.
    const brokenCatalog = { products: undefined } as unknown as Catalog;

    expect(() =>
      runCheckoutCommand('{"items":[{"sku":"APPLE","quantity":1}]}', brokenCatalog),
    ).toThrow(TypeError);
  });

  it("throws (rather than inventing an error code) when the input is not JSON", () => {
    expect(() => runCheckoutCommand("not json", tinyCatalog)).toThrow(SyntaxError);
    expect(() => runCheckoutCommand("", tinyCatalog)).toThrow(SyntaxError);
  });

  it.each([
    ["null", "expected a JSON object"],
    ["42", "expected a JSON object"],
    ['"text"', "expected a JSON object"],
    ["[]", "expected a JSON object"],
    ["{}", "'items' must be a list"],
    ['{"items": "APPLE"}', "'items' must be a list"],
    ['{"items": [1, 2]}', "every item needs a string 'sku'"],
    ['{"items": [null]}', "every item needs a string 'sku'"],
    ['{"items": ["APPLE"]}', "every item needs a string 'sku'"],
    ['{"items": [{"sku": 5}]}', "every item needs a string 'sku'"],
  ])("throws when the JSON %s is not a basket (%s)", (input, reason) => {
    expect(() => runCheckoutCommand(input, tinyCatalog)).toThrow(TypeError);
    expect(() => runCheckoutCommand(input, tinyCatalog)).toThrow(
      `The input is not a basket: ${reason}`,
    );
  });
});
