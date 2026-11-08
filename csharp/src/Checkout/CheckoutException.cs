namespace CheckoutSystem;

/// <summary>
/// Thrown by <see cref="Checkout.Calculate"/> when a basket cannot be checked out.
/// <see cref="Code"/> says why, using one of the six codes from the spec.
/// </summary>
public sealed class CheckoutException : Exception
{
    public CheckoutException(ErrorCode code)
        : base($"Checkout failed: {code.ToSpecCode()}")
    {
        Code = code;
    }

    /// <summary>Why the checkout failed.</summary>
    public ErrorCode Code { get; }
}
