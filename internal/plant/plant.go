package plant

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"time"
)

type Plant struct {
	PlantId          string    `json:"plantId"`
	UserId           string    `json:"userId"`
	Username         string    `json:"username"`
	UserProfileImage string    `json:"userProfileImage"`
	CommonName       string    `json:"commonName"`
	ScientificName   string    `json:"scientificName"`
	Toxicity         string    `json:"toxicity"`
	CreatedAt        time.Time `json:"createdAt"`
	Images           []string  `json:"images"`
	SearchTerms      []string  `json:"searchTerms"`
}

type ImageUrls struct {
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

//Store - this interface defines all the methods
// the service needs in order to operate
type Store interface {
	GetPlant(ctx context.Context, id string) (*Plant, error)
	GetPlantsByUserId(ctx context.Context, userId string) ([]Plant, error)
	AddPlant(ctx context.Context, p Plant, images []string) (*Plant, error)
	DeletePlant(ctx context.Context, id string) error
	UpdatePlant(ctx context.Context, id string, p Plant) (*Plant, error)
	GetCareLogsEntries(ctx context.Context, plantId string) ([]care.LogEntry, error)
}

// Service - is the struct on which out logic will
// be built upon
type Service struct {
	Store     Store
	BlobStore *session.Session
}

// NewService - returns a pointer to a new service
func NewService(store Store, blobStoreSession *session.Session) *Service {
	return &Service{
		Store:     store,
		BlobStore: blobStoreSession,
	}
}

func (s *Service) GetPlant(ctx context.Context, id string) (*Plant, error) {
	log.Info("Retrieving a plant with id: ", id)
	return s.Store.GetPlant(ctx, id)
}

func (s *Service) GetPlantsByUserId(ctx context.Context, id string) ([]Plant, error) {
	tag := "plant.GetPlantsByUserId"
	pl, err := s.Store.GetPlantsByUserId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Store.GetPlantsByUser in %s failed for %v", tag, err)
	}
	return pl, nil
}

func (s *Service) UpdatePlant(ctx context.Context, id string, updatedPlant Plant) (*Plant, error) {
	return s.Store.UpdatePlant(ctx, id, updatedPlant)
}

func (s *Service) DeletePlant(ctx context.Context, id string) error {
	log.Info("attempting to delete a plant")
	return s.Store.DeletePlant(ctx, id)
}

func (s *Service) AddPlant(ctx context.Context, newPlant Plant, images []string) (*Plant, error) {
	log.Info("attempting to add a new plant")
	return s.Store.AddPlant(ctx, newPlant, images)
}
