package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	Create(ctx context.Context, userID int, input domain.Quiz) (int, error)
	GetAll(ctx context.Context) ([]models.Quiz, error)
	GetByID(ctx context.Context, id int) (models.Quiz, error)
	SaveResult(ctx context.Context, userID int, quizID int, input domain.Result) (int, error)
	Update(ctx context.Context, userID, quizID int, input domain.Quiz) error
	Delete(ctx context.Context, userID, quizID int) error
	UploadImage(ctx context.Context, userID, quizID int, filename string) error
	DeleteImage(ctx context.Context, userID, quizID int) error
}
