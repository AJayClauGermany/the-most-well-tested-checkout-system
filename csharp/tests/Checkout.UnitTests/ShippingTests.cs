namespace CheckoutSystem.UnitTests;

/// <summary>Unit tests for shipping (rules §4.5).</summary>
public class ShippingTests
{
    [Fact]
    public void Parse_KnownMethods_AreAccepted()
    {
        Assert.Equal(ShippingMethod.Standard, Shipping.Parse("standard"));
        Assert.Equal(ShippingMethod.Express, Shipping.Parse("express"));
    }

    [Theory]
    [InlineData("overnight")]
    [InlineData("Express")]
    [InlineData("STANDARD")]
    [InlineData(" standard")]
    [InlineData("express ")]
    [InlineData("")]
    public void Parse_AnythingElse_ThrowsUnknownShippingMethod(string shipping)
    {
        // Shipping method names are matched exactly (lower case), unlike coupon codes.
        var exception = Assert.Throws<CheckoutException>(() => Shipping.Parse(shipping));

        Assert.Equal(ErrorCode.UnknownShippingMethod, exception.Code);
    }

    [Fact]
    public void Parse_UnknownMethod_NamesItInTheMessage()
    {
        var exception = Assert.Throws<CheckoutException>(() => Shipping.Parse("overnight"));

        Assert.Equal("Unknown shipping method: overnight", exception.Message);
    }

    [Fact]
    public void Constants_AreTheOnesFromTheRules()
    {
        Assert.Equal(499, Shipping.StandardCost);
        Assert.Equal(999, Shipping.ExpressCost);
        Assert.Equal(5000, Shipping.FreeStandardShippingFrom);
    }

    [Theory]
    [InlineData(0, 499)]
    [InlineData(4999, 499)] // one cent below the threshold: still paid
    [InlineData(5000, 0)] // exactly $50.00: free
    [InlineData(5001, 0)]
    [InlineData(100000, 0)]
    public void Cost_Standard_IsFreeFromFiftyDollars(int discountedSubtotal, int expected)
    {
        Assert.Equal(expected, Shipping.Cost(ShippingMethod.Standard, discountedSubtotal));
    }

    [Theory]
    [InlineData(0)]
    [InlineData(4999)]
    [InlineData(5000)]
    [InlineData(100000)]
    public void Cost_Express_AlwaysCostsNineNinetyNine(int discountedSubtotal)
    {
        Assert.Equal(999, Shipping.Cost(ShippingMethod.Express, discountedSubtotal));
    }
}
