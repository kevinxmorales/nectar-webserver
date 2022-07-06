package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
)

type UserService interface {
	GetUser(context.Context, string) (*user.User, error)
	GetUserByEmail(context.Context, string) (*user.User, error)
	AddUser(context.Context, user.User) (*user.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, user.User) (*user.User, error)
}

type PostUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

func convertUserRequestToUser(u PostUserRequest) user.User {
	return user.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var userRequest PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		h.SendBadRequestResponse(w, r, errors.New("unable to decode request"))
		return
	}
	validate := validator.New()
	if err := validate.Struct(userRequest); err != nil {
		h.SendBadRequestResponse(w, r, errors.New("not a valid user object"))
		return
	}
	convertedUser := convertUserRequestToUser(userRequest)
	insertedUser, err := h.UserService.AddUser(r.Context(), convertedUser)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, insertedUser)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		h.SendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	usr, err := h.UserService.GetUser(r.Context(), id)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, usr)
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	usr, err := h.UserService.GetUser(r.Context(), email)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, usr)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		h.SendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	var usr user.User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	usrUpdated, err := h.UserService.UpdateUser(r.Context(), id, usr)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, usrUpdated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		h.SendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	if err := h.UserService.DeleteUser(r.Context(), id); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	res := Response{Message: "successfully deleted user"}
	h.SendOkResponse(w, r, res)
}
