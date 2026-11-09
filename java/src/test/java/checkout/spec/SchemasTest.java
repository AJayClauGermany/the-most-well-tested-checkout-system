package checkout.spec;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.groups.Tuple.tuple;

import java.io.IOException;
import java.nio.file.Files;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests for the schema helper, using the hand-checked files in {@code spec/examples/}.
 *
 * <p>"Testing the test": a validator that accepted everything would make every contract test
 * pass for the wrong reason, so we also check that it rejects a broken document.
 */
class SchemasTest {

    @Test
    void acceptsTheExampleReceipt() throws IOException {
        String json = Files.readString(SpecFiles.specDir().resolve("examples/receipt.json"));

        assertThat(Schemas.problems("receipt.schema.json", Json.parse(json))).isEmpty();
    }

    @Test
    void acceptsTheExampleOrderFollowingTheRelativeRefToTheReceiptSchema() throws IOException {
        String json = Files.readString(SpecFiles.specDir().resolve("examples/order.json"));

        assertThat(Schemas.problems("order.schema.json", Json.parse(json))).isEmpty();
    }

    @Test
    void acceptsTheExampleCliError() throws IOException {
        String json = Files.readString(SpecFiles.specDir().resolve("examples/cli-error.json"));

        assertThat(Schemas.problems("cli-error.schema.json", Json.parse(json))).isEmpty();
    }

    @Test
    void rejectsABrokenReceiptWithOneProblemPerViolation() {
        // An extra field, a shipping cost that is not in the table, and a negative tax.
        String broken = """
                {"lines": [{"sku": "APPLE", "name": "Apple", "quantity": 1,
                            "unit_price": 50, "line_total": 50, "line_discount": 0}],
                 "subtotal": 50, "promotion_discount": 0, "coupon_discount": 0,
                 "shipping": 123, "tax": -1, "total": 172, "extra": true}
                """;

        assertThat(Schemas.problems("receipt.schema.json", Json.parse(broken)))
                .extracting(Schemas.Problem::path, Schemas.Problem::keyword)
                .containsExactlyInAnyOrder(
                        tuple("", "additionalProperties"),
                        tuple("/shipping", "enum"),
                        tuple("/tax", "minimum"));
    }

    @Test
    void rejectsAnOrderWhoseReceiptIsBrokenSoTheRefReallyIsFollowed() {
        String broken = """
                {"id": "ORDER-1", "created_at": "2026-01-15T10:30:00Z",
                 "receipt": {"lines": [], "subtotal": 0, "promotion_discount": 0,
                             "coupon_discount": 0, "shipping": 0, "tax": 0, "total": 0}}
                """;

        assertThat(Schemas.problems("order.schema.json", Json.parse(broken)))
                .extracting(Schemas.Problem::path)
                .containsExactly("/receipt/lines");
    }

    @Test
    void everySharedTestCaseBasketIsAValidBasketSoTheSpecIsConsistentWithItself() {
        for (SpecCase specCase : SpecFiles.loadAllCases()) {
            assertThat(Schemas.problems("basket.schema.json", specCase.basket())).as(specCase.id()).isEmpty();
        }
    }
}
