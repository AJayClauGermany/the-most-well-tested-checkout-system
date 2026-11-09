package checkout.cli;

/**
 * What the command line program prints and the exit code it ends with.
 *
 * @param output   one line of JSON: the receipt, or {@code {"error":"<CODE>"}}
 * @param exitCode {@code 0} for a receipt, {@code 1} for an error
 */
public record CommandResult(String output, int exitCode) {

    /** Exit code of a successful checkout. */
    public static final int SUCCESS = 0;
    /** Exit code of a checkout that failed with one of the six error codes. */
    public static final int CHECKOUT_FAILED = 1;
    /** Exit code when the input was not a basket at all (used by {@link Main} only). */
    public static final int BAD_INPUT = 2;
}
