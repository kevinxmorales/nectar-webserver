package plant

import (
	"context"
	"fmt"
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
	AddPlantImageWithId(ctx context.Context, plantId string, imageUri string) (string, error)
	DeletePlant(ctx context.Context, id string) error
	UpdatePlant(ctx context.Context, id string, p Plant) (*Plant, error)
	GetCareLogsEntries(ctx context.Context, plantId string) ([]care.LogEntry, error)
	DeletePlantImage(ctx context.Context, plantId string, uri string) error
}

type MessageQueue interface {
	PushToQueue(ctx context.Context, topic string, message []byte) error
}

type BlobStore interface {
	UploadToBlobStore(fileList []string, ctx context.Context) (resultUris []string, err error)
}

// Service - is the struct on which out logic will
// be built upon
type Service struct {
	Store        Store
	MessageQueue MessageQueue
	BlobStore    BlobStore
}

// NewService - returns a pointer to a new service
func NewService(store Store, blobService BlobStore, messageQueue MessageQueue) *Service {
	return &Service{
		Store:        store,
		BlobStore:    blobService,
		MessageQueue: messageQueue,
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

func (s *Service) UpdatePlant(ctx context.Context, id string, updatedPlant Plant, imagesToDelete []string) (*Plant, error) {
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

func (s *Service) AddPlantImage(ctx context.Context, uri string) (string, error) {
	tag := "plant.AddPlantImage"
	resultUris, err := s.BlobStore.UploadToBlobStore([]string{uri}, ctx)
	if err != nil {
		return "", fmt.Errorf("blob.UploadToBlobStore in %s failed for %v", tag, err)
	}
	resultUri := resultUris[0]
	return resultUri, nil
}

func (s *Service) AddPlantImageWithId(ctx context.Context, plantId string, uri string) (*Plant, string, error) {
	tag := "plant.AddImageToPlant"
	resultUris, err := s.BlobStore.UploadToBlobStore([]string{uri}, ctx)
	if err != nil {
		return nil, "", fmt.Errorf("blob.UploadToBlobStore in %s failed for %v", tag, err)
	}
	resultUri := resultUris[0]
	if _, err := s.Store.AddPlantImageWithId(ctx, plantId, resultUri); err != nil {
		return nil, "", fmt.Errorf("store.AddImageToPlant in %s failed for %v", tag, err)
	}
	p, err := s.Store.GetPlant(ctx, plantId)
	if err != nil {
		return nil, "", fmt.Errorf("store.GetPlant in %s failed for %v", tag, err)
	}
	return p, resultUri, nil
}

func (s *Service) DeletePlantImage(ctx context.Context, plantId string, uri string) error {
	log.Infof("Deleting image %s belonging to plant %s", uri, plantId)
	return s.Store.DeletePlantImage(ctx, plantId, uri)
}
