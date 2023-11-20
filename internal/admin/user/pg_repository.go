package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, input domain.SignUpRequest) (int, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, id int, input domain.User) error
	Delete(ctx context.Context, id int) error
}
