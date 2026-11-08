// Shipping (rules §4.5).

import { CheckoutError } from "./errors.ts";

export type ShippingMethod = "standard" | "express";

const STANDARD_COST = 499;
const EXPRESS_COST = 999;
/** Standard shipping is free from this discounted subtotal upwards. */
const FREE_SHIPPING_FROM = 5000;

/**
 * Turns the basket's `shipping` field into a known method. Missing means "standard".
 *
 * @throws {CheckoutError} UNKNOWN_SHIPPING_METHOD for anything else (rules §4.1 row 5).
 */
export function parseShippingMethod(shipping: string | undefined): ShippingMethod {
  if (shipping === undefined) {
    return "standard";
  }
  if (shipping === "standard" || shipping === "express") {
    return shipping;
  }
  throw new CheckoutError("UNKNOWN_SHIPPING_METHOD", `Unknown shipping method: ${shipping}`);
}

/** The shipping cost in cents for this method and this discounted subtotal. */
export function shippingCost(method: ShippingMethod, discountedSubtotal: number): number {
  if (method === "express") {
    return EXPRESS_COST;
  }
  return discountedSubtotal >= FREE_SHIPPING_FROM ? 0 : STANDARD_COST;
}
