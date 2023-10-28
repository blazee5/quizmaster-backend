package service

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
	userRepo "github.com/blazee5/testhub-backend/internal/user/repository"
)

type Service struct {
	repo userRepo.Repository
}

func NewService(repo userRepo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserById(ctx context.Context, userId int) (models.User, error) {
	return s.repo.GetUserById(ctx, userId)
}
