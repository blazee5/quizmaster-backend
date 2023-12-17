package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type ElasticRepository interface {
	CreateIndex(ctx context.Context, input models.Quiz) error
	SearchIndex(ctx context.Context, input, sortBy, sortDir string, offset, size int) (models.QuizList, error)
	UpdateIndex(ctx context.Context, id int, input models.Quiz) error
	DeleteIndex(ctx context.Context, ID int) error
}
