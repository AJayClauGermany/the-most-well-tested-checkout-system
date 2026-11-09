package checkout;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

/**
 * UNIT tests for rounding (rules section 7). Rounding is where checkout systems classically go
 * wrong, so it gets its own small method and its own tests.
 *
 * <p>This file also introduces two tools used all over the test suite:
 * <ul>
 *   <li><b>AssertJ</b>: {@code assertThat(actual).isEqualTo(expected)} reads like a sentence
 *       and prints a clear message when it fails. (The baseline tests use JUnit's own
 *       {@code assertEquals(expected, actual)}; both are fine, pick one per file.)</li>
 *   <li><b>{@code @ParameterizedTest}</b>: one test method, run once per row of a table.
 *       {@code @CsvSource} writes the table as comma-separated text; each row becomes the
 *       method's arguments.</li>
 * </ul>
 */
class MoneyTest {

    // The table from rules section 7, as a parameterised test: one row, one test run.
    @ParameterizedTest(name = "{2} cents rounds to {3}")
    @CsvSource({
        // amount, percent, exact value, expected
        "224, 10, 22.4, 22",
        "225, 10, 22.5, 23",
        "226, 10, 22.6, 23",
    })
    void roundsToTheNearestCentWithHalvesGoingUp(int amount, int percent, String exact, int expected) {
        assertThat(Money.roundHalfUp(amount, percent)).as("%s rounded", exact).isEqualTo(expected);
    }

    @Test
    void roundsAnExactHalfUpWhereBankersRoundingWouldGoToTheEvenNumber() {
        // 10% of $52.35 is exactly 523.5 cents. Math.round gives 524 too, but Java's
        // RoundingMode.HALF_EVEN would give 524 as well here - so also check 22.5 -> 23 (above),
        // where half-to-even gives 22. Both together pin the rule down.
        assertThat(Money.roundHalfUp(5235, 10)).isEqualTo(524);
    }

    @Test
    void doesNotRoundWhenTheResultIsAlreadyWhole() {
        assertThat(Money.roundHalfUp(500, 10)).isEqualTo(50);
        assertThat(Money.roundHalfUp(50, 8)).isEqualTo(4);
    }

    @Test
    void roundsEightPercentOfOneTwentyUpAndOfTwoFortyDown() {
        assertThat(Money.roundHalfUp(120, 8)).as("9.6 cents").isEqualTo(10);
        assertThat(Money.roundHalfUp(240, 8)).as("19.2 cents").isEqualTo(19);
    }

    @Test
    void givesZeroForAnAmountOfZeroAndTheWholeAmountForOneHundredPercent() {
        assertThat(Money.roundHalfUp(0, 10)).isZero();
        assertThat(Money.roundHalfUp(899, 100)).isEqualTo(899);
    }

    @Test
    void givesZeroForZeroPercent() {
        assertThat(Money.roundHalfUp(899, 0)).isZero();
    }

    @Test
    void staysExactWhereFloatingPointMathsWouldDrift() {
        // 10% of 3 cents is 0.3 - it must round to 0, not 1 (0.1 + 0.2 != 0.3 with doubles).
        assertThat(Money.roundHalfUp(3, 10)).isZero();
        // 8% of 1050 = 84.0 exactly; a double implementation might land on 83.99999 -> 83.
        assertThat(Money.roundHalfUp(1050, 8)).isEqualTo(84);
    }
}
