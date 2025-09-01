package cache

import (
	"sync"
	"tech-wb-L0/backend/domain"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]domain.Order
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]domain.Order),
	}
}

func (c *Cache) Get(orderUID string) (domain.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.store[orderUID]
	return order, ok
}

func (c *Cache) Set(orderUID string, order domain.Order) {
	c.mu.RLock()
	c.mu.RUnlock()
	c.store[orderUID] = order
}

func (c *Cache) RangeMap() int {
	return len(c.store)
}
