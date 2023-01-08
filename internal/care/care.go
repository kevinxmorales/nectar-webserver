package care

import (
	"context"
	"time"
)

type LogEntry struct {
	Id            string    `json:"id"`
	Date          time.Time `json:"date"`
	PlantId       string    `json:"plantId"`
	Notes         string    `json:"notes"`
	WasWatered    bool      `json:"wasWatered"`
	WasFertilized bool      `json:"wasFertilized"`
}

type Store interface {
	GetCareLogsEntries(ctx context.Context, plantId string) ([]LogEntry, error)
	AddCareLogEntry(ctx context.Context, entry LogEntry) (*LogEntry, error)
	DeleteCareLogEntry(ctx context.Context, logEntryId string) error
	UpdateCareLogEntry(ctx context.Context, logEntryId string, entry LogEntry) (*LogEntry, error)
}

// Service - is the struct on which our logic will
// be built upon
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetCareLogsEntries(ctx context.Context, plantId string) ([]LogEntry, error) {
	return s.Store.GetCareLogsEntries(ctx, plantId)
}

func (s *Service) AddCareLogEntry(ctx context.Context, entry LogEntry) (*LogEntry, error) {
	return s.Store.AddCareLogEntry(ctx, entry)
}

func (s *Service) DeleteCareLogEntry(ctx context.Context, logEntryId string) error {
	return s.Store.DeleteCareLogEntry(ctx, logEntryId)
}

func (s *Service) UpdateCareLogEntry(ctx context.Context, logEntryId string, entry LogEntry) (*LogEntry, error) {
	return s.Store.UpdateCareLogEntry(ctx, logEntryId, entry)
}
