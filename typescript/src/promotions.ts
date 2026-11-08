import type { Promotion } from "./types.ts";

/** The discount a promotion gives one receipt line, in cents (rules §4.2). `0` without a promotion. */
export function promotionDiscount(
  promotion: Promotion | null,
  quantity: number,
  unitPrice: number,
): number {
  throw new Error(
    `promotionDiscount(${JSON.stringify(promotion)}, ${quantity}, ${unitPrice}) is not implemented yet`,
  );
}
