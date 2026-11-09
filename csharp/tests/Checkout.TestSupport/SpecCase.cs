namespace CheckoutSystem.TestSupport;

/// <summary>
/// One shared test case from spec/test-cases/*.json: a basket plus EITHER the
/// expected totals OR the expected error code (never both).
/// </summary>
/// <param name="File">The JSON file the case came from, e.g. "01-subtotal.json".</param>
/// <param name="Id">Unique, readable name of the case, e.g. "subtotal-single-item".</param>
/// <param name="Description">Why the answer is what it is.</param>
/// <param name="Basket">The checkout input.</param>
/// <param name="Expected">The six expected totals, or <c>null</c> when an error is expected.</param>
/// <param name="ExpectedError">The expected error, or <c>null</c> when checkout should succeed.</param>
public sealed record SpecCase(
    string File,
    string Id,
    string Description,
    Basket Basket,
    Totals? Expected,
    ErrorCode? ExpectedError)
{
    /// <summary>Test runners and assertion messages show the id, which is easier to read than every field.</summary>
    public override string ToString() => Id;
}

/// <summary>
/// The six money totals of a receipt (without the lines). Comparing two
/// <see cref="Totals"/> records compares all six numbers at once, and a failing
/// assertion prints both records so you can see which number is wrong.
/// </summary>
public sealed record Totals(
    int Subtotal,
    int PromotionDiscount,
    int CouponDiscount,
    int Shipping,
    int Tax,
    int Total)
{
    public static Totals From(Receipt receipt) => new(
        receipt.Subtotal,
        receipt.PromotionDiscount,
        receipt.CouponDiscount,
        receipt.Shipping,
        receipt.Tax,
        receipt.Total);
}
