namespace CheckoutSystem;

/// <summary>The two kinds of promotion (rules §4.2).</summary>
public static class Promotions
{
    /// <summary>
    /// The discount a promotion gives one receipt line, in cents. <c>0</c> without a promotion.
    /// </summary>
    /// <remarks>
    /// A <c>switch</c> expression on the type: each kind of promotion gets its own formula,
    /// and anything else (no promotion at all) gives no discount.
    /// </remarks>
    public static int Discount(Promotion? promotion, int quantity, int unitPrice) => promotion switch
    {
        // For every complete group of `buy + free` items, `free` of them cost nothing.
        BuyXGetYFree p => quantity / (p.Buy + p.Free) * p.Free * unitPrice,

        // `percent`% off the whole line, but only from `min_quantity` items.
        BulkDiscount p when quantity >= p.MinQuantity => Money.RoundHalfUp(unitPrice * quantity, p.Percent),

        _ => 0,
    };
}
