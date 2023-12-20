package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/answer"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     answer.Repository
	quizRepo quizRepo.Repository
	tracer   trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo answer.Repository, quizRepo quizRepo.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo, tracer: tracer}
}

func (s *Service) Create(ctx context.Context, userID, quizID, questionID int) (int, error) {
	ctx, span := s.tracer.Start(ctx, "answerService.Create")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return 0, err
	}

	if quiz.UserID != userID {
		return 0, http_errors.PermissionDenied
	}

	return s.repo.Create(ctx, questionID)
}

func (s *Service) GetByQuestionID(ctx context.Context, quizID, questionID int) ([]models.AnswerInfo, error) {
	ctx, span := s.tracer.Start(ctx, "answerService.GetByQuestionID")
	defer span.End()

	_, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return nil, err
	}

	return s.repo.GetAnswersInfoByQuestionID(ctx, questionID)
}

func (s *Service) GetByQuestionIDForUser(ctx context.Context, quizID, questionID, userID int) ([]models.Answer, error) {
	ctx, span := s.tracer.Start(ctx, "answerService.GetByQuestionIDForUser")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return nil, err
	}

	if quiz.UserID != userID {
		return nil, http_errors.PermissionDenied
	}

	return s.repo.GetAnswersByQuestionID(ctx, questionID)
}

func (s *Service) Update(ctx context.Context, answerID, userID, quizID int, input domain.Answer) error {
	ctx, span := s.tracer.Start(ctx, "answerService.Update")
	defer span.End()

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
	ctx, span := s.tracer.Start(ctx, "answerService.Delete")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, answerID)
}

func (s *Service) ChangeOrder(ctx context.Context, userID, quizID, questionID int, input domain.ChangeAnswerOrder) error {
	ctx, span := s.tracer.Start(ctx, "answerService.ChangeOrder")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.ChangeOrder(ctx, questionID, input)

	if err != nil {
		return err
	}

	return nil
}
