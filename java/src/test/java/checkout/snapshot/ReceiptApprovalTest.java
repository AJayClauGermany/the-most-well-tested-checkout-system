package checkout.snapshot;

import static org.assertj.core.api.Assertions.assertThat;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Catalog;
import checkout.Checkout;
import checkout.cli.CheckoutCommand;
import checkout.spec.Json;
import checkout.spec.SpecFiles;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import org.approvaltests.Approvals;
import org.approvaltests.core.Options;
import org.junit.jupiter.api.Test;

/**
 * SNAPSHOT (approval) tests with ApprovalTests: did the output change since last time?
 *
 * <p>{@code Approvals.verify(text)} compares the text with the file
 * {@code ReceiptApprovalTest.<testName>.approved.json} next to this class. The first time (or
 * when the output differs) it writes a {@code *.received.json} file and fails; you look at
 * the difference, and if the change is intended you rename "received" to "approved" (or let a
 * diff tool do it). The approved files are committed - they are the snapshot.
 *
 * <p>The hand-checked receipt in {@code spec/examples/receipt.json} plays the role of an
 * approved output shared by all six languages, so the first test also compares against it.
 */
class ReceiptApprovalTest {

    private static final Catalog CATALOG = SpecFiles.loadCatalog();
    private static final Options JSON_FILE = new Options().forFile().withExtension(".json");

    /** The basket behind spec/examples/receipt.json (rules section 5). */
    private static final Basket EXAMPLE_BASKET = new Basket(
            List.of(new BasketItem("COFFEE", 3), new BasketItem("APPLE", 2)), "SAVE10", "standard");

    @Test
    void theExampleReceiptFromTheRules() {
        Approvals.verify(Json.toPrettyJson(Checkout.checkout(EXAMPLE_BASKET, CATALOG)), JSON_FILE);
    }

    @Test
    void theExampleReceiptMatchesTheHandCheckedSpecExampleLinesIncluded() throws IOException {
        Path example = SpecFiles.specDir().resolve("examples/receipt.json");

        assertThat(Json.toTree(Checkout.checkout(EXAMPLE_BASKET, CATALOG)))
                .isEqualTo(Json.parse(Files.readString(example)));
    }

    @Test
    void theExampleReceiptAsTheCommandPrintsIt() {
        Approvals.verify(CheckoutCommand.run(Json.toJson(EXAMPLE_BASKET), CATALOG).output(), JSON_FILE);
    }

    @Test
    void theBiggestSharedCaseFullEverything() {
        // Every product, both promotions, WELCOME5 and express shipping.
        Basket basket = new Basket(
                List.of(
                        new BasketItem("HEADPHONES", 1),
                        new BasketItem("COFFEE", 6),
                        new BasketItem("NOTEBOOK", 12),
                        new BasketItem("PEN", 5),
                        new BasketItem("BREAD", 1),
                        new BasketItem("APPLE", 3)),
                "WELCOME5",
                "express");

        Approvals.verify(Json.toPrettyJson(Checkout.checkout(basket, CATALOG)), JSON_FILE);
    }

    @Test
    void anErrorAsTheCommandPrintsIt() {
        Approvals.verify(CheckoutCommand.run("{\"items\":[]}", CATALOG).output(), JSON_FILE);
    }
}
