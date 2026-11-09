package checkout.store;

import checkout.Receipt;
import checkout.orders.Order;
import checkout.orders.OrderRepository;
import java.time.Clock;
import java.util.Optional;

/**
 * An order repository backed by a SQL database, through plain JDBC ({@code java.sql}).
 *
 * <p>Orders live in one table, {@code orders}, with an auto-numbered {@code number} column
 * that becomes the order number, so ids are sequential and survive closing and reopening
 * the database. The tests use <a href="https://h2database.com/">H2</a>, an embedded database
 * that needs no server: {@code jdbc:h2:file:/some/folder/orders} or {@code jdbc:h2:mem:}.
 */
public final class JdbcOrderRepository implements OrderRepository, AutoCloseable {

    private final Clock clock;

    /**
     * Connects to the database and creates the {@code orders} table if it does not exist yet.
     *
     * @param jdbcUrl a JDBC URL, e.g. {@code "jdbc:h2:file:/tmp/orders"}
     * @param clock   the time source for {@code created_at}
     * @throws StorageException when the database cannot be opened
     */
    public JdbcOrderRepository(String jdbcUrl, Clock clock) {
        this.clock = clock;
        throw new UnsupportedOperationException("not implemented yet " + jdbcUrl);
    }

    @Override
    public String save(Receipt receipt) {
        throw new UnsupportedOperationException("not implemented yet " + clock);
    }

    @Override
    public Optional<Order> find(String id) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * Closes the database connection. After this, {@link #save} and {@link #find} throw
     * {@link StorageException}.
     */
    @Override
    public void close() {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
