package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/user"
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

func (s *Service) GetByID(ctx context.Context, userID int) (models.UserInfo, error) {
	cachedUser, err := s.redisRepo.GetByIDCtx(ctx, strconv.Itoa(userID))

	if err != nil {
		s.log.Infof("cannot get user by id in redis: %v", err)
	}

	if cachedUser != nil {
		return *cachedUser, nil
	}

	user, err := s.repo.GetByID(ctx, userID)

	if err != nil {
		return models.UserInfo{}, err
	}

	if err := s.redisRepo.SetUserCtx(ctx, strconv.Itoa(user.User.ID), 600, &user); err != nil {
		s.log.Infof("error while save user to cache: %v", err)
	}

	return user, nil
}

func (s *Service) GetQuizzes(ctx context.Context, userID int) ([]models.Quiz, error) {
	return s.repo.GetQuizzes(ctx, userID)
}

func (s *Service) GetResults(ctx context.Context, userID int) ([]models.Quiz, error) {
	return s.repo.GetResults(ctx, userID)
}

func (s *Service) ChangeAvatar(ctx context.Context, userID int, file string) error {
	err := s.repo.ChangeAvatar(ctx, userID, file)

	if err != nil {
		return err
	}

	if err := s.redisRepo.DeleteUserCtx(ctx, strconv.Itoa(userID)); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, userID int, input domain.UpdateUser) error {
	err := s.repo.Update(ctx, userID, input)

	if err != nil {
		return err
	}

	if err := s.redisRepo.DeleteUserCtx(ctx, strconv.Itoa(userID)); err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	return s.repo.Delete(ctx, userID)
}
