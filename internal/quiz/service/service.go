package service

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/blazee5/testhub-backend/internal/quiz"
)

type Service struct {
	repo quiz.Repository
}

func NewService(repo quiz.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input domain.Quiz) (int, error) {
	return s.repo.Create(ctx, input)
}

func (s *Service) GetById(ctx context.Context, id int) (models.Quiz, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service) GetQuestionsById(ctx context.Context, id int) ([]models.Question, error) {
	return s.repo.GetQuestionsById(ctx, id, false)
}

func (s *Service) SaveResult(ctx context.Context, userId int, input domain.Result) (int, error) {
	return s.repo.SaveResult(ctx, userId, input)
}
