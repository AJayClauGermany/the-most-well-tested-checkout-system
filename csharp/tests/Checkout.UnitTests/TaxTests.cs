namespace CheckoutSystem.UnitTests;

/// <summary>Unit tests for tax (rules §4.6).</summary>
public class TaxTests
{
    [Fact]
    public void Percent_IsEight()
    {
        Assert.Equal(8, Tax.Percent);
    }

    [Theory]
    [InlineData(50, 4)] // 8% of $0.50 is exactly 4 cents
    [InlineData(120, 10)] // 9.6 cents rounds up
    [InlineData(240, 19)] // 19.2 cents rounds down
    [InlineData(625, 50)] // exactly 50 cents
    [InlineData(1706, 136)] // 136.48 rounds down
    [InlineData(11880, 950)] // 950.4 rounds down
    [InlineData(0, 0)] // no spend, no tax
    public void On_DiscountedSubtotal_IsEightPercentRoundedHalfUp(int amount, int expected)
    {
        Assert.Equal(expected, Tax.On(amount));
    }

    [Fact]
    public void On_PlainFractions_RoundTheUsualWay()
    {
        // 8% of a whole number of cents can never be exactly half a cent (x × 8 never ends in 50),
        // so the half-up rule is pinned down by the 10% coupon tests; here we check plain rounding.
        Assert.Equal(150, Tax.On(1875)); // exactly 150.00
        Assert.Equal(150, Tax.On(1881)); // 150.48
        Assert.Equal(151, Tax.On(1888)); // 151.04
    }
}
