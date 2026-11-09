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

    [Theory]
    [InlineData("EMPTY_BASKET", ErrorCode.EmptyBasket)]
    [InlineData("COUPON_MINIMUM_NOT_MET", ErrorCode.CouponMinimumNotMet)]
    public void FromSpecCode_KnownText_ReturnsTheCode(string specCode, ErrorCode expected)
    {
        Assert.Equal(expected, ErrorCodeExtensions.FromSpecCode(specCode));
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
        // It also hands the exception back, so we can check the message helps a human.
        var exception = Assert.Throws<ArgumentException>(() => ErrorCodeExtensions.FromSpecCode("NOT_A_CODE"));

        Assert.StartsWith("'NOT_A_CODE' is not a known error code", exception.Message, StringComparison.Ordinal);
    }

    [Fact]
    public void ToSpecCode_ValueOutsideTheEnum_Throws()
    {
        // C# lets you cast any int to an enum; the switch must not quietly accept it.
        var exception = Assert.Throws<ArgumentOutOfRangeException>(() => ((ErrorCode)42).ToSpecCode());

        Assert.StartsWith("Unknown error code", exception.Message, StringComparison.Ordinal);
    }
}

/// <summary>Unit tests for <see cref="CheckoutException"/>: it carries the code, and a useful message.</summary>
public class CheckoutExceptionTests
{
    [Fact]
    public void Constructor_CodeOnly_UsesAGenericMessage()
    {
        var exception = new CheckoutException(ErrorCode.UnknownCoupon);

        Assert.Equal(ErrorCode.UnknownCoupon, exception.Code);
        Assert.Equal("Checkout failed: UNKNOWN_COUPON", exception.Message);
    }

    [Fact]
    public void Constructor_CodeAndMessage_KeepsBoth()
    {
        var exception = new CheckoutException(ErrorCode.EmptyBasket, "The basket has no items");

        Assert.Equal(ErrorCode.EmptyBasket, exception.Code);
        Assert.Equal("The basket has no items", exception.Message);
    }

    [Fact]
    public void CheckoutException_IsANormalException_SoItCanBeThrownAndCaught()
    {
        var exception = new CheckoutException(ErrorCode.EmptyBasket, "The basket has no items");

        // Any `catch (Exception)` block will catch it, because it derives from Exception.
        Assert.IsAssignableFrom<Exception>(exception);
    }
}
