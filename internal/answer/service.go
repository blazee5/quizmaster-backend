package answer

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	Create(ctx context.Context, userID, quizID, questionID int) (int, error)
	GetByQuestionID(ctx context.Context, quizID, questionID int) ([]models.AnswerInfo, error)
	GetByQuestionIDForUser(ctx context.Context, quizID, questionID, userID int) ([]models.Answer, error)
	Update(ctx context.Context, answerID, userID, quizID int, input domain.Answer) error
	Delete(ctx context.Context, answerID, userID, quizID int) error
	ChangeOrder(ctx context.Context, userID, quizID, questionID int, input domain.ChangeAnswerOrder) error
}
