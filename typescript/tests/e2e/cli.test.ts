// End-to-end tests: run the real command line program (rules §8) as a separate process,
// exactly as a user (or the shared runner in e2e/) would, and look only at what comes out.
//
// Node.js runs .ts files directly since 22.18, so the program is `node src/cli.ts`.

import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

import { schemaErrors } from "../support/schemas.ts";
import { loadSpecCases } from "../support/spec.ts";

const CLI = fileURLToPath(new URL("../../src/cli.ts", import.meta.url));

/** Runs the program once, with `input` on standard input. */
function runCli(input: string): { stdout: string; stderr: string; exitCode: number | null } {
  const result = spawnSync(process.execPath, [CLI], { input, encoding: "utf8" });
  return { stdout: result.stdout, stderr: result.stderr, exitCode: result.status };
}

describe("node src/cli.ts", () => {
  // Every shared case, through the real program.
  it.each(loadSpecCases())("$id", (specCase) => {
    const result = runCli(`${JSON.stringify(specCase.basket)}\n`);

    if (specCase.expected_error !== undefined) {
      expect(result.exitCode).toBe(1);
      expect(JSON.parse(result.stdout)).toEqual({ error: specCase.expected_error });
    } else {
      expect(result.exitCode).toBe(0);
      const receipt: unknown = JSON.parse(result.stdout);
      expect(receipt).toMatchObject(specCase.expected ?? {});
      expect(schemaErrors("receipt.schema.json", receipt)).toEqual([]);
    }
  });

  it("prints exactly one line of JSON", () => {
    const result = runCli('{"items":[{"sku":"APPLE","quantity":1}]}');

    expect(result.stdout.endsWith("\n")).toBe(true);
    expect(result.stdout.trim().split("\n")).toHaveLength(1);
  });

  it("exits with 2 and explains on stderr when the input is not JSON", () => {
    const result = runCli("this is not json");

    expect(result.exitCode).toBe(2);
    expect(result.stdout).toBe("");
    expect(result.stderr).toMatch(/^checkout: /);
  });

  it("exits with 2 when the JSON is not a basket", () => {
    const result = runCli('{"hello": "world"}');

    expect(result.exitCode).toBe(2);
    expect(result.stderr).toContain("not a basket");
  });

  it("exits with 2 on empty input", () => {
    expect(runCli("").exitCode).toBe(2);
  });
});
