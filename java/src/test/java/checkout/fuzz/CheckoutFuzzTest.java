package checkout.fuzz;

import static org.assertj.core.api.Assertions.assertThat;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Catalog;
import checkout.Checkout;
import checkout.CheckoutException;
import checkout.ErrorCode;
import checkout.cli.CheckoutCommand;
import checkout.cli.CommandResult;
import checkout.spec.Json;
import checkout.spec.Schemas;
import checkout.spec.SpecFiles;
import com.code_intelligence.jazzer.api.FuzzedDataProvider;
import com.code_intelligence.jazzer.junit.FuzzTest;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import org.junit.jupiter.api.Tag;

/**
 * FUZZ tests with Jazzer: property P10 ("validation is total") attacked with random input.
 *
 * <p>Fuzzing throws random bytes at the code and watches for crashes. Jazzer's {@code @FuzzTest}
 * runs inside JUnit in two modes:
 * <ul>
 *   <li>by default (also in CI) as a <em>regression test</em>: it replays the seed inputs in
 *       {@code src/test/resources/checkout/fuzz/CheckoutFuzzTestInputs/} plus an empty input;</li>
 *   <li>with the environment variable {@code JAZZER_FUZZ=1} as a real fuzzer that mutates the
 *       inputs for {@code maxDuration}, guided by code coverage, and saves any crashing input
 *       as a new seed, so the crash becomes a regression test.</li>
 * </ul>
 * Property-based testing (jqwik) and fuzzing look alike; the difference is that a fuzzer is
 * guided by which branches of the code the input reaches, which makes it good at finding
 * parsing bugs.
 */
@Tag("fuzz")
class CheckoutFuzzTest {

    private static final Catalog CATALOG = SpecFiles.loadCatalog();

    /** Random bytes as stdin: the command must either answer properly or reject the input, never crash. */
    @FuzzTest(maxDuration = "20s")
    void theCommandNeverCrashesOnRandomInput(byte[] input) {
        String text = new String(input, StandardCharsets.UTF_8);
        try {
            CommandResult result = CheckoutCommand.run(text, CATALOG);
            assertThat(result.exitCode()).isIn(0, 1);
            Object parsed = Json.parse(result.output());
            assertThat(parsed).isNotNull();
            if (result.exitCode() == 0) {
                assertThat(Schemas.problems("receipt.schema.json", Json.parse(result.output()))).isEmpty();
            } else {
                assertThat(Schemas.problems("cli-error.schema.json", Json.parse(result.output()))).isEmpty();
            }
        } catch (IllegalArgumentException e) {
            // "Not a basket at all" - the one other allowed outcome (exit code 2 in Main).
            assertThat(e.getMessage()).startsWith("The input is not ");
        }
    }

    /** A structurally valid but otherwise random basket: checkout throws nothing but CheckoutException. */
    @FuzzTest(maxDuration = "20s")
    void checkoutOnlyEverThrowsCheckoutException(FuzzedDataProvider data) {
        List<BasketItem> items = new ArrayList<>();
        int lines = data.consumeInt(0, 6);
        for (int i = 0; i < lines; i++) {
            String sku = data.consumeBoolean() ? data.pickValue(skus()) : data.consumeString(8);
            items.add(new BasketItem(sku, data.consumeInt()));
        }
        String coupon = data.consumeBoolean() ? null : data.consumeString(10);
        String shipping = data.consumeBoolean() ? null : data.consumeString(10);
        Basket basket = new Basket(items, coupon, shipping);

        try {
            assertThat(Schemas.problems("receipt.schema.json", Checkout.checkout(basket, CATALOG))).isEmpty();
        } catch (CheckoutException e) {
            assertThat(e.getCode()).isIn((Object[]) ErrorCode.values());
        }
    }

    private static List<String> skus() {
        return CATALOG.products().stream().map(product -> product.sku()).toList();
    }
}
