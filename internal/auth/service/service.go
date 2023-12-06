package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.SugaredLogger
	repo   auth.Repository
	tracer trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo auth.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, tracer: tracer}
}

func (s *Service) SignUp(ctx context.Context, input domain.SignUpRequest) (int, error) {
	ctx, span := s.tracer.Start(ctx, "authService.SignUp")
	defer span.End()

	input.Password = authLib.GenerateHashPassword(input.Password)
	id, err := s.repo.CreateUser(ctx, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, err
}

func (s *Service) GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error) {
	ctx, span := s.tracer.Start(ctx, "authService.GenerateToken")
	defer span.End()

	input.Password = authLib.GenerateHashPassword(input.Password)
	user, err := s.repo.ValidateUser(ctx, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return "", err
	}

	return authLib.GenerateToken(user.ID, user.RoleID)
}
