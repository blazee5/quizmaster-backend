package repository

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
)

const (
	quizBucketName = "quizzes"
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

	bucketExists, err := s.client.BucketExists(ctx, quizBucketName)

	if err != nil {
		return err
	}

	if !bucketExists {
		err := s.client.MakeBucket(ctx, quizBucketName, minio.MakeBucketOptions{})

		if err != nil {
			return err
		}
	}

	_, err = s.client.PutObject(ctx, quizBucketName, fileName, file, file.Size(), options)

	if err != nil {
		return err
	}

	return nil
}

func (s *AWSRepository) DeleteFile(ctx context.Context, fileName string) error {
	if err := s.client.RemoveObject(ctx, quizBucketName, fileName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}
