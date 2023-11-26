package answer

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, questionID int, input domain.CreateAnswer) (int, error)
	GetByID(ctx context.Context, id int) (models.Answer, error)
	GetAnswersByQuestionID(ctx context.Context, id int) ([]models.Answer, error)
	Update(ctx context.Context, answerID int, input domain.Answer) error
	Delete(ctx context.Context, answerID int) error
	ChangeOrder(ctx context.Context, questionID int, input domain.ChangeAnswerOrder) error
}
