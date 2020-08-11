package domain

import "math"

type DiscountType int

const (
	None            DiscountType = iota
	BuyTwoGetOne    DiscountType = iota
	Reduce25Percent DiscountType = iota
)

// Product
type Product struct {
	Code     string
	Name     string
	Price    *Money
	Discount DiscountType
}

var Euro = Currency{
	Symbol:   "€",
	Decimals: 2,
}

func (p *Product) GetDiscountPriceIfAny(units int64) int64 {
	switch {
	case p.Discount == Reduce25Percent && units >= 3:
		// This adds a little overhead (Pow10 returns a float64) but it is necessary to keep the precision
		decimals := int64(math.Pow10(p.Price.Currency.Decimals))
		return (p.Price.Amount*units*decimals - p.Price.Amount*decimals*units/4) / decimals
	case p.Discount == BuyTwoGetOne:
		return p.Price.Amount * (units/2 + units%2)
	}
	return units * p.Price.Amount
}

// The product list should be provided by a database,
// but as we don't have one they have to ber hardcoded somewhere, namely here
func GetProducts() map[string]*Product {
	tenPowerDecimals := math.Pow10(Euro.Decimals)

	productMap := map[string]*Product{
		"PEN": {
			Code:     "PEN",
			Name:     "Lana Pen",
			Price:    &Money{Amount: int64(5 * tenPowerDecimals), Currency: &Euro},
			Discount: BuyTwoGetOne,
		},
		"TSHIRT": {
			Code:     "TSHIRT",
			Name:     "Lana T-Shirt",
			Price:    &Money{Amount: int64(20 * tenPowerDecimals), Currency: &Euro},
			Discount: Reduce25Percent,
		},
		"MUG": {
			Code:     "MUG",
			Name:     "Lana Coffee Mug",
			Price:    &Money{Amount: int64(15 * tenPowerDecimals / 2), Currency: &Euro},
			Discount: None,
		},
	}

	return productMap
}
