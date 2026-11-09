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
    public static int RoundHalfUp(int amount, int percent) =>
        throw new NotImplementedException();
}
