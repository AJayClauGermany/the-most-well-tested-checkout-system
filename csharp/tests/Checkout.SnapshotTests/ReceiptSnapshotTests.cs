using System.Text.Json;
using System.Text.Json.Nodes;

using CheckoutSystem.Cli;

namespace CheckoutSystem.SnapshotTests;

/// <summary>
/// SNAPSHOT (approval) TESTS: did the output change since last time?
///
/// <c>Verify(...)</c> writes what it received to <c>Snapshots/&lt;Class&gt;.&lt;Test&gt;.received.txt</c>
/// and compares it with the committed <c>.verified.txt</c> next to it. No verified file yet,
/// or a difference: the test fails and shows the two side by side. When the new output is
/// right, rename (or copy) the received file over the verified one and commit – that is the
/// "approval".
///
/// The hand-checked receipt in spec/examples/receipt.json plays the role of an approved
/// output shared by all six languages, so the first test compares with that directly.
/// </summary>
public class ReceiptSnapshotTests
{
    private static readonly Catalog Catalog = SpecLoader.LoadCatalog();

    /// <summary>The biggest shared case, "full-everything": every product, both promotions, WELCOME5, express.</summary>
    private static readonly Basket EverythingBasket = new(
        [
            new BasketItem("HEADPHONES", 1),
            new BasketItem("COFFEE", 6),
            new BasketItem("NOTEBOOK", 12),
            new BasketItem("PEN", 5),
            new BasketItem("BREAD", 1),
            new BasketItem("APPLE", 3),
        ],
        Coupon: "WELCOME5",
        Shipping: "express");

    [Fact]
    public void Calculate_ExampleBasket_MatchesTheHandCheckedExampleInTheSpec_LinesIncluded()
    {
        var approved = JsonNode.Parse(SpecLoader.ReadSpecFile("examples", "receipt.json"));

        var receipt = JsonSerializer.SerializeToNode(Checkout.Calculate(Fixtures.ExampleBasket, Catalog), CheckoutJson.Options);

        Assert.True(JsonNode.DeepEquals(approved, receipt), receipt?.ToJsonString());
    }

    [Fact]
    public Task Calculate_ExampleBasket_MatchesTheSnapshot()
    {
        var receipt = Checkout.Calculate(Fixtures.ExampleBasket, Catalog);

        // Verify serializes the record itself; the snapshot shows every field, PascalCase.
        return Verifier.Verify(receipt);
    }

    [Fact]
    public Task Run_ExampleBasket_PrintsExactlyThisJson()
    {
        var result = CheckoutCommand.Run(JsonSerializer.Serialize(Fixtures.ExampleBasket, CheckoutJson.Options), Catalog);

        // VerifyJson snapshots the JSON text pretty-printed, so a diff points at one field.
        return Verifier.VerifyJson(result.Output);
    }

    [Fact]
    public Task Calculate_EverythingBasket_MatchesTheSnapshot()
    {
        var receipt = Checkout.Calculate(EverythingBasket, Catalog);

        return Verifier.Verify(receipt);
    }

    [Fact]
    public Task Run_EmptyBasket_PrintsExactlyThisError()
    {
        var result = CheckoutCommand.Run("""{"items":[]}""", Catalog);

        return Verifier.Verify(result);
    }
}
