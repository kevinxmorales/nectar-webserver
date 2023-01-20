package care

import (
	"context"
)

type LogEntry struct {
	Id            string `json:"id"`
	PlantId       string `json:"plantId"`
	Notes         string `json:"notes"`
	CareDate      string `json:"careDate"`
	CreatedAt     string `json:"createdAt"`
	PlantImage    string `json:"plantImage"`
	PlantName     string `json:"plantName"`
	WasWatered    bool   `json:"wasWatered"`
	WasFertilized bool   `json:"wasFertilized"`
}

type Store interface {
	GetAllUsersCareLogEntries(ctx context.Context, userId string) ([]LogEntry, error)
	GetCareLogsEntries(ctx context.Context, plantId string) ([]LogEntry, error)
	AddCareLogEntry(ctx context.Context, entry LogEntry) (*LogEntry, error)
	DeleteCareLogEntry(ctx context.Context, logEntryId string) error
	UpdateCareLogEntry(ctx context.Context, logEntryId string, entry LogEntry) (*LogEntry, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetAllUsersCareLogs(ctx context.Context, userId string) ([]LogEntry, error) {
	return s.Store.GetAllUsersCareLogEntries(ctx, userId)
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
