package checkout.cli;

import java.net.URI;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Optional;

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
        return search(start).orElseThrow(() -> notFound(start.toString()));
    }

    /**
     * Finds the {@code spec/} folder from the current working directory, or failing that
     * from the location of this program's own classes (the jar file, or {@code target/classes}).
     *
     * @return the {@code spec/} folder
     * @throws IllegalStateException when neither place is inside the repository
     */
    public static Path find() {
        return find(Path.of("").toAbsolutePath(), programLocation());
    }

    /**
     * Looks above the working directory first, then above the program's location.
     *
     * @param workingDirectory where the program was started
     * @param programLocation  where the program's classes are (see {@link #programLocation()})
     * @return the {@code spec/} folder
     * @throws IllegalStateException when neither place is inside the repository
     */
    static Path find(Path workingDirectory, Path programLocation) {
        return search(workingDirectory)
                .or(() -> search(programLocation))
                .orElseThrow(() -> notFound(workingDirectory + " or " + programLocation));
    }

    /**
     * Where this program's own classes were loaded from: the jar file when run with
     * {@code java -jar}, or the {@code target/classes} folder when run from Maven or an IDE.
     *
     * @return the jar or classes folder
     */
    public static Path programLocation() {
        // A "file:/..." URL. URI.create + Path.of turn it into a path, decoding e.g. %20 into a space.
        String location = SpecFolder.class.getProtectionDomain().getCodeSource().getLocation().toString();
        return Path.of(URI.create(location));
    }

    private static Optional<Path> search(Path start) {
        for (Path folder = start.toAbsolutePath().normalize(); folder != null; folder = folder.getParent()) {
            Path spec = folder.resolve("spec");
            if (Files.isRegularFile(spec.resolve(MARKER_FILE))) {
                return Optional.of(spec);
            }
        }
        return Optional.empty();
    }

    private static IllegalStateException notFound(String where) {
        return new IllegalStateException("No spec/" + MARKER_FILE + " found in any folder above " + where);
    }
}
