package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	repo auth.Repository
}

func NewService(log *zap.SugaredLogger, repo auth.Repository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) SignUp(ctx context.Context, input domain.SignUpRequest) (int, error) {
	input.Password = authLib.GenerateHashPassword(input.Password)
	return s.repo.CreateUser(ctx, input)
}

func (s *Service) GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error) {
	input.Password = authLib.GenerateHashPassword(input.Password)
	user, err := s.repo.ValidateUser(ctx, input)

	if err != nil {
		return "", err
	}

	return authLib.GenerateToken(user.Id)
}
