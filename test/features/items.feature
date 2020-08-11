Feature: Buying items

  Scenario Outline:
    Given an online cart service
    And an existing empty cart
    When purchasing the following items: <items>
    Then the amount <total> matches
    Examples:
      | items                                      | total  |
      | PEN, TSHIRT, MUG                           | 32.50€ |
      | PEN, TSHIRT, PEN                           | 25.00€ |
      | TSHIRT, TSHIRT, TSHIRT, PEN, TSHIRT        | 65.00€ |
      | PEN, TSHIRT, PEN, PEN, MUG, TSHIRT, TSHIRT | 62.50€ |

  Scenario:
    Given an online cart service
    And an existing empty cart
    Then the amount 0.00€ matches
