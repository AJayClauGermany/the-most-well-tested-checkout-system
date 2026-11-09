package checkout.orders;

import checkout.Receipt;
import java.time.Clock;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

/**
 * Keeps orders in a plain list. The simplest repository there is - and the hand-written
 * <em>fake</em> the tests of the checkout itself use instead of files or a database.
 */
public final class InMemoryOrderRepository implements OrderRepository {

    private final List<Order> orders = new ArrayList<>();
    private final Clock clock;

    /**
     * A repository that stamps orders with the given clock's time.
     *
     * @param clock the time source for {@code created_at}; tests pass {@link Clock#fixed}
     */
    public InMemoryOrderRepository(Clock clock) {
        this.clock = clock;
    }

    /** A repository that uses the real clock. */
    public InMemoryOrderRepository() {
        this(Clock.systemUTC());
    }

    @Override
    public String save(Receipt receipt) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    @Override
    public Optional<Order> find(String id) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * How many orders are saved. Handy for tests ("the order repository is empty").
     *
     * @return the number of saved orders
     */
    public int count() {
        throw new UnsupportedOperationException("not implemented yet " + orders + clock);
    }
}
