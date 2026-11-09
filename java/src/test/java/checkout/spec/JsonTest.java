package checkout.spec;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;

/**
 * UNIT tests for the JSON helper: the field names and order of a serialized receipt are
 * part of the contract (snake_case, as in {@code spec/schema/receipt.schema.json}).
 */
class JsonTest {

    @Test
    void writesReceiptsWithSnakeCaseFieldNamesInDeclarationOrder() {
        assertThat(Json.toJson(Fixtures.ONE_APPLE_RECEIPT)).isEqualTo(
                "{\"lines\":[{\"sku\":\"APPLE\",\"name\":\"Apple\",\"quantity\":1,"
                        + "\"unit_price\":50,\"line_total\":50,\"line_discount\":0}],"
                        + "\"subtotal\":50,\"promotion_discount\":0,\"coupon_discount\":0,"
                        + "\"shipping\":499,\"tax\":4,\"total\":553}");
    }

    @Test
    void prettyPrintsOverSeveralLines() {
        assertThat(Json.toPrettyJson(Fixtures.ONE_APPLE_RECEIPT))
                .startsWith("{\n  \"lines\" : [ {")
                .contains("\n  \"total\" : 553\n}")
                .doesNotEndWith("\n");
    }

    @Test
    void parsesAndRoundTripsATree() {
        assertThat(Json.toTree(Fixtures.ONE_APPLE_RECEIPT).get("total").asInt()).isEqualTo(553);
        assertThat(Json.parse("{\"a\": [1, 2]}").get("a").size()).isEqualTo(2);
    }
}
