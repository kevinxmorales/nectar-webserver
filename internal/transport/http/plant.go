package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	errs "gitlab.com/kevinmorales/nectar-rest-api/internal/nectar_errors"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"net/http"
)

type PlantService interface {
	AddPlant(ctx context.Context, p plant.Plant, images []string) (*plant.Plant, error)
	GetPlant(ctx context.Context, id string) (*plant.Plant, error)
	GetPlantsByUserId(ctx context.Context, id string) ([]plant.Plant, error)
	UpdatePlant(ctx context.Context, id string, p plant.Plant) (*plant.Plant, error)
	DeletePlant(ctx context.Context, id string) error
}

type response struct {
	Message string `json:"message"`
}

type plantRequest struct {
	CommonName     string   `json:"commonName" validate:"required"`
	UserId         string   `json:"userId" validate:"required"`
	Images         []string `json:"images"`
	ScientificName string   `json:"scientificName"`
	Toxicity       string   `json:"toxicity"`
}

func (h *Handler) addPlant(w http.ResponseWriter, r *http.Request) {
	type addPlantResponse struct {
		ID string `json:"id"`
	}
	var pr plantRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, responseEntity{Message: "An unexpected error occurred"})
		return
	}
	validate := validator.New()
	if err := validate.Struct(pr); err != nil {
		res := responseEntity{Message: "Invalid request, could not create plant"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, res)
		return
	}
	p := plant.Plant{
		CommonName:     pr.CommonName,
		ScientificName: pr.ScientificName,
		Images:         pr.Images,
		UserId:         pr.UserId,
		Toxicity:       pr.Toxicity,
	}
	newPlant, err := h.PlantService.AddPlant(r.Context(), p, p.Images)
	if err != nil {
		switch e := err.(type) {
		case *errs.NoEntityError:
			log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusBadRequest))
			w.WriteHeader(http.StatusBadRequest)
			h.encodeJsonResponse(&w, responseEntity{Message: e.Message})
			return
		default:
			log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
			w.WriteHeader(http.StatusInternalServerError)
			h.encodeJsonResponse(&w, responseEntity{Message: "An unexpected error occurred"})
			return
		}
	}
	res := responseEntity{
		Content: addPlantResponse{ID: newPlant.PlantId},
		Message: "plant successfully created",
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, res)
	return
}

func (h *Handler) GetPlant(w http.ResponseWriter, r *http.Request) {
	type getPlantResponse struct {
		Plant plant.Plant `json:"plant"`
	}
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := uuid.Parse(id)
	if err != nil {
		res := responseEntity{Message: "Invalid id, please provide a valid plant id"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	p, err := h.PlantService.GetPlant(r.Context(), id)
	if err != nil {
		log.Error(err)
		res := responseEntity{Message: "An unexpected error occurred"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	res := responseEntity{
		Content: getPlantResponse{Plant: *p},
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) GetPlantsByUserId(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Plants []plant.Plant `json:"plants"`
	}
	vars := mux.Vars(r)
	id := vars["id"]
	plantList, err := h.PlantService.GetPlantsByUserId(r.Context(), id)
	if err != nil {
		res := responseEntity{Message: "An unexpected error occurred"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	res := responseEntity{
		Content: response{Plants: plantList},
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) UpdatePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var p plant.Plant
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		res := responseEntity{Message: "An unexpected error occurred"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	_, err := h.PlantService.UpdatePlant(r.Context(), id, p)
	if err != nil {
		res := responseEntity{Message: "An unexpected error occurred"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseEntity{Message: "Plant successfully updated"}); err != nil {
		panic(err)
	}
	return
}

func (h *Handler) DeletePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		res := responseEntity{Message: "No id provided, please provide an id for a plant"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	if err := h.PlantService.DeletePlant(r.Context(), id); err != nil {
		res := responseEntity{Message: "An unexpected error occurred"}
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseEntity{Message: "successfully deleted"}); err != nil {
		panic(err)
	}
	return
}
