package result

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	NewResult(ctx context.Context, userID int, quizID int) (int, error)
	SaveUserAnswer(ctx context.Context, userID, quizID int, input domain.UserAnswer) error
	GetResultsByQuizID(ctx context.Context, quizID int) ([]models.UsersResult, error)
	SubmitResult(ctx context.Context, userID, quizID int, input domain.SubmitResult) (models.UsersResult, error)
}
