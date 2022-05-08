package user

import (
	"context"
	"errors"
)

var errNotImplemented = errors.New("not implemented")

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Store interface {
	GetUser(context.Context, string) (User, error)
	AddUser(context.Context, User) (User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, User) (User, error)
}

type UserService struct {
	Store Store
}

// NewService - returns a pointer to a new user service
func NewService(store Store) *UserService {
	return &UserService{
		Store: store,
	}
}

func (s *UserService) GetUser(context.Context, string) (User, error) {
	return User{}, errNotImplemented
}
func (s *UserService) AddUser(context.Context, User) (User, error) {
	return User{}, errNotImplemented
}
func (s *UserService) DeleteUser(context.Context, string) error {
	return errNotImplemented
}
func (s *UserService) UpdateUser(context.Context, string, User) (User, error) {
	return User{}, errNotImplemented
}
