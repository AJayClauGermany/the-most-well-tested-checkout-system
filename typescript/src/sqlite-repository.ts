// An order repository backed by SQLite, using Node's built-in node:sqlite module.

import { DatabaseSync } from "node:sqlite";

import type { Clock, Order, OrderRepository } from "./orders.ts";
import { formatTimestamp, orderId, orderNumber, systemClock } from "./orders.ts";
import type { Receipt } from "./types.ts";

/** One row of the `orders` table. The receipt is stored as JSON text. */
interface OrderRow {
  number: number;
  created_at: string;
  receipt: string;
}

/**
 * Keeps orders in one table of an SQLite database file (or ":memory:").
 * The database row number becomes the order number, so ids are sequential and survive
 * closing and reopening the file.
 */
export class SqliteOrderRepository implements OrderRepository {
  private readonly db: DatabaseSync;
  private readonly clock: Clock;

  constructor(path: string, clock: Clock = systemClock) {
    this.db = new DatabaseSync(path);
    this.clock = clock;
    this.db.exec(`
      CREATE TABLE IF NOT EXISTS orders (
        number     INTEGER PRIMARY KEY AUTOINCREMENT,
        created_at TEXT NOT NULL,
        receipt    TEXT NOT NULL
      )
    `);
  }

  save(receipt: Receipt): string {
    // "?" placeholders keep the values apart from the SQL, so a receipt can never be mistaken
    // for a command (no SQL injection).
    const result = this.db
      .prepare("INSERT INTO orders (created_at, receipt) VALUES (?, ?)")
      .run(formatTimestamp(this.clock()), JSON.stringify(receipt));
    return orderId(Number(result.lastInsertRowid));
  }

  find(id: string): Order | undefined {
    const number = orderNumber(id);
    if (number === undefined) {
      return undefined;
    }
    const row = this.db
      .prepare("SELECT number, created_at, receipt FROM orders WHERE number = ?")
      .get(number) as OrderRow | undefined;
    if (row === undefined) {
      return undefined;
    }
    return {
      id: orderId(row.number),
      created_at: row.created_at,
      receipt: JSON.parse(row.receipt) as Receipt,
    };
  }

  /** Closes the database file. Call it when you are done (tests do, in afterEach). */
  close(): void {
    this.db.close();
  }
}
