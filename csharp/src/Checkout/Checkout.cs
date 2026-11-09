namespace CheckoutSystem;

/// <summary>
/// The entry point of the checkout system.
/// </summary>
public static class Checkout
{
    /// <summary>
    /// Prices a basket using the rules in spec/checkout-rules.md.
    /// </summary>
    /// <remarks>
    /// It reads like the hint in rules §4.7: validate, merge lines, price, coupon, shipping,
    /// tax, total. Each step is a small method with its own unit tests.
    /// </remarks>
    /// <param name="basket">What the customer wants to buy.</param>
    /// <param name="catalog">The products the shop sells (prices and promotions).</param>
    /// <returns>The receipt with every line and the six totals.</returns>
    /// <exception cref="CheckoutException">The basket is invalid (see rules §4.1).</exception>
    public static Receipt Calculate(Basket basket, Catalog catalog)
    {
        BasketValidator.Validate(basket, catalog); // rules §4.1 rows 1-4
        var method = Shipping.Parse(basket.Shipping); // row 5

        // After validation every merged quantity is from 1 to 99, so it fits in an int.
        var lines = ReceiptLines.Merge(basket.Items)
            .Select(item => ReceiptLines.Price(catalog.RequireProduct(item.Sku), (int)item.Quantity))
            .ToList();
        var subtotal = lines.Sum(line => line.LineTotal);
        var promotionDiscount = lines.Sum(line => line.LineDiscount);
        var afterPromotions = subtotal - promotionDiscount;

        var couponDiscount = Coupons.Discount(basket.Coupon, afterPromotions); // rows 6-7
        var discountedSubtotal = afterPromotions - couponDiscount;

        var shipping = Shipping.Cost(method, discountedSubtotal);
        var tax = Tax.On(discountedSubtotal);

        return new Receipt(
            lines,
            subtotal,
            promotionDiscount,
            couponDiscount,
            shipping,
            tax,
            Total: discountedSubtotal + shipping + tax);
    }

    /// <summary>
    /// Checks out, then saves the receipt as an order (rules §6). A failed checkout saves
    /// nothing: the <see cref="CheckoutException"/> is thrown before the repository is touched.
    /// </summary>
    public static PlacedOrder PlaceOrder(Basket basket, Catalog catalog, IOrderRepository repository)
    {
        var receipt = Calculate(basket, catalog);
        return new PlacedOrder(receipt, repository.Save(receipt));
    }
}

/// <summary>What <see cref="Checkout.PlaceOrder"/> gives back: the receipt and the id the repository gave it.</summary>
public sealed record PlacedOrder(Receipt Receipt, string OrderId);
