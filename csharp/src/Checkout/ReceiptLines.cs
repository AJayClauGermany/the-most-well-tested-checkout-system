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
    public static IReadOnlyList<BasketItem> Merge(IReadOnlyList<BasketItem> items)
    {
        // The list keeps the order of first appearance; the dictionary remembers where each SKU is.
        var merged = new List<BasketItem>();
        var positions = new Dictionary<string, int>(StringComparer.Ordinal);
        foreach (var item in items)
        {
            if (positions.TryGetValue(item.Sku, out var position))
            {
                merged[position] = merged[position] with { Quantity = merged[position].Quantity + item.Quantity };
            }
            else
            {
                positions[item.Sku] = merged.Count;
                merged.Add(item);
            }
        }

        return merged;
    }

    /// <summary>The receipt line for <paramref name="quantity"/> of <paramref name="product"/>, including its promotion discount.</summary>
    public static ReceiptLine Price(Product product, int quantity) => new(
        product.Sku,
        product.Name,
        quantity,
        product.UnitPrice,
        LineTotal: product.UnitPrice * quantity,
        LineDiscount: Promotions.Discount(product.Promotion, quantity, product.UnitPrice));
}
