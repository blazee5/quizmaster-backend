package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	GetById(ctx context.Context, userId int) (models.User, error)
	GetQuizzes(ctx context.Context, userId int) ([]models.Quiz, error)
	GetResults(ctx context.Context, userId int) ([]models.Quiz, error)
	ChangeAvatar(ctx context.Context, userId int, file string) error
	Update(ctx context.Context, userId int, input domain.UpdateUser) error
	Delete(ctx context.Context, userId int) error
}
