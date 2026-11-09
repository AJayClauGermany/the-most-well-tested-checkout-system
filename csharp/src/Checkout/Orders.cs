using System.Globalization;
using System.Text.RegularExpressions;

namespace CheckoutSystem;

/// <summary>A saved order, exactly as spec/schema/order.schema.json describes it (rules §6).</summary>
/// <param name="Id">"ORDER-1", "ORDER-2", … – sequential per repository.</param>
/// <param name="CreatedAt">When it was saved: UTC, ISO 8601, whole seconds, trailing Z ("2026-01-15T10:30:00Z").</param>
/// <param name="Receipt">The receipt exactly as it was saved.</param>
public sealed record Order(string Id, string CreatedAt, Receipt Receipt);

/// <summary>
/// Where orders are kept (rules §6). Implemented in memory (<see cref="InMemoryOrderRepository"/>),
/// as JSON files and in SQLite (both in the Checkout.Storage project).
///
/// Every implementation takes a <see cref="TimeProvider"/>: the rules demand that the time
/// source can be replaced, because <c>created_at</c> is the one value that changes from run
/// to run. Production uses <see cref="TimeProvider.System"/>; tests pass a fake that always
/// answers the same time.
/// </summary>
public interface IOrderRepository
{
    /// <summary>Saves a receipt and returns the new order's id.</summary>
    string Save(Receipt receipt);

    /// <summary>Looks an order up. <c>null</c> means "not found" – a normal answer, not an error.</summary>
    Order? Find(string id);
}

/// <summary>Order ids: "ORDER-1", "ORDER-2", … (rules §6).</summary>
public static partial class OrderIds
{
    /// <summary>The id of the n-th order saved into a repository: <c>Format(7)</c> is "ORDER-7".</summary>
    public static string Format(int number) =>
        string.Create(CultureInfo.InvariantCulture, $"ORDER-{number}");

    /// <summary>
    /// The inverse of <see cref="Format"/>: "ORDER-7" gives 7. Anything that is not an order
    /// id ("ORDER-07", "order-7", "7", "../secret") gives <c>null</c>.
    /// </summary>
    public static int? Parse(string id)
    {
        var match = OrderIdPattern().Match(id);
        return match.Success && int.TryParse(match.Groups[1].ValueSpan, CultureInfo.InvariantCulture, out var number)
            ? number
            : null;
    }

    /// <summary>
    /// "ORDER-" and a positive whole number with no leading zero, exactly as order.schema.json
    /// spells it. <c>\z</c> means "the very end" (<c>$</c> would also accept a newline at the end),
    /// and <c>[0-9]</c> rather than <c>\d</c>, which in .NET also matches digits from other scripts.
    /// </summary>
    [GeneratedRegex(@"\AORDER-([1-9][0-9]*)\z")]
    private static partial Regex OrderIdPattern();
}

/// <summary>The <c>created_at</c> format (rules §6).</summary>
public static class Timestamps
{
    /// <summary>Formats a moment the way <c>created_at</c> must look: UTC, whole seconds, trailing Z.</summary>
    public static string Format(DateTimeOffset moment) =>
        moment.ToUniversalTime().ToString("yyyy-MM-dd'T'HH:mm:ss'Z'", CultureInfo.InvariantCulture);
}

/// <summary>
/// Keeps orders in a plain list. The simplest repository there is – and the hand-written
/// FAKE that the unit tests and the BDD scenarios use instead of a real store.
/// </summary>
public sealed class InMemoryOrderRepository : IOrderRepository
{
    private readonly List<Order> _orders = [];
    private readonly TimeProvider _timeProvider;

    /// <param name="timeProvider">The clock. <c>null</c> means the real one, <see cref="TimeProvider.System"/>.</param>
    public InMemoryOrderRepository(TimeProvider? timeProvider = null)
    {
        _timeProvider = timeProvider ?? TimeProvider.System;
    }

    /// <summary>How many orders are saved. Handy for tests ("the order repository is empty").</summary>
    public int Count => _orders.Count;

    public string Save(Receipt receipt)
    {
        var id = OrderIds.Format(_orders.Count + 1);
        // A copy of the lines, so that whoever gave us the receipt cannot change the saved order afterwards.
        _orders.Add(new Order(id, Timestamps.Format(_timeProvider.GetUtcNow()), receipt with { Lines = [.. receipt.Lines] }));
        return id;
    }

    public Order? Find(string id) =>
        _orders.Find(order => order.Id == id);
}
