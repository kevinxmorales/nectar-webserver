package http

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
	"net/mail"
	"net/url"
)

type UserService interface {
	GetUser(context.Context, string) (user.User, error)
	GetUserByEmail(context.Context, string) (user.User, error)
	AddUser(context.Context, user.User) (user.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, user.User) (user.User, error)
}

type PostUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func convertUserRequestToUser(u PostUserRequest) user.User {
	return user.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var userRequest PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		log.Error(err)
		http.Error(w, "unable to decode request", http.StatusInternalServerError)
		return
	}
	validate := validator.New()
	err := validate.Struct(userRequest)
	if err != nil {
		log.Error(err)
		http.Error(w, "not a valid user object", http.StatusBadRequest)
		return
	}
	if !isValidEmail(userRequest.Email) {
		http.Error(w, "not a valid email address format", http.StatusBadRequest)
		return
	}
	convertedUser := convertUserRequestToUser(userRequest)
	insertedUser, err := h.UserService.AddUser(r.Context(), convertedUser)
	if err != nil {
		log.Error(err)
		return
	}
	if err := json.NewEncoder(w).Encode(insertedUser); err != nil {
		panic(err)
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Info("in GetUser", id)
	usr, err := h.UserService.GetUser(r.Context(), id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("after service call")
	if err := json.NewEncoder(w).Encode(usr); err != nil {
		panic(err)
	}
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encodedEmail := vars["email"]
	if encodedEmail == EMPTY {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	email, err := url.QueryUnescape(encodedEmail)
	if !isValidEmail(email) {
		http.Error(w, "not a valid email address format", http.StatusBadRequest)
		return
	}
	usr, err := h.UserService.GetUser(r.Context(), email)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(usr); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var usr user.User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		log.Error(err)
		http.Error(w, "unable to decode request", http.StatusInternalServerError)
	}
	usr, err := h.UserService.UpdateUser(r.Context(), id, usr)
	if err != nil {
		log.Error(err)
		http.Error(w, "unable to update user", http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(usr); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.UserService.DeleteUser(r.Context(), id)
	if err != nil {
		log.Error(err)
		http.Error(w, "unable to delete user", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(Response{Message: "successfully deleted user"}); err != nil {
		panic(err)
	}
}
