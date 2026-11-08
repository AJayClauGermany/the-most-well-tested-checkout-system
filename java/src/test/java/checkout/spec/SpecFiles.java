package checkout.spec;

import checkout.Catalog;
import checkout.Product;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Stream;
import tools.jackson.databind.DeserializationFeature;
import tools.jackson.databind.json.JsonMapper;

/**
 * Loads the shared, language-agnostic spec files from the repository's {@code spec/} folder.
 *
 * <p>This is test-support code (it lives under {@code src/test}), so it is never part of
 * the real program. The files are read straight from {@code spec/} - they are not copied
 * into the Java project - so every language is tested against exactly the same data.
 */
public final class SpecFiles {

    // Jackson's JsonMapper reads JSON text into Java records. One is enough for everything.
    // We make it strict about unknown fields, so a typo in a spec file (e.g. "quantiy")
    // makes the tests fail loudly instead of being silently ignored.
    private static final JsonMapper JSON = JsonMapper.builder()
            .enable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
            .build();

    private SpecFiles() {}

    /**
     * Finds the repository's {@code spec/} folder.
     *
     * <p>Maven runs the tests inside {@code java/}, but an IDE might run them from the
     * repository root. So we start at the current folder and walk up, one parent at a
     * time, until we find a folder that contains {@code spec/checkout-rules.md}.
     *
     * @return the path of the {@code spec/} folder
     */
    public static Path specDir() {
        Path start = Path.of("").toAbsolutePath();
        for (Path dir = start; dir != null; dir = dir.getParent()) {
            Path candidate = dir.resolve("spec");
            if (Files.isRegularFile(candidate.resolve("checkout-rules.md"))) {
                return candidate;
            }
        }
        throw new IllegalStateException("Could not find the spec/ folder in " + start + " or any parent folder");
    }

    /**
     * Loads the standard catalog, {@code spec/catalog.json}.
     *
     * @return the six standard products
     */
    public static Catalog loadCatalog() {
        CatalogFile file = JSON.readValue(specDir().resolve("catalog.json").toFile(), CatalogFile.class);
        return new Catalog(file.products());
    }

    /**
     * Loads every case from every {@code spec/test-cases/*.json} file,
     * sorted by file name and then in the order they appear in the file.
     *
     * @return all shared spec cases
     */
    public static List<SpecCase> loadAllCases() {
        List<SpecCase> allCases = new ArrayList<>();
        for (Path path : testCaseFiles()) {
            String fileName = path.getFileName().toString();
            TestCaseFile file = JSON.readValue(path.toFile(), TestCaseFile.class);
            for (SpecCase specCase : file.cases()) {
                // The JSON does not contain the file name, so we add it here.
                allCases.add(specCase.withFile(fileName));
            }
        }
        return allCases;
    }

    /** Returns the paths of all {@code spec/test-cases/*.json} files, sorted by name. */
    private static List<Path> testCaseFiles() {
        try (Stream<Path> files = Files.list(specDir().resolve("test-cases"))) {
            return files
                    .filter(path -> path.getFileName().toString().endsWith(".json"))
                    .sorted()
                    .toList();
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
    }

    /** The shape of {@code catalog.json}. */
    private record CatalogFile(String description, List<Product> products) {}

    /**
     * The shape of a {@code test-cases/*.json} file. The {@code "$schema"} field only helps
     * editors check the file against the JSON schema, so we ignore it.
     */
    @JsonIgnoreProperties({"$schema"})
    private record TestCaseFile(String description, List<SpecCase> cases) {}
}
