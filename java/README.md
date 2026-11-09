# Checkout – Java

The Java version of the checkout system described in [`../spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: done.** `Checkout.checkout(...)`, saving orders (in memory, as JSON files and in a
> SQL database through JDBC) and the command line program are implemented, and every row of the
> Java column of the [testing matrix](../docs/testing-matrix.md) is covered. It was built strictly
> test-first: the commit that adds the tests comes before the commit that makes them pass – have
> a look at the git history.

## Tools we use, and why

| Tool | What it does | Why this one |
|------|--------------|--------------|
| **Java 21** | The language | A long-term-support (LTS) release with *records*, *sealed interfaces* and *pattern matching in `switch`*, which make data types short and safe. |
| **[Maven](https://maven.apache.org/)** | Build tool: downloads libraries, compiles, runs tests | The most widely used Java build tool. Everything is configured in one file, [`pom.xml`](pom.xml). |
| **[JUnit 6](https://docs.junit.org/) (Jupiter)** | Test framework, including `@ParameterizedTest` | The standard Java test framework. JUnit 6 has one version for the platform and Jupiter. |
| **[AssertJ](https://assertj.github.io/doc/)** | Assertions that read like sentences: `assertThat(receipt.total()).isEqualTo(2441)` | The most popular assertion library; its messages say exactly what differed. |
| **Maven Surefire / Failsafe** | Run the *unit* tests (`*Test`) in `mvn test` and the *integration* tests (`*IT`) in `mvn verify`, and write JUnit XML reports | Come with Maven. |
| **[Cucumber-JVM](https://cucumber.io/docs/installation/java/)** 8 | Runs the Given/When/Then scenarios in `spec/features/` | The reference Gherkin implementation for Java. Version 8 is the one built for JUnit Platform 6. |
| **[Mockito](https://site.mockito.org/)** | Mock objects: "was `save` called?" | The standard Java mocking library. |
| **[H2](https://h2database.com/)** | An SQL database in one file (or in memory), used through plain JDBC | Runs inside the tests: nothing to install, no server. It plays the role SQLite plays in the other languages. |
| **[jqwik](https://jqwik.net/) 1.9.3** | Property-based testing | The only mainstream property-based library for Java. Pinned to 1.9.x on purpose: 1.10 added an unusual licence clause (see the matrix footnote). |
| **[networknt json-schema-validator](https://github.com/networknt/json-schema-validator)** | Checks JSON against `spec/schema/*.json` (contract tests) | Supports JSON Schema draft 2020-12, which the spec uses, and is built on Jackson 3 like our code. |
| **[ApprovalTests.Java](https://github.com/approvals/ApprovalTests.Java)** | Snapshot ("approval") testing of whole receipts | Saves the approved output as a plain file next to the test, reviewed like code. |
| **[Jazzer](https://github.com/CodeIntelligenceTesting/jazzer)** | Fuzzing from inside JUnit (`@FuzzTest`) | The coverage-guided fuzzer for the JVM; in normal runs it replays its saved inputs like any test. |
| **[ArchUnit](https://www.archunit.org/)** | Architecture tests: "the domain never depends on the CLI" | Checks the compiled classes against rules written in plain Java. |
| **[Jackson 3](https://github.com/FasterXML/jackson)** | Reads and writes JSON (the CLI, the JSON-file store) | The most common Java JSON library. Version 3 no longer forces you to catch `IOException`. |
| **[JaCoCo](https://www.jacoco.org/jacoco/)** | Code coverage, with a 100% threshold | The standard Java coverage tool; it is a Java agent, so it sees unit and integration tests. |
| **[PIT](https://pitest.org/)** | Mutation testing | The mutation testing tool for the JVM. `pitest-junit5-plugin` lets it run JUnit Platform tests. |
| **`javac -Xlint:all -Werror`** | Static check | The compiler itself reports every warning it knows; `-Werror` turns warnings into errors so they cannot pile up. |
| **[Error Prone](https://errorprone.info/)** | More static checks, inside the compiler | Google's bug-pattern checker: catches mistakes such as ignored return values or locale-dependent `toUpperCase()`. |
| **[Checkstyle](https://checkstyle.org/)** | Style checks (imports, braces, line length) | Checks without rewriting anything, so it fits next to the compiler checks. Its rules are in [`checkstyle.xml`](checkstyle.xml). (The matrix suggests Checkstyle *or* Spotless – we picked one.) |

Maven has no separate lockfile: every library and plugin version is written out exactly in `pom.xml`.

## Prerequisites

* **Java 21** (a JDK, not only a JRE). Check with `java -version`.
  Install e.g. [Eclipse Temurin](https://adoptium.net/) or via your package manager / [SDKMAN!](https://sdkman.io/).
* **Maven 3.9+**. Check with `mvn -v`. Install from [maven.apache.org](https://maven.apache.org/install.html),
  your package manager (`brew install maven`, `sudo apt install maven`, `choco install maven`) or SDKMAN!.
* For the shared end-to-end runner only: `bash` and [`jq`](https://jqlang.github.io/jq/).

The first build downloads the libraries into `~/.m2/repository`; later builds are fast and work offline.

## Running the tests and checks

Run these commands inside the `java/` folder. Add `-B` ("batch mode") for plain output without colours, as CI does.

| I want to… | Command |
|------------|---------|
| run **everything CI runs** except mutation testing: static checks, all tests, the coverage check | `mvn verify` |
| run only the **unit tests** (everything named `*Test`, including BDD, properties, contracts, snapshots, fuzz seeds) | `mvn test` |
| run **one test class** | `mvn test -Dtest=CouponsTest` |
| run **one test method** | `mvn test -Dtest='MoneyTest#givesZeroForZeroPercent'` |
| run **one kind** of test, by tag | `mvn test -Dgroups=property` (tags: `spec`, `contract`, `snapshot`, `property`, `fuzz`, `architecture`, `mock`) |
| run the **BDD** scenarios from `spec/features/` | `mvn test -Dtest=RunCucumberTest` → `target/cucumber-reports/cucumber.html` |
| run only the **integration and end-to-end** tests | `mvn verify -Djacoco.skip -Dtest=NONE -Dsurefire.failIfNoSpecifiedTests=false` (tags: `integration`, `e2e`) |
| see the **coverage** report (after `mvn verify`; fails below 100%) | open `target/site/jacoco/index.html` |
| run **mutation testing** (about five minutes) | `mvn -Pmutation verify` → `target/pit-reports/index.html` |
| … only mutation testing, when the tests already passed (about three minutes) | `mvn test-compile org.pitest:pitest-maven:mutationCoverage` |
| **fuzz** for real for 20 seconds (instead of replaying the saved inputs) | `JAZZER_FUZZ=1 mvn test -Djacoco.skip -Dtest='CheckoutFuzzTest#checkoutOnlyEverThrowsCheckoutException'` |
| only compile, with all static checks (`-Xlint`, Error Prone, Checkstyle) | `mvn test-compile` |
| run the **shared e2e runner** against our CLI (from the repository root) | `(cd java && mvn package -DskipTests) && e2e/run.sh java -jar java/target/checkout.jar` |
| try the **command line program** | `echo '{"items":[{"sku":"APPLE","quantity":1}]}' \| java -jar target/checkout.jar` |
| delete all build output | `mvn clean` |

`mvn verify` prints **434 tests** run by Surefire and **102 tests** run by Failsafe, all passing,
then "All coverage checks have been met". The shared runner prints `60 passed, 0 failed`.

Test reports end up in `target/surefire-reports/` and `target/failsafe-reports/` (JUnit XML,
which CI uploads). You can also run any test from your IDE (IntelliJ IDEA, VS Code with the Java
extensions, Eclipse): open `pom.xml` as a project and click the green ▶ next to a test.

## Folder layout

```
java/
├── pom.xml                          ← libraries, plugins, coverage threshold, mutation profile
├── checkstyle.xml                   ← the style rules
├── .mvn/jvm.config                  ← lets Error Prone into the compiler (see "Static checks")
├── src/main/java/checkout/          ← the domain: pure calculation, no files, no terminal
│   ├── Checkout.java                ← checkout(basket, catalog): the pipeline of rules §4
│   ├── Validation.java              ← rows 1-4 of rules §4.1
│   ├── Lines.java                   ← merging basket items into receipt lines, pricing a line
│   ├── Promotions.java              ← buy X get Y free, bulk discount
│   ├── Coupons.java                 ← SAVE10, WELCOME5
│   ├── ShippingMethod.java          ← standard / express (an enum)
│   ├── Tax.java                     ← 8%
│   ├── Money.java                   ← roundHalfUp, with whole numbers only
│   ├── Basket, BasketItem, Catalog, Product, Promotion, Receipt, ReceiptLine  ← records
│   ├── ErrorCode.java, CheckoutException.java
│   └── orders/                      ← saving orders (rules §6)
│       ├── OrderRepository.java     ← the interface: save / find
│       ├── InMemoryOrderRepository.java  ← a list in memory; also the hand-written fake
│       ├── Orders.java              ← placeOrder = checkout + save
│       └── Order, PlacedOrder, OrderIds, Timestamps
├── src/main/java/checkout/store/    ← order stores that touch the outside world
│   ├── JsonFileOrderRepository.java ← one ORDER-<n>.json file per order
│   └── JdbcOrderRepository.java     ← one SQL table, through JDBC (H2 in the tests)
├── src/main/java/checkout/cli/      ← the command line program (rules §8)
│   ├── CheckoutCommand.java         ← its logic: JSON text in, output + exit code out
│   ├── Main.java                    ← stdin → CheckoutCommand → stdout, exit code
│   ├── SpecFolder.java              ← finds ../spec/ by walking up the folder tree
│   └── StandardCatalog.java         ← loads spec/catalog.json
└── src/test/java/checkout/          ← all tests; the package says what kind they are
    ├── *Test.java                   ← unit tests of the domain, one class per step
    ├── SharedSpecCasesTest.java     ← the 60 shared cases from spec/test-cases/
    ├── orders/                      ← unit tests + test doubles (fake, stub clock, Mockito)
    ├── cli/                         ← unit tests of the CLI logic
    ├── store/                       ← *IT integration tests (temp folder, H2) + mock-JDBC unit test
    ├── e2e/                         ← CliIT: the real jar, started with ProcessBuilder
    ├── bdd/                         ← Cucumber runner + step definitions
    ├── property/                    ← P1…P14 with jqwik
    ├── contract/                    ← receipts and CLI errors against spec/schema/
    ├── snapshot/                    ← ApprovalTests + the committed *.approved.json files
    ├── fuzz/                        ← Jazzer @FuzzTest (+ seeds in src/test/resources)
    ├── architecture/                ← ArchUnit rules
    └── spec/                        ← test helpers: SpecFiles, Schemas, Json, Fixtures, Dollars
```

**Where does a new test go?** In Java, the *class name* decides how a test is run, not the folder:

| Kind of test | Class name ends in | Run by | Command |
|--------------|-------------------|--------|---------|
| Fast test: no real files, processes or databases | `Test` | Surefire | `mvn test` |
| Integration / end-to-end test: files, databases, other processes | `IT` | Failsafe | `mvn verify` |

The shared spec files are **read from the repository's `spec/` folder** at test time, never
copied. `SpecFiles.specDir()` (tests) and `SpecFolder.find()` (the program) find that folder by
walking up from the current directory, so everything works from Maven, from your IDE, and from
any working directory inside the repository. The jar also looks next to itself, so
`java -jar /path/to/java/target/checkout.jar` works from anywhere.

## The kinds of tests, and where to look

| Kind of test | Question it answers | Look at |
|---|---|---|
| **Unit** | Does this one method do the right thing? | `MoneyTest` is a good first read, then `ValidationTest`, `CouponsTest`, … |
| **Parameterised** | Does it work for all these examples? | `@ValueSource` in `ValidationTest`, `@CsvSource` in `PromotionsTest` and `TaxTest`, `@MethodSource` in `SharedSpecCasesTest` and `CheckoutCommandTest` |
| **Shared spec cases** | Does it match the 60 cases every language must pass? | `SharedSpecCasesTest` |
| **BDD** | Does it behave as the Given/When/Then scenarios say? | `bdd/CheckoutSteps.java` – each sentence from `spec/features/README.md` appears once; `bdd/RunCucumberTest.java` points Cucumber at `../spec/features` |
| **Test doubles** | How do I test without the real clock, files or database? Was `save()` really never called? | a fake (`InMemoryOrderRepository`), stub clocks (`Clock.fixed`), Mockito mocks in `orders/PlaceOrderMockitoTest` and `store/JdbcOrderRepositoryTest` |
| **Integration** | Do the repositories really write files and database rows? | `store/JsonFileOrderRepositoryIT`, `store/JdbcOrderRepositoryIT` (each test gets a `@TempDir`) |
| **Property-based** | Is it true for *every* basket, not just my examples? | `property/CheckoutPropertiesTest` – `p1_…` to `p14_…` |
| **Contract** | Is the output even the right shape? | `contract/ReceiptContractTest`, helper `spec/Schemas.java` |
| **Snapshot / approval** | Did the output change since last time? | `snapshot/ReceiptApprovalTest` and its `*.approved.json` files |
| **End-to-end** | Does the real program work from the outside? | `e2e/CliIT`, then the shared runner `e2e/run.sh` |
| **Fuzzing** | Can *any* input crash it? (property P10) | `fuzz/CheckoutFuzzTest` |
| **Architecture** | Does the code respect its own layers? | `architecture/ArchitectureTest` |
| **Coverage** | Which lines did the tests never run? | `mvn verify` → `target/site/jacoco/index.html` – the threshold is 100% of lines and branches |
| **Mutation** | Would the tests notice if the code were wrong? | `mvn -Pmutation verify` – see below |
| **Static checks** | Can we catch mistakes without running anything? | `mvn test-compile`: `-Xlint:all -Werror`, Error Prone, Checkstyle |

All tests run **in parallel** (`src/test/resources/junit-platform.properties`), which proves they
do not depend on each other: no shared state, no fixed order.

### Mutation testing, and why the code looks the way it does

PIT makes 123 small changes to `src/main` ("mutants": `>=` becomes `>`, `+` becomes `-`, a
`return` value becomes `0` or `true`, a method call is removed) and runs the tests against each
one. A mutant that survives means no test noticed the change. The score is currently **100%**
(123 of 123 killed), and the build fails below 95%. CI runs it on every pull request (about three
minutes).

Getting there shaped both the tests and the code:

* Boundaries are tested on both sides (`2499` / `2500`, `4999` / `5000`, `9` / `10`, `0` / `1`,
  `99` / `100`), because `>` vs `>=` is the classic surviving mutant.
* Error **messages** are tested, not only error codes – otherwise PIT shows that the message
  could be anything.
* PIT found one gap: no test noticed when the JSON-file store stopped checking that a file name
  ends in `.json`, because every other file name in the tests failed a later check too. The test
  now puts an `ORDER-7.yaml` (an extension of the same length) into the folder.
* `checkout.cli.Main` is excluded from coverage and mutation testing: it runs in a separate
  process (see `e2e/`), which an in-process tool cannot observe. Its logic is in
  `CheckoutCommand`, which *is* measured. The BDD, fuzz and e2e test classes are not used by PIT
  either (Cucumber and Jazzer have their own runners, and a mutant never reaches the jar); the
  other tests cover the same code.
* Some errors are almost impossible to cause for real (an SQL error while *closing* a
  connection). `JdbcOrderRepository` therefore also accepts a `Connection`, so a test can hand it a
  Mockito mock that fails on purpose.

### Static checks

* **`-Xlint:all -Werror`**: every compiler warning is an error.
* **Error Prone** runs inside `javac` as a plugin. It needs access to compiler internals that
  Java hides by default; [`.mvn/jvm.config`](.mvn/jvm.config) opens them (Maven reads that file
  automatically), and on JDK 21 it also needs `-XDaddTypeAnnotationsToSymbol=true` (in `pom.xml`).
* **Checkstyle** runs first in every build (`validate` phase), on the main code and the tests.

Getting the build warning-free with all these libraries took a few JVM options for the tests,
each explained next to the `test.argLine` property in `pom.xml`: Mockito is loaded as a Java
agent (JDK 21 warns when a library attaches one to itself), `-Xshare:off` (that agent turns off
class data sharing anyway), and `-XX:+EnableDynamicAgentLoading` (Jazzer has to attach its agent
at run time).

## How to add a new test

Say you want to check that `Money.roundHalfUp` rounds 10% of 5 cents (0.5 cents) up to 1 cent.

1. **Pick the kind of test.** It is one method with made-up numbers: a unit test, so open
   `src/test/java/checkout/MoneyTest.java` (or create a new class whose name ends in `Test`).
2. **Write the test**, following the *Arrange / Act / Assert* pattern:

   ```java
   @Test
   void roundsHalfACentUpToOneCent() {
       // Arrange
       int amount = 5;

       // Act
       int rounded = Money.roundHalfUp(amount, 10);

       // Assert
       assertThat(rounded).isEqualTo(1);
   }
   ```

3. **Run it:** `mvn test -Dtest='MoneyTest#roundsHalfACentUpToOneCent'`.
4. **Make sure it can fail.** Change `1` to `0` and check that the test goes red.
   A test that never fails does not test anything! (Mutation testing automates this idea.)
5. **Run everything** before you open a pull request: `mvn verify`, and
   `mvn -Pmutation verify` if you changed `src/main`.

Want to add a case for *every* language instead? Add it to
[`spec/test-cases/`](../spec/test-cases/) and it is picked up automatically by the shared spec
cases, contract, e2e and shared-runner tests (and update the count in `SpecFilesTest`).

## Things to know about this implementation

* **Whole numbers only.** `Money.roundHalfUp(amount, percent)` is `(amount × percent + 50) / 100`
  with integer division, as the hint in rules §7 suggests. No `double` anywhere near money.
* **Parse, don't validate.** `ShippingMethod.parse` turns the basket's free-form `shipping`
  text into the enum `STANDARD` or `EXPRESS`, or throws; after that the compiler knows the value is fine.
  Promotions are a `sealed` interface, so the `switch` in `Promotions.discount` is checked for
  completeness by the compiler.
* **The clock is a `java.time.Clock`.** Production uses `Clock.systemUTC()`; tests pass
  `Clock.fixed(...)`, a Mockito mock, or a tiny hand-written `Clock` they can move forward.
* **Repositories are an interface** (`save` and `find`), with three implementations. The
  in-memory one doubles as the test fake. Ids are `ORDER-<n>`; the JSON-file store takes the next
  number after the highest existing file, the JDBC store uses an identity column.
* **`placeOrder` = checkout + save.** It cannot save anything when checkout throws, which is
  what the Mockito tests and property P14 check.
* **The CLI is strict about the *shape* of its input.** `items` must be a list of objects with a
  text `sku` and a whole-number `quantity`; anything else (including `1.5`, which an `int` cannot
  hold – the spec's test cases README expects this of strongly typed languages) is "not a basket":
  exit code 2 and a message on stderr. Unknown products, quantities out of range, odd shipping
  methods and coupons are left to `checkout`, which reports the error codes.

### Versions that had to fit together

JUnit 6 is new, so each JUnit Platform engine was checked against it:

* **Cucumber-JVM 8.0.4** is the line built for JUnit 6 (its `cucumber-junit-platform-engine` is
  compiled against `junit-bom` 6.1.3); the 7.x line targets JUnit 5.
* **jqwik 1.9.3** was built against JUnit Platform 1.13, but runs unchanged on Platform 6.1.3.
* **PIT 1.30.2 with `pitest-junit5-plugin` 1.2.3** uses the JUnit Platform launcher found on the
  test classpath (6.1.3 here) and runs every engine above.
* **Jazzer 0.30.0** declares JUnit 5.9 but the JUnit BOM lifts it to 6.1.3; both its replay
  mode and real fuzzing (`JAZZER_FUZZ=1`) work.
* **AssertJ** is pinned to 3.27.7: 4.0.0-M1 is a milestone, not a release.

## Ideas for the next contribution

Tier 1, the Tier 2 rows that name a Java tool, and Tier 3 of the matrix are complete for Java,
so what is left is polish:

* Compare with another language once it is implemented: same spec, same kinds of tests,
  different idioms. That comparison is the point of the project.
* A [JMH](https://github.com/openjdk/jmh) micro-benchmark of `checkout` (the matrix marks
  benchmarks as optional for Java).
* Run the JaCoCo agent inside the jar started by the e2e tests too, so the coverage report also
  shows which lines of `Main` the end-to-end tests ran.
* Cucumber can run scenarios in parallel (`cucumber.execution.parallel.enabled`); the step
  definitions keep all state per scenario, so it should just work.
