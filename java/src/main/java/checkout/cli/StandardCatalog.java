package checkout.cli;

import checkout.Catalog;
import checkout.Product;
import java.nio.file.Path;
import java.util.List;
import tools.jackson.databind.json.JsonMapper;

/**
 * The standard catalog, {@code spec/catalog.json}, which the command line program always uses
 * (rules section 8).
 */
public final class StandardCatalog {

    private static final JsonMapper JSON = JsonMapper.builder().build();

    /** The shape of {@code catalog.json}: a description for humans, and the products. */
    private record CatalogFile(String description, List<Product> products) {}

    private StandardCatalog() {}

    /**
     * Where the standard catalog lives.
     *
     * @return {@code <repository>/spec/catalog.json}
     * @throws IllegalStateException when the {@code spec/} folder cannot be found (see {@link SpecFolder#find()})
     */
    public static Path path() {
        return SpecFolder.find().resolve(SpecFolder.MARKER_FILE);
    }

    /**
     * Reads the standard catalog from disk.
     *
     * @return the six standard products
     * @throws IllegalStateException when the {@code spec/} folder cannot be found
     * @throws tools.jackson.core.JacksonException when the file cannot be read or is not a catalog
     */
    public static Catalog load() {
        return new Catalog(JSON.readValue(path().toFile(), CatalogFile.class).products());
    }
}
