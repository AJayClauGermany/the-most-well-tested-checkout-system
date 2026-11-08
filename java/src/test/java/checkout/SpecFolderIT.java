package checkout;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

import checkout.spec.SpecFiles;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.stream.Stream;
import org.junit.jupiter.api.Test;

/**
 * An example INTEGRATION test. Class names ending in {@code IT} are run by Maven's
 * Failsafe plugin during {@code mvn verify} - after the unit tests and after the jar
 * has been built.
 *
 * <p>Integration tests check that our code works together with something outside it:
 * the file system, a database, another program... Here it is the file system: the
 * Gherkin files in {@code spec/features/} (which a future Cucumber test will use) must
 * be readable from disk.
 */
class SpecFolderIT {

    @Test
    void featureFilesCanBeReadFromDisk() throws IOException {
        List<Path> featureFiles;
        try (Stream<Path> files = Files.list(SpecFiles.specDir().resolve("features"))) {
            featureFiles = files.filter(path -> path.toString().endsWith(".feature")).toList();
        }

        assertFalse(featureFiles.isEmpty(), "expected some .feature files in spec/features");
        for (Path featureFile : featureFiles) {
            String text = Files.readString(featureFile);
            assertTrue(text.contains("Feature:"), featureFile + " should contain a 'Feature:' line");
        }
    }
}
