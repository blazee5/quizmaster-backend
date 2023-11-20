package service

import (
	"context"
	adminAuthRepo "github.com/blazee5/quizmaster-backend/internal/admin/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	repo adminAuthRepo.Repository
}

func NewService(log *zap.SugaredLogger, repo adminAuthRepo.Repository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error) {
	input.Password = authLib.GenerateHashPassword(input.Password)
	user, err := s.repo.ValidateUser(ctx, input)

	if err != nil {
		return "", err
	}

	return authLib.GenerateToken(user.Id, user.RoleId)
}
