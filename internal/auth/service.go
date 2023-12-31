package auth

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type Service interface {
	SignUp(ctx context.Context, input domain.SignUpRequest) (int, error)
	GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error)
	SendEmailCode(ctx context.Context, userID int, input domain.VerificationCode) error
	SendPasswordCode(ctx context.Context, input domain.VerificationCode) error
	ResetEmail(ctx context.Context, userID int, input domain.ResetEmailRequest) error
	ResetPassword(ctx context.Context, input domain.ResetPasswordRequest) error
}
