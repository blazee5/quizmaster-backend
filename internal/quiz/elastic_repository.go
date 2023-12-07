package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type ElasticRepository interface {
	CreateIndex(ctx context.Context, input models.Quiz) error
	//UpdateIndex(ctx context.Context, input domain.Quiz) error
	//DeleteIndex(ctx context.Context, input domain.Quiz) error
	SearchIndex(ctx context.Context, input string) ([]models.QuizInfo, error)
}
