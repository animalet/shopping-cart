Feature: Cart management
  Scenario:
    Given an online cart service
    And a non existing cart
    When purchasing any items
    Then cart is not found

  Scenario:
    Given an online cart service
    And an existing not empty cart
    When the cart is removed
    And purchasing any items
    Then cart is not found
