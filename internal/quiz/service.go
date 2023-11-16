package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	Create(ctx context.Context, userId int, input domain.Quiz) (int, error)
	GetAll(ctx context.Context) ([]models.Quiz, error)
	GetById(ctx context.Context, id int) (models.Quiz, error)
	SaveResult(ctx context.Context, userId int, quizId int, input domain.Result) (int, error)
	Delete(ctx context.Context, userId, quizId int) error
	UploadImage(ctx context.Context, userId, quizId int, filename string) error
	DeleteImage(ctx context.Context, userId, quizId int) error
}
