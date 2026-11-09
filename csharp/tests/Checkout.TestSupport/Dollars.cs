using System.Text.RegularExpressions;

namespace CheckoutSystem.TestSupport;

/// <summary>
/// Test helper: turns "$12.39" into 1239 cents – without ever touching a floating point
/// number (spec/features/README.md insists, and rules §1 explains why).
/// </summary>
public static partial class Dollars
{
    /// <summary>"$12.39" or "12.39" → 1239. "$5" → 500. "$5.5" → 550.</summary>
    /// <exception cref="FormatException">The text is not a dollar amount.</exception>
    public static int ToCents(string text)
    {
        var match = DollarAmount().Match(text.Trim());
        if (!match.Success)
        {
            throw new FormatException($"not a dollar amount: \"{text}\"");
        }

        var dollars = int.Parse(match.Groups["dollars"].Value, System.Globalization.CultureInfo.InvariantCulture);
        // "5" means 50 cents, so pad to two digits before parsing. No cents at all means "00".
        var cents = int.Parse(match.Groups["cents"].Value.PadRight(2, '0'), System.Globalization.CultureInfo.InvariantCulture);
        return (dollars * 100) + cents;
    }

    [GeneratedRegex(@"^\$?(?<dollars>\d+)(?:\.(?<cents>\d{1,2}))?$")]
    private static partial Regex DollarAmount();
}
