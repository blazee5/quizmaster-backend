package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/question"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.uber.org/zap"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     question.Repository
	quizRepo quizRepo.Repository
}

func NewService(log *zap.SugaredLogger, repo question.Repository, quizRepo quizRepo.Repository) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo}
}

func (s *Service) Create(ctx context.Context, userId, quizId int, input domain.Question) (int, error) {
	quiz, err := s.quizRepo.GetById(ctx, quizId)

	if err != nil {
		return 0, err
	}

	if quiz.UserId != userId {
		return 0, http_errors.PermissionDenied
	}

	return s.repo.CreateQuestion(ctx, quizId, input)
}

func (s *Service) GetQuestionsById(ctx context.Context, id int) ([]models.Question, error) {
	return s.repo.GetQuestionsById(ctx, id, false)
}

func (s *Service) Update(ctx context.Context, id, userId, quizId int, input domain.Question) error {
	quiz, err := s.quizRepo.GetById(ctx, quizId)

	if err != nil {
		return err
	}

	if quiz.UserId != userId {
		return http_errors.PermissionDenied
	}

	return s.repo.Update(ctx, id, input)
}

func (s *Service) Delete(ctx context.Context, id, userId, quizId int) error {
	quiz, err := s.quizRepo.GetById(ctx, quizId)

	if err != nil {
		return err
	}

	if quiz.UserId != userId {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, id)
}
