package checkout.cli;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import checkout.spec.SpecFiles;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

/**
 * UNIT tests for finding the {@code spec/} folder by walking up the folder tree.
 *
 * <p>{@code @TempDir} gives each test a fresh temporary folder that JUnit deletes afterwards -
 * the standard way to test code that looks at the file system without touching real files.
 */
class SpecFolderTest {

    @Test
    void findsTheSpecFolderFromTheCurrentWorkingDirectory() {
        // Maven runs the tests in java/, the IDE maybe at the repository root: both are inside.
        Path spec = SpecFolder.find(Path.of("").toAbsolutePath());

        assertThat(spec).isEqualTo(SpecFiles.specDir());
        assertThat(spec.resolve("catalog.json")).isRegularFile();
    }

    @Test
    void findsTheSpecFolderFromADeeplyNestedStart() {
        Path deep = SpecFiles.specDir().resolve("test-cases");

        assertThat(SpecFolder.find(deep)).isEqualTo(SpecFiles.specDir());
    }

    @Test
    void acceptsAFileAsTheStartingPointSuchAsTheJar() {
        Path file = SpecFiles.specDir().resolve("catalog.json");

        assertThat(SpecFolder.find(file)).isEqualTo(SpecFiles.specDir());
    }

    @Test
    void findsTheNearestSpecFolderWhenThereAreSeveral(@TempDir Path root) throws IOException {
        Files.createDirectories(root.resolve("spec"));
        Files.writeString(root.resolve("spec/catalog.json"), "{}");
        Path nested = root.resolve("a/b/spec");
        Files.createDirectories(nested);
        Files.writeString(nested.resolve("catalog.json"), "{}");

        assertThat(SpecFolder.find(root.resolve("a/b/c"))).isEqualTo(nested);
        assertThat(SpecFolder.find(root.resolve("a"))).isEqualTo(root.resolve("spec"));
    }

    @Test
    void ignoresASpecFolderWithoutTheCatalog(@TempDir Path root) throws IOException {
        Files.createDirectories(root.resolve("spec"));

        assertThatThrownBy(() -> SpecFolder.find(root))
                .isInstanceOf(IllegalStateException.class)
                .hasMessageContaining("spec/catalog.json")
                .hasMessageContaining(root.toString());
    }

    @Test
    void programLocationIsTheClassesFolderOrJarThisCodeRunsFrom() {
        Path location = SpecFolder.programLocation();

        // Under Maven that is java/target/classes - which is inside the repository, so walking
        // up from it finds the spec folder. Running from a jar works the same way.
        assertThat(location).exists();
        assertThat(SpecFolder.find(location)).isEqualTo(SpecFiles.specDir());
    }

    @Test
    void findWithoutArgumentsUsesTheWorkingDirectoryOrTheProgramLocation() {
        assertThat(SpecFolder.find()).isEqualTo(SpecFiles.specDir());
    }
}
