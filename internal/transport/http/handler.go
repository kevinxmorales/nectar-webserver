package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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
	h.Router.HandleFunc("/api/v1/plant/{id}", JWTAuth(h.GetPlant)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/plant/user/{id}", JWTAuth(h.GetPlantsByUserId)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/plant/{id}", JWTAuth(h.UpdatePlant)).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/plant/{id}", JWTAuth(h.DeletePlant)).Methods(http.MethodDelete)
	// User Endpoints
	h.Router.HandleFunc("/api/v1/user", h.PostUser).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/user/{id}", JWTAuth(h.GetUser)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/user/email/{email}", JWTAuth(h.GetUserByEmail)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/user/{id}", JWTAuth(h.UpdateUser)).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/user/{id}", JWTAuth(h.DeleteUser)).Methods(http.MethodDelete)

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

// ParseUrlQueryParams a function to parse url query params. The function accepts a URL and a slice of map keys that
//are expected in the query parameters. This function returns map that is safe to use with all expected keys in it/**
func (h *Handler) ParseUrlQueryParams(url *url.URL, paramMapKeys ...string) (map[string]string, error) {
	params := url.Query()
	mapValues := make(map[string]string)
	for _, key := range paramMapKeys {
		val, ok := params[key]
		//if key is not present OR the key is present and the value is an empty array
		if !ok || len(val) < 1 {
			return nil, errors.New(fmt.Sprintf("key, %s, not found in url query parameters", key))
		}
		mapValues[key] = val[0]
	}
	return mapValues, nil
}

func (h *Handler) ParseFilesFromMultiPartFormData(formData *multipart.Form, numberOfFiles int) ([]string, error) {
	log.Infof("number of files: %d", numberOfFiles)
	var fileNames []string
	for index := 0; index < numberOfFiles; index++ {
		fileName := fmt.Sprintf("image%d", index)
		files, ok := formData.File[fileName] // grab the filenames
		// loop through the files one by one
		err := func() error {
			if !ok {
				return errors.New(fmt.Sprintf("file not found: %s", fileName))
			}
			//Check if slice is empty
			if len(files) < 1 {
				return errors.New(fmt.Sprintf("unable to parse file supplied: %s", fileName))
			}
			file, err := files[0].Open()
			defer file.Close()
			if err != nil {
				return err
			}
			year, month, day := time.Now().Date()
			hour := time.Now().Hour()
			minute := time.Now().Minute()
			newFileName := fmt.Sprintf("/tmp/%d-%d-%d-T-%d-%d-%s", year, month, day, hour, minute, files[0].Filename)
			fmt.Println(newFileName)
			out, err := os.Create(newFileName)
			defer out.Close()
			if err != nil {
				return errors.New("unable to create the file for writing")
			}
			// file not files[i] !
			if _, err := io.Copy(out, file); err != nil {
				return err
			}
			fileNames = append(fileNames, newFileName)
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}
	return fileNames, nil
}

func (h *Handler) SendOkResponse(w http.ResponseWriter, r *http.Request, data any) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func (h *Handler) SendBadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Error(err)
	h.SendErrorResponse(w, r, http.StatusBadRequest)
}

func (h *Handler) SendServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Error(err)
	h.SendErrorResponse(w, r, http.StatusInternalServerError)
}

func (h *Handler) SendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info(fmt.Sprintf("unsuccessful request, status code: %d", statusCode))
	w.WriteHeader(statusCode)
}
