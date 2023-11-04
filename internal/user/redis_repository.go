package user

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
)

type RedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
