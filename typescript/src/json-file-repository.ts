// An order repository that writes one JSON file per order into a folder.

import type { Clock, Order, OrderRepository } from "./orders.ts";
import { systemClock } from "./orders.ts";
import type { Receipt } from "./types.ts";

/**
 * Keeps each order as `<folder>/ORDER-<n>.json`, in the shape of
 * spec/schema/order.schema.json. Ids continue where the files in the folder left off,
 * so reopening the same folder later keeps numbering correctly.
 */
export class JsonFileOrderRepository implements OrderRepository {
  constructor(folder: string, clock: Clock = systemClock) {
    void folder;
    void clock;
  }

  save(receipt: Receipt): string {
    throw new Error(`save(${receipt.total}) is not implemented yet`);
  }

  find(id: string): Order | undefined {
    throw new Error(`find(${id}) is not implemented yet`);
  }
}
