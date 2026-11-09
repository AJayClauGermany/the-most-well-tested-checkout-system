namespace CheckoutSystem.Storage;

/// <summary>
/// Keeps each order as <c>&lt;folder&gt;/ORDER-&lt;n&gt;.json</c>, in the shape of
/// spec/schema/order.schema.json. Ids continue where the files in the folder left off,
/// so reopening the same folder later keeps numbering correctly.
/// </summary>
public sealed class JsonFileOrderRepository : IOrderRepository
{
    /// <param name="folder">Where the files go. Created if it does not exist.</param>
    /// <param name="timeProvider">The clock. <c>null</c> means the real one.</param>
    public JsonFileOrderRepository(string folder, TimeProvider? timeProvider = null)
    {
        _ = folder;
        _ = timeProvider;
        throw new NotImplementedException();
    }

    public string Save(Receipt receipt) =>
        throw new NotImplementedException();

    public Order? Find(string id) =>
        throw new NotImplementedException();
}
