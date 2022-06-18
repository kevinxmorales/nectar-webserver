package plant

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

// Plant - a representation of a plant
type Plant struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	UserId     string      `json:"userId"`
	CategoryID string      `json:"categoryId"`
	Images     []ImageUrls `json:"images"`
	CreatedAt  time.Time   `json:"createdAt"`
}

type ImageUrls struct {
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

//Store - this interface defines all the methods
// the service needs in order to operate
type Store interface {
	GetPlant(context.Context, string) (*Plant, error)
	GetPlantsByUserId(context.Context, string) ([]Plant, error)
	AddPlant(context.Context, Plant) (*Plant, error)
	DeletePlant(context.Context, string) error
	UpdatePlant(context.Context, string, Plant) (*Plant, error)
}

// Service - is the struct on which out logic will
// be built upon
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetPlant(ctx context.Context, id string) (*Plant, error) {
	log.Info("Retrieving a plant with id: ", id)
	return s.Store.GetPlant(ctx, id)
}

func (s *Service) GetPlantsByUserId(ctx context.Context, id string) ([]Plant, error) {
	return s.Store.GetPlantsByUserId(ctx, id)
}

func (s *Service) UpdatePlant(ctx context.Context, ID string, updatedPlant Plant) (*Plant, error) {
	return s.Store.UpdatePlant(ctx, ID, updatedPlant)
}

func (s *Service) DeletePlant(ctx context.Context, id string) error {
	return s.Store.DeletePlant(ctx, id)
}

func (s *Service) PostPlant(ctx context.Context, newPlant Plant) (*Plant, error) {
	log.Info("attempting to add a new plant")
	return s.Store.AddPlant(ctx, newPlant)
}
