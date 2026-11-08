// Tests for the spec loader in tests/support/spec.ts.
//
// These are *integration tests*: instead of testing a function on its own with
// made-up data, they check that our code works together with something real -
// here, the actual JSON files on disk in the shared spec/ folder.

import { describe, expect, it } from "vitest";

import { findProduct } from "../../src/index.ts";
import { loadCatalog, loadSpecCases } from "../support/spec.ts";

describe("loadCatalog", () => {
  it("loads the six products of the standard catalog", () => {
    const catalog = loadCatalog();

    expect(catalog.products).toHaveLength(6);
  });

  it("loads COFFEE with its 'buy 2, get 1 free' promotion", () => {
    const coffee = findProduct(loadCatalog(), "COFFEE");

    expect(coffee?.promotion).toEqual({ type: "buy_x_get_y_free", buy: 2, free: 1 });
  });
});

describe("loadSpecCases", () => {
  const cases = loadSpecCases();

  it("loads all 60 shared test cases", () => {
    // If you add a case to spec/test-cases/, update this number too.
    expect(cases).toHaveLength(60);
  });

  it("gives every case a unique id", () => {
    const ids = cases.map((specCase) => specCase.id);
    const uniqueIds = new Set(ids);

    expect(uniqueIds.size).toBe(ids.length);
  });

  it("gives every case exactly one of 'expected' or 'expected_error'", () => {
    for (const specCase of cases) {
      const hasExpected = specCase.expected !== undefined;
      const hasExpectedError = specCase.expected_error !== undefined;

      // The second argument to expect() is a message shown if the check fails.
      expect(hasExpected !== hasExpectedError, `case ${specCase.id}`).toBe(true);
    }
  });
});
