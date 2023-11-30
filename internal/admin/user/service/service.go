package service

import (
	"context"
	adminUserRepo "github.com/blazee5/quizmaster-backend/internal/admin/user"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	repo adminUserRepo.Repository
}

func NewService(log *zap.SugaredLogger, repo adminUserRepo.Repository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, input domain.SignUpRequest) (int, error) {
	input.Password = authLib.GenerateHashPassword(input.Password)
	return s.repo.Create(ctx, input)
}

func (s *Service) GetUsers(ctx context.Context) ([]models.ShortUser, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateUser(ctx context.Context, id int, input domain.User) error {
	return s.repo.Update(ctx, id, input)
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
