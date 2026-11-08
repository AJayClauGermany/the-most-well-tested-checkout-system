using System.Text.Json.Serialization;

namespace CheckoutSystem;

/// <summary>
/// A promotion attached to a product (rules §4, step 3). There are exactly two kinds,
/// so <see cref="Promotion"/> is an abstract base record with one sub-record per kind.
/// Later code can tell them apart with a <c>switch</c>:
/// <code>
/// var discount = promotion switch
/// {
///     BuyXGetYFree p => ...,
///     BulkDiscount p => ...,
///     _ => 0,
/// };
/// </code>
/// The attributes tell System.Text.Json to use the JSON field <c>"type"</c> to pick the kind.
/// </summary>
[JsonPolymorphic(TypeDiscriminatorPropertyName = "type")]
[JsonDerivedType(typeof(BuyXGetYFree), "buy_x_get_y_free")]
[JsonDerivedType(typeof(BulkDiscount), "bulk_discount")]
public abstract record Promotion;

/// <summary>"Buy <paramref name="Buy"/>, get <paramref name="Free"/> free", e.g. buy 2 coffees, get 1 free.</summary>
public sealed record BuyXGetYFree(int Buy, int Free) : Promotion;

/// <summary>"<paramref name="Percent"/>% off the line when you buy <paramref name="MinQuantity"/> or more".</summary>
public sealed record BulkDiscount(int MinQuantity, int Percent) : Promotion;
