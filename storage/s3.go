package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	appConfig "github.com/brain-flowing-company/pprp-backend/config"
)

type Storage interface {
	Upload(string, io.Reader) (string, error)
}

type storageImpl struct {
	client     *s3.Client
	uploader   *manager.Uploader
	bucketName string
}

func New(appConfig *appConfig.Config) (Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &storageImpl{
		client:     client,
		uploader:   uploader,
		bucketName: appConfig.S3Bucket,
	}, nil
}

func (s *storageImpl) Upload(filename string, file io.Reader) (string, error) {
	result, err := s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
