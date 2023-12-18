package question

import "context"

type AWSRepository interface {
	SaveFile(ctx context.Context, fileName, contentType string, chunk []byte) error
	DeleteFile(ctx context.Context, fileName string) error
}
