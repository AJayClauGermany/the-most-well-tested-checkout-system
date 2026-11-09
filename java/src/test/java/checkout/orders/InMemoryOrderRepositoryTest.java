package checkout.orders;

import static checkout.spec.Fixtures.ONE_APPLE_RECEIPT;
import static org.assertj.core.api.Assertions.assertThat;

import checkout.Receipt;
import checkout.ReceiptLine;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneOffset;
import java.util.ArrayList;
import java.util.List;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests for the in-memory repository - which is also the hand-written FAKE used by other
 * tests. A fake is a test double with a real, simple implementation (a list instead of a
 * database). It is tested on its own so the tests that rely on it can trust it.
 *
 * <p>The clock is a STUB: {@code Clock.fixed(...)} always answers the same time, so
 * {@code created_at} can be compared exactly instead of "roughly now" (which would be flaky).
 */
class InMemoryOrderRepositoryTest {

    private static final Clock FIXED_CLOCK = Clock.fixed(Instant.parse("2026-01-15T10:30:00Z"), ZoneOffset.UTC);

    @Test
    void startsEmpty() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        assertThat(repository.count()).isZero();
        assertThat(repository.find("ORDER-1")).isEmpty();
    }

    @Test
    void numbersOrdersOneTwoThreeInTheOrderTheyAreSaved() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        List<String> ids = List.of(
                repository.save(ONE_APPLE_RECEIPT),
                repository.save(ONE_APPLE_RECEIPT),
                repository.save(ONE_APPLE_RECEIPT));

        assertThat(ids).containsExactly("ORDER-1", "ORDER-2", "ORDER-3");
        assertThat(repository.count()).isEqualTo(3);
    }

    @Test
    void findsASavedOrderWithItsReceiptAndTheTimeTheClockSaid() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        String id = repository.save(ONE_APPLE_RECEIPT);

        assertThat(repository.find(id)).contains(new Order("ORDER-1", "2026-01-15T10:30:00Z", ONE_APPLE_RECEIPT));
    }

    @Test
    void stampsEachOrderWithTheClocksTimeAtTheMomentOfSaving() {
        // A clock we can move by hand: a tiny hand-written stub.
        List<Instant> ticks = new ArrayList<>(List.of(
                Instant.parse("2026-01-15T10:30:00Z"), Instant.parse("2026-01-16T09:00:00Z")));
        Clock clock = new Clock() {
            @Override
            public ZoneOffset getZone() {
                return ZoneOffset.UTC;
            }

            @Override
            public Clock withZone(java.time.ZoneId zone) {
                return this;
            }

            @Override
            public Instant instant() {
                return ticks.remove(0);
            }
        };
        InMemoryOrderRepository repository = new InMemoryOrderRepository(clock);

        String first = repository.save(ONE_APPLE_RECEIPT);
        String second = repository.save(ONE_APPLE_RECEIPT);

        assertThat(repository.find(first).orElseThrow().createdAt()).isEqualTo("2026-01-15T10:30:00Z");
        assertThat(repository.find(second).orElseThrow().createdAt()).isEqualTo("2026-01-16T09:00:00Z");
        assertThat(ticks).as("the clock was read exactly once per save").isEmpty();
    }

    @Test
    void keepsTheReceiptAsItWasSaved() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);
        Receipt receipt = new Receipt(
                List.of(new ReceiptLine("PEN", "Pen", 2, 120, 240, 0)), 240, 0, 0, 499, 19, 758);

        String id = repository.save(receipt);

        // Records are immutable, so "equal" is as good as "the same": nothing can change it later.
        assertThat(repository.find(id).orElseThrow().receipt()).isEqualTo(receipt);
    }

    @Test
    void answersNotFoundForIdsThatWereNeverHandedOut() {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);
        repository.save(ONE_APPLE_RECEIPT);

        assertThat(repository.find("ORDER-2")).isEmpty();
        assertThat(repository.find("ORDER-0")).isEmpty();
        assertThat(repository.find("ORDER-99")).isEmpty();
        assertThat(repository.find("order-1")).isEmpty();
        assertThat(repository.find("")).isEmpty();
    }

    @Test
    void usesTheRealClockWhenNoneIsGiven() {
        Instant before = Instant.now().minusSeconds(1);
        InMemoryOrderRepository repository = new InMemoryOrderRepository();

        String id = repository.save(ONE_APPLE_RECEIPT);

        // Not an exact comparison (that would be flaky): the time must be "now, give or take".
        Instant createdAt = Instant.parse(repository.find(id).orElseThrow().createdAt());
        assertThat(createdAt).isBetween(before, Instant.now().plusSeconds(1));
    }
}
