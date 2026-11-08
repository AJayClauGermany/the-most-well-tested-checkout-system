import { requireProduct } from "./catalog.ts";
import { couponDiscount } from "./coupons.ts";
import { mergeItems, priceLine } from "./lines.ts";
import { parseShippingMethod, shippingCost } from "./shipping.ts";
import { taxOn } from "./tax.ts";
import type { Basket, Catalog, Receipt } from "./types.ts";
import { validateItems } from "./validation.ts";

/** Adds up a list of amounts. */
function sum(amounts: number[]): number {
  return amounts.reduce((total, amount) => total + amount, 0);
}

/**
 * Calculates the receipt for a basket, following spec/checkout-rules.md.
 *
 * It reads like the hint in rules §4.7: validate, merge lines, price, coupon, shipping, tax,
 * total. Each step is a small function with its own unit tests.
 *
 * @throws {CheckoutError} with one of the six error codes when the basket is invalid.
 */
export function checkout(basket: Basket, catalog: Catalog): Receipt {
  validateItems(basket, catalog); // rules §4.1 rows 1-4
  const method = parseShippingMethod(basket.shipping); // row 5

  const lines = mergeItems(basket.items).map((item) =>
    priceLine(requireProduct(catalog, item.sku), item.quantity),
  );
  const subtotal = sum(lines.map((line) => line.line_total));
  const promotionDiscount = sum(lines.map((line) => line.line_discount));
  const afterPromotions = subtotal - promotionDiscount;

  const coupon = couponDiscount(basket.coupon, afterPromotions); // rows 6-7
  const discountedSubtotal = afterPromotions - coupon;

  const shipping = shippingCost(method, discountedSubtotal);
  const tax = taxOn(discountedSubtotal);

  return {
    lines,
    subtotal,
    promotion_discount: promotionDiscount,
    coupon_discount: coupon,
    shipping,
    tax,
    total: discountedSubtotal + shipping + tax,
  };
}
