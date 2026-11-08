namespace CheckoutSystem;

/// <summary>
/// One product in the catalog. All money is whole cents: 899 means $8.99.
/// </summary>
/// <param name="Sku">Unique product code, e.g. "COFFEE". Matched exactly (case sensitive).</param>
/// <param name="Name">Human readable name, e.g. "Bag of Coffee Beans".</param>
/// <param name="UnitPrice">Price of one item, in cents.</param>
/// <param name="Promotion">The product's promotion, or <c>null</c> when it has none.</param>
public sealed record Product(string Sku, string Name, int UnitPrice, Promotion? Promotion = null);
