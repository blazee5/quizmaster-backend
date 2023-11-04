package user

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
)

type Repository interface {
	GetById(ctx context.Context, userId int) (models.User, error)
	GetQuizzes(ctx context.Context, userId int) ([]models.Quiz, error)
	GetResults(ctx context.Context, userId int) ([]models.Quiz, error)
	Delete(ctx context.Context, userId int) error
}
