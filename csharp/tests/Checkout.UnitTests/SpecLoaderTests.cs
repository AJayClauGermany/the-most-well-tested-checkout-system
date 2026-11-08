using CheckoutSystem.UnitTests.SpecSupport;

namespace CheckoutSystem.UnitTests;

/// <summary>
/// Tests for the test helper <see cref="SpecLoader"/>. Testing test code may sound odd,
/// but if the loader silently read zero cases, every spec-case test would "pass"
/// without checking anything. These tests make sure the shared spec is really loaded.
/// </summary>
public class SpecLoaderTests
{
    [Fact]
    public void LoadCatalog_StandardCatalog_HasSixProducts()
    {
        var catalog = SpecLoader.LoadCatalog();

        Assert.Equal(6, catalog.Products.Count);
    }

    [Fact]
    public void LoadCatalog_Coffee_HasBuyTwoGetOneFreePromotion()
    {
        var coffee = SpecLoader.LoadCatalog().FindBySku("COFFEE");

        Assert.NotNull(coffee);
        // Records compare by value, so we can compare the whole promotion in one go.
        Assert.Equal(new BuyXGetYFree(Buy: 2, Free: 1), coffee.Promotion);
    }

    [Fact]
    public void LoadCatalog_Notebook_HasBulkDiscountPromotion()
    {
        var notebook = SpecLoader.LoadCatalog().FindBySku("NOTEBOOK");

        Assert.NotNull(notebook);
        Assert.Equal(new BulkDiscount(MinQuantity: 10, Percent: 10), notebook.Promotion);
    }

    [Fact]
    public void LoadCases_AllFiles_LoadsSixtyCases()
    {
        var cases = SpecLoader.LoadCases();

        // Update this number when you add cases to spec/test-cases/.
        Assert.Equal(60, cases.Count);
    }

    [Fact]
    public void LoadCases_AllFiles_IdsAreUnique()
    {
        var ids = SpecLoader.LoadCases().Select(specCase => specCase.Id).ToList();

        Assert.Equal(ids.Count, ids.Distinct().Count());
    }

    [Fact]
    public void LoadCases_EveryCase_HasExactlyOneOfExpectedOrExpectedError()
    {
        // Assert.All runs the check for every item and reports every item that fails.
        Assert.All(SpecLoader.LoadCases(), specCase =>
        {
            var hasExpected = specCase.Expected is not null;
            var hasExpectedError = specCase.ExpectedError is not null;
            Assert.True(hasExpected != hasExpectedError, $"{specCase.Id} must have exactly one of expected / expected_error");
        });
    }

    [Fact]
    public void LoadCases_BasketWithoutShipping_DefaultsToStandard()
    {
        var specCase = SpecLoader.LoadCase("subtotal-single-item");

        Assert.Equal("standard", specCase.Basket.Shipping);
        Assert.Null(specCase.Basket.Coupon);
    }
}
