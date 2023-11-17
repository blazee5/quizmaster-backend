package question

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	CreateQuestion(ctx context.Context, quizId int) (int, error)
	GetQuestionsById(ctx context.Context, id int) ([]models.Question, error)
	GetQuestionsWithAnswers(ctx context.Context, id int) ([]models.QuestionWithAnswers, error)
	Update(ctx context.Context, id int, input domain.Question) error
	Delete(ctx context.Context, id int) error
	UploadImage(ctx context.Context, id int, filename string) error
	DeleteImage(ctx context.Context, id int) error
}
