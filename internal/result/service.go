package result

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type Service interface {
	NewResult(ctx context.Context, userID int, quizID int) (int, error)
	SaveUserAnswer(ctx context.Context, userID, quizID int, input domain.UserAnswer) error
}
