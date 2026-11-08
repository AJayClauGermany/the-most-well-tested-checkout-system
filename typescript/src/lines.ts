import { promotionDiscount } from "./promotions.ts";
import type { BasketItem, Product, ReceiptLine } from "./types.ts";

/** One item per distinct SKU, in order of first appearance, quantities added up (rules §4.2). */
export function mergeItems(items: BasketItem[]): BasketItem[] {
  // A Map remembers the order keys were first added - exactly the order the rules want.
  const quantities = new Map<string, number>();
  for (const item of items) {
    quantities.set(item.sku, (quantities.get(item.sku) ?? 0) + item.quantity);
  }
  return [...quantities].map(([sku, quantity]) => ({ sku, quantity }));
}

/** The receipt line for `quantity` of `product`, including its promotion discount. */
export function priceLine(product: Product, quantity: number): ReceiptLine {
  const lineTotal = product.unit_price * quantity;
  return {
    sku: product.sku,
    name: product.name,
    quantity,
    unit_price: product.unit_price,
    line_total: lineTotal,
    line_discount: promotionDiscount(product.promotion, quantity, product.unit_price),
  };
}
