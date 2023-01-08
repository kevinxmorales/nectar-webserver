package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/db"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (*user.User, error)
	GetUserById(ctx context.Context, firebaseId string) (*user.User, error)
	AddUser(ctx context.Context, u user.NewUserRequest) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, u user.UpdateUserRequest) (*user.User, error)
	CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var userRequest user.NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	validate := validator.New()
	if err := validate.Struct(userRequest); err != nil {
		res := responseEntity{Message: "Invalid request, please include all required fields"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	insertedUser, err := h.UserService.AddUser(r.Context(), userRequest)
	if err != nil {
		if errors.Is(err, db.DuplicateKeyError) {
			res := responseEntity{Message: fmt.Sprintf("This email is already registered: %s", userRequest.Email)}
			log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusBadRequest))
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}
			return
		}
		res := responseEntity{Message: fmt.Sprintf("This email is already registered: %s", userRequest.Email)}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	res := responseEntity{Content: insertedUser, Message: "account successfully created"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		res := responseEntity{Message: "Invalid id, please provide a valid id"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	usr, err := h.UserService.GetUser(r.Context(), id)
	if err != nil {
		res := responseEntity{Message: "Unexpected error, could not get user info"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	res := responseEntity{Content: usr, Message: "account successfully created"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	usr, err := h.UserService.GetUserById(r.Context(), id)
	if err != nil {
		res := responseEntity{Message: "Unexpected error, could not get user info"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	res := responseEntity{Content: usr, Message: "account successfully created"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) updateUserProfileImage(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id := vars["id"]
	return
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var usr user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		res := responseEntity{Message: "Unexpected error, could not get user info"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	usrUpdated, err := h.UserService.UpdateUser(r.Context(), id, usr)
	if err != nil {
		res := responseEntity{Message: "Unexpected error, could not get user info"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	res := responseEntity{Content: usrUpdated, Message: "user data successfully updated"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.UserService.DeleteUser(r.Context(), id); err != nil {
		res := responseEntity{Message: "Unexpected error, could not delete user info, please try again"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
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

func (h *Handler) checkIfUsernameIsTaken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Username string `json:"username"`
		IsTaken  bool   `json:"isTaken"`
	}
	usernameParam := "username"
	params, err := h.ParseUrlQueryParams(r.URL, usernameParam)
	if err != nil {
		res := responseEntity{Message: "Unexpected error"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d | error: %v", http.StatusInternalServerError, err))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	username := params[usernameParam]
	isTaken, err := h.UserService.CheckIfUsernameIsTaken(r.Context(), username)
	if err != nil {
		res := responseEntity{Message: "Unexpected error"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d | error: %v", http.StatusInternalServerError, err))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	res := response{
		Username: username,
		IsTaken:  isTaken,
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseEntity{Content: res}); err != nil {
		panic(err)
	}
	return
}
