// A micro-benchmark (run with `npm run bench`). The checkout is tiny, so one benchmark is
// enough to show the tool. Since Vitest 5 a benchmark is registered from inside a test, through
// the test context's `bench`: it runs the function many times and reports operations per second,
// and `bench.compare` prints the two side by side. This is not a performance project - see
// docs/testing-matrix.md.

import { describe, test } from "vitest";

import { checkout } from "../../src/index.ts";
import type { Basket } from "../../src/index.ts";
import { loadCatalog } from "../support/spec.ts";

const catalog = loadCatalog();

const oneApple: Basket = { items: [{ sku: "APPLE", quantity: 1 }] };

const everything: Basket = {
  items: [
    { sku: "HEADPHONES", quantity: 1 },
    { sku: "COFFEE", quantity: 6 },
    { sku: "NOTEBOOK", quantity: 12 },
    { sku: "PEN", quantity: 5 },
    { sku: "BREAD", quantity: 1 },
    { sku: "APPLE", quantity: 3 },
  ],
  coupon: "WELCOME5",
  shipping: "express",
};

describe("checkout", () => {
  test("a small and a big basket", async ({ bench }) => {
    await bench.compare(
      bench("one apple", () => {
        checkout(oneApple, catalog);
      }),
      bench("every product, both promotions, a coupon and express shipping", () => {
        checkout(everything, catalog);
      }),
    );
  });
});
