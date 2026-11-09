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
        throw new NotImplementedException();

    /// <summary>
    /// Throws a <see cref="CheckoutException"/> for rows 1-4 of rules §4.1, in that order.
    /// Does nothing when the items are fine.
    /// </summary>
    public static void Validate(Basket basket, Catalog catalog) =>
        throw new NotImplementedException();
}
