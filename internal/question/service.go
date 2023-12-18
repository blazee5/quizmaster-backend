package question

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"mime/multipart"
)

type Service interface {
	Create(ctx context.Context, userID, quizID int) (int, error)
	GetQuestionsByID(ctx context.Context, id int) ([]models.Question, error)
	Update(ctx context.Context, id, userID, quizID int, input domain.Question) error
	Delete(ctx context.Context, id, userID, quizID int) error
	UploadImage(ctx context.Context, id, userID, quizID int, fileHeader *multipart.FileHeader) error
	DeleteImage(ctx context.Context, id, userID, quizID int) error
	ChangeOrder(ctx context.Context, userId, quizId int, input domain.ChangeQuestionOrder) error
}
