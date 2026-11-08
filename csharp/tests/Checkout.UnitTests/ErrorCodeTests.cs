namespace CheckoutSystem.UnitTests;

/// <summary>
/// EXAMPLE PARAMETERISED TESTS.
///
/// A [Theory] is a test that runs once for every [InlineData] row. Use it when the
/// same check applies to many inputs – here, every error code from rules §9.
/// Each row shows up as its own test in the results, e.g.
/// ToSpecCode_EveryErrorCode_MatchesTheSpecSpelling(code: EmptyBasket, expected: "EMPTY_BASKET").
/// </summary>
public class ErrorCodeTests
{
    [Theory]
    [InlineData(ErrorCode.EmptyBasket, "EMPTY_BASKET")]
    [InlineData(ErrorCode.UnknownProduct, "UNKNOWN_PRODUCT")]
    [InlineData(ErrorCode.InvalidQuantity, "INVALID_QUANTITY")]
    [InlineData(ErrorCode.UnknownShippingMethod, "UNKNOWN_SHIPPING_METHOD")]
    [InlineData(ErrorCode.UnknownCoupon, "UNKNOWN_COUPON")]
    [InlineData(ErrorCode.CouponMinimumNotMet, "COUPON_MINIMUM_NOT_MET")]
    public void ToSpecCode_EveryErrorCode_MatchesTheSpecSpelling(ErrorCode code, string expected)
    {
        var actual = code.ToSpecCode();

        // Assert.Equal takes the EXPECTED value first, then the ACTUAL value.
        // Getting them the wrong way round makes failure messages confusing.
        Assert.Equal(expected, actual);
    }

    [Fact]
    public void ErrorCode_HasExactlyTheSixCodesFromTheSpec()
    {
        // If someone adds a seventh code, this test reminds them to update rules §9 too.
        Assert.Equal(6, Enum.GetValues<ErrorCode>().Length);
    }

    [Fact]
    public void FromSpecCode_UnknownText_Throws()
    {
        // Assert.Throws checks that the code inside the lambda (() => ...) throws
        // this exact exception type. The test fails if nothing is thrown.
        Assert.Throws<ArgumentException>(() => ErrorCodeExtensions.FromSpecCode("NOT_A_CODE"));
    }
}
