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

func (s *Service) Create(ctx context.Context, userID, quizID int) (int, error) {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return 0, err
	}

	if quiz.UserID != userID {
		return 0, http_errors.PermissionDenied
	}

	return s.repo.CreateQuestion(ctx, quizID)
}

func (s *Service) GetQuestionsByID(ctx context.Context, id int) ([]models.Question, error) {
	return s.repo.GetQuestionsByID(ctx, id)
}

func (s *Service) GetAllQuestionsByID(ctx context.Context, id, userID int) ([]models.QuestionWithAnswers, error) {
	quiz, err := s.quizRepo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if quiz.UserID != userID {
		return nil, http_errors.PermissionDenied
	}

	return s.repo.GetQuestionsWithAnswers(ctx, id)
}

func (s *Service) Update(ctx context.Context, id, userID, quizID int, input domain.Question) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.Update(ctx, id, input)
}

func (s *Service) Delete(ctx context.Context, id, userID, quizID int) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) UploadImage(ctx context.Context, id, userID, quizID int, filename string) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.UploadImage(ctx, id, filename)
}

func (s *Service) DeleteImage(ctx context.Context, id, userID, quizID int) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.DeleteImage(ctx, id)
}

func (s *Service) ChangeOrder(ctx context.Context, userID, quizID int, input domain.ChangeQuestionOrder) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	order := (input.FirstOrderID + input.SecondOrderID) / 2

	return s.repo.ChangeOrder(ctx, input.QuestionID, order)
}
