package checkout.store;

/**
 * Thrown when an order store cannot read or write (a file cannot be written, the database
 * connection is closed, ...).
 *
 * <p>It wraps the checked {@code IOException} / {@code SQLException} that Java's file and
 * database APIs throw, so callers of {@code save} and {@code find} are not forced to catch
 * exceptions that the in-memory repository can never throw.
 */
public class StorageException extends RuntimeException {

    private static final long serialVersionUID = 1L;

    /**
     * Creates the exception.
     *
     * @param message what went wrong
     * @param cause   the underlying file or database error
     */
    public StorageException(String message, Throwable cause) {
        super(message, cause);
    }
}
