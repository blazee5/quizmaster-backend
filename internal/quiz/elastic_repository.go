package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type ElasticRepository interface {
	CreateIndex(ctx context.Context, input domain.Quiz) error
	UpdateIndex(ctx context.Context, input domain.Quiz) error
	DeleteIndex(ctx context.Context, input domain.Quiz) error
	SearchIndex(ctx context.Context, input string) ([]domain.Quiz, error)
}
