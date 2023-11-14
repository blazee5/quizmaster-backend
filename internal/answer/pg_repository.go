package answer

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, input domain.Answer) (int, error)
	Update(ctx context.Context, answerId int, input domain.Answer) error
	Delete(ctx context.Context, answerId int) error
}
