package service

import (
	"context"
	"encoding/json"
	"github.com/blazee5/quizmaster-backend/internal/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/rabbitmq"
	userRepo "github.com/blazee5/quizmaster-backend/internal/user"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/blazee5/quizmaster-backend/lib/random"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     auth.Repository
	userRepo userRepo.Repository
	producer rabbitmq.QueueProducer
	tracer   trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo auth.Repository, userRepo userRepo.Repository, producer rabbitmq.QueueProducer, tracer trace.Tracer) *Service {
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

func (s *Service) SendEmailCode(ctx context.Context, userID int, input domain.VerificationCode) error {
	ctx, span := s.tracer.Start(ctx, "authService.SendEmailCode")
	defer span.End()

	user, err := s.userRepo.GetByID(ctx, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	code := random.GenerateVerificationCode(8)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.CreateVerificationCode(ctx, user.User.ID, "email", code, input.Email)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	email := domain.Email{
		Type:     "email",
		To:       input.Email,
		Username: user.User.Username,
		Code:     code,
	}

	bytes, err := json.Marshal(&email)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.producer.PublishMessage(ctx, bytes)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) SendPasswordCode(ctx context.Context, input domain.VerificationCode) error {
	ctx, span := s.tracer.Start(ctx, "authService.SendPasswordCode")
	defer span.End()

	user, err := s.userRepo.GetByEmail(ctx, input.Email)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	code := random.GenerateVerificationCode(8)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.CreateVerificationCode(ctx, user.ID, "password", code, input.Email)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	email := domain.Email{
		Type:     "password",
		To:       input.Email,
		Username: user.Username,
		Code:     code,
	}

	bytes, err := json.Marshal(&email)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.producer.PublishMessage(ctx, bytes)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) ResetEmail(ctx context.Context, userID int, input domain.ResetEmailRequest) error {
	ctx, span := s.tracer.Start(ctx, "authService.ResetPassword")
	defer span.End()

	code, err := s.repo.GetVerificationCode(ctx, input.Code, "email")

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if code.ExpireDate.Before(time.Now()) {
		return http_errors.ErrCodeExpired
	}

	err = s.repo.UpdateEmail(ctx, userID, code.Email)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.DeleteVerificationCode(ctx, code.ID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, input domain.ResetPasswordRequest) error {
	ctx, span := s.tracer.Start(ctx, "authService.ResetPassword")
	defer span.End()

	code, err := s.repo.GetVerificationCode(ctx, input.Code, "password")

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if code.ExpireDate.Before(time.Now()) {
		return http_errors.ErrCodeExpired
	}

	input.Password = authLib.GenerateHashPassword(input.Password)

	err = s.repo.UpdatePassword(ctx, code.UserID, input.Password)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.DeleteVerificationCode(ctx, code.ID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
