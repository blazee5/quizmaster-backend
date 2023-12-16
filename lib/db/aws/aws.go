package aws

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func NewAWSClient() *minio.Client {
	client, err := minio.New(os.Getenv("AWS_HOST"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("AWS_USER"), os.Getenv("AWS_PASSWORD"), os.Getenv("AWS_TOKEN")),
		Secure: false,
	})

	if err != nil {
		log.Fatalf("error while connect to minio s3: %v", err)
	}

	return client
}
