using System.Text.Json;

using CheckoutSystem.Cli;

namespace CheckoutSystem.TestSupport;

/// <summary>
/// Reads the shared, language-agnostic spec files from the repository's <c>spec/</c> folder.
/// The files are read where they are (never copied), so every language is tested
/// against exactly the same data.
/// </summary>
public static class SpecLoader
{
    /// <summary>
    /// The repository's <c>spec/</c> folder. Test runners start in different working
    /// directories, so we start from the folder of the compiled test program
    /// (e.g. <c>csharp/tests/Checkout.UnitTests/bin/Debug/net10.0/</c>) and walk up
    /// until we find a folder that contains <c>spec/catalog.json</c>. The search itself is
    /// shared with the command line program, see <see cref="SpecFolder"/>.
    /// </summary>
    public static string SpecDirectory { get; } = SpecFolder.Find(AppContext.BaseDirectory);

    /// <summary>Loads spec/catalog.json, the standard catalog used by every shared test case.</summary>
    public static Catalog LoadCatalog() =>
        ReadJson<Catalog>(Path.Combine(SpecDirectory, "catalog.json"));

    /// <summary>Loads every case from every spec/test-cases/*.json file, sorted by file name.</summary>
    public static IReadOnlyList<SpecCase> LoadCases()
    {
        var testCasesDirectory = Path.Combine(SpecDirectory, "test-cases");
        var files = Directory.GetFiles(testCasesDirectory, "*.json").Order(StringComparer.Ordinal);

        var cases = new List<SpecCase>();
        foreach (var path in files)
        {
            var file = ReadJson<CaseFileJson>(path);
            foreach (var json in file.Cases)
            {
                cases.Add(new SpecCase(
                    File: Path.GetFileName(path),
                    Id: json.Id,
                    Description: json.Description,
                    Basket: json.Basket,
                    Expected: json.Expected,
                    ExpectedError: json.ExpectedError is null ? null : ErrorCodeExtensions.FromSpecCode(json.ExpectedError)));
            }
        }

        return cases;
    }

    /// <summary>Finds a single case by its id, e.g. "subtotal-single-item".</summary>
    public static SpecCase LoadCase(string id) => LoadCases().Single(specCase => specCase.Id == id);

    /// <summary>The cases that expect a receipt (the valid baskets).</summary>
    public static IReadOnlyList<SpecCase> LoadValidCases() =>
        LoadCases().Where(specCase => specCase.Expected is not null).ToList();

    /// <summary>The cases that expect an error code (the invalid baskets).</summary>
    public static IReadOnlyList<SpecCase> LoadInvalidCases() =>
        LoadCases().Where(specCase => specCase.ExpectedError is not null).ToList();

    /// <summary>The text of a file in spec/, e.g. <c>ReadSpecFile("examples", "receipt.json")</c>.</summary>
    public static string ReadSpecFile(params string[] pathInsideSpec) =>
        File.ReadAllText(Path.Combine([SpecDirectory, .. pathInsideSpec]));

    private static T ReadJson<T>(string path)
    {
        using var stream = File.OpenRead(path);
        return JsonSerializer.Deserialize<T>(stream, CheckoutJson.Options)
            ?? throw new InvalidDataException($"{path} contains only 'null'");
    }

    // The shape of a spec/test-cases/*.json file. These records exist only to read the
    // JSON; SpecCase above is the friendlier type the tests use.
    private sealed record CaseFileJson(string Description, IReadOnlyList<CaseJson> Cases);

    private sealed record CaseJson(
        string Id,
        string Description,
        Basket Basket,
        Totals? Expected = null,
        string? ExpectedError = null);
}
