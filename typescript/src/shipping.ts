// Shipping (rules §4.5).

export type ShippingMethod = "standard" | "express";

/**
 * Turns the basket's `shipping` field into a known method. Missing means "standard".
 *
 * @throws {CheckoutError} UNKNOWN_SHIPPING_METHOD for anything else (rules §4.1 row 5).
 */
export function parseShippingMethod(shipping: string | undefined): ShippingMethod {
  throw new Error(`parseShippingMethod(${JSON.stringify(shipping)}) is not implemented yet`);
}

/** The shipping cost in cents for this method and this discounted subtotal. */
export function shippingCost(method: ShippingMethod, discountedSubtotal: number): number {
  throw new Error(`shippingCost(${method}, ${discountedSubtotal}) is not implemented yet`);
}
