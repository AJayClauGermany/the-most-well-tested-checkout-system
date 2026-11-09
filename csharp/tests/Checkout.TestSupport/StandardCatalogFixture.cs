namespace CheckoutSystem.TestSupport;

/// <summary>
/// The standard catalog and the 60 shared cases, loaded ONCE and shared by many tests.
///
/// xUnit creates a new instance of a test class for every test, so a field in the class
/// would re-read the JSON files for each of the 60 cases. A "fixture" is an object xUnit
/// creates once instead – per class with <c>IClassFixture&lt;T&gt;</c>, or once for a whole
/// group of classes with <c>ICollectionFixture&lt;T&gt;</c> and <c>[Collection("…")]</c>.
/// Tests in the same collection also run one after another, never in parallel; see
/// tests/README.md.
/// </summary>
public sealed class StandardCatalogFixture
{
    public Catalog Catalog { get; } = SpecLoader.LoadCatalog();

    public IReadOnlyList<SpecCase> Cases { get; } = SpecLoader.LoadCases();
}
