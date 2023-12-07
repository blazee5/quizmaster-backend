package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/internal/user"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strconv"
)

type Service struct {
	log           *zap.SugaredLogger
	repo          quiz.Repository
	quizRedisRepo quiz.RedisRepository
	userRedisRepo user.RedisRepository
	elasticRepo   quiz.ElasticRepository
	tracer        trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo quiz.Repository, quizRedisRepo quiz.RedisRepository, userRedisRepo user.RedisRepository, elasticRepo quiz.ElasticRepository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, quizRedisRepo: quizRedisRepo, userRedisRepo: userRedisRepo, elasticRepo: elasticRepo, tracer: tracer}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Quiz, error) {
	ctx, span := s.tracer.Start(ctx, "quizService.GetAll")
	defer span.End()

	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int) (models.Quiz, error) {
	ctx, span := s.tracer.Start(ctx, "quizService.GetByID")
	defer span.End()

	cachedQuiz, err := s.quizRedisRepo.GetByIDCtx(ctx, strconv.Itoa(id))

	if err != nil {
		s.log.Infof("cannot get quiz by id in redis: %v", err)
	}

	if cachedQuiz != nil {
		return *cachedQuiz, nil
	}

	quiz, err := s.repo.GetByID(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.Quiz{}, err
	}

	if err := s.quizRedisRepo.SetQuizCtx(ctx, strconv.Itoa(quiz.ID), 600, &quiz); err != nil {
		s.log.Infof("error while save quiz to cache: %v", err)
	}

	return quiz, nil
}

func (s *Service) Create(ctx context.Context, userID int, input domain.Quiz) (int, error) {
	ctx, span := s.tracer.Start(ctx, "quizService.Create")
	defer span.End()

	id, err := s.repo.Create(ctx, userID, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	if err := s.userRedisRepo.DeleteUserCtx(ctx, strconv.Itoa(userID)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	quiz := models.Quiz{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
	}

	if err := s.elasticRepo.CreateIndex(ctx, quiz); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (s *Service) Update(ctx context.Context, userID, quizID int, input domain.Quiz) error {
	ctx, span := s.tracer.Start(ctx, "quizService.Update")
	defer span.End()

	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.Update(ctx, quizID, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID, quizID int) error {
	ctx, span := s.tracer.Start(ctx, "quizService.Delete")
	defer span.End()

	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.Delete(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) Search(ctx context.Context, title string) ([]models.QuizInfo, error) {
	quizzes, err := s.elasticRepo.SearchIndex(ctx, title)

	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (s *Service) UploadImage(ctx context.Context, userID, quizID int, filename string) error {
	ctx, span := s.tracer.Start(ctx, "quizService.UploadImage")
	defer span.End()

	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.UploadImage(ctx, quizID, filename)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) DeleteImage(ctx context.Context, userID, quizID int) error {
	ctx, span := s.tracer.Start(ctx, "quizService.DeleteImage")
	defer span.End()

	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	err = s.repo.DeleteImage(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
