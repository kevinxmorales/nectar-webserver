package blob

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type S3Response struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
}

type Service struct {
	BlobSession *session.Session
}

func NewService() (*Service, error) {
	log.Info("initializing S3 Connection")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Info("FAILED to create S3 session: ", err.Error())
		return nil, err
	}
	log.Info("successfully created S3 session")

	return &Service{BlobSession: sess}, nil
}

func (s *Service) UploadToBlobStore(fileList []string, ctx context.Context) ([]string, error) {
	//Set up the S3 bucket
	bucket := os.Getenv("S3_BUCKET")
	acl := os.Getenv("AWS_ACL")

	//Set up the concurrency
	//Need wait group and mutex
	var wg sync.WaitGroup

	// The resulting s3 urls to the files
	s3Urls := make([]string, len(fileList))
	s3Errors := make([]error, len(fileList))

	// Iterate over the local files that need to be updated
	for i, pathOfFile := range fileList {
		handleClose := func(file *os.File) {
			if err := file.Close(); err != nil {
				log.Error(err)
				s3Errors[i] = err
			}
		}
		logAndStoreError := func(err error) {
			log.Error(err)
			s3Errors[i] = err
		}

		// Kick off goroutine with thread-safe function to upload to s3
		wg.Add(1)
		i := i
		go func(pathOfFile string, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
			}()
			file, err := os.Open(pathOfFile)
			if err != nil {
				logAndStoreError(err)
				return
			}
			defer func() {
				handleClose(file)
			}()
			uploader := s3manager.NewUploader(s.BlobSession)
			result, err := uploader.UploadWithContext(
				ctx,
				&s3manager.UploadInput{
					ACL:    aws.String(acl),
					Bucket: aws.String(bucket),
					Key:    aws.String(file.Name()),
					Body:   file,
				})
			if err != nil {
				logAndStoreError(err)
				return
			}
			s3Urls[i] = result.Location
			log.Info(fmt.Sprintf("Upload result: %+v\n", result))
		}(pathOfFile, &wg)
	}
	// Block until WaitGroup counter is zero, then return the s3 urls
	wg.Wait()

	//Check if any errors occurred
	for _, err := range s3Errors {
		if err != nil {
			return nil, err
		}
	}

	return s3Urls, nil
}

func (s *Service) DeleteFromBlobStore(fileName string) error {
	service := s3.New(s.BlobSession)
	bucket := os.Getenv("S3_BUCKET")
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}
	result, err := service.DeleteObject(input)
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			switch awsError.Code() {
			default:
				log.Error(awsError.Error())
				return awsError
			}
		} else {
			log.Error(err.Error())
			return err
		}
	}
	log.Info("Successfully deleted item from S3", result.String())
	return nil
}
