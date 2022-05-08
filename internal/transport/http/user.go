package http

import (
	"context"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
)

type UserService interface {
	GetUser(context.Context, string) (user.User, error)
	AddUser(context.Context, user.User) (user.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, user.User) (user.User, error)
}

type PostUserRequest struct {
	Name string `json:"name" validate:"required"`
}

func convertUserRequestToUser(u PostUserRequest) user.User {
	return user.User{}
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	return
}
