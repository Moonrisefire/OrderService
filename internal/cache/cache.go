package cache

import (
	"time"

	"OrderService/internal/models"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	c *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		c: cache.New(10*time.Minute, 15*time.Minute),
	}
}

func (cc *Cache) Get(id string) (models.Order, bool) {
	data, found := cc.c.Get(id)
	if !found {
		return models.Order{}, false
	}
	order, ok := data.(models.Order)
	return order, ok
}

func (cc *Cache) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			cc.c.DeleteExpired()
		}
	}()
}

func (cc *Cache) Set(id string, order models.Order) {
	cc.c.Set(id, order, cache.DefaultExpiration)
}

func (cc *Cache) Delete(id string) {
	cc.c.Delete(id)
}

func (cc *Cache) Flush() {
	cc.c.Flush()
}
