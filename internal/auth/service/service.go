package service

import (
	"context"
	"encoding/json"
	"github.com/blazee5/quizmaster-backend/internal/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/rabbitmq"
	userRepo "github.com/blazee5/quizmaster-backend/internal/user"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/blazee5/quizmaster-backend/lib/random"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     auth.Repository
	userRepo userRepo.Repository
	producer *rabbitmq.Producer
	tracer   trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo auth.Repository, userRepo userRepo.Repository, producer *rabbitmq.Producer, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, userRepo: userRepo, producer: producer, tracer: tracer}
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

func (s *Service) SendCode(ctx context.Context, userID int, input domain.VerificationCode) error {
	ctx, span := s.tracer.Start(ctx, "userService.SendCode")
	defer span.End()

	user, err := s.userRepo.GetByID(ctx, userID)

	if err != nil {
		return err
	}

	code := random.GenerateVerificationCode(8)

	if err != nil {
		return err
	}

	err = s.repo.CreateVerificationCode(ctx, userID, input.Type, code)

	if err != nil {
		return err
	}

	if input.Email == "" {
		input.Email = user.User.Email
	}

	email := domain.Email{
		Type:    input.Type,
		To:      input.Email,
		Message: code,
	}

	bytes, err := json.Marshal(&email)

	err = s.producer.PublishMessage(ctx, bytes)

	if err != nil {
		return err
	}

	return nil
}
