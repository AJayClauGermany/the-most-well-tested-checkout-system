package checkout.cli;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Catalog;
import checkout.Checkout;
import checkout.CheckoutException;
import java.util.ArrayList;
import java.util.List;
import tools.jackson.core.JacksonException;
import tools.jackson.databind.DeserializationFeature;
import tools.jackson.databind.JsonNode;
import tools.jackson.databind.json.JsonMapper;

/**
 * The command line program's logic (rules section 8), kept free of {@code System.in} and
 * {@code System.out} so it can be unit tested in-process. {@link Main} is the thin wrapper
 * that talks to the terminal.
 */
public final class CheckoutCommand {

    // FAIL_ON_TRAILING_TOKENS: '{"items":[]} garbage' is not one JSON value, so it is not a basket.
    private static final JsonMapper JSON = JsonMapper.builder()
            .enable(DeserializationFeature.FAIL_ON_TRAILING_TOKENS)
            .build();

    /** What the program prints when checkout fails: {@code {"error":"<CODE>"}} (cli-error.schema.json). */
    private record ErrorOutput(String error) {}

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
        Basket basket = parseBasket(input);
        try {
            String receipt = JSON.writeValueAsString(Checkout.checkout(basket, catalog));
            return new CommandResult(receipt, CommandResult.SUCCESS);
        } catch (CheckoutException e) {
            // Only CheckoutException becomes an error code. Anything else is a bug, and is
            // allowed to crash loudly instead of hiding behind a made-up error code.
            String output = JSON.writeValueAsString(new ErrorOutput(e.getCode().name()));
            return new CommandResult(output, CommandResult.CHECKOUT_FAILED);
        }
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
        JsonNode root = readJson(input);
        if (!root.isObject()) {
            throw notABasket("expected a JSON object");
        }
        JsonNode items = root.path("items");
        if (!items.isArray()) {
            throw notABasket("'items' must be a list");
        }
        List<BasketItem> basketItems = new ArrayList<>();
        for (JsonNode item : items) {
            // path() never returns null: a missing field is a "missing node", which is not a string.
            if (!item.path("sku").isString()) {
                throw notABasket("every item needs a text 'sku'");
            }
            // isInt(): a whole number that fits into a Java int. 1.5, "3" and 9999999999 do not.
            if (!item.path("quantity").isInt()) {
                throw notABasket("every item needs a whole-number 'quantity'");
            }
            basketItems.add(new BasketItem(item.path("sku").stringValue(), item.path("quantity").intValue()));
        }
        return new Basket(basketItems, optionalText(root, "coupon"), optionalText(root, "shipping"));
    }

    /** Reads JSON text, turning the parser's error into "not JSON". Empty input is not JSON either. */
    private static JsonNode readJson(String input) {
        JsonNode root;
        try {
            root = JSON.readTree(input);
        } catch (JacksonException e) {
            throw new IllegalArgumentException("The input is not JSON: " + e.getOriginalMessage(), e);
        }
        if (root.isMissingNode()) {
            throw new IllegalArgumentException("The input is not JSON: it is empty");
        }
        return root;
    }

    /** A field that may be missing or {@code null} (both give {@code null}), or else must be text. */
    private static String optionalText(JsonNode basket, String field) {
        JsonNode value = basket.path(field);
        if (value.isMissingNode() || value.isNull()) {
            return null;
        }
        if (!value.isString()) {
            throw notABasket("'" + field + "' must be text or null");
        }
        return value.stringValue();
    }

    private static IllegalArgumentException notABasket(String reason) {
        return new IllegalArgumentException("The input is not a basket: " + reason);
    }
}
