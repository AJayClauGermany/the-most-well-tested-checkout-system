// Test helper: small, hand-made test data for unit tests.
//
// Unit tests prefer tiny made-up data over the real catalog: when a test fails you can
// work out the right answer in your head. The shared spec cases use the real catalog.

import type { Catalog, Product, Receipt } from "../../src/index.ts";

export const APPLE: Product = { sku: "APPLE", name: "Apple", unit_price: 50, promotion: null };
export const PEN: Product = { sku: "PEN", name: "Pen", unit_price: 120, promotion: null };
export const COFFEE: Product = {
  sku: "COFFEE",
  name: "Bag of Coffee Beans",
  unit_price: 899,
  promotion: { type: "buy_x_get_y_free", buy: 2, free: 1 },
};
export const NOTEBOOK: Product = {
  sku: "NOTEBOOK",
  name: "Notebook",
  unit_price: 349,
  promotion: { type: "bulk_discount", min_quantity: 10, percent: 10 },
};
/** A $100.00 item: handy for reaching free shipping and coupon minimums quickly. */
export const TV: Product = { sku: "TV", name: "Television", unit_price: 10000, promotion: null };

/** Five products, two of them with a promotion. */
export const tinyCatalog: Catalog = { products: [APPLE, PEN, COFFEE, NOTEBOOK, TV] };

/** A valid receipt for one apple, used wherever "some receipt" is needed. */
export const oneAppleReceipt: Receipt = {
  lines: [
    { sku: "APPLE", name: "Apple", quantity: 1, unit_price: 50, line_total: 50, line_discount: 0 },
  ],
  subtotal: 50,
  promotion_discount: 0,
  coupon_discount: 0,
  shipping: 499,
  tax: 4,
  total: 553,
};
