package checkout.cli;

import checkout.Basket;
import checkout.Catalog;

/**
 * The command line program's logic (rules section 8), kept free of {@code System.in} and
 * {@code System.out} so it can be unit tested in-process. {@link Main} is the thin wrapper
 * that talks to the terminal.
 */
public final class CheckoutCommand {

    private CheckoutCommand() {}

    /**
     * Runs one checkout from the JSON text of a basket.
     *
     * <ul>
     *   <li>a valid basket gives the receipt as one line of JSON and exit code 0,</li>
     *   <li>an invalid basket gives {@code {"error":"<CODE>"}} and exit code 1,</li>
     *   <li>text that is not a basket at all throws (Main turns that into exit code 2).</li>
     * </ul>
     *
     * @param input   the text read from standard input
     * @param catalog the products to price the basket with
     * @return what to print and the exit code
     * @throws IllegalArgumentException when the input is not JSON, or is JSON but not a basket
     */
    public static CommandResult run(String input, Catalog catalog) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * Parses the JSON text of a basket. Only the shape that matters is checked - an object
     * with an {@code items} list of objects, each with a text {@code sku} and a whole-number
     * {@code quantity} - so that the checkout can do the real validation (unknown products,
     * quantities out of range, ...).
     *
     * @param input the JSON text
     * @return the basket
     * @throws IllegalArgumentException when the text is not JSON, or the JSON is not a basket
     */
    public static Basket parseBasket(String input) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
