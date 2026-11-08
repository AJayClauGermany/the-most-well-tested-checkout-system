using System.Text.Json;

namespace CheckoutSystem.UnitTests.SpecSupport;

/// <summary>
/// Reads the shared, language-agnostic spec files from the repository's <c>spec/</c> folder.
/// The files are read where they are (never copied), so every language is tested
/// against exactly the same data.
/// </summary>
public static class SpecLoader
{
    /// <summary>
    /// How JSON maps to our records: the spec uses snake_case (<c>unit_price</c>),
    /// C# uses PascalCase (<c>UnitPrice</c>).
    /// </summary>
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        // Fail loudly when a required field is missing or null, instead of quietly using null.
        // (Fields with a default value, like Basket.Shipping = "standard", may be left out.)
        RespectNullableAnnotations = true,
        RespectRequiredConstructorParameters = true,
        // Accept the promotion "type" field in any position, not only first.
        AllowOutOfOrderMetadataProperties = true,
    };

    /// <summary>
    /// The repository's <c>spec/</c> folder. Test runners start in different working
    /// directories, so we start from the folder of the compiled test program
    /// (e.g. <c>csharp/tests/Checkout.UnitTests/bin/Debug/net10.0/</c>) and walk up
    /// until we find a folder that contains <c>spec/catalog.json</c>.
    /// </summary>
    public static DirectoryInfo SpecDirectory { get; } = FindSpecDirectory();

    /// <summary>Loads spec/catalog.json, the standard catalog used by every shared test case.</summary>
    public static Catalog LoadCatalog() =>
        ReadJson<Catalog>(Path.Combine(SpecDirectory.FullName, "catalog.json"));

    /// <summary>Loads every case from every spec/test-cases/*.json file, sorted by file name.</summary>
    public static IReadOnlyList<SpecCase> LoadCases()
    {
        var testCasesDirectory = Path.Combine(SpecDirectory.FullName, "test-cases");
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

    private static T ReadJson<T>(string path)
    {
        using var stream = File.OpenRead(path);
        return JsonSerializer.Deserialize<T>(stream, JsonOptions)
            ?? throw new InvalidDataException($"{path} contains only 'null'");
    }

    private static DirectoryInfo FindSpecDirectory()
    {
        for (var directory = new DirectoryInfo(AppContext.BaseDirectory); directory is not null; directory = directory.Parent)
        {
            var candidate = new DirectoryInfo(Path.Combine(directory.FullName, "spec"));
            if (File.Exists(Path.Combine(candidate.FullName, "catalog.json")))
            {
                return candidate;
            }
        }

        throw new DirectoryNotFoundException(
            $"Could not find spec/catalog.json in {AppContext.BaseDirectory} or any of its parent folders.");
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
