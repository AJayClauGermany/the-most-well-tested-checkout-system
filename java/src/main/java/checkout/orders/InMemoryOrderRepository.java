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

    // "synchronized": only one thread at a time may save, so two saves can never get the same id.
    // (The tests run in parallel; a repository shared between threads must not break.)
    @Override
    public synchronized String save(Receipt receipt) {
        String id = OrderIds.of(orders.size() + 1);
        // Receipt is an immutable record (with an unmodifiable list of lines), so keeping the
        // very same object is safe: nobody can change the saved order afterwards.
        orders.add(new Order(id, Timestamps.now(clock), receipt));
        return id;
    }

    @Override
    public synchronized Optional<Order> find(String id) {
        return orders.stream().filter(order -> order.id().equals(id)).findFirst();
    }

    /**
     * How many orders are saved. Handy for tests ("the order repository is empty").
     *
     * @return the number of saved orders
     */
    public synchronized int count() {
        return orders.size();
    }
}
