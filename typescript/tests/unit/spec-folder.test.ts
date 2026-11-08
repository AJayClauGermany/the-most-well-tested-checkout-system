// Unit tests for finding the shared spec/ folder.

import { existsSync, mkdirSync, mkdtempSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

import { describe, expect, it } from "vitest";

import { findSpecFolder } from "../../src/index.ts";

describe("findSpecFolder", () => {
  it("finds the folder that holds catalog.json", () => {
    const folder = findSpecFolder();

    expect(folder.endsWith("spec")).toBe(true);
    expect(existsSync(join(folder, "catalog.json"))).toBe(true);
    expect(existsSync(join(folder, "checkout-rules.md"))).toBe(true);
  });

  it("finds the same folder when starting deeper inside the repository", () => {
    const folder = findSpecFolder();

    expect(findSpecFolder(join(folder, "test-cases"))).toBe(folder);
    expect(findSpecFolder(join(folder, "schema", "nested", "deeper"))).toBe(folder);
  });

  it("gives up with a clear error when there is no spec folder above the start", () => {
    expect(() => findSpecFolder(tmpdir())).toThrow("No spec/catalog.json found");
    expect(() => findSpecFolder("/")).toThrow("No spec/catalog.json found");
  });

  it("is not fooled by a spec folder without a catalog.json in it", () => {
    const folder = mkdtempSync(join(tmpdir(), "not-the-spec-"));
    mkdirSync(join(folder, "spec"));
    try {
      expect(() => findSpecFolder(join(folder, "deeper"))).toThrow("No spec/catalog.json found");
    } finally {
      rmSync(folder, { recursive: true, force: true });
    }
  });
});
