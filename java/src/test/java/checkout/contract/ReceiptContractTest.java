package checkout.contract;

import static checkout.spec.Fixtures.TINY_CATALOG;
import static checkout.spec.Fixtures.basketOf;
import static checkout.spec.Fixtures.item;
import static org.assertj.core.api.Assertions.assertThat;

import checkout.Catalog;
import checkout.Checkout;
import checkout.Receipt;
import checkout.cli.CheckoutCommand;
import checkout.cli.CommandResult;
import checkout.spec.Json;
import checkout.spec.Schemas;
import checkout.spec.SpecCase;
import checkout.spec.SpecFiles;
import java.util.stream.Stream;
import org.junit.jupiter.api.Named;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.MethodSource;

/**
 * CONTRACT tests: is the output even the right shape?
 *
 * <p>{@code spec/schema/receipt.schema.json} is the contract every language's receipt must
 * honour, and {@code cli-error.schema.json} the one for the command's error output. The
 * shared e2e runner only compares the six totals; this is where the shape of the whole
 * receipt (lines included, snake_case names, no extra fields) is held to the contract, using
 * networknt's JSON Schema validator.
 */
class ReceiptContractTest {

    private static final Catalog CATALOG = SpecFiles.loadCatalog();

    static Stream<Named<SpecCase>> successfulCases() {
        return SpecFiles.loadAllCases().stream()
                .filter(specCase -> specCase.expected() != null)
                .map(specCase -> Named.of(specCase.id(), specCase));
    }

    static Stream<Named<SpecCase>> failingCases() {
        return SpecFiles.loadAllCases().stream()
                .filter(specCase -> specCase.expectedError() != null)
                .map(specCase -> Named.of(specCase.id(), specCase));
    }

    @ParameterizedTest(name = "the receipt for {0} follows receipt.schema.json")
    @MethodSource("successfulCases")
    void everyReceiptOfTheSharedCasesFollowsTheSchema(SpecCase specCase) {
        Receipt receipt = Checkout.checkout(specCase.basket(), CATALOG);

        assertThat(Schemas.problems("receipt.schema.json", receipt)).isEmpty();
    }

    @ParameterizedTest(name = "the command's output for {0} follows cli-error.schema.json")
    @MethodSource("failingCases")
    void everyErrorTheCommandPrintsFollowsTheSchema(SpecCase specCase) {
        CommandResult result = CheckoutCommand.run(Json.toJson(specCase.basket()), CATALOG);

        assertThat(result.exitCode()).isEqualTo(1);
        assertThat(Schemas.problems("cli-error.schema.json", Json.parse(result.output()))).isEmpty();
    }

    @Test
    void aReceiptFromAMadeUpCatalogFollowsTheSchemaToo() {
        Receipt receipt = Checkout.checkout(basketOf(item("TV", 1), item("COFFEE", 6)), TINY_CATALOG);

        assertThat(Schemas.problems("receipt.schema.json", receipt)).isEmpty();
    }

    @Test
    void theCommandsReceiptOutputFollowsTheSchemaAsText() {
        // The command serializes the receipt itself; check the text it prints, not the Java object.
        CommandResult result = CheckoutCommand.run("{\"items\":[{\"sku\":\"COFFEE\",\"quantity\":3}]}", CATALOG);

        assertThat(Schemas.problems("receipt.schema.json", Json.parse(result.output()))).isEmpty();
    }
}
