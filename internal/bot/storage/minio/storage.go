package minio

import (
	"PDFHub/internal/bot/lib/e"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Minio      *minio.Client
	BucketName string
	ctx        context.Context
}

func New(endpoint string, accessKeyID string, secretAccessKey string, bucketName string, c context.Context) (
	*Storage,
	error,
) {
	minioClient, err := minio.New(
		endpoint, &minio.Options{
			Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		},
	)
	if err != nil {
		return nil, e.Wrap("can't create a storage", err)
	}

	return &Storage{
		Minio:      minioClient,
		BucketName: bucketName,
		ctx:        c,
	}, nil
}

func (s Storage) Save(file *http.Response) error {
	if s.existBucket(s.BucketName) == false {
		if err := s.createBucket(s.BucketName); err != nil {
			return e.Wrap("can't create a bucket", err)
		}
	}

	defer file.Body.Close()

	_, err := s.Minio.PutObject(
		s.ctx,
		s.BucketName,
		uuid.New().String(),
		file.Body,
		file.ContentLength,
		minio.PutObjectOptions{},
	)

	if err != nil {
		return e.Wrap("can't upload file", err)
	}

	return nil
}

func (s Storage) existBucket(bucketName string) bool {
	exists, err := s.Minio.BucketExists(s.ctx, bucketName)
	if err != nil {
		return false
	}
	return exists
}

func (s Storage) createBucket(bucketName string) error {
	err := s.Minio.MakeBucket(s.ctx, bucketName, minio.MakeBucketOptions{})

	if err != nil {
		return e.Wrap("can't create a bucket", err)
	}

	return nil
}
