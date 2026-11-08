package checkout;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import checkout.spec.SpecCase;
import checkout.spec.SpecFiles;
import checkout.spec.Totals;
import java.util.stream.Stream;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Named;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.MethodSource;

/**
 * A DATA-DRIVEN test: one test method, run once for every case in
 * {@code spec/test-cases/*.json} (the same cases every language uses).
 *
 * <p>This is test-driven development (TDD): the tests exist before the code.
 * They are disabled for now because {@link Checkout#checkout} is not written yet.
 * When you start implementing it, delete the {@code @Disabled} line and watch the
 * tests go from red to green, one case at a time.
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

    @Disabled("checkout() is not implemented yet - remove this @Disabled line when you start implementing")
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
