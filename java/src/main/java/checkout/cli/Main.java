package checkout.cli;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

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
     * <p>Exit codes: 0 receipt printed, 1 checkout failed (error code printed),
     * 2 the input was not a basket at all (a message on standard error, nothing on standard output).
     *
     * @param args ignored
     * @throws IOException when standard input cannot be read
     */
    public static void main(String[] args) throws IOException {
        String input = new String(System.in.readAllBytes(), StandardCharsets.UTF_8);
        int exitCode;
        try {
            CommandResult result = CheckoutCommand.run(input, StandardCatalog.load());
            // "\n" rather than println, so the output is the same on every operating system.
            System.out.print(result.output() + "\n");
            exitCode = result.exitCode();
        } catch (RuntimeException e) {
            // Not a CheckoutException (those became exit code 1 above): the input was not a
            // basket at all, or the catalog could not be read.
            System.err.print("checkout: " + e.getMessage() + "\n");
            exitCode = CommandResult.BAD_INPUT;
        }
        System.out.flush();
        System.err.flush();
        System.exit(exitCode);
    }
}
