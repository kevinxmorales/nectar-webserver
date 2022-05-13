package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const EMPTY = ""

type Handler struct {
	Router       *mux.Router
	UserService  UserService
	PlantService PlantService
	AuthService  AuthService
	Server       *http.Server
}

func NewHandler(plantService PlantService, userService UserService, authService AuthService) *Handler {
	h := &Handler{
		PlantService: plantService,
		UserService:  userService,
		AuthService:  authService,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)
	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am alive")
	})
	// Auth Endpoints
	h.Router.HandleFunc("/api/v1/auth", h.Login).Methods(http.MethodPost)
	// Plant Endpoints
	h.Router.HandleFunc("/api/v1/plant", JWTAuth(h.PostPlant)).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/plant/{id}", h.GetPlant).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/plant/{id}", JWTAuth(h.UpdatePlant)).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/plant/{id}", JWTAuth(h.DeletePlant)).Methods(http.MethodDelete)
	// User Endpoints
	h.Router.HandleFunc("/api/v1/user", h.PostUser).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/user/{id}", h.GetUser).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/user/email/{email}", h.GetUserByEmail).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/user/{id}", h.UpdateUser).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/user/{id}", h.DeleteUser).Methods(http.MethodDelete)

}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Info("shut down gracefully")
	return nil
}
