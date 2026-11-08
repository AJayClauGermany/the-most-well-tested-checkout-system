import type { Basket, Catalog } from "./types.ts";

export const MIN_QUANTITY = 1;
export const MAX_QUANTITY = 99;

/** True when `quantity` is a whole number from 1 to 99 (rules §4.1 row 3). */
export function isValidQuantity(quantity: number): boolean {
  throw new Error(`isValidQuantity(${quantity}) is not implemented yet`);
}

/**
 * Throws a CheckoutError for rows 1-4 of rules §4.1, in that order.
 * Does nothing when the items are fine.
 */
export function validateItems(basket: Basket, catalog: Catalog): void {
  throw new Error(
    `validateItems(${JSON.stringify(basket)}, ${catalog.products.length}) is not implemented yet`,
  );
}
