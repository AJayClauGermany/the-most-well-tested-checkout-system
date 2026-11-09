package checkout.cli;

import static org.assertj.core.api.Assertions.assertThat;

import checkout.Catalog;
import checkout.Product;
import checkout.spec.SpecFiles;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests for loading the standard catalog the command line program uses.
 */
class StandardCatalogTest {

    @Test
    void pointsAtSpecCatalogJsonAtTheRootOfTheRepository() {
        assertThat(StandardCatalog.path()).isEqualTo(SpecFiles.specDir().resolve("catalog.json"));
        // The repository root is the folder that also holds the rules.
        assertThat(StandardCatalog.path().getParent().resolve("checkout-rules.md")).isRegularFile();
    }

    @Test
    void loadsTheSixStandardProductsInCatalogOrder() {
        Catalog catalog = StandardCatalog.load();

        assertThat(catalog.products())
                .extracting(Product::sku)
                .containsExactly("APPLE", "BREAD", "COFFEE", "NOTEBOOK", "PEN", "HEADPHONES");
    }

    @Test
    void loadsTheSameCatalogAsTheTestSupportCode() {
        assertThat(StandardCatalog.load()).isEqualTo(SpecFiles.loadCatalog());
    }
}
