package auth

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
)

type Service interface {
	SignUp(ctx context.Context, input domain.SignUpRequest) (int, error)
	GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error)
}
