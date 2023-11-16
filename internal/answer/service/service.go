package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/answer"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.uber.org/zap"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     answer.Repository
	quizRepo quizRepo.Repository
}

func NewService(log *zap.SugaredLogger, repo answer.Repository, quizRepo quizRepo.Repository) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo}
}

func (s *Service) Create(ctx context.Context, userId, quizId, questionId int) (int, error) {
	quiz, err := s.quizRepo.GetById(ctx, quizId)

	if err != nil {
		return 0, err
	}

	if quiz.UserId != userId {
		return 0, http_errors.PermissionDenied
	}

	return s.repo.Create(ctx, questionId)
}

func (s *Service) Update(ctx context.Context, answerId, userId, quizId int, input domain.Answer) error {
	quiz, err := s.quizRepo.GetById(ctx, quizId)

	if err != nil {
		return err
	}

	if quiz.UserId != userId {
		return http_errors.PermissionDenied
	}

	return s.repo.Update(ctx, answerId, input)
}

func (s *Service) Delete(ctx context.Context, answerId, userId, quizId int) error {
	quiz, err := s.quizRepo.GetById(ctx, quizId)

	if err != nil {
		return err
	}

	if quiz.UserId != userId {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, answerId)
}
