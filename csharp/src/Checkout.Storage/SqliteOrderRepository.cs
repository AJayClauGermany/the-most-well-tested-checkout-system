using System.Globalization;
using System.Text.Json;

using Microsoft.Data.Sqlite;

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

    private readonly SqliteConnection _connection;
    private readonly TimeProvider _timeProvider;
    private bool _disposed;

    /// <param name="path">The database file (created if missing), or <see cref="InMemory"/>.</param>
    /// <param name="timeProvider">The clock. <c>null</c> means the real one.</param>
    public SqliteOrderRepository(string path, TimeProvider? timeProvider = null)
    {
        // No connection pooling: closing the repository really closes the file, so a test
        // (or a user on Windows) can delete it straight away.
        var connectionString = new SqliteConnectionStringBuilder { DataSource = path, Pooling = false }.ToString();
        _connection = new SqliteConnection(connectionString);
        _connection.Open();
        _timeProvider = timeProvider ?? TimeProvider.System;

        using var create = _connection.CreateCommand();
        create.CommandText = """
            CREATE TABLE IF NOT EXISTS orders (
                number     INTEGER PRIMARY KEY AUTOINCREMENT,
                created_at TEXT NOT NULL,
                receipt    TEXT NOT NULL
            )
            """;
        create.ExecuteNonQuery();
    }

    public string Save(Receipt receipt)
    {
        ObjectDisposedException.ThrowIf(_disposed, this);

        // "$name" parameters keep the values apart from the SQL, so a receipt can never be
        // mistaken for a command (no SQL injection).
        using var insert = _connection.CreateCommand();
        insert.CommandText = "INSERT INTO orders (created_at, receipt) VALUES ($created_at, $receipt) RETURNING number";
        insert.Parameters.AddWithValue("$created_at", Timestamps.Format(_timeProvider.GetUtcNow()));
        insert.Parameters.AddWithValue("$receipt", JsonSerializer.Serialize(receipt, CheckoutJson.Options));
        var number = Convert.ToInt32(insert.ExecuteScalar(), CultureInfo.InvariantCulture);
        return OrderIds.Format(number);
    }

    public Order? Find(string id)
    {
        ObjectDisposedException.ThrowIf(_disposed, this);

        if (OrderIds.Parse(id) is not { } number)
        {
            return null;
        }

        using var select = _connection.CreateCommand();
        select.CommandText = "SELECT created_at, receipt FROM orders WHERE number = $number";
        select.Parameters.AddWithValue("$number", number);
        using var row = select.ExecuteReader();
        if (!row.Read())
        {
            return null;
        }

        var receipt = JsonSerializer.Deserialize<Receipt>(row.GetString(1), CheckoutJson.Options)!;
        return new Order(id, row.GetString(0), receipt);
    }

    /// <summary>Closes the database. Saving or finding afterwards throws <see cref="ObjectDisposedException"/>.</summary>
    public void Dispose()
    {
        _disposed = true;
        _connection.Dispose();
    }
}
