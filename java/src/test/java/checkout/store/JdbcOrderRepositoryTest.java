package checkout.store;

import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.contains;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.Statement;
import java.time.Clock;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests of the JDBC repository's error handling, with a MOCK database connection.
 *
 * <p>The integration test {@code JdbcOrderRepositoryIT} uses a real H2 database. But some
 * failures - the disk filling up while the table is created, an error while closing - are
 * almost impossible to cause with a real database. A Mockito mock of {@link Connection} can
 * simply be told to fail, so we can check that such a failure is reported as a
 * {@link StorageException} with the original error attached, instead of being lost.
 */
@Tag("mock")
class JdbcOrderRepositoryTest {

    private final Connection connection = mock(Connection.class);

    @Test
    void reportsAFailureToCreateTheTableAsAStorageException() throws SQLException {
        SQLException diskFull = new SQLException("disk full");
        when(connection.createStatement()).thenThrow(diskFull);

        assertThatThrownBy(() -> new JdbcOrderRepository(connection, Clock.systemUTC()))
                .isInstanceOf(StorageException.class)
                .hasMessage("Could not create the orders table")
                .hasCause(diskFull);
    }

    @Test
    void closesTheConnectionItWasGiven() throws SQLException {
        when(connection.createStatement()).thenReturn(mock(Statement.class));
        JdbcOrderRepository repository = new JdbcOrderRepository(connection, Clock.systemUTC());

        repository.close();

        verify(connection).close();
    }

    @Test
    void reportsAFailureWhileClosingAsAStorageException() throws SQLException {
        when(connection.createStatement()).thenReturn(mock(Statement.class));
        SQLException broken = new SQLException("connection already broken");
        doThrow(broken).when(connection).close();
        JdbcOrderRepository repository = new JdbcOrderRepository(connection, Clock.systemUTC());

        assertThatThrownBy(repository::close)
                .isInstanceOf(StorageException.class)
                .hasMessage("Could not close the database")
                .hasCause(broken);
    }

    @Test
    void reportsAFailingQueryAsAStorageExceptionAndStillClosesTheStatement() throws SQLException {
        when(connection.createStatement()).thenReturn(mock(Statement.class));
        PreparedStatement select = mock(PreparedStatement.class);
        SQLException timeout = new SQLException("query timed out");
        when(connection.prepareStatement(anyString())).thenReturn(select);
        when(select.executeQuery()).thenThrow(timeout);
        JdbcOrderRepository repository = new JdbcOrderRepository(connection, Clock.systemUTC());

        assertThatThrownBy(() -> repository.find("ORDER-1"))
                .isInstanceOf(StorageException.class)
                .hasMessage("Could not look up ORDER-1")
                .hasCause(timeout);
        verify(select).close();
    }

    @Test
    void createsTheOrdersTableWhenItStarts() throws SQLException {
        Statement statement = mock(Statement.class);
        when(connection.createStatement()).thenReturn(statement);

        new JdbcOrderRepository(connection, Clock.systemUTC()).close();

        verify(statement).execute(contains("CREATE TABLE IF NOT EXISTS orders"));
        verify(statement).close();
    }
}
