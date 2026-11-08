// Test helper: loads the shared, language-independent spec files from the
// repository's top-level spec/ folder. We read the files where they are
// (we never copy them), so every language is tested against the same data.

import { readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

import { findSpecFolder } from "../../src/index.ts";
import type { Basket, Catalog, ErrorCode } from "../../src/index.ts";

/**
 * The absolute path of the spec/ folder, with a trailing slash.
 *
 * We find it relative to *this file* (`import.meta.url`), not relative to the current
 * working directory, so it works no matter which folder you run the tests from. The
 * search itself (walk up until a folder holds spec/catalog.json) is shared with the
 * command line program, see src/spec-folder.ts.
 */
export const SPEC_DIR = `${findSpecFolder(fileURLToPath(new URL(".", import.meta.url)))}/`;

/** The six totals a shared test case expects (all in cents). */
export interface ExpectedTotals {
  subtotal: number;
  promotion_discount: number;
  coupon_discount: number;
  shipping: number;
  tax: number;
  total: number;
}

/**
 * One case from spec/test-cases/*.json.
 * A case has EITHER `expected` (the totals) OR `expected_error` (an error code).
 */
export interface SpecCase {
  /** The file the case came from, e.g. "01-subtotal.json". */
  file: string;
  id: string;
  description: string;
  basket: Basket;
  expected?: ExpectedTotals;
  expected_error?: ErrorCode;
}

/** Reads and parses a JSON file. */
function readJson(path: string): unknown {
  return JSON.parse(readFileSync(path, "utf8"));
}

// Note on `as` below: `JSON.parse` cannot know what is inside a file, so we tell
// TypeScript the shape ourselves. That is safe enough here because CI checks every
// spec file against the JSON Schemas in spec/schema/.

/** Loads spec/catalog.json, the standard catalog used by every shared test case. */
export function loadCatalog(): Catalog {
  return readJson(join(SPEC_DIR, "catalog.json")) as Catalog;
}

/** Loads every case from every spec/test-cases/*.json file, in file order. */
export function loadSpecCases(): SpecCase[] {
  const folder = join(SPEC_DIR, "test-cases");
  const files = readdirSync(folder)
    .filter((name) => name.endsWith(".json"))
    .sort();

  const cases: SpecCase[] = [];
  for (const file of files) {
    const content = readJson(join(folder, file)) as { cases: Omit<SpecCase, "file">[] };
    for (const specCase of content.cases) {
      cases.push({ file, ...specCase });
    }
  }
  return cases;
}
