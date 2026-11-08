// The public API of the package. Import from here, e.g.
//   import { checkout, CheckoutError } from "../src/index.ts";
export { checkout } from "./checkout.ts";
export { findProduct } from "./catalog.ts";
export { CheckoutError, ERROR_CODES } from "./errors.ts";
export type { ErrorCode } from "./errors.ts";
export type {
  Basket,
  BasketItem,
  BulkDiscountPromotion,
  BuyXGetYFreePromotion,
  Catalog,
  Product,
  Promotion,
  Receipt,
  ReceiptLine,
} from "./types.ts";
