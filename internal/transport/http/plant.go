package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
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
	AddPlantImage(ctx context.Context, filePath string) (string, error)
	AddPlantImageWithId(ctx context.Context, plantId string, filePath string) (*plant.Plant, string, error)
	GetPlant(ctx context.Context, id string) (*plant.Plant, error)
	GetPlantsByUserId(ctx context.Context, id string) ([]plant.Plant, error)
	UpdatePlant(ctx context.Context, id string, p plant.Plant, imagesToDelete []string) (*plant.Plant, error)
	DeletePlant(ctx context.Context, id string) error
	DeletePlantImage(ctx context.Context, plantId string, uri string) error
}

type response struct {
	Message string `json:"message"`
}

type NewPlantRequest struct {
	CommonName     string   `json:"commonName" validate:"required"`
	UserId         string   `json:"userId" validate:"required"`
	Images         []string `json:"images"`
	ScientificName string   `json:"scientificName"`
	Toxicity       string   `json:"toxicity"`
}

type UpdatePlantRequest struct {
	CommonName     string   `json:"commonName" validate:"required"`
	UserId         string   `json:"userId" validate:"required"`
	Images         []string `json:"images"`
	ImagesToDelete []string `json:"imagesToDelete"`
	ScientificName string   `json:"scientificName"`
	Toxicity       string   `json:"toxicity"`
}

func (h *Handler) AddPlant(w http.ResponseWriter, r *http.Request) {
	type addPlantResponse struct {
		ID string `json:"id"`
	}
	var pr NewPlantRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	validate := validator.New()
	if err := validate.Struct(pr); err != nil {
		res := Response{Message: "Invalid request, could not create plant"}
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
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
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		switch e := err.(type) {
		case *errs.NoEntityError:
			w.WriteHeader(http.StatusBadRequest)
			h.encodeJsonResponse(&w, Response{Message: e.Message})
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
			return
		}
	}
	res := Response{
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
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "Invalid id, please provide a valid plant id"})
		return
	}
	p, err := h.PlantService.GetPlant(r.Context(), id)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	res := Response{
		Content: getPlantResponse{Plant: *p},
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, res)
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
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	res := Response{
		Content: response{Plants: plantList},
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, res)
	return
}

func (h *Handler) UpdatePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var up UpdatePlantRequest
	if err := json.NewDecoder(r.Body).Decode(&up); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	imagesToDelete := up.ImagesToDelete
	updatedPlant := plant.Plant{
		PlantId:        id,
		CommonName:     up.CommonName,
		ScientificName: up.ScientificName,
		Toxicity:       up.Toxicity,
		UserId:         up.UserId,
		Images:         up.Images,
	}
	_, err := h.PlantService.UpdatePlant(r.Context(), id, updatedPlant, imagesToDelete)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Message: "Plant successfully updated"})
	return
}

func (h *Handler) DeletePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "No id provided, please provide an id for a plant"})
		return
	}
	if err := h.PlantService.DeletePlant(r.Context(), id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Message: "successfully deleted"})
	return
}

func (h *Handler) AddImageToPlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		res := Response{Message: "No id provided, please provide an id for a plant"}
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}
	filePath, err := ParseImageFromRequestBody(r)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "Unexpected error, could not add plant image"})
		return
	}
	updatedPlant, fileUri, err := h.PlantService.AddPlantImageWithId(r.Context(), id, filePath)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "Unexpected error, could not add plant image"})
		return
	}
	content := struct {
		Plant plant.Plant `json:"plant"`
		Uri   string      `json:"imageUrl"`
	}{Plant: *updatedPlant, Uri: fileUri}
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: content})
	return
}

func (h *Handler) AddPlantImage(w http.ResponseWriter, r *http.Request) {
	filePath, err := ParseImageFromRequestBody(r)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "Unexpected error, could not add plant image"})
		return
	}
	fileUri, err := h.PlantService.AddPlantImage(r.Context(), filePath)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "Unexpected error, could not add plant image"})
		return
	}
	content := struct {
		Uri string `json:"imageUrl"`
	}{Uri: fileUri}
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: content})
	return
}

func (h *Handler) DeletePlantImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", "invalid id", http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "No id provided, please provide an id for a plant"})
		return
	}
	var deleteImgReq struct {
		Uri string `json:"uri"`
	}
	spew.Dump(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&deleteImgReq); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	if err := h.PlantService.DeletePlantImage(r.Context(), id, deleteImgReq.Uri); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: "successfully deleted"})
}
