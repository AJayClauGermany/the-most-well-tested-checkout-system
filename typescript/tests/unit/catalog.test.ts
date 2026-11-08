// EXAMPLE UNIT TEST - start reading here if you are new to Vitest!
//
// A *unit test* checks one small piece of code (a "unit") on its own:
// no files, no network, no other parts of the program. Unit tests are fast
// and, when one fails, it points straight at the problem.
//
// Vitest gives us three main tools:
//   describe(name, fn)  groups related tests (optional, but it keeps output tidy)
//   it(name, fn)        one test; the name should read like a sentence
//   expect(value)       an assertion: expect(actual).toBe(expected)

import { describe, expect, it } from "vitest";

import { findProduct } from "../../src/index.ts";
import type { Catalog } from "../../src/index.ts";

// A tiny catalog made up just for these tests. Small, hand-made test data makes
// it easy to see why a test passes or fails.
const catalog: Catalog = {
  products: [
    { sku: "APPLE", name: "Apple", unit_price: 50, promotion: null },
    {
      sku: "COFFEE",
      name: "Bag of Coffee Beans",
      unit_price: 899,
      promotion: { type: "buy_x_get_y_free", buy: 2, free: 1 },
    },
  ],
};

describe("findProduct", () => {
  it("finds a product by its sku", () => {
    // Arrange: get everything ready (here, the catalog above is all we need).
    const sku = "COFFEE";

    // Act: run the code we are testing.
    const product = findProduct(catalog, sku);

    // Assert: check the result.
    expect(product?.name).toBe("Bag of Coffee Beans");
    expect(product?.unit_price).toBe(899);
  });

  it("returns undefined for a sku that is not in the catalog", () => {
    expect(findProduct(catalog, "BANANA")).toBeUndefined();
  });

  it("matches skus exactly, so lower case 'apple' is not found", () => {
    // rules §4: "SKUs are matched exactly - they are case sensitive".
    expect(findProduct(catalog, "apple")).toBeUndefined();
  });
});
