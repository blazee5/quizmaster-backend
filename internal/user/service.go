package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"mime/multipart"
)

type Service interface {
	GetByID(ctx context.Context, userID int) (models.UserInfo, error)
	ChangeAvatar(ctx context.Context, userID int, fileHeader *multipart.FileHeader) error
	Update(ctx context.Context, userID int, input domain.UpdateUser) error
	Delete(ctx context.Context, userID int) error
}
