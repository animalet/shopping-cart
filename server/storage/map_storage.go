package storage

import (
	"shopping-cart/server/domain"
	"sync"
)

type mapStorage struct {
	backingMap map[string]*domain.Cart
	productMap map[string]*domain.Product
	mutex      *sync.RWMutex
}

// This storage implementation manages the carts in an in-memory map
func NewMapStorage() Storage {
	backingMap := make(map[string]*domain.Cart)
	return &mapStorage{
		backingMap: backingMap,
		productMap: domain.GetProducts(),
		mutex:      &sync.RWMutex{},
	}
}

func (mapStorage *mapStorage) DeleteCart(cartId string) bool {
	mapStorage.mutex.Lock()
	defer mapStorage.mutex.Unlock()
	_, exists := mapStorage.backingMap[cartId]
	if exists {
		delete(mapStorage.backingMap, cartId)
	}

	return exists
}

func (mapStorage *mapStorage) AddProductToCartIfAny(cartId, productCode string) (int, bool) {
	mapStorage.mutex.RLock()
	defer mapStorage.mutex.RUnlock()
	_, productExists := mapStorage.productMap[productCode]
	cart, cartExists := mapStorage.backingMap[cartId]
	currentAmount := 0
	if cartExists && productExists {
		if amount, existsItem := cart.Items[productCode]; !existsItem {
			currentAmount = 1
		} else {
			currentAmount = amount + 1
		}
		cart.Items[productCode] = currentAmount
	}
	return currentAmount, cartExists && productExists
}

func (mapStorage *mapStorage) GetCartContents(cartId string) (map[string]int, bool) {
	mapStorage.mutex.RLock()
	defer mapStorage.mutex.RUnlock()

	cart, exists := mapStorage.backingMap[cartId]
	if exists {
		return cart.Items, true
	}
	return nil, false
}

func (mapStorage *mapStorage) Store(cart *domain.Cart) {
	mapStorage.mutex.Lock()
	defer mapStorage.mutex.Unlock()
	mapStorage.backingMap[*cart.Id] = cart
}

func (mapStorage *mapStorage) GetProducts() map[string]*domain.Product {
	return mapStorage.productMap
}
