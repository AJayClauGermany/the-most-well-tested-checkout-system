using System.Globalization;

using Microsoft.Extensions.Time.Testing;

namespace CheckoutSystem.TestSupport;

/// <summary>
/// Clocks for tests. <see cref="FakeTimeProvider"/> (from Microsoft.Extensions.TimeProvider.Testing)
/// is a STUB for <see cref="TimeProvider"/>: <c>GetUtcNow()</c> answers whatever the test set,
/// and <c>Advance()</c> moves it forward on purpose. Nothing in a test ever waits for real time.
/// </summary>
public static class Clocks
{
    /// <summary>The moment most tests use: 10:30 on 15 January 2026, UTC – as in spec/examples/order.json.</summary>
    public const string ExampleTimestamp = "2026-01-15T10:30:00Z";

    /// <summary>A clock stopped at <paramref name="utcTimestamp"/>, e.g. "2026-01-15T10:30:00Z".</summary>
    public static FakeTimeProvider FixedAt(string utcTimestamp) =>
        new(DateTimeOffset.Parse(utcTimestamp, CultureInfo.InvariantCulture));

    /// <summary>A clock stopped at <see cref="ExampleTimestamp"/>.</summary>
    public static FakeTimeProvider Fixed() => FixedAt(ExampleTimestamp);
}
