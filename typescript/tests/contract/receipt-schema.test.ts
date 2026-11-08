// Contract tests: is the output even the right shape?
//
// spec/schema/receipt.schema.json is the contract every language's receipt must honour.
// Here ajv checks real receipts against it. The e2e runner only compares totals; this is
// where the shape of the whole receipt (lines included) is held to the contract.

import { describe, expect, it } from "vitest";

import { checkout, runCheckoutCommand } from "../../src/index.ts";
import { tinyCatalog } from "../support/fixtures.ts";
import { schemaErrors } from "../support/schemas.ts";
import { loadCatalog, loadSpecCases } from "../support/spec.ts";

const catalog = loadCatalog();
const validCases = loadSpecCases().filter((specCase) => specCase.expected !== undefined);

describe("receipt.schema.json", () => {
  it.each(validCases)("the receipt for $id follows the schema", (specCase) => {
    const receipt = checkout(specCase.basket, catalog);

    expect(schemaErrors("receipt.schema.json", receipt)).toEqual([]);
  });

  it("a receipt from a made-up catalog follows the schema too", () => {
    const receipt = checkout(
      {
        items: [
          { sku: "TV", quantity: 1 },
          { sku: "COFFEE", quantity: 6 },
        ],
      },
      tinyCatalog,
    );

    expect(schemaErrors("receipt.schema.json", receipt)).toEqual([]);
  });

  it("the schema check itself works: a broken receipt is rejected (testing the test)", () => {
    const receipt = checkout({ items: [{ sku: "APPLE", quantity: 1 }] }, tinyCatalog);
    const broken = { ...receipt, shipping: 123, extra: true, tax: -1 };

    expect(schemaErrors("receipt.schema.json", broken)).toEqual([
      "/ must NOT have additional properties",
      "/shipping must be equal to one of the allowed values",
      "/tax must be >= 0",
    ]);
  });
});

describe("cli-error.schema.json", () => {
  it.each(loadSpecCases().filter((specCase) => specCase.expected_error !== undefined))(
    "the command's output for $id follows the schema",
    (specCase) => {
      const result = runCheckoutCommand(JSON.stringify(specCase.basket), catalog);

      expect(result.exitCode).toBe(1);
      expect(schemaErrors("cli-error.schema.json", JSON.parse(result.output))).toEqual([]);
    },
  );
});

describe("basket.schema.json", () => {
  it("every shared test case's basket is a valid basket (the spec is consistent with itself)", () => {
    for (const specCase of loadSpecCases()) {
      expect(schemaErrors("basket.schema.json", specCase.basket), specCase.id).toEqual([]);
    }
  });
});
