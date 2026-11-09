namespace CheckoutSystem;

/// <summary>
/// The entry point of the checkout system.
/// </summary>
public static class Checkout
{
    /// <summary>
    /// Prices a basket using the rules in spec/checkout-rules.md.
    /// </summary>
    /// <param name="basket">What the customer wants to buy.</param>
    /// <param name="catalog">The products the shop sells (prices and promotions).</param>
    /// <returns>The receipt with every line and the six totals.</returns>
    /// <exception cref="CheckoutException">The basket is invalid (see rules §4.1).</exception>
    public static Receipt Calculate(Basket basket, Catalog catalog) =>
        throw new NotImplementedException();

    /// <summary>
    /// Checks out, then saves the receipt as an order (rules §6). A failed checkout saves
    /// nothing: the <see cref="CheckoutException"/> is thrown before the repository is touched.
    /// </summary>
    public static PlacedOrder PlaceOrder(Basket basket, Catalog catalog, IOrderRepository repository) =>
        throw new NotImplementedException();
}

/// <summary>What <see cref="Checkout.PlaceOrder"/> gives back: the receipt and the id the repository gave it.</summary>
public sealed record PlacedOrder(Receipt Receipt, string OrderId);
