package user

import "context"

type AWSRepository interface {
	SaveFile(ctx context.Context, fileName, contentType string, chunk []byte) error
}
