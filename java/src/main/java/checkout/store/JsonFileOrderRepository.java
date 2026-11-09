package checkout.store;

import checkout.Receipt;
import checkout.orders.Order;
import checkout.orders.OrderRepository;
import java.nio.file.Path;
import java.time.Clock;
import java.util.Optional;

/**
 * An order repository that writes one JSON file per order into a folder:
 * {@code <folder>/ORDER-1.json}, {@code <folder>/ORDER-2.json}, ... in the shape of
 * {@code spec/schema/order.schema.json}.
 *
 * <p>Ids continue where the files in the folder left off, so reopening the same folder
 * later keeps numbering correctly.
 */
public final class JsonFileOrderRepository implements OrderRepository {

    private final Path folder;
    private final Clock clock;

    /**
     * Opens (and if needed creates) the folder.
     *
     * @param folder where the {@code ORDER-<n>.json} files live
     * @param clock  the time source for {@code created_at}
     * @throws StorageException when the folder cannot be created
     */
    public JsonFileOrderRepository(Path folder, Clock clock) {
        this.folder = folder;
        this.clock = clock;
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * Opens the folder, using the real clock.
     *
     * @param folder where the {@code ORDER-<n>.json} files live
     */
    public JsonFileOrderRepository(Path folder) {
        this(folder, Clock.systemUTC());
    }

    @Override
    public String save(Receipt receipt) {
        throw new UnsupportedOperationException("not implemented yet " + folder + clock);
    }

    @Override
    public Optional<Order> find(String id) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
