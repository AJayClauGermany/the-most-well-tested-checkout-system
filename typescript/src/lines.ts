import type { BasketItem, Product, ReceiptLine } from "./types.ts";

/** One item per distinct SKU, in order of first appearance, quantities added up (rules §4.2). */
export function mergeItems(items: BasketItem[]): BasketItem[] {
  throw new Error(`mergeItems(${JSON.stringify(items)}) is not implemented yet`);
}

/** The receipt line for `quantity` of `product`, including its promotion discount. */
export function priceLine(product: Product, quantity: number): ReceiptLine {
  throw new Error(`priceLine(${product.sku}, ${quantity}) is not implemented yet`);
}
