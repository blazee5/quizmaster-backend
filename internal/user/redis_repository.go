package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*models.UserInfo, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.UserInfo) error
	DeleteUserCtx(ctx context.Context, key string) error
}
