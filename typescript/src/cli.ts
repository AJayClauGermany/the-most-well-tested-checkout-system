// The command line program (rules §8): basket JSON in on stdin, receipt JSON out on stdout.
//
//   echo '{"items":[{"sku":"APPLE","quantity":1}]}' | node src/cli.ts
//
// All the logic is in command.ts; this file only connects it to the terminal, which is why
// it is left out of the in-process coverage and mutation reports and tested end-to-end instead.

import { readFileSync } from "node:fs";

import { loadStandardCatalog, runCheckoutCommand } from "./command.ts";

const input = readFileSync(0, "utf8"); // file descriptor 0 is standard input
try {
  const result = runCheckoutCommand(input, loadStandardCatalog());
  process.stdout.write(`${result.output}\n`);
  process.exitCode = result.exitCode;
} catch (error) {
  // Not a CheckoutError (those become exit code 1 above): the input was not a basket at all.
  process.stderr.write(`checkout: ${error instanceof Error ? error.message : String(error)}\n`);
  process.exitCode = 2;
}
