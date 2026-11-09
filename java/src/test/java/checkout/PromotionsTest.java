package checkout;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

/**
 * UNIT tests for the two kinds of promotion (rules section 4.2), mostly as PARAMETERISED
 * tests: the same check for a whole table of quantities, with the boundary (2/3 coffees,
 * 9/10 notebooks) on both sides.
 */
class PromotionsTest {

    private static final Promotion BUY_TWO_GET_ONE = new Promotion.BuyXGetYFree(2, 1);
    private static final Promotion TEN_PERCENT_OFF_TEN = new Promotion.BulkDiscount(10, 10);

    @Test
    void isZeroWhenTheProductHasNoPromotion() {
        assertThat(Promotions.discount(null, 50, 899)).isZero();
    }

    @Nested
    class BuyXGetYFree {

        // Coffee is buy 2 get 1 free at $8.99.
        @ParameterizedTest(name = "{0} coffees -> {1} free")
        @CsvSource({
            "1, 0",
            "2, 0",
            "3, 1",
            "4, 1",
            "5, 1",
            "6, 2",
            "99, 33",
        })
        void givesOneFreeItemPerCompleteGroupOfThree(int quantity, int free) {
            assertThat(Promotions.discount(BUY_TWO_GET_ONE, quantity, 899)).isEqualTo(free * 899);
        }

        @Test
        void usesThePromotionsOwnNumbersNotTwoAndOne() {
            // Buy 1 get 1 free: every second item is free.
            Promotion buyOneGetOne = new Promotion.BuyXGetYFree(1, 1);
            assertThat(Promotions.discount(buyOneGetOne, 4, 100)).isEqualTo(200);
            assertThat(Promotions.discount(buyOneGetOne, 5, 100)).isEqualTo(200);
            // Buy 3 get 2 free: groups of 5, 2 of each group free.
            Promotion buyThreeGetTwo = new Promotion.BuyXGetYFree(3, 2);
            assertThat(Promotions.discount(buyThreeGetTwo, 4, 100)).isZero();
            assertThat(Promotions.discount(buyThreeGetTwo, 5, 100)).isEqualTo(200);
            assertThat(Promotions.discount(buyThreeGetTwo, 10, 100)).isEqualTo(400);
        }

        @Test
        void multipliesByTheUnitPrice() {
            assertThat(Promotions.discount(BUY_TWO_GET_ONE, 3, 100)).isEqualTo(100);
            assertThat(Promotions.discount(BUY_TWO_GET_ONE, 3, 1)).isEqualTo(1);
        }
    }

    @Nested
    class BulkDiscount {

        // Notebooks are 10% off from 10 at $3.49.
        @ParameterizedTest(name = "{0} notebooks -> {1} cents off ({2})")
        @CsvSource({
            "9, 0, one short of the minimum",
            "10, 349, exactly the minimum: 10% of $34.90",
            "11, 384, 383.9 cents rounds up",
            "15, 524, exactly 523.5 cents rounds UP",
            "99, 3455, 10% of $345.51 is 3455.1",
        })
        void takesTenPercentOffFromTenNotebooks(int quantity, int discount, String why) {
            assertThat(Promotions.discount(TEN_PERCENT_OFF_TEN, quantity, 349)).as(why).isEqualTo(discount);
        }

        @Test
        void usesThePromotionsOwnMinimumAndPercentage() {
            Promotion quarterOffFromThree = new Promotion.BulkDiscount(3, 25);

            assertThat(Promotions.discount(quarterOffFromThree, 2, 100)).isZero();
            assertThat(Promotions.discount(quarterOffFromThree, 3, 100)).isEqualTo(75);
            assertThat(Promotions.discount(quarterOffFromThree, 3, 101)).as("75.75 rounds up").isEqualTo(76);
        }
    }
}
