namespace CheckoutSystem;

/// <summary>
/// Why a checkout failed (rules §9). C# enum members use PascalCase;
/// <see cref="ErrorCodeExtensions.ToSpecCode"/> gives the exact spelling used by the spec.
/// </summary>
public enum ErrorCode
{
    /// <summary>The basket has no items.</summary>
    EmptyBasket,

    /// <summary>An item's SKU is not in the catalog.</summary>
    UnknownProduct,

    /// <summary>A quantity is not a whole number from 1 to 99.</summary>
    InvalidQuantity,

    /// <summary>Shipping is not "standard" or "express".</summary>
    UnknownShippingMethod,

    /// <summary>The coupon code does not exist.</summary>
    UnknownCoupon,

    /// <summary>The coupon needs a bigger spend.</summary>
    CouponMinimumNotMet,
}

public static class ErrorCodeExtensions
{
    /// <summary>
    /// The code exactly as the spec writes it, e.g. <c>ErrorCode.EmptyBasket</c> → <c>"EMPTY_BASKET"</c>.
    /// </summary>
    public static string ToSpecCode(this ErrorCode code) => code switch
    {
        ErrorCode.EmptyBasket => "EMPTY_BASKET",
        ErrorCode.UnknownProduct => "UNKNOWN_PRODUCT",
        ErrorCode.InvalidQuantity => "INVALID_QUANTITY",
        ErrorCode.UnknownShippingMethod => "UNKNOWN_SHIPPING_METHOD",
        ErrorCode.UnknownCoupon => "UNKNOWN_COUPON",
        ErrorCode.CouponMinimumNotMet => "COUPON_MINIMUM_NOT_MET",
        _ => throw new ArgumentOutOfRangeException(nameof(code), code, "Unknown error code"),
    };

    /// <summary>
    /// The opposite of <see cref="ToSpecCode"/>: <c>"EMPTY_BASKET"</c> → <c>ErrorCode.EmptyBasket</c>.
    /// </summary>
    /// <exception cref="ArgumentException">The text is not one of the spec's error codes.</exception>
    public static ErrorCode FromSpecCode(string specCode)
    {
        foreach (var code in Enum.GetValues<ErrorCode>())
        {
            if (code.ToSpecCode() == specCode)
            {
                return code;
            }
        }

        throw new ArgumentException($"'{specCode}' is not a known error code", nameof(specCode));
    }
}
