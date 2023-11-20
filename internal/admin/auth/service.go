package user

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type Service interface {
	GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error)
}
