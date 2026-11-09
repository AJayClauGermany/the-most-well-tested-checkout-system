namespace CheckoutSystem.UnitTests;

/// <summary>
/// Unit tests for the two kinds of promotion (rules §4.2). The tables are [Theory]s with
/// [InlineData]: the same check, many inputs, each reported as its own test.
/// </summary>
public class PromotionsTests
{
    private static readonly Promotion BuyTwoGetOne = new BuyXGetYFree(Buy: 2, Free: 1);
    private static readonly Promotion TenPercentOffFromTen = new BulkDiscount(MinQuantity: 10, Percent: 10);

    [Fact]
    public void Discount_NoPromotion_IsZero()
    {
        Assert.Equal(0, Promotions.Discount(promotion: null, quantity: 50, unitPrice: 899));
    }

    // Coffee is buy 2 get 1 free at $8.99.
    [Theory]
    [InlineData(1, 0)]
    [InlineData(2, 0)]
    [InlineData(3, 1)]
    [InlineData(4, 1)]
    [InlineData(5, 1)]
    [InlineData(6, 2)]
    [InlineData(99, 33)]
    public void Discount_BuyTwoGetOneFree_GivesOneFreePerCompleteGroupOfThree(int quantity, int freeItems)
    {
        Assert.Equal(freeItems * 899, Promotions.Discount(BuyTwoGetOne, quantity, unitPrice: 899));
    }

    [Theory]
    [InlineData(1, 1, 4, 200)] // buy 1 get 1: every second item is free
    [InlineData(1, 1, 5, 200)]
    [InlineData(3, 2, 4, 0)] // buy 3 get 2: groups of 5, 2 of each group free
    [InlineData(3, 2, 5, 200)]
    [InlineData(3, 2, 10, 400)]
    public void Discount_BuyXGetYFree_UsesThePromotionsOwnNumbers(int buy, int free, int quantity, int expected)
    {
        var promotion = new BuyXGetYFree(buy, free);

        Assert.Equal(expected, Promotions.Discount(promotion, quantity, unitPrice: 100));
    }

    [Fact]
    public void Discount_BuyXGetYFree_MultipliesByTheUnitPrice()
    {
        Assert.Equal(100, Promotions.Discount(BuyTwoGetOne, 3, unitPrice: 100));
        Assert.Equal(1, Promotions.Discount(BuyTwoGetOne, 3, unitPrice: 1));
    }

    // Notebooks are 10% off from 10 at $3.49.
    [Theory]
    [InlineData(9, 0)] // one short of the minimum
    [InlineData(10, 349)] // exactly the minimum: 10% of $34.90
    [InlineData(11, 384)] // 383.9 cents rounds up
    [InlineData(15, 524)] // exactly 523.5 cents rounds UP
    [InlineData(99, 3455)] // 10% of $345.51 is 3455.1
    public void Discount_BulkDiscount_TakesThePercentOffFromTheMinimumQuantity(int quantity, int expected)
    {
        Assert.Equal(expected, Promotions.Discount(TenPercentOffFromTen, quantity, unitPrice: 349));
    }

    [Theory]
    [InlineData(2, 100, 0)]
    [InlineData(3, 100, 75)]
    [InlineData(3, 101, 76)] // 75.75 → 76
    public void Discount_BulkDiscount_UsesThePromotionsOwnMinimumAndPercent(int quantity, int unitPrice, int expected)
    {
        var quarterOffFromThree = new BulkDiscount(MinQuantity: 3, Percent: 25);

        Assert.Equal(expected, Promotions.Discount(quarterOffFromThree, quantity, unitPrice));
    }
}
