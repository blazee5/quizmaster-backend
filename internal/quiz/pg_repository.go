package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (models.Quiz, error)
	GetAll(ctx context.Context, sortBy, sortDir string, page, size int) (models.QuizList, error)
	Create(ctx context.Context, userID int, input domain.Quiz) (models.Quiz, error)
	Update(ctx context.Context, quizID int, input domain.Quiz) (models.Quiz, error)
	Delete(ctx context.Context, quizID int) error
	UploadImage(ctx context.Context, id int, filename string) error
	DeleteImage(ctx context.Context, id int) error
}
