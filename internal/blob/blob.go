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
	Bucket      string
	Acl         string
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
		log.Errorf("FAILED to create S3 session: %s", err.Error())
		return nil, err
	}
	log.Info("successfully created S3 session")

	s := &Service{
		BlobSession: sess,
		Bucket:      os.Getenv("S3_BUCKET"),
		Acl:         os.Getenv("AWS_ACL"),
	}
	return s, nil
}

func (s *Service) UploadToBlobStore(fileList []string, ctx context.Context) ([]string, error) {

	//Set up the concurrency
	var wg sync.WaitGroup

	// The resulting s3 urls to the files
	s3Urls := make([]string, len(fileList))
	s3Errors := make([]error, len(fileList))

	handleClose := func(file *os.File, index int) {
		if err := file.Close(); err != nil {
			log.Error(err)
			s3Errors[index] = err
		}
	}
	uploadImage := func(pathOfFile string, index int, wg *sync.WaitGroup) {
		defer wg.Done()

		file, err := os.Open(pathOfFile)
		if err != nil {
			s3Errors[index] = err
			return
		}
		defer handleClose(file, index)

		uploader := s3manager.NewUploader(s.BlobSession)
		input := &s3manager.UploadInput{
			ACL:    aws.String(s.Acl),
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(file.Name()),
			Body:   file,
		}
		result, err := uploader.UploadWithContext(ctx, input)
		if err != nil {
			s3Errors[index] = err
			return
		}
		s3Urls[index] = result.Location
		log.Info(fmt.Sprintf("Upload result: %+v\n", result))
	}

	// Iterate over the local files that need to be updated
	for i, pathOfFile := range fileList {
		// Kick off goroutine with thread-safe function to upload to s3
		wg.Add(1)
		go uploadImage(pathOfFile, i, &wg)
	}
	// Block until WaitGroup counter is zero, then return the s3 urls
	wg.Wait()

	//Check if any errors occurred
	for _, err := range s3Errors {
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return s3Urls, nil
}

func (s *Service) DeleteFromBlobStore(fileName string) error {
	service := s3.New(s.BlobSession)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fileName),
	}
	result, err := service.DeleteObject(input)
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			log.Errorf("error occurred on AWS side of transaction: %v", awsError)
		}
		return err
	}
	log.Info("Successfully deleted item from S3", result.String())
	return nil
}
