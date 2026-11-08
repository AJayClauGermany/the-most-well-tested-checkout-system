import { roundHalfUp } from "./money.ts";
import type { Promotion } from "./types.ts";

/** The discount a promotion gives one receipt line, in cents (rules §4.2). `0` without a promotion. */
export function promotionDiscount(
  promotion: Promotion | null,
  quantity: number,
  unitPrice: number,
): number {
  if (promotion === null) {
    return 0;
  }

  switch (promotion.type) {
    case "buy_x_get_y_free": {
      // For every complete group of `buy + free` items, `free` of them cost nothing.
      const groupSize = promotion.buy + promotion.free;
      const completeGroups = Math.floor(quantity / groupSize);
      const freeItems = completeGroups * promotion.free;
      return freeItems * unitPrice;
    }
    case "bulk_discount": {
      // `percent`% off the whole line, but only from `min_quantity` items.
      if (quantity < promotion.min_quantity) {
        return 0;
      }
      return roundHalfUp(unitPrice * quantity, promotion.percent);
    }
  }
}
