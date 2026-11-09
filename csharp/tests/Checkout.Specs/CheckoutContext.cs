using Microsoft.Extensions.Time.Testing;

namespace CheckoutSystem.Specs;

/// <summary>
/// The state of one scenario, carried from step to step: the basket being built, the
/// receipt or error that came out, the order repository. Reqnroll creates one per scenario
/// and hands it to the step definition class ("context injection"), so scenarios can never
/// leak state into each other.
/// </summary>
public sealed class CheckoutContext
{
    public Catalog Catalog { get; set; } = new([]);

    public List<BasketItem> Items { get; } = [];

    public string? Coupon { get; set; }

    public string Shipping { get; set; } = "standard";

    public Receipt? Receipt { get; set; }

    public CheckoutException? Error { get; set; }

    /// <summary>The clock every repository in this scenario reads. "Given the clock is fixed at …" changes it.</summary>
    public FakeTimeProvider Clock { get; } = Clocks.FixedAt("2026-01-01T00:00:00Z");

    public InMemoryOrderRepository? Repository { get; set; }

    public string? LastOrderId { get; set; }

    public Basket Basket => new(Items, Coupon, Shipping);

    /// <summary>The receipt, or a clear failure when the scenario did not check out successfully.</summary>
    public Receipt ReceiptOrFail()
    {
        Assert.Null(Error);
        Assert.NotNull(Receipt);
        return Receipt;
    }

    /// <summary>The repository, or a clear failure when the scenario forgot "Given an empty order repository".</summary>
    public InMemoryOrderRepository RepositoryOrFail()
    {
        Assert.NotNull(Repository);
        return Repository;
    }
}
