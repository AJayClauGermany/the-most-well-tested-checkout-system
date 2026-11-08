// The command line program's logic (rules §8), kept free of process.stdin/stdout so it can
// be unit tested in-process. src/cli.ts is the thin wrapper that talks to the terminal.

import type { Catalog } from "./types.ts";

/** Where the standard catalog lives: spec/catalog.json at the root of the repository. */
export const STANDARD_CATALOG_PATH = new URL("../../spec/catalog.json", import.meta.url);

/** Reads spec/catalog.json. */
export function loadStandardCatalog(): Catalog {
  throw new Error("loadStandardCatalog() is not implemented yet");
}

/** What the program prints and the exit code it ends with. */
export interface CommandResult {
  output: string;
  exitCode: 0 | 1;
}

/**
 * Runs one checkout from the JSON text of a basket.
 *
 * - a valid basket gives the receipt as JSON and exit code 0,
 * - an invalid basket gives `{"error": "<CODE>"}` and exit code 1,
 * - text that is not a basket at all throws (the CLI turns that into exit code 2).
 */
export function runCheckoutCommand(input: string, catalog: Catalog): CommandResult {
  throw new Error(
    `runCheckoutCommand(${input.length}, ${catalog.products.length}) is not implemented yet`,
  );
}
