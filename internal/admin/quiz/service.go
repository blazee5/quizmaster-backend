package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	CreateQuiz(ctx context.Context, userId int, input domain.Quiz) (int, error)
	GetQuizzes(ctx context.Context) ([]models.Quiz, error)
	UpdateQuiz(ctx context.Context, id int, input domain.Quiz) error
	DeleteQuiz(ctx context.Context, id int) error
}
