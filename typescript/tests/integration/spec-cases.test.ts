// The shared, data-driven test cases from spec/test-cases/*.json.
//
// This is a *data-driven* (or "table") test: we write the test once, and
// Vitest runs it once for every case in the JSON files. Each run is named after
// the case id, e.g. "promo-coffee-three-bags-one-free".
//
// These 60 cases are the acceptance tests of the whole project: every language runs
// exactly the same ones, so the six implementations can be compared.

import { describe, expect, it } from "vitest";

import { CheckoutError, checkout } from "../../src/index.ts";
import { loadCatalog, loadSpecCases } from "../support/spec.ts";

const catalog = loadCatalog();
const cases = loadSpecCases();

describe("shared spec cases", () => {
  // "$id" in the test name is replaced by each case's id.
  it.each(cases)("$id", (specCase) => {
    if (specCase.expected_error !== undefined) {
      // This case describes an invalid basket: checkout must throw a
      // CheckoutError whose `code` is the expected error code.
      expect(() => checkout(specCase.basket, catalog)).toThrow(CheckoutError);
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
