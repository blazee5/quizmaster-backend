package user

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
)

type Service interface {
	GetById(ctx context.Context, userId int) (models.User, error)
	GetQuizzes(ctx context.Context, userId int) ([]models.Quiz, error)
	GetResults(ctx context.Context, userId int) ([]models.Quiz, error)
	ChangeAvatar(ctx context.Context, userId int, file string) error
	Delete(ctx context.Context, userId int) error
}
