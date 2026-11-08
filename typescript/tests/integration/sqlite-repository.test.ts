// Integration tests for the SQLite order repository, using Node's built-in node:sqlite.
//
// A real database, but no server to install: SQLite is one file. Each test gets its own
// file in a temporary folder and closes the database afterwards.

import { mkdtempSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { DatabaseSync } from "node:sqlite";

import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { checkout, SqliteOrderRepository } from "../../src/index.ts";
import { oneAppleReceipt, tinyCatalog } from "../support/fixtures.ts";
import { schemaErrors } from "../support/schemas.ts";

const fixedClock = () => new Date("2026-01-15T10:30:00Z");

let folder: string;
let path: string;
const open: SqliteOrderRepository[] = [];

/** Opens a repository and remembers it, so afterEach can close it. */
function openRepository(clock = fixedClock): SqliteOrderRepository {
  const repository = new SqliteOrderRepository(path, clock);
  open.push(repository);
  return repository;
}

/** Closes a repository early, inside a test (and forgets it, so afterEach does not close it twice). */
function closeNow(repository: SqliteOrderRepository): void {
  open.splice(open.indexOf(repository), 1);
  repository.close();
}

beforeEach(() => {
  folder = mkdtempSync(join(tmpdir(), "checkout-sqlite-"));
  path = join(folder, "orders.db");
});

afterEach(() => {
  for (const repository of open.splice(0)) {
    repository.close();
  }
  rmSync(folder, { recursive: true, force: true });
});

describe("SqliteOrderRepository", () => {
  it("saves and finds an order", () => {
    const repository = openRepository();
    const receipt = checkout(
      { items: [{ sku: "COFFEE", quantity: 3 }], coupon: "SAVE10" },
      tinyCatalog,
    );

    const id = repository.save(receipt);

    expect(id).toBe("ORDER-1");
    expect(repository.find(id)).toEqual({
      id: "ORDER-1",
      created_at: "2026-01-15T10:30:00Z",
      receipt,
    });
  });

  it("gives back an order that follows order.schema.json", () => {
    const repository = openRepository();
    const id = repository.save(oneAppleReceipt);

    expect(schemaErrors("order.schema.json", repository.find(id))).toEqual([]);
  });

  it("numbers orders sequentially", () => {
    const repository = openRepository();

    const ids = [
      repository.save(oneAppleReceipt),
      repository.save(oneAppleReceipt),
      repository.save(oneAppleReceipt),
    ];

    expect(ids).toEqual(["ORDER-1", "ORDER-2", "ORDER-3"]);
  });

  it("continues numbering after the database is closed and opened again", () => {
    const first = openRepository();
    first.save(oneAppleReceipt);
    first.save(oneAppleReceipt);
    closeNow(first);

    const second = openRepository();

    expect(second.save(oneAppleReceipt)).toBe("ORDER-3");
    expect(second.find("ORDER-2")?.id).toBe("ORDER-2");
  });

  it("really writes to the database file (checked with a separate connection)", () => {
    const repository = openRepository();
    repository.save(oneAppleReceipt);
    closeNow(repository);

    const db = new DatabaseSync(path);
    const row = db.prepare("SELECT COUNT(*) AS n FROM orders").get() as { n: number };
    db.close();

    expect(row.n).toBe(1);
  });

  it("refuses to save or find once it is closed, instead of silently doing nothing", () => {
    const repository = openRepository();
    closeNow(repository);

    expect(() => repository.save(oneAppleReceipt)).toThrow("database is not open");
    expect(() => repository.find("ORDER-1")).toThrow("database is not open");
  });

  it("works in memory too (':memory:' is SQLite's name for a throw-away database)", () => {
    const repository = new SqliteOrderRepository(":memory:", fixedClock);
    open.push(repository);

    expect(repository.save(oneAppleReceipt)).toBe("ORDER-1");
    expect(repository.find("ORDER-1")?.receipt).toEqual(oneAppleReceipt);
  });

  it.each(["ORDER-2", "ORDER-0", "ORDER-01", "ORDER-x", "1", "", "ORDER-1; DROP TABLE orders"])(
    "answers 'not found' for %j without throwing",
    (id) => {
      const repository = openRepository();
      repository.save(oneAppleReceipt);

      expect(repository.find(id)).toBeUndefined();
      expect(repository.find("ORDER-1")).toBeDefined(); // and the table is still there
    },
  );

  it("stamps each order with the clock's time at the moment of saving", () => {
    let now = new Date("2026-01-15T10:30:00Z");
    const repository = openRepository(() => now);

    const first = repository.save(oneAppleReceipt);
    now = new Date("2026-01-16T09:00:00Z");
    const second = repository.save(oneAppleReceipt);

    expect(repository.find(first)?.created_at).toBe("2026-01-15T10:30:00Z");
    expect(repository.find(second)?.created_at).toBe("2026-01-16T09:00:00Z");
  });
});
