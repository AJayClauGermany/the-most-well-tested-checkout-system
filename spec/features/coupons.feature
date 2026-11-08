Feature: Coupons
  A shopper may enter one coupon code.
  SAVE10 takes 10% off. WELCOME5 takes $5.00 off when spending at least $25.00.

  Background:
    Given the standard catalog

  Scenario: SAVE10 takes 10% off
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 10       |
    And the customer uses the coupon "SAVE10"
    When the customer checks out
    Then the coupon discount is $0.50

  Scenario: Half a cent is rounded up
    Given a basket containing:
      | sku   | quantity |
      | BREAD | 1        |
    And the customer uses the coupon "SAVE10"
    When the customer checks out
    Then the coupon discount is $0.23

  Scenario: Coupon codes ignore case and surrounding spaces
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 10       |
    And the customer uses the coupon "  save10 "
    When the customer checks out
    Then the coupon discount is $0.50

  Scenario: Coupons apply to the price after promotions
    Given a basket containing:
      | sku    | quantity |
      | COFFEE | 3        |
    And the customer uses the coupon "SAVE10"
    When the customer checks out
    Then the promotion discount is $8.99
    And the coupon discount is $1.80

  Scenario: WELCOME5 at exactly the minimum spend
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 50       |
    And the customer uses the coupon "WELCOME5"
    When the customer checks out
    Then the coupon discount is $5.00

  Scenario: WELCOME5 below the minimum spend
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 49       |
    And the customer uses the coupon "WELCOME5"
    When the customer checks out
    Then checkout fails with the error "COUPON_MINIMUM_NOT_MET"

  Scenario: An unknown coupon is rejected
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    And the customer uses the coupon "FREESTUFF"
    When the customer checks out
    Then checkout fails with the error "UNKNOWN_COUPON"
