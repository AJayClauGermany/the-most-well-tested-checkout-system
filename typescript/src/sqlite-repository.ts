// An order repository backed by SQLite, using Node's built-in node:sqlite module.

import type { Clock, Order, OrderRepository } from "./orders.ts";
import { systemClock } from "./orders.ts";
import type { Receipt } from "./types.ts";

/**
 * Keeps orders in one table of an SQLite database file (or ":memory:").
 * The database row number becomes the order number, so ids are sequential and survive
 * closing and reopening the file.
 */
export class SqliteOrderRepository implements OrderRepository {
  constructor(path: string, clock: Clock = systemClock) {
    void path;
    void clock;
  }

  save(receipt: Receipt): string {
    throw new Error(`save(${receipt.total}) is not implemented yet`);
  }

  find(id: string): Order | undefined {
    throw new Error(`find(${id}) is not implemented yet`);
  }

  /** Closes the database file. Call it when you are done (tests do, in afterEach). */
  close(): void {
    throw new Error("close() is not implemented yet");
  }
}
