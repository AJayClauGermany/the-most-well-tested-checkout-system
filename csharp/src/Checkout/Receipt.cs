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
    int Total)
{
    /// <summary>
    /// Two receipts are equal when all six totals AND all lines are equal.
    /// </summary>
    /// <remarks>
    /// A record compares its fields by value – but a list field is compared by reference
    /// ("is it the very same list object?"), which is never what we want for a receipt.
    /// So this record spells out its own equality and compares the lines one by one.
    /// </remarks>
    public bool Equals(Receipt? other) =>
        other is not null
        && Totals == other.Totals
        && Lines.SequenceEqual(other.Lines);

    /// <summary>Equal receipts have equal totals, so the totals alone make a valid hash code.</summary>
    public override int GetHashCode() => Totals.GetHashCode();

    /// <summary>The six amounts as one value, so they can be compared in one go.</summary>
    private (int, int, int, int, int, int) Totals => (Subtotal, PromotionDiscount, CouponDiscount, Shipping, Tax, Total);
}

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
