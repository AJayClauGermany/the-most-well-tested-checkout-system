// Tax (rules §4.6).

export const TAX_PERCENT = 8;

/** 8% of `discountedSubtotal`, rounded half up. Shipping is never passed in here. */
export function taxOn(discountedSubtotal: number): number {
  throw new Error(`taxOn(${discountedSubtotal}) is not implemented yet`);
}
