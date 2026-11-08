// Saving orders (rules §6): the order type, the clock, the repository interface and
// the simplest repository there is - a list in memory.

import { checkout } from "./checkout.ts";
import type { Basket, Catalog, Receipt } from "./types.ts";

/** A saved order, exactly as spec/schema/order.schema.json describes it. */
export interface Order {
  id: string;
  /** UTC, ISO 8601, whole seconds, trailing Z: "2026-01-15T10:30:00Z". */
  created_at: string;
  receipt: Receipt;
}

/**
 * Anything that answers "what time is it?". Tests pass one that always says the same time -
 * the rules (§6) demand that the time source can be replaced, because `created_at` is the one
 * value that changes from run to run.
 */
export type Clock = () => Date;

/** The real clock. */
export const systemClock: Clock = () => new Date();

/** Formats a date the way `created_at` must look: whole seconds, UTC, trailing Z. */
export function formatTimestamp(date: Date): string {
  // toISOString() gives "2026-01-15T10:30:00.999Z"; cut off the milliseconds.
  return `${date.toISOString().slice(0, 19)}Z`;
}

const ORDER_ID_PREFIX = "ORDER-";

/** The id of the n-th order saved into a repository: "ORDER-1", "ORDER-2", ... */
export function orderId(n: number): string {
  return `${ORDER_ID_PREFIX}${n}`;
}

/** The inverse of orderId: "ORDER-7" gives 7. Anything that is not an order id gives `undefined`. */
export function orderNumber(id: string): number | undefined {
  // A positive whole number with no leading zero, exactly as order.schema.json spells it.
  const match = /^ORDER-([1-9][0-9]*)$/.exec(id);
  return match === null ? undefined : Number(match[1]);
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
  private readonly orders: Order[] = [];
  private readonly clock: Clock;

  constructor(clock: Clock = systemClock) {
    this.clock = clock;
  }

  save(receipt: Receipt): string {
    const id = orderId(this.orders.length + 1);
    this.orders.push({
      id,
      created_at: formatTimestamp(this.clock()),
      // A copy, so that whoever gave us the receipt cannot change the saved order afterwards.
      receipt: structuredClone(receipt),
    });
    return id;
  }

  find(id: string): Order | undefined {
    return this.orders.find((order) => order.id === id);
  }

  /** How many orders are saved. Handy for tests ("the order repository is empty"). */
  count(): number {
    return this.orders.length;
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
  const receipt = checkout(basket, catalog);
  return { receipt, orderId: repository.save(receipt) };
}
