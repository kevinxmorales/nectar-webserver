package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/nectar_errors"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (*user.User, error)
	GetUserById(ctx context.Context, id string) (*user.User, error)
	AddUser(ctx context.Context, u user.NewUserRequest) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, u user.UpdateUserRequest) (*user.User, error)
	UpdateUserProfileImage(ctx context.Context, filePath string, id string) (string, error)
	CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest user.NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	validate := validator.New()
	if err := validate.Struct(userRequest); err != nil {
		res := Response{Message: "Invalid request, please include all required fields"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, res)
		return
	}
	insertedUser, err := h.UserService.AddUser(r.Context(), userRequest)
	if err != nil {
		switch err {
		case nectar_errors.BadRequestError{}:
			res := Response{Message: fmt.Sprintf("Provided email or username is already registered: %s", userRequest.Email)}
			log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
			w.WriteHeader(http.StatusBadRequest)
			h.encodeJsonResponse(&w, res)
			return
		default:
			res := Response{Message: fmt.Sprintf("An unexpected error occurred")}
			log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
			w.WriteHeader(http.StatusInternalServerError)
			h.encodeJsonResponse(&w, res)
			return
		}
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	res := Response{Content: insertedUser, Message: "account successfully created"}
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, res)
	return
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s, status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: fmt.Sprintf("invalid id %s", id)})
		return
	}
	usr, err := h.UserService.GetUser(r.Context(), id)
	if err != nil {
		res := Response{Message: "Unexpected error, could not get user info"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, res)
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	res := Response{Content: usr, Message: "account successfully created"}
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, res)
	return
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s, status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: fmt.Sprintf("invalid id %s", id)})
		return
	}
	usr, err := h.UserService.GetUserById(r.Context(), id)
	if err != nil {
		res := Response{Message: "Unexpected error, could not get user info"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, res)
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	res := Response{Content: usr, Message: "account successfully created"}
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, res)
	return
}

func (h *Handler) UpdateUserProfileImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s, status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: fmt.Sprintf("invalid id %s", id)})
		return
	}
	filePath, err := ParseImageFromRequestBody(r)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "Unexpected error, could not update user profile image"})
		return
	}
	fileUri, err := h.UserService.UpdateUserProfileImage(r.Context(), filePath, id)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "Unexpected error, could not update user profile image"})
		return
	}
	content := struct {
		uri string
	}{uri: fileUri}
	h.encodeJsonResponse(&w, Response{Content: content})
	return
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s, status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: fmt.Sprintf("invalid id %s", id)})
		return
	}
	var usr user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		res := Response{Message: "Unexpected error, could not get user info"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, res)
		return
	}
	usrUpdated, err := h.UserService.UpdateUser(r.Context(), id, usr)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "Unexpected error, could not get user info"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: usrUpdated, Message: "user data successfully updated"})
	return
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s, status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: fmt.Sprintf("invalid id %s", id)})
		return
	}
	if err := h.UserService.DeleteUser(r.Context(), id); err != nil {
		res := Response{Message: "Unexpected error, could not delete user info, please try again"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	res := response{Message: "successfully deleted user"}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) CheckIfUsernameIsTaken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Username string `json:"username"`
		IsTaken  bool   `json:"isTaken"`
	}
	usernameParam := "username"
	params, err := h.ParseUrlQueryParams(r.URL, usernameParam)
	if err != nil {
		res := Response{Message: "Unexpected error"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, res)
		return
	}
	username := params[usernameParam]
	isTaken, err := h.UserService.CheckIfUsernameIsTaken(r.Context(), username)
	if err != nil {
		res := Response{Message: "Unexpected error"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, res)
		return
	}
	res := response{
		Username: username,
		IsTaken:  isTaken,
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: res})
}
