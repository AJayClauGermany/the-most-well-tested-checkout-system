// The public API of the package. Import from here, e.g.
//   import { checkout, CheckoutError } from "../src/index.ts";
export { checkout } from "./checkout.ts";
export { findProduct, requireProduct } from "./catalog.ts";
export { CheckoutError, ERROR_CODES } from "./errors.ts";
export type { ErrorCode } from "./errors.ts";
export { roundHalfUp } from "./money.ts";
export { isValidQuantity, MAX_QUANTITY, MIN_QUANTITY, validateItems } from "./validation.ts";
export { mergeItems, priceLine } from "./lines.ts";
export { promotionDiscount } from "./promotions.ts";
export { couponDiscount, normalizeCouponCode } from "./coupons.ts";
export { parseShippingMethod, shippingCost } from "./shipping.ts";
export type { ShippingMethod } from "./shipping.ts";
export { TAX_PERCENT, taxOn } from "./tax.ts";
export {
  formatTimestamp,
  InMemoryOrderRepository,
  orderId,
  orderNumber,
  placeOrder,
  systemClock,
} from "./orders.ts";
export type { Clock, Order, OrderRepository } from "./orders.ts";
export { JsonFileOrderRepository } from "./json-file-repository.ts";
export { SqliteOrderRepository } from "./sqlite-repository.ts";
export { loadStandardCatalog, runCheckoutCommand, STANDARD_CATALOG_PATH } from "./command.ts";
export type { CommandResult } from "./command.ts";
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
