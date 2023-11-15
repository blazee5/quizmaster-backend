package question

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	Create(ctx context.Context, userId, quizId int) (int, error)
	GetQuestionsById(ctx context.Context, id int) ([]models.Question, error)
	Update(ctx context.Context, id, userId, quizId int, input domain.Question) error
	Delete(ctx context.Context, id, userId, quizId int) error
	UploadImage(ctx context.Context, id, userId, quizId int) error
	DeleteImage(ctx context.Context, id, userId, quizId int) error
}
