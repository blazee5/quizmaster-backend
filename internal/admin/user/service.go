package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Service interface {
	CreateUser(ctx context.Context, input domain.SignUpRequest) (int, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, id int, input domain.User) error
	DeleteUser(ctx context.Context, id int) error
}
