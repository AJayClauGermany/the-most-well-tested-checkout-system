package checkout.spec;

import com.networknt.schema.Error;
import com.networknt.schema.Schema;
import com.networknt.schema.SchemaLocation;
import com.networknt.schema.SchemaRegistry;
import com.networknt.schema.SpecificationVersion;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import tools.jackson.databind.JsonNode;

/**
 * Test helper: JSON Schema validation with networknt's json-schema-validator.
 *
 * <p>The schemas in {@code spec/schema/} are the <em>contract</em> every language's output
 * must honour. A "contract test" checks a real receipt (or saved order, or CLI error)
 * against them. The schemas use JSON Schema draft 2020-12.
 *
 * <p>{@code order.schema.json} refers to {@code receipt.schema.json} by its relative name.
 * To make that work without a web server, every schema file is registered under a made-up
 * base IRI ({@code https://checkout.example/schema/<name>}), so a relative {@code $ref}
 * resolves to another registered schema.
 */
public final class Schemas {

    /** The validator's name for a schema. The address never goes on the network. */
    private static final String BASE_IRI = "https://checkout.example/schema/";

    private static final List<String> SCHEMA_FILES = List.of(
            "receipt.schema.json",
            "order.schema.json",
            "cli-error.schema.json",
            "basket.schema.json",
            "catalog.schema.json");

    private static final SchemaRegistry REGISTRY = SchemaRegistry.withDefaultDialect(
            SpecificationVersion.DRAFT_2020_12,
            builder -> builder.schemaLoader(loader -> loader.resourceLoaders(
                    resources -> resources.resources(readAllSchemas()))));

    private Schemas() {}

    /**
     * One problem the validator found.
     *
     * @param path    where in the document, e.g. {@code "/shipping"} ({@code ""} is the root)
     * @param keyword which schema rule failed, e.g. {@code "minimum"} or {@code "additionalProperties"}
     * @param message the validator's explanation
     */
    public record Problem(String path, String keyword, String message) {}

    /**
     * Validates a JSON tree against one of the shared schemas.
     *
     * @param schemaFile the file name in {@code spec/schema/}, e.g. {@code "receipt.schema.json"}
     * @param json       the document to check
     * @return an empty list when the document is valid, otherwise one problem per violation
     *         (an empty list is easy to assert with {@code isEmpty()} and prints nicely when it is not)
     */
    public static List<Problem> problems(String schemaFile, JsonNode json) {
        Schema schema = REGISTRY.getSchema(SchemaLocation.of(BASE_IRI + schemaFile));
        List<Error> errors = schema.validate(json);
        return errors.stream()
                .map(error -> new Problem(
                        error.getInstanceLocation().toString(), error.getKeyword(), error.getMessage()))
                .toList();
    }

    /**
     * Validates any object (a receipt, an order, ...) by serializing it to JSON first.
     *
     * @param schemaFile the file name in {@code spec/schema/}
     * @param value      the object to check
     * @return the problems, empty when valid
     */
    public static List<Problem> problems(String schemaFile, Object value) {
        return problems(schemaFile, Json.toTree(value));
    }

    /** Reads every schema file from the spec folder, keyed by the IRI the validator will ask for. */
    private static Map<String, String> readAllSchemas() {
        Path schemaDir = SpecFiles.specDir().resolve("schema");
        Map<String, String> schemas = new HashMap<>();
        for (String file : SCHEMA_FILES) {
            try {
                schemas.put(BASE_IRI + file, Files.readString(schemaDir.resolve(file)));
            } catch (IOException e) {
                throw new UncheckedIOException(e);
            }
        }
        return schemas;
    }
}
