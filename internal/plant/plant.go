package plant

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Category struct {
	Color string `json:"color"`
	Icon  string `json:"icon"`
	Label string `json:"label"`
	ID    uint   `json:"id"`
}

// Plant - a representation of a plant
type Plant struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Images     []string   `json:"images"`
	UserId     string     `json:"userId"`
	Categories []Category `json:"categories"`
}

var ErrFetchingPlant = errors.New("failed to fetch plant by id")

//Store - this interface defines all the methods
// the service needs in order to operate
type Store interface {
	GetPlant(context.Context, string) (Plant, error)
	AddPlant(context.Context, Plant) (Plant, error)
	DeletePlant(context.Context, string) error
	UpdatePlant(context.Context, string, Plant) (Plant, error)
}

// PlantService - is the struct on which out logic will
// be built upon
type PlantService struct {
	Store Store
}

// NewService - returns a pointer to a new service
func NewService(store Store) *PlantService {
	return &PlantService{
		Store: store,
	}
}

func (s *PlantService) GetPlant(ctx context.Context, id string) (Plant, error) {
	log.Info("Retrieving a plant with id: ", id)
	plant, err := s.Store.GetPlant(ctx, id)
	if err != nil {
		fmt.Println(err)
		return Plant{}, ErrFetchingPlant
	}
	return plant, nil
}

func (s *PlantService) UpdatePlant(ctx context.Context, ID string, updatedPlant Plant) (Plant, error) {
	plant, err := s.Store.UpdatePlant(ctx, ID, updatedPlant)
	if err != nil {
		log.Error("error updating plant")
		return Plant{}, err
	}
	return plant, nil
}

func (s *PlantService) DeletePlant(ctx context.Context, id string) error {
	return s.Store.DeletePlant(ctx, id)
}

func (s *PlantService) PostPlant(ctx context.Context, newPlant Plant) (Plant, error) {
	log.Info("attempting to add a new plant")
	insertedPlant, err := s.Store.AddPlant(ctx, newPlant)
	if err != nil {
		return Plant{}, err
	}
	return insertedPlant, nil
}
