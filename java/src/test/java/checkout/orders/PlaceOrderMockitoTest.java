package checkout.orders;

import static checkout.spec.Fixtures.ONE_APPLE_RECEIPT;
import static checkout.spec.Fixtures.TINY_CATALOG;
import static checkout.spec.Fixtures.basketOf;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;

import checkout.CheckoutException;
import java.time.Clock;
import java.time.Instant;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

/**
 * TEST DOUBLES with Mockito: a MOCK is an object that records how it was called, so the test
 * can check the <em>interaction</em> ("save was called exactly once with this receipt",
 * "save was never called") instead of the end state.
 *
 * <p>Mockito's inline mock maker (the default since Mockito 5) can mock classes too, not
 * only interfaces - here {@code java.time.Clock}. On JDK 21 Mockito is loaded as a Java agent
 * (see {@code argLine} in {@code pom.xml}), otherwise the JVM prints a warning about
 * libraries attaching agents to themselves at run time.
 */
@Tag("mock")
class PlaceOrderMockitoTest {

    @Test
    void aFailedCheckoutNeverCallsSave() {
        OrderRepository repository = mock(OrderRepository.class);

        assertThatThrownBy(() -> Orders.placeOrder(basketOf(), TINY_CATALOG, repository))
                .isInstanceOf(CheckoutException.class);
        assertThatThrownBy(() -> Orders.placeOrder(basketOf("BANANA", 1), TINY_CATALOG, repository))
                .isInstanceOf(CheckoutException.class);

        verify(repository, never()).save(any());
        verifyNoInteractions(repository);
    }

    @Test
    void aSuccessfulCheckoutSavesExactlyTheReceiptItReturnsOnce() {
        OrderRepository repository = mock(OrderRepository.class);
        // A canned answer: whatever is saved, the repository says it is ORDER-7.
        when(repository.save(any())).thenReturn("ORDER-7");

        PlacedOrder placed = Orders.placeOrder(basketOf("APPLE", 1), TINY_CATALOG, repository);

        verify(repository, times(1)).save(placed.receipt());
        verify(repository, never()).find(any());
        assertThat(placed.orderId()).as("the id comes from the repository, not from placeOrder").isEqualTo("ORDER-7");
        assertThat(placed.receipt()).isEqualTo(ONE_APPLE_RECEIPT);
    }

    @Test
    void theRepositoryAsksTheClockOncePerSave() {
        Clock clock = mock(Clock.class);
        when(clock.instant()).thenReturn(Instant.parse("2026-01-15T10:30:00Z"));
        InMemoryOrderRepository repository = new InMemoryOrderRepository(clock);

        repository.save(ONE_APPLE_RECEIPT);
        repository.save(ONE_APPLE_RECEIPT);

        verify(clock, times(2)).instant();
        assertThat(repository.find("ORDER-2").orElseThrow().createdAt()).isEqualTo("2026-01-15T10:30:00Z");
    }

    @Test
    void aFailedCheckoutDoesNotEvenLookAtTheClock() {
        Clock clock = mock(Clock.class);
        InMemoryOrderRepository repository = new InMemoryOrderRepository(clock);

        assertThatThrownBy(() -> Orders.placeOrder(basketOf(), TINY_CATALOG, repository))
                .isInstanceOf(CheckoutException.class);

        verifyNoInteractions(clock);
    }
}
