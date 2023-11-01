package quiz

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, input domain.Quiz) (int, error)
	GetById(ctx context.Context, id int) (models.Quiz, error)
	GetQuestionsById(ctx context.Context, id int, includeIsCorrect bool) ([]models.Question, error)
	SaveResult(ctx context.Context, userId int, input domain.Result) (int, error)
}
