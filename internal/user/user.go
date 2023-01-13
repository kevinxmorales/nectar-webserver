package user

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/validation"
	"strings"
)

type NewUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	ImageUrl string `json:"imageUrl" validate:"required"`
}

type User struct {
	Id         string   `json:"id"`
	PlantCount uint     `json:"plantCount"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
	ImageUrl   string   `json:"image_url"`
	Following  []string `json:"following"`
}

type Store interface {
	GetUser(ctx context.Context, id string) (*User, error)
	GetUserById(ctx context.Context, id string) (*User, error)
	AddUser(ctx context.Context, u User) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, u User) (*User, error)
	CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error)
}

type AuthClient interface {
	CreateUser(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error)
}

type Service struct {
	Store      Store
	AuthClient AuthClient
}

// NewService - returns a pointer to a new user service
func NewService(store Store, authClient AuthClient) *Service {
	return &Service{
		Store:      store,
		AuthClient: authClient,
	}
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.Store.GetUser(ctx, id)
}

func (s *Service) GetUserById(ctx context.Context, id string) (*User, error) {
	tag := "user.GetUserById"
	u, err := s.Store.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Store.GetUserByAuthId in %s failed for %v", tag, err)
	}
	u.Following = []string{}
	return u, nil
}

func (s *Service) AddUser(ctx context.Context, u NewUserRequest) (*User, error) {
	tag := "user.AddUser"
	if err := validation.IsValidEmail(u.Email); err != nil {
		return nil, fmt.Errorf("email validation for email %s in %s, failed for %v", u.Email, tag, err)
	}
	newUserId := uuid.NewV4().String()
	params := (&auth.UserToCreate{}).
		UID(newUserId).
		Email(u.Email).
		EmailVerified(false).
		Password(u.Password).
		Disabled(false)
	if _, err := s.AuthClient.CreateUser(ctx, params); err != nil {
		return nil, err
	}
	newUser := User{
		Id:       newUserId,
		Name:     u.Name,
		Email:    u.Email,
		Username: strings.Split(u.Email, "@")[0],
	}
	return s.Store.AddUser(ctx, newUser)
}
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.Store.DeleteUser(ctx, id)

}
func (s *Service) UpdateUser(ctx context.Context, id string, usr UpdateUserRequest) (*User, error) {
	u := User{
		Id:       id,
		Username: usr.Username,
		Name:     usr.Name,
		ImageUrl: usr.ImageUrl,
		Email:    usr.Email,
	}
	return s.Store.UpdateUser(ctx, id, u)
}

func (s *Service) CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error) {
	return s.Store.CheckIfUsernameIsTaken(ctx, username)
}
