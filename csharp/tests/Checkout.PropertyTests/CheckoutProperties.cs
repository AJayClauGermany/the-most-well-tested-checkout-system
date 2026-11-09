using FsCheck;
using FsCheck.Fluent;
using FsCheck.Xunit;

namespace CheckoutSystem.PropertyTests;

/// <summary>
/// PROPERTY-BASED TESTS with FsCheck.
///
/// An example test says "THIS basket gives THAT total". A property says something that must
/// be true for EVERY basket, and FsCheck generates 100 random baskets trying to break it –
/// then shrinks any failure to the smallest example. The properties are listed once for all
/// six languages in spec/properties.md; each test is named after its id.
///
/// [Property] replaces [Fact]: FsCheck fills in the parameters from the generators in
/// <see cref="Generators"/>. The method either asserts (like a normal test) or returns a
/// <see cref="Property"/>.
/// </summary>
[Properties(Arbitrary = [typeof(Generators)])]
public class SuccessfulCheckoutProperties
{
    private static readonly Catalog Catalog = Generators.Catalog;
    private static readonly int[] ShippingPrices = [0, 499, 999];

    [Property(DisplayName = "P1: total == subtotal − promotion_discount − coupon_discount + shipping + tax")]
    public void P1_TheTotalsAddUp(ValidBasket valid)
    {
        var r = Checkout.Calculate(valid.Basket, Catalog);

        Assert.Equal(r.Subtotal - r.PromotionDiscount - r.CouponDiscount + r.Shipping + r.Tax, r.Total);
    }

    [Property(DisplayName = "P2: nothing is negative and a discount never exceeds what it discounts")]
    public void P2_NothingIsNegativeAndDiscountsNeverExceedWhatTheyDiscount(ValidBasket valid)
    {
        var r = Checkout.Calculate(valid.Basket, Catalog);

        Assert.All(new[] { r.Subtotal, r.PromotionDiscount, r.CouponDiscount, r.Shipping, r.Tax, r.Total }, amount => Assert.True(amount >= 0));
        Assert.All(r.Lines, line =>
        {
            Assert.True(line.UnitPrice >= 0);
            Assert.True(line.LineTotal >= 0);
            Assert.InRange(line.LineDiscount, 0, line.LineTotal);
        });
        Assert.True(r.PromotionDiscount <= r.Subtotal);
        Assert.True(r.CouponDiscount <= r.Subtotal - r.PromotionDiscount);
    }

    [Property(DisplayName = "P3: the receipt lines explain the totals")]
    public void P3_TheReceiptLinesExplainTheTotals(ValidBasket valid)
    {
        var r = Checkout.Calculate(valid.Basket, Catalog);

        Assert.Equal(r.Lines.Sum(line => line.LineTotal), r.Subtotal);
        Assert.Equal(r.Lines.Sum(line => line.LineDiscount), r.PromotionDiscount);
        Assert.All(r.Lines, line => Assert.Equal(line.UnitPrice * line.Quantity, line.LineTotal));
    }

    [Property(DisplayName = "P4: shuffling the items gives the same totals")]
    public Property P4_ShufflingTheItemsGivesTheSameTotals(ValidBasket valid)
    {
        // Gen.Shuffle gives a random permutation; Prop.ForAll makes FsCheck drive it (so it shrinks).
        var shuffled = Gen.Shuffle(valid.Basket.Items.ToArray()).Select(items => valid.Basket with { Items = items });

        return Prop.ForAll(Arb.From(shuffled), basket =>
            Assert.Equal(Totals.From(Checkout.Calculate(valid.Basket, Catalog)), Totals.From(Checkout.Calculate(basket, Catalog))));
    }

    [Property(DisplayName = "P5: splitting one line into two lines of the same SKU gives the same receipt")]
    public Property P5_SplittingALineIsInvisible(ValidBasket valid, PositiveInt pick, PositiveInt cut)
    {
        var items = valid.Basket.Items;
        var candidates = Enumerable.Range(0, items.Count).Where(i => items[i].Quantity >= 2).ToList();

        // Baskets whose lines all have quantity 1 cannot be split: FsCheck skips ("discards") them.
        return (candidates.Count > 0).Implies(() =>
        {
            var index = candidates[pick.Get % candidates.Count];
            var line = items[index];
            var first = 1 + (cut.Get % (line.Quantity - 1)); // 1 .. quantity-1
            var split = items.Take(index)
                .Append(new BasketItem(line.Sku, first))
                .Concat(items.Skip(index + 1))
                .Append(new BasketItem(line.Sku, line.Quantity - first))
                .ToList();

            Assert.Equal(Checkout.Calculate(valid.Basket, Catalog), Checkout.Calculate(valid.Basket with { Items = split }, Catalog));
        });
    }

    [Property(DisplayName = "P6: the shipping table is applied exactly")]
    public void P6_TheShippingTableIsAppliedExactly(ValidBasket valid)
    {
        var r = Checkout.Calculate(valid.Basket, Catalog);
        var discounted = r.Subtotal - r.PromotionDiscount - r.CouponDiscount;

        Assert.Contains(r.Shipping, ShippingPrices);
        if (valid.Basket.Shipping == "express")
        {
            Assert.Equal(999, r.Shipping);
        }
        else
        {
            Assert.Equal(discounted >= 5000 ? 0 : 499, r.Shipping);
        }
    }

    [Property(DisplayName = "P7: tax == (discounted_subtotal × 8 + 50) ÷ 100 with integer division")]
    public void P7_TaxIsEightPercentRoundedHalfUp(ValidBasket valid)
    {
        var r = Checkout.Calculate(valid.Basket, Catalog);
        var discounted = r.Subtotal - r.PromotionDiscount - r.CouponDiscount;

        Assert.Equal(((discounted * 8) + 50) / 100, r.Tax);
    }

    [Property(DisplayName = "P8: writing the coupon in another case or with surrounding spaces gives the same receipt")]
    public Property P8_CouponsAreCaseInsensitiveAndTrimmed(ValidBasket valid)
    {
        var code = Coupons.Normalize(valid.Basket.Coupon);

        // Only baskets WITH a coupon say anything about coupon spelling; the rest are discarded.
        return (code is not null).Implies(() =>
            Prop.ForAll(Arb.From(Generators.Sloppy(code!)), sloppy =>
                Assert.Equal(Checkout.Calculate(valid.Basket, Catalog), Checkout.Calculate(valid.Basket with { Coupon = sloppy }, Catalog))));
    }

    [Property(DisplayName = "P9: one line per distinct SKU, in order of first appearance, each with quantity 1–99")]
    public void P9_LinesAreMergedAndOrdered(ValidBasket valid)
    {
        var r = Checkout.Calculate(valid.Basket, Catalog);
        var distinctInOrder = valid.Basket.Items.Select(item => item.Sku).Distinct();

        Assert.Equal(distinctInOrder, r.Lines.Select(line => line.Sku));
        Assert.All(r.Lines, line =>
        {
            Assert.InRange(line.Quantity, 1, 99);
            Assert.Equal(valid.Basket.Items.Where(item => item.Sku == line.Sku).Sum(item => item.Quantity), line.Quantity);
        });
    }
}

[Properties(Arbitrary = [typeof(Generators)])]
public class AnyInputProperties
{
    private static readonly Catalog Catalog = Generators.Catalog;

    [Property(DisplayName = "P10: for ANY basket, checkout gives a schema-valid receipt or exactly one of the six error codes – never a crash", MaxTest = 500)]
    public void P10_ValidationIsTotal(AnyBasketAtAll any)
    {
        try
        {
            var receipt = Checkout.Calculate(any.Basket, Catalog);
            Assert.Empty(SchemaValidator.ErrorsFor("receipt.schema.json", receipt));
        }
        catch (CheckoutException exception)
        {
            Assert.Contains(exception.Code, Enum.GetValues<ErrorCode>());
        }
    }

    [Property(DisplayName = "P11: round_half_up is within half a cent of the exact value, and exact when the value is whole")]
    public void P11_RoundingIsCorrectForEveryInput(NonNegativeInt amount, Percent percent)
    {
        var cents = amount.Get % 10_000_000;

        var rounded = Money.RoundHalfUp(cents, percent.Value);
        var exactTimes100 = (long)cents * percent.Value; // the exact value is this ÷ 100

        Assert.True(Math.Abs((rounded * 100L) - exactTimes100) <= 50);
        if (exactTimes100 % 100 == 0)
        {
            Assert.Equal(exactTimes100, rounded * 100L);
        }

        if (exactTimes100 % 100 == 50)
        {
            Assert.Equal(exactTimes100 + 50, rounded * 100L); // halves go UP
        }
    }
}

[Properties(Arbitrary = [typeof(Generators)])]
public class SavingOrdersProperties
{
    private static readonly Catalog Catalog = Generators.Catalog;

    [Property(DisplayName = "P12: after saving n orders the ids are exactly ORDER-1 … ORDER-n, in that order")]
    public void P12_IdsAreSequential(ValidBasket[] baskets)
    {
        var repository = new InMemoryOrderRepository(Clocks.Fixed());

        var ids = baskets.Select(valid => Checkout.PlaceOrder(valid.Basket, Catalog, repository).OrderId);

        Assert.Equal(Enumerable.Range(1, baskets.Length).Select(n => $"ORDER-{n}"), ids);
    }

    [Property(DisplayName = "P13: find(save(receipt)) gives back the same receipt and the clock's time")]
    public void P13_NothingIsLostOnTheWayInAndOut(ValidBasket valid, DateTime when)
    {
        var moment = new DateTimeOffset(DateTime.SpecifyKind(when, DateTimeKind.Utc));
        var repository = new InMemoryOrderRepository(new Microsoft.Extensions.Time.Testing.FakeTimeProvider(moment));
        var receipt = Checkout.Calculate(valid.Basket, Catalog);

        var order = repository.Find(repository.Save(receipt));

        Assert.NotNull(order);
        Assert.Equal(receipt, order.Receipt);
        Assert.Equal(moment.ToString("yyyy-MM-dd'T'HH:mm:ss'Z'", System.Globalization.CultureInfo.InvariantCulture), order.CreatedAt);
    }

    [Property(DisplayName = "P14: a failed checkout saves nothing, and find of an unknown id is 'not found'")]
    public Property P14_FailuresLeaveNoTrace(AnyBasketAtAll any, string randomId)
    {
        return (Generators.TryCheckout(any.Basket) is null).Implies(() =>
        {
            var repository = new InMemoryOrderRepository(Clocks.Fixed());

            Assert.Throws<CheckoutException>(() => Checkout.PlaceOrder(any.Basket, Catalog, repository));

            Assert.Equal(0, repository.Count);
            Assert.Null(repository.Find(randomId ?? ""));
            Assert.Null(repository.Find("ORDER-1"));
        });
    }
}
