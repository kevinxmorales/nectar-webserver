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

func NewBlobStoreSession() (*session.Session, error) {
	log.Info("Initializing S3 Connection")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Info("Error creating S3 session: ", err.Error())
		return nil, err
	}
	log.Info("Successfully created S3 session")

	return sess, nil
}

func UploadToBlobStore(fileList []string, ctx context.Context) ([]string, error) {
	sess, err := NewBlobStoreSession()
	if err != nil {
		return nil, err
	}
	//Set up the S3 bucket
	bucket := os.Getenv("S3_BUCKET")
	acl := os.Getenv("AWS_ACL")

	//Set up the concurrency
	//Need wait group and mutex
	var wg sync.WaitGroup

	// The resulting s3 urls to the files
	s3Urls := make([]string, len(fileList))

	// Iterate over the local files that need to be updated
	for i, pathOfFile := range fileList {
		// Kick off goroutine with thread-safe function to upload to s3
		wg.Add(1)
		i := i
		go func(pathOfFile string, store *session.Session, wg *sync.WaitGroup) {
			// Defer the
			defer func() {
				wg.Done()
			}()
			file, _ := os.Open(pathOfFile)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Error(err)
				}
			}(file)
			uploader := s3manager.NewUploader(sess)
			result, err := uploader.UploadWithContext(
				ctx,
				&s3manager.UploadInput{
					ACL:    aws.String(acl),
					Bucket: aws.String(bucket),
					Key:    aws.String(file.Name()),
					Body:   file,
				})
			if err != nil {
				log.Error(err)
				return
			}
			s3Urls[i] = result.Location
			log.Info(fmt.Sprintf("Upload result: %+v\n", result))
		}(pathOfFile, sess, &wg)
	}
	// Block until WaitGroup counter is zero, then return the s3 urls
	wg.Wait()
	return s3Urls, nil
}

func UploadToBlobStore2(fileList []string, ctx context.Context, sess *session.Session) ([]string, error) {
	sess, err := NewBlobStoreSession()
	if err != nil {
		return nil, err
	}
	//Set up the S3 bucket
	bucket := os.Getenv("S3_BUCKET")
	acl := os.Getenv("AWS_ACL")

	//Set up the concurrency
	//Need wait group and mutex
	var wg sync.WaitGroup

	// The resulting s3 urls to the files
	s3Urls := make([]string, len(fileList))

	// Iterate over the local files that need to be updated
	for i, pathOfFile := range fileList {
		// Kick off goroutine with thread-safe function to upload to s3
		wg.Add(1)
		i := i
		go func(pathOfFile string, store *session.Session, wg *sync.WaitGroup) {
			// Defer the
			defer func() {
				wg.Done()
			}()
			file, _ := os.Open(pathOfFile)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Error(err)
				}
			}(file)
			uploader := s3manager.NewUploader(sess)
			result, err := uploader.UploadWithContext(
				ctx,
				&s3manager.UploadInput{
					ACL:    aws.String(acl),
					Bucket: aws.String(bucket),
					Key:    aws.String(file.Name()),
					Body:   file,
				})
			if err != nil {
				log.Error(err)
				return
			}
			s3Urls[i] = result.Location
			log.Info(fmt.Sprintf("Upload result: %+v\n", result))
		}(pathOfFile, sess, &wg)
	}
	// Block until WaitGroup counter is zero, then return the s3 urls
	wg.Wait()
	return s3Urls, nil
}

func DeleteFromS3(fileName string) error {
	sess, err := NewBlobStoreSession()
	if err != nil {
		log.Info("Error creating S3 session", err.Error())
		return err
	}
	service := s3.New(sess)
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
