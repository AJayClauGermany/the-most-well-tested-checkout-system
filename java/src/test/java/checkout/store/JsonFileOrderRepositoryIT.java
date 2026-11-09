package checkout.store;

import static checkout.spec.Fixtures.ONE_APPLE_RECEIPT;
import static checkout.spec.Fixtures.TINY_CATALOG;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Checkout;
import checkout.Receipt;
import checkout.orders.Order;
import checkout.spec.Json;
import checkout.spec.Schemas;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneOffset;
import java.util.List;
import java.util.stream.Stream;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

/**
 * INTEGRATION tests for the JSON-file order repository: real files in a real temporary folder.
 *
 * <p>Unlike a unit test, this touches the file system, so each test gets its own fresh folder
 * from JUnit's {@code @TempDir}, deleted afterwards. The stored files must follow
 * {@code spec/schema/order.schema.json} - a contract test inside an integration test.
 * The class name ends in {@code IT}, so Failsafe runs it during {@code mvn verify}.
 */
class JsonFileOrderRepositoryIT {

    private static final Clock FIXED_CLOCK = Clock.fixed(Instant.parse("2026-01-15T10:30:00Z"), ZoneOffset.UTC);

    @TempDir
    Path folder;

    @Test
    void writesOrder1JsonIntoTheFolderOnTheFirstSave() {
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);

        String id = repository.save(ONE_APPLE_RECEIPT);

        assertThat(id).isEqualTo("ORDER-1");
        assertThat(folder.resolve("ORDER-1.json")).isRegularFile();
    }

    @Test
    void writesAFileThatFollowsOrderSchemaJsonExactly() throws IOException {
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);
        repository.save(ONE_APPLE_RECEIPT);

        String stored = Files.readString(folder.resolve("ORDER-1.json"));

        assertThat(Schemas.problems("order.schema.json", Json.parse(stored))).isEmpty();
        assertThat(Json.parse(stored))
                .isEqualTo(Json.toTree(new Order("ORDER-1", "2026-01-15T10:30:00Z", ONE_APPLE_RECEIPT)));
    }

    @Test
    void writesHumanReadableJsonIndentedAndEndingInANewline() throws IOException {
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);
        repository.save(ONE_APPLE_RECEIPT);

        String text = Files.readString(folder.resolve("ORDER-1.json"));

        assertThat(text).contains("\n  ").endsWith("\n");
    }

    @Test
    void readsASavedOrderBackExactlyAsSaved() {
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);
        Receipt receipt = Checkout.checkout(
                new Basket(List.of(new BasketItem("COFFEE", 3)), "SAVE10", null), TINY_CATALOG);

        String id = repository.save(receipt);

        assertThat(repository.find(id)).contains(new Order("ORDER-1", "2026-01-15T10:30:00Z", receipt));
    }

    @Test
    void numbersOrdersSequentially() throws IOException {
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);

        List<String> ids = List.of(
                repository.save(ONE_APPLE_RECEIPT),
                repository.save(ONE_APPLE_RECEIPT),
                repository.save(ONE_APPLE_RECEIPT));

        assertThat(ids).containsExactly("ORDER-1", "ORDER-2", "ORDER-3");
        assertThat(fileNames()).containsExactly("ORDER-1.json", "ORDER-2.json", "ORDER-3.json");
    }

    @Test
    void continuesNumberingWhenTheSameFolderIsOpenedAgainLater() {
        JsonFileOrderRepository first = new JsonFileOrderRepository(folder, FIXED_CLOCK);
        first.save(ONE_APPLE_RECEIPT);
        first.save(ONE_APPLE_RECEIPT);

        JsonFileOrderRepository second = new JsonFileOrderRepository(folder, FIXED_CLOCK);

        assertThat(second.save(ONE_APPLE_RECEIPT)).isEqualTo("ORDER-3");
        assertThat(second.find("ORDER-1")).map(Order::id).contains("ORDER-1");
    }

    @Test
    void numbersAfterTheHighestExistingOrderEvenWithGapsAndTwoDigitNumbers() throws IOException {
        writeOrderFile("ORDER-3");
        writeOrderFile("ORDER-12");
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);

        // Not ORDER-4: alphabetically "ORDER-3" > "ORDER-12", numerically 12 > 3.
        assertThat(repository.save(ONE_APPLE_RECEIPT)).isEqualTo("ORDER-13");
    }

    @Test
    void ignoresFilesInTheFolderThatAreNotOrders() throws IOException {
        Files.writeString(folder.resolve("README.txt"), "not an order");
        Files.writeString(folder.resolve("ORDER-9.json.bak"), "not an order either");
        Files.writeString(folder.resolve("ORDER-9"), "no .json extension, so not an order");
        // ".yaml" is exactly as long as ".json": proves the extension is really checked
        // (mutation testing found that no test did).
        Files.writeString(folder.resolve("ORDER-7.yaml"), "another extension, so not an order");
        Files.writeString(folder.resolve("ORDER-09.json"), "a leading zero is not an order id");
        Files.createDirectory(folder.resolve("ORDER-8.json"));
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);

        assertThat(repository.save(ONE_APPLE_RECEIPT)).isEqualTo("ORDER-1");
    }

    @Test
    void createsTheFolderIfItDoesNotExistYet() {
        Path nested = folder.resolve("orders/2026");
        JsonFileOrderRepository repository = new JsonFileOrderRepository(nested, FIXED_CLOCK);

        repository.save(ONE_APPLE_RECEIPT);

        assertThat(nested.resolve("ORDER-1.json")).isRegularFile();
    }

    @Test
    void answersNotFoundForIdsThatHaveNoFileWithoutThrowing() {
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);
        repository.save(ONE_APPLE_RECEIPT);

        assertThat(repository.find("ORDER-2")).isEmpty();
        assertThat(repository.find("ORDER-0")).isEmpty();
        assertThat(repository.find("nonsense")).isEmpty();
    }

    @Test
    void doesNotLetAnIdEscapeTheFolderNoPathTricks() throws IOException {
        Files.writeString(folder.resolve("secret.json"), "{\"id\": \"x\", \"created_at\": \"x\", \"receipt\": {}}");
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder.resolve("orders"), FIXED_CLOCK);

        assertThat(repository.find("../secret")).isEmpty();
        assertThat(repository.find("../secret.json")).isEmpty();
    }

    @Test
    void stampsEachOrderWithTheClocksTimeAtTheMomentOfSaving() {
        Instant[] now = {Instant.parse("2026-01-15T10:30:00Z")};
        Clock adjustable = new Clock() {
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
                return now[0];
            }
        };
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, adjustable);

        String first = repository.save(ONE_APPLE_RECEIPT);
        now[0] = Instant.parse("2026-01-16T09:00:00Z");
        String second = repository.save(ONE_APPLE_RECEIPT);

        assertThat(repository.find(first)).map(Order::createdAt).contains("2026-01-15T10:30:00Z");
        assertThat(repository.find(second)).map(Order::createdAt).contains("2026-01-16T09:00:00Z");
    }

    @Test
    void usesTheRealClockWhenNoneIsGiven() {
        Instant before = Instant.now().minusSeconds(1);
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder);

        String id = repository.save(ONE_APPLE_RECEIPT);

        Instant createdAt = Instant.parse(repository.find(id).orElseThrow().createdAt());
        assertThat(createdAt).isBetween(before, Instant.now().plusSeconds(1));
    }

    @Test
    void reportsAFolderThatCannotBeCreatedAsAStorageException() throws IOException {
        Path file = folder.resolve("a-file-not-a-folder");
        Files.writeString(file, "in the way");

        assertThatThrownBy(() -> new JsonFileOrderRepository(file.resolve("orders"), FIXED_CLOCK))
                .isInstanceOf(StorageException.class)
                .hasMessageContaining("orders");
    }

    @Test
    void reportsAFolderThatDisappearedAsAStorageException() throws IOException {
        Path orders = folder.resolve("orders");
        JsonFileOrderRepository repository = new JsonFileOrderRepository(orders, FIXED_CLOCK);
        Files.delete(orders);

        assertThatThrownBy(() -> repository.save(ONE_APPLE_RECEIPT))
                .isInstanceOf(StorageException.class)
                .hasMessageContaining("Could not list the order folder");
    }

    @Test
    void neverOverwritesWhatIsInTheWayOfTheNextOrderFile() throws IOException {
        // A folder called ORDER-1.json is not an order (so the next id is ORDER-1),
        // but it is in the way: saving must fail loudly, not replace or lose anything.
        Files.createDirectory(folder.resolve("ORDER-1.json"));
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);

        assertThatThrownBy(() -> repository.save(ONE_APPLE_RECEIPT))
                .isInstanceOf(StorageException.class)
                .hasMessageContaining("Could not write")
                .hasMessageContaining("ORDER-1.json");
        assertThat(folder.resolve("ORDER-1.json")).isDirectory();
    }

    @Test
    void reportsAnUnreadableOrderFileAsAStorageException() throws IOException {
        Files.writeString(folder.resolve("ORDER-1.json"), "this is not json");
        JsonFileOrderRepository repository = new JsonFileOrderRepository(folder, FIXED_CLOCK);

        assertThatThrownBy(() -> repository.find("ORDER-1"))
                .isInstanceOf(StorageException.class)
                .hasMessageContaining("ORDER-1.json");
    }

    private void writeOrderFile(String id) throws IOException {
        Files.writeString(folder.resolve(id + ".json"),
                Json.toJson(new Order(id, "2026-01-15T10:30:00Z", ONE_APPLE_RECEIPT)));
    }

    private List<String> fileNames() throws IOException {
        try (Stream<Path> files = Files.list(folder)) {
            return files.map(path -> path.getFileName().toString()).sorted().toList();
        }
    }
}
