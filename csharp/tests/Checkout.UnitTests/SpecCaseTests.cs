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
///
/// These 60 cases are the acceptance tests of the whole project: every language runs
/// exactly the same ones, so the six implementations can be compared.
/// </summary>
[Collection(StandardCatalogCollectionDefinition.Name)]
public class SpecCaseTests(StandardCatalogFixture spec)
{
    public static TheoryData<string> CaseIds() =>
        new(SpecLoader.LoadCases().Select(specCase => specCase.Id));

    [Theory]
    [MemberData(nameof(CaseIds))]
    public void Calculate_SharedSpecCase_MatchesExpectedTotalsOrError(string caseId)
    {
        // Arrange
        var specCase = spec.Cases.Single(c => c.Id == caseId);

        if (specCase.ExpectedError is { } expectedError)
        {
            // Act + Assert: checkout must fail with the expected error code.
            var exception = Assert.Throws<CheckoutException>(() => Checkout.Calculate(specCase.Basket, spec.Catalog));
            Assert.Equal(expectedError, exception.Code);
        }
        else
        {
            // Act
            var receipt = Checkout.Calculate(specCase.Basket, spec.Catalog);

            // Assert: compare all six totals at once.
            Assert.Equal(specCase.Expected, Totals.From(receipt));
        }
    }
}

/// <summary>
/// Declares the "standard catalog" collection: every test class marked
/// <c>[Collection(StandardCatalogCollectionDefinition.Name)]</c> shares ONE <see cref="StandardCatalogFixture"/>
/// (the spec files are read once) and runs in sequence with the others in the collection.
/// </summary>
[CollectionDefinition(Name)]
public sealed class StandardCatalogCollectionDefinition : ICollectionFixture<StandardCatalogFixture>
{
    public const string Name = "standard catalog";
}
