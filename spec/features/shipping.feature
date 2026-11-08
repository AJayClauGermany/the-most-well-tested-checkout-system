Feature: Shipping
  Standard shipping costs $4.99 and is free when the price after discounts is at least $50.00.
  Express shipping always costs $9.99.

  Background:
    Given the standard catalog

  Scenario: Standard shipping is used when none is chosen
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    When the customer checks out
    Then the shipping cost is $4.99

  Scenario: Standard shipping is free at exactly $50.00
    Given a basket containing:
      | sku   | quantity |
      | BREAD | 20       |
      | APPLE | 10       |
    And the customer chooses "standard" shipping
    When the customer checks out
    Then the shipping cost is $0.00

  Scenario: A coupon can push the order below the free shipping threshold
    Given a basket containing:
      | sku   | quantity |
      | BREAD | 20       |
      | APPLE | 10       |
    And the customer uses the coupon "SAVE10"
    When the customer checks out
    Then the shipping cost is $4.99

  Scenario: Express shipping is never free
    Given a basket containing:
      | sku        | quantity |
      | HEADPHONES | 1        |
    And the customer chooses "express" shipping
    When the customer checks out
    Then the shipping cost is $9.99

  Scenario: An unknown shipping method is rejected
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    And the customer chooses "overnight" shipping
    When the customer checks out
    Then checkout fails with the error "UNKNOWN_SHIPPING_METHOD"
