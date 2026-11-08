Feature: Rejecting invalid baskets
  When a basket is invalid, checkout stops with an error code.
  If several things are wrong, the first problem in the rules' check order is reported.

  Background:
    Given the standard catalog

  Scenario: An empty basket
    Given an empty basket
    When the customer checks out
    Then checkout fails with the error "EMPTY_BASKET"

  Scenario: A product that is not in the catalog
    Given a basket containing:
      | sku    | quantity |
      | BANANA | 1        |
    When the customer checks out
    Then checkout fails with the error "UNKNOWN_PRODUCT"

  Scenario Outline: Quantities must be from 1 to 99
    Given a basket containing:
      | sku   | quantity   |
      | APPLE | <quantity> |
    When the customer checks out
    Then checkout fails with the error "INVALID_QUANTITY"

    Examples:
      | quantity |
      | 0        |
      | -1       |
      | 100      |

  Scenario: Merged lines may not go over 99
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 50       |
      | APPLE | 50       |
    When the customer checks out
    Then checkout fails with the error "INVALID_QUANTITY"

  Scenario: Unknown products are reported before bad quantities
    Given a basket containing:
      | sku    | quantity |
      | APPLE  | 0        |
      | BANANA | 1        |
    When the customer checks out
    Then checkout fails with the error "UNKNOWN_PRODUCT"
