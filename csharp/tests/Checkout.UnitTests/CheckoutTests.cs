namespace CheckoutSystem.UnitTests;

/// <summary>
/// Unit tests of <see cref="Checkout.Calculate"/> as a whole, with the tiny made-up catalog
/// from Checkout.TestSupport. The 60 shared spec cases (<see cref="SpecCaseTests"/>) cover
/// the standard catalog; these tests check the pieces fit together and the things the totals
/// alone cannot show, such as the receipt lines.
/// </summary>
public class CheckoutTests
{
    private static readonly Catalog Catalog = Fixtures.TinyCatalog;

    [Fact]
    public void Calculate_OneApple_GivesTheWholeReceipt()
    {
        var receipt = Checkout.Calculate(Fixtures.BasketOf("APPLE", 1), Catalog);

        // 50 + $4.99 shipping + 4 tax = $5.53. Records compare by value, lines included.
        Assert.Equal(Fixtures.OneAppleReceipt, receipt);
    }

    [Fact]
    public void Calculate_SameSkuOnTwoLines_GivesOneLinePerSkuInOrderOfFirstAppearance()
    {
        var basket = new Basket([new BasketItem("PEN", 1), new BasketItem("APPLE", 2), new BasketItem("PEN", 3)]);

        var receipt = Checkout.Calculate(basket, Catalog);

        Assert.Equal(["PEN", "APPLE"], receipt.Lines.Select(line => line.Sku));
        Assert.Equal([4, 2], receipt.Lines.Select(line => line.Quantity));
    }

    [Fact]
    public void Calculate_TwoPromotions_WorksTheTotalsOutFromTheLines()
    {
        // 3 coffees ($26.97, one free: -$8.99) and 10 notebooks ($34.90, 10% off: -$3.49)
        var basket = new Basket([new BasketItem("COFFEE", 3), new BasketItem("NOTEBOOK", 10)]);

        var receipt = Checkout.Calculate(basket, Catalog);

        Assert.Equal(2697 + 3490, receipt.Subtotal);
        Assert.Equal(899 + 349, receipt.PromotionDiscount);
        Assert.Equal(0, receipt.CouponDiscount);
        Assert.Equal(499, receipt.Shipping); // $49.39 after promotions: not free
        Assert.Equal(395, receipt.Tax); // 8% of 4939 = 395.12
        Assert.Equal(4939 + 499 + 395, receipt.Total);
    }

    [Fact]
    public void Calculate_WithCoupon_TakesItOffAfterPromotionsAndBeforeShippingAndTax()
    {
        // TV $100.00 with SAVE10: $90.00 → free shipping, tax 8% of $90.00 = $7.20
        var receipt = Checkout.Calculate(Fixtures.BasketOf("TV", 1, coupon: "SAVE10"), Catalog);

        Assert.Equal(new Totals(Subtotal: 10000, PromotionDiscount: 0, CouponDiscount: 1000, Shipping: 0, Tax: 720, Total: 9720), Totals.From(receipt));
    }

    [Fact]
    public void Calculate_ExpressShipping_IsChargedOnTopAndNotTaxed()
    {
        var receipt = Checkout.Calculate(Fixtures.BasketOf("TV", 1, shipping: "express"), Catalog);

        Assert.Equal(999, receipt.Shipping);
        Assert.Equal(800, receipt.Tax);
        Assert.Equal(10000 + 999 + 800, receipt.Total);
    }

    [Fact]
    public void Calculate_FreeShippingThreshold_IsJudgedAfterTheCoupon()
    {
        // $50.00 of pens is free shipping; SAVE10 pulls it to $45.00, which is not.
        var withoutCoupon = Checkout.Calculate(new Basket([new BasketItem("PEN", 40), new BasketItem("APPLE", 4)]), Catalog);
        var withCoupon = Checkout.Calculate(new Basket([new BasketItem("PEN", 40), new BasketItem("APPLE", 4)], Coupon: "SAVE10"), Catalog);

        Assert.Equal(0, withoutCoupon.Shipping);
        Assert.Equal(499, withCoupon.Shipping);
    }

    // The first failing row of rules §4.1 is the one reported, even when several rows fail.
    public static TheoryData<string, Basket, ErrorCode> BasketsFailingSeveralRows() => new()
    {
        { "empty basket before anything else", new Basket([], Coupon: "NOPE", Shipping: "overnight"), ErrorCode.EmptyBasket },
        { "unknown product before bad quantity, shipping and coupon", new Basket([new BasketItem("APPLE", 0), new BasketItem("BANANA", 1)], "NOPE", "x"), ErrorCode.UnknownProduct },
        { "bad quantity before shipping and coupon", new Basket([new BasketItem("APPLE", 0)], "NOPE", "x"), ErrorCode.InvalidQuantity },
        { "merged quantity over 99 before shipping and coupon", new Basket([new BasketItem("APPLE", 60), new BasketItem("APPLE", 40)], "NOPE", "x"), ErrorCode.InvalidQuantity },
        { "bad shipping before bad coupon", new Basket([new BasketItem("APPLE", 1)], "NOPE", "x"), ErrorCode.UnknownShippingMethod },
        { "bad shipping before an unmet coupon minimum", new Basket([new BasketItem("APPLE", 1)], "WELCOME5", "x"), ErrorCode.UnknownShippingMethod },
        { "unknown coupon", new Basket([new BasketItem("APPLE", 1)], "NOPE"), ErrorCode.UnknownCoupon },
        { "coupon minimum judged after promotions: 3 coffees are $17.98 after the free one", new Basket([new BasketItem("COFFEE", 3)], "WELCOME5"), ErrorCode.CouponMinimumNotMet },
    };

    [Theory]
    [MemberData(nameof(BasketsFailingSeveralRows))]
    public void Calculate_BasketFailingSeveralRows_ReportsTheFirstRowOfTheRules(string why, Basket basket, ErrorCode expected)
    {
        var exception = Assert.Throws<CheckoutException>(() => Checkout.Calculate(basket, Catalog));

        Assert.True(expected == exception.Code, $"{why}: expected {expected} but got {exception.Code}");
    }
}
