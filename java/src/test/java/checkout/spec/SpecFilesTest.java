package checkout.spec;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import checkout.Catalog;
import checkout.Product;
import checkout.Promotion;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import org.junit.jupiter.api.Test;

/**
 * Tests for the test-support code: does {@link SpecFiles} read the shared spec correctly?
 *
 * <p>Yes, test code can need tests too! If the loader silently lost a case, the
 * shared spec test would happily pass without ever checking it.
 */
class SpecFilesTest {

    @Test
    void standardCatalogHasSixProducts() {
        Catalog catalog = SpecFiles.loadCatalog();

        assertEquals(6, catalog.products().size());
    }

    @Test
    void coffeeIsBuyTwoGetOneFree() {
        Product coffee = SpecFiles.loadCatalog().findBySku("COFFEE").orElseThrow();

        assertEquals(899, coffee.unitPrice());
        assertEquals(new Promotion.BuyXGetYFree(2, 1), coffee.promotion());
    }

    @Test
    void notebookHasABulkDiscount() {
        Product notebook = SpecFiles.loadCatalog().findBySku("NOTEBOOK").orElseThrow();

        assertEquals(new Promotion.BulkDiscount(10, 10), notebook.promotion());
    }

    @Test
    void loadsAllSixtySharedCases() {
        List<SpecCase> cases = SpecFiles.loadAllCases();

        assertEquals(60, cases.size());
    }

    @Test
    void caseIdsAreUnique() {
        Set<String> seen = new HashSet<>();

        for (SpecCase specCase : SpecFiles.loadAllCases()) {
            // Set.add returns false when the id was already in the set.
            assertTrue(seen.add(specCase.id()), "duplicate id: " + specCase.id());
        }
    }

    @Test
    void everyCaseExpectsEitherTotalsOrAnError() {
        for (SpecCase specCase : SpecFiles.loadAllCases()) {
            boolean hasTotals = specCase.expected() != null;
            boolean hasError = specCase.expectedError() != null;

            // Exactly one of the two: "true != false" or "false != true".
            assertTrue(hasTotals != hasError, specCase.id() + " must have exactly one of expected / expected_error");
        }
    }
}
