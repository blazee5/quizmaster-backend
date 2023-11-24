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

func (s *Service) Create(ctx context.Context, userID, quizID, questionID int) (int, error) {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return 0, err
	}

	if quiz.UserID != userID {
		return 0, http_errors.PermissionDenied
	}

	return s.repo.Create(ctx, questionID)
}

func (s *Service) Update(ctx context.Context, answerID, userID, quizID int, input domain.Answer) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.Update(ctx, answerID, input)
}

func (s *Service) Delete(ctx context.Context, answerID, userID, quizID int) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, answerID)
}

func (s *Service) ChangeOrder(ctx context.Context, userID, quizID int, input domain.ChangeAnswerOrder) error {
	//TODO implement me
	panic("implement me")
}
