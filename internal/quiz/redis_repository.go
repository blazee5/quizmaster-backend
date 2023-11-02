package quiz

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
)

type RedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.Quiz, error)
	SetQuizCtx(ctx context.Context, key string, seconds int, user *models.Quiz) error
	DeleteQuizCtx(ctx context.Context, key string) error
}
