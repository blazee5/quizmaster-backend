package quiz

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/models"
)

type Repository interface {
	GetById(ctx context.Context, id int) (models.Quiz, error)
	GetAll(ctx context.Context) ([]models.Quiz, error)
	GetQuestionsById(ctx context.Context, id int, includeIsCorrect bool) ([]models.Question, error)
	Create(ctx context.Context, input domain.Quiz) (int, error)
	SaveResult(ctx context.Context, userId, quizId int, input domain.Result) (int, error)
	Delete(ctx context.Context, quizId int) error
}
