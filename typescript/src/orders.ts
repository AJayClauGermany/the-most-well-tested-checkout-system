// Saving orders (rules §6): the order type, the clock, the repository interface and
// the simplest repository there is - a list in memory.

import type { Basket, Catalog, Receipt } from "./types.ts";

/** A saved order, exactly as spec/schema/order.schema.json describes it. */
export interface Order {
  id: string;
  /** UTC, ISO 8601, whole seconds, trailing Z: "2026-01-15T10:30:00Z". */
  created_at: string;
  receipt: Receipt;
}

/** Anything that answers "what time is it?". Tests pass one that always says the same time. */
export type Clock = () => Date;

/** The real clock. */
export const systemClock: Clock = () => new Date();

/** Formats a date the way `created_at` must look: whole seconds, UTC, trailing Z. */
export function formatTimestamp(date: Date): string {
  throw new Error(`formatTimestamp(${date.toISOString()}) is not implemented yet`);
}

/** The id of the n-th order saved into a repository: "ORDER-1", "ORDER-2", ... */
export function orderId(n: number): string {
  throw new Error(`orderId(${n}) is not implemented yet`);
}

/** The inverse of orderId: "ORDER-7" gives 7. Anything that is not an order id gives `undefined`. */
export function orderNumber(id: string): number | undefined {
  throw new Error(`orderNumber(${id}) is not implemented yet`);
}

/** Where orders are kept (rules §6). Implemented in memory, as JSON files and in SQLite. */
export interface OrderRepository {
  /** Saves a receipt and returns the new order's id. */
  save(receipt: Receipt): string;
  /** Looks an order up. `undefined` means "not found" - a normal answer, not an error. */
  find(id: string): Order | undefined;
}

/** Keeps orders in a plain array. Doubles as the hand-written fake for tests. */
export class InMemoryOrderRepository implements OrderRepository {
  constructor(clock: Clock = systemClock) {
    void clock;
  }

  save(receipt: Receipt): string {
    throw new Error(`save(${receipt.total}) is not implemented yet`);
  }

  find(id: string): Order | undefined {
    throw new Error(`find(${id}) is not implemented yet`);
  }

  /** How many orders are saved. Handy for tests ("the order repository is empty"). */
  count(): number {
    throw new Error("count() is not implemented yet");
  }
}

/**
 * Checks out, then saves the receipt as an order. A failed checkout saves nothing:
 * the CheckoutError is thrown before the repository is touched.
 */
export function placeOrder(
  basket: Basket,
  catalog: Catalog,
  repository: OrderRepository,
): { receipt: Receipt; orderId: string } {
  throw new Error(
    `placeOrder(${basket.items.length}, ${catalog.products.length}, ${typeof repository}) is not implemented yet`,
  );
}
