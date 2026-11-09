namespace CheckoutSystem.UnitTests;

/// <summary>
/// Unit tests for basket validation (rules §4.1, rows 1 to 4) and the order errors are
/// reported in. The check order matters: a basket can be wrong in several ways at once.
/// Boundaries are tested on BOTH sides (0 and 1, 99 and 100): a `&gt;` written as `&gt;=`
/// is the classic bug, and the one mutation testing looks for first.
/// </summary>
public class BasketValidatorTests
{
    private static readonly Catalog Catalog = Fixtures.TinyCatalog;

    [Fact]
    public void Limits_AreTheOnesFromTheRules()
    {
        Assert.Equal(1, BasketValidator.MinQuantity);
        Assert.Equal(99, BasketValidator.MaxQuantity);
    }

    [Theory]
    [InlineData(1)]
    [InlineData(2)]
    [InlineData(50)]
    [InlineData(98)]
    [InlineData(99)]
    public void IsValidQuantity_OneToNinetyNine_IsTrue(int quantity)
    {
        Assert.True(BasketValidator.IsValidQuantity(quantity));
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-1)]
    [InlineData(100)]
    [InlineData(101)]
    [InlineData(1000)]
    public void IsValidQuantity_OutsideOneToNinetyNine_IsFalse(int quantity)
    {
        Assert.False(BasketValidator.IsValidQuantity(quantity));
    }

    [Fact]
    public void Validate_NormalBasket_DoesNothing()
    {
        var basket = new Basket([new BasketItem("APPLE", 1), new BasketItem("COFFEE", 99)]);

        // No exception means the test passes; there is nothing to assert on a void method.
        BasketValidator.Validate(basket, Catalog);
    }

    [Fact]
    public void Validate_EmptyBasket_ThrowsEmptyBasket()
    {
        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(new Basket([]), Catalog));

        Assert.Equal(ErrorCode.EmptyBasket, exception.Code);
        Assert.Equal("The basket has no items", exception.Message);
    }

    [Fact]
    public void Validate_UnknownProduct_ThrowsUnknownProductNamingTheSku()
    {
        var basket = new Basket([new BasketItem("APPLE", 1), new BasketItem("BANANA", 1)]);

        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(basket, Catalog));

        Assert.Equal(ErrorCode.UnknownProduct, exception.Code);
        Assert.Equal("Unknown product: BANANA", exception.Message);
    }

    [Fact]
    public void Validate_TwoUnknownProducts_ReportsTheEarliestOne()
    {
        var basket = new Basket([new BasketItem("BANANA", 1), new BasketItem("CHERRY", 1)]);

        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(basket, Catalog));

        Assert.Equal("Unknown product: BANANA", exception.Message);
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-1)]
    [InlineData(100)]
    public void Validate_QuantityOutOfRange_ThrowsInvalidQuantity(int quantity)
    {
        var basket = Fixtures.BasketOf("APPLE", quantity);

        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(basket, Catalog));

        Assert.Equal(ErrorCode.InvalidQuantity, exception.Code);
        Assert.Equal($"Invalid quantity for APPLE: {quantity}", exception.Message);
    }

    [Fact]
    public void Validate_UnknownProductOnALaterLineThanABadQuantity_StillReportsUnknownProduct()
    {
        // Row 2 (unknown product) wins over row 3 (bad quantity) regardless of line position.
        var basket = new Basket([new BasketItem("APPLE", 0), new BasketItem("BANANA", 1)]);

        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(basket, Catalog));

        Assert.Equal(ErrorCode.UnknownProduct, exception.Code);
    }

    [Fact]
    public void Validate_UnknownProductOnAnEarlierLineThanABadQuantity_ReportsUnknownProduct()
    {
        var basket = new Basket([new BasketItem("BANANA", 1), new BasketItem("APPLE", 0)]);

        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(basket, Catalog));

        Assert.Equal(ErrorCode.UnknownProduct, exception.Code);
    }

    [Fact]
    public void Validate_MergedQuantityOverNinetyNine_ThrowsInvalidQuantity()
    {
        // Each line is fine on its own (50 and 50), but one SKU may not go over 99 in total (row 4).
        var basket = new Basket([new BasketItem("APPLE", 50), new BasketItem("APPLE", 50)]);

        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(basket, Catalog));

        Assert.Equal(ErrorCode.InvalidQuantity, exception.Code);
        Assert.Equal("Too many APPLE in total: 100", exception.Message);
    }

    [Fact]
    public void Validate_MergedQuantityExactlyNinetyNine_IsFine()
    {
        var basket = new Basket([new BasketItem("APPLE", 50), new BasketItem("APPLE", 49)]);

        BasketValidator.Validate(basket, Catalog);
    }

    [Fact]
    public void Validate_BadQuantityAndMergedOverflow_ReportsTheSingleLineFirst()
    {
        // Row 3 (a single bad line) is checked before row 4 (the merged total).
        var basket = new Basket([new BasketItem("APPLE", 100), new BasketItem("APPLE", 1)]);

        var exception = Assert.Throws<CheckoutException>(() => BasketValidator.Validate(basket, Catalog));

        Assert.Equal("Invalid quantity for APPLE: 100", exception.Message);
    }

    [Fact]
    public void Validate_BadShippingAndCoupon_AreNotItsJob()
    {
        // Rows 5-7 are checked elsewhere (Shipping and Coupons), so this must pass here.
        var basket = Fixtures.BasketOf("APPLE", 1, coupon: "NOPE", shipping: "overnight");

        BasketValidator.Validate(basket, Catalog);
    }
}
