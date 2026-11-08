// Basket validation (rules §4.1, rows 1 to 4). Rows 5 to 7 (shipping and coupon) are checked
// by shipping.ts and coupons.ts, because the coupon check needs the price after promotions.

import { requireProduct } from "./catalog.ts";
import { CheckoutError } from "./errors.ts";
import { mergeItems } from "./lines.ts";
import type { Basket, Catalog } from "./types.ts";

export const MIN_QUANTITY = 1;
export const MAX_QUANTITY = 99;

/** True when `quantity` is a whole number from 1 to 99 (rules §4.1 row 3). */
export function isValidQuantity(quantity: number): boolean {
  return Number.isInteger(quantity) && quantity >= MIN_QUANTITY && quantity <= MAX_QUANTITY;
}

/**
 * Throws a CheckoutError for rows 1-4 of rules §4.1, in that order.
 * Does nothing when the items are fine.
 */
export function validateItems(basket: Basket, catalog: Catalog): void {
  // Row 1: no items at all.
  if (basket.items.length === 0) {
    throw new CheckoutError("EMPTY_BASKET", "The basket has no items");
  }

  // Row 2: every SKU must exist. All SKUs are checked before any quantity, so an unknown
  // product on line 3 still wins over a bad quantity on line 1.
  for (const item of basket.items) {
    requireProduct(catalog, item.sku);
  }

  // Row 3: every quantity must be a whole number from 1 to 99.
  for (const item of basket.items) {
    if (!isValidQuantity(item.quantity)) {
      throw new CheckoutError(
        "INVALID_QUANTITY",
        `Invalid quantity for ${item.sku}: ${item.quantity}`,
      );
    }
  }

  // Row 4: the lines of one SKU, added together, may not go over 99 either.
  for (const line of mergeItems(basket.items)) {
    if (line.quantity > MAX_QUANTITY) {
      throw new CheckoutError(
        "INVALID_QUANTITY",
        `Too many ${line.sku} in total: ${line.quantity}`,
      );
    }
  }
}
