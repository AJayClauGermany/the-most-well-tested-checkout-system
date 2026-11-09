package checkout.cli;

import static checkout.spec.Fixtures.TINY_CATALOG;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Catalog;
import checkout.spec.Json;
import java.util.List;
import java.util.stream.Stream;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

/**
 * UNIT tests for the command line program's logic (rules section 8), run in-process.
 * The real program ({@link Main}) is tested from the outside in {@code checkout.e2e.CliIT}.
 */
class CheckoutCommandTest {

    private static final String ONE_APPLE = "{\"items\":[{\"sku\":\"APPLE\",\"quantity\":1}]}";

    @Nested
    class Run {

        @Test
        void printsTheReceiptAsJsonAndExitsWithZeroForAValidBasket() {
            CommandResult result = CheckoutCommand.run(ONE_APPLE, TINY_CATALOG);

            assertThat(result.exitCode()).isEqualTo(CommandResult.SUCCESS);
            assertThat(Json.parse(result.output())).isEqualTo(Json.parse("""
                    {"lines": [{"sku": "APPLE", "name": "Apple", "quantity": 1,
                                "unit_price": 50, "line_total": 50, "line_discount": 0}],
                     "subtotal": 50, "promotion_discount": 0, "coupon_discount": 0,
                     "shipping": 499, "tax": 4, "total": 553}
                    """));
        }

        @Test
        void printsAnErrorObjectAndExitsWithOneForAnInvalidBasket() {
            CommandResult result = CheckoutCommand.run("{\"items\":[]}", TINY_CATALOG);

            assertThat(result)
                    .isEqualTo(new CommandResult("{\"error\":\"EMPTY_BASKET\"}", CommandResult.CHECKOUT_FAILED));
        }

        @Test
        void reportsEveryKindOfErrorTheSameWay() {
            CommandResult result = CheckoutCommand.run(
                    "{\"items\":[{\"sku\":\"APPLE\",\"quantity\":1}],\"coupon\":\"NOPE\"}", TINY_CATALOG);

            assertThat(result).isEqualTo(new CommandResult("{\"error\":\"UNKNOWN_COUPON\"}", 1));
        }

        @Test
        void printsCompactJsonOnASingleLine() {
            CommandResult result = CheckoutCommand.run(ONE_APPLE, TINY_CATALOG);

            assertThat(result.output()).doesNotContain("\n").doesNotContain(" ");
        }

        @Test
        void acceptsSurroundingWhitespaceAsAShellPipelineAddsANewline() {
            CommandResult result = CheckoutCommand.run("  " + ONE_APPLE + "\n", TINY_CATALOG);

            assertThat(result.exitCode()).isZero();
        }

        @Test
        void passesTheCouponAndShippingThrough() {
            CommandResult result = CheckoutCommand.run(
                    "{\"items\":[{\"sku\":\"TV\",\"quantity\":1}],\"coupon\":\"save10\",\"shipping\":\"express\"}",
                    TINY_CATALOG);

            assertThat(Json.parse(result.output()).get("coupon_discount").asInt()).isEqualTo(1000);
            assertThat(Json.parse(result.output()).get("shipping").asInt()).isEqualTo(999);
        }

        @Test
        void neverHidesABugBehindAnErrorCodeOtherExceptionsPassStraightThrough() {
            // A broken catalog makes checkout throw a NullPointerException, not a CheckoutException.
            Catalog brokenCatalog = new Catalog(null);

            assertThatThrownBy(() -> CheckoutCommand.run(ONE_APPLE, brokenCatalog))
                    .isInstanceOf(NullPointerException.class);
        }

        @Test
        void throwsRatherThanInventingAnErrorCodeWhenTheInputIsNotJson() {
            assertThatThrownBy(() -> CheckoutCommand.run("not json", TINY_CATALOG))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessageStartingWith("The input is not JSON");
            assertThatThrownBy(() -> CheckoutCommand.run("", TINY_CATALOG))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessageStartingWith("The input is not JSON");
        }
    }

    @Nested
    class ParseBasket {

        @Test
        void readsItemsCouponAndShipping() {
            Basket basket = CheckoutCommand.parseBasket(
                    "{\"items\":[{\"sku\":\"COFFEE\",\"quantity\":3},{\"sku\":\"APPLE\",\"quantity\":2}],"
                            + "\"coupon\":\"SAVE10\",\"shipping\":\"express\"}");

            assertThat(basket).isEqualTo(new Basket(
                    List.of(new BasketItem("COFFEE", 3), new BasketItem("APPLE", 2)), "SAVE10", "express"));
        }

        @Test
        void fillsInTheDefaultsWhenCouponAndShippingAreMissingOrNull() {
            assertThat(CheckoutCommand.parseBasket(ONE_APPLE))
                    .isEqualTo(new Basket(List.of(new BasketItem("APPLE", 1)), null, "standard"));
            assertThat(CheckoutCommand.parseBasket("{\"items\":[{\"sku\":\"APPLE\",\"quantity\":1}],\"coupon\":null}"))
                    .isEqualTo(new Basket(List.of(new BasketItem("APPLE", 1)), null, "standard"));
        }

        @Test
        void keepsInvalidValuesForTheCheckoutToReport() {
            // Not the command's job: unknown products, quantities out of range, odd shipping.
            Basket basket = CheckoutCommand.parseBasket(
                    "{\"items\":[{\"sku\":\"BANANA\",\"quantity\":-5}],\"shipping\":\"overnight\"}");

            assertThat(basket).isEqualTo(new Basket(List.of(new BasketItem("BANANA", -5)), null, "overnight"));
        }

        /**
         * JSON that is not a basket, and why. A {@code @MethodSource} instead of {@code @CsvSource}
         * because JSON is full of quotes and commas, which CSV would try to interpret.
         */
        static Stream<Arguments> notBaskets() {
            String objectExpected = "expected a JSON object";
            String itemsList = "'items' must be a list";
            String textSku = "every item needs a text 'sku'";
            String wholeQuantity = "every item needs a whole-number 'quantity'";
            return Stream.of(
                    Arguments.of("null", objectExpected),
                    Arguments.of("42", objectExpected),
                    Arguments.of("\"text\"", objectExpected),
                    Arguments.of("[]", objectExpected),
                    Arguments.of("{}", itemsList),
                    Arguments.of("{\"items\": \"APPLE\"}", itemsList),
                    Arguments.of("{\"items\": null}", itemsList),
                    Arguments.of("{\"items\": [1, 2]}", textSku),
                    Arguments.of("{\"items\": [null]}", textSku),
                    Arguments.of("{\"items\": [\"APPLE\"]}", textSku),
                    Arguments.of("{\"items\": [{\"sku\": 5}]}", textSku),
                    Arguments.of("{\"items\": [{\"sku\": \"APPLE\"}]}", wholeQuantity),
                    Arguments.of("{\"items\": [{\"sku\": \"APPLE\", \"quantity\": 1.5}]}", wholeQuantity),
                    Arguments.of("{\"items\": [{\"sku\": \"APPLE\", \"quantity\": \"3\"}]}", wholeQuantity),
                    Arguments.of("{\"items\": [{\"sku\": \"APPLE\", \"quantity\": 9999999999}]}", wholeQuantity),
                    Arguments.of("{\"items\": [], \"coupon\": 10}", "'coupon' must be text or null"),
                    Arguments.of("{\"items\": [], \"shipping\": [\"express\"]}", "'shipping' must be text or null"));
        }

        @ParameterizedTest(name = "rejects {0} ({1})")
        @MethodSource("notBaskets")
        void rejectsJsonThatIsNotABasket(String input, String reason) {
            assertThatThrownBy(() -> CheckoutCommand.parseBasket(input))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("The input is not a basket: " + reason);
        }

        @Test
        void rejectsTextThatIsNotJsonWithTheParsersExplanation() {
            assertThatThrownBy(() -> CheckoutCommand.parseBasket("{\"items\": [}"))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessageStartingWith("The input is not JSON: ");
        }
    }
}
