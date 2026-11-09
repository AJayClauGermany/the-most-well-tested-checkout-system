package checkout;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

/**
 * UNIT tests for tax (rules section 4.6).
 */
class TaxTest {

    @Test
    void isEightPercent() {
        assertThat(Tax.TAX_PERCENT).isEqualTo(8);
    }

    @ParameterizedTest(name = "taxes {0} cents as {1} cents ({2})")
    @CsvSource({
        "50, 4, 8% of $0.50 is exactly 4 cents",
        "120, 10, 9.6 cents rounds up",
        "240, 19, 19.2 cents rounds down",
        "625, 50, exactly 50 cents",
        "1706, 136, 136.48 rounds down",
        "11880, 950, 950.4 rounds down",
        "0, 0, no spend - no tax",
    })
    void taxesEightPercentRoundedHalfUp(int amount, int tax, String why) {
        assertThat(Tax.taxOn(amount)).as(why).isEqualTo(tax);
    }

    @Test
    void roundsPlainFractionsTheUsualWay() {
        // 8% of a whole number of cents can never be exactly half a cent (x * 8 never ends in 50),
        // so the half-up rule is pinned down by the 10% coupon tests; here we check plain rounding.
        assertThat(Tax.taxOn(1875)).as("exactly 150.00").isEqualTo(150);
        assertThat(Tax.taxOn(1881)).as("150.48").isEqualTo(150);
        assertThat(Tax.taxOn(1888)).as("151.04").isEqualTo(151);
    }
}
