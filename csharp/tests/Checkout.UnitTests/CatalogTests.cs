
namespace CheckoutSystem.UnitTests;

/// <summary>
/// EXAMPLE UNIT TESTS – start reading here!
///
/// A unit test checks one small piece of code (a "unit") on its own, without files,
/// networks or databases. These tests check <see cref="Catalog.FindBySku"/>.
///
/// xUnit basics:
/// * A test is a public method marked with [Fact]. xUnit finds and runs it for you.
/// * The class does not need any attribute. xUnit creates a NEW instance of the class
///   for every test, so tests cannot accidentally share state.
/// * Test names follow the common .NET pattern Method_Scenario_ExpectedResult, so a
///   failing test already tells you what broke.
/// * Every test has three parts: Arrange (set up), Act (do the thing), Assert (check).
/// </summary>
public class CatalogTests
{
    // A tiny, made-up catalog. Unit tests do not need the real spec/catalog.json:
    // a small catalog written right here makes the test easy to read on its own.
    private readonly Catalog _catalog = new(
    [
        new Product("APPLE", "Apple", UnitPrice: 50),
        new Product("COFFEE", "Bag of Coffee Beans", UnitPrice: 899, new BuyXGetYFree(Buy: 2, Free: 1)),
    ]);

    [Fact]
    public void FindBySku_KnownSku_ReturnsTheProduct()
    {
        // Arrange: everything the test needs (here: the catalog above and a SKU).
        var sku = "COFFEE";

        // Act: call the code we are testing – exactly once.
        var product = _catalog.FindBySku(sku);

        // Assert: check the result. Assert.NotNull fails the test when product is null.
        Assert.NotNull(product);
        Assert.Equal("Bag of Coffee Beans", product.Name);
        Assert.Equal(899, product.UnitPrice);
    }

    [Fact]
    public void FindBySku_UnknownSku_ReturnsNull()
    {
        var product = _catalog.FindBySku("BANANA");

        Assert.Null(product);
    }

    [Fact]
    public void FindBySku_SkuInWrongCase_ReturnsNull()
    {
        // Rules §4 step 1: SKUs are case sensitive, so "apple" is NOT "APPLE".
        var product = _catalog.FindBySku("apple");

        Assert.Null(product);
    }
}
