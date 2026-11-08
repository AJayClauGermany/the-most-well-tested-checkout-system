// Finding the shared spec/ folder at the root of the repository.

import { existsSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

/**
 * The repository's `spec/` folder: walks up from `start` (by default, this very file) until it
 * reaches a folder that contains `spec/catalog.json`. Looking upwards (instead of hard-coding `../../spec`) means
 * the code still finds the spec when it is run from a copy somewhere deeper inside the
 * repository - which is exactly what a mutation testing tool does with it.
 *
 * @throws {Error} when no folder above `start` contains `spec/catalog.json`.
 */
export function findSpecFolder(start: string = fileURLToPath(import.meta.url)): string {
  let folder = start;
  while (!existsSync(join(folder, "spec", "catalog.json"))) {
    const parent = dirname(folder);
    if (parent === folder) {
      throw new Error(`No spec/catalog.json found in any folder above ${start}`);
    }
    folder = parent;
  }
  return join(folder, "spec");
}
