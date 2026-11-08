namespace CheckoutSystem;

/// <summary>
/// The checkout input (rules §3).
/// </summary>
/// <param name="Items">The lines of the basket, in the order the customer added them.</param>
/// <param name="Coupon">A coupon code, or <c>null</c> for no coupon.</param>
/// <param name="Shipping">"standard" or "express". Left out means "standard".</param>
public sealed record Basket(
    IReadOnlyList<BasketItem> Items,
    string? Coupon = null,
    string Shipping = "standard");

/// <summary>
/// One line of a basket. The quantity is NOT checked here: the checkout validates it,
/// so tests can create invalid baskets (0, -1, 100) on purpose.
/// </summary>
public sealed record BasketItem(string Sku, int Quantity);
