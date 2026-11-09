using CheckoutSystem.Storage;

using Microsoft.Data.Sqlite;

namespace CheckoutSystem.IntegrationTests;

/// <summary>
/// INTEGRATION TESTS for the SQLite order repository (Microsoft.Data.Sqlite).
///
/// A real database, but no server to install: SQLite is one file. Each test gets its own
/// file in a temporary folder; the repository is disposed (which closes the database) before
/// the folder is deleted.
/// </summary>
public sealed class SqliteOrderRepositoryTests : IDisposable
{
    private static readonly Receipt Receipt = Fixtures.OneAppleReceipt;

    private readonly TemporaryFolder _folder = new();
    private readonly List<SqliteOrderRepository> _open = [];
    private readonly string _path;

    public SqliteOrderRepositoryTests()
    {
        _path = _folder.PathOf("orders.db");
    }

    public void Dispose()
    {
        foreach (var repository in _open)
        {
            repository.Dispose();
        }

        _folder.Dispose();
    }

    [Fact]
    public void SaveAndFind_RoundTripsAnOrder()
    {
        var repository = Open(Clocks.FixedAt("2026-01-15T10:30:00Z"));
        var receipt = Checkout.Calculate(Fixtures.BasketOf("COFFEE", 3, coupon: "SAVE10"), Fixtures.TinyCatalog);

        var id = repository.Save(receipt);

        Assert.Equal("ORDER-1", id);
        Assert.Equal(new Order("ORDER-1", "2026-01-15T10:30:00Z", receipt), repository.Find(id));
    }

    [Fact]
    public void Find_SavedOrder_FollowsOrderSchemaJson()
    {
        var repository = Open();
        var id = repository.Save(Receipt);

        Assert.Empty(SchemaValidator.ErrorsFor("order.schema.json", repository.Find(id)));
    }

    [Fact]
    public void Save_ThreeTimes_NumbersTheOrdersInSequence()
    {
        var repository = Open();

        var ids = new[] { repository.Save(Receipt), repository.Save(Receipt), repository.Save(Receipt) };

        Assert.Equal(["ORDER-1", "ORDER-2", "ORDER-3"], ids);
    }

    [Fact]
    public void New_OnAnExistingDatabase_ContinuesNumbering()
    {
        var first = Open();
        first.Save(Receipt);
        first.Save(Receipt);
        CloseNow(first);

        var second = Open();

        Assert.Equal("ORDER-3", second.Save(Receipt));
        Assert.Equal("ORDER-2", second.Find("ORDER-2")?.Id);
    }

    [Fact]
    public void Save_ReallyWritesToTheDatabaseFile_CheckedWithASeparateConnection()
    {
        var repository = Open();
        repository.Save(Receipt);
        CloseNow(repository);

        using var connection = new SqliteConnection($"Data Source={_path}");
        connection.Open();
        using var command = connection.CreateCommand();
        command.CommandText = "SELECT COUNT(*) FROM orders";

        Assert.Equal(1L, command.ExecuteScalar());
    }

    [Fact]
    public void SaveAndFind_AfterDispose_ThrowInsteadOfSilentlyDoingNothing()
    {
        var repository = Open();
        CloseNow(repository);

        Assert.Throws<ObjectDisposedException>(() => repository.Save(Receipt));
        Assert.Throws<ObjectDisposedException>(() => repository.Find("ORDER-1"));
    }

    [Fact]
    public void New_InMemory_WorksWithoutAFile()
    {
        using var repository = new SqliteOrderRepository(SqliteOrderRepository.InMemory, Clocks.Fixed());

        Assert.Equal("ORDER-1", repository.Save(Receipt));
        Assert.Equal(Receipt, repository.Find("ORDER-1")?.Receipt);
    }

    [Theory]
    [InlineData("ORDER-2")]
    [InlineData("ORDER-0")]
    [InlineData("ORDER-01")]
    [InlineData("ORDER-x")]
    [InlineData("1")]
    [InlineData("")]
    [InlineData("ORDER-1; DROP TABLE orders")]
    public void Find_IdThatIsNotASavedOrder_IsNotFoundWithoutThrowing(string id)
    {
        var repository = Open();
        repository.Save(Receipt);

        Assert.Null(repository.Find(id));
        Assert.NotNull(repository.Find("ORDER-1")); // and the table is still there
    }

    [Fact]
    public void Save_StampsEachOrderWithTheClocksTimeAtTheMomentOfSaving()
    {
        var clock = Clocks.FixedAt("2026-01-15T10:30:00Z");
        var repository = Open(clock);

        var first = repository.Save(Receipt);
        clock.Advance(TimeSpan.FromHours(22.5));
        var second = repository.Save(Receipt);

        Assert.Equal("2026-01-15T10:30:00Z", repository.Find(first)?.CreatedAt);
        Assert.Equal("2026-01-16T09:00:00Z", repository.Find(second)?.CreatedAt);
    }

    [Fact]
    public void New_WithoutAClock_UsesTheRealOne()
    {
        using var repository = new SqliteOrderRepository(_path);

        var id = repository.Save(Receipt);

        var createdAt = DateTimeOffset.Parse(repository.Find(id)!.CreatedAt, System.Globalization.CultureInfo.InvariantCulture);
        Assert.InRange(createdAt, DateTimeOffset.UtcNow.AddMinutes(-1), DateTimeOffset.UtcNow.AddMinutes(1));
    }

    /// <summary>Opens a repository on the test's database file and remembers it, so Dispose can close it.</summary>
    private SqliteOrderRepository Open(TimeProvider? timeProvider = null)
    {
        var repository = new SqliteOrderRepository(_path, timeProvider ?? Clocks.Fixed());
        _open.Add(repository);
        return repository;
    }

    /// <summary>Closes a repository early, inside a test (and forgets it, so Dispose does not close it twice).</summary>
    private void CloseNow(SqliteOrderRepository repository)
    {
        _open.Remove(repository);
        repository.Dispose();
    }
}
