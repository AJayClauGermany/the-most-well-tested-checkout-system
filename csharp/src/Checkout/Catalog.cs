namespace CheckoutSystem;

/// <summary>
/// The list of products the shop sells (rules §2). The checkout receives the catalog
/// as an input, so tests can build a tiny made-up catalog instead of the standard one.
/// </summary>
public sealed record Catalog(IReadOnlyList<Product> Products)
{
    /// <summary>
    /// Finds a product by its SKU. SKUs are case sensitive: "apple" does not find "APPLE".
    /// </summary>
    /// <returns>The product, or <c>null</c> when no product has this SKU.</returns>
    public Product? FindBySku(string sku) =>
        Products.FirstOrDefault(product => product.Sku == sku);
}
