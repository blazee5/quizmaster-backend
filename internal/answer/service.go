package answer

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userId, quizId int, input domain.Answer) (int, error)
	Update(ctx context.Context, answerId, userId, quizId int, input domain.Answer) error
	Delete(ctx context.Context, answerId, userId, quizId int) error
}
