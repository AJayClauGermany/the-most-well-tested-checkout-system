namespace CheckoutSystem;

/// <summary>The two shipping methods the shop offers (rules §4.5).</summary>
public enum ShippingMethod
{
    Standard,
    Express,
}

/// <summary>Shipping (rules §4.5).</summary>
public static class Shipping
{
    public const int StandardCost = 499;
    public const int ExpressCost = 999;

    /// <summary>Standard shipping is free from this discounted subtotal upwards.</summary>
    public const int FreeStandardShippingFrom = 5000;

    /// <summary>
    /// Turns the basket's <c>shipping</c> text into a known method. "Parse, don't validate":
    /// after this, the type system knows the value is one of the two real methods.
    /// </summary>
    /// <exception cref="CheckoutException">UNKNOWN_SHIPPING_METHOD for anything but "standard" or "express" (rules §4.1 row 5).</exception>
    public static ShippingMethod Parse(string shipping) =>
        throw new NotImplementedException();

    /// <summary>The shipping cost in cents for this method and this discounted subtotal.</summary>
    public static int Cost(ShippingMethod method, int discountedSubtotal) =>
        throw new NotImplementedException();
}
