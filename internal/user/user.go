package user

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Store interface {
	GetUser(context.Context, string) (User, error)
	GetUserByEmail(context.Context, string) (User, error)
	AddUser(context.Context, User) (User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, User) (User, error)
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

func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
	log.Info("in Service.GetUser")
	usr, err := s.Store.GetUser(ctx, id)
	if err != nil {
		log.Info("exiting Service.GetUser...")
		return User{}, err
	}
	log.Info("exiting Service.GetUser...")
	return usr, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (User, error) {
	usr, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	return usr, nil
}

func (s *Service) AddUser(ctx context.Context, usr User) (User, error) {
	newUser, err := s.Store.AddUser(ctx, usr)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.Store.DeleteUser(ctx, id)

}
func (s *Service) UpdateUser(ctx context.Context, id string, usr User) (User, error) {
	updatedUser, err := s.Store.UpdateUser(ctx, id, usr)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}
