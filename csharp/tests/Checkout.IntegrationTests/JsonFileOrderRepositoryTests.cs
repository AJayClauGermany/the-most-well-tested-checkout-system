using System.Text.Json;
using System.Text.Json.Nodes;

using CheckoutSystem.Storage;

namespace CheckoutSystem.IntegrationTests;

/// <summary>
/// INTEGRATION TESTS for the JSON-file order repository: real files in a real temporary folder.
///
/// Unlike a unit test, this touches the file system, so each test gets its own fresh folder
/// (<see cref="TemporaryFolder"/>) and removes it afterwards. The stored files must follow
/// spec/schema/order.schema.json – a contract test inside an integration test.
/// </summary>
public sealed class JsonFileOrderRepositoryTests : IDisposable
{
    private static readonly Receipt Receipt = Fixtures.OneAppleReceipt;

    private readonly TemporaryFolder _folder = new();

    public void Dispose() => _folder.Dispose();

    [Fact]
    public void Save_FirstOrder_WritesOrder1JsonIntoTheFolder()
    {
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());

        var id = repository.Save(Receipt);

        Assert.Equal("ORDER-1", id);
        Assert.True(File.Exists(_folder.PathOf("ORDER-1.json")));
    }

    [Fact]
    public void Save_WritesAFileThatFollowsOrderSchemaJsonExactly()
    {
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.FixedAt("2026-01-15T10:30:00Z"));
        repository.Save(Receipt);

        var stored = File.ReadAllText(_folder.PathOf("ORDER-1.json"));

        Assert.Empty(SchemaValidator.Errors("order.schema.json", stored));
        var expected = JsonSerializer.SerializeToNode(new Order("ORDER-1", "2026-01-15T10:30:00Z", Receipt), CheckoutJson.Options);
        Assert.True(JsonNode.DeepEquals(expected, JsonNode.Parse(stored)), stored);
    }

    [Fact]
    public void Save_WritesHumanReadableJson_IndentedAndEndingInANewline()
    {
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());
        repository.Save(Receipt);

        var text = File.ReadAllText(_folder.PathOf("ORDER-1.json"));

        Assert.Contains("\n  ", text, StringComparison.Ordinal);
        Assert.EndsWith("\n", text, StringComparison.Ordinal);
    }

    [Fact]
    public void Find_SavedOrder_ReadsItBackExactlyAsSaved()
    {
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.FixedAt("2026-01-15T10:30:00Z"));
        var receipt = Checkout.Calculate(Fixtures.BasketOf("COFFEE", 3, coupon: "SAVE10"), Fixtures.TinyCatalog);

        var id = repository.Save(receipt);

        Assert.Equal(new Order("ORDER-1", "2026-01-15T10:30:00Z", receipt), repository.Find(id));
    }

    [Fact]
    public void Save_ThreeTimes_NumbersTheOrdersAndFilesInSequence()
    {
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());

        var ids = new[] { repository.Save(Receipt), repository.Save(Receipt), repository.Save(Receipt) };

        Assert.Equal(["ORDER-1", "ORDER-2", "ORDER-3"], ids);
        Assert.Equal(["ORDER-1.json", "ORDER-2.json", "ORDER-3.json"], Directory.GetFiles(_folder.Path).Select(Path.GetFileName).Order());
    }

    [Fact]
    public void New_OnAFolderWithOrders_ContinuesNumbering()
    {
        var first = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());
        first.Save(Receipt);
        first.Save(Receipt);

        var second = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());

        Assert.Equal("ORDER-3", second.Save(Receipt));
        Assert.Equal("ORDER-1", second.Find("ORDER-1")?.Id);
    }

    [Fact]
    public void Save_AfterGapsAndTwoDigitNumbers_NumbersAfterTheHighest()
    {
        WriteOrderFile("ORDER-3");
        WriteOrderFile("ORDER-12");
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());

        // Not ORDER-4: alphabetically "3" > "12", but numerically 12 is the highest.
        Assert.Equal("ORDER-13", repository.Save(Receipt));
    }

    [Fact]
    public void Save_FolderWithOtherFiles_IgnoresEverythingThatIsNotAnOrder()
    {
        File.WriteAllText(_folder.PathOf("README.txt"), "not an order");
        File.WriteAllText(_folder.PathOf("ORDER-9.json.bak"), "not an order either");
        File.WriteAllText(_folder.PathOf("ORDER-9"), "no .json extension, so not an order");
        File.WriteAllText(_folder.PathOf("ORDER-09.json"), "a leading zero is not an order id");
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());

        Assert.Equal("ORDER-1", repository.Save(Receipt));
    }

    [Fact]
    public void New_FolderDoesNotExistYet_CreatesIt()
    {
        var nested = _folder.PathOf("orders", "2026");
        var repository = new JsonFileOrderRepository(nested, Clocks.Fixed());

        repository.Save(Receipt);

        Assert.True(File.Exists(Path.Combine(nested, "ORDER-1.json")));
    }

    [Theory]
    [InlineData("ORDER-2")]
    [InlineData("ORDER-0")]
    [InlineData("nonsense")]
    public void Find_IdWithoutAFile_IsNotFoundWithoutThrowing(string id)
    {
        var repository = new JsonFileOrderRepository(_folder.Path, Clocks.Fixed());
        repository.Save(Receipt);

        Assert.Null(repository.Find(id));
    }

    [Theory]
    [InlineData("../secret")]
    [InlineData("../secret.json")]
    public void Find_IdWithPathTricks_CannotEscapeTheFolder(string id)
    {
        File.WriteAllText(_folder.PathOf("secret.json"), """{"id":"x","created_at":"x","receipt":{}}""");
        var repository = new JsonFileOrderRepository(_folder.PathOf("orders"), Clocks.Fixed());

        Assert.Null(repository.Find(id));
    }

    [Fact]
    public void Save_StampsEachOrderWithTheClocksTimeAtTheMomentOfSaving()
    {
        var clock = Clocks.FixedAt("2026-01-15T10:30:00Z");
        var repository = new JsonFileOrderRepository(_folder.Path, clock);

        var first = repository.Save(Receipt);
        clock.SetUtcNow(DateTimeOffset.Parse("2026-01-16T09:00:00Z", System.Globalization.CultureInfo.InvariantCulture));
        var second = repository.Save(Receipt);

        Assert.Equal("2026-01-15T10:30:00Z", repository.Find(first)?.CreatedAt);
        Assert.Equal("2026-01-16T09:00:00Z", repository.Find(second)?.CreatedAt);
    }

    [Fact]
    public void New_WithoutAClock_UsesTheRealOne()
    {
        var repository = new JsonFileOrderRepository(_folder.Path);

        var id = repository.Save(Receipt);

        var createdAt = DateTimeOffset.Parse(repository.Find(id)!.CreatedAt, System.Globalization.CultureInfo.InvariantCulture);
        Assert.InRange(createdAt, DateTimeOffset.UtcNow.AddMinutes(-1), DateTimeOffset.UtcNow.AddMinutes(1));
    }

    private void WriteOrderFile(string id)
    {
        var order = new Order(id, "2026-01-15T10:30:00Z", Receipt);
        File.WriteAllText(_folder.PathOf($"{id}.json"), JsonSerializer.Serialize(order, CheckoutJson.Options));
    }
}
