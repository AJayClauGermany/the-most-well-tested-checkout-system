# Behaviour-driven scenarios (Gherkin)

These `.feature` files are written in [Gherkin](https://cucumber.io/docs/gherkin/reference/),
the *Given / When / Then* language used by Cucumber and its relatives:

| Language   | Typical tool                       |
|------------|------------------------------------|
| Python     | `pytest-bdd` or `behave`           |
| TypeScript | `@cucumber/cucumber`               |
| Go         | `godog`                            |
| Java       | Cucumber-JVM                       |
| C#         | Reqnroll                           |
| Rust       | `cucumber` crate                   |

Each language reads these files **directly from this folder** (do not copy them), so
every language is held to exactly the same scenarios.

## Step vocabulary

To keep the step definitions small, the features only use these steps.
If you need a new step, add it here as well.

| Step                                                             | Meaning |
|------------------------------------------------------------------|---------|
| `Given the standard catalog`                                     | Use [`catalog.json`](../catalog.json) |
| `Given an empty basket`                                          | A basket with no items |
| `Given a basket containing:` + a `sku` / `quantity` table        | Add these lines in order |
| `Given the customer uses the coupon "<code>"`                    | Set the coupon |
| `Given the customer chooses "<method>" shipping`                 | Set the shipping method |
| `When the customer checks out`                                   | Run the checkout and remember the receipt **or** the error |
| `Then the subtotal is $<amount>`                                 | |
| `Then the promotion discount is $<amount>`                       | |
| `Then the coupon discount is $<amount>`                          | |
| `Then the shipping cost is $<amount>`                            | |
| `Then the tax is $<amount>`                                      | |
| `Then the total is $<amount>`                                    | |
| `Then the receipt has <n> lines`                                 | |
| `Then the receipt line for "<sku>" has a quantity of <n>`        | |
| `Then checkout fails with the error "<CODE>"`                    | |
| `Given an empty order repository`                                | A fresh in-memory repository |
| `When the customer checks out and saves the order`               | Check out, then save the receipt if it succeeded |
| `Then the order is saved with the id "<id>"`                     | The id of the most recently saved order |
| `Then the saved order "<id>" has a total of $<amount>`           | |
| `Then the order repository is empty`                             | |
| `Then the order "<id>" cannot be found`                          | |
| `Given the clock is fixed at "<timestamp>"`                      | The repository's clock always answers this UTC time (rules §6). Orders saved afterwards get it as `created_at` |
| `Then the saved order "<id>" was created at "<timestamp>"`       | |

Timestamps are written exactly as they appear in a saved order (`2026-01-15T10:30:00Z`);
compare them as strings.

Amounts are written as dollars and cents (`$12.39`). Convert them to whole cents
(`1239`) before comparing – never parse them into a floating point number.
