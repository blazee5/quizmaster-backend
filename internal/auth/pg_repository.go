package auth

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, input domain.SignUpRequest) (int, error)
	ValidateUser(ctx context.Context, input domain.SignInRequest) (models.User, error)
	UpdateEmail(ctx context.Context, userID int, email string) error
	UpdatePassword(ctx context.Context, userID int, password string) error
	CreateVerificationCode(ctx context.Context, userID int, codeType, code, email string) error
	GetVerificationCode(ctx context.Context, code, codeType string) (models.VerificationCode, error)
	DeleteVerificationCode(ctx context.Context, id int) error
}
