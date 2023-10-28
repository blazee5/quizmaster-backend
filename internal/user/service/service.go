package service

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/blazee5/testhub-backend/internal/user"
	userRepo "github.com/blazee5/testhub-backend/internal/user/repository"
)

type Service struct {
	repo user.Repository
}

func NewService(repo userRepo.Repository) *Service {
	return &Service{repo: &repo}
}

func (s *Service) GetUserById(ctx context.Context, userId int) (models.User, error) {
	return s.repo.GetUserById(ctx, userId)
}
