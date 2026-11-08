// The six error codes from rules §9, and the error that checkout throws.

/** Every reason a checkout can fail. The spelling must match the spec exactly. */
export type ErrorCode =
  | "EMPTY_BASKET"
  | "UNKNOWN_PRODUCT"
  | "INVALID_QUANTITY"
  | "UNKNOWN_SHIPPING_METHOD"
  | "UNKNOWN_COUPON"
  | "COUPON_MINIMUM_NOT_MET";

/** All error codes, in the order rules §9 lists them. Handy for loops and tests. */
export const ERROR_CODES: readonly ErrorCode[] = [
  "EMPTY_BASKET",
  "UNKNOWN_PRODUCT",
  "INVALID_QUANTITY",
  "UNKNOWN_SHIPPING_METHOD",
  "UNKNOWN_COUPON",
  "COUPON_MINIMUM_NOT_MET",
];

/**
 * Thrown by `checkout()` when the basket cannot be checked out.
 *
 * In TypeScript the natural way to report a failure is to throw an error.
 * Callers can tell the reasons apart by looking at `error.code`.
 */
export class CheckoutError extends Error {
  readonly code: ErrorCode;

  constructor(code: ErrorCode, message: string = code) {
    super(message);
    this.name = "CheckoutError";
    this.code = code;
  }
}
