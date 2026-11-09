package checkout;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.util.Optional;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.junit.jupiter.params.provider.NullAndEmptySource;
import org.junit.jupiter.params.provider.ValueSource;

/**
 * UNIT tests for coupons (rules section 4.4). The minimum spend of WELCOME5 is tested on
 * both sides of the boundary (2499 / 2500).
 */
class CouponsTest {

    /** Asserts that the coupon fails with exactly this error code. */
    private static void assertFailsWith(String coupon, int afterPromotions, ErrorCode code) {
        assertThatThrownBy(() -> Coupons.discount(coupon, afterPromotions))
                .isInstanceOf(CheckoutException.class)
                .extracting("code")
                .isEqualTo(code);
    }

    @Nested
    class NormalizeCode {

        @ParameterizedTest(name = "turns ''{0}'' into {1}")
        @CsvSource({
            "SAVE10, SAVE10",
            "save10, SAVE10",
            "'  save10 ', SAVE10",
            "'\tWelcome5\n', WELCOME5",
            "nope, NOPE",
        })
        void trimsAndUpperCases(String input, String expected) {
            assertThat(Coupons.normalizeCode(input)).contains(expected);
        }

        // @NullAndEmptySource adds null and "" to the values below.
        @ParameterizedTest(name = "treats ''{0}'' as no coupon")
        @NullAndEmptySource
        @ValueSource(strings = {"   ", "\t\n"})
        void treatsMissingAndBlankAsNoCoupon(String input) {
            assertThat(Coupons.normalizeCode(input)).isEqualTo(Optional.empty());
        }
    }

    @Nested
    class Discount {

        @Test
        void isZeroWithoutACoupon() {
            assertThat(Coupons.discount(null, 10000)).isZero();
            assertThat(Coupons.discount("", 10000)).isZero();
            assertThat(Coupons.discount("  ", 10000)).isZero();
        }

        @Nested
        class Save10 {

            @Test
            void takesTenPercentOffRoundedHalfUp() {
                assertThat(Coupons.discount("SAVE10", 500)).isEqualTo(50);
                assertThat(Coupons.discount("SAVE10", 225)).as("22.5 -> 23").isEqualTo(23);
                assertThat(Coupons.discount("SAVE10", 224)).as("22.4 -> 22").isEqualTo(22);
                assertThat(Coupons.discount("SAVE10", 1798)).as("179.8 -> 180").isEqualTo(180);
            }

            @Test
            void hasNoMinimumSpend() {
                assertThat(Coupons.discount("SAVE10", 1)).as("0.1 -> 0").isZero();
                assertThat(Coupons.discount("SAVE10", 5)).as("0.5 -> 1").isEqualTo(1);
            }

            @Test
            void ignoresCaseAndSurroundingSpaces() {
                assertThat(Coupons.discount(" save10 ", 500)).isEqualTo(50);
            }
        }

        @Nested
        class Welcome5 {

            @Test
            void takesExactlyFiveDollarsOff() {
                assertThat(Coupons.discount("WELCOME5", 2500)).isEqualTo(500);
                assertThat(Coupons.discount("WELCOME5", 99999)).isEqualTo(500);
            }

            @Test
            void worksAtExactlyTheMinimumAndFailsOneCentBelowIt() {
                assertThat(Coupons.discount("WELCOME5", 2500)).isEqualTo(500);
                assertFailsWith("WELCOME5", 2499, ErrorCode.COUPON_MINIMUM_NOT_MET);
            }

            @Test
            void ignoresCaseAndSurroundingSpaces() {
                assertThat(Coupons.discount("welcome5", 2500)).isEqualTo(500);
            }
        }

        @Test
        void rejectsACodeItDoesNotKnow() {
            assertFailsWith("FREESTUFF", 10000, ErrorCode.UNKNOWN_COUPON);
        }

        @Test
        void explainsWhatWentWrongInTheErrorMessage() {
            assertThatThrownBy(() -> Coupons.discount("freestuff", 10000))
                    .hasMessage("Unknown coupon: FREESTUFF");
            assertThatThrownBy(() -> Coupons.discount("WELCOME5", 2499))
                    .hasMessage("WELCOME5 needs a spend of at least 2500 cents, got 2499");
        }

        @Test
        void doesNotAcceptACodeThatMerelyStartsOrEndsWithARealCode() {
            assertFailsWith("SAVE100", 10000, ErrorCode.UNKNOWN_COUPON);
            assertFailsWith("WELCOME50", 10000, ErrorCode.UNKNOWN_COUPON);
            assertFailsWith("XSAVE10", 10000, ErrorCode.UNKNOWN_COUPON);
        }
    }
}
