package service

import (
	"context"
	adminQuizRepo "github.com/blazee5/quizmaster-backend/internal/admin/quiz"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.SugaredLogger
	repo   adminQuizRepo.Repository
	tracer trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo adminQuizRepo.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, tracer: tracer}
}

func (s *Service) CreateQuiz(ctx context.Context, userID int, input domain.Quiz) (int, error) {
	ctx, span := s.tracer.Start(ctx, "admin.quizService.CreateQuiz")
	defer span.End()

	id, err := s.repo.Create(ctx, userID, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (s *Service) GetQuizzes(ctx context.Context) ([]models.Quiz, error) {
	ctx, span := s.tracer.Start(ctx, "admin.quizService.GetQuizzes")
	defer span.End()

	quizzes, err := s.repo.GetQuizzes(ctx)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return quizzes, nil
}

func (s *Service) UpdateQuiz(ctx context.Context, id int, input domain.Quiz) error {
	ctx, span := s.tracer.Start(ctx, "admin.quizService.UpdateQuiz")
	defer span.End()

	err := s.repo.Update(ctx, id, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) DeleteQuiz(ctx context.Context, id int) error {
	ctx, span := s.tracer.Start(ctx, "admin.quizService.DeleteQuiz")
	defer span.End()

	err := s.repo.Delete(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
