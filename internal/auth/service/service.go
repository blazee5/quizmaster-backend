package service

import (
	"context"
	authRepo "github.com/blazee5/testhub-backend/internal/auth/repository"
	"github.com/blazee5/testhub-backend/internal/domain"
	auth "github.com/blazee5/testhub-backend/lib/auth"
)

type Service struct {
	repo authRepo.Repository
}

func NewService(repo authRepo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignUp(ctx context.Context, input domain.SignUpRequest) (int, error) {
	input.Password = auth.GenerateHashPassword(input.Password)
	return s.repo.CreateUser(ctx, input)
}

func (s *Service) GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error) {
	input.Password = auth.GenerateHashPassword(input.Password)
	user, err := s.repo.ValidateUser(ctx, input)

	if err != nil {
		return "", err
	}

	return auth.GenerateToken(user.Id)
}
