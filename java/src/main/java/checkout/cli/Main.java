package checkout.cli;

/**
 * The command line program (rules section 8): basket JSON in on standard input, receipt
 * JSON out on standard output.
 *
 * <pre>
 * echo '{"items":[{"sku":"APPLE","quantity":1}]}' | java -jar target/checkout.jar
 * </pre>
 *
 * <p>All the logic is in {@link CheckoutCommand}; this class only connects it to the
 * terminal. That is why it is left out of the coverage and mutation reports (they only see
 * code running inside the test process) and tested end-to-end instead, in
 * {@code checkout.e2e.CliIT} and by the shared runner {@code e2e/run.sh}.
 */
public final class Main {

    private Main() {}

    /**
     * Reads one basket from standard input and prints the result.
     *
     * @param args ignored
     */
    public static void main(String[] args) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
