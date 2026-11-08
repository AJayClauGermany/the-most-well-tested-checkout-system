Feature: Product promotions
  Some products have a promotion that lowers their price automatically.
  COFFEE is "buy 2, get 1 free" and NOTEBOOK is "10% off when you buy 10 or more".

  Background:
    Given the standard catalog

  Scenario Outline: Buy 2 coffees, get 1 free
    Given a basket containing:
      | sku    | quantity   |
      | COFFEE | <quantity> |
    When the customer checks out
    Then the promotion discount is $<discount>

    Examples:
      | quantity | discount |
      | 1        | 0.00     |
      | 2        | 0.00     |
      | 3        | 8.99     |
      | 5        | 8.99     |
      | 6        | 17.98    |

  Scenario Outline: Bulk discount on notebooks
    Given a basket containing:
      | sku      | quantity   |
      | NOTEBOOK | <quantity> |
    When the customer checks out
    Then the promotion discount is $<discount>

    Examples:
      | quantity | discount |
      | 9        | 0.00     |
      | 10       | 3.49     |
      | 11       | 3.84     |
      | 15       | 5.24     |

  Scenario: Coffee on separate lines still counts towards the promotion
    Given a basket containing:
      | sku    | quantity |
      | COFFEE | 2        |
      | COFFEE | 1        |
    When the customer checks out
    Then the promotion discount is $8.99

  Scenario: Products without a promotion are never discounted
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 10       |
      | PEN   | 10       |
    When the customer checks out
    Then the promotion discount is $0.00
