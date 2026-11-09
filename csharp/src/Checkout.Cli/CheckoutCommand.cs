using System.Text.Json;

namespace CheckoutSystem.Cli;

/// <summary>What the program prints and the exit code it ends with.</summary>
public sealed record CommandResult(string Output, int ExitCode);

/// <summary>
/// The command line program's logic (rules §8), kept free of Console.In/Out so it can be
/// unit tested in-process. Program.cs is the thin wrapper that talks to the terminal.
/// </summary>
public static class CheckoutCommand
{
    /// <summary>Where the standard catalog lives: spec/catalog.json at the root of the repository.</summary>
    public static string StandardCatalogPath =>
        Path.Combine(SpecFolder.FindForProgram(), "catalog.json");

    /// <summary>Reads spec/catalog.json.</summary>
    public static Catalog LoadStandardCatalog() =>
        JsonSerializer.Deserialize<Catalog>(File.ReadAllText(StandardCatalogPath), CheckoutJson.Options)!;

    /// <summary>
    /// Runs one checkout from the JSON text of a basket.
    /// <list type="bullet">
    /// <item>a valid basket gives the receipt as JSON and exit code 0,</item>
    /// <item>an invalid basket gives <c>{"error":"&lt;CODE&gt;"}</c> and exit code 1,</item>
    /// <item>text that is not a basket at all throws a <see cref="JsonException"/> (the program turns that into exit code 2).</item>
    /// </list>
    /// </summary>
    public static CommandResult Run(string input, Catalog catalog)
    {
        var basket = ParseBasket(input);
        try
        {
            return new CommandResult(JsonSerializer.Serialize(Checkout.Calculate(basket, catalog), CheckoutJson.Options), 0);
        }
        catch (CheckoutException error)
        {
            return new CommandResult(JsonSerializer.Serialize(new CliError(error.Code.ToSpecCode()), CheckoutJson.Options), 1);
        }
    }

    /// <summary>
    /// Reads the JSON text of a basket. System.Text.Json does the shape checks for us: "items"
    /// must be there, every item needs a string "sku" and a whole-number "quantity". Only the
    /// shape is checked here; whether the values make sense is the checkout's job.
    /// </summary>
    /// <exception cref="JsonException">The text is not JSON, or the JSON is not a basket.</exception>
    private static Basket ParseBasket(string input)
    {
        Basket? basket;
        try
        {
            basket = JsonSerializer.Deserialize<Basket>(input, CheckoutJson.Options);
        }
        catch (JsonException error)
        {
            throw new JsonException($"The input is not a basket: {error.Message}", error);
        }

        // `null` is valid JSON, and so is a list with a `null` item in it - but neither is a basket.
        if (basket is null || basket.Items.Contains(null!))
        {
            throw new JsonException("The input is not a basket: expected an object with a list of items");
        }

        return basket;
    }

    /// <summary>What the program prints on failure: <c>{"error":"EMPTY_BASKET"}</c> (spec/schema/cli-error.schema.json).</summary>
    private sealed record CliError(string Error);
}
