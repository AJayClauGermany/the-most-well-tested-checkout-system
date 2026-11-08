// Coupons (rules §4.4).

/**
 * The coupon code as the shop knows it: trimmed and upper case.
 * Missing, `null`, empty or blank all mean "no coupon" and give `undefined`.
 */
export function normalizeCouponCode(coupon: string | null | undefined): string | undefined {
  throw new Error(`normalizeCouponCode(${JSON.stringify(coupon)}) is not implemented yet`);
}

/**
 * The coupon discount in cents for `afterPromotions` cents of shopping.
 *
 * @throws {CheckoutError} UNKNOWN_COUPON or COUPON_MINIMUM_NOT_MET (rules §4.1 rows 6 and 7).
 */
export function couponDiscount(coupon: string | null | undefined, afterPromotions: number): number {
  throw new Error(
    `couponDiscount(${JSON.stringify(coupon)}, ${afterPromotions}) is not implemented yet`,
  );
}
