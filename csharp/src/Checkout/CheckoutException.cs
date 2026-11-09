namespace CheckoutSystem;

/// <summary>
/// Thrown by <see cref="Checkout.Calculate"/> when a basket cannot be checked out.
/// <see cref="Code"/> says why, using one of the six codes from the spec; the message
/// adds the detail a human wants to see ("Unknown product: BANANA").
///
/// In C# the natural way to report a failure is to throw an exception. Callers tell the
/// reasons apart by looking at <c>exception.Code</c>, never by parsing the message.
/// </summary>
public sealed class CheckoutException : Exception
{
    /// <summary>Fails with a code and a generic message, e.g. "Checkout failed: EMPTY_BASKET".</summary>
    public CheckoutException(ErrorCode code)
        : this(code, $"Checkout failed: {code.ToSpecCode()}")
    {
    }

    /// <summary>Fails with a code and a message that explains what exactly went wrong.</summary>
    public CheckoutException(ErrorCode code, string message)
        : base(message)
    {
        Code = code;
    }

    /// <summary>Why the checkout failed.</summary>
    public ErrorCode Code { get; }
}
