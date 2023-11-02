package service

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/blazee5/testhub-backend/internal/quiz"
	"go.uber.org/zap"
	"strconv"
)

type Service struct {
	repo      quiz.Repository
	redisRepo quiz.RedisRepository
	log       *zap.SugaredLogger
}

func NewService(repo quiz.Repository, redisRepo quiz.RedisRepository, log *zap.SugaredLogger) *Service {
	return &Service{repo: repo, redisRepo: redisRepo, log: log}
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

func (s *Service) SaveResult(ctx context.Context, userId int, input domain.Result) (int, error) {
	return s.repo.SaveResult(ctx, userId, input)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
