namespace CheckoutSystem.UnitTests;

/// <summary>
/// Unit tests for comparing receipts. A C# record compares its fields by value, but a list
/// field only by reference ("the very same list object?"). <see cref="Receipt"/> therefore
/// spells out its own equality; these tests pin it down from both sides – equal when it
/// should be, and NOT equal when any one part differs.
/// </summary>
public class ReceiptTests
{
    private static readonly Receipt Original = Fixtures.OneAppleReceipt;

    [Fact]
    public void Equals_SameValuesInADifferentListObject_IsEqual()
    {
        var copy = Original with { Lines = [.. Original.Lines] };

        Assert.NotSame(Original.Lines, copy.Lines);
        Assert.Equal(Original, copy);
        Assert.Equal(Original.GetHashCode(), copy.GetHashCode());
    }

    [Fact]
    public void Equals_DifferentLine_IsNotEqual()
    {
        var other = Original with { Lines = [Original.Lines[0] with { Quantity = 2 }] };

        Assert.NotEqual(Original, other);
    }

    [Fact]
    public void Equals_ExtraLine_IsNotEqual()
    {
        var other = Original with { Lines = [.. Original.Lines, Original.Lines[0]] };

        Assert.NotEqual(Original, other);
    }

    [Theory]
    [InlineData("subtotal")]
    [InlineData("promotion")]
    [InlineData("coupon")]
    [InlineData("shipping")]
    [InlineData("tax")]
    [InlineData("total")]
    public void Equals_AnyOneTotalDifferent_IsNotEqual(string which)
    {
        var other = which switch
        {
            "subtotal" => Original with { Subtotal = Original.Subtotal + 1 },
            "promotion" => Original with { PromotionDiscount = Original.PromotionDiscount + 1 },
            "coupon" => Original with { CouponDiscount = Original.CouponDiscount + 1 },
            "shipping" => Original with { Shipping = Original.Shipping + 1 },
            "tax" => Original with { Tax = Original.Tax + 1 },
            _ => Original with { Total = Original.Total + 1 },
        };

        Assert.NotEqual(Original, other);
    }

    [Fact]
    public void Equals_Null_IsNotEqual()
    {
        Assert.False(Original.Equals(null));
    }
}
