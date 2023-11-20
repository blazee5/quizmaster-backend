package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	ValidateUser(ctx context.Context, input domain.SignInRequest) (models.User, error)
}
