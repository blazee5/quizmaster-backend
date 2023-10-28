package auth

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, input domain.SignUpRequest) (int, error)
	ValidateUser(ctx context.Context, input domain.SignInRequest) (models.User, error)
}
