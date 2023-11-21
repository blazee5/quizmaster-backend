package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, userId int, input domain.Quiz) (int, error)
	GetQuizzes(ctx context.Context) ([]models.Quiz, error)
	Update(ctx context.Context, id int, input domain.Quiz) error
	Delete(ctx context.Context, id int) error
}
