namespace CheckoutSystem.UnitTests;

/// <summary>
/// Test code deserves tests too: the BDD steps rely on <see cref="Dollars.ToCents"/>,
/// so a bug there would hide bugs in the checkout.
/// </summary>
public class DollarsTests
{
    [Theory]
    [InlineData("$0.00", 0)]
    [InlineData("$0.50", 50)]
    [InlineData("$12.39", 1239)]
    [InlineData("$5", 500)]
    [InlineData("$5.5", 550)]
    [InlineData("  $8.99 ", 899)]
    [InlineData("0.04", 4)]
    public void ToCents_DollarAmount_GivesWholeCents(string text, int cents)
    {
        Assert.Equal(cents, Dollars.ToCents(text));
    }

    [Theory]
    [InlineData("ten dollars")]
    [InlineData("$1.234")]
    [InlineData("$-1.00")]
    public void ToCents_TextThatIsNotAnAmount_Throws(string text)
    {
        var exception = Assert.Throws<FormatException>(() => Dollars.ToCents(text));

        Assert.Contains("not a dollar amount", exception.Message, StringComparison.Ordinal);
    }
}
