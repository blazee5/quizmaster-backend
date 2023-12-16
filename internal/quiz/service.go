package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"mime/multipart"
)

type Service interface {
	Create(ctx context.Context, userID int, input domain.Quiz) (int, error)
	GetAll(ctx context.Context) ([]models.Quiz, error)
	GetByID(ctx context.Context, id int) (models.Quiz, error)
	Update(ctx context.Context, userID, quizID int, input domain.Quiz) error
	Delete(ctx context.Context, userID, quizID int) error
	Search(ctx context.Context, title string) ([]models.QuizInfo, error)
	UploadImage(ctx context.Context, userID, quizID int, fileHeader *multipart.FileHeader) error
	DeleteImage(ctx context.Context, userID, quizID int) error
}
