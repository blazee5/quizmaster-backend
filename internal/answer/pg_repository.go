package answer

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, questionID int) (int, error)
	Update(ctx context.Context, answerID int, input domain.Answer) error
	Delete(ctx context.Context, answerID int) error
}
