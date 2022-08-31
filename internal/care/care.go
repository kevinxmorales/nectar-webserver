package care

import (
	"context"
	"time"
)

type LogEntry struct {
	Id            int       `json:"id"`
	Date          time.Time `json:"date"`
	PlantId       int       `json:"plantId"`
	Notes         string    `json:"notes"`
	WasWatered    bool      `json:"wasWatered"`
	WasFertilized bool      `json:"wasFertilized"`
}

type Store interface {
	GetCareLogsEntries(ctx context.Context, plantId int) ([]LogEntry, error)
	AddCareLogEntry(ctx context.Context, entry LogEntry) (*LogEntry, error)
	DeleteCareLogEntry(ctx context.Context, logEntryId int) error
	UpdateCareLogEntry(ctx context.Context, logEntryId int, entry LogEntry) (*LogEntry, error)
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

func (s *Service) GetCareLogsEntries(ctx context.Context, plantId int) ([]LogEntry, error) {
	return s.Store.GetCareLogsEntries(ctx, plantId)
}

func (s *Service) AddCareLogEntry(ctx context.Context, entry LogEntry) (*LogEntry, error) {
	return s.Store.AddCareLogEntry(ctx, entry)
}

func (s *Service) DeleteCareLogEntry(ctx context.Context, logEntryId int) error {
	return s.Store.DeleteCareLogEntry(ctx, logEntryId)
}

func (s *Service) UpdateCareLogEntry(ctx context.Context, logEntryId int, entry LogEntry) (*LogEntry, error) {
	return s.Store.UpdateCareLogEntry(ctx, logEntryId, entry)
}
