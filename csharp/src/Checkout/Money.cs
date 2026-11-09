namespace CheckoutSystem;

/// <summary>
/// Money helpers (rules §1 and §7). All amounts are whole cents, so everything here is
/// integer maths – no <c>double</c>, no <c>decimal</c>, no surprises.
/// </summary>
public static class Money
{
    /// <summary>
    /// <paramref name="percent"/>% of <paramref name="amount"/>, rounded half up to a whole
    /// cent (rules §7). <c>RoundHalfUp(225, 10)</c> is 23 (22.5 rounds UP);
    /// <c>RoundHalfUp(224, 10)</c> is 22.
    /// </summary>
    /// <remarks>
    /// Done with whole numbers only, as the hint in rules §7 suggests: <c>amount × percent</c>
    /// is the exact answer in hundredths of a cent; adding 50 hundredths and dividing by 100
    /// (C#'s <c>/</c> on whole numbers rounds down for positive numbers) is the same as
    /// rounding half up. The multiplication is done in <c>long</c> so that large amounts
    /// cannot overflow.
    /// <code>
    /// RoundHalfUp(225, 10)  →  (2250 + 50) / 100  →  23      (22.5 rounds UP)
    /// RoundHalfUp(224, 10)  →  (2240 + 50) / 100  →  22.9 → 22
    /// </code>
    /// </remarks>
    public static int RoundHalfUp(int amount, int percent) =>
        (int)((((long)amount * percent) + 50) / 100);
}
