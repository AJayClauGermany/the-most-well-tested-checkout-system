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
    /// <exception cref="CheckoutException">The basket is invalid (see rules §4, step 1).</exception>
    public static Receipt Calculate(Basket basket, Catalog catalog)
    {
        // Test-driven development: the tests come first. Remove this line and
        // implement the rules step by step, making one skipped test pass at a time.
        throw new NotImplementedException("Checkout.Calculate is not implemented yet – see spec/checkout-rules.md");
    }
}
