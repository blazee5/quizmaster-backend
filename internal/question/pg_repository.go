package question

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	CreateQuestion(ctx context.Context, quizId int, input domain.Question) (int, error)
	GetQuestionsById(ctx context.Context, id int, includeIsCorrect bool) ([]models.Question, error)
	Update(ctx context.Context, id int, input domain.Question) error
	Delete(ctx context.Context, id int) error
}
