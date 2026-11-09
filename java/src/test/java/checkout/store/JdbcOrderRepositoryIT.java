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
import checkout.spec.Schemas;
import java.nio.file.Path;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneOffset;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

/**
 * INTEGRATION tests for the JDBC order repository, using H2 - a real SQL database that lives in
 * one file (or in memory), with no server to install.
 *
 * <p>Each test gets its own database file in a fresh {@code @TempDir} folder. The URL
 * {@code jdbc:h2:file:<path>} makes H2 create the file; {@code jdbc:h2:mem:} is a throw-away
 * database that disappears when the connection closes.
 */
class JdbcOrderRepositoryIT {

    private static final Clock FIXED_CLOCK = Clock.fixed(Instant.parse("2026-01-15T10:30:00Z"), ZoneOffset.UTC);

    @TempDir
    Path folder;

    private String fileUrl() {
        return "jdbc:h2:file:" + folder.resolve("orders").toAbsolutePath();
    }

    @Test
    void savesAndFindsAnOrder() {
        try (JdbcOrderRepository repository = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK)) {
            Receipt receipt = Checkout.checkout(
                    new Basket(List.of(new BasketItem("COFFEE", 3)), "SAVE10", null), TINY_CATALOG);

            String id = repository.save(receipt);

            assertThat(id).isEqualTo("ORDER-1");
            assertThat(repository.find(id)).contains(new Order("ORDER-1", "2026-01-15T10:30:00Z", receipt));
        }
    }

    @Test
    void givesBackAnOrderThatFollowsOrderSchemaJson() {
        try (JdbcOrderRepository repository = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK)) {
            String id = repository.save(ONE_APPLE_RECEIPT);

            assertThat(Schemas.problems("order.schema.json", repository.find(id).orElseThrow())).isEmpty();
        }
    }

    @Test
    void numbersOrdersSequentially() {
        try (JdbcOrderRepository repository = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK)) {
            List<String> ids = List.of(
                    repository.save(ONE_APPLE_RECEIPT),
                    repository.save(ONE_APPLE_RECEIPT),
                    repository.save(ONE_APPLE_RECEIPT));

            assertThat(ids).containsExactly("ORDER-1", "ORDER-2", "ORDER-3");
        }
    }

    @Test
    void continuesNumberingAfterTheDatabaseIsClosedAndOpenedAgain() {
        try (JdbcOrderRepository first = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK)) {
            first.save(ONE_APPLE_RECEIPT);
            first.save(ONE_APPLE_RECEIPT);
        }

        try (JdbcOrderRepository second = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK)) {
            assertThat(second.save(ONE_APPLE_RECEIPT)).isEqualTo("ORDER-3");
            assertThat(second.find("ORDER-2")).map(Order::id).contains("ORDER-2");
        }
    }

    @Test
    void reallyWritesToTheDatabaseFileCheckedWithASeparateConnection() throws SQLException {
        try (JdbcOrderRepository repository = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK)) {
            repository.save(ONE_APPLE_RECEIPT);
        }

        try (Connection connection = DriverManager.getConnection(fileUrl());
                Statement statement = connection.createStatement();
                ResultSet rows = statement.executeQuery("SELECT COUNT(*) FROM orders")) {
            rows.next();
            assertThat(rows.getInt(1)).isEqualTo(1);
        }
    }

    @Test
    void refusesToSaveOrFindOnceItIsClosedInsteadOfSilentlyDoingNothing() {
        JdbcOrderRepository repository = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK);
        repository.close();

        assertThatThrownBy(() -> repository.save(ONE_APPLE_RECEIPT)).isInstanceOf(StorageException.class);
        assertThatThrownBy(() -> repository.find("ORDER-1")).isInstanceOf(StorageException.class);
    }

    @Test
    void worksInMemoryToo() {
        // "jdbc:h2:mem:" is H2's name for a private throw-away database.
        try (JdbcOrderRepository repository = new JdbcOrderRepository("jdbc:h2:mem:", FIXED_CLOCK)) {
            assertThat(repository.save(ONE_APPLE_RECEIPT)).isEqualTo("ORDER-1");
            assertThat(repository.find("ORDER-1")).map(Order::receipt).contains(ONE_APPLE_RECEIPT);
        }
    }

    @ParameterizedTest(name = "answers 'not found' for ''{0}'' without throwing")
    @ValueSource(strings = {"ORDER-2", "ORDER-0", "ORDER-01", "ORDER-x", "1", "", "ORDER-1; DROP TABLE orders"})
    void answersNotFoundForIdsThatWereNeverHandedOut(String id) {
        try (JdbcOrderRepository repository = new JdbcOrderRepository(fileUrl(), FIXED_CLOCK)) {
            repository.save(ONE_APPLE_RECEIPT);

            assertThat(repository.find(id)).isEmpty();
            assertThat(repository.find("ORDER-1")).as("the table is still there").isPresent();
        }
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

        try (JdbcOrderRepository repository = new JdbcOrderRepository(fileUrl(), adjustable)) {
            String first = repository.save(ONE_APPLE_RECEIPT);
            now[0] = Instant.parse("2026-01-16T09:00:00Z");
            String second = repository.save(ONE_APPLE_RECEIPT);

            assertThat(repository.find(first)).map(Order::createdAt).contains("2026-01-15T10:30:00Z");
            assertThat(repository.find(second)).map(Order::createdAt).contains("2026-01-16T09:00:00Z");
        }
    }

    @Test
    void reportsADatabaseThatCannotBeOpenedAsAStorageException() {
        assertThatThrownBy(() -> new JdbcOrderRepository("jdbc:nonsense:nowhere", FIXED_CLOCK))
                .isInstanceOf(StorageException.class)
                .hasMessageContaining("jdbc:nonsense:nowhere");
    }
}
