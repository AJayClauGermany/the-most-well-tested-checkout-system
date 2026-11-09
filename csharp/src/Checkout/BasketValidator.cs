namespace CheckoutSystem;

/// <summary>
/// Basket validation: rows 1 to 4 of rules §4.1. Rows 5 to 7 (shipping and coupon) are
/// checked by <see cref="Shipping"/> and <see cref="Coupons"/>, because the coupon check
/// needs the price after promotions.
/// </summary>
public static class BasketValidator
{
    public const int MinQuantity = 1;
    public const int MaxQuantity = 99;

    /// <summary>True when <paramref name="quantity"/> is from 1 to 99 (rules §4.1 row 3).</summary>
    public static bool IsValidQuantity(int quantity) =>
        quantity is >= MinQuantity and <= MaxQuantity;

    /// <summary>
    /// Throws a <see cref="CheckoutException"/> for rows 1-4 of rules §4.1, in that order.
    /// Does nothing when the items are fine.
    /// </summary>
    public static void Validate(Basket basket, Catalog catalog)
    {
        // Row 1: no items at all.
        if (basket.Items.Count == 0)
        {
            throw new CheckoutException(ErrorCode.EmptyBasket, "The basket has no items");
        }

        // Row 2: every SKU must exist. All SKUs are checked before any quantity, so an unknown
        // product on line 3 still wins over a bad quantity on line 1.
        foreach (var item in basket.Items)
        {
            catalog.RequireProduct(item.Sku);
        }

        // Row 3: every quantity must be a whole number from 1 to 99.
        foreach (var item in basket.Items)
        {
            if (!IsValidQuantity(item.Quantity))
            {
                throw new CheckoutException(ErrorCode.InvalidQuantity, $"Invalid quantity for {item.Sku}: {item.Quantity}");
            }
        }

        // Row 4: the lines of one SKU, added together, may not go over 99 either.
        foreach (var line in ReceiptLines.Merge(basket.Items))
        {
            if (line.Quantity > MaxQuantity)
            {
                throw new CheckoutException(ErrorCode.InvalidQuantity, $"Too many {line.Sku} in total: {line.Quantity}");
            }
        }
    }
}
