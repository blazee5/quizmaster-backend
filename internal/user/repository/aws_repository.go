package repository

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
)

const (
	avatarsBucketName = "avatars"
)

type AWSRepository struct {
	client *minio.Client
}

func NewAWSRepository(client *minio.Client) *AWSRepository {
	return &AWSRepository{client: client}
}

func (s *AWSRepository) SaveFile(ctx context.Context, fileName, contentType string, chunk []byte) error {
	options := minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	file := bytes.NewReader(chunk)

	bucketExists, err := s.client.BucketExists(ctx, avatarsBucketName)

	if err != nil {
		return err
	}

	if !bucketExists {
		err := s.client.MakeBucket(ctx, avatarsBucketName, minio.MakeBucketOptions{})

		if err != nil {
			return err
		}
	}

	_, err = s.client.PutObject(ctx, avatarsBucketName, fileName, file, file.Size(), options)

	if err != nil {
		return err
	}

	return nil
}
