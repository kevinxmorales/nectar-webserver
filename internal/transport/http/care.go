package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"net/http"
	"time"
)

type CareService interface {
	GetAllUsersCareLogs(ctx context.Context, userId string) ([]care.LogEntry, error)
	GetCareLogsEntries(ctx context.Context, plantId string) ([]care.LogEntry, error)
	AddCareLogEntry(ctx context.Context, entry care.LogEntry) (*care.LogEntry, error)
	DeleteCareLogEntry(ctx context.Context, logEntryId string) error
	UpdateCareLogEntry(ctx context.Context, logEntryId string, entry care.LogEntry) (*care.LogEntry, error)
}

type CareLogEntryRequest struct {
	PlantId       string `json:"plantId"`
	Notes         string `json:"notes"`
	CareDate      string `json:"careDate"`
	WasWatered    bool   `json:"wasWatered"`
	WasFertilized bool   `json:"wasFertilized"`
}

func convertRequestToLogEntry(request CareLogEntryRequest) care.LogEntry {
	return care.LogEntry{
		PlantId:       request.PlantId,
		Notes:         request.Notes,
		WasWatered:    request.WasWatered,
		WasFertilized: request.WasFertilized,
	}
}

func (h *Handler) GetAllUsersCareLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "Invalid request, please include care log id"})
		return
	}
	entries, err := h.CareService.GetAllUsersCareLogs(r.Context(), id)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: entries})
	return
}

func (h *Handler) GetCareLogsEntries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "Invalid request, please include care log id"})
		return
	}
	entries, err := h.CareService.GetCareLogsEntries(r.Context(), id)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: entries})
	return
}

func (h *Handler) AddCareLogEntry(w http.ResponseWriter, r *http.Request) {
	var logEntryRequest CareLogEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&logEntryRequest); err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("before %s", logEntryRequest.CareDate))
	date, err := time.Parse("12-31-2022", logEntryRequest.CareDate)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("after %s", date))
	logEntry := convertRequestToLogEntry(logEntryRequest)
	insertedEntry, err := h.CareService.AddCareLogEntry(r.Context(), logEntry)
	if err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: insertedEntry})
	return
}

func (h *Handler) UpdateCareLogEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "Invalid request, please include care log id"})
		return
	}
	var logEntryRequest CareLogEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&logEntryRequest); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	validate := validator.New()
	if err := validate.Struct(logEntryRequest); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	entry := convertRequestToLogEntry(logEntryRequest)
	updatedEntry, err := h.CareService.UpdateCareLogEntry(r.Context(), id, entry)
	if err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: updatedEntry})
	return
}

func (h *Handler) DeleteCareLogEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		h.encodeJsonResponse(&w, Response{Message: "Invalid request, please include care log id"})
		return
	}
	if err := h.CareService.DeleteCareLogEntry(r.Context(), id); err != nil {
		log.Errorf(fmt.Sprintf("unsuccessful request, reason: %s,status code: %d", err.Error(), http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, Response{Message: "An unexpected error occurred"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, Response{Content: "entry successfully deleted"})
	return
}
