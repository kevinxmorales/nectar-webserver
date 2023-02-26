package health

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type Store interface {
	CheckDbHealth(ctx context.Context) error
}

type Cache interface {
	CheckCacheHealth(ctx context.Context) error
}

type Service struct {
	Store Store
	Cache Cache
}

func NewService(store Store, cache Cache) *Service {
	return &Service{
		Store: store,
		Cache: cache,
	}
}

func (s *Service) CheckDbHealth(ctx context.Context) error {
	log.Info("Retrieving database health")
	return s.Store.CheckDbHealth(ctx)
}

func (s *Service) CheckCacheHealth(ctx context.Context) error {
	log.Info("Retrieving database health")
	return s.Cache.CheckCacheHealth(ctx)
}
