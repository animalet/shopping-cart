package storage

import "shopping-cart/server/domain"

type Storage interface {
	DeleteCart(cartId string) bool
	AddProductToCartIfAny(cartId, productCode string) (int, bool)
	GetCartContents(cartId string) (map[string]int, bool)
	Store(cart *domain.Cart)
	GetProducts() map[string]*domain.Product
}
