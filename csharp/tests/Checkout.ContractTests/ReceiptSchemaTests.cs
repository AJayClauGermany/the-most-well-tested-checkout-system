using System.Text.Json;
using System.Text.Json.Nodes;

using CheckoutSystem.Cli;

namespace CheckoutSystem.ContractTests;

/// <summary>
/// CONTRACT TESTS: is the output even the right shape?
///
/// spec/schema/receipt.schema.json is the contract every language's receipt must honour.
/// Here JsonSchema.Net checks real receipts against it. The shared e2e runner only compares
/// totals; this is where the shape of the whole receipt (lines included) is held to the contract.
/// </summary>
public class ReceiptSchemaTests
{
    private static readonly Catalog Catalog = SpecLoader.LoadCatalog();

    public static TheoryData<string> ValidCaseIds() =>
        new(SpecLoader.LoadValidCases().Select(specCase => specCase.Id));

    [Theory]
    [MemberData(nameof(ValidCaseIds))]
    public void Calculate_EveryValidSpecCase_GivesAReceiptThatFollowsTheSchema(string caseId)
    {
        var receipt = Checkout.Calculate(SpecLoader.LoadCase(caseId).Basket, Catalog);

        Assert.Empty(SchemaValidator.ErrorsFor("receipt.schema.json", receipt));
    }

    [Fact]
    public void Calculate_MadeUpCatalog_FollowsTheSchemaToo()
    {
        var receipt = Checkout.Calculate(new Basket([new BasketItem("TV", 1), new BasketItem("COFFEE", 6)]), Fixtures.TinyCatalog);

        Assert.Empty(SchemaValidator.ErrorsFor("receipt.schema.json", receipt));
    }

    [Fact]
    public void Serialize_Receipt_UsesTheSnakeCaseFieldNamesFromTheSpec()
    {
        var json = JsonSerializer.SerializeToNode(Fixtures.OneAppleReceipt, CheckoutJson.Options)!.AsObject();

        Assert.Equal(["lines", "subtotal", "promotion_discount", "coupon_discount", "shipping", "tax", "total"], json.Select(property => property.Key));
        Assert.Equal(["sku", "name", "quantity", "unit_price", "line_total", "line_discount"], json["lines"]![0]!.AsObject().Select(property => property.Key));
    }

    [Fact]
    public void SchemaValidator_BrokenReceipt_IsRejected_TestingTheTest()
    {
        // A schema check that accepted anything would make every contract test meaningless.
        var broken = JsonSerializer.SerializeToNode(Fixtures.OneAppleReceipt, CheckoutJson.Options)!.AsObject();
        broken["shipping"] = 123;
        broken["tax"] = -1;
        broken["extra"] = true;

        var errors = SchemaValidator.Errors("receipt.schema.json", broken);

        Assert.Equal(3, errors.Count);
        Assert.Contains(errors, error => error.StartsWith("/shipping ", StringComparison.Ordinal));
        Assert.Contains(errors, error => error.StartsWith("/tax ", StringComparison.Ordinal));
        Assert.Contains(errors, error => error.Contains("additionalProperties", StringComparison.Ordinal));
    }

    [Fact]
    public void SchemaValidator_ExampleReceiptFromTheSpec_IsValid()
    {
        Assert.Empty(SchemaValidator.Errors("receipt.schema.json", SpecLoader.ReadSpecFile("examples", "receipt.json")));
        Assert.Empty(SchemaValidator.Errors("order.schema.json", SpecLoader.ReadSpecFile("examples", "order.json")));
        Assert.Empty(SchemaValidator.Errors("cli-error.schema.json", SpecLoader.ReadSpecFile("examples", "cli-error.json")));
    }
}

/// <summary>The CLI's error output must follow spec/schema/cli-error.schema.json (rules §8).</summary>
public class CliErrorSchemaTests
{
    private static readonly Catalog Catalog = SpecLoader.LoadCatalog();

    public static TheoryData<string> InvalidCaseIds() =>
        new(SpecLoader.LoadInvalidCases().Select(specCase => specCase.Id));

    [Theory]
    [MemberData(nameof(InvalidCaseIds))]
    public void Run_EveryInvalidSpecCase_PrintsAnErrorThatFollowsTheSchema(string caseId)
    {
        var basket = JsonSerializer.Serialize(SpecLoader.LoadCase(caseId).Basket, CheckoutJson.Options);

        var result = CheckoutCommand.Run(basket, Catalog);

        Assert.Equal(1, result.ExitCode);
        Assert.Empty(SchemaValidator.Errors("cli-error.schema.json", result.Output));
    }

    [Fact]
    public void SchemaValidator_UnknownErrorCode_IsRejected_TestingTheTest()
    {
        Assert.NotEmpty(SchemaValidator.Errors("cli-error.schema.json", """{"error":"SOMETHING_ELSE"}"""));
        Assert.NotEmpty(SchemaValidator.Errors("cli-error.schema.json", """{"error":"EMPTY_BASKET","hint":"?"}"""));
    }
}

/// <summary>The spec must be consistent with itself: every shared case's basket is a valid basket.</summary>
public class BasketSchemaTests
{
    [Fact]
    public void EverySpecCaseBasket_FollowsBasketSchemaJson()
    {
        Assert.All(SpecLoader.LoadCases(), specCase =>
            Assert.Empty(SchemaValidator.ErrorsFor("basket.schema.json", specCase.Basket)));
    }

    [Fact]
    public void Serialize_Basket_WritesTheOptionalFieldsTheWayTheSchemaAllows()
    {
        var json = JsonSerializer.SerializeToNode(new Basket([new BasketItem("APPLE", 1)]), CheckoutJson.Options);

        // "coupon" is written as JSON null (System.Text.Json gives a C# null for it), which the schema allows.
        Assert.True(json!.AsObject().ContainsKey("coupon"));
        Assert.Null(json["coupon"]);
        Assert.Equal("standard", json["shipping"]!.GetValue<string>());
    }
}
