namespace CheckoutSystem;

/// <summary>The two kinds of promotion (rules §4.2).</summary>
public static class Promotions
{
    /// <summary>
    /// The discount a promotion gives one receipt line, in cents. <c>0</c> without a promotion.
    /// </summary>
    public static int Discount(Promotion? promotion, int quantity, int unitPrice) =>
        throw new NotImplementedException();
}
