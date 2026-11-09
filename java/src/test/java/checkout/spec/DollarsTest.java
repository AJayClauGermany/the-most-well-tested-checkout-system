package checkout.spec;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

/**
 * UNIT tests for a test helper. Test code deserves tests too: the BDD steps rely on
 * {@link Dollars#toCents}, so a bug there would hide bugs in the checkout.
 */
class DollarsTest {

    @ParameterizedTest(name = "turns ''{0}'' into {1} cents")
    @CsvSource({
        "$0.00, 0",
        "$0.50, 50",
        "$12.39, 1239",
        "$5, 500",
        "$5.5, 550",
        "'  $8.99 ', 899",
        "0.04, 4",
    })
    void convertsDollarsToWholeCents(String text, int cents) {
        assertThat(Dollars.toCents(text)).isEqualTo(cents);
    }

    @Test
    void refusesTextThatIsNotAnAmount() {
        assertThatThrownBy(() -> Dollars.toCents("ten dollars"))
                .isInstanceOf(IllegalArgumentException.class)
                .hasMessageContaining("not a dollar amount");
        assertThatThrownBy(() -> Dollars.toCents("$1.234")).hasMessageContaining("not a dollar amount");
        assertThatThrownBy(() -> Dollars.toCents("-$1.00")).hasMessageContaining("not a dollar amount");
    }
}
