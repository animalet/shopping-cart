package server

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"shopping-cart/server/domain"
	"shopping-cart/server/storage"
)

// This is the service with the business logic
type CartService struct {
	Storage storage.Storage
}

func (m *CartService) Create() string {
	var (
		cartId  = uuid.New().String()
		itemMap = make(map[string]int)
		cart    = domain.Cart{Id: &cartId, Items: itemMap}
	)
	m.Storage.Store(&cart)
	log.Printf("Cart created: %s", cartId)
	return cartId
}

func (m *CartService) Add(cartId, code string) bool {
	currentAmount, found := m.Storage.AddProductToCartIfAny(cartId, code)
	if found {
		log.Printf("Product %s added to cart (current amount is %d): %s", code, currentAmount, cartId)
	}
	return found
}

func (m *CartService) Remove(cartId string) bool {
	found := m.Storage.DeleteCart(cartId)
	if found {
		log.Printf("Cart removed: %s", cartId)
	}
	return found
}

func (m *CartService) Total(cartId string) (*domain.Money, bool) {
	var (
		items, found       = m.Storage.GetCartContents(cartId)
		total        int64 = 0
	)

	var currency *domain.Currency
	if found {
		products := m.Storage.GetProducts()
		for k, v := range items {
			product := products[k]
			currency = mergeCurrencies(currency, product.Price.Currency)
			total += product.GetDiscountPriceIfAny(int64(v))
		}
		log.Printf("Total retrieved for cart %s: %d", cartId, total)
	}
	if currency == nil {
		currency = &domain.Euro // Default currency is Euro. Needed due to multicurrency lack of support
	}
	return &domain.Money{
		Amount:   total,
		Currency: currency,
	}, found
}

func mergeCurrencies(currency *domain.Currency, productCurrency *domain.Currency) *domain.Currency {
	if currency == nil {
		currency = productCurrency
	} else if currency != productCurrency {
		// Check for multicurrency: a good implementation should convert currencies at this point, instead of failing
		panic(fmt.Sprintf("Mismatching currencies when summing totals: %v, %v", currency, productCurrency))
	}
	return currency
}
