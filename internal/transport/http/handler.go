package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type Response struct {
	Content any    `json:"content"`
	Message string `json:"message"`
}

type Handler struct {
	Router *mux.Router

	AuthService   AuthService
	CareService   CareService
	HealthService HealthService
	PlantService  PlantService
	UserService   UserService
	Server        *http.Server
}

// NewHandler - returns a pointer to a http handler
// Need to give the handler all the different services
func NewHandler(
	plantService PlantService,
	userService UserService,
	careService CareService,
	authService AuthService,
	healthService HealthService) *Handler {

	//Create the http handler
	h := &Handler{
		PlantService:  plantService,
		UserService:   userService,
		CareService:   careService,
		AuthService:   authService,
		HealthService: healthService,
	}

	h.Router = mux.NewRouter()
	h.mapRoutes()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)
	port := os.Getenv("PORT")
	address := fmt.Sprintf("0.0.0.0:%s", port)
	h.Server = &http.Server{
		Addr:    address,
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", h.healthCheck).Methods(http.MethodGet)
	// Plant Endpoints
	h.Router.HandleFunc("/api/v1/plant", h.JWTAuth(h.AddPlant)).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/plant/image", h.JWTAuth(h.AddPlantImage)).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/plant/{id}", h.JWTAuth(h.GetPlant)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/plant/user/{id}", h.JWTAuth(h.GetPlantsByUserId)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/plant/{id}", h.JWTAuth(h.UpdatePlant)).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/plant/{id}", h.JWTAuth(h.DeletePlant)).Methods(http.MethodDelete)
	// User Endpoints
	h.Router.HandleFunc("/api/v1/user", h.CreateUser).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/user/{id}", h.JWTAuth(h.GetUser)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/user/id/{id}", h.JWTAuth(h.GetUserById)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/user/id/{id}", h.JWTAuth(h.UpdateUser)).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/user/id/{id}/image", h.JWTAuth(h.UpdateUserProfileImage)).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/user/id/{id}", h.JWTAuth(h.DeleteUser)).Methods(http.MethodDelete)
	h.Router.HandleFunc("/api/v1/user/username-check/is-taken", h.CheckIfUsernameIsTaken).Methods(http.MethodGet)

	//Plant Care Log Endpoints
	h.Router.HandleFunc("/api/v1/plant-care", h.JWTAuth(h.AddCareLogEntry)).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/plant-care/{id}", h.JWTAuth(h.GetCareLogsEntries)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/plant-care/{id}", h.JWTAuth(h.UpdateCareLogEntry)).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/plant-care/{id}", h.JWTAuth(h.DeleteCareLogEntry)).Methods(http.MethodDelete)
	h.Router.HandleFunc("/api/v1/plant-care/user/{id}", h.JWTAuth(h.GetAllUsersCareLogs)).Methods(http.MethodGet)

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
			if err != nil {
				return err
			}
			defer file.Close()
			year, month, day := time.Now().Date()
			hour := time.Now().Hour()
			minute := time.Now().Minute()
			newFileName := fmt.Sprintf("/tmp/%d-%d-%d-T-%d-%d-%s", year, month, day, hour, minute, files[0].Filename)
			fmt.Println(newFileName)
			out, err := os.Create(newFileName)
			if err != nil {
				return errors.New("unable to create the file for writing")
			}
			defer out.Close()
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

func (h *Handler) ParseImagesFromRequestBody(request *http.Request, numImages int) ([]image.Image, error) {
	// Create variables to hold image data
	var images []image.Image

	// Read the request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the request body into a map
	var bodyData map[string]interface{}
	if err := json.Unmarshal(body, &bodyData); err != nil {
		return nil, err
	}

	// Iterate over the three images in the request body
	for i := 0; i < numImages; i++ {
		// Get the image data from the body
		imageData := bodyData[fmt.Sprintf("image%d", i)].(string)

		// Decode the base64 image
		img, err := base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			return images, err
		}

		// Decode the image
		imgReader := bytes.NewReader(img)
		imgDecoded, _, err := image.Decode(imgReader)
		if err != nil {
			return images, err
		}
		// Append the image to the image array
		images = append(images, imgDecoded)
	}

	// Return the images
	return images, nil
}

func ParseImageFromRequestBody(r *http.Request) (string, error) {
	// parse the multipart form in the request
	if err := r.ParseMultipartForm(1024); err != nil {
		return "", err
	}

	// get the file from the form
	file, _, err := r.FormFile("image")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// read the file into a buffer
	var buf bytes.Buffer
	io.Copy(&buf, file)

	year, month, day := time.Now().Date()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	// write the buffer to a file
	fileName := fmt.Sprintf("/tmp/%d-%d-%d-T-%d-%d-%s.jpeg", year, month, day, hour, minute, uuid.NewV4().String())
	if err = ioutil.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
		return "", err
	}
	return fileName, nil
}

func (h *Handler) encodeJsonResponse(w *http.ResponseWriter, res Response) {
	if err := json.NewEncoder(*w).Encode(res); err != nil {
		panic(err)
	}
}
