package question

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	Create(ctx context.Context, userId, quizId int, input domain.Question) (int, error)
	GetQuestionsById(ctx context.Context, id int) ([]models.Question, error)
}
