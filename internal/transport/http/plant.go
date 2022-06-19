package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/blob"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/plant"
	"io"
	"net/http"
	"os"
	"time"
)

type PlantService interface {
	PostPlant(context.Context, plant.Plant) (*plant.Plant, error)
	GetPlant(context.Context, string) (*plant.Plant, error)
	GetPlantsByUserId(context.Context, string) ([]plant.Plant, error)
	UpdatePlant(context.Context, string, plant.Plant) (*plant.Plant, error)
	DeletePlant(context.Context, string) error
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
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	formdata := r.MultipartForm

	var fileNames []string
	for index := 0; index < 3; index++ {
		fileName := fmt.Sprintf("image%d", index)
		files, ok := formdata.File[fileName] // grab the filenames
		// loop through the files one by one
		func() {
			if !ok {
				log.Error(fmt.Sprintf("file not found: %s", fileName))
				return
			}
			file, err := files[0].Open()
			defer file.Close()
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			year, month, day := time.Now().Date()
			hour := time.Now().Hour()
			minute := time.Now().Minute()
			newFileName := fmt.Sprintf("/tmp/%d-%d-%d-T-%d-%d-%s", year, month, day, hour, minute, files[0].Filename)
			fmt.Println(newFileName)
			out, err := os.Create(newFileName)
			defer out.Close()
			if err != nil {
				fmt.Fprintf(w, "Unable to create the file for writing")
				return
			}
			/*
				image, err := jpeg.Decode(file)
				opt := jpeg.Options{
					Quality: 90,
				}
				err = jpeg.Encode(out, image, &opt)
				if err != nil {
					fmt.Fprintln(w, err)
					return
				}

			*/

			_, err = io.Copy(out, file) // file not files[i] !
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

			//fmt.Fprintf(w, "Files uploaded successfully : ")
			//fmt.Fprintf(w, files[0].Filename+"\n")
			fileNames = append(fileNames, newFileName)
		}()
	}

	fileUrls, err := blob.UploadToBlobStore(fileNames, r.Context())
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, fileUrls)
}

func (h *Handler) GetPlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	p, err := h.PlantService.GetPlant(r.Context(), id)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, p)
}

func (h *Handler) GetPlantsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	log.Info(fmt.Sprintf("Attempting to get all plants that belong to user with id: %s", id))
	p, err := h.PlantService.GetPlantsByUserId(r.Context(), id)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, p)
}

func (h *Handler) UpdatePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	var p plant.Plant
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		send500Response(w, r, err)
		return
	}
	pl, err := h.PlantService.UpdatePlant(r.Context(), id, p)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	sendOkResponse(w, r, pl)
}

func (h *Handler) DeletePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == EMPTY {
		sendBadRequestResponse(w, r, errors.New("no id was supplied with this request"))
		return
	}
	err := h.PlantService.DeletePlant(r.Context(), id)
	if err != nil {
		send500Response(w, r, err)
		return
	}
	res := Response{Message: "successfully deleted"}
	sendOkResponse(w, r, res)
}
