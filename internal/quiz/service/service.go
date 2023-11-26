package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/internal/user"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.uber.org/zap"
	"strconv"
)

type Service struct {
	log           *zap.SugaredLogger
	repo          quiz.Repository
	quizRedisRepo quiz.RedisRepository
	userRedisRepo user.RedisRepository
	elasticRepo   quiz.ElasticRepository
}

func NewService(log *zap.SugaredLogger, repo quiz.Repository, quizRedisRepo quiz.RedisRepository, userRedisRepo user.RedisRepository, elasticRepo quiz.ElasticRepository) *Service {
	return &Service{log: log, repo: repo, quizRedisRepo: quizRedisRepo, userRedisRepo: userRedisRepo, elasticRepo: elasticRepo}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Quiz, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int) (models.Quiz, error) {
	cachedQuiz, err := s.quizRedisRepo.GetByIDCtx(ctx, strconv.Itoa(id))

	if err != nil {
		s.log.Infof("cannot get quiz by id in redis: %v", err)
	}

	if cachedQuiz != nil {
		return *cachedQuiz, nil
	}

	quiz, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return models.Quiz{}, err
	}

	if err := s.quizRedisRepo.SetQuizCtx(ctx, strconv.Itoa(quiz.ID), 600, &quiz); err != nil {
		s.log.Infof("error while save quiz to cache: %v", err)
	}

	return quiz, nil
}

func (s *Service) Create(ctx context.Context, userID int, input domain.Quiz) (int, error) {
	id, err := s.repo.Create(ctx, userID, input)

	if err != nil {
		return 0, err
	}

	if err := s.userRedisRepo.DeleteUserCtx(ctx, strconv.Itoa(userID)); err != nil {
		return 0, err
	}

	if err := s.elasticRepo.CreateIndex(ctx, input); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) Update(ctx context.Context, userID, quizID int, input domain.Quiz) error {
	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.Update(ctx, quizID, input)
}

func (s *Service) Delete(ctx context.Context, userID, quizID int) error {
	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, quizID)
}

func (s *Service) UploadImage(ctx context.Context, userID, quizID int, filename string) error {
	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.UploadImage(ctx, quizID, filename)
}

func (s *Service) DeleteImage(ctx context.Context, userID, quizID int) error {
	quiz, err := s.repo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID {
		return http_errors.PermissionDenied
	}

	return s.repo.DeleteImage(ctx, quizID)
}
