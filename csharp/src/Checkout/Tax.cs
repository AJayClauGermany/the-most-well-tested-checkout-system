namespace CheckoutSystem;

/// <summary>Tax (rules §4.6).</summary>
public static class Tax
{
    public const int Percent = 8;

    /// <summary>8% of <paramref name="discountedSubtotal"/>, rounded half up. Shipping is never passed in here.</summary>
    public static int On(int discountedSubtotal) =>
        throw new NotImplementedException();
}
