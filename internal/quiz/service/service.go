package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/internal/user"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
)

type Service struct {
	log           *zap.SugaredLogger
	repo          quizRepo.Repository
	quizRedisRepo quizRepo.RedisRepository
	userRedisRepo user.RedisRepository
	elasticRepo   quizRepo.ElasticRepository
	awsRepo       quizRepo.AWSRepository
	tracer        trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo quizRepo.Repository, quizRedisRepo quizRepo.RedisRepository, userRedisRepo user.RedisRepository, elasticRepo quizRepo.ElasticRepository, awsRepo quizRepo.AWSRepository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, quizRedisRepo: quizRedisRepo, userRedisRepo: userRedisRepo, elasticRepo: elasticRepo, awsRepo: awsRepo, tracer: tracer}
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

	quiz := models.QuizInfo{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
	}

	if err = s.elasticRepo.CreateIndex(ctx, quiz); err != nil {
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

	if err = s.repo.Update(ctx, quizID, input); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if err = s.elasticRepo.UpdateIndex(ctx, quizID, models.QuizInfo{ID: quizID, Title: input.Title, Description: input.Description}); err != nil {
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

	err = s.elasticRepo.DeleteIndex(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) Search(ctx context.Context, title string) ([]models.QuizInfo, error) {
	ctx, span := s.tracer.Start(ctx, "quizService.SearchByTitle")
	defer span.End()

	quizzes, err := s.elasticRepo.SearchIndex(ctx, title)

	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (s *Service) UploadImage(ctx context.Context, userID, quizID int, fileHeader *multipart.FileHeader) error {
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

	file, err := fileHeader.Open()

	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	contentType := http.DetectContentType(bytes)

	uuid, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	fileName := quiz.Image

	if fileName == "" {
		fileName = uuid.String() + filepath.Ext(fileHeader.Filename)
	}

	err = s.awsRepo.SaveFile(ctx, fileName, contentType, bytes)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.UploadImage(ctx, quizID, fileName)

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

	err = s.awsRepo.DeleteFile(ctx, quiz.Image)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.DeleteImage(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
