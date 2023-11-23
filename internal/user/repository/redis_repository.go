package repository

import (
	"context"
	"encoding/json"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserRedisRepo struct {
	redisClient *redis.Client
}

func NewUserRedisRepo(redisClient *redis.Client) *UserRedisRepo {
	return &UserRedisRepo{redisClient: redisClient}
}

func (repo *UserRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.UserInfo, error) {
	userBytes, err := repo.redisClient.Get(ctx, "user:"+key).Bytes()

	if err != nil {
		return nil, err
	}

	var user *models.UserInfo

	if err = json.Unmarshal(userBytes, &user); err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRedisRepo) SetUserCtx(ctx context.Context, key string, seconds int, user *models.UserInfo) error {
	userBytes, err := json.Marshal(user)

	if err != nil {
		return err
	}

	if err := repo.redisClient.Set(ctx, "user:"+key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *UserRedisRepo) DeleteUserCtx(ctx context.Context, key string) error {
	if err := repo.redisClient.Del(ctx, "user:"+key).Err(); err != nil {
		return err
	}

	return nil
}
