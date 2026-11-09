package checkout.spec;

import tools.jackson.databind.JsonNode;
import tools.jackson.databind.json.JsonMapper;

/**
 * Test helper: JSON in and out, with one shared Jackson mapper.
 *
 * <p>The receipt and order records carry {@code @JsonProperty("snake_case")} annotations,
 * so {@link #toJson} produces exactly the field names the spec schemas expect.
 */
public final class Json {

    /** One mapper is enough for everything; Jackson mappers are thread-safe. */
    public static final JsonMapper MAPPER = JsonMapper.builder().build();

    private Json() {}

    /**
     * Serializes any object (a receipt, an order, ...) to compact JSON on one line.
     *
     * @param value the object
     * @return its JSON text
     */
    public static String toJson(Object value) {
        return MAPPER.writeValueAsString(value);
    }

    /**
     * Serializes any object to indented JSON, ending without a newline.
     *
     * @param value the object
     * @return its pretty-printed JSON text
     */
    public static String toPrettyJson(Object value) {
        return MAPPER.writerWithDefaultPrettyPrinter().writeValueAsString(value);
    }

    /**
     * Parses JSON text into a tree ({@link JsonNode}), the form the schema validator wants.
     *
     * @param json the JSON text
     * @return the tree
     */
    public static JsonNode parse(String json) {
        return MAPPER.readTree(json);
    }

    /**
     * Turns any object into a JSON tree, by serializing it and parsing it back.
     *
     * @param value the object
     * @return the tree
     */
    public static JsonNode toTree(Object value) {
        return parse(toJson(value));
    }
}
