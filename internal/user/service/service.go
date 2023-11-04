package service

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/blazee5/testhub-backend/internal/user"
	"go.uber.org/zap"
	"strconv"
)

type Service struct {
	log       *zap.SugaredLogger
	repo      user.Repository
	redisRepo user.RedisRepository
}

func NewService(log *zap.SugaredLogger, repo user.Repository, redisRepo user.RedisRepository) *Service {
	return &Service{log: log, repo: repo, redisRepo: redisRepo}
}

func (s *Service) GetUserById(ctx context.Context, userId int) (models.User, error) {
	cachedUser, err := s.redisRepo.GetByIdCtx(ctx, strconv.Itoa(userId))

	if err != nil {
		s.log.Infof("cannot get quiz by id in redis: %v", err)
	}

	if cachedUser != nil {
		return *cachedUser, nil
	}

	user, err := s.repo.GetUserById(ctx, userId)

	if err != nil {
		return models.User{}, err
	}

	if err := s.redisRepo.SetUserCtx(ctx, strconv.Itoa(user.Id), 600, &user); err != nil {
		s.log.Infof("error while save user to cache: %v", err)
	}

	return user, nil
}
