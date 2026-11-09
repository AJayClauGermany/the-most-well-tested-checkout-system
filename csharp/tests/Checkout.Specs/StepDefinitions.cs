using System.Globalization;

using Reqnroll;

namespace CheckoutSystem.Specs;

/// <summary>
/// A dollar amount from a feature file ("$12.39"), already converted to whole cents.
/// A dedicated type (rather than a plain <c>int</c>) makes the conversion below apply to
/// money and only to money: "the receipt has 2 lines" keeps its ordinary integer.
/// </summary>
public readonly record struct Money(int Cents);

/// <summary>
/// STEP DEFINITIONS: the glue between the Given/When/Then sentences in spec/features/*.feature
/// and the checkout code. Every sentence used by the features is listed in
/// spec/features/README.md; each one appears exactly once here.
///
/// The attribute text is a regular expression (anchored with ^ and $); the groups in
/// parentheses become the method's parameters.
/// </summary>
[Binding]
public sealed class StepDefinitions(CheckoutContext context)
{
    // --- Converting step arguments -------------------------------------------------------

    /// <summary>"$12.39" → 1239 cents, with integer maths only (spec/features/README.md).</summary>
    [StepArgumentTransformation(@"\$(\d+(?:\.\d{1,2})?)")]
    public static Money DollarsToMoney(string dollars) => new(Dollars.ToCents(dollars));

    // --- Given ---------------------------------------------------------------------------

    [Given(@"^the standard catalog$")]
    public void GivenTheStandardCatalog() => context.Catalog = SpecLoader.LoadCatalog();

    [Given(@"^an empty basket$")]
    public void GivenAnEmptyBasket() => context.Items.Clear();

    [Given(@"^a basket containing:$")]
    public void GivenABasketContaining(DataTable table)
    {
        foreach (var row in table.Rows)
        {
            context.Items.Add(new BasketItem(row["sku"], int.Parse(row["quantity"], CultureInfo.InvariantCulture)));
        }
    }

    [Given(@"^the customer uses the coupon ""(.*)""$")]
    public void GivenTheCustomerUsesTheCoupon(string code) => context.Coupon = code;

    [Given(@"^the customer chooses ""(.*)"" shipping$")]
    public void GivenTheCustomerChoosesShipping(string method) => context.Shipping = method;

    [Given(@"^an empty order repository$")]
    public void GivenAnEmptyOrderRepository() => context.Repository = new InMemoryOrderRepository(context.Clock);

    [Given(@"^the clock is fixed at ""(.*)""$")]
    public void GivenTheClockIsFixedAt(string timestamp) =>
        context.Clock.SetUtcNow(DateTimeOffset.Parse(timestamp, CultureInfo.InvariantCulture));

    // --- When ----------------------------------------------------------------------------

    [When(@"^the customer checks out$")]
    public void WhenTheCustomerChecksOut()
    {
        try
        {
            context.Receipt = Checkout.Calculate(context.Basket, context.Catalog);
        }
        catch (CheckoutException error)
        {
            context.Error = error;
        }
    }

    [When(@"^the customer checks out and saves the order$")]
    public void WhenTheCustomerChecksOutAndSavesTheOrder()
    {
        try
        {
            var placed = Checkout.PlaceOrder(context.Basket, context.Catalog, context.RepositoryOrFail());
            context.Receipt = placed.Receipt;
            context.LastOrderId = placed.OrderId;
        }
        catch (CheckoutException error)
        {
            context.Error = error;
        }
    }

    // --- Then: the receipt ---------------------------------------------------------------

    [Then(@"^the subtotal is (\$\d+\.\d{2})$")]
    public void ThenTheSubtotalIs(Money amount) => Assert.Equal(amount.Cents, context.ReceiptOrFail().Subtotal);

    [Then(@"^the promotion discount is (\$\d+\.\d{2})$")]
    public void ThenThePromotionDiscountIs(Money amount) => Assert.Equal(amount.Cents, context.ReceiptOrFail().PromotionDiscount);

    [Then(@"^the coupon discount is (\$\d+\.\d{2})$")]
    public void ThenTheCouponDiscountIs(Money amount) => Assert.Equal(amount.Cents, context.ReceiptOrFail().CouponDiscount);

    [Then(@"^the shipping cost is (\$\d+\.\d{2})$")]
    public void ThenTheShippingCostIs(Money amount) => Assert.Equal(amount.Cents, context.ReceiptOrFail().Shipping);

    [Then(@"^the tax is (\$\d+\.\d{2})$")]
    public void ThenTheTaxIs(Money amount) => Assert.Equal(amount.Cents, context.ReceiptOrFail().Tax);

    [Then(@"^the total is (\$\d+\.\d{2})$")]
    public void ThenTheTotalIs(Money amount) => Assert.Equal(amount.Cents, context.ReceiptOrFail().Total);

    [Then(@"^the receipt has (\d+) lines$")]
    public void ThenTheReceiptHasLines(int count) => Assert.Equal(count, context.ReceiptOrFail().Lines.Count);

    [Then(@"^the receipt line for ""(.*)"" has a quantity of (\d+)$")]
    public void ThenTheReceiptLineForHasAQuantityOf(string sku, int quantity)
    {
        var line = context.ReceiptOrFail().Lines.SingleOrDefault(l => l.Sku == sku);

        Assert.NotNull(line);
        Assert.Equal(quantity, line.Quantity);
    }

    [Then(@"^checkout fails with the error ""(.*)""$")]
    public void ThenCheckoutFailsWithTheError(string code)
    {
        Assert.Null(context.Receipt);
        Assert.NotNull(context.Error);
        Assert.Equal(code, context.Error.Code.ToSpecCode());
    }

    // --- Then: saved orders --------------------------------------------------------------

    [Then(@"^the order is saved with the id ""(.*)""$")]
    public void ThenTheOrderIsSavedWithTheId(string id) => Assert.Equal(id, context.LastOrderId);

    [Then(@"^the saved order ""(.*)"" has a total of (\$\d+\.\d{2})$")]
    public void ThenTheSavedOrderHasATotalOf(string id, Money amount)
    {
        var order = context.RepositoryOrFail().Find(id);

        Assert.NotNull(order);
        Assert.Equal(amount.Cents, order.Receipt.Total);
    }

    [Then(@"^the order repository is empty$")]
    public void ThenTheOrderRepositoryIsEmpty() => Assert.Equal(0, context.RepositoryOrFail().Count);

    [Then(@"^the order ""(.*)"" cannot be found$")]
    public void ThenTheOrderCannotBeFound(string id) => Assert.Null(context.RepositoryOrFail().Find(id));

    [Then(@"^the saved order ""(.*)"" was created at ""(.*)""$")]
    public void ThenTheSavedOrderWasCreatedAt(string id, string timestamp)
    {
        var order = context.RepositoryOrFail().Find(id);

        Assert.NotNull(order);
        // Timestamps are compared as text, exactly as they appear in a saved order.
        Assert.Equal(timestamp, order.CreatedAt);
    }
}
