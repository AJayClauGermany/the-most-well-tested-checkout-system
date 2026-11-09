using System.Text.Json;

namespace CheckoutSystem;

/// <summary>
/// How our records map to the JSON in spec/: the spec uses snake_case (<c>unit_price</c>),
/// C# uses PascalCase (<c>UnitPrice</c>). Everything that reads or writes checkout JSON
/// (the CLI, the file and SQLite stores, the tests) uses these options, so the field names
/// are right everywhere – and defined only once.
/// </summary>
public static class CheckoutJson
{
    public static JsonSerializerOptions Options { get; } = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        // Fail loudly when a required field is missing or null, instead of quietly using null.
        // (Fields with a default value, like Basket.Shipping = "standard", may be left out.)
        RespectNullableAnnotations = true,
        RespectRequiredConstructorParameters = true,
        // Accept the promotion "type" field in any position, not only first.
        AllowOutOfOrderMetadataProperties = true,
    };

    /// <summary>The same options, but written out indented and human readable (for files on disk).</summary>
    public static JsonSerializerOptions IndentedOptions { get; } = new(Options) { WriteIndented = true };
}
