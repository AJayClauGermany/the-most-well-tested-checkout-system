namespace CheckoutSystem;

/// <summary>Coupons (rules §4.4): SAVE10 and WELCOME5.</summary>
public static class Coupons
{
    /// <summary>
    /// The coupon code as the shop knows it: trimmed and upper case.
    /// Missing, empty or blank all mean "no coupon" and give <c>null</c>.
    /// </summary>
    public static string? Normalize(string? coupon) =>
        throw new NotImplementedException();

    /// <summary>
    /// The coupon discount in cents for <paramref name="afterPromotions"/> cents of shopping.
    /// </summary>
    /// <exception cref="CheckoutException">UNKNOWN_COUPON or COUPON_MINIMUM_NOT_MET (rules §4.1 rows 6 and 7).</exception>
    public static int Discount(string? coupon, int afterPromotions) =>
        throw new NotImplementedException();
}
