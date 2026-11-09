using System.Text.Json;
using System.Text.Json.Nodes;

using CheckoutSystem.Cli;

namespace CheckoutSystem.UnitTests;

/// <summary>
/// Unit tests for the command line program's logic (rules §8), run in-process.
/// The real program (Program.cs) is tested from the outside in Checkout.EndToEndTests.
/// </summary>
public class CheckoutCommandTests
{
    private static readonly Catalog Catalog = Fixtures.TinyCatalog;

    [Fact]
    public void StandardCatalogPath_PointsAtSpecCatalogJsonAtTheRepositoryRoot()
    {
        var path = CheckoutCommand.StandardCatalogPath;

        Assert.EndsWith(Path.Combine("spec", "catalog.json"), path, StringComparison.Ordinal);
        // The repository root is the folder that also holds the rules.
        var root = Path.GetDirectoryName(Path.GetDirectoryName(path))!;
        Assert.True(File.Exists(Path.Combine(root, "spec", "checkout-rules.md")));
    }

    [Fact]
    public void LoadStandardCatalog_LoadsTheSixStandardProducts()
    {
        var catalog = CheckoutCommand.LoadStandardCatalog();

        Assert.Equal(["APPLE", "BREAD", "COFFEE", "NOTEBOOK", "PEN", "HEADPHONES"], catalog.Products.Select(product => product.Sku));
    }

    [Fact]
    public void Run_ValidBasket_PrintsTheReceiptAsJsonAndExitsWithZero()
    {
        var result = CheckoutCommand.Run("""{"items":[{"sku":"APPLE","quantity":1}]}""", Catalog);

        Assert.Equal(0, result.ExitCode);
        var expected = JsonNode.Parse("""
            {"lines":[{"sku":"APPLE","name":"Apple","quantity":1,"unit_price":50,"line_total":50,"line_discount":0}],
             "subtotal":50,"promotion_discount":0,"coupon_discount":0,"shipping":499,"tax":4,"total":553}
            """);
        Assert.True(JsonNode.DeepEquals(expected, JsonNode.Parse(result.Output)), result.Output);
    }

    [Fact]
    public void Run_InvalidBasket_PrintsTheErrorCodeAndExitsWithOne()
    {
        var result = CheckoutCommand.Run("""{"items":[]}""", Catalog);

        Assert.Equal(new CommandResult("""{"error":"EMPTY_BASKET"}""", 1), result);
    }

    [Fact]
    public void Run_EveryKindOfError_IsReportedTheSameWay()
    {
        var result = CheckoutCommand.Run("""{"items":[{"sku":"APPLE","quantity":1}],"coupon":"NOPE"}""", Catalog);

        Assert.Equal(new CommandResult("""{"error":"UNKNOWN_COUPON"}""", 1), result);
    }

    [Fact]
    public void Run_ValidBasket_PrintsCompactJsonOnASingleLine()
    {
        var result = CheckoutCommand.Run("""{"items":[{"sku":"APPLE","quantity":1}]}""", Catalog);

        Assert.DoesNotContain('\n', result.Output);
        Assert.DoesNotContain(": ", result.Output);
    }

    [Fact]
    public void Run_SurroundingWhitespace_IsAccepted()
    {
        // A shell pipeline adds a newline after the JSON; that must not matter.
        var result = CheckoutCommand.Run("  {\"items\":[{\"sku\":\"APPLE\",\"quantity\":1}]}\n", Catalog);

        Assert.Equal(0, result.ExitCode);
    }

    [Fact]
    public void Run_MissingShippingAndNullCoupon_MeansStandardAndNoCoupon()
    {
        var result = CheckoutCommand.Run("""{"items":[{"sku":"APPLE","quantity":1}],"coupon":null}""", Catalog);

        Assert.Equal(0, result.ExitCode);
        Assert.Equal(499, JsonNode.Parse(result.Output)?["shipping"]?.GetValue<int>());
    }

    [Theory]
    [InlineData("not json")]
    [InlineData("")]
    [InlineData("{\"items\": [")]
    public void Run_TextThatIsNotJson_ThrowsInsteadOfInventingAnErrorCode(string input)
    {
        Assert.Throws<JsonException>(() => CheckoutCommand.Run(input, Catalog));
    }

    [Theory]
    [InlineData("null")]
    [InlineData("42")]
    [InlineData("\"text\"")]
    [InlineData("[]")]
    [InlineData("{}")]
    [InlineData("""{"items": "APPLE"}""")]
    [InlineData("""{"items": [1, 2]}""")]
    [InlineData("""{"items": [null]}""")]
    [InlineData("""{"items": ["APPLE"]}""")]
    [InlineData("""{"items": [{"sku": 5}]}""")]
    [InlineData("""{"items": [{"sku": "APPLE"}]}""")]
    [InlineData("""{"items": [{"sku": "APPLE", "quantity": 1.5}]}""")]
    [InlineData("""{"items": [{"sku": "APPLE", "quantity": "1"}]}""")]
    public void Run_JsonThatIsNotABasket_ThrowsInsteadOfInventingAnErrorCode(string input)
    {
        Assert.Throws<JsonException>(() => CheckoutCommand.Run(input, Catalog));
    }
}
