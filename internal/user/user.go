package user

import (
	"context"
)

type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	PlantCount uint   `json:"plantCount"`
}

type Store interface {
	GetUser(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	AddUser(context.Context, User) (*User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, User) (*User, error)
}

type Service struct {
	Store Store
}

// NewService - returns a pointer to a new user service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.Store.GetUser(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.Store.GetUserByEmail(ctx, email)
}

func (s *Service) AddUser(ctx context.Context, usr User) (*User, error) {
	return s.Store.AddUser(ctx, usr)
}
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.Store.DeleteUser(ctx, id)

}
func (s *Service) UpdateUser(ctx context.Context, id string, usr User) (*User, error) {
	return s.Store.UpdateUser(ctx, id, usr)
}
