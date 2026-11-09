package checkout.architecture;

import static com.tngtech.archunit.lang.syntax.ArchRuleDefinition.classes;
import static com.tngtech.archunit.lang.syntax.ArchRuleDefinition.noClasses;

import checkout.orders.OrderRepository;
import com.tngtech.archunit.core.domain.JavaClasses;
import com.tngtech.archunit.core.importer.ImportOption;
import com.tngtech.archunit.junit.AnalyzeClasses;
import com.tngtech.archunit.junit.ArchTag;
import com.tngtech.archunit.junit.ArchTest;

/**
 * ARCHITECTURE tests with ArchUnit: rules about how the code is allowed to depend on itself.
 *
 * <p>The checkout is split into layers by package:
 * <ul>
 *   <li>{@code checkout} and {@code checkout.orders} - the <em>domain</em>: pure calculation and the
 *       repository interface. No files, no databases, no terminal.</li>
 *   <li>{@code checkout.store} - repositories that talk to the file system and a database.</li>
 *   <li>{@code checkout.cli} - the command line program.</li>
 * </ul>
 * An ordinary test cannot notice when somebody adds {@code import checkout.cli.Main} to the
 * domain by accident; this test can. ArchUnit reads the compiled classes (the main code only,
 * not the tests) once, and hands them to every {@code @ArchTest} method, which is run by
 * ArchUnit's own JUnit Platform engine.
 */
@ArchTag("architecture")
@AnalyzeClasses(packages = "checkout", importOptions = ImportOption.DoNotIncludeTests.class)
class ArchitectureTest {

    private static final String[] DOMAIN = {"checkout", "checkout.orders"};

    @ArchTest
    void theDomainDoesNotDependOnTheCliOrTheStores(JavaClasses classes) {
        noClasses()
                .that().resideInAnyPackage(DOMAIN)
                .should().dependOnClassesThat().resideInAnyPackage("checkout.cli..", "checkout.store..")
                .because("the calculation must stay usable without a terminal, files or a database")
                .check(classes);
    }

    @ArchTest
    void theDomainDoesNoInputOrOutput(JavaClasses classes) {
        noClasses()
                .that().resideInAnyPackage(DOMAIN)
                .should().dependOnClassesThat()
                .resideInAnyPackage("java.io..", "java.nio..", "java.sql..", "java.net..")
                .because("the domain only calculates; checkout.store and checkout.cli do the reading and writing")
                .check(classes);
    }

    @ArchTest
    void theStoresDoNotDependOnTheCli(JavaClasses classes) {
        noClasses()
                .that().resideInAPackage("checkout.store..")
                .should().dependOnClassesThat().resideInAPackage("checkout.cli..")
                .because("a repository must not know how the program is started")
                .check(classes);
    }

    @ArchTest
    void onlyTheCliTalksToTheTerminal(JavaClasses classes) {
        noClasses()
                .that().resideOutsideOfPackage("checkout.cli..")
                .should().accessClassesThat().belongToAnyOf(System.class)
                .because("System.in / System.out / System.exit belong to the command line program only")
                .check(classes);
    }

    @ArchTest
    void repositoriesAreNamedAfterTheInterface(JavaClasses classes) {
        classes()
                .that().implement(OrderRepository.class)
                .should().haveSimpleNameEndingWith("OrderRepository")
                .because("a consistent name makes the three implementations easy to find")
                .check(classes);
    }
}
