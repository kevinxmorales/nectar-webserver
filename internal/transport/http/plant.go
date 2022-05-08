package http

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"net/http"
)

type PlantService interface {
	PostPlant(ctx context.Context, newPlant plant.Plant) (plant.Plant, error)
	GetPlant(ctx context.Context, ID string) (plant.Plant, error)
	UpdatePlant(ctx context.Context, ID string, newPlant plant.Plant) (plant.Plant, error)
	DeletePlant(ctx context.Context, ID string) error
}

type Response struct {
	Message string `json:"message"`
}

type PostPlantRequest struct {
	Name   string `json:"name" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

func convertPlantRequestToPlant(p PostPlantRequest) plant.Plant {
	return plant.Plant{
		Name:   p.Name,
		UserId: p.UserID,
	}
}

func (h *Handler) PostPlant(w http.ResponseWriter, r *http.Request) {
	var plantRequest PostPlantRequest
	if err := json.NewDecoder(r.Body).Decode(&plantRequest); err != nil {
		return
	}
	validate := validator.New()
	err := validate.Struct(plantRequest)
	if err != nil {
		http.Error(w, "not a valid plant object", http.StatusBadRequest)
		return
	}
	convertedPlant := convertPlantRequestToPlant(plantRequest)
	insertedPlant, err := h.PlantService.PostPlant(r.Context(), convertedPlant)
	if err != nil {
		log.Print(err)
		return
	}
	if err := json.NewEncoder(w).Encode(insertedPlant); err != nil {
		panic(err)
		return
	}

}

func (h *Handler) GetPlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := h.PlantService.GetPlant(r.Context(), id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdatePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var p plant.Plant
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := h.PlantService.UpdatePlant(r.Context(), id, p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
}

func (h *Handler) DeletePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.PlantService.DeletePlant(r.Context(), id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(Response{Message: "successfully deleted"})
	if err != nil {
		panic(err)
	}
}
