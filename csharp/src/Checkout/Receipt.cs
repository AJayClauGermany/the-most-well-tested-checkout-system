namespace CheckoutSystem;

/// <summary>
/// The checkout output (rules §5). All amounts are whole cents.
/// </summary>
/// <param name="Lines">One line per product, in order of first appearance in the basket.</param>
/// <param name="Subtotal">Sum of all line totals.</param>
/// <param name="PromotionDiscount">Sum of all line discounts.</param>
/// <param name="CouponDiscount">Discount from the coupon, or 0.</param>
/// <param name="Shipping">Shipping cost (0 when shipping is free).</param>
/// <param name="Tax">8% of the discounted subtotal.</param>
/// <param name="Total">What the customer pays.</param>
public sealed record Receipt(
    IReadOnlyList<ReceiptLine> Lines,
    int Subtotal,
    int PromotionDiscount,
    int CouponDiscount,
    int Shipping,
    int Tax,
    int Total);

/// <summary>
/// One priced line of a receipt, after duplicate basket lines were merged.
/// </summary>
/// <param name="LineTotal">UnitPrice × Quantity, before the promotion.</param>
/// <param name="LineDiscount">How much the product's promotion takes off this line.</param>
public sealed record ReceiptLine(
    string Sku,
    string Name,
    int Quantity,
    int UnitPrice,
    int LineTotal,
    int LineDiscount);
