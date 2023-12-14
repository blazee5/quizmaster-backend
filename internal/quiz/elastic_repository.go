package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type ElasticRepository interface {
	CreateIndex(ctx context.Context, input models.QuizInfo) error
	SearchIndex(ctx context.Context, input string) ([]models.QuizInfo, error)
	UpdateIndex(ctx context.Context, id int, input models.QuizInfo) error
	DeleteIndex(ctx context.Context, ID int) error
}
