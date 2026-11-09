namespace CheckoutSystem;

/// <summary>Coupons (rules §4.4): SAVE10 and WELCOME5.</summary>
public static class Coupons
{
    /// <summary>
    /// The coupon code as the shop knows it: trimmed and upper case.
    /// Missing, empty or blank all mean "no coupon" and give <c>null</c>.
    /// </summary>
    public static string? Normalize(string? coupon) =>
        string.IsNullOrWhiteSpace(coupon) ? null : coupon.Trim().ToUpperInvariant();

    /// <summary>
    /// The coupon discount in cents for <paramref name="afterPromotions"/> cents of shopping.
    /// </summary>
    /// <exception cref="CheckoutException">UNKNOWN_COUPON or COUPON_MINIMUM_NOT_MET (rules §4.1 rows 6 and 7).</exception>
    public static int Discount(string? coupon, int afterPromotions)
    {
        var code = Normalize(coupon);
        if (code is null)
        {
            return 0;
        }

        // The table from rules §4.4: what each coupon takes off, and the least the customer
        // must spend (after promotions) to use it.
        var (discount, minimumSpend) = code switch
        {
            "SAVE10" => (Money.RoundHalfUp(afterPromotions, 10), 0),
            "WELCOME5" => (500, 2500),
            _ => throw new CheckoutException(ErrorCode.UnknownCoupon, $"Unknown coupon: {code}"),
        };

        if (afterPromotions < minimumSpend)
        {
            throw new CheckoutException(
                ErrorCode.CouponMinimumNotMet,
                $"{code} needs a spend of at least {minimumSpend} cents, got {afterPromotions}");
        }

        return discount;
    }
}
