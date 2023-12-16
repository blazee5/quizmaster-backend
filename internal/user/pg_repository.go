package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	GetByID(ctx context.Context, userID int) (models.UserInfo, error)
	GetAvatarByID(ctx context.Context, userID int) (string, error)
	GetQuizzes(ctx context.Context, userID int) ([]models.Quiz, error)
	GetResults(ctx context.Context, userID int) ([]models.Quiz, error)
	ChangeAvatar(ctx context.Context, userID int, file string) error
	Update(ctx context.Context, userID int, input domain.UpdateUser) error
	Delete(ctx context.Context, userID int) error
}
