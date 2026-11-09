namespace CheckoutSystem.UnitTests;

/// <summary>Unit tests for coupons (rules §4.4).</summary>
public class CouponsTests
{
    [Theory]
    [InlineData("SAVE10", "SAVE10")]
    [InlineData("save10", "SAVE10")]
    [InlineData("  save10 ", "SAVE10")]
    [InlineData("\tWelcome5\n", "WELCOME5")]
    [InlineData("nope", "NOPE")]
    public void Normalize_AnyText_TrimsAndUpperCases(string input, string expected)
    {
        Assert.Equal(expected, Coupons.Normalize(input));
    }

    [Theory]
    [InlineData(null)]
    [InlineData("")]
    [InlineData("   ")]
    [InlineData("\t\n")]
    public void Normalize_MissingOrBlank_MeansNoCoupon(string? input)
    {
        Assert.Null(Coupons.Normalize(input));
    }

    [Fact]
    public void Discount_NoCoupon_IsZero()
    {
        Assert.Equal(0, Coupons.Discount(null, 10000));
        Assert.Equal(0, Coupons.Discount("", 10000));
        Assert.Equal(0, Coupons.Discount("  ", 10000));
    }

    [Theory]
    [InlineData(500, 50)]
    [InlineData(225, 23)] // 22.5 → 23
    [InlineData(224, 22)] // 22.4 → 22
    [InlineData(1798, 180)] // 179.8 → 180
    [InlineData(1, 0)] // no minimum spend: 0.1 → 0
    [InlineData(5, 1)] // 0.5 → 1
    public void Discount_Save10_TakesTenPercentOffRoundedHalfUp(int afterPromotions, int expected)
    {
        Assert.Equal(expected, Coupons.Discount("SAVE10", afterPromotions));
    }

    [Fact]
    public void Discount_Save10InAnyCaseWithSpaces_IsTheSameCoupon()
    {
        Assert.Equal(50, Coupons.Discount(" save10 ", 500));
    }

    [Fact]
    public void Discount_Welcome5_TakesExactlyFiveDollarsOff()
    {
        Assert.Equal(500, Coupons.Discount("WELCOME5", 2500));
        Assert.Equal(500, Coupons.Discount("WELCOME5", 99999));
        Assert.Equal(500, Coupons.Discount("welcome5", 2500));
    }

    [Fact]
    public void Discount_Welcome5OneCentBelowTheMinimum_ThrowsCouponMinimumNotMet()
    {
        // Boundary: 2500 works (above), 2499 does not.
        var exception = Assert.Throws<CheckoutException>(() => Coupons.Discount("WELCOME5", 2499));

        Assert.Equal(ErrorCode.CouponMinimumNotMet, exception.Code);
        Assert.Equal("WELCOME5 needs a spend of at least 2500 cents, got 2499", exception.Message);
    }

    [Fact]
    public void Discount_UnknownCode_ThrowsUnknownCoupon()
    {
        var exception = Assert.Throws<CheckoutException>(() => Coupons.Discount("freestuff", 10000));

        Assert.Equal(ErrorCode.UnknownCoupon, exception.Code);
        Assert.Equal("Unknown coupon: FREESTUFF", exception.Message);
    }

    [Theory]
    [InlineData("SAVE100")]
    [InlineData("WELCOME50")]
    [InlineData("SAVE1")]
    public void Discount_CodeThatMerelyStartsWithARealCode_IsUnknown(string coupon)
    {
        var exception = Assert.Throws<CheckoutException>(() => Coupons.Discount(coupon, 10000));

        Assert.Equal(ErrorCode.UnknownCoupon, exception.Code);
    }
}
