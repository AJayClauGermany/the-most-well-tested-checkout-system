namespace CheckoutSystem;

/// <summary>
/// Turning basket items into receipt lines (rules §4.2): merging duplicate SKUs and
/// pricing one line.
/// </summary>
public static class ReceiptLines
{
    /// <summary>
    /// One item per distinct SKU, in order of first appearance, quantities added up.
    /// <c>[APPLE x1, PEN x2, APPLE x3]</c> gives <c>[APPLE x4, PEN x2]</c>.
    /// </summary>
    public static IReadOnlyList<BasketItem> Merge(IReadOnlyList<BasketItem> items) =>
        throw new NotImplementedException();

    /// <summary>The receipt line for <paramref name="quantity"/> of <paramref name="product"/>, including its promotion discount.</summary>
    public static ReceiptLine Price(Product product, int quantity) =>
        throw new NotImplementedException();
}
