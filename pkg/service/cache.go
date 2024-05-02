package service

import (
	"L0"
	"L0/pkg/repository"
	"errors"
)

type CacheService struct {
	repo  *repository.Repository
	Cache map[string]L0.Order
}

func NewCacheService(repo *repository.Repository) (*CacheService, error) {
	cacheService := &CacheService{repo: repo}
	if err := cacheService.NewCache(); err != nil {
		
		return nil, err

	}
	return cacheService, nil
}

func (c *CacheService) NewCache() error { //заполняет кэш из postgres
	cache := map[string]L0.Order{}

	orders, err := c.repo.RecoverCache()

	if err != nil {
		return err
	}

	for _, order := range orders {
		cache[order.OrderUID] = order
	}
	c.Cache = cache
	return nil
}

func (c *CacheService) AddOrder(orderUID string, order L0.Order) { //просто добавляем заказ в кэш
	
	c.Cache[orderUID] = order
}

func (c *CacheService) GetCache(uid string) (L0.Order, error) { 
	c.repo.RecoverCache()
	order := c.Cache[uid]
	if len(order.OrderUID) == 0 {
		return order, errors.New("no order by this uid")
	}
	return order, nil
}
