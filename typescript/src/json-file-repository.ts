// An order repository that writes one JSON file per order into a folder.

import { existsSync, mkdirSync, readdirSync, readFileSync, writeFileSync } from "node:fs";
import { basename, join } from "node:path";

import type { Clock, Order, OrderRepository } from "./orders.ts";
import { formatTimestamp, orderId, orderNumber, systemClock } from "./orders.ts";
import type { Receipt } from "./types.ts";

/**
 * Keeps each order as `<folder>/ORDER-<n>.json`, in the shape of
 * spec/schema/order.schema.json. Ids continue where the files in the folder left off,
 * so reopening the same folder later keeps numbering correctly.
 */
export class JsonFileOrderRepository implements OrderRepository {
  private readonly folder: string;
  private readonly clock: Clock;

  constructor(folder: string, clock: Clock = systemClock) {
    this.folder = folder;
    this.clock = clock;
    mkdirSync(folder, { recursive: true });
  }

  save(receipt: Receipt): string {
    const order: Order = {
      id: orderId(this.highestOrderNumber() + 1),
      created_at: formatTimestamp(this.clock()),
      receipt,
    };
    writeFileSync(this.pathOf(order.id), `${JSON.stringify(order, null, 2)}\n`);
    return order.id;
  }

  find(id: string): Order | undefined {
    // Only accept real order ids: "../secret" must not read files outside the folder.
    if (orderNumber(id) === undefined) {
      return undefined;
    }
    const path = this.pathOf(id);
    if (!existsSync(path)) {
      return undefined;
    }
    return JSON.parse(readFileSync(path, "utf8")) as Order;
  }

  private pathOf(id: string): string {
    return join(this.folder, `${id}.json`);
  }

  /** The highest order number among the `ORDER-<n>.json` files in the folder, or 0 when there are none. */
  private highestOrderNumber(): number {
    const numbers = readdirSync(this.folder)
      .filter((file) => file.endsWith(".json"))
      .map((file) => orderNumber(basename(file, ".json")))
      .filter((number) => number !== undefined);
    return Math.max(0, ...numbers);
  }
}
