namespace CheckoutSystem.TestSupport;

/// <summary>
/// Small, hand-made test data for unit tests.
///
/// Unit tests prefer tiny made-up data over the real catalog: when a test fails you can
/// work out the right answer in your head. The shared spec cases use the real catalog.
/// </summary>
public static class Fixtures
{
    public static readonly Product Apple = new("APPLE", "Apple", UnitPrice: 50);
    public static readonly Product Pen = new("PEN", "Pen", UnitPrice: 120);
    public static readonly Product Coffee = new("COFFEE", "Bag of Coffee Beans", UnitPrice: 899, new BuyXGetYFree(Buy: 2, Free: 1));
    public static readonly Product Notebook = new("NOTEBOOK", "Notebook", UnitPrice: 349, new BulkDiscount(MinQuantity: 10, Percent: 10));

    /// <summary>A $100.00 item: handy for reaching free shipping and coupon minimums quickly.</summary>
    public static readonly Product Tv = new("TV", "Television", UnitPrice: 10000);

    /// <summary>Five products, two of them with a promotion.</summary>
    public static readonly Catalog TinyCatalog = new([Apple, Pen, Coffee, Notebook, Tv]);

    /// <summary>A valid receipt for one apple, used wherever "some receipt" is needed.</summary>
    public static readonly Receipt OneAppleReceipt = new(
        Lines: [new ReceiptLine("APPLE", "Apple", Quantity: 1, UnitPrice: 50, LineTotal: 50, LineDiscount: 0)],
        Subtotal: 50,
        PromotionDiscount: 0,
        CouponDiscount: 0,
        Shipping: 499,
        Tax: 4,
        Total: 553);

    /// <summary>The basket behind spec/examples/receipt.json (rules §5): 3 coffees, 2 apples, SAVE10.</summary>
    public static readonly Basket ExampleBasket = new(
        [new BasketItem("COFFEE", 3), new BasketItem("APPLE", 2)],
        Coupon: "SAVE10",
        Shipping: "standard");

    /// <summary>A basket with one line. Saves a few keystrokes in tests.</summary>
    public static Basket BasketOf(string sku, int quantity, string? coupon = null, string shipping = "standard") =>
        new([new BasketItem(sku, quantity)], coupon, shipping);
}
