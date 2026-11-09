namespace CheckoutSystem.UnitTests;

/// <summary>
/// Unit tests for rounding (rules §7). Rounding is where checkout systems classically go
/// wrong, so it gets its own small function and its own tests.
/// </summary>
public class MoneyTests
{
    // The table from rules §7, as a parameterised test: one row, one test.
    [Theory]
    [InlineData(224, 10, 22)] // 22.4 rounds down
    [InlineData(225, 10, 23)] // 22.5 rounds UP
    [InlineData(226, 10, 23)] // 22.6 rounds up
    public void RoundHalfUp_TableFromTheRules_RoundsToTheNearestCentWithHalvesUp(int amount, int percent, int expected)
    {
        Assert.Equal(expected, Money.RoundHalfUp(amount, percent));
    }

    [Fact]
    public void RoundHalfUp_ExactHalf_RoundsUpWhereBankersRoundingWouldGoToEven()
    {
        // 10% of $52.35 is exactly 523.5 cents. Math.Round (banker's rounding) gives 524 too,
        // but for 22.5 (above) it would give 22 - both together pin the rule down.
        Assert.Equal(524, Money.RoundHalfUp(5235, 10));
    }

    [Fact]
    public void RoundHalfUp_WholeResult_IsNotChanged()
    {
        Assert.Equal(50, Money.RoundHalfUp(500, 10));
        Assert.Equal(4, Money.RoundHalfUp(50, 8));
    }

    [Fact]
    public void RoundHalfUp_EightPercent_RoundsNineAndAHalfUpAndNineteenPointTwoDown()
    {
        Assert.Equal(10, Money.RoundHalfUp(120, 8)); // 9.6 → 10
        Assert.Equal(19, Money.RoundHalfUp(240, 8)); // 19.2 → 19
    }

    [Fact]
    public void RoundHalfUp_ZeroAmountOrHundredPercent_GivesZeroOrTheWholeAmount()
    {
        Assert.Equal(0, Money.RoundHalfUp(0, 10));
        Assert.Equal(899, Money.RoundHalfUp(899, 100));
    }

    [Fact]
    public void RoundHalfUp_ZeroPercent_GivesZero()
    {
        Assert.Equal(0, Money.RoundHalfUp(899, 0));
    }

    [Fact]
    public void RoundHalfUp_AmountsWhereFloatingPointDrifts_StaysExact()
    {
        // 0.1 + 0.2 != 0.3 with doubles. 10% of 3 cents is 0.3 - it must round to 0, not 1.
        Assert.Equal(0, Money.RoundHalfUp(3, 10));
        // 8% of 1050 = 84.0 exactly; a double implementation might land on 83.99999.
        Assert.Equal(84, Money.RoundHalfUp(1050, 8));
    }
}
