package checkout.store;

import checkout.Receipt;
import checkout.orders.Order;
import checkout.orders.OrderIds;
import checkout.orders.OrderRepository;
import checkout.orders.Timestamps;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.time.Clock;
import java.util.Optional;
import java.util.OptionalInt;
import java.util.stream.Stream;
import tools.jackson.core.JacksonException;
import tools.jackson.databind.json.JsonMapper;

/**
 * An order repository that writes one JSON file per order into a folder:
 * {@code <folder>/ORDER-1.json}, {@code <folder>/ORDER-2.json}, ... in the shape of
 * {@code spec/schema/order.schema.json}.
 *
 * <p>Ids continue where the files in the folder left off, so reopening the same folder
 * later keeps numbering correctly.
 */
public final class JsonFileOrderRepository implements OrderRepository {

    private static final JsonMapper JSON = JsonMapper.builder().build();

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
        try {
            Files.createDirectories(folder);
        } catch (IOException e) {
            throw new StorageException("Could not create the order folder " + folder, e);
        }
    }

    /**
     * Opens the folder, using the real clock.
     *
     * @param folder where the {@code ORDER-<n>.json} files live
     */
    public JsonFileOrderRepository(Path folder) {
        this(folder, Clock.systemUTC());
    }

    // "synchronized": two threads saving at once must not pick the same number.
    @Override
    public synchronized String save(Receipt receipt) {
        Order order = new Order(OrderIds.of(highestOrderNumber() + 1), Timestamps.now(clock), receipt);
        Path path = pathOf(order.id());
        String json = JSON.writerWithDefaultPrettyPrinter().writeValueAsString(order) + "\n";
        try {
            // CREATE_NEW: never overwrite an existing order, even if another program wrote it meanwhile.
            Files.writeString(path, json, StandardOpenOption.CREATE_NEW);
        } catch (IOException e) {
            throw new StorageException("Could not write " + path, e);
        }
        return order.id();
    }

    @Override
    public Optional<Order> find(String id) {
        // Only accept real order ids: "../secret" must not read files outside the folder.
        if (OrderIds.numberOf(id).isEmpty()) {
            return Optional.empty();
        }
        Path path = pathOf(id);
        if (!Files.isRegularFile(path)) {
            return Optional.empty();
        }
        try {
            return Optional.of(JSON.readValue(path.toFile(), Order.class));
        } catch (JacksonException e) {
            throw new StorageException("Could not read " + path, e);
        }
    }

    private Path pathOf(String id) {
        return folder.resolve(id + ".json");
    }

    /** The highest order number among the {@code ORDER-<n>.json} files in the folder, or 0 when there are none. */
    private int highestOrderNumber() {
        try (Stream<Path> files = Files.list(folder)) {
            return files
                    .filter(Files::isRegularFile)
                    .map(path -> path.getFileName().toString())
                    .filter(name -> name.endsWith(".json"))
                    .map(name -> OrderIds.numberOf(name.substring(0, name.length() - ".json".length())))
                    .filter(OptionalInt::isPresent)
                    .mapToInt(OptionalInt::getAsInt)
                    .max()
                    .orElse(0);
        } catch (IOException e) {
            throw new StorageException("Could not list the order folder " + folder, e);
        }
    }
}
