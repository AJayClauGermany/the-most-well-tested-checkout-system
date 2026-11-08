Feature: Saving orders
  After a successful checkout the order can be saved and found again later.

  Background:
    Given the standard catalog
    And an empty order repository

  Scenario: Saving an order gives it an id
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    When the customer checks out and saves the order
    Then the order is saved with the id "ORDER-1"
    And the saved order "ORDER-1" has a total of $5.53

  Scenario: Each saved order gets the next id
    Given a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    When the customer checks out and saves the order
    And the customer checks out and saves the order
    Then the order is saved with the id "ORDER-2"

  Scenario: A failed checkout saves nothing
    Given an empty basket
    When the customer checks out and saves the order
    Then checkout fails with the error "EMPTY_BASKET"
    And the order repository is empty

  Scenario: A saved order remembers when it was saved
    Given the clock is fixed at "2026-01-15T10:30:00Z"
    And a basket containing:
      | sku   | quantity |
      | APPLE | 1        |
    When the customer checks out and saves the order
    Then the saved order "ORDER-1" was created at "2026-01-15T10:30:00Z"

  Scenario: Looking up an order that does not exist
    Then the order "ORDER-99" cannot be found
