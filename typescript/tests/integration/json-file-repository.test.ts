// Integration tests for the JSON-file order repository: real files in a real temporary folder.
//
// Unlike a unit test, this touches the file system, so each test gets its own fresh folder
// (os.tmpdir() + mkdtempSync) and removes it afterwards. The stored files must follow
// spec/schema/order.schema.json - a contract test inside an integration test.

import { existsSync, mkdtempSync, readdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { checkout, JsonFileOrderRepository } from "../../src/index.ts";
import { oneAppleReceipt, tinyCatalog } from "../support/fixtures.ts";
import { schemaErrors } from "../support/schemas.ts";

const fixedClock = () => new Date("2026-01-15T10:30:00Z");

let folder: string;

beforeEach(() => {
  folder = mkdtempSync(join(tmpdir(), "checkout-orders-"));
});

afterEach(() => {
  rmSync(folder, { recursive: true, force: true });
});

describe("JsonFileOrderRepository", () => {
  it("writes ORDER-1.json into the folder on the first save", () => {
    const repository = new JsonFileOrderRepository(folder, fixedClock);

    const id = repository.save(oneAppleReceipt);

    expect(id).toBe("ORDER-1");
    expect(existsSync(join(folder, "ORDER-1.json"))).toBe(true);
  });

  it("writes a file that follows order.schema.json exactly", () => {
    const repository = new JsonFileOrderRepository(folder, fixedClock);
    repository.save(oneAppleReceipt);

    const stored: unknown = JSON.parse(readFileSync(join(folder, "ORDER-1.json"), "utf8"));

    expect(schemaErrors("order.schema.json", stored)).toEqual([]);
    expect(stored).toEqual({
      id: "ORDER-1",
      created_at: "2026-01-15T10:30:00Z",
      receipt: oneAppleReceipt,
    });
  });

  it("writes human-readable JSON (indented, ending in a newline)", () => {
    const repository = new JsonFileOrderRepository(folder, fixedClock);
    repository.save(oneAppleReceipt);

    const text = readFileSync(join(folder, "ORDER-1.json"), "utf8");

    expect(text).toContain("\n  ");
    expect(text.endsWith("\n")).toBe(true);
  });

  it("reads a saved order back, exactly as saved", () => {
    const repository = new JsonFileOrderRepository(folder, fixedClock);
    const receipt = checkout(
      { items: [{ sku: "COFFEE", quantity: 3 }], coupon: "SAVE10" },
      tinyCatalog,
    );

    const id = repository.save(receipt);

    expect(repository.find(id)).toEqual({
      id: "ORDER-1",
      created_at: "2026-01-15T10:30:00Z",
      receipt,
    });
  });

  it("numbers orders sequentially", () => {
    const repository = new JsonFileOrderRepository(folder, fixedClock);

    const ids = [
      repository.save(oneAppleReceipt),
      repository.save(oneAppleReceipt),
      repository.save(oneAppleReceipt),
    ];

    expect(ids).toEqual(["ORDER-1", "ORDER-2", "ORDER-3"]);
    expect(readdirSync(folder).sort()).toEqual(["ORDER-1.json", "ORDER-2.json", "ORDER-3.json"]);
  });

  it("continues numbering when the same folder is opened again later", () => {
    const first = new JsonFileOrderRepository(folder, fixedClock);
    first.save(oneAppleReceipt);
    first.save(oneAppleReceipt);

    const second = new JsonFileOrderRepository(folder, fixedClock);

    expect(second.save(oneAppleReceipt)).toBe("ORDER-3");
    expect(second.find("ORDER-1")?.id).toBe("ORDER-1");
  });

  it("numbers after the highest existing order, even with gaps and 2-digit numbers", () => {
    writeFileSync(
      join(folder, "ORDER-3.json"),
      JSON.stringify({
        id: "ORDER-3",
        created_at: "2026-01-15T10:30:00Z",
        receipt: oneAppleReceipt,
      }),
    );
    writeFileSync(
      join(folder, "ORDER-12.json"),
      JSON.stringify({
        id: "ORDER-12",
        created_at: "2026-01-15T10:30:00Z",
        receipt: oneAppleReceipt,
      }),
    );
    const repository = new JsonFileOrderRepository(folder, fixedClock);

    expect(repository.save(oneAppleReceipt)).toBe("ORDER-13"); // not ORDER-4 (alphabetical would say 3 > 12)
  });

  it("ignores files in the folder that are not orders", () => {
    writeFileSync(join(folder, "README.txt"), "not an order");
    writeFileSync(join(folder, "ORDER-9.json.bak"), "not an order either");
    const repository = new JsonFileOrderRepository(folder, fixedClock);

    expect(repository.save(oneAppleReceipt)).toBe("ORDER-1");
  });

  it("creates the folder if it does not exist yet", () => {
    const nested = join(folder, "orders", "2026");
    const repository = new JsonFileOrderRepository(nested, fixedClock);

    repository.save(oneAppleReceipt);

    expect(existsSync(join(nested, "ORDER-1.json"))).toBe(true);
  });

  it("answers 'not found' for ids that have no file, without throwing", () => {
    const repository = new JsonFileOrderRepository(folder, fixedClock);
    repository.save(oneAppleReceipt);

    expect(repository.find("ORDER-2")).toBeUndefined();
    expect(repository.find("ORDER-0")).toBeUndefined();
    expect(repository.find("nonsense")).toBeUndefined();
  });

  it("does not let an id escape the folder (no path tricks)", () => {
    writeFileSync(
      join(folder, "secret.json"),
      JSON.stringify({ id: "x", created_at: "x", receipt: {} }),
    );
    const repository = new JsonFileOrderRepository(join(folder, "orders"), fixedClock);

    expect(repository.find("../secret")).toBeUndefined();
    expect(repository.find("../secret.json")).toBeUndefined();
  });

  it("stamps each order with the clock's time at the moment of saving", () => {
    let now = new Date("2026-01-15T10:30:00Z");
    const repository = new JsonFileOrderRepository(folder, () => now);

    const first = repository.save(oneAppleReceipt);
    now = new Date("2026-01-16T09:00:00Z");
    const second = repository.save(oneAppleReceipt);

    expect(repository.find(first)?.created_at).toBe("2026-01-15T10:30:00Z");
    expect(repository.find(second)?.created_at).toBe("2026-01-16T09:00:00Z");
  });
});
