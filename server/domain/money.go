package domain

import (
	"fmt"
	"strconv"
)

type Money struct {
	Amount   int64
	Currency *Currency
}

type Currency struct {
	Symbol   string
	Decimals int
}

// Converts money to string format considering its decimal and integer positions and appending the symbol
func (m *Money) String() string {
	str := []rune(fmt.Sprintf("%0"+strconv.Itoa(m.Currency.Decimals+1)+"d%s", m.Amount, m.Currency.Symbol))
	length := len(str)
	formatted := fmt.Sprintf("%s.%s", string(str[0:length-m.Currency.Decimals-1]), string(str[length-1-m.Currency.Decimals:length]))
	return formatted
}
