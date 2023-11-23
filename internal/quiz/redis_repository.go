package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*models.Quiz, error)
	SetQuizCtx(ctx context.Context, key string, seconds int, quiz *models.Quiz) error
	DeleteQuizCtx(ctx context.Context, key string) error
}
