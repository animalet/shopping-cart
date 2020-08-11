package storage

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	consul "github.com/hashicorp/consul/api"
	"log"
	"os"
	"shopping-cart/server/domain"
)

type redisStorage struct {
	rdb        redis.UniversalClient
	productMap map[string]*domain.Product
}

// This storage implementation manages the carts in a redis server.
func NewRedisStorage(consulAddr string) Storage {
	return &redisStorage{getRedisClientIfScalabilityEnabled(consulAddr), domain.GetProducts()}
}

func (r *redisStorage) DeleteCart(cartId string) bool {
	if removedNum, err := r.rdb.Del(cartId).Result(); err != nil {
		panic(fmt.Sprintf("Error while deleting a key in Redis %s: %s", r.rdb, err))
	} else {
		return removedNum > 0
	}
}

func (r *redisStorage) AddProductToCartIfAny(cartId, productCode string) (int, bool) {
	currentAmount := 0
	_, productExists := r.productMap[productCode]
	cartExists := false

	// We use optimistic locking because it is unlikely that
	// two threads access the same cart at the same time from the same client
	err := r.rdb.Watch(func(tx *redis.Tx) error {
		get, err := tx.Get(cartId).Bytes()
		if err != nil {
			return err
		} else {
			cartExists = true
		}

		items := make(map[string]int)
		if err = json.Unmarshal(get, &items); err != nil {
			return err
		}
		if amount, existsItem := items[productCode]; !existsItem {
			currentAmount = 1
		} else {
			currentAmount = amount + 1
		}
		items[productCode] = currentAmount
		_, err = tx.Pipelined(func(pipeliner redis.Pipeliner) error {
			if b, err := json.Marshal(items); err == nil {
				pipeliner.Set(cartId, b, 0)
				return nil
			} else {
				return err
			}
		})
		return err
	}, cartId)
	if err != nil {
		panic(fmt.Sprintf("Could not update cart in Redis %s: %s", r.rdb, err))
	}
	return currentAmount, cartExists && productExists
}

func (r *redisStorage) GetCartContents(cartId string) (map[string]int, bool) {
	if get, err := r.rdb.Get(cartId).Bytes(); err != nil && err != redis.Nil {
		panic(fmt.Sprintf("Error geting cart %s from Redis %s: %s", err, r.rdb, err))
	} else {
		if len(get) == 0 {
			return nil, false
		}
		var items map[string]int
		if err := json.Unmarshal(get, &items); err != nil {
			panic(fmt.Sprintf("Error decoding data from cart %s in Redis %s: %s", cartId, r.rdb, err))
		} else {
			return items, true
		}
	}
}

func (r *redisStorage) Store(cart *domain.Cart) {
	if b, err := json.Marshal(cart.Items); err == nil {
		r.rdb.Set(*cart.Id, b, -1)
	} else {
		panic(fmt.Sprintf("Could not store cart in Redis: %s", err))
	}
}

func (r *redisStorage) GetProducts() map[string]*domain.Product {
	return r.productMap
}

// If running in an scalable environment (Consul and Redis available) we try to instance Redis through Consul.
// Environment variable CONSUL_ADDR expected for this to happen
func getRedisClientIfScalabilityEnabled(consulAddr string) redis.UniversalClient {
	redisName, found := os.LookupEnv("REDIS_NAME")
	if !found {
		redisName = "redis"
	}

	client, err := consul.NewClient(&consul.Config{
		Address: consulAddr,
		Scheme:  "http",
	})
	if err != nil {
		panic(fmt.Sprintf("Couldn't instantiate Consul client: %s", err))
	}

	if catalog, _, err := client.Catalog().Service(redisName, "", &consul.QueryOptions{}); err != nil {
		log.Printf("Error occured when resolving redis service in consul: %s", err)
	} else {
		redisAddr := make([]string, 0)
		for _, node := range catalog {
			redisAddr = append(redisAddr, fmt.Sprintf("%s:%d", node.ServiceAddress, node.ServicePort))
		}
		rdb := redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: redisAddr,
		})

		if pong, err := rdb.Ping().Result(); err == nil {
			log.Printf("Connected to Redis: %s", pong)
			return rdb
		} else {
			log.Printf("Couldn't connect to Redis: %s", err)
		}
	}

	panic("Couldn't connect to Redis")
}
