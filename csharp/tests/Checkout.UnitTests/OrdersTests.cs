using Microsoft.Extensions.Time.Testing;

using NSubstitute;

namespace CheckoutSystem.UnitTests;

/// <summary>
/// Unit tests for saving orders (rules §6), and the project's showcase of TEST DOUBLES.
///
/// A test double stands in for a real dependency:
/// * a FAKE is a simple working version – <see cref="InMemoryOrderRepository"/> is one,
/// * a STUB gives canned answers – <see cref="FakeTimeProvider"/> is a clock that says what we tell it,
/// * a MOCK records how it was called – NSubstitute builds one from <see cref="IOrderRepository"/>,
///   so we can prove that <c>Save</c> was never called.
/// </summary>
public class TimestampsTests
{
    [Fact]
    public void Format_UtcMoment_WritesWholeSecondsAndATrailingZ()
    {
        var moment = new DateTimeOffset(2026, 1, 15, 10, 30, 0, TimeSpan.Zero);

        Assert.Equal("2026-01-15T10:30:00Z", Timestamps.Format(moment));
    }

    [Fact]
    public void Format_Milliseconds_AreDroppedNotRounded()
    {
        var moment = new DateTimeOffset(2026, 1, 15, 10, 30, 0, 999, TimeSpan.Zero);

        Assert.Equal("2026-01-15T10:30:00Z", Timestamps.Format(moment));
    }

    [Fact]
    public void Format_OtherTimeZone_IsConvertedToUtc()
    {
        var moment = new DateTimeOffset(2026, 1, 15, 12, 30, 0, TimeSpan.FromHours(2));

        Assert.Equal("2026-01-15T10:30:00Z", Timestamps.Format(moment));
    }
}

public class OrderIdsTests
{
    [Theory]
    [InlineData(1, "ORDER-1")]
    [InlineData(2, "ORDER-2")]
    [InlineData(10, "ORDER-10")]
    [InlineData(1234, "ORDER-1234")]
    public void FormatAndParse_AreInverses(int number, string id)
    {
        Assert.Equal(id, OrderIds.Format(number));
        Assert.Equal(number, OrderIds.Parse(id));
    }

    [Theory]
    [InlineData("")]
    [InlineData("ORDER-")]
    [InlineData("ORDER-0")]
    [InlineData("ORDER-01")]
    [InlineData("ORDER--1")]
    [InlineData("ORDER-1.5")]
    [InlineData("ORDER-1 ")]
    [InlineData(" ORDER-1")]
    [InlineData("MYORDER-1")]
    [InlineData("ORDER-1x")]
    [InlineData("order-1")]
    [InlineData("ORDER-x")]
    [InlineData("1")]
    [InlineData("ORDER_1")]
    [InlineData("../secret")]
    public void Parse_TextThatIsNotAnOrderId_GivesNull(string id)
    {
        Assert.Null(OrderIds.Parse(id));
    }
}

public class InMemoryOrderRepositoryTests
{
    private static readonly Receipt Receipt = Fixtures.OneAppleReceipt;

    [Fact]
    public void New_Repository_StartsEmpty()
    {
        var repository = new InMemoryOrderRepository(Clocks.Fixed());

        Assert.Equal(0, repository.Count);
        Assert.Null(repository.Find("ORDER-1"));
    }

    [Fact]
    public void Save_ThreeTimes_NumbersTheOrdersInSequence()
    {
        var repository = new InMemoryOrderRepository(Clocks.Fixed());

        var ids = new[] { repository.Save(Receipt), repository.Save(Receipt), repository.Save(Receipt) };

        Assert.Equal(["ORDER-1", "ORDER-2", "ORDER-3"], ids);
        Assert.Equal(3, repository.Count);
    }

    [Fact]
    public void Find_SavedOrder_HasTheReceiptAndTheTimeTheClockSaid()
    {
        var repository = new InMemoryOrderRepository(Clocks.FixedAt("2026-01-15T10:30:00Z"));

        var id = repository.Save(Receipt);

        Assert.Equal(new Order("ORDER-1", "2026-01-15T10:30:00Z", Receipt), repository.Find(id));
    }

    [Fact]
    public void Save_ReadsTheClockAtTheMomentOfSaving()
    {
        // The stub clock is moved forward between the two saves, so each order gets its own time.
        var clock = Clocks.FixedAt("2026-01-15T10:30:00Z");
        var repository = new InMemoryOrderRepository(clock);

        var first = repository.Save(Receipt);
        clock.Advance(TimeSpan.FromHours(22.5));
        var second = repository.Save(Receipt);

        Assert.Equal("2026-01-15T10:30:00Z", repository.Find(first)?.CreatedAt);
        Assert.Equal("2026-01-16T09:00:00Z", repository.Find(second)?.CreatedAt);
    }

    [Theory]
    [InlineData("ORDER-2")]
    [InlineData("ORDER-0")]
    [InlineData("ORDER-99")]
    [InlineData("order-1")]
    [InlineData("")]
    public void Find_IdThatWasNeverHandedOut_IsNotFound(string id)
    {
        var repository = new InMemoryOrderRepository(Clocks.Fixed());
        repository.Save(Receipt);

        Assert.Null(repository.Find(id));
    }

    [Fact]
    public void New_WithoutAClock_UsesTheRealOne()
    {
        var before = DateTimeOffset.UtcNow;
        var repository = new InMemoryOrderRepository();

        var id = repository.Save(Receipt);

        // We cannot know the exact second, but it must be "now": between before and after.
        var createdAt = DateTimeOffset.Parse(repository.Find(id)!.CreatedAt, System.Globalization.CultureInfo.InvariantCulture);
        Assert.InRange(createdAt, before.AddSeconds(-1), DateTimeOffset.UtcNow.AddSeconds(1));
    }
}

public class PlaceOrderTests
{
    private static readonly Catalog Catalog = Fixtures.TinyCatalog;

    [Fact]
    public void PlaceOrder_ValidBasket_ChecksOutAndSavesTheReceipt()
    {
        var repository = new InMemoryOrderRepository(Clocks.Fixed());

        var placed = Checkout.PlaceOrder(Fixtures.BasketOf("APPLE", 1), Catalog, repository);

        Assert.Equal("ORDER-1", placed.OrderId);
        Assert.Equal(Fixtures.OneAppleReceipt, placed.Receipt);
        Assert.Equal(Fixtures.OneAppleReceipt, repository.Find("ORDER-1")?.Receipt);
    }

    [Fact]
    public void PlaceOrder_ValidBasket_SavesExactlyWhatItReturns_CheckedWithAMock()
    {
        // Substitute.For creates a MOCK: an object that implements the interface and records
        // every call. .Returns(...) tells it what to answer.
        var repository = Substitute.For<IOrderRepository>();
        repository.Save(Arg.Any<Receipt>()).Returns("ORDER-7");

        var placed = Checkout.PlaceOrder(Fixtures.BasketOf("TV", 1), Catalog, repository);

        repository.Received(1).Save(placed.Receipt);
        repository.DidNotReceive().Find(Arg.Any<string>());
        Assert.Equal("ORDER-7", placed.OrderId); // the id comes from the repository, not from PlaceOrder
    }

    [Fact]
    public void PlaceOrder_InvalidBasket_NeverCallsSave()
    {
        // Property P14 in spec/properties.md: a failed checkout saves nothing.
        var repository = Substitute.For<IOrderRepository>();

        Assert.Throws<CheckoutException>(() => Checkout.PlaceOrder(new Basket([]), Catalog, repository));
        var unknown = Assert.Throws<CheckoutException>(() => Checkout.PlaceOrder(Fixtures.BasketOf("BANANA", 1), Catalog, repository));

        Assert.Equal(ErrorCode.UnknownProduct, unknown.Code);
        repository.DidNotReceive().Save(Arg.Any<Receipt>());
    }

    [Fact]
    public void PlaceOrder_InvalidBasket_LeavesTheFakeRepositoryEmpty()
    {
        var repository = new InMemoryOrderRepository(Clocks.Fixed());

        Assert.Throws<CheckoutException>(() => Checkout.PlaceOrder(Fixtures.BasketOf("APPLE", 0), Catalog, repository));

        Assert.Equal(0, repository.Count);
        Assert.Null(repository.Find("ORDER-1"));
    }
}
