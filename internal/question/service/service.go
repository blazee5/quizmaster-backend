package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/question"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     question.Repository
	quizRepo quizRepo.Repository
	tracer   trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo question.Repository, quizRepo quizRepo.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo, tracer: tracer}
}

func (s *Service) Create(ctx context.Context, userID, quizID int) (int, error) {
	ctx, span := s.tracer.Start(ctx, "questionService.Create")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	if quiz.UserID != userID {
		return 0, http_errors.PermissionDenied
	}

	id, err := s.repo.CreateQuestion(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (s *Service) GetQuestionsByID(ctx context.Context, id int) ([]models.Question, error) {
	ctx, span := s.tracer.Start(ctx, "questionService.GetQuestionsByID")
	defer span.End()

	questions, err := s.repo.GetQuestionsByQuizID(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return questions, nil
}

func (s *Service) Update(ctx context.Context, id, userID, quizID int, input domain.Question) error {
	ctx, span := s.tracer.Start(ctx, "questionService.Update")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.Update(ctx, id, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id, userID, quizID int) error {
	ctx, span := s.tracer.Start(ctx, "questionService.Delete")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.Delete(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) UploadImage(ctx context.Context, id, userID, quizID int, filename string) error {
	ctx, span := s.tracer.Start(ctx, "questionService.UploadImage")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.UploadImage(ctx, id, filename)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) DeleteImage(ctx context.Context, id, userID, quizID int) error {
	ctx, span := s.tracer.Start(ctx, "questionService.DeleteImage")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.DeleteImage(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) ChangeOrder(ctx context.Context, userID, quizID int, input domain.ChangeQuestionOrder) error {
	ctx, span := s.tracer.Start(ctx, "questionService.ChangeOrder")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.ChangeOrder(ctx, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
