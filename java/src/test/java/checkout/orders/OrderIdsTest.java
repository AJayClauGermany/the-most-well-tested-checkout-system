package checkout.orders;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.OptionalInt;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.junit.jupiter.params.provider.ValueSource;

/**
 * UNIT tests for order ids (rules section 6): {@code "ORDER-1"}, {@code "ORDER-2"}, ...
 */
class OrderIdsTest {

    @ParameterizedTest(name = "number {0} <-> id {1}")
    @CsvSource({
        "1, ORDER-1",
        "2, ORDER-2",
        "10, ORDER-10",
        "1234, ORDER-1234",
    })
    void numbersAndIdsConvertBothWays(int n, String id) {
        assertThat(OrderIds.of(n)).isEqualTo(id);
        assertThat(OrderIds.numberOf(id)).isEqualTo(OptionalInt.of(n));
    }

    @ParameterizedTest(name = "does not recognise ''{0}'' as an order id")
    @ValueSource(strings = {
        "", "ORDER-", "ORDER-0", "ORDER-01", "ORDER--1", "ORDER-1.5", "ORDER-1 ", " ORDER-1",
        "MYORDER-1", "ORDER-1x", "order-1", "ORDER-x", "1", "ORDER_1", "ORDER-99999999999",
    })
    void rejectsEverythingThatIsNotExactlyAnOrderId(String id) {
        assertThat(OrderIds.numberOf(id)).isEmpty();
    }
}
