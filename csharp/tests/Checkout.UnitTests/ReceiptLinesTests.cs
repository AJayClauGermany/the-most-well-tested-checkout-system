namespace CheckoutSystem.UnitTests;

/// <summary>Unit tests for merging basket items into receipt lines and pricing them (rules §4.2).</summary>
public class ReceiptLinesTests
{
    [Fact]
    public void Merge_DifferentSkus_KeepsThemApartInBasketOrder()
    {
        var merged = ReceiptLines.Merge([new BasketItem("PEN", 2), new BasketItem("APPLE", 1)]);

        // Records compare by value, so two lists of records can be compared element by element.
        Assert.Equal([new BasketItem("PEN", 2), new BasketItem("APPLE", 1)], merged);
    }

    [Fact]
    public void Merge_SameSkuTwice_AddsUpTheQuantitiesAtTheFirstAppearance()
    {
        // The example from rules §4.2: [APPLE x1, PEN x2, APPLE x3] gives [APPLE x4, PEN x2].
        var merged = ReceiptLines.Merge([new BasketItem("APPLE", 1), new BasketItem("PEN", 2), new BasketItem("APPLE", 3)]);

        Assert.Equal([new BasketItem("APPLE", 4), new BasketItem("PEN", 2)], merged);
    }

    [Fact]
    public void Merge_ThreeLinesOfOneSku_GivesOneLine()
    {
        var merged = ReceiptLines.Merge([new BasketItem("COFFEE", 1), new BasketItem("COFFEE", 1), new BasketItem("COFFEE", 1)]);

        Assert.Equal([new BasketItem("COFFEE", 3)], merged);
    }

    [Fact]
    public void Merge_NoItems_GivesNoLines()
    {
        Assert.Empty(ReceiptLines.Merge([]));
    }

    [Fact]
    public void Price_ProductWithoutPromotion_CopiesTheProductAndMultiplies()
    {
        var line = ReceiptLines.Price(Fixtures.Apple, quantity: 3);

        Assert.Equal(new ReceiptLine("APPLE", "Apple", Quantity: 3, UnitPrice: 50, LineTotal: 150, LineDiscount: 0), line);
    }

    [Fact]
    public void Price_ProductWithPromotion_IncludesThePromotionDiscount()
    {
        Assert.Equal(899, ReceiptLines.Price(Fixtures.Coffee, 3).LineDiscount);
        Assert.Equal(2697, ReceiptLines.Price(Fixtures.Coffee, 3).LineTotal);
        Assert.Equal(349, ReceiptLines.Price(Fixtures.Notebook, 10).LineDiscount);
        Assert.Equal(3490, ReceiptLines.Price(Fixtures.Notebook, 10).LineTotal);
    }

    [Fact]
    public void Price_BelowThePromotionThreshold_GivesNoDiscount()
    {
        Assert.Equal(0, ReceiptLines.Price(Fixtures.Coffee, 2).LineDiscount);
        Assert.Equal(0, ReceiptLines.Price(Fixtures.Notebook, 9).LineDiscount);
    }
}
