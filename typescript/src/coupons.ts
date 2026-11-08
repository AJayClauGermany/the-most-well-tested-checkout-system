// Coupons (rules §4.4).

import { CheckoutError } from "./errors.ts";
import { roundHalfUp } from "./money.ts";

/** What the shop knows about one coupon code. */
interface Coupon {
  /** The discount in cents, given the price after promotions. */
  discount: (afterPromotions: number) => number;
  /** The least the customer must spend (after promotions) to use it. */
  minimumSpend: number;
}

/** The coupons that exist, by their normalized code. */
const COUPONS: ReadonlyMap<string, Coupon> = new Map<string, Coupon>([
  ["SAVE10", { discount: (afterPromotions) => roundHalfUp(afterPromotions, 10), minimumSpend: 0 }],
  ["WELCOME5", { discount: () => 500, minimumSpend: 2500 }],
]);

/**
 * The coupon code as the shop knows it: trimmed and upper case.
 * Missing, `null`, empty or blank all mean "no coupon" and give `undefined`.
 */
export function normalizeCouponCode(coupon: string | null | undefined): string | undefined {
  const code = (coupon ?? "").trim().toUpperCase();
  return code === "" ? undefined : code;
}

/**
 * The coupon discount in cents for `afterPromotions` cents of shopping.
 *
 * @throws {CheckoutError} UNKNOWN_COUPON or COUPON_MINIMUM_NOT_MET (rules §4.1 rows 6 and 7).
 */
export function couponDiscount(coupon: string | null | undefined, afterPromotions: number): number {
  const code = normalizeCouponCode(coupon);
  if (code === undefined) {
    return 0;
  }

  const known = COUPONS.get(code);
  if (known === undefined) {
    throw new CheckoutError("UNKNOWN_COUPON", `Unknown coupon: ${code}`);
  }
  if (afterPromotions < known.minimumSpend) {
    throw new CheckoutError(
      "COUPON_MINIMUM_NOT_MET",
      `${code} needs a spend of at least ${known.minimumSpend} cents, got ${afterPromotions}`,
    );
  }
  return known.discount(afterPromotions);
}
