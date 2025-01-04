package cache

import (
	"nqrm/wbtechlvl0/order_services/internal/model"
	"sync"
)

type CacheStorage struct {
	mu    sync.RWMutex
	cache map[string]*model.Order
}

func NewCacheStorage() *CacheStorage {
	orderCache := &CacheStorage{cache: make(map[string]*model.Order)}

	return orderCache
}

func (c *CacheStorage) Get(orderUID string) (*model.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	order, ok := c.cache[orderUID]
	return order, ok
}

func (c *CacheStorage) Set(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[order.OrderUID] = order
}
