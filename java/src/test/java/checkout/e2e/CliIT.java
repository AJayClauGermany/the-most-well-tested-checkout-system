package checkout.e2e;

import static org.assertj.core.api.Assertions.assertThat;

import checkout.spec.Json;
import checkout.spec.Schemas;
import checkout.spec.SpecCase;
import checkout.spec.SpecFiles;
import checkout.spec.Totals;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.stream.Stream;
import org.junit.jupiter.api.Named;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.MethodSource;
import tools.jackson.databind.JsonNode;

/**
 * END-TO-END tests: run the real command line program (rules section 8) as a separate process,
 * exactly as a user (or the shared runner in {@code e2e/run.sh}) would, and look only at what
 * comes out: standard output, standard error and the exit code.
 *
 * <p>The program is {@code java -jar target/checkout.jar}. Failsafe runs this class after
 * {@code mvn package} has built that jar, and passes its path as the system property
 * {@code checkout.jar}. Because the test only sees stdout, stderr and the exit code, it would
 * pass unchanged for a program written in any language.
 */
@Tag("e2e")
class CliIT {

    /** What one run of the program produced. */
    record CliResult(String stdout, String stderr, int exitCode) {}

    private static Path jar() {
        return Path.of(System.getProperty("checkout.jar", "target/checkout.jar")).toAbsolutePath();
    }

    /** The same java that runs this test. */
    private static String java() {
        return ProcessHandle.current().info().command().orElse("java");
    }

    /** Runs the program once, with {@code input} on standard input, in the given working directory. */
    private static CliResult run(String input, Path workingDirectory) throws IOException, InterruptedException {
        ProcessBuilder builder = new ProcessBuilder(java(), "-jar", jar().toString())
                .directory(workingDirectory.toFile());
        // If these are set (some CI machines and proxies do), every JVM announces them on stderr
        // ("Picked up JAVA_TOOL_OPTIONS: ..."), which is not the program's output.
        builder.environment().remove("JAVA_TOOL_OPTIONS");
        builder.environment().remove("JDK_JAVA_OPTIONS");
        Process process = builder.start();
        try (var stdin = process.getOutputStream()) {
            stdin.write(input.getBytes(StandardCharsets.UTF_8));
        }
        String stdout = new String(process.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
        String stderr = new String(process.getErrorStream().readAllBytes(), StandardCharsets.UTF_8);
        int exitCode = process.waitFor();
        return new CliResult(stdout, stderr, exitCode);
    }

    private static CliResult run(String input) throws IOException, InterruptedException {
        return run(input, Path.of("").toAbsolutePath());
    }

    static Stream<Named<SpecCase>> specCases() {
        return SpecFiles.loadAllCases().stream().map(specCase -> Named.of(specCase.id(), specCase));
    }

    @Test
    void theJarExists() {
        assertThat(jar()).as("run 'mvn package' first").isRegularFile();
    }

    // Every shared case, through the real program.
    @ParameterizedTest(name = "{0}")
    @MethodSource("specCases")
    void answersEverySharedCaseLikeTheSpecSays(SpecCase specCase) throws IOException, InterruptedException {
        CliResult result = run(Json.toJson(specCase.basket()) + "\n");

        if (specCase.expectedError() != null) {
            assertThat(result.exitCode()).as("stderr: %s", result.stderr()).isEqualTo(1);
            assertThat(Json.parse(result.stdout()))
                    .isEqualTo(Json.parse("{\"error\":\"" + specCase.expectedError().name() + "\"}"));
        } else {
            assertThat(result.exitCode()).as("stderr: %s", result.stderr()).isZero();
            JsonNode receipt = Json.parse(result.stdout());
            assertThat(Json.MAPPER.treeToValue(receipt, Totals.class)).isEqualTo(specCase.expected());
            assertThat(Schemas.problems("receipt.schema.json", receipt)).isEmpty();
        }
    }

    @Test
    void printsExactlyOneLineOfJson() throws IOException, InterruptedException {
        CliResult result = run("{\"items\":[{\"sku\":\"APPLE\",\"quantity\":1}]}");

        assertThat(result.stdout()).endsWith("\n");
        assertThat(result.stdout().strip().lines().toList()).hasSize(1);
        assertThat(result.stderr()).isEmpty();
    }

    @Test
    void exitsWith2AndExplainsOnStderrWhenTheInputIsNotJson() throws IOException, InterruptedException {
        CliResult result = run("this is not json");

        assertThat(result.exitCode()).isEqualTo(2);
        assertThat(result.stdout()).isEmpty();
        assertThat(result.stderr()).startsWith("checkout: ").contains("not JSON");
    }

    @Test
    void exitsWith2WhenTheJsonIsNotABasket() throws IOException, InterruptedException {
        CliResult result = run("{\"hello\": \"world\"}");

        assertThat(result.exitCode()).isEqualTo(2);
        assertThat(result.stdout()).isEmpty();
        assertThat(result.stderr()).contains("not a basket");
    }

    @Test
    void exitsWith2OnEmptyInput() throws IOException, InterruptedException {
        assertThat(run("").exitCode()).isEqualTo(2);
    }

    @Test
    void findsTheStandardCatalogFromTheJarsLocationWhenRunFromOutsideTheRepository(@TempDir Path elsewhere)
            throws IOException, InterruptedException {
        // No spec/ folder above the working directory here - the program must look next to the jar.
        assertThat(Files.exists(elsewhere.resolve("spec"))).isFalse();

        CliResult result = run("{\"items\":[{\"sku\":\"APPLE\",\"quantity\":1}]}\n", elsewhere);

        assertThat(result.exitCode()).as("stderr: %s", result.stderr()).isZero();
        assertThat(Json.parse(result.stdout()).get("total").asInt()).isEqualTo(553);
    }

    @Test
    void usesTheStandardCatalogNotSomeOtherOne() throws IOException, InterruptedException {
        // HEADPHONES only exist in spec/catalog.json: $59.99 -> free standard shipping, tax $4.80.
        CliResult result = run("{\"items\":[{\"sku\":\"HEADPHONES\",\"quantity\":1}]}");

        assertThat(Json.MAPPER.treeToValue(Json.parse(result.stdout()), Totals.class))
                .isEqualTo(new Totals(5999, 0, 0, 0, 480, 6479));
    }

    @Test
    void theSharedCasesCoverEveryErrorCode() {
        // Not a test of the program, but of the test data: the e2e runs above exercise every exit path.
        List<String> codes = SpecFiles.loadAllCases().stream()
                .filter(specCase -> specCase.expectedError() != null)
                .map(specCase -> specCase.expectedError().name())
                .distinct()
                .toList();

        assertThat(codes).containsExactlyInAnyOrder(
                "EMPTY_BASKET", "UNKNOWN_PRODUCT", "INVALID_QUANTITY",
                "UNKNOWN_SHIPPING_METHOD", "UNKNOWN_COUPON", "COUPON_MINIMUM_NOT_MET");
    }
}
