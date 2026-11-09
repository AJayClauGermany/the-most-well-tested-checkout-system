using System.Text.Json;

namespace CheckoutSystem.Storage;

/// <summary>
/// Keeps each order as <c>&lt;folder&gt;/ORDER-&lt;n&gt;.json</c>, in the shape of
/// spec/schema/order.schema.json. Ids continue where the files in the folder left off,
/// so reopening the same folder later keeps numbering correctly.
/// </summary>
public sealed class JsonFileOrderRepository : IOrderRepository
{
    private readonly string _folder;
    private readonly TimeProvider _timeProvider;

    /// <param name="folder">Where the files go. Created if it does not exist.</param>
    /// <param name="timeProvider">The clock. <c>null</c> means the real one.</param>
    public JsonFileOrderRepository(string folder, TimeProvider? timeProvider = null)
    {
        _folder = folder;
        _timeProvider = timeProvider ?? TimeProvider.System;
        Directory.CreateDirectory(folder);
    }

    public string Save(Receipt receipt)
    {
        var order = new Order(
            OrderIds.Format(HighestOrderNumber() + 1),
            Timestamps.Format(_timeProvider.GetUtcNow()),
            receipt);
        File.WriteAllText(PathOf(order.Id), JsonSerializer.Serialize(order, CheckoutJson.IndentedOptions) + "\n");
        return order.Id;
    }

    public Order? Find(string id)
    {
        // Only accept real order ids: "../secret" must not read files outside the folder.
        if (OrderIds.Parse(id) is null || !File.Exists(PathOf(id)))
        {
            return null;
        }

        return JsonSerializer.Deserialize<Order>(File.ReadAllText(PathOf(id)), CheckoutJson.Options);
    }

    private string PathOf(string id) => Path.Combine(_folder, $"{id}.json");

    /// <summary>The highest order number among the <c>ORDER-&lt;n&gt;.json</c> files in the folder, or 0 when there are none.</summary>
    private int HighestOrderNumber() =>
        Directory.EnumerateFiles(_folder, "*.json")
            .Select(path => OrderIds.Parse(Path.GetFileNameWithoutExtension(path)) ?? 0)
            .DefaultIfEmpty(0)
            .Max();
}
