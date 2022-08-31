package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/db"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
	"strconv"
)

type UserService interface {
	GetUser(ctx context.Context, id int) (*user.User, error)
	GetUserByAuthId(ctx context.Context, firebaseId string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	AddUser(ctx context.Context, u user.User) (*user.User, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, id int, u user.User) (*user.User, error)
	CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error)
}

type uniqueUsernameResponse struct {
	Username string `json:"username"`
	IsTaken  bool   `json:"isTaken"`
}

type PostUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	AuthId   string `json:"authId" validate:"required"`
}

func convertUserRequestToUser(u PostUserRequest) user.User {
	return user.User{
		Name:     u.Name,
		Email:    u.Email,
		Username: u.Username,
		AuthId:   u.AuthId,
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
		if errors.Is(err, db.DuplicateKeyError) {
			h.SendForbiddenResponse(w, r, err)
			return
		}
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, insertedUser)
	return
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	usr, err := h.UserService.GetUser(r.Context(), id)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, usr)
	return
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	usr, err := h.UserService.GetUserByEmail(r.Context(), email)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, usr)
	return
}

func (h *Handler) GetUserByAuthId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firebaseId := vars["id"]
	usr, err := h.UserService.GetUserByAuthId(r.Context(), firebaseId)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	var message string
	if usr == nil {
		message = "No user found with the given id"
	}
	resEntity := &ResponseEntity{
		Content:    usr,
		Messages:   []string{message},
		HttpStatus: http.StatusOK,
	}
	h.SendOkResponse(w, r, resEntity)
	return
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	var usr user.User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	usr.Id = id
	usrUpdated, err := h.UserService.UpdateUser(r.Context(), id, usr)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, usrUpdated)
	return
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	if err := h.UserService.DeleteUser(r.Context(), id); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	res := Response{Message: "successfully deleted user"}
	h.SendOkResponse(w, r, res)
	return
}

func (h *Handler) CheckIfUsernameIsTaken(w http.ResponseWriter, r *http.Request) {
	usernameParam := "username"
	params, err := h.ParseUrlQueryParams(r.URL, usernameParam)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	username := params[usernameParam]
	isTaken, err := h.UserService.CheckIfUsernameIsTaken(r.Context(), username)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	res := &uniqueUsernameResponse{
		Username: username,
		IsTaken:  isTaken,
	}
	resEntity := ResponseEntity{
		Content:    res,
		HttpStatus: http.StatusOK,
		Messages:   nil,
	}
	h.SendOkResponse(w, r, resEntity)
	return
}
