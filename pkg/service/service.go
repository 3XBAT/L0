package service

import (
	"L0"
	"L0/pkg/repository"

	"github.com/sirupsen/logrus"
)

type Cache interface {
	NewCache() error
	AddOrder(orderUID string, order L0.Order)
	GetCache(uid string) (L0.Order, error)
}

type Service struct {
	Cache
}

func NewService(repo *repository.Repository) (*Service, error) {
	cacheServcie, err := NewCacheService(repo)
	if err != nil {
		logrus.Errorf("Error with cacheService : %s", err.Error())
	}
	return &Service{Cache: cacheServcie}, nil
}


