package service

import (
	"context"
	adminUserRepo "github.com/blazee5/quizmaster-backend/internal/admin/user"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.SugaredLogger
	repo   adminUserRepo.Repository
	tracer trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo adminUserRepo.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, tracer: tracer}
}

func (s *Service) CreateUser(ctx context.Context, input domain.SignUpRequest) (int, error) {
	ctx, span := s.tracer.Start(ctx, "admin.userService.CreateUser")
	defer span.End()

	input.Password = authLib.GenerateHashPassword(input.Password)
	id, err := s.repo.Create(ctx, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]models.ShortUser, error) {
	ctx, span := s.tracer.Start(ctx, "admin.userService.GetUsers")
	defer span.End()

	users, err := s.repo.GetAll(ctx)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return users, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int, input domain.User) error {
	ctx, span := s.tracer.Start(ctx, "admin.userService.UpdateUser")
	defer span.End()

	err := s.repo.Update(ctx, id, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	ctx, span := s.tracer.Start(ctx, "admin.userService.DeleteUser")
	defer span.End()

	err := s.repo.Delete(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
