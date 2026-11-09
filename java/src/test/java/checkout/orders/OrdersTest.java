package checkout.orders;

import static checkout.spec.Fixtures.TINY_CATALOG;
import static checkout.spec.Fixtures.basketOf;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import checkout.Checkout;
import checkout.CheckoutException;
import checkout.ErrorCode;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneOffset;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests for {@link Orders#placeOrder} = checkout + save, using the hand-written FAKE
 * repository. {@link PlaceOrderMockitoTest} checks the same behaviour with a Mockito MOCK -
 * compare the two to see the difference between a fake and a mock.
 */
class OrdersTest {

    private static final Clock FIXED_CLOCK = Clock.fixed(Instant.parse("2026-01-15T10:30:00Z"), ZoneOffset.UTC);

    @Test
    void checksOutAndSavesTheReceiptReturningBothTheReceiptAndTheNewId() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        PlacedOrder placed = Orders.placeOrder(basketOf("APPLE", 1), TINY_CATALOG, repository);

        assertThat(placed.orderId()).isEqualTo("ORDER-1");
        assertThat(placed.receipt()).isEqualTo(Checkout.checkout(basketOf("APPLE", 1), TINY_CATALOG));
        assertThat(repository.find("ORDER-1").orElseThrow().receipt()).isEqualTo(placed.receipt());
    }

    @Test
    void leavesTheRepositoryEmptyAfterAFailedCheckout() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        assertThatThrownBy(() -> Orders.placeOrder(basketOf("APPLE", 0), TINY_CATALOG, repository))
                .isInstanceOf(CheckoutException.class)
                .extracting("code")
                .isEqualTo(ErrorCode.INVALID_QUANTITY);

        assertThat(repository.count()).isZero();
        assertThat(repository.find("ORDER-1")).isEmpty();
    }

    @Test
    void numbersSeveralPlacedOrdersInSequence() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        Orders.placeOrder(basketOf("APPLE", 1), TINY_CATALOG, repository);
        PlacedOrder second = Orders.placeOrder(basketOf("PEN", 2), TINY_CATALOG, repository);

        assertThat(second.orderId()).isEqualTo("ORDER-2");
    }
}
