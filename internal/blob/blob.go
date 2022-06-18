package blob

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"os"
)

type S3Response struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
}

func NewS3Session() (*session.Session, error) {
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
		return sess, err
	}
	log.Info("Successfully created S3 session")
	return sess, nil
}

func UploadToS3(file string) (S3Response, error) {
	sess, err := NewS3Session()
	if err != nil {
		return S3Response{}, err
	}
	bucket := os.Getenv("S3_BUCKET")
	acl := os.Getenv("AWS_ACL")
	f, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return S3Response{}, err
	}
	defer f.Close()
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String(acl),
		Bucket: aws.String(bucket),
		Key:    aws.String(f.Name()),
		Body:   f,
	})
	if err != nil {
		return S3Response{}, err
	}
	log.Info(fmt.Sprintf("Upload result: %+v\n", result))
	response := S3Response{
		URL:      result.Location,
		FileName: file,
	}
	return response, nil
}

func DownloadFromS3(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Info("Error creating temporary file for S3 download: ", err)
		return file, err
	}
	defer file.Close()
	sess, err := NewS3Session()
	if err != nil {
		return file, err
	}
	bucket := os.Getenv("S3_BUCKET")
	downloader := s3manager.NewDownloader(sess)
	// number of bytes downloaded or error
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileName),
		},
	)
	if err != nil {
		log.Info("Error downloading from S3: ", err.Error())
		return file, err
	}
	log.Println("Successfully downloaded")
	return file, nil
}

func DeleteFromS3(fileName string) error {
	sess, err := NewS3Session()
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
