package checkout;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import checkout.spec.SpecCase;
import checkout.spec.SpecFiles;
import checkout.spec.Totals;
import java.util.stream.Stream;
import org.junit.jupiter.api.Named;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.MethodSource;

/**
 * A DATA-DRIVEN test: one test method, run once for every case in
 * {@code spec/test-cases/*.json} (the same cases every language uses).
 *
 * <p>This is test-driven development (TDD): the tests existed before the code. They were
 * disabled while {@link Checkout#checkout} was not written yet; removing the
 * {@code @Disabled} line was the first step of the implementation, and the 60 cases went
 * from red to green one by one.
 *
 * <p>{@code @MethodSource} is the third way to feed a parameterised test (next to
 * {@code @ValueSource} and {@code @CsvSource}): a method supplies the arguments, here by
 * reading the JSON files.
 */
class SharedSpecCasesTest {

    /**
     * Supplies the cases. {@code Named.of(id, case)} gives each run the case id as
     * its name, so the test report shows e.g. {@code promo-coffee-three-bags-one-free}.
     */
    static Stream<Named<SpecCase>> specCases() {
        return SpecFiles.loadAllCases().stream()
                .map(specCase -> Named.of(specCase.id(), specCase));
    }

    @ParameterizedTest(name = "{0}")
    @MethodSource("specCases")
    void checkoutMatchesTheSharedSpec(SpecCase specCase) {
        Catalog catalog = SpecFiles.loadCatalog();
        String why = specCase.id() + " (" + specCase.file() + "): " + specCase.description();

        if (specCase.expectedError() != null) {
            // This case expects the checkout to fail with a specific error code.
            CheckoutException error = assertThrows(
                    CheckoutException.class,
                    () -> Checkout.checkout(specCase.basket(), catalog),
                    why);
            assertEquals(specCase.expectedError(), error.getCode(), why);
        } else {
            // This case expects a receipt: compare all six totals at once.
            Receipt receipt = Checkout.checkout(specCase.basket(), catalog);
            assertEquals(specCase.expected(), Totals.of(receipt), why);
        }
    }
}
