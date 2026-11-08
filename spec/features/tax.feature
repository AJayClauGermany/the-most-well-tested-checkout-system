Feature: Tax
  Tax is 8% of the price after all discounts, rounded to the nearest cent.
  Shipping is not taxed.

  Background:
    Given the standard catalog

  Scenario Outline: Tax is rounded to the nearest cent
    Given a basket containing:
      | sku   | quantity   |
      | <sku> | <quantity> |
    When the customer checks out
    Then the tax is $<tax>

    Examples:
      | sku   | quantity | tax  |
      | APPLE | 1        | 0.04 |
      | PEN   | 1        | 0.10 |
      | PEN   | 2        | 0.19 |

  Scenario: Shipping is not taxed
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    And the customer chooses "express" shipping
    When the customer checks out
    Then the tax is $0.04
    And the total is $10.53

  Scenario: Tax is charged on the price after discounts
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 10       |
    And the customer uses the coupon "SAVE10"
    When the customer checks out
    Then the tax is $0.36
