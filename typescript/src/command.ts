// The command line program's logic (rules §8), kept free of process.stdin/stdout so it can
// be unit tested in-process. src/cli.ts is the thin wrapper that talks to the terminal.

import { readFileSync } from "node:fs";
import { join } from "node:path";

import { checkout } from "./checkout.ts";
import { CheckoutError } from "./errors.ts";
import { findSpecFolder } from "./spec-folder.ts";
import type { Basket, Catalog } from "./types.ts";

/** Where the standard catalog lives: spec/catalog.json at the root of the repository. */
export const STANDARD_CATALOG_PATH = join(findSpecFolder(), "catalog.json");

/** Reads spec/catalog.json. */
export function loadStandardCatalog(): Catalog {
  return JSON.parse(readFileSync(STANDARD_CATALOG_PATH, "utf8")) as Catalog;
}

/** What the program prints and the exit code it ends with. */
export interface CommandResult {
  output: string;
  exitCode: 0 | 1;
}

/**
 * Parses the JSON text of a basket. Only the shape that matters is checked - `items` must
 * be a list of objects with a string `sku` - so that checkout() can do the real validation.
 *
 * @throws {SyntaxError} when the text is not JSON; {TypeError} when the JSON is not a basket.
 */
function parseBasket(input: string): Basket {
  const parsed: unknown = JSON.parse(input);
  if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
    throw new TypeError("The input is not a basket: expected a JSON object");
  }
  const { items } = parsed as { items?: unknown };
  if (!Array.isArray(items)) {
    throw new TypeError("The input is not a basket: 'items' must be a list");
  }
  for (const item of items as Array<{ sku?: unknown } | null>) {
    // `?.` makes this safe for null, numbers and strings too: their `sku` is simply undefined.
    if (typeof item?.sku !== "string") {
      throw new TypeError("The input is not a basket: every item needs a string 'sku'");
    }
  }
  return parsed as Basket;
}

/**
 * Runs one checkout from the JSON text of a basket.
 *
 * - a valid basket gives the receipt as JSON and exit code 0,
 * - an invalid basket gives `{"error": "<CODE>"}` and exit code 1,
 * - text that is not a basket at all throws (the CLI turns that into exit code 2).
 */
export function runCheckoutCommand(input: string, catalog: Catalog): CommandResult {
  const basket = parseBasket(input);
  try {
    return { output: JSON.stringify(checkout(basket, catalog)), exitCode: 0 };
  } catch (error) {
    if (error instanceof CheckoutError) {
      return { output: JSON.stringify({ error: error.code }), exitCode: 1 };
    }
    throw error;
  }
}
