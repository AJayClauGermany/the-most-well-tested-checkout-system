using CheckoutSystem.UnitTests.SpecSupport;

namespace CheckoutSystem.UnitTests;

/// <summary>
/// DATA-DRIVEN TESTS from the shared spec: one test per case in spec/test-cases/*.json.
///
/// [MemberData] feeds the [Theory] with the rows returned by <see cref="CaseIds"/>, so
/// each case id becomes its own test, e.g.
/// Calculate_SharedSpecCase_MatchesExpectedTotalsOrError(caseId: "subtotal-single-item").
///
/// We pass only the id (a plain string): it is short, readable in the test results,
/// and lets you run a single case by name.
/// </summary>
public class SpecCaseTests
{
    public static TheoryData<string> CaseIds() =>
        new(SpecLoader.LoadCases().Select(specCase => specCase.Id));

    // TDD "red tests waiting": these tests are written before the code. They are
    // skipped for now so the build stays green. Delete the `Skip = "..."` text when you
    // start implementing Checkout.Calculate, then make the cases pass one by one.
    //
    // Why is Skip on [MemberData] and not on [Theory]? A skipped [Theory] is reported
    // as ONE skipped test. Skipping the data instead reports every case separately
    // (60 skipped tests, each with its id), so you can see exactly what is waiting.
    [Theory]
    [MemberData(nameof(CaseIds), Skip = "Checkout.Calculate() is not implemented yet – remove this Skip when you start implementing")]
    public void Calculate_SharedSpecCase_MatchesExpectedTotalsOrError(string caseId)
    {
        // Arrange
        var specCase = SpecLoader.LoadCase(caseId);
        var catalog = SpecLoader.LoadCatalog();

        if (specCase.ExpectedError is { } expectedError)
        {
            // Act + Assert: checkout must fail with the expected error code.
            var exception = Assert.Throws<CheckoutException>(() => Checkout.Calculate(specCase.Basket, catalog));
            Assert.Equal(expectedError, exception.Code);
        }
        else
        {
            // Act
            var receipt = Checkout.Calculate(specCase.Basket, catalog);

            // Assert: compare all six totals at once.
            Assert.Equal(specCase.Expected, Totals.From(receipt));
        }
    }
}
