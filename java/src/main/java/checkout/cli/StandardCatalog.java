package checkout.cli;

import checkout.Catalog;
import java.nio.file.Path;

/**
 * The standard catalog, {@code spec/catalog.json}, which the command line program always uses
 * (rules section 8).
 */
public final class StandardCatalog {

    private StandardCatalog() {}

    /**
     * Where the standard catalog lives.
     *
     * @return {@code <repository>/spec/catalog.json}
     * @throws IllegalStateException when the {@code spec/} folder cannot be found (see {@link SpecFolder#find()})
     */
    public static Path path() {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * Reads the standard catalog from disk.
     *
     * @return the six standard products
     * @throws IllegalStateException when the file cannot be found or read
     */
    public static Catalog load() {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
