package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"net/http"
	"os"
	"strconv"
)

type PlantService interface {
	PostPlant(ctx context.Context, p plant.Plant, fileNames []string) (*plant.Plant, error)
	GetPlant(ctx context.Context, id int) (*plant.Plant, error)
	GetPlantsByUserId(ctx context.Context, id int) ([]plant.Plant, error)
	UpdatePlant(ctx context.Context, id int, p plant.Plant) (*plant.Plant, error)
	DeletePlant(ctx context.Context, id int) error
}

type Response struct {
	Message string `json:"message"`
}

func (h *Handler) PostPlant(w http.ResponseWriter, r *http.Request) {
	plantName, categoryId, userId, numberOfImages := "plantName", "catId", "userId", "numImages"
	params, err := h.ParseUrlQueryParams(r.URL, plantName, categoryId, userId, numberOfImages)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	//Get the number of images being uploaded
	numImagesStr := params[numberOfImages]
	numImages := 0
	if _, err := fmt.Sscan(numImagesStr, &numImages); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}

	if err := r.ParseMultipartForm(200000); err != nil {
		fmt.Fprintln(w, err)
		return
	}

	formData := r.MultipartForm
	fileNames, err := h.ParseFilesFromMultiPartFormData(formData, numImages)
	// remove files that were created for upload
	defer func() {
		for _, file := range fileNames {
			if err := os.Remove(file); err != nil {
				log.Error(err)
			}
		}
	}()
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	userID, err := strconv.Atoi(params[userId])
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	p := plant.Plant{
		Name:       params[plantName],
		UserId:     userID,
		FileNames:  fileNames,
		CategoryID: params[categoryId],
	}
	newPlant, err := h.PlantService.PostPlant(r.Context(), p, fileNames)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, newPlant)
	return
}

func (h *Handler) GetPlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	p, err := h.PlantService.GetPlant(r.Context(), id)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, p)
	return
}

func (h *Handler) GetPlantsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	plantList, err := h.PlantService.GetPlantsByUserId(r.Context(), id)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	res := ResponseEntity{
		Content:    plantList,
		HttpStatus: http.StatusOK,
		Messages:   []string{},
	}
	h.SendOkResponse(w, r, res)
	return
}

func (h *Handler) UpdatePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	var p plant.Plant
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	pl, err := h.PlantService.UpdatePlant(r.Context(), id, p)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, pl)
	return
}

func (h *Handler) DeletePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	if err := h.PlantService.DeletePlant(r.Context(), id); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	res := Response{Message: "successfully deleted"}
	h.SendOkResponse(w, r, res)
	return
}
