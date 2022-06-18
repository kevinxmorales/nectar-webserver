package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
	"net/mail"
	"net/url"
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

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var userRequest PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		sendBadRequestResponse(w, r, errors.New("unable to decode request"))
		return
	}
	validate := validator.New()
	err := validate.Struct(userRequest)
	if err != nil {
		sendBadRequestResponse(w, r, errors.New("not a valid user object"))
		return
	}
	if !isValidEmail(userRequest.Email) {
		sendBadRequestResponse(w, r, errors.New("not a valid email address format"))
		return
	}
	convertedUser := convertUserRequestToUser(userRequest)
	insertedUser, err := h.UserService.AddUser(r.Context(), convertedUser)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, insertedUser)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	usr, err := h.UserService.GetUser(r.Context(), id)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, usr)
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encodedEmail := vars["email"]
	if encodedEmail == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	email, err := url.QueryUnescape(encodedEmail)
	if !isValidEmail(email) {
		sendBadRequestResponse(w, r, errors.New("not a valid email address format"))
		return
	}
	usr, err := h.UserService.GetUser(r.Context(), email)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, usr)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	var usr user.User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		send500Response(w, r, err)
		return
	}
	usrUpdated, err := h.UserService.UpdateUser(r.Context(), id, usr)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, usrUpdated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	err := h.UserService.DeleteUser(r.Context(), id)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	res := Response{Message: "successfully deleted user"}
	sendOkResponse(w, r, res)
}
