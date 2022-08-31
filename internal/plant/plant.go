package plant

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/blob"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/care"
	"time"
)

// Plant - a representation of a plant
type Plant struct {
	Id             int             `json:"plantId"`
	IdStr          string          `json:"id"`
	UserId         int             `json:"userId"`
	Username       string          `json:"username"`
	Name           string          `json:"name"`
	CategoryID     string          `json:"categoryId"`
	Images         []ImageUrls     `json:"images"`
	CreatedAt      time.Time       `json:"createdAt"`
	CareLogEntries []care.LogEntry `json:"careLogEntries"`
	FileNames      []string        `json:"-"`
}

type ImageUrls struct {
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

//Store - this interface defines all the methods
// the service needs in order to operate
type Store interface {
	GetPlant(ctx context.Context, id int) (*Plant, error)
	GetPlantsByUserId(ctx context.Context, userId int) ([]Plant, error)
	AddPlant(ctx context.Context, p Plant) (*Plant, error)
	DeletePlant(ctx context.Context, id int) error
	UpdatePlant(ctx context.Context, id int, p Plant) (*Plant, error)
	GetCareLogsEntries(ctx context.Context, plantId int) ([]care.LogEntry, error)
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

func (s *Service) GetPlant(ctx context.Context, id int) (*Plant, error) {
	log.Info("Retrieving a plant with id: ", id)
	return s.Store.GetPlant(ctx, id)
}

func (s *Service) GetPlantsByUserId(ctx context.Context, id int) ([]Plant, error) {
	tag := "plant.GetPlantsByUserId"
	plantList, err := s.Store.GetPlantsByUserId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Store.GetPlantsByUser in %s failed for %v", tag, err)
	}
	idList := make([]int, len(plantList))
	for i, plant := range plantList {
		idList[i] = plant.Id
		entries, err := s.Store.GetCareLogsEntries(ctx, plant.Id)
		if err != nil {
			return nil, fmt.Errorf("Store.GetCareLogsEntries in %s failed for %v", tag, err)
		}
		plantList[i].CareLogEntries = entries
	}
	return plantList, nil
}

func (s *Service) UpdatePlant(ctx context.Context, id int, updatedPlant Plant) (*Plant, error) {
	return s.Store.UpdatePlant(ctx, id, updatedPlant)
}

func (s *Service) DeletePlant(ctx context.Context, id int) error {
	return s.Store.DeletePlant(ctx, id)
}

func (s *Service) PostPlant(ctx context.Context, newPlant Plant, images []string) (*Plant, error) {
	log.Info("attempting to add a new plant")
	//TODO create a lower resolution image for each image in the images array for the thumbnail url
	/*
		thumbnailUrls := make([]string, len(images))
		//Something like this
		for i, v := range images {
			// 1. resize image at lower resolution
			// 2. add to thumbnailUrls[i]
		}
	*/

	//Upload all plant images to blob store
	s3Urls, err := blob.UploadToBlobStore(images, ctx, s.BlobStore)
	if err != nil {
		return nil, fmt.Errorf("blob.UploadToBlobStore in plant.PostPlant failed for %v", err)
	}
	newPlant.FileNames = s3Urls
	return s.Store.AddPlant(ctx, newPlant)
}
