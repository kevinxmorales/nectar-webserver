package health

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type Store interface {
	CheckDbHealth(ctx context.Context) error
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) CheckDbHealth(ctx context.Context) error {
	log.Info("Retrieving database health")
	return s.Store.CheckDbHealth(ctx)
}
