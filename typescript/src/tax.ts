// Tax (rules §4.6).

import { roundHalfUp } from "./money.ts";

export const TAX_PERCENT = 8;

/** 8% of `discountedSubtotal`, rounded half up. Shipping is never passed in here. */
export function taxOn(discountedSubtotal: number): number {
  return roundHalfUp(discountedSubtotal, TAX_PERCENT);
}
