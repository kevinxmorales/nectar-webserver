package plant

import (
	"context"
	"errors"
	"fmt"
)

type Category struct {
	Color string
	Icon  string
	Label string
	ID    uint
}

// Plant - a representation of a plant
type Plant struct {
	ID         string
	Name       string
	Images     []string
	UserId     string
	Categories []Category
}

var (
	ErrFetchingPlant  = errors.New("failed to fetch plant by id")
	ErrNotImplemented = errors.New("not implemented")
)

//Store - this interface defines all the methods
// the service needs in order to operate
type Store interface {
	GetPlant(context.Context, string) (Plant, error)
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

func (s *Service) GetPlant(ctx context.Context, id string) (Plant, error) {
	fmt.Println("Retrieving a plant with id: ", id)
	plant, err := s.Store.GetPlant(ctx, id)
	if err != nil {
		fmt.Println(err)
		return Plant{}, ErrFetchingPlant
	}
	return plant, nil
}

func (s *Service) UpdatePlant(ctx context.Context, updatedPlant Plant) error {
	return ErrNotImplemented
}

func (s *Service) DeletePlant(ctx context.Context, id string) error {
	return ErrNotImplemented
}

func (s *Service) CreatePlant(ctx context.Context, newPlant Plant) (Plant, error) {
	return Plant{}, ErrNotImplemented
}
