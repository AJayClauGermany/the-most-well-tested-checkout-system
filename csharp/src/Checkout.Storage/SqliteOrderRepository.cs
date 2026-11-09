namespace CheckoutSystem.Storage;

/// <summary>
/// Keeps orders in one table of an SQLite database file. The database row number becomes
/// the order number, so ids are sequential and survive closing and reopening the file.
/// <see cref="Dispose"/> closes the database; use it in a <c>using</c> statement.
/// </summary>
public sealed class SqliteOrderRepository : IOrderRepository, IDisposable
{
    /// <summary>SQLite's name for a throw-away database that lives only in memory.</summary>
    public const string InMemory = ":memory:";

    /// <param name="path">The database file (created if missing), or <see cref="InMemory"/>.</param>
    /// <param name="timeProvider">The clock. <c>null</c> means the real one.</param>
    public SqliteOrderRepository(string path, TimeProvider? timeProvider = null)
    {
        _ = path;
        _ = timeProvider;
        throw new NotImplementedException();
    }

    public string Save(Receipt receipt) =>
        throw new NotImplementedException();

    public Order? Find(string id) =>
        throw new NotImplementedException();

    public void Dispose() =>
        throw new NotImplementedException();
}
