package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.uber.org/zap"
	"strconv"
)

type Service struct {
	log       *zap.SugaredLogger
	repo      quiz.Repository
	redisRepo quiz.RedisRepository
}

func NewService(log *zap.SugaredLogger, repo quiz.Repository, redisRepo quiz.RedisRepository) *Service {
	return &Service{log: log, repo: repo, redisRepo: redisRepo}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Quiz, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetById(ctx context.Context, id int) (models.Quiz, error) {
	cachedQuiz, err := s.redisRepo.GetByIdCtx(ctx, strconv.Itoa(id))

	if err != nil {
		s.log.Infof("cannot get quiz by id in redis: %v", err)
	}

	if cachedQuiz != nil {
		return *cachedQuiz, nil
	}

	quiz, err := s.repo.GetById(ctx, id)

	if err != nil {
		return models.Quiz{}, err
	}

	if err := s.redisRepo.SetQuizCtx(ctx, strconv.Itoa(quiz.Id), 600, &quiz); err != nil {
		s.log.Infof("error while save quiz to cache: %v", err)
	}

	return quiz, nil
}

func (s *Service) Create(ctx context.Context, input domain.Quiz) (int, error) {
	return s.repo.Create(ctx, input)
}

func (s *Service) GetQuestionsById(ctx context.Context, id int) ([]models.Question, error) {
	return s.repo.GetQuestionsById(ctx, id, false)
}

func (s *Service) SaveResult(ctx context.Context, userId, quizId int, input domain.Result) (int, error) {
	_, err := s.repo.GetById(ctx, quizId)

	if err != nil {
		return 0, err
	}

	return s.repo.SaveResult(ctx, userId, quizId, input)
}

func (s *Service) Delete(ctx context.Context, userId, quizId int) error {
	quiz, err := s.repo.GetById(ctx, quizId)

	if err != nil {
		return err
	}

	if quiz.UserId != userId {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, quizId)
}
