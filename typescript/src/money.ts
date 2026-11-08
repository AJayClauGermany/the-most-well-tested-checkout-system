// Money helpers (rules §1 and §7). All amounts are whole cents.

/**
 * `percent`% of `amount`, rounded half up to a whole cent (rules §7).
 *
 * Done with whole numbers only, as the hint in rules §7 suggests: `amount × percent` is the
 * exact answer in hundredths of a cent; adding 50 hundredths and dividing by 100 (rounding
 * down) is the same as rounding half up. No floating point, so no `0.1 + 0.2` surprises.
 *
 *   roundHalfUp(225, 10)  →  (2250 + 50) / 100  →  23      (22.5 rounds UP)
 *   roundHalfUp(224, 10)  →  (2240 + 50) / 100  →  22.9 → 22
 */
export function roundHalfUp(amount: number, percent: number): number {
  return Math.floor((amount * percent + 50) / 100);
}
