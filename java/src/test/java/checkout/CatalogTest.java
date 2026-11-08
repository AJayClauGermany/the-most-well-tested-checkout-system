package checkout;

// "import static" lets us write assertEquals(...) instead of Assertions.assertEquals(...).
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.util.List;
import java.util.Optional;
import org.junit.jupiter.api.Test;

/**
 * An example UNIT test: it checks one small piece of code (the {@link Catalog} class)
 * on its own, quickly, without files, networks or databases.
 *
 * <p>How JUnit works, in short:
 * <ul>
 *   <li>Maven's Surefire plugin runs every class whose name ends in {@code Test}.</li>
 *   <li>Every method marked with {@code @Test} is one test. JUnit creates a fresh
 *       {@code CatalogTest} object for each test, so tests cannot affect each other.</li>
 *   <li>A test passes when it finishes without an exception. The {@code assert...}
 *       methods throw an exception (and so fail the test) when something is wrong.</li>
 * </ul>
 *
 * <p>Each test follows the same three steps: <b>Arrange</b> (set things up),
 * <b>Act</b> (call the code under test) and <b>Assert</b> (check the result).
 * Test method names describe the behaviour they check, so a failing test's name
 * already tells you what broke.
 */
class CatalogTest {

    // A tiny, made-up catalog. Small test data makes tests easy to read.
    private final Catalog catalog = new Catalog(List.of(
            new Product("APPLE", "Apple", 50),
            new Product("PEN", "Pen", 120)));

    @Test
    void findsAProductBySku() {
        // Arrange: the catalog above is all we need.

        // Act
        Optional<Product> found = catalog.findBySku("PEN");

        // Assert: assertEquals(expected, actual) - the expected value comes first!
        assertEquals(Optional.of(new Product("PEN", "Pen", 120)), found);
    }

    @Test
    void findsNothingForAnUnknownSku() {
        Optional<Product> found = catalog.findBySku("BANANA");

        assertTrue(found.isEmpty(), "BANANA is not in the catalog");
    }

    @Test
    void skusAreCaseSensitive() {
        // Rules section 4.1: "apple" is an unknown product.
        Optional<Product> found = catalog.findBySku("apple");

        assertTrue(found.isEmpty(), "'apple' must not find 'APPLE'");
    }
}
