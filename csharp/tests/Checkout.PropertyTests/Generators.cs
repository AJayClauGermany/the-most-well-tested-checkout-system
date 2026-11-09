using FsCheck;
using FsCheck.Fluent;

namespace CheckoutSystem.PropertyTests;

/// <summary>
/// Generators ("arbitraries") for random baskets, following the "Generating baskets" section
/// of spec/properties.md. FsCheck's <c>[Property(Arbitrary = [typeof(Generators)])]</c> looks
/// here for a method returning <c>Arbitrary&lt;T&gt;</c> whenever a test asks for a <c>T</c>.
/// </summary>
public static class Generators
{
    public static readonly Catalog Catalog = SpecLoader.LoadCatalog();

    private static readonly string[] Skus = Catalog.Products.Select(product => product.Sku).ToArray();

    /// <summary>A sensible item: a real SKU, a quantity from 1 to 99.</summary>
    public static Gen<BasketItem> Item() =>
        from sku in Gen.Elements(Skus)
        from quantity in Gen.Choose(1, 99)
        select new BasketItem(sku, quantity);

    /// <summary>SAVE10 or WELCOME5 written any old way: mixed case, spaces around it.</summary>
    public static Gen<string> SloppyCoupon() =>
        from code in Gen.Elements("SAVE10", "WELCOME5")
        from sloppy in Sloppy(code)
        select sloppy;

    /// <summary><paramref name="code"/> in random case with random spaces around it.</summary>
    public static Gen<string> Sloppy(string code) =>
        from upper in Gen.ArrayOf(ArbMap.Default.GeneratorFor<bool>(), code.Length)
        from before in Gen.Choose(0, 3)
        from after in Gen.Choose(0, 3)
        select new string(' ', before)
               + string.Concat(code.Select((letter, i) => upper[i] ? char.ToUpperInvariant(letter) : char.ToLowerInvariant(letter)))
               + new string(' ', after);

    /// <summary>Missing, null-ish, a real code, a sloppy real code, or a random string.</summary>
    private static Gen<string?> Coupon() =>
        Gen.Frequency(
            (3, Gen.Constant<string?>(null)),
            (2, Gen.Constant<string?>("SAVE10")),
            (2, Gen.Constant<string?>("WELCOME5")),
            (2, SloppyCoupon().Select(coupon => (string?)coupon)),
            (1, ArbMap.Default.GeneratorFor<NonNull<string>>().Select(s => (string?)s.Get)));

    /// <summary>Missing (= standard), standard, express, or something wrong.</summary>
    private static Gen<string> ShippingMethod() =>
        Gen.Frequency(
            (4, Gen.Constant("standard")),
            (3, Gen.Constant("express")),
            (1, Gen.Elements("overnight", "Express", "")));

    /// <summary>A whole basket: 1 to 8 items (the same SKU may repeat), maybe a coupon, maybe a shipping method.</summary>
    public static Gen<Basket> Basket() =>
        from count in Gen.Choose(1, 8)
        from items in Gen.ListOf(Item(), count)
        from coupon in Coupon()
        from shipping in ShippingMethod()
        select new Basket(items.ToList(), coupon, shipping);

    /// <summary>Anything at all that still fits the Basket type: nonsense SKUs, silly quantities, random strings.</summary>
    public static Gen<Basket> AnyBasket() =>
        from count in Gen.Choose(0, 6)
        from items in Gen.ListOf(AnyItem(), count)
        from coupon in Gen.OneOf(Coupon(), ArbMap.Default.GeneratorFor<string>().Select(s => (string?)s))
        from shipping in Gen.OneOf(ShippingMethod(), ArbMap.Default.GeneratorFor<NonNull<string>>().Select(s => s.Get))
        select new Basket(items.ToList(), coupon, shipping);

    private static Gen<BasketItem> AnyItem() =>
        from sku in Gen.OneOf(Gen.Elements(Skus), ArbMap.Default.GeneratorFor<NonNull<string>>().Select(s => s.Get))
        from quantity in Gen.OneOf(Gen.Choose(-5, 150), ArbMap.Default.GeneratorFor<int>())
        select new BasketItem(sku, quantity);

    /// <summary>Runs checkout; <c>null</c> means the basket was invalid (a CheckoutException was thrown).</summary>
    public static Receipt? TryCheckout(Basket basket)
    {
        try
        {
            return Checkout.Calculate(basket, Catalog);
        }
        catch (CheckoutException)
        {
            return null;
        }
    }

    // --- The Arbitrary<T> methods FsCheck looks up by return type ---------------------------

    /// <summary>Only baskets that check out successfully – what P1 to P9 are about.</summary>
    public static Arbitrary<ValidBasket> ValidBaskets() =>
        Arb.From(Basket().Where(basket => TryCheckout(basket) is not null).Select(basket => new ValidBasket(basket)));

    public static Arbitrary<AnyBasketAtAll> AnyBaskets() =>
        Arb.From(AnyBasket().Select(basket => new AnyBasketAtAll(basket)));

    public static Arbitrary<Percent> Percents() =>
        Arb.From(Gen.Choose(0, 100).Select(percent => new Percent(percent)));
}

/// <summary>A basket that checks out successfully (wrapper type so FsCheck knows which generator to use).</summary>
public sealed record ValidBasket(Basket Basket);

/// <summary>Any basket at all, valid or not.</summary>
public sealed record AnyBasketAtAll(Basket Basket);

/// <summary>A percentage from 0 to 100.</summary>
public readonly record struct Percent(int Value);
