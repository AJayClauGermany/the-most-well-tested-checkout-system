// Snapshot tests: did the output change since last time?
//
// toMatchSnapshot() saves the value to tests/snapshot/__snapshots__/ on the first run
// and compares against that file on later runs. toMatchInlineSnapshot() keeps the saved
// value right here in the test. If a change is intended, run `npx vitest -u` to update.
//
// The hand-checked receipt in spec/examples/receipt.json plays the role of an
// "approved" output shared by all six languages.

import { readFileSync } from "node:fs";
import { join } from "node:path";

import { describe, expect, it } from "vitest";

import { checkout, runCheckoutCommand } from "../../src/index.ts";
import type { Basket } from "../../src/index.ts";
import { loadCatalog, SPEC_DIR } from "../support/spec.ts";

const catalog = loadCatalog();

/** The basket behind spec/examples/receipt.json (rules §5). */
const exampleBasket: Basket = {
  items: [
    { sku: "COFFEE", quantity: 3 },
    { sku: "APPLE", quantity: 2 },
  ],
  coupon: "SAVE10",
  shipping: "standard",
};

describe("the example receipt from the rules", () => {
  it("matches the hand-checked spec/examples/receipt.json, lines included", () => {
    const approved: unknown = JSON.parse(
      readFileSync(join(SPEC_DIR, "examples", "receipt.json"), "utf8"),
    );

    expect(checkout(exampleBasket, catalog)).toEqual(approved);
  });

  it("matches the saved snapshot", () => {
    expect(checkout(exampleBasket, catalog)).toMatchSnapshot();
  });

  it("is printed by the command exactly like this", () => {
    expect(runCheckoutCommand(JSON.stringify(exampleBasket), catalog).output).toMatchSnapshot();
  });
});

describe("the biggest shared case, full-everything", () => {
  it("matches the saved snapshot (every product, both promotions, WELCOME5, express)", () => {
    const basket: Basket = {
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

    expect(checkout(basket, catalog)).toMatchSnapshot();
  });
});

describe("inline snapshots", () => {
  it("one pen, written out as the receipt lines", () => {
    const receipt = checkout({ items: [{ sku: "PEN", quantity: 1 }] }, catalog);

    expect(receipt.lines).toMatchInlineSnapshot(`
      [
        {
          "line_discount": 0,
          "line_total": 120,
          "name": "Pen",
          "quantity": 1,
          "sku": "PEN",
          "unit_price": 120,
        },
      ]
    `);
  });

  it("an error, as the command prints it", () => {
    expect(runCheckoutCommand('{"items":[]}', catalog)).toMatchInlineSnapshot(`
      {
        "exitCode": 1,
        "output": "{"error":"EMPTY_BASKET"}",
      }
    `);
  });
});
