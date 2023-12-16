package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
)

type Service struct {
	log       *zap.SugaredLogger
	repo      user.Repository
	redisRepo user.RedisRepository
	awsRepo   user.AWSRepository
}

func NewService(log *zap.SugaredLogger, repo user.Repository, redisRepo user.RedisRepository, awsRepo user.AWSRepository) *Service {
	return &Service{log: log, repo: repo, redisRepo: redisRepo, awsRepo: awsRepo}
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

func (s *Service) ChangeAvatar(ctx context.Context, userID int, fileHeader *multipart.FileHeader) error {
	avatar, err := s.repo.GetAvatarByID(ctx, userID)

	if err != nil {
		return err
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

	if avatar == "" {
		avatar = uuid.String() + filepath.Ext(fileHeader.Filename)
	}

	err = s.awsRepo.SaveFile(ctx, avatar, contentType, bytes)

	if err != nil {
		return err
	}

	err = s.repo.ChangeAvatar(ctx, userID, avatar)

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
