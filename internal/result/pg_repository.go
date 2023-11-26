package result

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (models.Result, error)
	GetByUserID(ctx context.Context, id int) (models.Result, error)
	GetUserAnswerByID(ctx context.Context, answerID, resultID int) (bool, error)
	NewResult(ctx context.Context, userID, quizID int) (int, error)
	UpdateResult(ctx context.Context, id, userID, score int) error
	SaveUserAnswer(ctx context.Context, userID, questionID, answerID, resultID int, answerText string) error
}
