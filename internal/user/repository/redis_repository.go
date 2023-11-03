package repository

import (
	"context"
	"encoding/json"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type UserRedisRepo struct {
	redisClient *redis.Client
}

func NewUserRedisRepo(redisClient *redis.Client) *UserRedisRepo {
	return &UserRedisRepo{redisClient: redisClient}
}

func (repo *UserRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.User, error) {
	userBytes, err := repo.redisClient.Get(ctx, key).Bytes()

	if err != nil {
		return nil, err
	}

	var user *models.User

	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, err
	}

	return user, nil
}

//func (repo *UserRedisRepo) SetUserCtx(ctx context.Context, key string)
