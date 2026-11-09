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
 *
 * <p>Run only the scenarios with {@code mvn test -Dtest=RunCucumberTest}. (Tags cannot select
 * them: a scenario only has the tags written in its {@code .feature} file, and the shared files
 * have none.) {@code failIfNoTests = false}: when you run only some kinds of tests by tag
 * ({@code mvn test -Dgroups=property}), this suite is left empty, which is fine. That the feature
 * files are really found is checked by {@code checkout.SpecFolderIT} and by the scenario count in
 * the test report.
 */
@Suite(failIfNoTests = false)
@IncludeEngines("cucumber")
@SelectDirectories("../spec/features")
@ConfigurationParameter(key = GLUE_PROPERTY_NAME, value = "checkout.bdd")
class RunCucumberTest {}
