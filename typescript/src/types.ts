// The data types of the checkout, mirroring the JSON files in spec/.
//
// Why are the field names in snake_case (`unit_price`) instead of the usual
// TypeScript camelCase (`unitPrice`)? Because the JSON in spec/ uses snake_case.
// Keeping the same names means we can use parsed JSON directly, with no code to
// translate field names back and forth (and no bugs in that code!).
//
// All money is a whole number of cents: 199 means $1.99 (rules §1).

// ---------------------------------------------------------------------------
// The catalog (rules §2)
// ---------------------------------------------------------------------------

/** "Buy X, get Y free": for every `buy + free` items, `free` of them cost nothing. */
export interface BuyXGetYFreePromotion {
  type: "buy_x_get_y_free";
  buy: number;
  free: number;
}

/** "Bulk discount": `percent`% off the line when you buy at least `min_quantity`. */
export interface BulkDiscountPromotion {
  type: "bulk_discount";
  min_quantity: number;
  percent: number;
}

/**
 * A promotion is one of the two kinds above.
 *
 * This is a "discriminated union": the `type` field tells TypeScript which kind
 * it is. Inside `if (promotion.type === "bulk_discount") { ... }` TypeScript knows
 * that `promotion.min_quantity` exists.
 */
export type Promotion = BuyXGetYFreePromotion | BulkDiscountPromotion;

export interface Product {
  sku: string;
  name: string;
  /** Price of one item, in cents. */
  unit_price: number;
  /** `null` means "no promotion" (that is how the JSON spells it). */
  promotion: Promotion | null;
}

export interface Catalog {
  products: Product[];
}

// ---------------------------------------------------------------------------
// The basket: the input (rules §3)
// ---------------------------------------------------------------------------

export interface BasketItem {
  sku: string;
  /** Should be 1 to 99, but invalid values are allowed here so we can test validation. */
  quantity: number;
}

export interface Basket {
  items: BasketItem[];
  /** A coupon code. Missing, `null` or blank all mean "no coupon". */
  coupon?: string | null;
  /**
   * "standard" or "express"; missing means "standard".
   * It is a plain `string` (not `"standard" | "express"`) on purpose: a basket can
   * contain a wrong value such as "overnight", and checkout must report it.
   */
  shipping?: string;
}

// ---------------------------------------------------------------------------
// The receipt: the output (rules §5)
// ---------------------------------------------------------------------------

export interface ReceiptLine {
  sku: string;
  name: string;
  quantity: number;
  unit_price: number;
  line_total: number;
  line_discount: number;
}

export interface Receipt {
  lines: ReceiptLine[];
  subtotal: number;
  promotion_discount: number;
  coupon_discount: number;
  shipping: number;
  tax: number;
  total: number;
}
