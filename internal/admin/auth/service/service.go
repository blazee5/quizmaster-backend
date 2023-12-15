package service

import (
	"context"
	adminAuthRepo "github.com/blazee5/quizmaster-backend/internal/admin/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.SugaredLogger
	repo   adminAuthRepo.Repository
	tracer trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo adminAuthRepo.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, tracer: tracer}
}

func (s *Service) GenerateToken(ctx context.Context, input domain.SignInRequest) (string, error) {
	ctx, span := s.tracer.Start(ctx, "admin.authService.GenerateToken")
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
