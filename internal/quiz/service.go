package quiz

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/models"
)

type Service interface {
	Create(ctx context.Context, input domain.Quiz) (int, error)
	GetById(ctx context.Context, id int) (models.Quiz, error)
	GetQuestionsById(ctx context.Context, id int) ([]models.Question, error)
	SaveResult(ctx context.Context, userId int, input domain.Result) (int, error)
}
