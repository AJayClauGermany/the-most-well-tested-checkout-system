package checkout.spec;

import checkout.Basket;
import checkout.ErrorCode;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * One shared test case from {@code spec/test-cases/*.json}: a basket plus
 * <em>either</em> the expected totals <em>or</em> the expected error code.
 *
 * @param file          the JSON file the case came from, e.g. {@code "02-promotions.json"}
 * @param id            unique, readable name of the case, e.g. {@code "promo-coffee-three-bags-one-free"}
 * @param description   explains why the answer is what it is
 * @param basket        the checkout input
 * @param expected      the expected totals, or {@code null} when an error is expected
 * @param expectedError the expected error, or {@code null} when totals are expected
 */
public record SpecCase(
        String file,
        String id,
        String description,
        Basket basket,
        Totals expected,
        @JsonProperty("expected_error") ErrorCode expectedError) {

    /**
     * Returns a copy of this case with the file name filled in.
     *
     * @param fileName the JSON file the case came from
     * @return the same case, remembering its file
     */
    SpecCase withFile(String fileName) {
        return new SpecCase(fileName, id, description, basket, expected, expectedError);
    }
}
