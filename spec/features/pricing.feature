Feature: Pricing a basket
  As a shopper
  I want the checkout to add up what is in my basket
  So that I know how much I will pay

  Background:
    Given the standard catalog

  Scenario: A single item
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    When the customer checks out
    Then the subtotal is $0.50
    And the shipping cost is $4.99
    And the tax is $0.04
    And the total is $5.53

  Scenario: Quantity multiplies the unit price
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 3        |
    When the customer checks out
    Then the subtotal is $1.50

  Scenario: Different products are added together
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 2        |
      | BREAD | 1        |
      | PEN   | 3        |
    When the customer checks out
    Then the subtotal is $6.85
    And the total is $12.39

  Scenario: The same product on two lines is merged into one line
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
      | PEN   | 2        |
      | APPLE | 3        |
    When the customer checks out
    Then the receipt has 2 lines
    And the receipt line for "APPLE" has a quantity of 4
    And the subtotal is $4.40
