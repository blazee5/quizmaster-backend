package service

import (
	"context"
	adminQuizRepo "github.com/blazee5/quizmaster-backend/internal/admin/quiz"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	repo adminQuizRepo.Repository
}

func NewService(log *zap.SugaredLogger, repo adminQuizRepo.Repository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) CreateQuiz(ctx context.Context, userId int, input domain.Quiz) (int, error) {
	return s.repo.Create(ctx, userId, input)
}

func (s *Service) GetQuizzes(ctx context.Context) ([]models.Quiz, error) {
	return s.repo.GetQuizzes(ctx)
}

func (s *Service) UpdateQuiz(ctx context.Context, id int, input domain.Quiz) error {
	return s.repo.Update(ctx, id, input)
}

func (s *Service) DeleteQuiz(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
