// The shared, data-driven test cases from spec/test-cases/*.json.
//
// This is a *data-driven* (or "table") test: we write the test once, and
// Vitest runs it once for every case in the JSON files. Each run is named after
// the case id, e.g. "promo-coffee-three-bags-one-free".
//
// TDD in action: these tests describe the finished checkout, but checkout()
// is not written yet, so they are SKIPPED for now. They still type-check, so
// they are ready and waiting. To start implementing:
//   1. change `it.skip.each` below to `it.each` (and drop "skipped: ..." from the title),
//   2. run `npm test` and watch all 60 tests fail (red),
//   3. write code in src/ until they pass (green), one case at a time.

import { describe, expect, it } from "vitest";

import { checkout } from "../../src/index.ts";
import { loadCatalog, loadSpecCases } from "../support/spec.ts";

const catalog = loadCatalog();
const cases = loadSpecCases();

describe("shared spec cases (skipped: checkout() is not implemented yet)", () => {
  // SKIPPED: checkout() is not implemented yet - remove this skip when you start implementing.
  // "$id" in the test name is replaced by each case's id.
  it.skip.each(cases)("$id", (specCase) => {
    if (specCase.expected_error !== undefined) {
      // This case describes an invalid basket: checkout must throw a
      // CheckoutError whose `code` is the expected error code.
      expect(() => checkout(specCase.basket, catalog)).toThrow(
        expect.objectContaining({ code: specCase.expected_error }),
      );
    } else {
      // This case describes a valid basket: compare the six totals.
      // toMatchObject only checks the fields we list, so the receipt's
      // `lines` are ignored here.
      const receipt = checkout(specCase.basket, catalog);
      expect(receipt).toMatchObject(specCase.expected ?? {});
    }
  });
});
