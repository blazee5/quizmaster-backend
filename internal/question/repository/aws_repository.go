package repository

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
)

const (
	quizzesBucketName = "quizzes"
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

	bucketExists, err := s.client.BucketExists(ctx, quizzesBucketName)

	if err != nil {
		return err
	}

	if !bucketExists {
		err := s.client.MakeBucket(ctx, quizzesBucketName, minio.MakeBucketOptions{})

		if err != nil {
			return err
		}
	}

	_, err = s.client.PutObject(ctx, quizzesBucketName, fileName, file, file.Size(), options)

	if err != nil {
		return err
	}

	return nil
}

func (s *AWSRepository) DeleteFile(ctx context.Context, fileName string) error {
	if err := s.client.RemoveObject(ctx, quizzesBucketName, fileName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}
