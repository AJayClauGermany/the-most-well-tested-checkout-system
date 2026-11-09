package checkout;

import static checkout.spec.Fixtures.TINY_CATALOG;
import static checkout.spec.Fixtures.basketOf;
import static checkout.spec.Fixtures.item;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatCode;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.util.List;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

/**
 * UNIT tests for basket validation (rules section 4.1, rows 1 to 4) and the order errors are
 * reported in. The check order matters: a basket can be wrong in several ways at once.
 *
 * <p>Boundaries are tested on both sides (0/1 and 99/100): "{@code >} instead of {@code >=}"
 * is the classic off-by-one bug, and the one mutation testing finds most often.
 */
class ValidationTest {

    /** Asserts that validating the basket fails with exactly this error code. */
    private static void assertFailsWith(Basket basket, ErrorCode code) {
        assertThatThrownBy(() -> Validation.validateItems(basket, TINY_CATALOG))
                .isInstanceOf(CheckoutException.class)
                .extracting("code")
                .isEqualTo(code);
    }

    @Nested
    class IsValidQuantity {

        @Test
        void hasTheLimitsFromTheRules() {
            assertThat(Validation.MIN_QUANTITY).isEqualTo(1);
            assertThat(Validation.MAX_QUANTITY).isEqualTo(99);
        }

        // @ValueSource: the simplest parameterised test, one value per run.
        @ParameterizedTest(name = "accepts {0}")
        @ValueSource(ints = {1, 2, 50, 98, 99})
        void acceptsWholeNumbersFromOneToNinetyNine(int quantity) {
            assertThat(Validation.isValidQuantity(quantity)).isTrue();
        }

        @ParameterizedTest(name = "rejects {0}")
        @ValueSource(ints = {0, -1, 100, 101, 1000, Integer.MIN_VALUE, Integer.MAX_VALUE})
        void rejectsEverythingOutsideOneToNinetyNine(int quantity) {
            assertThat(Validation.isValidQuantity(quantity)).isFalse();
        }
    }

    @Nested
    class ValidateItems {

        @Test
        void acceptsANormalBasket() {
            Basket basket = basketOf(item("APPLE", 1), item("COFFEE", 99));

            assertThatCode(() -> Validation.validateItems(basket, TINY_CATALOG)).doesNotThrowAnyException();
        }

        @Test
        void row1AnEmptyBasket() {
            assertFailsWith(basketOf(), ErrorCode.EMPTY_BASKET);
        }

        @Test
        void row2ASkuThatIsNotInTheCatalog() {
            assertFailsWith(basketOf("BANANA", 1), ErrorCode.UNKNOWN_PRODUCT);
        }

        @Test
        void row2SkusAreCaseSensitive() {
            assertFailsWith(basketOf("apple", 1), ErrorCode.UNKNOWN_PRODUCT);
        }

        @Test
        void row2AnUnknownProductOnAnyLineNotOnlyTheFirst() {
            assertFailsWith(
                    basketOf(item("APPLE", 1), item("PEN", 1), item("BANANA", 1)),
                    ErrorCode.UNKNOWN_PRODUCT);
        }

        @ParameterizedTest(name = "row 3: a quantity of {0}")
        @ValueSource(ints = {0, -1, 100})
        void row3ABadQuantity(int quantity) {
            assertFailsWith(basketOf("APPLE", quantity), ErrorCode.INVALID_QUANTITY);
        }

        @Test
        void row3ABadQuantityOnALaterLine() {
            assertFailsWith(basketOf(item("APPLE", 1), item("PEN", 0)), ErrorCode.INVALID_QUANTITY);
        }

        @Test
        void row4LinesOfOneSkuThatAddUpToMoreThanNinetyNine() {
            assertFailsWith(basketOf(item("APPLE", 50), item("APPLE", 50)), ErrorCode.INVALID_QUANTITY);
        }

        @Test
        void row4LinesOfOneSkuThatAddUpToExactlyNinetyNineAreFine() {
            Basket basket = basketOf(item("APPLE", 50), item("PEN", 99), item("APPLE", 49));

            assertThatCode(() -> Validation.validateItems(basket, TINY_CATALOG)).doesNotThrowAnyException();
        }

        @Test
        void row4NinetyNineOfOneSkuAndNinetyNineOfAnotherIsFineTheLimitIsPerSku() {
            Basket basket = basketOf(item("APPLE", 99), item("PEN", 99));

            assertThatCode(() -> Validation.validateItems(basket, TINY_CATALOG)).doesNotThrowAnyException();
        }
    }

    @Nested
    class ErrorMessages {

        // The code is what programs look at; the message is for the human reading a log.
        @Test
        void explainTheProblem() {
            assertThatThrownBy(() -> Validation.validateItems(basketOf(), TINY_CATALOG))
                    .hasMessage("The basket has no items");
            assertThatThrownBy(() -> Validation.validateItems(basketOf("BANANA", 1), TINY_CATALOG))
                    .hasMessage("Unknown product: BANANA");
            assertThatThrownBy(() -> Validation.validateItems(basketOf("APPLE", 0), TINY_CATALOG))
                    .hasMessage("Invalid quantity for APPLE: 0");
            assertThatThrownBy(() -> Validation.validateItems(
                            basketOf(item("APPLE", 60), item("APPLE", 40)), TINY_CATALOG))
                    .hasMessage("Too many APPLE in total: 100");
        }
    }

    @Nested
    class CheckOrder {

        @Test
        void anUnknownProductBeatsABadQuantityEvenOnALaterLine() {
            assertFailsWith(basketOf(item("APPLE", 0), item("BANANA", 1)), ErrorCode.UNKNOWN_PRODUCT);
        }

        @Test
        void aBadSingleQuantityBeatsABadMergedQuantity() {
            // Both rows 3 and 4 apply; either way the code is INVALID_QUANTITY.
            assertFailsWith(basketOf(item("APPLE", 100), item("APPLE", 1)), ErrorCode.INVALID_QUANTITY);
        }

        @Test
        void ignoresTheCouponAndShippingMethodThoseAreCheckedLaterByCheckout() {
            Basket basket = new Basket(List.of(item("APPLE", 1)), "NOPE", "overnight");

            assertThatCode(() -> Validation.validateItems(basket, TINY_CATALOG)).doesNotThrowAnyException();
        }
    }
}
