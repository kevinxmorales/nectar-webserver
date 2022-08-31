package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"net/http"
	"strconv"
)

type CareService interface {
	GetCareLogsEntries(ctx context.Context, plantId int) ([]care.LogEntry, error)
	AddCareLogEntry(ctx context.Context, entry care.LogEntry) (*care.LogEntry, error)
	DeleteCareLogEntry(ctx context.Context, logEntryId int) error
	UpdateCareLogEntry(ctx context.Context, logEntryId int, entry care.LogEntry) (*care.LogEntry, error)
}

type CareLogEntryRequest struct {
	PlantId       int    `json:"plantId" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
	WasWatered    bool   `json:"wasWatered" validate:"required"`
	WasFertilized bool   `json:"wasFertilized" validate:"required"`
}

func convertRequestToLogEntry(request CareLogEntryRequest) care.LogEntry {
	return care.LogEntry{
		PlantId:       request.PlantId,
		Notes:         request.Notes,
		WasWatered:    request.WasWatered,
		WasFertilized: request.WasFertilized,
	}
}

func (h *Handler) GetCareLogsEntries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	entries, err := h.CareService.GetCareLogsEntries(r.Context(), id)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, entries)
	return
}

func (h *Handler) AddCareLogEntry(w http.ResponseWriter, r *http.Request) {
	var logEntryRequest CareLogEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&logEntryRequest); err != nil {
		h.SendBadRequestResponse(w, r, errors.New("unable to decode request"))
		return
	}
	validate := validator.New()
	if err := validate.Struct(logEntryRequest); err != nil {
		h.SendBadRequestResponse(w, r, errors.New("not a valid log entry object"))
		return
	}
	logEntry := convertRequestToLogEntry(logEntryRequest)
	insertedEntry, err := h.CareService.AddCareLogEntry(r.Context(), logEntry)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, insertedEntry)
	return
}

func (h *Handler) UpdateCareLogEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	var logEntryRequest CareLogEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&logEntryRequest); err != nil {
		h.SendBadRequestResponse(w, r, errors.New("unable to decode request"))
		return
	}
	validate := validator.New()
	if err := validate.Struct(logEntryRequest); err != nil {
		h.SendBadRequestResponse(w, r, errors.New("not a valid log entry object"))
		return
	}
	entry := convertRequestToLogEntry(logEntryRequest)
	updatedEntry, err := h.CareService.UpdateCareLogEntry(r.Context(), id, entry)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	h.SendOkResponse(w, r, updatedEntry)
	return
}

func (h *Handler) DeleteCareLogEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendBadRequestResponse(w, r, err)
		return
	}
	err = h.CareService.DeleteCareLogEntry(r.Context(), id)
	if err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	message := "entry successfully deleted"
	h.SendOkResponse(w, r, message)
	return
}
