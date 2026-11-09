package checkout.bdd;

import static io.cucumber.junit.platform.engine.Constants.GLUE_PROPERTY_NAME;

import org.junit.platform.suite.api.ConfigurationParameter;
import org.junit.platform.suite.api.IncludeEngines;
import org.junit.platform.suite.api.SelectDirectories;
import org.junit.platform.suite.api.Suite;

/**
 * BDD (behaviour-driven development) tests: runs the Gherkin scenarios in the repository's
 * {@code spec/features/*.feature} files with Cucumber.
 *
 * <p>This class has no test methods. It is a JUnit Platform {@code @Suite} that hands the
 * work to the Cucumber engine:
 * <ul>
 *   <li>{@code @SelectDirectories("../spec/features")} - the feature files are read straight
 *       from the shared {@code spec/} folder (relative to {@code java/}, where Maven runs), never
 *       copied, so every language is held to exactly the same scenarios;</li>
 *   <li>{@code cucumber.glue} - the package with the step definitions, {@link CheckoutSteps}.</li>
 * </ul>
 * Surefire picks this class up because its name ends in {@code Test}. The other Cucumber settings
 * (naming, output) are in {@code src/test/resources/junit-platform.properties}.
 */
@Suite
@IncludeEngines("cucumber")
@SelectDirectories("../spec/features")
@ConfigurationParameter(key = GLUE_PROPERTY_NAME, value = "checkout.bdd")
class RunCucumberTest {}
