package plant

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/blob"
	"time"
)

// Plant - a representation of a plant
type Plant struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	UserID     string      `json:"userId"`
	CategoryID string      `json:"categoryId"`
	Images     []ImageUrls `json:"images"`
	CreatedAt  time.Time   `json:"createdAt"`
	FileNames  []string    `json:"-"`
}

type ImageUrls struct {
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

//Store - this interface defines all the methods
// the service needs in order to operate
type Store interface {
	GetPlant(context.Context, string) (*Plant, error)
	GetPlantsByUserId(context.Context, string) ([]Plant, error)
	AddPlant(context.Context, Plant) (*Plant, error)
	DeletePlant(context.Context, string) error
	UpdatePlant(context.Context, string, Plant) (*Plant, error)
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
	return s.Store.GetPlantsByUserId(ctx, id)
}

func (s *Service) UpdatePlant(ctx context.Context, ID string, updatedPlant Plant) (*Plant, error) {
	return s.Store.UpdatePlant(ctx, ID, updatedPlant)
}

func (s *Service) DeletePlant(ctx context.Context, id string) error {
	return s.Store.DeletePlant(ctx, id)
}

func (s *Service) PostPlant(ctx context.Context, newPlant Plant, images []string) (*Plant, error) {
	log.Info("attempting to add a new plant")
	//Upload all plant images to blob store
	s3Urls, err := blob.UploadToBlobStore2(images, ctx, s.BlobStore)
	if err != nil {
		return nil, err
	}
	newPlant.FileNames = s3Urls
	return s.Store.AddPlant(ctx, newPlant)
}
