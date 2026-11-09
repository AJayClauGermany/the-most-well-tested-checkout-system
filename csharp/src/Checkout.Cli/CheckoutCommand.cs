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
        throw new NotImplementedException();

    /// <summary>Reads spec/catalog.json.</summary>
    public static Catalog LoadStandardCatalog() =>
        throw new NotImplementedException();

    /// <summary>
    /// Runs one checkout from the JSON text of a basket.
    /// <list type="bullet">
    /// <item>a valid basket gives the receipt as JSON and exit code 0,</item>
    /// <item>an invalid basket gives <c>{"error":"&lt;CODE&gt;"}</c> and exit code 1,</item>
    /// <item>text that is not a basket at all throws a <see cref="JsonException"/> (the program turns that into exit code 2).</item>
    /// </list>
    /// </summary>
    public static CommandResult Run(string input, Catalog catalog) =>
        throw new NotImplementedException();
}
