# Checkout – Java

The Java version of the checkout system described in [`../spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: tests first, code later.** The domain types and the test setup are ready, but
> `Checkout.checkout(...)` only throws `UnsupportedOperationException("not implemented yet")`.
> The 60 shared spec cases are already wired up and waiting (disabled) for an implementation.
> That is test-driven development (TDD): write down what "correct" means *before* writing the code.

## Tools we use, and why

| Tool | What it does | Why this one |
|------|--------------|--------------|
| **Java 21** | The language | A long-term-support (LTS) release with *records* and *sealed interfaces*, which make data types short and safe. |
| **[Maven](https://maven.apache.org/)** | Build tool: downloads libraries, compiles, runs tests | The most widely used Java build tool. Everything is configured in one file, [`pom.xml`](pom.xml). |
| **[JUnit 6](https://docs.junit.org/) (Jupiter)** | Test framework | The standard Java test framework. `@Test`, `assertEquals`, `@ParameterizedTest`. |
| **Maven Surefire** | Runs the *unit* tests (`*Test.java`) in `mvn test` | Comes with Maven. |
| **Maven Failsafe** | Runs the *integration* tests (`*IT.java`) in `mvn verify` | The standard way to keep slow/outside-world tests separate from fast unit tests. |
| **[Jackson 3](https://github.com/FasterXML/jackson)** | Reads JSON into Java records | The most common Java JSON library. Version 3 no longer forces you to catch `IOException`. |
| **`javac -Xlint:all -Werror`** | Static check | The compiler itself reports every warning it knows; `-Werror` turns warnings into errors so they cannot pile up. |

Maven has no separate lockfile: every library and plugin version is written out exactly in `pom.xml`.

## Prerequisites

* **Java 21** (a JDK, not only a JRE). Check with `java -version`.
  Install e.g. [Eclipse Temurin](https://adoptium.net/) or via your package manager / [SDKMAN!](https://sdkman.io/).
* **Maven 3.9+**. Check with `mvn -v`. Install from [maven.apache.org](https://maven.apache.org/install.html),
  your package manager (`brew install maven`, `sudo apt install maven`, `choco install maven`) or SDKMAN!.

The first build downloads the libraries into `~/.m2/repository`; later builds are fast and work offline.

## Running the tests

Run these commands inside the `java/` folder.

| I want to… | Command |
|------------|---------|
| run **everything** (compile + unit tests + integration tests) – this is what CI does | `mvn verify` |
| run only the **unit tests** | `mvn test` |
| run **one test class** | `mvn test -Dtest=CatalogTest` |
| run **one test method** | `mvn test -Dtest='CatalogTest#findsAProductBySku'` |
| run **one integration test** (skipping unit tests) | `mvn verify -Dtest=NONE -Dsurefire.failIfNoSpecifiedTests=false -Dit.test=SpecFolderIT` |
| skip the integration tests | `mvn verify -DskipITs` |
| only compile (and run the static check) | `mvn test-compile` |
| delete all build output | `mvn clean` |

Add `-B` ("batch mode") for plain output without colours, as CI does.
Test reports end up in `target/surefire-reports/` and `target/failsafe-reports/`.

You can also run any test from your IDE (IntelliJ IDEA, VS Code with the Java extensions, Eclipse):
open `pom.xml` as a project and click the green ▶ next to a test.

Expected result today: **12 unit tests** (11 passed, 1 skipped – the shared spec test) and **1 integration test** passed.

## Folder layout

```
java/
├── pom.xml                         ← build file: Java version, libraries, plugins
├── src/main/java/checkout/         ← the real program
│   ├── Checkout.java               ← checkout(basket, catalog) → Receipt   (not implemented yet)
│   ├── Basket.java, BasketItem.java
│   ├── Catalog.java, Product.java, Promotion.java
│   ├── Receipt.java, ReceiptLine.java
│   └── ErrorCode.java, CheckoutException.java
└── src/test/java/checkout/         ← all tests (Maven keeps them apart from the real program)
    ├── CatalogTest.java            ← example UNIT test – start reading here!
    ├── BasketTest.java             ← unit test
    ├── SharedSpecCasesTest.java    ← one test per spec/test-cases case (disabled until implemented)
    ├── SpecFolderIT.java           ← example INTEGRATION test (name ends in IT)
    ├── spec/                       ← test helpers that load ../spec/*.json (+ their own tests)
    ├── bdd/                        ← future home of Cucumber BDD tests (see package-info.java)
    └── e2e/                        ← future home of end-to-end tests (see package-info.java)
```

**Where does a new test go?** In Java, the *class name* decides how a test is run, not the folder:

| Kind of test | Class name ends in | Run by | Command |
|--------------|-------------------|--------|---------|
| Unit test (fast, no files/network/processes) | `Test` | Surefire | `mvn test` |
| Integration / end-to-end test (files, processes, databases…) | `IT` | Failsafe | `mvn verify` |

Tests usually live in the **same package** as the code they test (`checkout`), mirroring
`src/main/java`. That lets them use package-private methods too.

The shared spec files are **read from the repository's `spec/` folder** at test time, never
copied. `SpecFiles.specDir()` finds that folder by walking up from the current directory,
so the tests work from Maven, from your IDE, and from any working directory inside the repo.

## The domain types in 30 seconds

* All types are **records**: immutable data classes where Java writes the constructor,
  getters, `equals`, `hashCode` and `toString` for you.
* JSON uses `snake_case` (`unit_price`), Java uses `camelCase` (`unitPrice`).
  `@JsonProperty("unit_price")` connects the two.
* All money is an `int` number of **cents** (`899` = $8.99). Never `double`!
* `Promotion` is a **sealed interface** with two records, `Promotion.BuyXGetYFree` and
  `Promotion.BulkDiscount`. "Sealed" means no other kinds can exist, so a `switch` over a
  promotion can handle both cases and the compiler checks you did not forget one.
  (A single record with every field and a `type` string would also work, but then half the
  fields would be meaningless for each kind – easy to mix up.)
* `Basket.shipping()` is `"standard"` when the basket does not say otherwise.
* Errors are reported by throwing `CheckoutException`; `getCode()` returns one of the six
  `ErrorCode`s from rules §9.

## How to add a new test

Say you want to check that `Catalog.findBySku` finds the *first* product of the standard catalog.

1. Open (or create) the test class next to the code it tests:
   `src/test/java/checkout/CatalogTest.java`. The class name must end in `Test`.
2. Add a method with the `@Test` annotation. Its name says what should happen:

   ```java
   @Test
   void findsTheFirstProductOfTheStandardCatalog() {
       // Arrange
       Catalog catalog = SpecFiles.loadCatalog();   // import checkout.spec.SpecFiles;

       // Act
       Product apple = catalog.findBySku("APPLE").orElseThrow();

       // Assert – expected value first, actual value second
       assertEquals("Apple", apple.name());
       assertEquals(50, apple.unitPrice());
   }
   ```

3. Run just that test: `mvn test -Dtest='CatalogTest#findsTheFirstProductOfTheStandardCatalog'`.
4. Make it fail on purpose once (change `50` to `51`) to see that it really checks something,
   then change it back.
5. Run everything with `mvn verify` before you open a pull request.

**Starting the implementation?** Delete the `@Disabled(...)` line in `SharedSpecCasesTest.java`.
You will see 60 failing tests, one per spec case, named after the case id. Make them pass one by one.
Each step of the rules (validate, merge, promotions, coupon, shipping, tax) deserves its own small
unit tests too.

## Testing ideas we'd love contributions for

Each of these shows off a different kind of testing. Pick one, open a pull request!

* **Nicer assertions** – [AssertJ](https://assertj.github.io/doc/): `assertThat(receipt.total()).isEqualTo(2441)`.
* **BDD** – run `../spec/features/*.feature` with [Cucumber-JVM](https://cucumber.io/docs/installation/java/)
  (`cucumber-java` + `cucumber-junit-platform-engine`); step definitions go in `checkout.bdd`.
* **Property-based testing** – [jqwik](https://jqwik.net/): e.g. "the total is never negative",
  "merging lines never changes the subtotal", "`roundHalfUp` matches `BigDecimal`'s `HALF_UP`".
* **Mocks and fakes** – [Mockito](https://site.mockito.org/) for the order repository (rules §6),
  or a hand-written in-memory fake to compare with.
* **Integration tests with a real database** – a JDBC order repository tested against
  [Testcontainers](https://java.testcontainers.org/) or an H2 in-memory database (`*IT.java`).
* **End-to-end tests** – a `main` method implementing the CLI of rules §8, then `*IT` tests in
  `checkout.e2e` that start it with `ProcessBuilder` and feed it every spec basket.
* **Snapshot / approval testing** – [ApprovalTests.Java](https://github.com/approvals/ApprovalTests.Java)
  for whole receipts as JSON.
* **Code coverage** – [JaCoCo](https://www.jacoco.org/jacoco/) Maven plugin, with a report in CI.
* **Mutation testing** – [PIT (pitest)](https://pitest.org/): do our tests notice when the code is subtly broken?
* **Architecture tests** – [ArchUnit](https://www.archunit.org/): e.g. "the domain does not depend on the CLI".
* **More static checks** – [Error Prone](https://errorprone.info/), [Checkstyle](https://checkstyle.org/),
  or formatting with [Spotless](https://github.com/diffplug/spotless) + google-java-format.
* **JSON Schema validation** – check `Receipt` JSON output against a schema with
  [networknt json-schema-validator](https://github.com/networknt/json-schema-validator).
