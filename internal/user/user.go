package user

import (
	"context"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/validation"
	"net/url"
)

type User struct {
	Id         int    `json:"id"`
	PlantCount uint   `json:"plantCount"`
	Name       string `json:"name"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Username   string `json:"username"`
	AuthId     string `json:"authId"`
	ImageUrl   string `json:"image_url"`
	Following  []int  `json:"following"`
}

type Store interface {
	GetUser(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByAuthId(ctx context.Context, firebaseId string) (*User, error)
	AddUser(ctx context.Context, u User) (*User, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, id int, u User) (*User, error)
	CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error)
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

func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
	return s.Store.GetUser(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, encodedEmail string) (*User, error) {
	tag := "user.GetUserByEmail"
	email, err := url.QueryUnescape(encodedEmail)
	if err != nil {
		return nil, fmt.Errorf("url.QueryUnescape in %s failed for %v", tag, err)
	}
	if err := validation.IsValidEmail(email); err != nil {
		return nil, fmt.Errorf("email validation for email %s in %s failed for %v", email, tag, err)
	}
	return s.Store.GetUserByEmail(ctx, email)
}

func (s *Service) GetUserByAuthId(ctx context.Context, firebaseId string) (*User, error) {
	tag := "user.GetUserByAuthId"
	u, err := s.Store.GetUserByAuthId(ctx, firebaseId)
	if err != nil {
		return nil, fmt.Errorf("Store.GetUserByAuthId in %s failed for %v", tag, err)
	}
	u.Following = []int{}
	return u, nil
}

func (s *Service) AddUser(ctx context.Context, usr User) (*User, error) {
	tag := "user.AddUser"
	if err := validation.IsValidEmail(usr.Email); err != nil {
		return nil, fmt.Errorf("email validation for email %s in %s, failed for %v", usr.Email, tag, err)
	}
	return s.Store.AddUser(ctx, usr)
}
func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.Store.DeleteUser(ctx, id)

}
func (s *Service) UpdateUser(ctx context.Context, id int, usr User) (*User, error) {
	return s.Store.UpdateUser(ctx, id, usr)
}

func (s *Service) CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error) {
	return s.Store.CheckIfUsernameIsTaken(ctx, username)
}
