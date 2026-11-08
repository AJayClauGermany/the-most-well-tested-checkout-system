package checkout;

/**
 * Thrown by {@link Checkout#checkout} when a basket cannot be checked out.
 *
 * <p>In Java, the usual way to say "this went wrong" is to throw an exception.
 * The {@link ErrorCode} says <em>what</em> went wrong, so tests can check it:
 *
 * <pre>{@code
 * CheckoutException error = assertThrows(CheckoutException.class, () -> Checkout.checkout(basket, catalog));
 * assertEquals(ErrorCode.EMPTY_BASKET, error.getCode());
 * }</pre>
 *
 * <p>It extends {@link RuntimeException} (an "unchecked" exception), so callers
 * are not forced to write {@code try/catch} or {@code throws} everywhere.
 */
public class CheckoutException extends RuntimeException {

    // Exceptions can be "serialized" (turned into bytes). Java wants a version
    // number for that, and our compiler settings turn a missing one into an error.
    private static final long serialVersionUID = 1L;

    private final ErrorCode code;

    /**
     * Creates an exception for the given error code.
     *
     * @param code what went wrong
     */
    public CheckoutException(ErrorCode code) {
        super(code.name());
        this.code = code;
    }

    /**
     * Returns what went wrong.
     *
     * @return the error code from rules section 9
     */
    public ErrorCode getCode() {
        return code;
    }
}
