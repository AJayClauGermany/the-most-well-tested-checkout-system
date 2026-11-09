package checkout.cli;

import java.nio.file.Path;

/**
 * Finds the repository's shared {@code spec/} folder, which holds the standard catalog.
 */
public final class SpecFolder {

    /** The file that proves a folder is the real {@code spec/} folder. */
    public static final String MARKER_FILE = "catalog.json";

    private SpecFolder() {}

    /**
     * Walks up from {@code start} (a folder, or a file such as a jar) until it reaches a
     * folder that contains {@code spec/catalog.json}.
     *
     * <p>Looking upwards, instead of hard-coding {@code ../spec}, means the program finds the
     * spec from any working directory inside the repository and from wherever the jar is.
     *
     * @param start where to begin looking
     * @return the {@code spec/} folder
     * @throws IllegalStateException when no folder above {@code start} contains {@code spec/catalog.json}
     */
    public static Path find(Path start) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * Finds the {@code spec/} folder from the current working directory, or failing that
     * from the location of this program's own classes (the jar file, or {@code target/classes}).
     *
     * @return the {@code spec/} folder
     * @throws IllegalStateException when neither place is inside the repository
     */
    public static Path find() {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * Where this program's own classes were loaded from: the jar file when run with
     * {@code java -jar}, or the {@code target/classes} folder when run from Maven or an IDE.
     *
     * @return the jar or classes folder
     */
    public static Path programLocation() {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
